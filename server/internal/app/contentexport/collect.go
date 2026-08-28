package contentexport

import (
	"context"
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/article"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/moment"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/page"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/sysconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/taxonomy"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/thinking"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	domainthinking "github.com/grtsinry43/grtblog-v2/server/internal/domain/thinking"
)

// exportPageSize 是批量加载内容时的分页大小。必须 >= 1（GORM Limit(0) 会返回空集）。
const exportPageSize = 100

// Snapshot 是一次导出任务加载到的全量内容快照。
type Snapshot struct {
	Articles   []*content.Article
	Moments    []*content.Moment
	Pages      []*content.Page
	Thinkings  []*domainthinking.Thinking
	Categories []*content.ArticleCategory
	Columns    []*content.MomentColumn
	Tags       []*content.Tag
	SiteTZ     *time.Location
	SiteName   string
	SiteURL    string
	// PublicHost 是站点 public_url 的主机名，用于把老文章里的
	// https://<本站域名>/uploads/... 绝对引用判定为站内图片。
	PublicHost string
}

// Collector 通过现有内容应用服务批量加载全部内容（含未发布/禁用）。
type Collector struct {
	articleSvc  *article.Service
	momentSvc   *moment.Service
	pageSvc     *page.Service
	thinkingSvc *thinking.Service
	categorySvc *taxonomy.CategoryService
	columnSvc   *taxonomy.ColumnService
	tagSvc      *taxonomy.TagService
	sysCfg      *sysconfig.Service
}

func NewCollector(
	articleSvc *article.Service,
	momentSvc *moment.Service,
	pageSvc *page.Service,
	thinkingSvc *thinking.Service,
	categorySvc *taxonomy.CategoryService,
	columnSvc *taxonomy.ColumnService,
	tagSvc *taxonomy.TagService,
	sysCfg *sysconfig.Service,
) *Collector {
	return &Collector{
		articleSvc:  articleSvc,
		momentSvc:   momentSvc,
		pageSvc:     pageSvc,
		thinkingSvc: thinkingSvc,
		categorySvc: categorySvc,
		columnSvc:   columnSvc,
		tagSvc:      tagSvc,
		sysCfg:      sysCfg,
	}
}

// Collect 加载 admin 全集：Published/Enabled/Builtin 过滤一律传 nil。
func (c *Collector) Collect(ctx context.Context) (*Snapshot, error) {
	snap := &Snapshot{
		Articles:   make([]*content.Article, 0),
		Moments:    make([]*content.Moment, 0),
		Pages:      make([]*content.Page, 0),
		Thinkings:  make([]*domainthinking.Thinking, 0),
		Categories: make([]*content.ArticleCategory, 0),
		Columns:    make([]*content.MomentColumn, 0),
		Tags:       make([]*content.Tag, 0),
		SiteTZ:     time.UTC,
	}

	for pg := 1; ; pg++ {
		items, _, err := c.articleSvc.ListArticles(ctx, content.ArticleListOptionsInternal{Page: pg, PageSize: exportPageSize})
		if err != nil {
			return nil, err
		}
		snap.Articles = append(snap.Articles, items...)
		if len(items) < exportPageSize {
			break
		}
	}

	for pg := 1; ; pg++ {
		items, _, err := c.momentSvc.ListMoments(ctx, content.MomentListOptionsInternal{Page: pg, PageSize: exportPageSize})
		if err != nil {
			return nil, err
		}
		snap.Moments = append(snap.Moments, items...)
		if len(items) < exportPageSize {
			break
		}
	}

	for pg := 1; ; pg++ {
		items, _, err := c.pageSvc.ListPages(ctx, content.PageListOptionsInternal{Page: pg, PageSize: exportPageSize})
		if err != nil {
			return nil, err
		}
		snap.Pages = append(snap.Pages, items...)
		if len(items) < exportPageSize {
			break
		}
	}

	for offset := 0; ; offset += exportPageSize {
		items, _, err := c.thinkingSvc.List(ctx, exportPageSize, offset)
		if err != nil {
			return nil, err
		}
		snap.Thinkings = append(snap.Thinkings, items...)
		if len(items) < exportPageSize {
			break
		}
	}

	categories, err := c.categorySvc.List(ctx)
	if err != nil {
		return nil, err
	}
	snap.Categories = categories
	columns, err := c.columnSvc.List(ctx)
	if err != nil {
		return nil, err
	}
	snap.Columns = columns
	tags, err := c.tagSvc.List(ctx)
	if err != nil {
		return nil, err
	}
	snap.Tags = tags

	if c.sysCfg != nil {
		snap.SiteTZ = c.sysCfg.Timezone(ctx)
		if info, infoErr := c.sysCfg.WebsiteInfo(ctx); infoErr == nil {
			snap.SiteName = info["website_name"]
			snap.SiteURL = info["public_url"]
			snap.PublicHost = hostFromURL(info["public_url"])
		}
	}

	return snap, nil
}
