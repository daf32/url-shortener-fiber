package core_server

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func DecodeAndValidateRequest(c fiber.Ctx, dest any) error {
	if err := c.Bind().Body(dest); err != nil {
		return fmt.Errorf(
			"decode and validate: %v: %w",
			err,
			ErrInvalidArgument,
		)
	}

	return nil
}
