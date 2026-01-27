package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/thinking"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerThinkingPublicRoutes(v2 fiber.Router, deps Dependencies) {
	thinkingHandler := newThinkingHandler(deps)

	publicGroup := v2.Group("/thinkings")
	publicGroup.Get("/", thinkingHandler.ListThinkings)
}

func registerThinkingAuthRoutes(v2 fiber.Router, deps Dependencies) {
	thinkingHandler := newThinkingHandler(deps)
	authGroup := v2.Group("/thinkings", middleware.RequireAuth(deps.JWTManager), middleware.RequireAdmin())
	authGroup.Post("/", thinkingHandler.CreateThinking)
	authGroup.Put("/", thinkingHandler.UpdateThinking)
	authGroup.Delete("/:id", thinkingHandler.DeleteThinking)
}

func newThinkingHandler(deps Dependencies) *handler.ThinkingHandler {
	thinkingRepo := persistence.NewThinkingRepository(deps.DB)
	commentRepo := persistence.NewCommentRepository(deps.DB)
	thinkingSvc := thinking.NewService(thinkingRepo, commentRepo)
	return handler.NewThinkingHandler(thinkingSvc)
}
