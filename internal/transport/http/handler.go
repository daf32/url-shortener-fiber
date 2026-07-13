package transport

import (
	core_server "github.com/daf32/url-shortener-fiber/internal/core/server"
	"github.com/gofiber/fiber/v3"
)

type shortenerRequest struct {
	URL string `json:"url" validate:"required"`
}

type shortenResponse struct {
	Code        string `json:"code"`
	Short       string `json:"short"`
	OriginalURL string `json:"original_url"`
}

func (h *ShortenerHTTPHandler) Shorten(c fiber.Ctx) error {
	responseHandler := core_server.NewHTTPResponseHandler(h.log, c)

	var req shortenerRequest

	if err := core_server.DecodeAndValidateRequest(c, &req); err != nil {
		return responseHandler.ErrorResponse(
			err,
			"failed to decode and validate request",
		)
	}

	link, err := h.svc.Shorten(c.Context(), req.URL)
	if err != nil {
		return responseHandler.ErrorResponse(
			err,
			"failed to generate url code",
		)
	}

	resp := shortenResponse{
		Code:        link.Code,
		Short:       h.baseURL + "/" + link.Code,
		OriginalURL: link.OriginalURL,
	}
	return responseHandler.JSONResponse(resp, fiber.StatusCreated)
}

func (h *ShortenerHTTPHandler) Resolve(c fiber.Ctx) error {
	responseHandler := core_server.NewHTTPResponseHandler(h.log, c)
	code := c.Params("code")

	link, err := h.svc.Resolve(c.Context(), code)
	if err != nil {
		return responseHandler.ErrorResponse(
			err,
			"failed to resolve link",
		)
	}

	return responseHandler.Redirect(link.OriginalURL)
}
