package transport

import (
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func NewRouter(h *HTTPHandler, logWriter io.Writer) *fiber.App {
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
	app.Post("/shorten", h.Shorten)
	app.Get("/:code", h.Resolve)

	return app
}
