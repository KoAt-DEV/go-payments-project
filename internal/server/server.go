package server

import (
	"context"
	"fmt"
	"time"

	"go-payments-portfolio-project/internal/config"
	"go-payments-portfolio-project/internal/database"
	ratelimit "go-payments-portfolio-project/internal/limiter"
	"go-payments-portfolio-project/internal/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

type Server struct {
	app     *fiber.App
	cfg     *config.Config
	log     *zerolog.Logger
	limiter *ratelimit.RateLimiter
}

func New(cfg *config.Config, log *zerolog.Logger, pgPool *database.PgPool, redisClient *database.RedisClient) *Server {
	rateLimiter := ratelimit.New(redisClient.Client)

	app := fiber.New(fiber.Config{
		AppName:      "go-payments-api",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	app.Use(middleware.RequestID())

	requestLogger := log.With().Str("component", "http").Logger()
	app.Use(middleware.NewRequestLoggerMiddleware(middleware.RequestLoggerConfig{
		Logger:      &requestLogger,
		SlowRequest: 500 * time.Millisecond,
	}))

	rateLimiterLogger := log.With().Str("component", "rate_limiter").Logger()
	app.Use(middleware.NewRateLimiterMiddleware(middleware.RateLimiterConfig{
		Limiter: rateLimiter,
		Burst:   100,
		Rate:    20.0,
		Logger:  &rateLimiterLogger,
	}))

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"service":   "go-payments-api",
			"env":       cfg.App.Env,
			"version":   "v1.0.0",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
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

	return &Server{
		app:     app,
		cfg:     cfg,
		log:     log,
		limiter: rateLimiter,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.cfg.App.Port)
	s.log.Info().Int("port", s.cfg.App.Port).Msg("Starting HTTP server")
	return s.app.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
