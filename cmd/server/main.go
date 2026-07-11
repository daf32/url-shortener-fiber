package main

import (
	"log"

	"github.com/daf32/url-shortener-fiber/internal/repository"
	"github.com/daf32/url-shortener-fiber/internal/service"
	"github.com/daf32/url-shortener-fiber/internal/transport/http"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func main() {
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})

	repo := repository.NewMemoryRepo()
	svc := service.NewShortenerService(repo)
	httpHandler := http.NewHTTPHanlder(svc)

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("URL Shortener is running")
	})

	app.Post("/shorten", httpHandler.Shorten)

	app.Get("/:code", httpHandler.Resolve)

	log.Fatal(app.Listen(":3000"))
}
