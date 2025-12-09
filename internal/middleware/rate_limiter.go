package middleware

import (
	"fmt"
	"strings"

	ratelimit "go-payments-portfolio-project/internal/limiter"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

type RateLimiterConfig struct {
	Limiter *ratelimit.RateLimiter
	Burst   int
	Rate    float64
	Logger  *zerolog.Logger
}

func NewRateLimiterMiddleware(cfg RateLimiterConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		key := "ip:" + strings.ReplaceAll(c.IP(), ":", "_")
		allowed, wait := cfg.Limiter.Allow(c.Context(), key, cfg.Burst, cfg.Rate)
		if !allowed {
			c.Set("Retry-After", fmt.Sprintf("%d", int(wait.Seconds()+1)))
			c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "too_many_request",
				"message":     "Rate limit exceeded",
				"retry_after": int(wait.Seconds() + 1),
			})
		}

		return c.Next()
	}
}
