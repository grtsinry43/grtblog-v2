package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	appdiscovery "github.com/grtsinry43/grtblog-v2/server/internal/app/discovery"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/friendlink"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/friendtimeline"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/globalnotification"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/home"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/htmlsnapshot"
	applike "github.com/grtsinry43/grtblog-v2/server/internal/app/like"
	apprss "github.com/grtsinry43/grtblog-v2/server/internal/app/rss"
	appsearch "github.com/grtsinry43/grtblog-v2/server/internal/app/search"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerPublicRoutes(v2 fiber.Router, deps Dependencies, websiteInfoHandler *handler.WebsiteInfoHandler, htmlSnapshotSvc *htmlsnapshot.Service, navMenuHandler *handler.NavMenuHandler) {
	public := v2.Group("/public")
	discoveryHandler := handler.NewDiscoveryHandler(appdiscovery.NewService(persistence.NewDiscoveryRepository(deps.DB), deps.SysConfig))
	public.Get("/discovery/catalog", discoveryHandler.Catalog)
	public.Get("/discovery/resource/*", discoveryHandler.Resource)
	ownerStatusHandler := handler.NewOwnerStatusHandler(deps.OwnerStatus)
	public.Get("/owner-status", ownerStatusHandler.GetStatus)
	v2.Get("/onlineStatus", ownerStatusHandler.GetStatus)

	public.Get("/website-info", websiteInfoHandler.PublicList)
	public.Get("/nav-menus", navMenuHandler.ListPublic)
	public.Get("/tags", newTaxonomyHandler(deps).ListPublicTags)

	articleHandler := newArticleHandler(deps)
	public.Get("/articles/recent", articleHandler.ListRecentPublicArticles)

	momentHandler := newMomentHandler(deps)
	public.Get("/moments/recent", momentHandler.ListRecentPublicMoments)

	homeSvc := home.NewService(deps.DB, deps.Redis, deps.Config.Redis.Prefix, deps.SysConfig)
	homeHandler := handler.NewHomeHandler(homeSvc)
	public.Get("/home/activity-pulse", homeHandler.GetActivityPulse)
	public.Get("/home/inspiration-stats", homeHandler.GetInspirationStats)
	public.Get("/home/timeline-by-year", homeHandler.GetTimelineByYear)

	friendLinkRepo := persistence.NewFriendLinkRepository(deps.DB)
	friendLinkSvc := friendlink.NewLinkService(friendLinkRepo)
	friendLinkHandler := handler.NewFriendLinkPublicHandler(friendLinkSvc, deps.SysConfig)
	public.Get("/friend-links", friendLinkHandler.ListPublic)
	public.Get("/friend-links/apply-config", friendLinkHandler.GetApplyConfig)
	friendTimelineSvc := friendtimeline.NewService(persistence.NewFederatedPostCacheRepository(deps.DB), deps.Redis, deps.Config.Redis.Prefix)
	friendTimelineHandler := handler.NewFriendTimelineHandler(friendTimelineSvc)
	public.Get("/friend-timeline", friendTimelineHandler.ListPublic)

	globalNotificationRepo := persistence.NewGlobalNotificationRepository(deps.DB)
	globalNotificationSvc := globalnotification.NewService(globalNotificationRepo, deps.EventBus)
	globalNotificationHandler := handler.NewGlobalNotificationHandler(globalNotificationSvc)
	public.Get("/global-notifications", globalNotificationHandler.ListPublicActive)

	searchRepo := persistence.NewSearchRepository(deps.DB)
	searchSvc := appsearch.NewService(searchRepo, deps.Redis, deps.Config.Redis.Prefix)
	searchHandler := handler.NewSearchHandler(searchSvc)
	searchGuard := public.Group("/search", newRateLimiter(deps, limiter.Config{
		Max:        30,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			handler.Audit(c, "search.rate_limited", map[string]any{"ip": c.IP()})
			return response.NewBizErrorWithMsg(response.TooManyRequests, "")
		},
	}))
	searchGuard.Get("/", searchHandler.SiteSearch)

	rssSvc := apprss.NewService(
		persistence.NewContentRepository(deps.DB),
		persistence.NewThinkingRepository(deps.DB),
		deps.SysConfig,
		persistence.NewIdentityRepository(deps.DB),
	)
	rssAccessSvc := newRSSAccessAnalyticsService(deps)
	rssHandler := handler.NewRSSHandler(rssSvc, rssAccessSvc)
	public.Get("/rss.xml", rssHandler.GetFeed)
	public.Get("/feed", rssHandler.GetFeed)

	if deps.Analytics != nil {
		analyticsHandler := handler.NewAnalyticsHandler(deps.Analytics)
		viewGuard := public.Group("/analytics", newRateLimiter(deps, limiter.Config{
			Max:        60,
			Expiration: time.Minute,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				handler.Audit(c, "analytics.view.rate_limited", map[string]any{"ip": c.IP()})
				return response.NewBizErrorWithMsg(response.TooManyRequests, "")
			},
		}))
		viewGuard.Post("/view", analyticsHandler.TrackView)
	}
	likeRepo := persistence.NewLikeRepository(deps.DB)
	likeSvc := applike.NewService(likeRepo)
	likeHandler := handler.NewLikeHandler(likeSvc)
	likeGuard := public.Group("/analytics", newRateLimiter(deps, limiter.Config{
		Max:        30,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			handler.Audit(c, "analytics.like.rate_limited", map[string]any{"ip": c.IP()})
			return response.NewBizErrorWithMsg(response.TooManyRequests, "")
		},
	}))
	likeGuard.Post("/like", likeHandler.TrackLike)
}
