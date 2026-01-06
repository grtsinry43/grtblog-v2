package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/auth"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/friendlink"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerUserRoutes(v2 fiber.Router, deps Dependencies, websiteInfoHandler *handler.WebsiteInfoHandler) {
	authenticated := v2.Group("", middleware.RequireAuth(deps.JWTManager))

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
}
