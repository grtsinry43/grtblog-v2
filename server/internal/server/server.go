package server

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/router"
)

// Server wraps Fiber with configuration and dependencies.
type Server struct {
	cfg config.Config
	db  *gorm.DB
	app *fiber.App
}

// New builds a Fiber server with registered routes and middlewares.
func New(cfg config.Config, db *gorm.DB) *Server {
	app := fiber.New(fiber.Config{
		AppName:          cfg.App.Name,
		EnablePrintRoutes: cfg.App.Env == "development",
	})

	router.Register(app, router.Dependencies{
		DB:     db,
		Config: cfg,
	})

	return &Server{
		cfg: cfg,
		db:  db,
		app: app,
	}
}

// Start launches the Fiber HTTP server.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.cfg.App.Port)
	return s.app.Listen(addr)
}

// Shutdown gracefully stops Fiber.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

// App exposes the underlying Fiber instance for testing.
func (s *Server) App() *fiber.App {
	return s.app
}
