package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   Server
	Logger   Logger
	Limiter  Limiter
	Postgres Postgres
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

type Postgres struct {
	Host     string `env:"POSTGRES_HOST"     env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT"     env-default:"5432"`
	User     string `env:"POSTGRES_USER"     env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName   string `env:"POSTGRES_DB"       env-required:"true"`
	SSLMode  string `env:"POSTGRES_SSLMODE"  env-default:"disable"`
}

func (p Postgres) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode,
	)
}

func Load() (*Config, error) {
	var cfg Config

	if _, err := os.Stat(".env"); err == nil {
		if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
			return nil, fmt.Errorf("read .env: %w", err)
		}
		return &cfg, nil
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read env: %w", err)
	}
	return &cfg, nil
}
