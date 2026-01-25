package handler

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/federationconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/federation"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/identity"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type FederationPostHandler struct {
	contentRepo   content.Repository
	userRepo      identity.Repository
	postCacheRepo federation.FederatedPostCacheRepository
	cfgSvc        *federationconfig.Service
}

func NewFederationPostHandler(contentRepo content.Repository, userRepo identity.Repository, postCacheRepo federation.FederatedPostCacheRepository, cfgSvc *federationconfig.Service) *FederationPostHandler {
	return &FederationPostHandler{
		contentRepo:   contentRepo,
		userRepo:      userRepo,
		postCacheRepo: postCacheRepo,
		cfgSvc:        cfgSvc,
	}
}

// GetPostDetail returns a single post with optional related posts.
// @Summary 联合文章详情
// @Tags Federation
// @Accept json
// @Produce json
// @Param id path string true "文章 ID 或短链接"
// @Success 200 {object} contract.FederationPostDetailResp
// @Router /api/federation/posts/{id} [get]
func (h *FederationPostHandler) GetPostDetail(c *fiber.Ctx) error {
	if h.cfgSvc != nil {
		if settings, err := h.cfgSvc.Settings(c.Context()); err == nil {
			if !settings.Enabled {
				return response.NewBizError(response.NotFound)
			}
		}
	}
	rawID := strings.TrimSpace(c.Params("id"))
	if rawID == "" {
		return response.NewBizError(response.ParamsError)
	}

	article, err := h.resolveArticle(c, rawID)
	if err != nil {
		if errors.Is(err, content.ErrArticleNotFound) {
			return response.NewBizError(response.NotFound)
		}
		return response.NewBizErrorWithCause(response.ServerError, "文章获取失败", err)
	}
	if !article.IsPublished {
		return response.NewBizError(response.NotFound)
	}

	baseURL := resolveFederationBaseURL(c, h.cfgSvc)
	post := h.buildPostResp(c.Context(), baseURL, article)

	related := h.relatedPosts(c, baseURL, article)

	resp := contract.FederationPostDetailResp{
		Post:         post,
		RelatedPosts: related,
	}
	return response.Success(c, resp)
}

func (h *FederationPostHandler) resolveArticle(c *fiber.Ctx, rawID string) (*content.Article, error) {
	if numericID, err := strconv.ParseInt(rawID, 10, 64); err == nil {
		return h.contentRepo.GetArticleByID(c.Context(), numericID)
	}
	return h.contentRepo.GetArticleByShortURL(c.Context(), rawID)
}

func (h *FederationPostHandler) buildPostResp(ctx context.Context, baseURL string, article *content.Article) contract.FederationPostResp {
	var authorName string
	var avatar *string
	if author, err := h.userRepo.FindByID(ctx, article.AuthorID); err == nil && author != nil {
		authorName = author.Nickname
		if authorName == "" {
			authorName = author.Username
		}
		if author.Avatar != "" {
			avatar = &author.Avatar
		}
	}
	return contract.FederationPostResp{
		ID:             article.ShortURL,
		URL:            baseURL + "/posts/" + article.ShortURL,
		Title:          article.Title,
		Summary:        article.Summary,
		ContentPreview: article.LeadIn,
		Author: contract.FederationPostAuthorResp{
			Name:   authorName,
			Avatar: avatar,
		},
		PublishedAt:   article.CreatedAt,
		UpdatedAt:     &article.UpdatedAt,
		CoverImage:    article.Cover,
		Language:      nil,
		AllowCitation: true,
		AllowComment:  true,
	}
}

func (h *FederationPostHandler) relatedPosts(c *fiber.Ctx, baseURL string, article *content.Article) []contract.FederationPostResp {
	const limit = 6

	local := make([]contract.FederationPostResp, 0, limit)
	tags, err := h.contentRepo.GetTagsByArticleID(c.Context(), article.ID)
	if err == nil && len(tags) > 0 {
		tagID := tags[0].ID
		items, _, err := h.contentRepo.ListPublicArticles(c.Context(), content.ArticleListOptions{
			Page:     1,
			PageSize: limit + 1,
			TagID:    &tagID,
		})
		if err == nil {
			for _, item := range items {
				if item.ID == article.ID {
					continue
				}
				local = append(local, h.buildPostResp(c.Context(), baseURL, item))
				if len(local) >= limit {
					break
				}
			}
		}
	}

	if len(local) >= limit {
		return local
	}

	remote, err := h.postCacheRepo.ListRecent(c.Context(), limit)
	if err != nil {
		return local
	}
	for _, item := range remote {
		if len(local) >= limit {
			break
		}
		local = append(local, mapRemotePostToResp(item))
	}
	return local
}

func mapRemotePostToResp(item federation.FederatedPostCache) contract.FederationPostResp {
	author := contract.FederationPostAuthorResp{Name: ""}
	var payload struct {
		Name   string  `json:"name"`
		URL    *string `json:"url,omitempty"`
		Avatar *string `json:"avatar,omitempty"`
	}
	if err := json.Unmarshal(item.Author, &payload); err == nil {
		author.Name = payload.Name
		author.URL = payload.URL
		author.Avatar = payload.Avatar
	}
	id := item.URL
	if item.RemotePostID != nil && *item.RemotePostID != "" {
		id = *item.RemotePostID
	}
	return contract.FederationPostResp{
		ID:             id,
		URL:            item.URL,
		Title:          item.Title,
		Summary:        item.Summary,
		ContentPreview: item.ContentPreview,
		Author:         author,
		PublishedAt:    item.PublishedAt,
		UpdatedAt:      item.UpdatedAt,
		CoverImage:     item.CoverImage,
		Language:       item.Language,
		AllowCitation:  item.AllowCitation,
		AllowComment:   item.AllowComment,
	}
}
