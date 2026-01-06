package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerAdminRoutes(v2 fiber.Router, deps Dependencies, websiteInfoHandler *handler.WebsiteInfoHandler) {
	adminGroup := v2.Group("", middleware.RequireAuth(deps.JWTManager), middleware.RequireAdmin())

	websiteInfo := adminGroup.Group("/website-info")
	websiteInfo.Get("", websiteInfoHandler.List)
	websiteInfo.Post("", websiteInfoHandler.Create)
	websiteInfo.Put("/:key", websiteInfoHandler.Update)
	websiteInfo.Delete("/:key", websiteInfoHandler.Delete)

	oauthRepo := persistence.NewOAuthProviderRepository(deps.DB)
	adminOAuth := handler.NewAdminOAuthHandler(oauthRepo)
	admin := adminGroup.Group("/admin")
	admin.Get("/oauth-providers", adminOAuth.List)
	admin.Post("/oauth-providers", adminOAuth.Create)
	admin.Put("/oauth-providers/:key", adminOAuth.Update)
	admin.Delete("/oauth-providers/:key", adminOAuth.Delete)

	logHandler := handler.NewAdminLogHandler("storage/logs/app.log", 200)
	adminLogs := adminGroup.Group("/admin")
	adminLogs.Get("/logs", logHandler.List)
}
