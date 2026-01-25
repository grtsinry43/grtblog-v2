package router

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/comment"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/clientinfo"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/geoip"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerCommentPublicRoutes(v2 fiber.Router, deps Dependencies) {
	commentHandler := newCommentHandler(deps)

	publicGroup := v2.Group("/comments")
	publicGroup.Get("/areas/:areaId", commentHandler.ListCommentTree)
	publicGroup.Post("/areas/:areaId/visitor", commentHandler.CreateCommentVisitor)
}

func registerCommentAuthRoutes(v2 fiber.Router, deps Dependencies) {
	commentHandler := newCommentHandler(deps)

	authGroup := v2.Group("/comments", middleware.RequireAuth(deps.JWTManager))
	authGroup.Post("/areas/:areaId", commentHandler.CreateCommentLogin)
}

func newCommentHandler(deps Dependencies) *handler.CommentHandler {
	commentRepo := persistence.NewCommentRepository(deps.DB)
	identityRepo := persistence.NewIdentityRepository(deps.DB)
	clientInfoResolver := clientinfo.NewUAParser()

	var geoResolver comment.GeoIPResolver
	if deps.Config.GeoIP.DBPath != "" {
		geoip.EnsureDatabasesAsync(
			context.Background(),
			deps.Config.GeoIP.DBPath,
			deps.Config.GeoIP.DownloadURL,
			deps.Config.GeoIP.ASNPath,
			deps.Config.GeoIP.ASNURL,
			log.Printf,
		)
		resolver, err := geoip.NewResolver(deps.Config.GeoIP.DBPath, deps.Config.GeoIP.ASNPath)
		if err != nil {
			log.Printf("geoip resolver init failed: %v", err)
			geoResolver = geoip.NewLazyResolver(deps.Config.GeoIP.DBPath, deps.Config.GeoIP.ASNPath)
		} else {
			geoResolver = resolver
		}
	}

	commentSvc := comment.NewService(commentRepo, identityRepo, clientInfoResolver, geoResolver)
	return handler.NewCommentHandler(commentSvc)
}
