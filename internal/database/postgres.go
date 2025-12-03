package database

import (
	"context"
	"fmt"
	"go-payments-portfolio-project/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PgPool struct {
	*pgxpool.Pool
	log *zerolog.Logger
}

func NewPostgres(cfg *config.Config, log *zerolog.Logger) (*PgPool, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	pgConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	pgConfig.MaxConns = 25
	pgConfig.MinConns = 5
	pgConfig.MaxConnLifetime = time.Minute
	pgConfig.MaxConnIdleTime = 30 * time.Minute
	pgConfig.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	log.Info().
		Str("host", cfg.Postgres.Host).
		Int("port", cfg.Postgres.Port).
		Int("max_conns", int(pgConfig.MaxConns)).
		Msg("Postgres connected successfully!")

	return &PgPool{Pool: pool, log: log}, nil
}

func (p *PgPool) Shutdown(ctx context.Context) error {
	p.Pool.Close()

	p.log.Info().Msg("Postgres pool closed gracefully")

	return nil
}
