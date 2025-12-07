package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/friendlink"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/websiteinfo"
	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

// Dependencies collects the shared instances that handlers require.
type Dependencies struct {
	DB     *gorm.DB
	Config config.Config
}

// Register wires up all HTTP endpoints with middlewares.
func Register(app *fiber.App, deps Dependencies) {
	healthHandler := handler.NewHealthHandler(deps.Config.App)

	app.Get("/health/liveness", healthHandler.Liveness)
	app.Get("/health/readiness", healthHandler.Readiness)

	websiteInfoRepo := persistence.NewWebsiteInfoRepository(deps.DB)
	websiteInfoSvc := websiteinfo.NewService(websiteInfoRepo)
	websiteInfoHandler := handler.NewWebsiteInfoHandler(websiteInfoSvc)

	app.Get("/website-info", websiteInfoHandler.List)
	app.Post("/website-info", websiteInfoHandler.Create)
	app.Put("/website-info/:key", websiteInfoHandler.Update)
	app.Delete("/website-info/:key", websiteInfoHandler.Delete)

	friendLinkRepo := persistence.NewFriendLinkApplicationRepository(deps.DB)
	friendLinkSvc := friendlink.NewService(friendLinkRepo)
	friendLinkHandler := handler.NewFriendLinkHandler(friendLinkSvc)

	app.Post("/friend-links/applications", friendLinkHandler.SubmitApplication)
}
