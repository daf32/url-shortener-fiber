package http

import "github.com/daf32/url-shortener-fiber/internal/domain"

type ShortenerService interface {
	Shorten(url string) (domain.Link, error)
	Resolve(code string) (domain.Link, error)
}

type HTTPHanlder struct {
	svc ShortenerService
}

func NewHTTPHanlder(svc ShortenerService) *HTTPHanlder {
	return &HTTPHanlder{
		svc: svc,
	}
}
