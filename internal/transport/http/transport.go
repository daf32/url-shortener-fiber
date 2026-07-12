package transport

import (
	"context"
	"time"

	"github.com/daf32/url-shortener-fiber/internal/core/domain"
)

type ShortenerService interface {
	Shorten(ctx context.Context, url string) (domain.Link, error)
	Resolve(ctx context.Context, code string) (domain.Link, error)
}

type LimiterConfig struct {
	Max        int
	Expiration time.Duration
}

type ShortenerHTTPHandler struct {
	svc     ShortenerService
	baseURL string
	limiter LimiterConfig
}

func NewShortenerHTTPHandler(
	svc ShortenerService,
	baseURL string,
	limiter LimiterConfig,
) *ShortenerHTTPHandler {
	return &ShortenerHTTPHandler{
		svc:     svc,
		baseURL: baseURL,
		limiter: limiter,
	}
}
