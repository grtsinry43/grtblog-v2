package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/adminnotification"
	appai "github.com/grtsinry43/grtblog-v2/server/internal/app/ai"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/analytics"
	backupapp "github.com/grtsinry43/grtblog-v2/server/internal/app/backup"
	appcomment "github.com/grtsinry43/grtblog-v2/server/internal/app/comment"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/email"
	appEvent "github.com/grtsinry43/grtblog-v2/server/internal/app/event"
	appfed "github.com/grtsinry43/grtblog-v2/server/internal/app/federation"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/friendlink"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/health"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/htmlsnapshot"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/isr"
	mediaapp "github.com/grtsinry43/grtblog-v2/server/internal/app/media"
	appnav "github.com/grtsinry43/grtblog-v2/server/internal/app/navigation"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/observability"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/ownerstatus"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/sysconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/telemetry"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/webhook"
	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	infraevent "github.com/grtsinry43/grtblog-v2/server/internal/infra/event"
	fedinfra "github.com/grtsinry43/grtblog-v2/server/internal/infra/federation"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/metrics"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/jwt"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/turnstile"
	"github.com/grtsinry43/grtblog-v2/server/internal/ws"
	"github.com/redis/go-redis/v9"
)

// Dependencies collects the shared instances that handlers require.
type Dependencies struct {
	DB                   *gorm.DB
	Config               config.Config
	JWTManager           *jwt.Manager
	Turnstile            *turnstile.Client
	SysConfig            *sysconfig.Service
	EventBus             appEvent.Bus
	Redis                *redis.Client
	Analytics            *analytics.Service
	HTTPStats            *metrics.HTTPStats
	Observability        *observability.Service
	HTMLSnapshot         *htmlsnapshot.Service
	ISR                  *isr.Service
	OwnerStatus          *ownerstatus.Service
	HealthState          *health.State
	HealthChecker        *health.Checker
	FedSync              *appfed.SyncWorker
	Telemetry            *telemetry.Service
	FederationHTTPClient *http.Client
	Backup               *backupapp.Service
	MediaGate            *mediaapp.MutationGate
}

// Register wires up all HTTP endpoints with middlewares.
func Register(app *fiber.App, deps Dependencies) {
	// Health state machine.
	isDev := strings.ToLower(deps.Config.App.Env) == "development"
	if deps.HealthState == nil {
		deps.HealthState = health.NewState(isDev)
	}

	healthHandler := handler.NewHealthHandler(deps.Config.App, deps.DB, deps.Redis, deps.HealthState)

	app.Get("/health/liveness", healthHandler.Liveness)
	app.Get("/health/readiness", healthHandler.Readiness)
	app.Static("/uploads", filepath.Join("storage", "uploads"))

	api := app.Group("/api")
	v2 := api.Group("/v2")

	eventBus := deps.EventBus
	if eventBus == nil {
		eventBus = infraevent.NewInMemoryBus()
	}
	sysCfgSvc := deps.SysConfig
	if sysCfgSvc == nil {
		sysCfgRepo := persistence.NewSysConfigRepository(deps.DB)
		sysCfgSvc = sysconfig.NewService(sysCfgRepo, deps.Config.Turnstile, eventBus)
	}
	deps.SysConfig = sysCfgSvc
	wsManager := ws.NewManager(ws.Config{
		CacheSize:       3,
		RoomTTL:         30 * time.Second,
		CleanupInterval: 5 * time.Second,
		MessageTTL:      60 * time.Second,
	})
	ws.RegisterArticleUpdateSubscriber(eventBus, wsManager)
	ws.RegisterMomentUpdateSubscriber(eventBus, wsManager)
	ws.RegisterPageUpdateSubscriber(eventBus, wsManager)
	ws.RegisterNotificationSubscriber(eventBus, wsManager)
	ws.RegisterGlobalNotificationSubscriber(eventBus, wsManager)
	ws.RegisterHealthSubscriber(eventBus, wsManager)

	// Late-inject wsManager into telemetry (created before router).
	if deps.Telemetry != nil {
		deps.Telemetry.SetWSManager(wsManager)
	}

	// Create health checker (will be started by server.Start).
	if deps.HealthChecker == nil {
		deps.HealthChecker = health.NewChecker(deps.HealthState, deps.DB, deps.Redis, sysCfgSvc, eventBus, 0, deps.Config.App.HTMLSnapshotBaseURL)
	}

	webhookSettings, err := sysCfgSvc.WebhookSettings(context.Background())
	if err != nil {
		log.Printf("webhook settings error: %v", err)
	}
	webhookRepo := persistence.NewWebhookRepository(deps.DB)
	webhookSender := webhook.NewSender(webhookRepo, webhookSettings.Timeout, sysCfgSvc)
	webhookDispatcher := webhook.NewDispatcher(webhookRepo, webhookSender, webhookSettings.Workers, webhookSettings.QueueSize)
	webhookSvc := webhook.NewService(webhookRepo, webhookSender)
	webhook.RegisterSubscribers(eventBus, webhookDispatcher)

	emailSettings, err := sysCfgSvc.EmailSettings(context.Background())
	if err != nil {
		log.Printf("email settings error: %v", err)
	}
	emailRepo := persistence.NewEmailRepository(deps.DB)
	emailSender := email.NewSender(sysCfgSvc)
	emailDispatcher := email.NewDispatcher(
		emailRepo,
		emailSender,
		sysCfgSvc,
		emailSettings.Workers,
		emailSettings.QueueSize,
		emailSettings.MaxRetries,
		2*time.Second,
	)
	email.RegisterSubscribers(eventBus, emailDispatcher)

	contentRepo := persistence.NewContentRepository(deps.DB)
	albumRepo := persistence.NewAlbumRepository(deps.DB)
	thinkingRepo := persistence.NewThinkingRepository(deps.DB)
	ws.RegisterSiteActivitySubscriber(
		eventBus,
		wsManager,
		contentRepo,
		thinkingRepo,
		persistence.NewCommentRepository(deps.DB),
		albumRepo,
	)
	htmlSnapshotSvc := deps.HTMLSnapshot
	if htmlSnapshotSvc == nil {
		htmlSnapshotSvc = htmlsnapshot.NewService(contentRepo, deps.Config.App.HTMLSnapshotBaseURL, deps.Redis, deps.Config.Redis.Prefix)
	}
	deps.HTMLSnapshot = htmlSnapshotSvc
	isrSvc := deps.ISR
	if isrSvc == nil {
		isrSvc = isr.NewService(deps.Redis, deps.Config.Redis.Prefix, htmlSnapshotSvc, contentRepo, albumRepo, thinkingRepo, sysCfgSvc)
	}
	deps.ISR = isrSvc
	isr.RegisterArticleSubscribers(eventBus, isrSvc)
	isr.RegisterMomentSubscribers(eventBus, isrSvc)
	isr.RegisterPageSubscribers(eventBus, isrSvc)
	isr.RegisterThinkingSubscribers(eventBus, isrSvc)
	isr.RegisterAlbumSubscribers(eventBus, isrSvc)
	isr.RegisterFriendLinkSubscribers(eventBus, isrSvc)
	isr.RegisterFriendTimelineSubscribers(eventBus, isrSvc)
	isr.RegisterLayoutSubscribers(eventBus, isrSvc)
	isr.RegisterTagContentCacheSubscribers(eventBus, deps.Redis, deps.Config.Redis.Prefix)
	deps.Observability = observability.NewService(deps.DB, deps.Redis, deps.Config.Redis.Prefix, eventBus, deps.HTTPStats, wsManager, htmlSnapshotSvc, isrSvc)
	ownerStatusSvc := deps.OwnerStatus
	if ownerStatusSvc == nil {
		ownerStatusSvc = ownerstatus.NewService(wsManager)
	}
	deps.OwnerStatus = ownerStatusSvc

	fedInstanceRepo := persistence.NewFederationInstanceRepository(deps.DB)
	fedOutboundRepo := persistence.NewOutboundDeliveryRepository(deps.DB)
	var fedCache fedinfra.Cache
	if deps.Redis != nil {
		fedCache = fedinfra.NewRedisCache(deps.Redis, deps.Config.Redis.Prefix)
	}
	fedHTTPClient := federationHTTPClient(deps)
	fedResolver := fedinfra.NewResolver(fedHTTPClient, fedCache)
	fedOutbound := appfed.NewOutboundService(sysCfgSvc, fedResolver, fedInstanceRepo, fedHTTPClient)
	fedDelivery := appfed.NewDeliveryService(
		fedOutboundRepo,
		fedOutbound,
		persistence.NewFriendLinkRepository(deps.DB),
		eventBus,
	)
	appfed.RegisterSubscribers(eventBus, fedDelivery)
	friendlink.RegisterFederationSubscribers(eventBus, fedInstanceRepo, persistence.NewFriendLinkRepository(deps.DB), fedResolver, deps.FedSync)
	adminNotifRepo := persistence.NewAdminNotificationRepository(deps.DB)
	adminNotifSvc := adminnotification.NewService(adminNotifRepo, eventBus)
	adminnotification.RegisterSubscribers(eventBus, adminNotifSvc, contentRepo, persistence.NewIdentityRepository(deps.DB))

	// AI event-driven moderation
	aiRepo := persistence.NewAIRepository(deps.DB)
	aiSvc := appai.NewService(aiRepo, sysCfgSvc)
	commentSvcForAI := appcomment.NewService(
		persistence.NewCommentRepository(deps.DB),
		persistence.NewIdentityRepository(deps.DB),
		persistence.NewFriendLinkRepository(deps.DB),
		sysCfgSvc, nil, nil, eventBus,
	)
	appai.RegisterSubscribers(eventBus, aiSvc, commentSvcForAI)

	websiteInfoHandler := handler.NewWebsiteInfoHandler(sysCfgSvc)

	navMenuRepo := persistence.NewNavMenuRepository(deps.DB)
	navMenuSvc := appnav.NewService(navMenuRepo, eventBus)
	navMenuHandler := handler.NewNavMenuHandler(navMenuSvc)

	analyticsSvc := deps.Analytics
	if analyticsSvc == nil {
		analyticsSvc = analytics.NewService(deps.Config, deps.DB, deps.Redis)
	}
	deps.Analytics = analyticsSvc

	registerPublicRoutes(v2, deps, websiteInfoHandler, htmlSnapshotSvc, navMenuHandler)
	registerEmailPublicRoutes(v2, deps, sysCfgSvc)
	registerAuthRoutes(v2, deps, sysCfgSvc)
	deps.EventBus = eventBus
	registerWSRoutes(v2, wsManager, deps)
	registerArticlePublicRoutes(v2, deps)
	registerMomentPublicRoutes(v2, deps)
	registerAlbumPublicRoutes(v2, deps)
	registerThinkingPublicRoutes(v2, deps)
	registerPagePublicRoutes(v2, deps)
	registerTaxonomyPublicRoutes(v2, deps)
	registerCommentPublicRoutes(v2, deps)
	registerUserRoutes(v2, deps, websiteInfoHandler)
	registerArticleAuthRoutes(v2, deps)
	registerMomentAuthRoutes(v2, deps)
	registerAlbumAuthRoutes(v2, deps)
	registerThinkingAuthRoutes(v2, deps)
	registerPageAuthRoutes(v2, deps)
	registerCommentAuthRoutes(v2, deps)
	registerAdminRoutes(v2, deps, websiteInfoHandler, navMenuHandler, sysCfgSvc, wsManager, aiSvc)
	registerBackupRoutes(v2, deps)
	registerContentExportRoutes(v2, deps)
	registerTaxonomyAdminRoutes(v2, deps)
	registerWebhookAdminRoutes(v2, deps, webhookSvc)

	if isDev {
		docsHandler := handler.NewDocsHandler("docs/swagger.json")
		app.Get("/docs/openapi.json", docsHandler.OpenAPI)
		app.Get("/docs", docsHandler.Scalar)
	}

	registerFederationRoutes(app, deps)
	registerInternalRoutes(app, deps)
	registerAdminSPA(app)
}

func federationHTTPClient(deps Dependencies) *http.Client {
	if deps.FederationHTTPClient != nil {
		return deps.FederationHTTPClient
	}
	return fedinfra.NewSafeHTTPClient(10 * time.Second)
}

// registerAdminSPA serves the admin Vue SPA with client-side routing support.
func registerAdminSPA(app *fiber.App) {
	const dir = "admin"

	app.Get("/admin/*", func(c *fiber.Ctx) error {
		sub := c.Params("*")
		if sub != "" {
			// Normalize to a relative path under admin/ to avoid path escape.
			clean := strings.TrimPrefix(path.Clean("/"+sub), "/")
			fp := filepath.Join(dir, clean)
			if info, err := os.Stat(fp); err == nil && !info.IsDir() {
				return c.SendFile(fp)
			}
		}
		// SPA fallback: serve index.html for client-side routing
		return c.SendFile(filepath.Join(dir, "index.html"))
	})
}
