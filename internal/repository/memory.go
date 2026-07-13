package repository

import (
	"context"
	"sync"
	"time"

	"github.com/daf32/url-shortener-fiber/internal/core/domain"
	core_server "github.com/daf32/url-shortener-fiber/internal/core/server"
)

type MemoryRepo struct {
	mu    sync.Mutex
	links map[string]domain.Link
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{links: make(map[string]domain.Link)}
}

func (r *MemoryRepo) Save(ctx context.Context, code, url string) (domain.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	link := domain.NewLink(int64(len(r.links))+1, code, url, time.Now())
	r.links[code] = link

	return link, nil
}

func (r *MemoryRepo) Get(ctx context.Context, code string) (domain.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	link, ok := r.links[code]
	if !ok {
		return domain.Link{}, core_server.ErrNotFound
	}

	return link, nil
}
