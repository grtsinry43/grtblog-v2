package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/database"
	appserver "github.com/grtsinry43/grtblog-v2/server/internal/server"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	srv := appserver.New(cfg, db)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.Start(); err != nil && err != fiber.ErrServerClosed {
			log.Fatalf("fiber server failed: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
