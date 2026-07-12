package domain

import "time"

type Link struct {
	ID          int64     `db:"id"           json:"-"`
	Code        string    `db:"code"         json:"code"`
	OriginalURL string    `db:"original_url" json:"original_url"`
	CreatedAt   time.Time `db:"created_at"   json:"created_at"`
}

func NewLink(
	id int64,
	code string,
	originalURL string,
	createdAt time.Time,
) Link {
	return Link{
		ID:          id,
		Code:        code,
		OriginalURL: originalURL,
		CreatedAt:   createdAt,
	}
}
