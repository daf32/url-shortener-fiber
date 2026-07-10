package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("URL Shortener is running")
	})

	app.Post("/shorten", func(c fiber.Ctx) error {
		return c.SendString("here we will create a short link")
	})

	app.Get("/:code", func(c fiber.Ctx) error {
		code := c.Params("code")
		return c.SendString("you asked for code: " + code)
	})

	log.Fatal(app.Listen(":3000"))
}
