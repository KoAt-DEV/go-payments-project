package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func RequestID() fiber.Handler {
	return func(c fiber.Ctx) error {
		reqID := uuid.NewString()

		c.Set("X-Request-ID", reqID)
		c.Locals("request_id", reqID)

		return c.Next()
	}
}
