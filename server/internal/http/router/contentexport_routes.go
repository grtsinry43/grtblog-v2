package router

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/article"
	contentexportapp "github.com/grtsinry43/grtblog-v2/server/internal/app/contentexport"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/moment"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/page"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/taxonomy"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/thinking"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

// registerContentExportRoutes 注册内置内容导出功能。所有服务在本地构造
// （与各资源路由 helper 的既有风格一致），不需要修改 Dependencies。
func registerContentExportRoutes(v2 fiber.Router, deps Dependencies) {
	contentRepo := persistence.NewContentRepository(deps.DB)
	commentRepo := persistence.NewCommentRepository(deps.DB)
	identityRepo := persistence.NewIdentityRepository(deps.DB)
	thinkingRepo := persistence.NewThinkingRepository(deps.DB)

	articleSvc := article.NewService(contentRepo, commentRepo, deps.EventBus)
	momentSvc := moment.NewService(contentRepo, commentRepo, deps.EventBus)
	pageSvc := page.NewService(contentRepo, commentRepo, deps.EventBus)
	thinkingSvc := thinking.NewService(thinkingRepo, commentRepo, deps.EventBus)
	categorySvc := taxonomy.NewCategoryService(persistence.NewArticleCategoryRepository(deps.DB))
	columnSvc := taxonomy.NewColumnService(persistence.NewMomentColumnRepository(deps.DB))
	tagSvc := taxonomy.NewTagService(persistence.NewTagRepository(deps.DB))

	collector := contentexportapp.NewCollector(articleSvc, momentSvc, pageSvc, thinkingSvc, categorySvc, columnSvc, tagSvc, deps.SysConfig)
	mapper := contentexportapp.NewMapper(articleSvc, momentSvc, pageSvc, thinkingSvc, contentRepo, commentRepo, identityRepo, deps.SysConfig)
	svc := contentexportapp.NewService(deps.Config.Export, deps.Config.Backup.UploadDir, persistence.NewContentExportRepository(deps.DB), collector, mapper)
	if err := svc.Initialize(context.Background()); err != nil {
		log.Printf("[contentexport] initialize failed: %v", err)
		return
	}

	h := handler.NewContentExportHandler(svc)
	// 下载路由在鉴权组外：ticket 即凭证（浏览器导航无法携带 Authorization 头）。
	v2.Get("/exports/download", h.Download)

	adminTokenRepo := persistence.NewAdminTokenRepository(deps.DB)
	authMiddleware := middleware.RequireAuth(deps.JWTManager, identityRepo, adminTokenRepo)
	adminMiddleware := middleware.RequireAdmin(identityRepo)
	exports := v2.Group("/admin/exports", authMiddleware, adminMiddleware)
	exports.Get("", h.List)
	exports.Post("", h.Create)
	exports.Get("/:id", h.Get)
	exports.Delete("/:id", h.Delete)
	exports.Post("/:id/download-ticket", h.IssueDownloadTicket)
}
