package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	appEvent "github.com/grtsinry43/grtblog-v2/server/internal/app/event"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/sysconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/websiteinfo"
	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	infraevent "github.com/grtsinry43/grtblog-v2/server/internal/infra/event"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/jwt"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/turnstile"
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

	api := app.Group("/api")
	v2 := api.Group("/v2")

	sysCfgSvc := deps.SysConfig
	if sysCfgSvc == nil {
		sysCfgRepo := persistence.NewSysConfigRepository(deps.DB)
		sysCfgSvc = sysconfig.NewService(sysCfgRepo, deps.Config.Turnstile)
	}
	eventBus := deps.EventBus
	if eventBus == nil {
		eventBus = infraevent.NewInMemoryBus()
	}

	websiteInfoRepo := persistence.NewWebsiteInfoRepository(deps.DB)
	websiteInfoSvc := websiteinfo.NewService(websiteInfoRepo)
	websiteInfoHandler := handler.NewWebsiteInfoHandler(websiteInfoSvc)

	registerPublicRoutes(v2, deps, websiteInfoHandler)
	registerAuthRoutes(v2, deps, sysCfgSvc)
	deps.EventBus = eventBus
	registerArticlePublicRoutes(v2, deps)
	registerTaxonomyPublicRoutes(v2, deps)
	registerUserRoutes(v2, deps, websiteInfoHandler)
	registerArticleAuthRoutes(v2, deps)
	registerAdminRoutes(v2, deps, websiteInfoHandler)
	registerTaxonomyAdminRoutes(v2, deps)

	docsHandler := handler.NewDocsHandler("docs/swagger.json")
	app.Get("/docs/openapi.json", docsHandler.OpenAPI)
	app.Get("/docs", docsHandler.Scalar)
}
