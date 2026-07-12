package service

import (
	"context"
	"math/rand/v2"

	"github.com/daf32/url-shortener-fiber/internal/domain"
)

type LinkRepository interface {
	Save(ctx context.Context, code, url string) (domain.Link, error)
	Get(ctx context.Context, code string) (domain.Link, error)
}

type ShortenerService struct {
	repo LinkRepository
}

func NewShortenerService(repo LinkRepository) *ShortenerService {
	return &ShortenerService{repo: repo}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

func (s *ShortenerService) Shorten(ctx context.Context, url string) (domain.Link, error) {
	code := generateCode(6)
	link, err := s.repo.Save(ctx, code, url)
	if err != nil {
		return domain.Link{}, err
	}
	return link, nil
}

func (s *ShortenerService) Resolve(ctx context.Context, code string) (domain.Link, error) {
	return s.repo.Get(ctx, code)
}
