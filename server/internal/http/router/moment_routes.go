package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/moment"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerMomentPublicRoutes(v2 fiber.Router, deps Dependencies) {
	momentHandler := newMomentHandler(deps)

	publicGroup := v2.Group("/moments")
	publicGroup.Get("/", momentHandler.ListMoments)                        // GET /api/v2/moments
	publicGroup.Get("/:id", momentHandler.GetMoment)                       // GET /api/v2/moments/123
	publicGroup.Get("/short/:shortUrl", momentHandler.GetMomentByShortURL) // GET /api/v2/moments/short/abc123
	publicGroup.Post("/:id/latest", momentHandler.CheckMomentLatest)       // POST /api/v2/moments/123/latest
}

func registerMomentAuthRoutes(v2 fiber.Router, deps Dependencies) {
	momentHandler := newMomentHandler(deps)

	authGroup := v2.Group("/moments", middleware.RequireAuth(deps.JWTManager))
	authGroup.Post("/", momentHandler.CreateMoment)      // POST /api/v2/moments
	authGroup.Put("/:id", momentHandler.UpdateMoment)    // PUT /api/v2/moments/123
	authGroup.Delete("/:id", momentHandler.DeleteMoment) // DELETE /api/v2/moments/123
}

func newMomentHandler(deps Dependencies) *handler.MomentHandler {
	contentRepo := persistence.NewContentRepository(deps.DB)
	identityRepo := persistence.NewIdentityRepository(deps.DB)
	momentSvc := moment.NewService(contentRepo, deps.EventBus)
	return handler.NewMomentHandler(momentSvc, contentRepo, identityRepo)
}
