package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/federationconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/identity"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type FederationTimelineHandler struct {
	contentRepo content.Repository
	userRepo    identity.Repository
	cfgSvc      *federationconfig.Service
}

func NewFederationTimelineHandler(contentRepo content.Repository, userRepo identity.Repository, cfgSvc *federationconfig.Service) *FederationTimelineHandler {
	return &FederationTimelineHandler{contentRepo: contentRepo, userRepo: userRepo, cfgSvc: cfgSvc}
}

// ListTimelinePosts returns published articles for federation timeline.
// @Summary 联合时间线
// @Tags Federation
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param per_page query int false "每页数量"
// @Param since query string false "起始时间 RFC3339"
// @Param until query string false "结束时间 RFC3339"
// @Success 200 {object} contract.FederationTimelineResp
// @Router /api/federation/timeline/posts [get]
func (h *FederationTimelineHandler) ListTimelinePosts(c *fiber.Ctx) error {
	if h.cfgSvc != nil {
		if settings, err := h.cfgSvc.Settings(c.Context()); err == nil {
			if !settings.Enabled {
				return response.NewBizError(response.NotFound)
			}
		}
	}
	page := parseIntQuery(c, "page", 1)
	size := parseIntQuery(c, "per_page", 20)
	if size > 100 {
		size = 100
	}

	since := parseTimeQuery(c, "since")
	until := parseTimeQuery(c, "until")

	articles, total, err := h.contentRepo.ListPublicArticlesForFederation(c.Context(), since, until, page, size)
	if err != nil {
		return err
	}

	baseURL := resolveFederationBaseURL(c, h.cfgSvc)
	items := make([]contract.FederationPostResp, len(articles))
	userCache := make(map[int64]*identity.User)
	for i, article := range articles {
		author, ok := userCache[article.AuthorID]
		if !ok {
			user, err := h.userRepo.FindByID(c.Context(), article.AuthorID)
			if err == nil {
				author = user
				userCache[article.AuthorID] = user
			}
		}
		authorName := ""
		var avatar *string
		if author != nil {
			authorName = author.Nickname
			if authorName == "" {
				authorName = author.Username
			}
			if author.Avatar != "" {
				avatar = &author.Avatar
			}
		}
		items[i] = contract.FederationPostResp{
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

	resp := contract.FederationTimelineResp{
		Items: items,
		Total: total,
		Page:  page,
		Size:  size,
	}
	return response.Success(c, resp)
}

func parseIntQuery(c *fiber.Ctx, key string, fallback int) int {
	if raw := c.Query(key); raw != "" {
		if val, err := strconv.Atoi(raw); err == nil {
			return val
		}
	}
	return fallback
}

func parseTimeQuery(c *fiber.Ctx, key string) *time.Time {
	raw := c.Query(key)
	if raw == "" {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil
	}
	return &parsed
}

func resolveFederationBaseURL(c *fiber.Ctx, svc *federationconfig.Service) string {
	if svc != nil {
		if settings, err := svc.Settings(c.Context()); err == nil && strings.TrimSpace(settings.InstanceURL) != "" {
			return strings.TrimRight(settings.InstanceURL, "/")
		}
	}
	scheme := "https"
	if c.Protocol() != "" {
		scheme = c.Protocol()
	}
	return scheme + "://" + c.Hostname()
}
