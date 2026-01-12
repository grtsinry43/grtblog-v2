package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerPublicRoutes(v2 fiber.Router, deps Dependencies, websiteInfoHandler *handler.WebsiteInfoHandler) {
	public := v2.Group("/public")
	public.Get("/website-info", websiteInfoHandler.PublicList)

	htmlSnapshotHandler := handler.NewHTMLSnapshotHandler(persistence.NewContentRepository(deps.DB))
	public.Post("/html/posts/refresh", htmlSnapshotHandler.RefreshPostsHTML)
}
