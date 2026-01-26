package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/htmlsnapshot"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
)

func registerPublicRoutes(v2 fiber.Router, deps Dependencies, websiteInfoHandler *handler.WebsiteInfoHandler, htmlSnapshotSvc *htmlsnapshot.Service, navMenuHandler *handler.NavMenuHandler) {
	public := v2.Group("/public")
	public.Get("/website-info", websiteInfoHandler.PublicList)
	public.Get("/nav-menus", navMenuHandler.ListPublic)

	htmlSnapshotHandler := handler.NewHTMLSnapshotHandler(htmlSnapshotSvc)
	public.Post("/html/posts/refresh", htmlSnapshotHandler.RefreshPostsHTML)
}
