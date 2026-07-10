package main

import (
	"log"
	"math/rand/v2"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var (
	mtx       sync.Mutex
	linsk     = make(map[string]string)
	Validator = validator.New()
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func generateCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}

	return string(b)
}

type shortenerRequest struct {
	URL string `json:"url" validate:"required"`
}

func main() {
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("URL Shortener is running")
	})

	app.Post("/shorten", func(c fiber.Ctx) error {
		var req shortenerRequest

		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		code := generateCode(6)

		mtx.Lock()
		linsk[code] = req.URL
		mtx.Unlock()

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"code":     code,
			"short":    c.BaseURL() + "/" + code,
			"original": req.URL,
		})
	})

	app.Get("/:code", func(c fiber.Ctx) error {
		code := c.Params("code")

		mtx.Lock()
		original, ok := linsk[code]
		mtx.Unlock()

		if !ok {
			return c.Status(fiber.StatusNotFound).JSON(
				fiber.Map{"error": "code not found"},
			)
		}

		return c.Redirect().To(original)
	})

	log.Fatal(app.Listen(":3000"))
}
