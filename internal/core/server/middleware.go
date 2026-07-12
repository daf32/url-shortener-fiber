package core_server

import (
	"time"

	core_logger "github.com/daf32/url-shortener-fiber/internal/core/logger"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func RequestLogger(log *core_logger.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		log = log.With(
			zap.String("request_id", c.RequestID()),
			zap.String("url", c.FullURL()),
		)
		
		before := time.Now()
		log.Debug(
			">>> incoming HTTP request",
			zap.String("http_method", c.Method()),
			zap.Time("time", before.UTC()),
		)
		
		err := c.Next()

		log.Debug(
			"<<< done HTTP request",
			zap.Int("status code", c.Response().StatusCode()),
			zap.Duration("latency", time.Since(before)),
		)

		return err
	}
}
