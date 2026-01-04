package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/article"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerArticleRoutes(v2 fiber.Router, deps Dependencies) {
	// 创建repository和service
	contentRepo := persistence.NewContentRepository(deps.DB)
	articleSvc := article.NewService(contentRepo)
	articleHandler := handler.NewArticleHandler(articleSvc)

	// 公开路由（无需认证）
	publicGroup := v2.Group("/articles")
	publicGroup.Get("/", articleHandler.ListArticles)                        // GET /api/v2/articles
	publicGroup.Get("/:id", articleHandler.GetArticle)                       // GET /api/v2/articles/123
	publicGroup.Get("/short/:shortUrl", articleHandler.GetArticleByShortURL) // GET /api/v2/articles/short/abc123

	// 需要认证的路由
	authGroup := v2.Group("/articles", middleware.RequireAuth(deps.JWTManager))
	authGroup.Post("/", articleHandler.CreateArticle)      // POST /api/v2/articles
	authGroup.Put("/:id", articleHandler.UpdateArticle)    // PUT /api/v2/articles/123
	authGroup.Delete("/:id", articleHandler.DeleteArticle) // DELETE /api/v2/articles/123
}
