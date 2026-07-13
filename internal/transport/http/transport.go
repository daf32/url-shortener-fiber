package transport

import (
	"context"
	"time"

	"github.com/daf32/url-shortener-fiber/internal/core/domain"
	core_logger "github.com/daf32/url-shortener-fiber/internal/core/logger"
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
	log     *core_logger.Logger
}

func NewShortenerHTTPHandler(
	svc ShortenerService,
	baseURL string,
	limiter LimiterConfig,
	log *core_logger.Logger,
) *ShortenerHTTPHandler {
	return &ShortenerHTTPHandler{
		svc:     svc,
		baseURL: baseURL,
		limiter: limiter,
		log:     log,
	}
}
