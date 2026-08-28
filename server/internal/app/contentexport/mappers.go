package contentexport

// 本文件复刻了 4 个 handler 私有的详情映射器，使导出 meta.json 与 admin API
// 详情响应（即旧 node 导出脚本的 meta）完全一致：
//   - toArticleResp  -> internal/http/handler/article_handler.go
//   - toMomentResp   -> internal/http/handler/moment_handler.go
//   - toPageResp     -> internal/http/handler/page_handler.go
//   - toThinkingResp -> internal/http/handler/thinking_handler.go
// 若上述映射器发生变更，请同步本文件（mappers_test.go 的 golden 测试可防漂移）。

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/jinzhu/copier"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/article"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/moment"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/page"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/sysconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/thinking"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/comment"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/identity"
	domainthinking "github.com/grtsinry43/grtblog-v2/server/internal/domain/thinking"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
)

// Mapper 持有复刻映射器所需的全部公开依赖。
type Mapper struct {
	articleSvc  *article.Service
	momentSvc   *moment.Service
	pageSvc     *page.Service
	thinkingSvc *thinking.Service
	contentRepo content.Repository
	commentRepo comment.CommentRepository
	userRepo    identity.Repository
	sysCfg      *sysconfig.Service
}

func NewMapper(
	articleSvc *article.Service,
	momentSvc *moment.Service,
	pageSvc *page.Service,
	thinkingSvc *thinking.Service,
	contentRepo content.Repository,
	commentRepo comment.CommentRepository,
	userRepo identity.Repository,
	sysCfg *sysconfig.Service,
) *Mapper {
	return &Mapper{
		articleSvc:  articleSvc,
		momentSvc:   momentSvc,
		pageSvc:     pageSvc,
		thinkingSvc: thinkingSvc,
		contentRepo: contentRepo,
		commentRepo: commentRepo,
		userRepo:    userRepo,
		sysCfg:      sysCfg,
	}
}

// ArticleResp 复刻 ArticleHandler.toArticleResp。
func (m *Mapper) ArticleResp(ctx context.Context, art *content.Article) (*contract.ArticleResp, error) {
	tags, err := m.articleSvc.GetArticleTags(ctx, art.ID)
	if err != nil {
		return nil, err
	}
	metrics, err := m.articleSvc.GetArticleMetrics(ctx, art.ID)
	if err != nil {
		return nil, err
	}

	var resp contract.ArticleResp
	if err := copier.Copy(&resp, art); err != nil {
		return nil, err
	}
	resp.TOC = mapTOCNodes(art.TOC)
	resp.ExtInfo = jsonRawFromBytes(art.ExtInfo)
	resp.AllowComment = m.allowCommentPtr(ctx, art.CommentID)
	resp.FediverseObjectURL = m.fediverseObjectURL(ctx, art)

	if art.CategoryID != nil {
		category, catErr := m.contentRepo.GetCategoryByID(ctx, *art.CategoryID)
		if catErr == nil && category != nil {
			resp.CategoryName = category.Name
			if category.ShortURL != nil {
				resp.CategoryShortURL = *category.ShortURL
			}
		}
	}

	if len(tags) > 0 {
		resp.Tags = make([]contract.TagResp, len(tags))
		for i, tag := range tags {
			if err := copier.Copy(&resp.Tags[i], tag); err != nil {
				return nil, err
			}
		}
	}

	if metrics != nil {
		var metricsResp contract.MetricsResp
		if err := copier.Copy(&metricsResp, metrics); err != nil {
			return nil, err
		}
		resp.Metrics = &metricsResp
	}

	return &resp, nil
}

// MomentResp 复刻 MomentHandler.toMomentResp；tz 为站点时区（每个任务只取一次）。
func (m *Mapper) MomentResp(ctx context.Context, tz *time.Location, momentItem *content.Moment) (*contract.MomentResp, error) {
	topics, err := m.momentSvc.GetMomentTopics(ctx, momentItem.ID)
	if err != nil {
		return nil, err
	}
	metrics, err := m.momentSvc.GetMomentMetrics(ctx, momentItem.ID)
	if err != nil {
		return nil, err
	}

	resp := contract.MomentResp{
		ID:                         momentItem.ID,
		Title:                      momentItem.Title,
		Summary:                    momentItem.Summary,
		AISummary:                  momentItem.AISummary,
		TOC:                        mapTOCNodes(momentItem.TOC),
		Content:                    momentItem.Content,
		ContentHash:                momentItem.ContentHash,
		AuthorID:                   momentItem.AuthorID,
		Image:                      splitImages(momentItem.Image),
		ActivityPubObjectID:        momentItem.ActivityPubObjectID,
		ActivityPubLastPublishedAt: momentItem.ActivityPubLastPublishedAt,
		ColumnID:                   momentItem.ColumnID,
		CommentID:                  momentItem.CommentID,
		ShortURL:                   momentItem.ShortURL,
		IsPublished:                momentItem.IsPublished,
		IsTop:                      momentItem.IsTop,
		IsHot:                      momentItem.IsHot,
		AllowComment:               m.allowCommentPtr(ctx, momentItem.CommentID),
		IsOriginal:                 momentItem.IsOriginal,
		ExtInfo:                    jsonRawFromBytes(momentItem.ExtInfo),
		ContentUpdatedAt:           momentItem.ContentUpdatedAt,
		CreatedAt:                  momentItem.CreatedAt.In(tz),
		UpdatedAt:                  momentItem.UpdatedAt,
	}

	if momentItem.ColumnID != nil {
		column, colErr := m.contentRepo.GetColumnByID(ctx, *momentItem.ColumnID)
		if colErr == nil && column != nil {
			resp.ColumnName = column.Name
			if column.ShortURL != nil {
				resp.ColumnShortURL = *column.ShortURL
			}
		}
	}

	if len(topics) > 0 {
		resp.Topics = make([]contract.TagResp, len(topics))
		for i, topic := range topics {
			if err := copier.Copy(&resp.Topics[i], topic); err != nil {
				return nil, err
			}
		}
	}

	if metrics != nil {
		var metricsResp contract.MetricsResp
		if err := copier.Copy(&metricsResp, metrics); err != nil {
			return nil, err
		}
		resp.Metrics = &metricsResp
	}

	return &resp, nil
}

// PageResp 复刻 PageHandler.toPageResp。
func (m *Mapper) PageResp(ctx context.Context, pageItem *content.Page) (*contract.PageResp, error) {
	metrics, err := m.pageSvc.GetPageMetrics(ctx, pageItem.ID)
	if err != nil {
		return nil, err
	}

	resp := contract.PageResp{
		ID:               pageItem.ID,
		Title:            pageItem.Title,
		Description:      pageItem.Description,
		AISummary:        pageItem.AISummary,
		TOC:              mapTOCNodes(pageItem.TOC),
		Content:          pageItem.Content,
		ContentHash:      pageItem.ContentHash,
		CommentID:        pageItem.CommentID,
		ShortURL:         pageItem.ShortURL,
		IsEnabled:        pageItem.IsEnabled,
		IsBuiltin:        pageItem.IsBuiltin,
		IsHot:            false,
		AllowComment:     m.allowCommentPtr(ctx, pageItem.CommentID),
		ExtInfo:          jsonRawFromBytes(pageItem.ExtInfo),
		ContentUpdatedAt: pageItem.ContentUpdatedAt,
		Metrics:          &contract.MetricsResp{Views: 0, Likes: 0, Comments: 0},
		CreatedAt:        pageItem.CreatedAt,
		UpdatedAt:        pageItem.UpdatedAt,
	}

	if metrics != nil {
		resp.Metrics.Views = metrics.Views
		resp.Metrics.Likes = metrics.Likes
		resp.Metrics.Comments = metrics.Comments
	}

	return &resp, nil
}

// ThinkingResp 复刻 ThinkingHandler.toThinkingResp。
func (m *Mapper) ThinkingResp(ctx context.Context, t *domainthinking.Thinking) (*contract.ThinkingResp, error) {
	resp := &contract.ThinkingResp{
		ID:                         t.ID,
		CommentID:                  t.CommentID,
		Content:                    t.Content,
		AuthorID:                   t.AuthorID,
		ActivityPubObjectID:        t.ActivityPubObjectID,
		ActivityPubLastPublishedAt: t.ActivityPubLastPublishedAt,
		IsHot:                      false,
		AllowComment:               m.allowCommentID(ctx, t.CommentID),
		Metrics: contract.ThinkingMetrics{
			Views:    t.Metrics.Views,
			Likes:    t.Metrics.Likes,
			Comments: t.Metrics.Comments,
		},
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
	if m.userRepo != nil {
		user, err := m.userRepo.FindByID(ctx, t.AuthorID)
		if err != nil {
			if !errors.Is(err, identity.ErrUserNotFound) {
				return nil, err
			}
		} else if user != nil {
			resp.AuthorName = user.Nickname
			resp.Avatar = user.Avatar
		}
	}
	return resp, nil
}

func (m *Mapper) allowCommentPtr(ctx context.Context, areaID *int64) bool {
	if m.commentRepo == nil || areaID == nil || *areaID <= 0 {
		return true
	}
	area, err := m.commentRepo.GetAreaByID(ctx, *areaID)
	if err != nil || area == nil {
		return false
	}
	return !area.IsClosed
}

func (m *Mapper) allowCommentID(ctx context.Context, areaID int64) bool {
	if m.commentRepo == nil || areaID <= 0 {
		return true
	}
	area, err := m.commentRepo.GetAreaByID(ctx, areaID)
	if err != nil || area == nil {
		return false
	}
	return !area.IsClosed
}

func (m *Mapper) fediverseObjectURL(ctx context.Context, art *content.Article) *string {
	if m.sysCfg == nil || art == nil {
		return nil
	}
	settings, err := m.sysCfg.ActivityPubSettings(ctx)
	if err != nil || !settings.Enabled {
		return nil
	}
	if art.ActivityPubObjectID == nil {
		return nil
	}
	objectURL := strings.TrimSpace(*art.ActivityPubObjectID)
	if objectURL == "" {
		return nil
	}
	return &objectURL
}

func mapTOCNodes(nodes []content.TOCNode) []contract.TOCNode {
	result := make([]contract.TOCNode, len(nodes))
	for i, node := range nodes {
		result[i] = contract.TOCNode{
			Name:     node.Name,
			Anchor:   node.Anchor,
			Children: mapTOCNodes(node.Children),
		}
	}
	return result
}

func jsonRawFromBytes(value []byte) *contract.JSONRaw {
	trimmed := bytes.TrimSpace(value)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil
	}
	copied := append([]byte(nil), trimmed...)
	raw := contract.JSONRaw(copied)
	return &raw
}

func splitImages(input *string) []string {
	if input == nil {
		return []string{}
	}
	trimmed := strings.TrimSpace(*input)
	if trimmed == "" {
		return []string{}
	}
	parts := strings.Split(trimmed, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" {
			continue
		}
		out = append(out, item)
	}
	return out
}

// MetaNode 是写入 meta.json / flatten 文档头部的包装结构，
// 与旧 node 导出脚本的字段保持一致。
type MetaNode struct {
	Kind       string    `json:"kind"`
	ID         int64     `json:"id"`
	RoutePath  string    `json:"routePath"`
	SourcePath string    `json:"sourcePath"`
	ExportedAt time.Time `json:"exportedAt"`
	Metadata   any       `json:"metadata"`
}

// 以下 *Meta 类型利用"浅层同名字段遮蔽"把 content 从 JSON 中剔除：
// 外层 Content 声明为 json:"content,omitempty" 的 nil 指针，占据 content 这个
// 键名（深度 0 优先于内嵌结构体深度 1 的同名字段），nil + omitempty 即被省略，
// 其余字段保持原结构体顺序输出。（注意 `json:"-"` 不占键名，无法遮蔽深层字段。）

type articleMeta struct {
	*contract.ArticleResp
	Content *string `json:"content,omitempty"`
}

type momentMeta struct {
	*contract.MomentResp
	Content *string `json:"content,omitempty"`
}

type pageMeta struct {
	*contract.PageResp
	Content *string `json:"content,omitempty"`
}

type thinkingMeta struct {
	*contract.ThinkingResp
	Content *string `json:"content,omitempty"`
}

func marshalIndent(value any) ([]byte, error) {
	return json.MarshalIndent(value, "", "  ")
}
