package content

import "context"

// Repository 定义文章相关的持久化操作
type Repository interface {
	// Article 相关操作
	CreateArticle(ctx context.Context, article *Article) error
	GetArticleByID(ctx context.Context, id int64) (*Article, error)
	GetArticleByShortURL(ctx context.Context, shortURL string) (*Article, error)
	UpdateArticle(ctx context.Context, article *Article) error
	DeleteArticle(ctx context.Context, id int64) error
	ListArticles(ctx context.Context, options ArticleListOptionsInternal) ([]*Article, int64, error)
	ListPublicArticles(ctx context.Context, options ArticleListOptions) ([]*Article, int64, error)

	// ArticleCategory 相关操作
	CreateCategory(ctx context.Context, category *ArticleCategory) error
	GetCategoryByID(ctx context.Context, id int64) (*ArticleCategory, error)
	ListCategories(ctx context.Context) ([]*ArticleCategory, error)
	UpdateCategory(ctx context.Context, category *ArticleCategory) error
	DeleteCategory(ctx context.Context, id int64) error

	// Tag 相关操作
	CreateTag(ctx context.Context, tag *Tag) error
	GetTagByID(ctx context.Context, id int64) (*Tag, error)
	GetTagByName(ctx context.Context, name string) (*Tag, error)
	ListTags(ctx context.Context) ([]*Tag, error)
	UpdateTag(ctx context.Context, tag *Tag) error
	DeleteTag(ctx context.Context, id int64) error

	// ArticleTag 关联操作
	AddTagsToArticle(ctx context.Context, articleID int64, tagIDs []int64) error
	SyncTagsToArticle(ctx context.Context, articleID int64, tagIDs []int64) error // 这里因为每次提交更新才触发，所以不用 remove 而是 sync
	GetTagsByArticleID(ctx context.Context, articleID int64) ([]*Tag, error)

	// Metrics 相关操作（这里用于统计交互信息，保证原子操作）
	UpdateArticleViews(ctx context.Context, articleID int64) error
	GetArticleMetrics(ctx context.Context, articleID int64) (*ArticleMetrics, error)
}
