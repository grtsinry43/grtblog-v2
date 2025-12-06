package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
)

// Dependencies collects the shared instances that handlers require.
type Dependencies struct {
	DB     *gorm.DB
	Config config.Config
}

// Register wires up all HTTP endpoints with middlewares.
func Register(app *fiber.App, deps Dependencies) {
	healthHandler := handler.NewHealthHandler(deps.Config.App)

	app.Get("/health/liveness", healthHandler.Liveness)
	app.Get("/health/readiness", healthHandler.Readiness)
}
