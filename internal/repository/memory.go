package repository

import (
	"sync"

	"github.com/daf32/url-shortener-fiber/internal/domain"
)

type MemoryRepo struct {
	mu    sync.Mutex
	links map[string]domain.Link
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{links: make(map[string]domain.Link)}
}

func (r *MemoryRepo) Save(code, url string) (domain.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	link := domain.NewLink(code, url)
	r.links[code] = link

	return link, nil
}

func (r *MemoryRepo) Get(code string) (domain.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	link, ok := r.links[code]
	if !ok {
		return domain.Link{}, domain.ErrNotFound
	}

	return link, nil
}
