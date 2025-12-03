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
	"go-payments-portfolio-project/internal/database"
	"go-payments-portfolio-project/internal/logger"
	"go-payments-portfolio-project/internal/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load cinfig: %w", err))
	}

	log := logger.New(logger.Options{
		Environment: cfg.App.Env,
		ServiceName: "gopayments",
	})

	log.Info().Msg("Config loaded successfully!")

	pgPool, err := database.NewPostgres(cfg, log)
	if err != nil {
		log.Fatal().Msg("Failed to connect to Postgres")
	}

	redisClient, err := database.NewRedis(cfg, log)
	if err != nil {
		log.Fatal().Msg("Failed to connect to redis")
	}

	app := fiber.New(fiber.Config{
		AppName:      "go-payments-api",
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

	app.Get("/ping-db", func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := pgPool.Ping(ctx); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "fail",
				"message": "Postgres not reachable",
				"error":   err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Postgres reachable",
		})
	})

	app.Get("/ping-redis", func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := redisClient.Ping(ctx).Err(); err != nil {
			return c.JSON(fiber.Map{
				"status":  "failed",
				"message": "Redis is not reachable",
				"error":   err.Error(),
			})

		}

		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Redis reachable",
		})
	})

	addr := fmt.Sprintf(":%d", cfg.App.Port)
	log.Info().Int("port", cfg.App.Port).Msg("Starting Fiber server...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	srvErr := make(chan error, 1)

	go func() {
		srvErr <- app.Listen(addr)
	}()

	var shutdownErr error

	select {
	case err := <-srvErr:
		if err != nil {
			log.Fatal().Msg("Failed to start the server!")
		}
	case sig := <-quit:
		log.Info().Str("signal", sig.String()).Msg("Shutdown signal received...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error().Err(err).Msg("Fiber server forced to shutdown! Continuing cleanup...")
		shutdownErr = err
	}

	if err := pgPool.Shutdown(ctx); err != nil {
		log.Warn().Err(err).Msg("Postgres shutdown incomplete")
		shutdownErr = err
	}

	if err := redisClient.Shutdown(ctx); err != nil {
		log.Warn().Err(err).Msg("Redis shutdown incomplete")
		shutdownErr = err
	}

	if shutdownErr != nil {
		log.Fatal().Err(shutdownErr).Msg("Server exited with error during shutdown")
	}

	log.Info().Msg("Server stopped gracefully!")
}
