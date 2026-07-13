package transport

import (
	"errors"

	"github.com/daf32/url-shortener-fiber/internal/core/domain"
	core_logger "github.com/daf32/url-shortener-fiber/internal/core/logger"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type HTTPResponseHandler struct {
	log *core_logger.Logger
	c   fiber.Ctx
}

func NewHTTPResponseHandler(
	log *core_logger.Logger,
	c fiber.Ctx,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		c:   c,
	}
}

func (h *HTTPResponseHandler) JSONResponse(
	data any,
	statusCode int,
) error {
	err := h.c.Status(statusCode).JSON(data)
	if err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}

	return err
}

func (h *HTTPResponseHandler) Redirect(
	location string,
) error {
	err := h.c.Redirect().To(location)
	if err != nil {
		h.log.Error("redirect", zap.String("location", location), zap.Error(err))
	}

	return err
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) error {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, domain.ErrCodeExists):
		statusCode = fiber.StatusConflict
		logFunc = h.log.Warn
	case errors.Is(err, domain.ErrNotFound):
		statusCode = fiber.StatusNotFound
		logFunc = h.log.Debug
	case errors.Is(err, domain.ErrInvalidArgument):
		statusCode = fiber.StatusBadRequest
		logFunc = h.log.Warn
	default:
		statusCode = fiber.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	return h.c.Status(statusCode).JSON(ErrorResponse{Message: msg})
}
