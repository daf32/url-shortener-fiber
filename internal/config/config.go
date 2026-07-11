package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server  Server
	Logger  Logger
	Limiter Limiter
}

type Server struct {
	Port    string `env:"SERVER_PORT"     env-default:"3000"`
	BaseURL string `env:"SERVER_BASE_URL" env-default:"http://localhost:3000"`
}

type Logger struct {
	Folder string `env:"LOG_FOLDER" env-default:"logs"`
}

type Limiter struct {
	Max        int           `env:"LIMITER_MAX"        env-default:"10"`
	Expiration time.Duration `env:"LIMITER_EXPIRATION" env-default:"1m"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	return &cfg, nil
}
