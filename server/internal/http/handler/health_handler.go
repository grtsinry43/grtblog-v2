package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/config"
)

// HealthHandler exposes lightweight probe endpoints for uptime monitoring.
type HealthHandler struct {
	cfg config.AppConfig
}

func NewHealthHandler(cfg config.AppConfig) *HealthHandler {
	return &HealthHandler{cfg: cfg}
}

func (h *HealthHandler) Liveness(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "alive",
		"app":    h.cfg.Name,
		"env":    h.cfg.Env,
	})
}

func (h *HealthHandler) Readiness(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ready",
		"time":   time.Now().UTC(),
	})
}
