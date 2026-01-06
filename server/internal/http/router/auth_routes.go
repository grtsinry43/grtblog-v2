package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/auth"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/sysconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerAuthRoutes(v2 fiber.Router, deps Dependencies, sysCfgSvc *sysconfig.Service) {
	identityRepo := persistence.NewIdentityRepository(deps.DB)
	oauthRepo := persistence.NewOAuthProviderRepository(deps.DB)
	var stateStore auth.StateStore
	if deps.Redis != nil {
		stateStore = auth.NewRedisStateStore(deps.Redis, deps.Config.Redis.Prefix)
	}
	authSvc := auth.NewService(identityRepo, oauthRepo, deps.JWTManager, stateStore, deps.Config.Auth)
	authHandler := handler.NewAuthHandler(authSvc, sysCfgSvc, deps.Turnstile)
	oauthHandler := handler.NewOAuthHandler(authSvc, deps.Config.Auth.OAuthStateTTL)

	authGroup := v2.Group("/auth", limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			handler.Audit(c, "auth.rate_limited", map[string]any{"ip": c.IP()})
			return response.NewBizErrorWithMsg(response.TooManyRequests, "")
		},
	}))
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Get("/init-state", authHandler.InitState)
	authGroup.Get("/providers", oauthHandler.ListProviders)
	authGroup.Get("/providers/:provider/authorize", oauthHandler.Authorize)
	authGroup.Post("/providers/:provider/callback", oauthHandler.Callback)
}
