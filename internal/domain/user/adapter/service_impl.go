package adapter

import (
	"context"
	"errors"
	"fmt"
	"go-payments-portfolio-project/internal/config"
	"go-payments-portfolio-project/internal/domain/user"
	"go-payments-portfolio-project/internal/metrics"
	"go-payments-portfolio-project/internal/utils"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrPasswordHashingFailed  = errors.New("password hashing failed")
	ErrTokenGenerationFailed  = errors.New("token generation failed")
)

type ServiceImpl struct {
	repo   user.Repository
	cfg    *config.Config
	logger *zerolog.Logger
}

func NewService(repo user.Repository, cfg *config.Config, logger *zerolog.Logger) user.Service {
	return &ServiceImpl{repo: repo, cfg: cfg, logger: logger}
}

func (s *ServiceImpl) Register(ctx context.Context, req user.RegisterRequest) (*user.RegisterResponse, error) {
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if existing != nil {
		s.logger.Error().Err(err).Msg("Email already registered")
		metrics.RegisterErrorsWithReason.WithLabelValues("already_registered_email").Inc()
		return nil, ErrEmailAlreadyRegistered
	}

	startPasswordHash := time.Now()
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error().Err(err).Msg("Password hashing failed")
		metrics.RegisterErrorsWithReason.WithLabelValues("password_hash_failed").Inc()
		return nil, ErrPasswordHashingFailed
	}
	metrics.RegisterPasswordHashDurationSeconds.Observe(time.Since(startPasswordHash).Seconds())

	newUser := &user.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         req.Role,
		CoffeeCount:  0,
	}

	if newUser.Role == "" {
		newUser.Role = "user"
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		s.logger.Error().
			Err(err).
			Str("error_type", fmt.Sprintf("%T", err)).
			Str("full_error_msg", err.Error()).
			Msg("Repo Create failed during registration")

		return nil, err
	}

	accessToken, refreshToken, err := utils.GenerateTokens(newUser, s.cfg)
	if err != nil {
		s.logger.Error().Err(err).Msg("Token generation failed")
		metrics.RegisterErrorsWithReason.WithLabelValues("token_generation_failed").Inc()
		return nil, ErrTokenGenerationFailed
	}

	resp := &user.RegisterResponse{
		ID:           newUser.ID,
		Email:        newUser.Email,
		Role:         newUser.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	s.logger.Info().
		Str("user_id", newUser.ID.String()).
		Str("email", newUser.Email).
		Msg("User created successfully")

	return resp, nil
}
