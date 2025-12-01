package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

type RequestLoggerConfig struct {
	Logger      *zerolog.Logger
	SlowRequest time.Duration
}

func NewRequestLoggerMiddleware(cfg RequestLoggerConfig) fiber.Handler {
	l := cfg.Logger
	if l == nil {
		panic("RequestLoggerMiddleware: log is nil")
	}

	if cfg.SlowRequest == 0 {
		cfg.SlowRequest = 500 * time.Millisecond
	}

	return func(c fiber.Ctx) error {
		start := time.Now()

		reqID := c.Locals("request_id").(string)

		err := c.Next()

		duration := time.Since(start)

		entry := l.With().
			Str("request_id", reqID).
			Str("path", string(c.Path())).
			Str("method", string(c.Method())).
			Int("status", c.Response().StatusCode()).
			Dur("duration", duration).
			Str("client_ip", c.IP()).
			Str("user_agent", string(c.Get("User-Agent"))).
			Logger()

		if err != nil {
			entry.Error().Err(err).Msg("request failed")
			return err
		}

		if duration > cfg.SlowRequest {
			entry.Warn().Msg("slow request")
			return nil
		}

		entry.Info().Msg("request completed")

		return nil
	}
}
