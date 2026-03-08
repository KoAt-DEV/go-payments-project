package adapter

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"

	"go-payments-portfolio-project/internal/domain/user"
	"go-payments-portfolio-project/internal/metrics"
)

type Handler struct {
	service   user.Service
	logger    *zerolog.Logger
	validator *validator.Validate
}

func NewHandler(service user.Service, logger *zerolog.Logger, validator *validator.Validate) *Handler {
	return &Handler{service: service, logger: logger, validator: validator}
}

func (h *Handler) Register(c fiber.Ctx) error {
	start := time.Now()
	var req user.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		h.logger.Warn().Err(err).Msg("Body binding failed")
		metrics.RegisterErrorsWithReason.WithLabelValues("binding_error").Inc()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		validationErrs := err.(validator.ValidationErrors)
		h.logger.Warn().Err(err).Msg("Validation failed")
		metrics.RegisterErrorsWithReason.WithLabelValues("validation_error").Inc()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation failed",
			"details": validationErrs.Translate(nil),
		})
	}

	resp, err := h.service.Register(c.Context(), req)
	if err != nil {
		metrics.RegisterRequestTotal.WithLabelValues("failed").Inc()
		metrics.TotalRegisterRequestTime.Observe(time.Since(start).Seconds())
		switch {
		case errors.Is(err, ErrEmailAlreadyRegistered):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "email already in use",
			})
		case errors.Is(err, ErrPasswordHashingFailed):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "error with password generation",
			})
		case errors.Is(err, ErrTokenGenerationFailed):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "error with token generation",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	metrics.RegisterRequestTotal.WithLabelValues("success").Inc()
	metrics.RegisterSuccessTotal.Inc()
	metrics.TotalRegisterRequestTime.Observe(time.Since(start).Seconds())
	return c.Status(fiber.StatusCreated).JSON(resp)
}
