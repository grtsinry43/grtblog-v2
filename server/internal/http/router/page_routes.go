package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/page"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerPagePublicRoutes(v2 fiber.Router, deps Dependencies) {
	pageHandler := newPageHandler(deps)

	publicGroup := v2.Group("/pages")
	publicGroup.Get("/", pageHandler.ListPages)                        // GET /api/v2/pages
	publicGroup.Get("/:id", pageHandler.GetPage)                       // GET /api/v2/pages/123
	publicGroup.Get("/short/:shortUrl", pageHandler.GetPageByShortURL) // GET /api/v2/pages/short/abc123
	publicGroup.Post("/:id/latest", pageHandler.CheckPageLatest)       // POST /api/v2/pages/123/latest
}

func registerPageAuthRoutes(v2 fiber.Router, deps Dependencies) {
	pageHandler := newPageHandler(deps)

	authGroup := v2.Group("/pages", middleware.RequireAuth(deps.JWTManager))
	authGroup.Post("/", pageHandler.CreatePage)      // POST /api/v2/pages
	authGroup.Put("/:id", pageHandler.UpdatePage)    // PUT /api/v2/pages/123
	authGroup.Delete("/:id", pageHandler.DeletePage) // DELETE /api/v2/pages/123
}

func newPageHandler(deps Dependencies) *handler.PageHandler {
	contentRepo := persistence.NewContentRepository(deps.DB)
	pageSvc := page.NewService(contentRepo, deps.EventBus)
	return handler.NewPageHandler(pageSvc)
}
