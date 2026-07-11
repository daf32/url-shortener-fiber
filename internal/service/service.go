package service

import (
	"math/rand/v2"

	"github.com/daf32/url-shortener-fiber/internal/domain"
)

type LinkRepository interface {
	Save(code, url string) (domain.Link, error)
	Get(code string) (domain.Link, error)
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

func (s *ShortenerService) Shorten(url string) (domain.Link, error) {
	code := generateCode(6)
	link, err := s.repo.Save(code, url)
	if err != nil {
		return domain.Link{}, err
	}
	return link, nil
}

func (s *ShortenerService) Resolve(code string) (domain.Link, error) {
	return s.repo.Get(code)
}
