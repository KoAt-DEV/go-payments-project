package database

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"

	"go-payments-portfolio-project/internal/config"
)

var (
	//go:embed migrations/*.sql
	migrationsFS embed.FS
)

func RunMigrations(ctx context.Context, pg *PgPool, cfg *config.Config, log *zerolog.Logger) error {
	if !cfg.AutoMigrate {
		log.Info().Msg("Auto-migrations disabled → skipping")
		return nil
	}

	sqlDB := stdlib.OpenDBFromPool(pg.Pool)
	if sqlDB == nil {
		return errors.New("failed to create sql.DB from pgxpool")
	}

	defer sqlDB.Close()

	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}
	goose.SetVerbose(cfg.App.Env == "development")

	log.Info().Msg("Starting database migrations...")

	migrationCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	if err := goose.StatusContext(migrationCtx, sqlDB, "migrations"); err != nil {
		log.Warn().Err(err).Msg("Could not get migration status (continuing anyway)")
	}
	if err := goose.UpContext(migrationCtx, sqlDB, "migrations"); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Info().Msg("Database migrations completed successfully")
	return nil

}
