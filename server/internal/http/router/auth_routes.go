package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/auth"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/sysconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerAuthRoutes(v2 fiber.Router, deps Dependencies, sysCfgSvc *sysconfig.Service) {
	identityRepo := persistence.NewIdentityRepository(deps.DB)
	authSvc := auth.NewService(identityRepo, deps.JWTManager, deps.Config.Auth)
	authHandler := handler.NewAuthHandler(authSvc, sysCfgSvc, deps.Turnstile)

	authGroup := v2.Group("/auth", limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}))
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
}
