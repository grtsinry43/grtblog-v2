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

	"github.com/grtsinry43/grtblog-v2/server/internal/app/moment"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type MomentHandler struct {
	svc         *moment.Service
	contentRepo content.Repository
	userRepo    identity.Repository
}

func NewMomentHandler(svc *moment.Service, contentRepo content.Repository, userRepo identity.Repository) *MomentHandler {
	return &MomentHandler{
		svc:         svc,
		contentRepo: contentRepo,
		userRepo:    userRepo,
	}
}

// CreateMoment godoc
// @Summary 创建手记
// @Tags Moment
// @Accept json
// @Produce json
// @Param request body contract.CreateMomentReq true "创建手记参数"
// @Success 200 {object} contract.MomentResp
// @Security BearerAuth
// @Router /moments [post]
// @Security JWTAuth
func (h *MomentHandler) CreateMoment(c *fiber.Ctx) error {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return response.ErrorFromBiz[any](c, response.NotLogin)
	}

	var req contract.CreateMomentReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}

	cmd := moment.CreateMomentCmd{
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		Image:       joinImages(req.Image),
		ColumnID:    req.ColumnID,
		TopicIDs:    req.TopicIDs,
		ShortURL:    req.ShortURL,
		IsPublished: req.IsPublished,
		IsTop:       req.IsTop,
		IsHot:       req.IsHot,
		IsOriginal:  req.IsOriginal,
		CreatedAt:   req.CreatedAt,
	}

	createdMoment, err := h.svc.CreateMoment(c.Context(), claims.UserID, cmd)
	if err != nil {
		if errors.Is(err, content.ErrMomentShortURLExists) {
			return response.NewBizErrorWithMsg(response.ParamsError, "短链接已存在")
		}
		if errors.Is(err, content.ErrColumnNotFound) {
			return response.NewBizErrorWithMsg(response.ParamsError, "分区不存在")
		}
		if errors.Is(err, content.ErrTagNotFound) {
			return response.NewBizErrorWithMsg(response.ParamsError, "话题不存在")
		}
		return err
	}

	momentResponse, err := h.toMomentResp(c.Context(), createdMoment)
	if err != nil {
		return err
	}

	Audit(c, "moment.create", map[string]any{
		"momentId": createdMoment.ID,
		"title":    createdMoment.Title,
		"userId":   claims.UserID,
	})

	return response.SuccessWithMessage(c, momentResponse, "手记创建成功")
}

// UpdateMoment godoc
// @Summary 更新手记
// @Tags Moment
// @Accept json
// @Produce json
// @Param id path int true "手记ID"
// @Param request body contract.UpdateMomentReq true "更新手记参数"
// @Success 200 {object} contract.MomentResp
// @Security BearerAuth
// @Router /moments/{id} [put]
// @Security JWTAuth
func (h *MomentHandler) UpdateMoment(c *fiber.Ctx) error {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return response.ErrorFromBiz[any](c, response.NotLogin)
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的手记ID")
	}

	var req contract.UpdateMomentReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}

	cmd := moment.UpdateMomentCmd{
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		Image:       joinImages(req.Image),
		ColumnID:    req.ColumnID,
		TopicIDs:    req.TopicIDs,
		ShortURL:    req.ShortURL,
		IsPublished: req.IsPublished,
		IsTop:       req.IsTop,
		IsHot:       req.IsHot,
		IsOriginal:  req.IsOriginal,
	}
	cmd.ID = id

	updatedMoment, err := h.svc.UpdateMoment(c.Context(), cmd)
	if err != nil {
		if errors.Is(err, content.ErrMomentShortURLExists) {
			return response.NewBizErrorWithMsg(response.ParamsError, "短链接已存在")
		}
		if errors.Is(err, content.ErrColumnNotFound) {
			return response.NewBizErrorWithMsg(response.ParamsError, "分区不存在")
		}
		if errors.Is(err, content.ErrTagNotFound) {
			return response.NewBizErrorWithMsg(response.ParamsError, "话题不存在")
		}
		return err
	}

	momentResponse, err := h.toMomentResp(c.Context(), updatedMoment)
	if err != nil {
		return err
	}

	Audit(c, "moment.update", map[string]any{
		"momentId": updatedMoment.ID,
		"title":    updatedMoment.Title,
		"userId":   claims.UserID,
	})

	return response.SuccessWithMessage(c, momentResponse, "手记更新成功")
}

// GetMoment godoc
// @Summary 获取手记详情
// @Tags Moment
// @Produce json
// @Param id path int true "手记ID"
// @Security BearerAuth
// @Success 200 {object} contract.MomentResp
// @Router /moments/{id} [get]
func (h *MomentHandler) GetMoment(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的手记ID")
	}

	momentItem, err := h.svc.GetMomentByID(c.Context(), id)
	if err != nil {
		return err
	}

	momentResponse, err := h.toMomentResp(c.Context(), momentItem)
	if err != nil {
		return err
	}

	return response.Success(c, momentResponse)
}

// GetMomentByShortURL godoc
// @Summary 根据短链接获取手记
// @Tags Moment
// @Produce json
// @Param shortUrl path string true "短链接"
// @Success 200 {object} contract.MomentResp
// @Router /moments/short/{shortUrl} [get]
func (h *MomentHandler) GetMomentByShortURL(c *fiber.Ctx) error {
	shortURL := c.Params("shortUrl")
	if shortURL == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "短链接不能为空")
	}

	momentItem, err := h.svc.GetMomentByShortURL(c.Context(), shortURL)
	if err != nil {
		return err
	}

	momentResponse, err := h.toMomentResp(c.Context(), momentItem)
	if err != nil {
		return err
	}

	return response.Success(c, momentResponse)
}

// ListMoments godoc
// @Summary 获取手记列表
// @Tags Moment
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param columnId query int false "分区ID"
// @Param topicId query int false "话题ID"
// @Param search query string false "搜索关键词"
// @Success 200 {object} contract.MomentListResp
// @Router /moments [get]
func (h *MomentHandler) ListMoments(c *fiber.Ctx) error {
	query := contract.ListMomentsReq{
		Page:     1,
		PageSize: 10,
	}

	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		query.Page = page
	}
	if pageSize, err := strconv.Atoi(c.Query("pageSize", "10")); err == nil && pageSize > 0 && pageSize <= 100 {
		query.PageSize = pageSize
	}
	if columnID, err := strconv.ParseInt(c.Query("columnId"), 10, 64); err == nil {
		query.ColumnID = &columnID
	}
	if topicID, err := strconv.ParseInt(c.Query("topicId"), 10, 64); err == nil {
		query.TopicID = &topicID
	}
	if search := c.Query("search"); search != "" {
		query.Search = &search
	}

	_, hasAuth := middleware.GetClaims(c)
	if hasAuth {
		if publishedStr := c.Query("published"); publishedStr != "" {
			if published, err := strconv.ParseBool(publishedStr); err == nil {
				query.Published = &published
			}
		}
	} else {
		published := true
		query.Published = &published
	}

	moments, total, err := h.svc.ListMoments(c.Context(), content.MomentListOptionsInternal(query))
	if err != nil {
		return err
	}

	momentResponses := make([]contract.MomentListItemResp, len(moments))
	for i, item := range moments {
		resp, err := h.toMomentListItemResp(c.Context(), item)
		if err != nil {
			return err
		}
		momentResponses[i] = *resp
	}

	listResponse := contract.MomentListResp{
		Items: momentResponses,
		Total: total,
		Page:  query.Page,
		Size:  query.PageSize,
	}

	return response.Success(c, listResponse)
}

// CheckMomentLatest godoc
// @Summary 校验手记是否最新
// @Tags Moment
// @Accept json
// @Produce json
// @Param id path int true "手记ID"
// @Param request body contract.CheckMomentLatestReq true "手记版本校验参数"
// @Success 200 {object} contract.CheckMomentLatestResp
// @Router /moments/{id}/latest [post]
func (h *MomentHandler) CheckMomentLatest(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的手记ID")
	}

	var req contract.CheckMomentLatestReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}

	momentItem, err := h.svc.GetMomentByID(c.Context(), id)
	if errors.Is(err, content.ErrMomentNotFound) {
		return response.NewBizErrorWithMsg(response.NotFound, "手记不存在")
	} else if err != nil {
		return err
	}

	if req.Hash == momentItem.ContentHash {
		return response.Success(c, contract.CheckMomentLatestResp{
			Latest: true,
			MomentContentPayload: contract.MomentContentPayload{
				ContentHash: momentItem.ContentHash,
			},
		})
	}

	return response.Success(c, contract.CheckMomentLatestResp{
		Latest: false,
		MomentContentPayload: contract.MomentContentPayload{
			ContentHash: momentItem.ContentHash,
			Title:       momentItem.Title,
			Summary:     momentItem.Summary,
			TOC:         mapMomentTOCNodes(momentItem.TOC),
			Content:     momentItem.Content,
		},
	})
}

// DeleteMoment godoc
// @Summary 删除手记
// @Tags Moment
// @Produce json
// @Param id path int true "手记ID"
// @Success 200 {object} any
// @Security BearerAuth
// @Router /moments/{id} [delete]
// @Security JWTAuth
func (h *MomentHandler) DeleteMoment(c *fiber.Ctx) error {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return response.ErrorFromBiz[any](c, response.NotLogin)
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的手记ID")
	}

	if err := h.svc.DeleteMoment(c.Context(), id); err != nil {
		return err
	}

	Audit(c, "moment.delete", map[string]any{
		"momentId": id,
		"userId":   claims.UserID,
	})

	return response.SuccessWithMessage[any](c, nil, "手记删除成功")
}

func (h *MomentHandler) toMomentResp(ctx context.Context, momentItem *content.Moment) (*contract.MomentResp, error) {
	topics, err := h.svc.GetMomentTopics(ctx, momentItem.ID)
	if err != nil {
		return nil, err
	}

	metrics, err := h.svc.GetMomentMetrics(ctx, momentItem.ID)
	if err != nil {
		return nil, err
	}

	resp := contract.MomentResp{
		ID:          momentItem.ID,
		Title:       momentItem.Title,
		Summary:     momentItem.Summary,
		AISummary:   momentItem.AISummary,
		TOC:         mapMomentTOCNodes(momentItem.TOC),
		Content:     momentItem.Content,
		ContentHash: momentItem.ContentHash,
		AuthorID:    momentItem.AuthorID,
		Image:       splitImages(momentItem.Image),
		ColumnID:    momentItem.ColumnID,
		ShortURL:    momentItem.ShortURL,
		IsPublished: momentItem.IsPublished,
		IsTop:       momentItem.IsTop,
		IsHot:       momentItem.IsHot,
		IsOriginal:  momentItem.IsOriginal,
		CreatedAt:   momentItem.CreatedAt,
		UpdatedAt:   momentItem.UpdatedAt,
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

func (h *MomentHandler) toMomentListItemResp(ctx context.Context, momentItem *content.Moment) (*contract.MomentListItemResp, error) {
	topics, err := h.svc.GetMomentTopics(ctx, momentItem.ID)
	if err != nil {
		return nil, err
	}

	metrics, err := h.svc.GetMomentMetrics(ctx, momentItem.ID)
	if err != nil {
		return nil, err
	}

	resp := contract.MomentListItemResp{
		ID:         momentItem.ID,
		Title:      momentItem.Title,
		ShortURL:   momentItem.ShortURL,
		Summary:    momentItem.Summary,
		IsTop:      momentItem.IsTop,
		IsHot:      momentItem.IsHot,
		IsOriginal: momentItem.IsOriginal,
		CreatedAt:  momentItem.CreatedAt,
		UpdatedAt:  momentItem.UpdatedAt,
		Topics:     []string{},
		Image:      splitImages(momentItem.Image),
	}

	if metrics != nil {
		resp.Views = metrics.Views
		resp.Likes = metrics.Likes
		resp.Comments = metrics.Comments
	}

	if len(topics) > 0 {
		topicNames := make([]string, len(topics))
		for i, topic := range topics {
			topicNames[i] = topic.Name
		}
		resp.Topics = topicNames
	}

	if momentItem.ColumnID != nil {
		column, err := h.contentRepo.GetColumnByID(ctx, *momentItem.ColumnID)
		if err != nil {
			if !errors.Is(err, content.ErrColumnNotFound) {
				return nil, err
			}
		} else if column != nil {
			resp.ColumnName = column.Name
			if column.ShortURL != nil {
				resp.ColumnShortURL = *column.ShortURL
			}
		}
	}

	if h.userRepo != nil {
		user, err := h.userRepo.FindByID(ctx, momentItem.AuthorID)
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

func mapMomentTOCNodes(nodes []content.TOCNode) []contract.TOCNode {
	result := make([]contract.TOCNode, len(nodes))
	for i, node := range nodes {
		result[i] = contract.TOCNode{
			Name:     node.Name,
			Anchor:   node.Anchor,
			Children: mapMomentTOCNodes(node.Children),
		}
	}
	return result
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

func joinImages(images []string) *string {
	if len(images) == 0 {
		return nil
	}
	out := make([]string, 0, len(images))
	for _, img := range images {
		item := strings.TrimSpace(img)
		if item == "" {
			continue
		}
		out = append(out, item)
	}
	if len(out) == 0 {
		return nil
	}
	joined := strings.Join(out, ",")
	return &joined
}
