package router

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	appfed "github.com/grtsinry43/grtblog-v2/server/internal/app/federation"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/federationconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/sysconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	fedinfra "github.com/grtsinry43/grtblog-v2/server/internal/infra/federation"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerAdminRoutes(v2 fiber.Router, deps Dependencies, websiteInfoHandler *handler.WebsiteInfoHandler, navMenuHandler *handler.NavMenuHandler, sysCfgSvc *sysconfig.Service) {
	adminGroup := v2.Group("", middleware.RequireAuth(deps.JWTManager), middleware.RequireAdmin())

	websiteInfo := adminGroup.Group("/website-info")
	websiteInfo.Get("", websiteInfoHandler.List)
	websiteInfo.Post("", websiteInfoHandler.Create)
	websiteInfo.Put("/:key", websiteInfoHandler.Update)
	websiteInfo.Delete("/:key", websiteInfoHandler.Delete)

	navMenus := adminGroup.Group("/admin/nav-menus")
	navMenus.Get("", navMenuHandler.ListAdmin)
	navMenus.Post("", navMenuHandler.Create)
	navMenus.Put("/reorder", navMenuHandler.Reorder)
	navMenus.Put("/:id", navMenuHandler.Update)
	navMenus.Delete("/:id", navMenuHandler.Delete)

	oauthRepo := persistence.NewOAuthProviderRepository(deps.DB)
	adminOAuth := handler.NewAdminOAuthHandler(oauthRepo)
	admin := adminGroup.Group("/admin")
	admin.Get("/oauth-providers", adminOAuth.List)
	admin.Post("/oauth-providers", adminOAuth.Create)
	admin.Put("/oauth-providers/:key", adminOAuth.Update)
	admin.Delete("/oauth-providers/:key", adminOAuth.Delete)

	if sysCfgSvc != nil {
		sysConfigHandler := handler.NewSysConfigHandler(sysCfgSvc)
		admin.Get("/sysconfig", sysConfigHandler.ListSysConfig)
		admin.Put("/sysconfig", sysConfigHandler.UpdateSysConfig)
	}

	fedCfgRepo := persistence.NewFederationConfigRepository(deps.DB)
	fedCfgSvc := federationconfig.NewService(fedCfgRepo)
	fedCfgHandler := handler.NewFederationConfigHandler(fedCfgSvc)
	admin.Get("/federation/config", fedCfgHandler.ListFederationConfig)
	admin.Put("/federation/config", fedCfgHandler.UpdateFederationConfig)

	contentRepo := persistence.NewContentRepository(deps.DB)
	instanceRepo := persistence.NewFederationInstanceRepository(deps.DB)
	var cache fedinfra.Cache
	if deps.Redis != nil {
		cache = fedinfra.NewRedisCache(deps.Redis, deps.Config.Redis.Prefix)
	}
	resolver := fedinfra.NewResolver(&http.Client{Timeout: 10 * time.Second}, cache)
	outbound := appfed.NewOutboundService(fedCfgSvc, resolver, instanceRepo)
	federationAdminHandler := handler.NewFederationAdminHandler(fedCfgSvc, contentRepo, outbound, resolver)
	admin.Post("/federation/friendlinks/request", federationAdminHandler.RequestFriendLink)
	admin.Post("/federation/citations/request", federationAdminHandler.SendCitation)
	admin.Post("/federation/mentions/notify", federationAdminHandler.SendMention)
	admin.Get("/federation/remote/check", federationAdminHandler.CheckRemote)

	logHandler := handler.NewAdminLogHandler("storage/logs/app.log", 200)
	adminLogs := adminGroup.Group("/admin")
	adminLogs.Get("/logs", logHandler.List)
}
