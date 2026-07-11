package main

import (
	"github.com/daf32/url-shortener-fiber/internal/config"
	"github.com/daf32/url-shortener-fiber/internal/logger"
	"github.com/daf32/url-shortener-fiber/internal/repository"
	"github.com/daf32/url-shortener-fiber/internal/service"
	transport "github.com/daf32/url-shortener-fiber/internal/transport/http"
	"github.com/gofiber/fiber/v3/log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewMemoryRepo()
	svc := service.NewShortenerService(repo)
	handler := transport.NewHTTPHandler(svc, cfg.Server.BaseURL)

	logWriter, err := logger.NewLogWriter(cfg.Logger.Folder)
	if err != nil {
		log.Fatal(err)
	}
	app := transport.NewRouter(handler, logWriter, transport.RouterConfig{
		LimiterMax:        cfg.Limiter.Max,
		LimiterExpiration: cfg.Limiter.Expiration,
	})

	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
