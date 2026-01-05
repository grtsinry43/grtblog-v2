package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/auth"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/friendlink"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerProtectedRoutes(v2 fiber.Router, deps Dependencies, websiteInfoHandler *handler.WebsiteInfoHandler) {
	authenticated := v2.Group("", middleware.RequireAuth(deps.JWTManager))

	// 当前登录人角色/权限
	identityRepo := persistence.NewIdentityRepository(deps.DB)
	oauthRepo := persistence.NewOAuthProviderRepository(deps.DB)
	authSvc := auth.NewService(identityRepo, oauthRepo, deps.JWTManager, nil, deps.Config.Auth)
	authHandler := handler.NewAuthHandler(authSvc, nil, nil)
	authenticated.Get("/auth/access-info", authHandler.AccessInfo)
	authenticated.Put("/auth/profile", authHandler.UpdateProfile)
	authenticated.Put("/auth/password", authHandler.ChangePassword)
	authenticated.Get("/auth/oauth-bindings", authHandler.ListOAuthBindings)

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

	// Admin OAuth provider 管理
	adminOAuth := handler.NewAdminOAuthHandler(oauthRepo)
	admin := authenticated.Group("/admin")
	admin.Get("/oauth-providers", middleware.RequirePermission(deps.Enforcer, "oauth:read"), adminOAuth.List)
	admin.Post("/oauth-providers", adminOAuth.Create)
	admin.Put("/oauth-providers/:key", adminOAuth.Update)
	admin.Delete("/oauth-providers/:key", adminOAuth.Delete)

	logHandler := handler.NewAdminLogHandler("storage/logs/app.log", 200)
	adminLogs := authenticated.Group("/admin", middleware.RequirePermission(deps.Enforcer, "log:read"))
	adminLogs.Get("/logs", logHandler.List)
}
