package transport

import (
	"context"

	"github.com/daf32/url-shortener-fiber/internal/domain"
)

type ShortenerService interface {
	Shorten(ctx context.Context, url string) (domain.Link, error)
	Resolve(ctx context.Context, code string) (domain.Link, error)
}

type HTTPHandler struct {
	svc     ShortenerService
	baseURL string
}

func NewHTTPHandler(svc ShortenerService, baseURL string) *HTTPHandler {
	return &HTTPHandler{
		svc:     svc,
		baseURL: baseURL,
	}
}
