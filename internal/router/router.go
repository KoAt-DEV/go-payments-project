package router

import (
	"go-payments-portfolio-project/internal/domain/user/adapter"
	ratelimit "go-payments-portfolio-project/internal/limiter"
	"go-payments-portfolio-project/internal/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

func SetupRoutes(app *fiber.App, userHandler *adapter.Handler, rateLimiter *ratelimit.RateLimiter, log *zerolog.Logger) {
	api := app.Group("/api/v1")

	auth := api.Group("/auth")

	rateLimiterLogger := log.With().Str("component", "rate_limiter").Logger()

	rlConfig := middleware.RateLimiterConfig{
		Limiter: rateLimiter,
		Burst:   5,
		Rate:    1.0 / 60.0,
		Logger:  &rateLimiterLogger,
	}

	auth.Post("/register", middleware.NewRateLimiterMiddleware(rlConfig), userHandler.Register)

}
