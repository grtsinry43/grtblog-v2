package handler

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"

	"github.com/gofiber/fiber/v2"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/identity"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/article"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type ArticleHandler struct {
	svc         *article.Service
	contentRepo content.Repository
	userRepo    identity.Repository
}

func NewArticleHandler(svc *article.Service, contentRepo content.Repository, userRepo identity.Repository) *ArticleHandler {
	return &ArticleHandler{
		svc:         svc,
		contentRepo: contentRepo,
		userRepo:    userRepo,
	}
}

// CreateArticle godoc
// @Summary 创建文章
// @Tags Article
// @Accept json
// @Produce json
// @Param request body contract.CreateArticleReq true "创建文章参数"
// @Success 200 {object} contract.ArticleResp
// @Security BearerAuth
// @Router /articles [post]
// @Security JWTAuth
func (h *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return response.ErrorFromBiz[any](c, response.NotLogin)
	}

	var req contract.CreateArticleReq
	if err := c.BodyParser(&req); err != nil {
		println(err.Error())
		return response.NewBizErrorWithMsg(response.ParamsError, "请求体解析失败")
	}

	var cmd article.CreateArticleCmd
	if err := copier.Copy(&cmd, req); err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "请求体映射失败")
	}

	createdArticle, err := h.svc.CreateArticle(c.Context(), claims.UserID, cmd)
	if err != nil {
		if errors.Is(err, content.ErrArticleShortURLExists) {
			return response.NewBizErrorWithMsg(response.ParamsError, "短链接已存在")
		}
		if errors.Is(err, content.ErrCategoryNotFound) {
			return response.NewBizErrorWithMsg(response.ParamsError, "分类不存在")
		}
		if errors.Is(err, content.ErrTagNotFound) {
			return response.NewBizErrorWithMsg(response.ParamsError, "标签不存在")
		}
		return err
	}

	articleResponse, err := h.toArticleResp(c.Context(), createdArticle)
	if err != nil {
		return err
	}

	Audit(c, "article.create", map[string]any{
		"articleId": createdArticle.ID,
		"title":     createdArticle.Title,
		"userId":    claims.UserID,
	})

	return response.SuccessWithMessage(c, articleResponse, "文章创建成功")
}

// UpdateArticle godoc
// @Summary 更新文章
// @Tags Article
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param request body contract.UpdateArticleReq true "更新文章参数"
// @Success 200 {object} contract.ArticleResp
// @Security BearerAuth
// @Router /articles/{id} [put]
// @Security JWTAuth
func (h *ArticleHandler) UpdateArticle(c *fiber.Ctx) error {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return response.ErrorFromBiz[any](c, response.NotLogin)
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的文章ID")
	}

	var req contract.UpdateArticleReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "请求体解析失败")
	}

	var cmd article.UpdateArticleCmd
	if err := copier.Copy(&cmd, req); err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "请求体映射失败")
	}
	cmd.ID = id

	updatedArticle, err := h.svc.UpdateArticle(c.Context(), cmd)
	if err != nil {
		if errors.Is(err, content.ErrArticleShortURLExists) {
			return response.NewBizErrorWithMsg(response.ParamsError, "短链接已存在")
		}
		if errors.Is(err, content.ErrCategoryNotFound) {
			return response.NewBizErrorWithMsg(response.ParamsError, "分类不存在")
		}
		if errors.Is(err, content.ErrTagNotFound) {
			return response.NewBizErrorWithMsg(response.ParamsError, "标签不存在")
		}
		return err
	}

	articleResponse, err := h.toArticleResp(c.Context(), updatedArticle)
	if err != nil {
		return err
	}

	Audit(c, "article.update", map[string]any{
		"articleId": updatedArticle.ID,
		"title":     updatedArticle.Title,
		"userId":    claims.UserID,
	})

	return response.SuccessWithMessage(c, articleResponse, "文章更新成功")
}

// GetArticle godoc
// @Summary 获取文章详情
// @Tags Article
// @Produce json
// @Param id path int true "文章ID"
// @Security BearerAuth
// @Success 200 {object} contract.ArticleResp
// @Router /articles/{id} [get]
func (h *ArticleHandler) GetArticle(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的文章ID")
	}

	article, err := h.svc.GetArticleByID(c.Context(), id)
	if err != nil {
		return err
	}

	articleResponse, err := h.toArticleResp(c.Context(), article)
	if err != nil {
		return err
	}

	return response.Success(c, articleResponse)
}

// GetArticleByShortURL godoc
// @Summary 根据短链接获取文章
// @Tags Article
// @Produce json
// @Param shortUrl path string true "短链接"
// @Success 200 {object} contract.ArticleResp
// @Router /articles/short/{shortUrl} [get]
func (h *ArticleHandler) GetArticleByShortURL(c *fiber.Ctx) error {
	shortURL := c.Params("shortUrl")
	if shortURL == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "短链接不能为空")
	}

	article, err := h.svc.GetArticleByShortURL(c.Context(), shortURL)
	if err != nil {
		return err
	}

	articleResponse, err := h.toArticleResp(c.Context(), article)
	if err != nil {
		return err
	}

	return response.Success(c, articleResponse)
}

// ListArticles godoc
// @Summary 获取文章列表
// @Tags Article
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param categoryId query int false "分类ID"
// @Param tagId query int false "标签ID"
// @Param search query string false "搜索关键词"
// @Success 200 {object} contract.ArticleListResp
// @Router /articles [get]
func (h *ArticleHandler) ListArticles(c *fiber.Ctx) error {
	query := contract.ListArticlesReq{
		Page:     1,
		PageSize: 10,
	}

	// 解析查询参数
	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		query.Page = page
	}
	if pageSize, err := strconv.Atoi(c.Query("pageSize", "10")); err == nil && pageSize > 0 && pageSize <= 100 {
		query.PageSize = pageSize
	}
	if categoryID, err := strconv.ParseInt(c.Query("categoryId"), 10, 64); err == nil {
		query.CategoryID = &categoryID
	}
	if tagID, err := strconv.ParseInt(c.Query("tagId"), 10, 64); err == nil {
		query.TagID = &tagID
	}
	if search := c.Query("search"); search != "" {
		query.Search = &search
	}

	// 只有管理员可以查看未发布的文章
	_, hasAuth := middleware.GetClaims(c)
	if hasAuth {
		// TODO: 检查是否有管理员权限
		if publishedStr := c.Query("published"); publishedStr != "" {
			if published, err := strconv.ParseBool(publishedStr); err == nil {
				query.Published = &published
			}
		}
	} else {
		// 非登录用户只能看发布的文章
		published := true
		query.Published = &published
	}

	articles, total, err := h.svc.ListArticles(c.Context(), content.ArticleListOptionsInternal(query))
	if err != nil {
		return err
	}

	// 转换为响应DTO
	articleResponses := make([]contract.ArticleListItemResp, len(articles))
	for i, art := range articles {
		resp, err := h.toArticleListItemResp(c.Context(), art)
		if err != nil {
			return err
		}
		articleResponses[i] = *resp
	}

	listResponse := contract.ArticleListResp{
		Items: articleResponses,
		Total: total,
		Page:  query.Page,
		Size:  query.PageSize,
	}

	return response.Success(c, listResponse)
}

// DeleteArticle godoc
// @Summary 删除文章
// @Tags Article
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} any
// @Security BearerAuth
// @Router /articles/{id} [delete]
// @Security JWTAuth
func (h *ArticleHandler) DeleteArticle(c *fiber.Ctx) error {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return response.ErrorFromBiz[any](c, response.NotLogin)
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的文章ID")
	}

	if err := h.svc.DeleteArticle(c.Context(), id); err != nil {
		return err
	}

	Audit(c, "article.delete", map[string]any{
		"articleId": id,
		"userId":    claims.UserID,
	})

	return response.SuccessWithMessage[any](c, nil, "文章删除成功")
}

func (h *ArticleHandler) toArticleResp(ctx context.Context, article *content.Article) (*contract.ArticleResp, error) {
	tags, err := h.svc.GetArticleTags(ctx, article.ID)
	if err != nil {
		return nil, err
	}

	metrics, err := h.svc.GetArticleMetrics(ctx, article.ID)
	if err != nil {
		return nil, err
	}

	var resp contract.ArticleResp
	if err := copier.Copy(&resp, article); err != nil {
		return nil, err
	}
	resp.TOC = mapTOCNodes(article.TOC)

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

func (h *ArticleHandler) toArticleListItemResp(ctx context.Context, article *content.Article) (*contract.ArticleListItemResp, error) {
	tags, err := h.svc.GetArticleTags(ctx, article.ID)
	if err != nil {
		return nil, err
	}

	metrics, err := h.svc.GetArticleMetrics(ctx, article.ID)
	if err != nil {
		return nil, err
	}

	resp := contract.ArticleListItemResp{
		ID:        article.ID,
		Title:     article.Title,
		ShortURL:  article.ShortURL,
		Summary:   article.Summary,
		IsTop:     article.IsTop,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}

	if article.Cover != nil {
		resp.Cover = *article.Cover
	}

	if metrics != nil {
		resp.Views = metrics.Views
		resp.Likes = metrics.Likes
		resp.Comments = metrics.Comments
	}

	if len(tags) > 0 {
		tagNames := make([]string, len(tags))
		for i, tag := range tags {
			tagNames[i] = tag.Name
		}
		resp.Tags = strings.Join(tagNames, ",")
	}

	if article.CategoryID != nil {
		category, err := h.contentRepo.GetCategoryByID(ctx, *article.CategoryID)
		if err != nil {
			if !errors.Is(err, content.ErrCategoryNotFound) {
				return nil, err
			}
		} else if category != nil {
			resp.CategoryName = category.Name
			if category.ShortURL != nil {
				resp.CategoryShortURL = *category.ShortURL
			}
		}
	}

	if h.userRepo != nil {
		user, err := h.userRepo.FindByID(ctx, article.AuthorID)
		if err != nil {
			if !errors.Is(err, identity.ErrUserNotFound) {
				return nil, err
			}
		} else if user != nil {
			resp.AuthorName = user.Nickname
			resp.Avatar = user.Avatar
		}
	}

	return &resp, nil
}
func mapTOCNodes(nodes []content.TOCNode) []contract.TOCNode {
	if len(nodes) == 0 {
		return nil
	}
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
