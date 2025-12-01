package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"

	"go-payments-portfolio-project/internal/config"
	"go-payments-portfolio-project/internal/logger"
	"go-payments-portfolio-project/internal/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("config load failed: %w", err))
	}

	log := logger.New(logger.Options{
		Environment: cfg.App.Env,
		ServiceName: "gopayments",
	})
	log.Info().Msg("Config loaded successfully!")

	app := fiber.New(fiber.Config{
		AppName:      "gopayments-api v1.0.0",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
	})

	app.Use(middleware.RequestID())
	app.Use(middleware.NewRequestLoggerMiddleware(middleware.RequestLoggerConfig{
		Logger:      log,
		SlowRequest: 500 * time.Millisecond,
	}))

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"env":       cfg.App.Env,
			"timestamp": time.Now().UTC(),
		})
	})

	addr := fmt.Sprintf(":%d", cfg.App.Port)
	log.Info().Int("port", cfg.App.Port).Msg("Starting Fiber server...")

	go func() {
		if err := app.Listen(addr); err != nil {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server stopped gracefully")
}
