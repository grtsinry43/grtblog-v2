package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/sysconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/websiteinfo"
	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/jwt"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/rbac"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/turnstile"
)

// Dependencies collects the shared instances that handlers require.
type Dependencies struct {
	DB         *gorm.DB
	Config     config.Config
	JWTManager *jwt.Manager
	Enforcer   *rbac.Enforcer
	Turnstile  *turnstile.Client
	SysConfig  *sysconfig.Service
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

	websiteInfoRepo := persistence.NewWebsiteInfoRepository(deps.DB)
	websiteInfoSvc := websiteinfo.NewService(websiteInfoRepo)
	websiteInfoHandler := handler.NewWebsiteInfoHandler(websiteInfoSvc)

	registerAuthRoutes(v2, deps, sysCfgSvc)
	registerPublicRoutes(v2, websiteInfoHandler)
	registerProtectedRoutes(v2, deps, websiteInfoHandler)

	docsHandler := handler.NewDocsHandler("docs/swagger.json")
	app.Get("/docs/openapi.json", docsHandler.OpenAPI)
	app.Get("/docs", docsHandler.Scalar)
}
