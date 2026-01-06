package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/sysconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/router"
	infraevent "github.com/grtsinry43/grtblog-v2/server/internal/infra/event"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/jwt"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/turnstile"
)

// Server wraps Fiber with configuration and dependencies.
type Server struct {
	cfg     config.Config
	db      *gorm.DB
	app     *fiber.App
	logFile *os.File
}

// New builds a Fiber server with registered routes and middlewares.
func New(cfg config.Config, db *gorm.DB) *Server {
	logFile := initLogging()
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
			if reqID, ok := c.Locals("requestId").(string); ok && reqID != "" {
				log.Printf("[req:%s] unhandled error %s %s: %v", reqID, c.Method(), c.Path(), err)
			} else {
				log.Printf("unhandled error %s %s: %v", c.Method(), c.Path(), err)
			}
			return response.ErrorFromBiz[any](c, response.ServerError)
		},
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			reqID, _ := c.Locals("requestId").(string)
			stack := debug.Stack()
			if reqID != "" {
				log.Printf("[panic] req=%s %s %s: %v\n%s", reqID, c.Method(), c.Path(), e, stack)
			} else {
				log.Printf("[panic] %s %s: %v\n%s", c.Method(), c.Path(), e, stack)
			}
		},
	}))

	jwtManager := jwt.NewManager(cfg.Auth)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	turnstileClient := turnstile.NewClient(cfg.Turnstile)
	eventBus := infraevent.NewInMemoryBus()
	sysCfgRepo := persistence.NewSysConfigRepository(db)
	sysCfgSvc := sysconfig.NewService(sysCfgRepo, cfg.Turnstile)

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
		DB:         db,
		Config:     cfg,
		JWTManager: jwtManager,
		Turnstile:  turnstileClient,
		SysConfig:  sysCfgSvc,
		EventBus:   eventBus,
		Redis:      redisClient,
	})

	return &Server{
		cfg:     cfg,
		db:      db,
		app:     app,
		logFile: logFile,
	}
}

// Start launches the Fiber HTTP server.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.cfg.App.Port)
	return s.app.Listen(addr)
}

// Shutdown gracefully stops Fiber.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.logFile != nil {
		_ = s.logFile.Close()
	}
	return s.app.ShutdownWithContext(ctx)
}

// App exposes the underlying Fiber instance for testing.
func (s *Server) App() *fiber.App {
	return s.app
}

// initLogging sets a file logger under storage/logs/app.log while keeping stdout.
func initLogging() *os.File {
	logDir := filepath.Join("storage", "logs")
	_ = os.MkdirAll(logDir, 0o755)
	logPath := filepath.Join(logDir, "app.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Printf("failed to open log file: %v", err)
		return nil
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.LUTC)
	return f
}
