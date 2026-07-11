package domain

import "github.com/gofiber/fiber/v3"

type ErrorResponse fiber.Map

func NewErrorResponse(err error, msg string) ErrorResponse {
	return ErrorResponse{
		"err": err.Error(),
		"message": msg,
	}
}