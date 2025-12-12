package handler

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"

	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"

	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

// DocsHandler 提供 OpenAPI JSON 以及 Scalar UI。
type DocsHandler struct {
	specPath string
}

func NewDocsHandler(specPath string) *DocsHandler {
	return &DocsHandler{specPath: specPath}
}

// OpenAPI 返回 swagger/openapi JSON 文档。
func (h *DocsHandler) OpenAPI(c *fiber.Ctx) error {
	data, err := os.ReadFile(h.specPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":   response.ServerError.Code,
			"bizErr": response.ServerError.BizErr,
			"msg":    fmt.Sprintf("读取文档失败: %v", err),
			"data":   nil,
			"meta":   fiber.Map{},
		})
	}
	return c.Type("json").Send(data)
}

// Scalar 渲染 Scalar UI，加载同路径下的 openapi json。
func (h *DocsHandler) Scalar(c *fiber.Ctx) error {
	specURL := c.BaseURL() + "/docs/openapi.json"
	if specURL == "/docs/openapi.json" {
		specURL = "/docs/openapi.json"
	}
	html, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: specURL,
		CustomOptions: scalar.CustomOptions{
			PageTitle: "grtblog API V2",
		},
		DarkMode: true,
	})
	if err != nil {
		return response.ErrorWithMsg[any](c, response.ServerError, fmt.Sprintf("渲染文档失败: %v", err))
	}
	return c.Type("html", "utf-8").SendString(html)
}
