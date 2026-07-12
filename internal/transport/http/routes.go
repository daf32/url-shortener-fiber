package transport

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func (h *ShortenerHTTPHandler) RegisterRoutes(r fiber.Router) {
	r.Get("/", func(c fiber.Ctx) error {
		return c.SendString("URL Shortener is running")
	})

	// Handlers run in the order given, so the limiter must precede h.Shorten.
	r.Post("/shorten", limiter.New(limiter.Config{
		Max:        h.limiter.Max,
		Expiration: h.limiter.Expiration,
	}), h.Shorten)

	r.Get("/:code", h.Resolve)
}
