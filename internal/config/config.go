package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		Port int    `env:"PORT" envDefault:"3000"`
		Env  string `env:"ENV" envDefault:"development"`
	} `envPrefix:"APP_"`

	Postgres struct {
		User     string `env:"USER"`
		Password string `env:"PASSWORD"`
		DBName   string `env:"DB"`
		Host     string `env:"HOST"`
		Port     int    `env:"PORT" envDefault:"5434"`
		SSLMode  string `env:"SSL_MODE" envDefault:"disable"`
	} `envPrefix:"POSTGRES_"`

	Redis struct {
		Addr         string `env:"ADDR" envDefault:"localhost:6379"`
		Password     string `env:"PASSWORD" envDefault:""`
		DB           int    `env:"DB" envDefault:"0"`
		PoolSize     int    `env:"POOL_SIZE" envDefault:"20"`
		MinIdleConns int    `env:"MIN_IDLE_CONNS" envDefault:"10"`
	} `envPrefix:"REDIS_"`

	JWT struct {
		Secret        string        `env:"SECRET"`
		AccesExpire   time.Duration `env:"ACCESS_EXPIRE" envDefault:"1h"`
		RefreshExpire time.Duration `env:"REFRESH_EXPIRE" envDefault:"720h"`
	} `envPrefix:"JWT_"`

	Metrics struct {
		NameSpace string `env:"NAMESPACE" envDefault:"gopayments"`
		SubSystem string `env:"SUB_SYSTEM" envDefault:"api"`
	} `envPrefix:"METRICS_"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil

}
