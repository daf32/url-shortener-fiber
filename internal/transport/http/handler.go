package http

import (
	"errors"

	"github.com/daf32/url-shortener-fiber/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type shortenerRequest struct {
	URL string `json:"url" validate:"required"`
}

type shortenResponse struct {
	Code     string `json:"code"`
	Short    string `json:"short"`
	Original string `json:"original"`
}

func (h *HTTPHanlder) Shorten(c fiber.Ctx) error {
	var req shortenerRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorResponse(
			err,
			"failed to decode and validate request",
		))
	}

	link, err := h.svc.Shorten(req.URL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorResponse(
			err,
			"failed to generate url code",
		))
	}

	resp := shortenResponse{
		Code: link.Code,
		Short: c.BaseURL() + "/" + link.Code,
		Original: link.Original,
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *HTTPHanlder) Resolve(c fiber.Ctx) error {
	code := c.Params("code")

	link, err := h.svc.Resolve(code)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(domain.NewErrorResponse(
				err,
				"code not found",
			))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorResponse(
			err,
			"failed to get link",
		))
	}

	return c.Redirect().To(link.Original)
}
