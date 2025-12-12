package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
)

func registerPublicRoutes(v2 fiber.Router, websiteInfoHandler *handler.WebsiteInfoHandler) {
	public := v2.Group("/public")
	public.Get("/website-info", websiteInfoHandler.PublicList)
}
