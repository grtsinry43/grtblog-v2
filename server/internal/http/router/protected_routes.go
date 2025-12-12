package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/friendlink"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerProtectedRoutes(v2 fiber.Router, deps Dependencies, websiteInfoHandler *handler.WebsiteInfoHandler) {
	authenticated := v2.Group("", middleware.RequireAuth(deps.JWTManager))

	friendLinkRepo := persistence.NewFriendLinkApplicationRepository(deps.DB)
	friendLinkSvc := friendlink.NewService(friendLinkRepo)
	friendLinkHandler := handler.NewFriendLinkHandler(friendLinkSvc)
	friendLinks := authenticated.Group("/friend-links")
	friendLinks.Post("/applications", friendLinkHandler.SubmitApplication)

	websiteInfo := authenticated.Group("/website-info")
	websiteInfo.Get("", middleware.RequirePermission(deps.Enforcer, "config:read"), websiteInfoHandler.List)
	websiteInfo.Post("", middleware.RequirePermission(deps.Enforcer, "config:write"), websiteInfoHandler.Create)
	websiteInfo.Put("/:key", middleware.RequirePermission(deps.Enforcer, "config:write"), websiteInfoHandler.Update)
	websiteInfo.Delete("/:key", middleware.RequirePermission(deps.Enforcer, "config:write"), websiteInfoHandler.Delete)
}
