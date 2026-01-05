package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/article"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type ArticleHandler struct {
	svc *article.Service
}

func NewArticleHandler(svc *article.Service) *ArticleHandler {
	return &ArticleHandler{svc: svc}
}

// CreateArticle godoc
// @Summary 创建文章
// @Tags Article
// @Accept json
// @Produce json
// @Param request body article.CreateArticleCommand true "创建文章参数"
// @Success 200 {object} article.ViewArticleResponse
// @Security BearerAuth
// @Router /articles [post]
// @Security JWTAuth
func (h *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return response.ErrorFromBiz[any](c, response.NotLogin)
	}

	var cmd article.CreateArticleCommand
	if err := c.BodyParser(&cmd); err != nil {
		println(err.Error())
		return response.NewBizErrorWithMsg(response.ParamsError, "请求体解析失败")
	}

	createdArticle, err := h.svc.CreateArticle(c.Context(), claims.UserID, cmd)
	if err != nil {
		return err
	}

	articleResponse, err := h.svc.ToResponse(c.Context(), createdArticle)
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
// @Param request body article.UpdateArticleCommand true "更新文章参数"
// @Success 200 {object} article.ViewArticleResponse
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

	var cmd article.UpdateArticleCommand
	if err := c.BodyParser(&cmd); err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "请求体解析失败")
	}

	// 设置ID
	cmd.ID = id

	updatedArticle, err := h.svc.UpdateArticle(c.Context(), cmd)
	if err != nil {
		return err
	}

	articleResponse, err := h.svc.ToResponse(c.Context(), updatedArticle)
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
// @Success 200 {object} article.ViewArticleResponse
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

	articleResponse, err := h.svc.ToResponse(c.Context(), article)
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
// @Success 200 {object} article.ViewArticleResponse
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

	articleResponse, err := h.svc.ToResponse(c.Context(), article)
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
// @Success 200 {object} article.ListArticleResponse
// @Router /articles [get]
func (h *ArticleHandler) ListArticles(c *fiber.Ctx) error {
	query := article.ListArticlesQuery{
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
	articleResponses := make([]article.ViewArticleResponse, len(articles))
	for i, art := range articles {
		resp, err := h.svc.ToResponse(c.Context(), art)
		if err != nil {
			return err
		}
		articleResponses[i] = *resp
	}

	listResponse := article.ListArticleResponse{
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
