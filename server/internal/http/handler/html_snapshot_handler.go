package handler

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/htmlsnapshot"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type HTMLSnapshotHandler struct {
	service *htmlsnapshot.Service
}

func NewHTMLSnapshotHandler(service *htmlsnapshot.Service) *HTMLSnapshotHandler {
	return &HTMLSnapshotHandler{
		service: service,
	}
}

// RefreshPostsHTML godoc
// @Summary 刷新文章 HTML 缓存
// @Tags Public
// @Produce json
// @Success 200 {object} any
// @Router /public/html/posts/refresh [post]
func (h *HTMLSnapshotHandler) RefreshPostsHTML(c *fiber.Ctx) error {
	go func() {
		if h.service == nil {
			return
		}
		if err := h.service.RefreshPostsHTML(context.Background()); err != nil {
			log.Printf("[html-snapshot] generate posts html failed: %v", err)
		}
	}()

	return response.SuccessWithMessage[any](c, nil, "ok")
}
