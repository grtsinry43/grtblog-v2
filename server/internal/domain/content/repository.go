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

	// MomentColumn 相关操作
	CreateColumn(ctx context.Context, column *MomentColumn) error
	GetColumnByID(ctx context.Context, id int64) (*MomentColumn, error)
	ListColumns(ctx context.Context) ([]*MomentColumn, error)
	UpdateColumn(ctx context.Context, column *MomentColumn) error
	DeleteColumn(ctx context.Context, id int64) error

	// Tag 相关操作
	CreateTag(ctx context.Context, tag *Tag) error
	GetTagByID(ctx context.Context, id int64) (*Tag, error)
	GetTagByName(ctx context.Context, name string) (*Tag, error)
	ListTags(ctx context.Context) ([]*Tag, error)
	UpdateTag(ctx context.Context, tag *Tag) error
	DeleteTag(ctx context.Context, id int64) error
	TagIDsExist(ctx context.Context, ids []int64) (bool, error)

	// ArticleTag 关联操作
	AddTagsToArticle(ctx context.Context, articleID int64, tagIDs []int64) error
	SyncTagsToArticle(ctx context.Context, articleID int64, tagIDs []int64) error // 这里因为每次提交更新才触发，所以不用 remove 而是 sync
	GetTagsByArticleID(ctx context.Context, articleID int64) ([]*Tag, error)

	// MomentTopic 关联操作
	AddTopicsToMoment(ctx context.Context, momentID int64, tagIDs []int64) error
	SyncTopicsToMoment(ctx context.Context, momentID int64, tagIDs []int64) error
	GetTopicsByMomentID(ctx context.Context, momentID int64) ([]*Tag, error)

	// Metrics 相关操作（这里用于统计交互信息，保证原子操作）
	UpdateArticleViews(ctx context.Context, articleID int64) error
	GetArticleMetrics(ctx context.Context, articleID int64) (*ArticleMetrics, error)

	UpdateMomentViews(ctx context.Context, momentID int64) error
	GetMomentMetrics(ctx context.Context, momentID int64) (*MomentMetrics, error)

	// Moment 相关操作
	CreateMoment(ctx context.Context, moment *Moment) error
	GetMomentByID(ctx context.Context, id int64) (*Moment, error)
	GetMomentByShortURL(ctx context.Context, shortURL string) (*Moment, error)
	UpdateMoment(ctx context.Context, moment *Moment) error
	DeleteMoment(ctx context.Context, id int64) error
	ListMoments(ctx context.Context, options MomentListOptionsInternal) ([]*Moment, int64, error)
	ListPublicMoments(ctx context.Context, options MomentListOptions) ([]*Moment, int64, error)

	// Page 相关操作
	CreatePage(ctx context.Context, page *Page) error
	GetPageByID(ctx context.Context, id int64) (*Page, error)
	GetPageByShortURL(ctx context.Context, shortURL string) (*Page, error)
	UpdatePage(ctx context.Context, page *Page) error
	DeletePage(ctx context.Context, id int64) error
	ListPages(ctx context.Context, options PageListOptionsInternal) ([]*Page, int64, error)
	ListPublicPages(ctx context.Context, options PageListOptions) ([]*Page, int64, error)

	UpdatePageViews(ctx context.Context, pageID int64) error
	GetPageMetrics(ctx context.Context, pageID int64) (*PageMetrics, error)
}
