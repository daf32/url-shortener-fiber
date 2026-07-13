package core_server

import (
	"time"

	core_logger "github.com/daf32/url-shortener-fiber/internal/core/logger"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func RequestLogger(log *core_logger.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		reqLog := log.With( // ← локальная, своя на каждый запрос
			zap.String("request_id", c.RequestID()),
			zap.String("url", c.FullURL()),
		)

		before := time.Now()
		reqLog.Debug(">>> incoming HTTP request",
			zap.String("http_method", c.Method()),
		)

		err := c.Next()

		reqLog.Debug("<<< done HTTP request",
			zap.Int("status_code", c.Response().StatusCode()),
			zap.Duration("latency", time.Since(before)),
		)

		return err
	}
}
