package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

// HealthHandler exposes lightweight probe endpoints for uptime monitoring.
type HealthHandler struct {
	cfg config.AppConfig
}

func NewHealthHandler(cfg config.AppConfig) *HealthHandler {
	return &HealthHandler{cfg: cfg}
}

func (h *HealthHandler) Liveness(c *fiber.Ctx) error {
	data := struct {
		Status string `json:"status"`
		App    string `json:"app"`
		Env    string `json:"env"`
	}{
		Status: "alive",
		App:    h.cfg.Name,
		Env:    h.cfg.Env,
	}

	return response.Success(c, data)
}

func (h *HealthHandler) Readiness(c *fiber.Ctx) error {
	data := struct {
		Status string    `json:"status"`
		Time   time.Time `json:"time"`
	}{
		Status: "ready",
		Time:   time.Now().UTC(),
	}

	return response.SuccessWithMessage(c, data, "ready")
}
