package handler

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	mediaapp "github.com/grtsinry43/grtblog-v2/server/internal/app/media"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/media"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type UploadHandler struct {
	svc *mediaapp.Service
}

func NewUploadHandler(svc *mediaapp.Service) *UploadHandler {
	return &UploadHandler{svc: svc}
}

// UploadFile godoc
// @Summary 上传文件
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param type formData string true "上传类型: picture|file"
// @Success 200 {object} contract.UploadFileRespEnvelope
// @Security BearerAuth
// @Router /upload [post]
func (h *UploadHandler) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil || file == nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "file 不能为空")
	}
	fileType := c.FormValue("type")
	if strings.TrimSpace(fileType) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "type 不能为空")
	}

	result, err := h.svc.Upload(c.Context(), file, fileType)
	if err != nil {
		if errors.Is(err, media.ErrInvalidUploadType) {
			return response.NewBizErrorWithMsg(response.ParamsError, "type 仅支持 picture 或 file")
		}
		return response.NewBizErrorWithCause(response.ServerError, "文件上传失败", err)
	}

	resp := contract.ToUploadFileResp(result.File, !result.Created)
	msg := "上传成功"
	if !result.Created {
		msg = "文件已存在，返回已上传结果"
	}
	return response.SuccessWithMessage(c, resp, msg)
}

// ListUploads godoc
// @Summary 获取上传文件列表
// @Tags Upload
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} contract.UploadFileListRespEnvelope
// @Security BearerAuth
// @Router /uploads [get]
func (h *UploadHandler) ListUploads(c *fiber.Ctx) error {
	page := 1
	pageSize := 10
	if val, err := strconv.Atoi(c.Query("page", "1")); err == nil && val > 0 {
		page = val
	}
	if val, err := strconv.Atoi(c.Query("pageSize", "10")); err == nil && val > 0 && val <= 100 {
		pageSize = val
	}

	result, err := h.svc.List(c.Context(), page, pageSize)
	if err != nil {
		return response.NewBizErrorWithCause(response.ServerError, "获取文件列表失败", err)
	}

	items := make([]contract.UploadFileResp, len(result.Items))
	for i, file := range result.Items {
		items[i] = contract.ToUploadFileResp(file, false)
	}

	resp := contract.UploadFileListResp{
		Items: items,
		Total: result.Total,
		Page:  result.Page,
		Size:  result.Size,
	}
	return response.Success(c, resp)
}

// RenameUpload godoc
// @Summary 修改上传文件名
// @Tags Upload
// @Accept json
// @Produce json
// @Param id path int true "文件ID"
// @Param request body contract.UploadRenameReq true "新文件名"
// @Success 200 {object} contract.UploadFileRespEnvelope
// @Security BearerAuth
// @Router /upload/{id} [put]
func (h *UploadHandler) RenameUpload(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的文件ID")
	}

	var req contract.UploadRenameReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	if strings.TrimSpace(req.Name) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "文件名不能为空")
	}

	updated, err := h.svc.Rename(c.Context(), id, req.Name)
	if err != nil {
		if errors.Is(err, media.ErrUploadFileNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "文件不存在")
		}
		return response.NewBizErrorWithCause(response.ParamsError, "文件重命名失败", err)
	}

	return response.SuccessWithMessage(c, contract.ToUploadFileResp(*updated, false), "文件名已更新")
}

// DeleteUpload godoc
// @Summary 删除上传文件
// @Tags Upload
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} any
// @Security BearerAuth
// @Router /upload/{id} [delete]
func (h *UploadHandler) DeleteUpload(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的文件ID")
	}

	_, err = h.svc.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, media.ErrUploadFileNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "文件不存在")
		}
		return response.NewBizErrorWithCause(response.ServerError, "删除文件失败", err)
	}

	return response.SuccessWithMessage[any](c, nil, "文件已删除")
}

// DownloadUpload godoc
// @Summary 下载上传文件
// @Tags Upload
// @Produce application/octet-stream
// @Param id path int true "文件ID"
// @Success 200 {file} any
// @Security BearerAuth
// @Router /upload/{id}/download [get]
func (h *UploadHandler) DownloadUpload(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的文件ID")
	}

	file, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, media.ErrUploadFileNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "文件不存在")
		}
		return response.NewBizErrorWithCause(response.ServerError, "获取文件失败", err)
	}

	diskPath, err := h.svc.ResolveDiskPath(file.Path)
	if err != nil {
		return response.NewBizErrorWithCause(response.ServerError, "文件路径解析失败", err)
	}
	if _, err := os.Stat(diskPath); err != nil {
		if os.IsNotExist(err) {
			return response.NewBizErrorWithMsg(response.NotFound, "文件不存在")
		}
		return response.NewBizErrorWithCause(response.ServerError, "读取文件失败", err)
	}

	return c.Download(diskPath, file.Name)
}
