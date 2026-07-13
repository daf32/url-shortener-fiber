package transport

import (
	"fmt"

	"github.com/daf32/url-shortener-fiber/internal/core/domain"
	"github.com/gofiber/fiber/v3"
)

func DecodeAndValidateRequest(c fiber.Ctx, dest any) error {
	if err := c.Bind().Body(dest); err != nil {
		return fmt.Errorf(
			"decode and validate: %v: %w",
			err,
			domain.ErrInvalidArgument,
		)
	}

	return nil
}
