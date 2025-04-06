package storage

import "errors"

var ErrNotFound = errors.New("URL not found")

type Storage interface {
	Save(shortURL, originalURL string) error
	Get(shortURL string) (string, error)
	FindByOriginalURL(originalURL string) (string, error)
}
