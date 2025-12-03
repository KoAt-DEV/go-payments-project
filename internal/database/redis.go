package database

import (
	"context"
	"fmt"
	"go-payments-portfolio-project/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type RedisClient struct {
	*redis.Client
	log *zerolog.Logger
}

func NewRedis(cfg *config.Config, log *zerolog.Logger) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		PoolTimeout:  30 * time.Second,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	log.Info().
		Str("address", cfg.Redis.Addr).
		Int("db", cfg.Redis.DB).
		Int("pool_size", cfg.Redis.PoolSize).
		Msg("Redis connected successfully!")

	return &RedisClient{Client: client, log: log}, nil
}

func (r RedisClient) Shutdown(ctx context.Context) error {
	err := r.Client.Close()
	if err != nil {
		r.log.Warn().Err(err).Msg("Redis shutdown error")
	}
	return err
}
