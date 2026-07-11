package transport

import "github.com/daf32/url-shortener-fiber/internal/domain"

type ShortenerService interface {
	Shorten(url string) (domain.Link, error)
	Resolve(code string) (domain.Link, error)
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
