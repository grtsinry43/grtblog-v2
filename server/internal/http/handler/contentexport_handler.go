package handler

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"

	contentexportapp "github.com/grtsinry43/grtblog-v2/server/internal/app/contentexport"
	exportdomain "github.com/grtsinry43/grtblog-v2/server/internal/domain/contentexport"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type ContentExportHandler struct{ svc *contentexportapp.Service }

type createContentExportRequest struct {
	Mode string `json:"mode"`
}

func NewContentExportHandler(svc *contentexportapp.Service) *ContentExportHandler {
	return &ContentExportHandler{svc: svc}
}

func (h *ContentExportHandler) List(c *fiber.Ctx) error {
	items, err := h.svc.List(c.UserContext())
	if err != nil {
		return response.NewBizErrorWithCause(response.ServerError, "获取导出列表失败", err)
	}
	return response.Success(c, items)
}

func (h *ContentExportHandler) Create(c *fiber.Ctx) error {
	var req createContentExportRequest
	// 允许空 body（默认 both）。
	if len(c.Body()) > 0 {
		if err := c.BodyParser(&req); err != nil {
			return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
		}
	}
	item, err := h.svc.Create(c.UserContext(), req.Mode)
	if err != nil {
		if errors.Is(err, exportdomain.ErrExportRunning) {
			return response.NewBizErrorWithMsg(response.ParamsError, "已有导出任务正在运行")
		}
		return response.NewBizErrorWithCause(response.ParamsError, "创建导出任务失败", err)
	}
	return response.SuccessWithMessage(c, item, "导出任务已创建")
}

func (h *ContentExportHandler) Get(c *fiber.Ctx) error {
	item, err := h.svc.Get(c.UserContext(), c.Params("id"))
	if err != nil {
		return mapExportError(err, "获取导出详情失败")
	}
	return response.Success(c, item)
}

func (h *ContentExportHandler) Delete(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.UserContext(), c.Params("id")); err != nil {
		return mapExportError(err, "删除导出失败")
	}
	return response.SuccessWithMessage(c, fiber.Map{"id": c.Params("id")}, "导出已删除")
}

func (h *ContentExportHandler) IssueDownloadTicket(c *fiber.Ctx) error {
	token, expiresAt, err := h.svc.IssueDownloadTicket(c.UserContext(), c.Params("id"))
	if err != nil {
		return mapExportError(err, "生成下载链接失败")
	}
	path := "/api/v2/exports/download?ticket=" + url.QueryEscape(token)
	return response.Success(c, fiber.Map{"url": path, "expiresAt": expiresAt})
}

func (h *ContentExportHandler) Download(c *fiber.Ctx) error {
	token := strings.TrimSpace(c.Query("ticket"))
	if token == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "下载凭证不能为空")
	}
	item, path, err := h.svc.ResolveDownload(c.UserContext(), token)
	if err != nil {
		if errors.Is(err, exportdomain.ErrInvalidTicket) {
			return response.NewBizErrorWithMsg(response.Unauthorized, "下载链接无效或已过期")
		}
		return response.NewBizErrorWithCause(response.ServerError, "读取导出文件失败", err)
	}
	c.Set(fiber.HeaderContentType, "application/gzip")
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%q", item.Filename))
	c.Set(fiber.HeaderCacheControl, "private, no-store")
	return c.SendFile(path)
}

func mapExportError(err error, message string) error {
	switch {
	case errors.Is(err, exportdomain.ErrNotFound):
		return response.NewBizErrorWithMsg(response.NotFound, "导出记录不存在")
	case errors.Is(err, exportdomain.ErrExportRunning):
		return response.NewBizErrorWithMsg(response.ParamsError, "导出任务仍在运行")
	default:
		return response.NewBizErrorWithCause(response.ServerError, message, err)
	}
}
