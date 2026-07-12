package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/daf32/url-shortener-fiber/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

type ShortenerRepo struct {
	db *DB
}

func NewShortenerRepo(
	db *DB,
) *ShortenerRepo {
	return &ShortenerRepo{db: db}
}

const (
	uniqueViolationCode = "23505"
)

func (r *ShortenerRepo) Save(ctx context.Context, code, url string) (domain.Link, error) {
	const query = `
		INSERT INTO links (code, original_url)
		VALUES ($1, $2)
		RETURNING id, code, original_url, created_at
	`

	var linkModel LinkModel

	err := r.db.GetContext(ctx, &linkModel, query, code, url)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return domain.Link{}, domain.ErrCodeExists
		}

		return domain.Link{}, fmt.Errorf("save link: %w", err)
	}

	linkDomain := linkDomainFromModel(linkModel)

	return linkDomain, nil
}

func (r *ShortenerRepo) Get(ctx context.Context, code string) (domain.Link, error) {
	const query = `
		SELECT id, code, original_url, created_at
		FROM links
		WHERE code=$1
	`

	var linkModel LinkModel

	err := r.db.GetContext(ctx, &linkModel, query, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Link{}, domain.ErrNotFound
		}

		return domain.Link{}, fmt.Errorf("get link: %w", err)
	}

	linkDomain := linkDomainFromModel(linkModel)

	return linkDomain, nil
}
