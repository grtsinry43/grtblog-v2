package server

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
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
		AppName:           cfg.App.Name,
		EnablePrintRoutes: cfg.App.Env == "development",

		// 核心：全局错误处理，自动把业务错误包装成统一响应
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// 1. 我们自己抛出的业务错误：*response.AppError
			if ae, ok := err.(*response.AppError); ok {
				return response.ErrorWithMsg[any](c, ae.Biz, ae.Message)
			}

			// 2. Fiber 内置错误（比如 fiber.ErrNotFound / ErrMethodNotAllowed）
			if fe, ok := err.(*fiber.Error); ok {
				// 这里可以按需映射到你的 BizError
				switch fe.Code {
				case fiber.StatusNotFound:
					return response.ErrorFromBiz[any](c, response.NotFound)
				case fiber.StatusMethodNotAllowed:
					return response.ErrorFromBiz[any](c, response.MethodNotAllowed)
				default:
					// 其他 HTTP 错误，统一当作 SERVER_ERROR 或自定义映射
					return response.ErrorFromBiz[any](c, response.ServerError)
				}
			}

			// 3. 其他未识别错误，统一视为服务器内部错误
			// TODO: 这里可以打日志：log.Errorf("unhandled error: %v", err)
			return response.ErrorFromBiz[any](c, response.ServerError)
		},
	})

	// 中间件：为每个请求附加 requestId（Meta 用）
	app.Use(func(c *fiber.Ctx) error {
		if c.Locals("requestId") == nil {
			reqID := c.Get("X-Request-ID")
			if reqID == "" {
				reqID = uuid.NewString()
			}
			c.Locals("requestId", reqID)
		}
		return c.Next()
	})

	// 注册路由
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
