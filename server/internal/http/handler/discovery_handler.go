package handler

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/discovery"
	domain "github.com/grtsinry43/grtblog-v2/server/internal/domain/discovery"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type DiscoveryHandler struct{ service *discovery.Service }

// DiscoveryCatalogEnvelope documents the standard public API response.
type DiscoveryCatalogEnvelope struct {
	Code   int               `json:"code"`
	BizErr string            `json:"bizErr"`
	Msg    string            `json:"msg"`
	Data   discovery.Catalog `json:"data"`
	Meta   response.Meta     `json:"meta"`
}

func NewDiscoveryHandler(service *discovery.Service) *DiscoveryHandler {
	return &DiscoveryHandler{service}
}

// Catalog godoc
// @Summary 公开站点导航目录
// @Tags Discovery
// @Produce json
// @Success 200 {object} DiscoveryCatalogEnvelope
// @Router /public/discovery/catalog [get]
func (h *DiscoveryHandler) Catalog(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-store")
	catalog, err := h.service.Catalog(c.Context())
	if err != nil {
		return discoveryError(c, err)
	}
	return response.Success(c, catalog)
}

// Resource godoc
// @Summary 站点地图、robots、LLM 目录和公开 Markdown 正文
// @Tags Discovery
// @Produce plain
// @Param path path string true "公开资源路径，如 sitemap.xml 或 posts/example/index.md"
// @Param page query int false "LLM 目录页码"
// @Success 200 {string} string
// @Failure 404 {string} string
// @Failure 503 {string} string
// @Router /public/discovery/resource/{path} [get]
func (h *DiscoveryHandler) Resource(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-cache")
	c.Set("X-Content-Type-Options", "nosniff")
	path, err := url.PathUnescape(c.Params("*"))
	if err != nil {
		return discoveryError(c, domain.ErrNotFound)
	}
	path = "/" + strings.TrimPrefix(path, "/")
	if path == "/robots.txt" {
		body, err := h.service.Robots(c.Context())
		if err != nil {
			return discoveryError(c, err)
		}
		return sendDiscovery(c, []byte(body), "text/plain; charset=utf-8")
	}
	if strings.HasSuffix(path, "/index.md") {
		body, canonical, err := h.service.Document(c.Context(), strings.TrimSuffix(path, "index.md"))
		if err != nil {
			return discoveryError(c, err)
		}
		c.Set("Link", fmt.Sprintf("<%s>; rel=\"canonical\", </llms.txt>; rel=\"describedby\"", canonical))
		c.Set("X-Robots-Tag", "noindex")
		return sendDiscovery(c, []byte(body), "text/markdown; charset=utf-8")
	}
	// Reject unrelated paths before performing any catalog queries.
	scope := ""
	switch path {
	case "/sitemap.xml", "/llms.txt":
	case "/posts/llms.txt":
		scope = "posts"
	case "/moments/llms.txt":
		scope = "moments"
	case "/pages/llms.txt":
		scope = "pages"
	default:
		if !strings.HasPrefix(path, "/sitemaps/") || !strings.HasSuffix(path, ".xml") {
			return discoveryError(c, domain.ErrNotFound)
		}
	}
	catalog, err := h.service.Catalog(c.Context())
	if err != nil {
		return discoveryError(c, err)
	}
	if path == "/sitemap.xml" || strings.HasPrefix(path, "/sitemaps/") {
		page := 0
		if path != "/sitemap.xml" {
			raw := strings.TrimSuffix(strings.TrimPrefix(path, "/sitemaps/"), ".xml")
			page, err = strconv.Atoi(raw)
			if err != nil || page < 1 {
				return discoveryError(c, domain.ErrNotFound)
			}
		}
		body, err := catalog.Sitemap(page)
		if err != nil {
			return discoveryError(c, err)
		}
		return sendDiscovery(c, body, "application/xml; charset=utf-8")
	}
	page, err := discovery.ParsePage(c.Query("page"))
	if err != nil {
		return discoveryError(c, err)
	}
	body, err := catalog.LLMS(scope, page)
	if err != nil {
		return discoveryError(c, err)
	}
	c.Set("Link", "</llms.txt>; rel=\"describedby\"")
	return sendDiscovery(c, []byte(body), "text/markdown; charset=utf-8")
}

func sendDiscovery(c *fiber.Ctx, body []byte, contentType string) error {
	c.Set("Content-Type", contentType)
	etag := fmt.Sprintf("\"%x\"", sha256.Sum256(body))
	c.Set("ETag", etag)
	for _, candidate := range strings.Split(c.Get("If-None-Match"), ",") {
		candidate = strings.TrimSpace(candidate)
		if candidate == "*" || strings.TrimPrefix(candidate, "W/") == etag {
			return c.SendStatus(fiber.StatusNotModified)
		}
	}
	return c.Send(body)
}

func discoveryError(c *fiber.Ctx, err error) error {
	c.Set("Cache-Control", "no-store")
	c.Set("Content-Type", "text/plain; charset=utf-8")
	if errors.Is(err, domain.ErrNotFound) {
		return c.Status(404).SendString("Not Found")
	}
	log.Printf("[discovery] %s: %v", c.Path(), err)
	return c.Status(503).SendString("Discovery temporarily unavailable")
}
