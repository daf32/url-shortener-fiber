package main

import (
	"context"
	"os/signal"
	"syscall"

	core_config "github.com/daf32/url-shortener-fiber/internal/core/config"
	core_logger "github.com/daf32/url-shortener-fiber/internal/core/logger"
	core_server "github.com/daf32/url-shortener-fiber/internal/core/server"
	"github.com/daf32/url-shortener-fiber/internal/repository"

	"github.com/daf32/url-shortener-fiber/internal/service"
	transport "github.com/daf32/url-shortener-fiber/internal/transport/http"
	"github.com/gofiber/fiber/v3/log"

	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	cfg, err := core_config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	db, err := repository.NewDB(
		ctx,
		cfg.Postgres.DSN(),
		cfg.Postgres.MaxOpenConns,
		cfg.Postgres.MaxIdleConns,
		cfg.Postgres.ConnMaxLifetime,
		cfg.Postgres.ConnMaxIdleTime,
	)
	if err != nil {
		log.Fatal("failed to init postgres connection pool: ", err)
	}
	defer db.Close()

	repo := repository.NewShortenerRepo(db)
	svc := service.NewShortenerService(repo)
	handler := transport.NewShortenerHTTPHandler(
		svc,
		cfg.Server.BaseURL,
		transport.LimiterConfig{
			Max:        cfg.Limiter.Max,
			Expiration: cfg.Limiter.Expiration,
		},
	)

	logWriter, err := core_logger.NewLogWriter(cfg.Logger.Folder)
	if err != nil {
		log.Fatal(err)
	}


	httpServer := core_server.NewHTTPServer(
		core_server.ServerConfig{
			Port:            cfg.Server.Port,
			BaseURL:         cfg.Server.BaseURL,
			ShutdownTimeout: cfg.Server.ShutdownTimeout,
		},
		nil,
		core_server.WithMiddleware(
			recover.New(),
			logger.New(logger.Config{
				Stream: logWriter,
			}),
		),
		core_server.WithRoutes("/", handler),
	)

	if err := httpServer.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
