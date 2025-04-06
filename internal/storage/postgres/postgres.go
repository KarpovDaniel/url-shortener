package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"strings"
	"url-shortener/internal/storage"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func (s *Postgres) Save(shortURL, originalURL string) (string, error) {
	query := squirrel.Insert("urls").
		Columns("short_url", "original_url").
		Values(shortURL, originalURL)

	_, err := query.RunWith(s.db).ExecContext(context.Background())
	if err == nil {
		return shortURL, nil
	}
	if strings.Contains(err.Error(), "urls_original_url_key") {
		query := squirrel.Select("short_url").
			From("urls").
			Where(squirrel.Eq{"original_url": originalURL})

		err = query.RunWith(s.db).QueryRowContext(context.Background()).Scan(&shortURL)
		if err == nil {
			return shortURL, nil
		}
	}
	return "", err
}

func (s *Postgres) Get(shortURL string) (string, error) {
	var originalURL string
	query := squirrel.Select("original_url").
		From("urls").
		Where(squirrel.Eq{"short_url": shortURL})

	row := query.RunWith(s.db).QueryRowContext(context.Background())
	err := row.Scan(&originalURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrNotFound
	}
	return originalURL, err
}
