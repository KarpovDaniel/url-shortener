package postgres

import (
	"context"
	"database/sql"
	"errors"

	"url-shortener/internal/storage"
	"url-shortener/internal/storage/postgres/db"
)

type PostgresStorage struct {
	queries *db.Queries
}

func NewPostgresStorage(dbConn *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		queries: db.New(dbConn),
	}
}

func (s *PostgresStorage) Save(shortURL, originalURL string) error {
	err := s.queries.SaveURL(context.Background(), db.SaveURLParams{
		ShortUrl:    shortURL,
		OriginalUrl: originalURL,
	})
	return err
}

func (s *PostgresStorage) Get(shortURL string) (string, error) {
	originalURL, err := s.queries.GetURL(context.Background(), shortURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrNotFound
	}
	return originalURL, err
}

func (s *PostgresStorage) FindByOriginalURL(originalURL string) (string, error) {
	shortURL, err := s.queries.FindByOriginalURL(context.Background(), originalURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrNotFound
	}
	return shortURL, err
}
