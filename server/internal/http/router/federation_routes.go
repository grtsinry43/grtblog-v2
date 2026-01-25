package router

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/federationconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/federation"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerFederationRoutes(app *fiber.App, deps Dependencies) {
	cfgRepo := persistence.NewFederationConfigRepository(deps.DB)
	cfgSvc := federationconfig.NewService(cfgRepo)
	instanceRepo := persistence.NewFederationInstanceRepository(deps.DB)
	linkRepo := persistence.NewFriendLinkRepository(deps.DB)
	appRepo := persistence.NewFriendLinkApplicationRepository(deps.DB)
	contentRepo := persistence.NewContentRepository(deps.DB)
	userRepo := persistence.NewIdentityRepository(deps.DB)
	citationRepo := persistence.NewFederatedCitationRepository(deps.DB)
	mentionRepo := persistence.NewFederatedMentionRepository(deps.DB)
	postCacheRepo := persistence.NewFederatedPostCacheRepository(deps.DB)

	var cache federation.Cache
	if deps.Redis != nil {
		cache = federation.NewRedisCache(deps.Redis, deps.Config.Redis.Prefix)
	}
	resolver := federation.NewResolver(&http.Client{Timeout: 10 * time.Second}, cache)
	verifier := federation.NewVerifier(resolver, 5*time.Minute)

	wellKnownHandler := handler.NewFederationWellKnownHandler(cfgSvc, deps.Config.App)
	app.Get("/.well-known/blog-federation/manifest.json", wellKnownHandler.Manifest)
	app.Get("/.well-known/blog-federation/public-key.json", wellKnownHandler.PublicKey)
	app.Get("/.well-known/blog-federation/endpoints.json", wellKnownHandler.Endpoints)

	federationGroup := app.Group("/api/federation")
	friendLinkHandler := handler.NewFederationFriendLinkHandler(cfgSvc, instanceRepo, linkRepo, appRepo, resolver, verifier)
	federationGroup.Post("/friendlinks/request", friendLinkHandler.RequestFriendLink)

	timelineHandler := handler.NewFederationTimelineHandler(contentRepo, userRepo, cfgSvc)
	federationGroup.Get("/timeline/posts", timelineHandler.ListTimelinePosts)

	postHandler := handler.NewFederationPostHandler(contentRepo, userRepo, postCacheRepo, cfgSvc)
	federationGroup.Get("/posts/:id", postHandler.GetPostDetail)

	citationHandler := handler.NewFederationCitationHandler(cfgSvc, contentRepo, instanceRepo, citationRepo, linkRepo, resolver, verifier)
	federationGroup.Post("/citations/request", citationHandler.RequestCitation)

	mentionHandler := handler.NewFederationMentionHandler(cfgSvc, instanceRepo, mentionRepo, userRepo, resolver, verifier)
	federationGroup.Post("/mentions/notify", mentionHandler.NotifyMention)
}
