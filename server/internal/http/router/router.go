package router

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	appEvent "github.com/grtsinry43/grtblog-v2/server/internal/app/event"
	appfed "github.com/grtsinry43/grtblog-v2/server/internal/app/federation"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/federationconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/htmlsnapshot"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/sysconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/webhook"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/websiteinfo"
	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	infraevent "github.com/grtsinry43/grtblog-v2/server/internal/infra/event"
	fedinfra "github.com/grtsinry43/grtblog-v2/server/internal/infra/federation"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/jwt"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/turnstile"
	"github.com/grtsinry43/grtblog-v2/server/internal/ws"
	"github.com/redis/go-redis/v9"
)

// Dependencies collects the shared instances that handlers require.
type Dependencies struct {
	DB         *gorm.DB
	Config     config.Config
	JWTManager *jwt.Manager
	Turnstile  *turnstile.Client
	SysConfig  *sysconfig.Service
	EventBus   appEvent.Bus
	Redis      *redis.Client
}

// Register wires up all HTTP endpoints with middlewares.
func Register(app *fiber.App, deps Dependencies) {
	healthHandler := handler.NewHealthHandler(deps.Config.App)

	app.Get("/health/liveness", healthHandler.Liveness)
	app.Get("/health/readiness", healthHandler.Readiness)
	app.Static("/uploads", filepath.Join("storage", "uploads"))

	api := app.Group("/api")
	v2 := api.Group("/v2")

	sysCfgSvc := deps.SysConfig
	if sysCfgSvc == nil {
		sysCfgRepo := persistence.NewSysConfigRepository(deps.DB)
		sysCfgSvc = sysconfig.NewService(sysCfgRepo, deps.Config.Turnstile)
	}
	deps.SysConfig = sysCfgSvc
	eventBus := deps.EventBus
	if eventBus == nil {
		eventBus = infraevent.NewInMemoryBus()
	}
	wsManager := ws.NewManager(ws.Config{
		CacheSize:       3,
		RoomTTL:         30 * time.Second,
		CleanupInterval: 5 * time.Second,
	})
	ws.RegisterArticleUpdateSubscriber(eventBus, wsManager)
	ws.RegisterMomentUpdateSubscriber(eventBus, wsManager)
	ws.RegisterPageUpdateSubscriber(eventBus, wsManager)

	webhookSettings, err := sysCfgSvc.WebhookSettings(context.Background())
	if err != nil {
		log.Printf("webhook settings error: %v", err)
	}
	webhookRepo := persistence.NewWebhookRepository(deps.DB)
	webhookSender := webhook.NewSender(webhookRepo, webhookSettings.Timeout)
	webhookDispatcher := webhook.NewDispatcher(webhookRepo, webhookSender, webhookSettings.Workers, webhookSettings.QueueSize)
	webhookSvc := webhook.NewService(webhookRepo, webhookSender)
	webhook.RegisterSubscribers(eventBus, webhookDispatcher)

	contentRepo := persistence.NewContentRepository(deps.DB)
	htmlSnapshotSvc := htmlsnapshot.NewService(contentRepo, "")
	htmlsnapshot.RegisterArticleUpdateSubscriber(eventBus, htmlSnapshotSvc)

	fedCfgRepo := persistence.NewFederationConfigRepository(deps.DB)
	fedCfgSvc := federationconfig.NewService(fedCfgRepo)
	var fedCache fedinfra.Cache
	if deps.Redis != nil {
		fedCache = fedinfra.NewRedisCache(deps.Redis, deps.Config.Redis.Prefix)
	}
	fedResolver := fedinfra.NewResolver(&http.Client{Timeout: 10 * time.Second}, fedCache)
	fedOutbound := appfed.NewOutboundService(fedCfgSvc, fedResolver)
	appfed.RegisterSubscribers(eventBus, fedOutbound)

	websiteInfoRepo := persistence.NewWebsiteInfoRepository(deps.DB)
	websiteInfoSvc := websiteinfo.NewService(websiteInfoRepo)
	websiteInfoHandler := handler.NewWebsiteInfoHandler(websiteInfoSvc)

	registerPublicRoutes(v2, deps, websiteInfoHandler, htmlSnapshotSvc)
	registerAuthRoutes(v2, deps, sysCfgSvc)
	deps.EventBus = eventBus
	registerWSRoutes(v2, wsManager)
	registerArticlePublicRoutes(v2, deps)
	registerMomentPublicRoutes(v2, deps)
	registerPagePublicRoutes(v2, deps)
	registerTaxonomyPublicRoutes(v2, deps)
	registerUserRoutes(v2, deps, websiteInfoHandler)
	registerArticleAuthRoutes(v2, deps)
	registerMomentAuthRoutes(v2, deps)
	registerPageAuthRoutes(v2, deps)
	registerAdminRoutes(v2, deps, websiteInfoHandler, sysCfgSvc)
	registerTaxonomyAdminRoutes(v2, deps)
	registerWebhookAdminRoutes(v2, deps, webhookSvc)

	docsHandler := handler.NewDocsHandler("docs/swagger.json")
	app.Get("/docs/openapi.json", docsHandler.OpenAPI)
	app.Get("/docs", docsHandler.Scalar)

	registerFederationRoutes(app, deps)
}
