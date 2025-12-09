package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-payments-portfolio-project/internal/config"
	"go-payments-portfolio-project/internal/database"
	"go-payments-portfolio-project/internal/logger"
	"go-payments-portfolio-project/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
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

	srv := server.New(cfg, log, pgPool, redisClient)

	log.Info().Int("port", cfg.App.Port).Msg("Starting Fiber server...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	startErr := make(chan error, 1)

	go func() {
		startErr <- srv.Start()
	}()

	var shutdownErr error

	select {
	case err := <-startErr:
		if err != nil {
			log.Fatal().Msg("Failed to start the server!")
		}
	case sig := <-quit:
		log.Info().Str("signal", sig.String()).Msg("Shutdown signal received...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server forced shutdown")
	} else {
		log.Info().Msg("HTTP server stopped cleanly")
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
