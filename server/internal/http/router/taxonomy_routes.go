package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/taxonomy"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerTaxonomyPublicRoutes(v2 fiber.Router, deps Dependencies) {
	taxHandler := newTaxonomyHandler(deps)

	v2.Get("/categories", taxHandler.ListCategories)
	v2.Get("/columns", taxHandler.ListColumns)
	v2.Get("/tags", taxHandler.ListTags)
}

func registerTaxonomyAdminRoutes(v2 fiber.Router, deps Dependencies) {
	taxHandler := newTaxonomyHandler(deps)
	admin := v2.Group("", middleware.RequireAuth(deps.JWTManager), middleware.RequireAdmin())

	admin.Post("/admin/categories", taxHandler.CreateCategory)
	admin.Put("/admin/categories/:id", taxHandler.UpdateCategory)
	admin.Delete("/admin/categories/:id", taxHandler.DeleteCategory)

	admin.Post("/admin/columns", taxHandler.CreateColumn)
	admin.Put("/admin/columns/:id", taxHandler.UpdateColumn)
	admin.Delete("/admin/columns/:id", taxHandler.DeleteColumn)

	admin.Post("/admin/tags", taxHandler.CreateTag)
	admin.Put("/admin/tags/:id", taxHandler.UpdateTag)
	admin.Delete("/admin/tags/:id", taxHandler.DeleteTag)
}

func newTaxonomyHandler(deps Dependencies) *handler.TaxonomyHandler {
	categoryRepo := persistence.NewArticleCategoryRepository(deps.DB)
	columnRepo := persistence.NewMomentColumnRepository(deps.DB)
	tagRepo := persistence.NewTagRepository(deps.DB)

	categorySvc := taxonomy.NewCategoryService(categoryRepo)
	columnSvc := taxonomy.NewColumnService(columnRepo)
	tagSvc := taxonomy.NewTagService(tagRepo)

	return handler.NewTaxonomyHandler(categorySvc, columnSvc, tagSvc)
}
