package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_config "github.com/daf32/url-shortener-fiber/internal/core/config"
	core_logger "github.com/daf32/url-shortener-fiber/internal/core/logger"
	core_server "github.com/daf32/url-shortener-fiber/internal/core/server"
	"github.com/daf32/url-shortener-fiber/internal/repository"
	"go.uber.org/zap"

	"github.com/daf32/url-shortener-fiber/internal/service"
	transport "github.com/daf32/url-shortener-fiber/internal/transport/http"
	"github.com/gofiber/fiber/v3/log"

	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
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

	log, err := core_logger.NewLogger(core_logger.LoggerConfig{
		Folder: cfg.Logger.Folder,
		Level:  cfg.Logger.Level,
	})
	if err != nil {
		fmt.Println("failed to init application logger: ", err)
		os.Exit(1)
	}
	defer log.Close()

	db, err := repository.NewDB(
		ctx,
		cfg.Postgres.DSN(),
		cfg.Postgres.MaxOpenConns,
		cfg.Postgres.MaxIdleConns,
		cfg.Postgres.ConnMaxLifetime,
		cfg.Postgres.ConnMaxIdleTime,
	)
	if err != nil {
		log.Fatal("failed to init postgres connection pool: ", zap.Error(err))
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

	httpServer := core_server.NewHTTPServer(
		core_server.ServerConfig{
			Port:            cfg.Server.Port,
			BaseURL:         cfg.Server.BaseURL,
			ShutdownTimeout: cfg.Server.ShutdownTimeout,
		},
		log,
		core_server.WithMiddleware(
			requestid.New(),
			recover.New(),
			core_server.RequestLogger(log),
		),
		core_server.WithRoutes("/", handler),
	)

	if err := httpServer.Run(ctx); err != nil {
		log.Error("HTTP server run error", zap.Error(err))
	}
}
