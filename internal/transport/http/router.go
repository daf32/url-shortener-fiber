package transport

import (
	"io"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

type RouterConfig struct {
	LimiterMax        int
	LimiterExpiration time.Duration
}

func NewRouter(h *HTTPHandler, logWriter io.Writer, cfg RouterConfig) *fiber.App {
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})

	app.Use(recover.New())
	app.Use(logger.New(
		logger.Config{
			Stream: logWriter,
		},
	))

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("URL Shortener is running")
	})
	app.Post("/shorten", limiter.New(limiter.Config{
		Max:        cfg.LimiterMax,
		Expiration: cfg.LimiterExpiration,
	}), h.Shorten)

	app.Get("/:code", h.Resolve)

	return app
}
