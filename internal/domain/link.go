package domain

import "time"

type Link struct {
	ID          int64
	Code        string
	OriginalURL string
	CreatedAt   time.Time
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
