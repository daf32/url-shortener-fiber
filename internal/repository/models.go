package repository

import (
	"time"

	"github.com/daf32/url-shortener-fiber/internal/core/domain"
)

type LinkModel struct {
	ID          int64     `db:"id"`
	Code        string    `db:"code"`
	OriginalURL string    `db:"original_url"`
	CreatedAt   time.Time `db:"created_at"`
}

func linkDomainFromModel(linkModel LinkModel) domain.Link {
	return domain.NewLink(
		linkModel.ID,
		linkModel.Code,
		linkModel.OriginalURL,
		linkModel.CreatedAt,
	)
}
