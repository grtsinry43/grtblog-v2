package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/comment"
	domaincomment "github.com/grtsinry43/grtblog-v2/server/internal/domain/comment"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type CommentHandler struct {
	svc *comment.Service
}

func NewCommentHandler(svc *comment.Service) *CommentHandler {
	return &CommentHandler{svc: svc}
}

// CreateCommentLogin godoc
// @Summary 创建评论（登录用户）
// @Tags Comment
// @Accept json
// @Produce json
// @Param areaId path int true "评论区ID"
// @Param request body contract.CreateCommentLoginReq true "创建评论参数"
// @Success 200 {object} contract.CreateCommentResp
// @Security JWTAuth
// @Router /comments/areas/{areaId} [post]
func (h *CommentHandler) CreateCommentLogin(c *fiber.Ctx) error {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return response.ErrorFromBiz[any](c, response.NotLogin)
	}

	areaID, err := parseInt64Param(c, "areaId")
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的评论区ID")
	}

	var req contract.CreateCommentLoginReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}

	var cmd comment.CreateCommentLoginCmd
	if err := copier.Copy(&cmd, req); err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "请求体映射失败")
	}
	cmd.AreaID = areaID

	meta := comment.RequestMeta{
		IP:        c.IP(),
		UserAgent: c.Get("User-Agent", ""),
	}
	created, err := h.svc.CreateCommentLogin(c.Context(), claims.UserID, cmd, meta)
	if err != nil {
		return h.mapCommentError(c, err)
	}
	resp := toCreateCommentResp(created)
	return response.SuccessWithMessage(c, resp, "评论创建成功")
}

// CreateCommentVisitor godoc
// @Summary 创建评论（访客）
// @Tags Comment
// @Accept json
// @Produce json
// @Param areaId path int true "评论区ID"
// @Param request body contract.CreateCommentVisitorReq true "创建评论参数"
// @Success 200 {object} contract.CreateCommentResp
// @Router /comments/areas/{areaId}/visitor [post]
func (h *CommentHandler) CreateCommentVisitor(c *fiber.Ctx) error {
	areaID, err := parseInt64Param(c, "areaId")
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的评论区ID")
	}

	var req contract.CreateCommentVisitorReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}

	cmd := comment.CreateCommentVisitorCmd{
		AreaID:   areaID,
		Content:  req.Content,
		ParentID: req.ParentID,
	}
	if req.NickName != nil {
		cmd.NickName = *req.NickName
	}
	if req.Email != nil {
		cmd.Email = *req.Email
	}
	cmd.Website = req.Website

	meta := comment.RequestMeta{
		IP:        c.IP(),
		UserAgent: c.Get("User-Agent", ""),
	}
	created, err := h.svc.CreateCommentVisitor(c.Context(), cmd, meta)
	if err != nil {
		return h.mapCommentError(c, err)
	}
	resp := toCreateCommentResp(created)
	return response.SuccessWithMessage(c, resp, "评论创建成功")
}

// ListCommentTree godoc
// @Summary 获取评论树
// @Tags Comment
// @Produce json
// @Param areaId path int true "评论区ID"
// @Success 200 {object} []contract.CommentNodeResp
// @Router /comments/areas/{areaId} [get]
func (h *CommentHandler) ListCommentTree(c *fiber.Ctx) error {
	areaID, err := parseInt64Param(c, "areaId")
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的评论区ID")
	}

	nodes, err := h.svc.ListCommentTree(c.Context(), areaID)
	if err != nil {
		return h.mapCommentError(c, err)
	}

	resp := make([]contract.CommentNodeResp, len(nodes))
	for i, node := range nodes {
		resp[i] = toCommentNodeResp(node)
	}
	return response.Success(c, resp)
}

func (h *CommentHandler) mapCommentError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domaincomment.ErrCommentAreaNotFound):
		return response.NewBizErrorWithMsg(response.NotFound, "评论区不存在")
	case errors.Is(err, domaincomment.ErrCommentParentNotFound):
		return response.NewBizErrorWithMsg(response.ParamsError, "父评论不存在")
	case errors.Is(err, domaincomment.ErrCommentTooDeep):
		return response.NewBizErrorWithMsg(response.ParamsError, "评论层级过深")
	case errors.Is(err, domaincomment.ErrCommentContentEmpty):
		return response.NewBizErrorWithMsg(response.ParamsError, "评论内容不能为空")
	case errors.Is(err, domaincomment.ErrCommentAreaClosed):
		return response.NewBizErrorWithMsg(response.ParamsError, "评论区已关闭")
	default:
		return err
	}
}

func parseInt64Param(c *fiber.Ctx, name string) (int64, error) {
	return strconv.ParseInt(c.Params(name), 10, 64)
}

func toCreateCommentResp(entity *domaincomment.Comment) contract.CreateCommentResp {
	return contract.CreateCommentResp{
		ID:        entity.ID,
		AreaID:    entity.AreaID,
		Content:   entity.Content,
		NickName:  entity.NickName,
		Location:  entity.Location,
		Platform:  entity.Platform,
		Browser:   entity.Browser,
		Website:   entity.Website,
		IsOwner:   entity.IsOwner,
		IsFriend:  entity.IsFriend,
		IsAuthor:  entity.IsAuthor,
		IsViewed:  entity.IsViewed,
		IsTop:     entity.IsTop,
		ParentID:  entity.ParentID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}

func toCommentNodeResp(node *comment.CommentNode) contract.CommentNodeResp {
	resp := contract.CommentNodeResp{
		ID:        node.Comment.ID,
		AreaID:    node.Comment.AreaID,
		Content:   node.Comment.Content,
		NickName:  node.Comment.NickName,
		Location:  node.Comment.Location,
		Platform:  node.Comment.Platform,
		Browser:   node.Comment.Browser,
		Website:   node.Comment.Website,
		IsOwner:   node.Comment.IsOwner,
		IsFriend:  node.Comment.IsFriend,
		IsAuthor:  node.Comment.IsAuthor,
		IsViewed:  node.Comment.IsViewed,
		IsTop:     node.Comment.IsTop,
		ParentID:  node.Comment.ParentID,
		CreatedAt: node.Comment.CreatedAt,
		UpdatedAt: node.Comment.UpdatedAt,
		DeletedAt: node.Comment.DeletedAt,
	}
	if len(node.Children) > 0 {
		resp.Children = make([]contract.CommentNodeResp, len(node.Children))
		for i, child := range node.Children {
			resp.Children[i] = toCommentNodeResp(child)
		}
	}
	return resp
}
