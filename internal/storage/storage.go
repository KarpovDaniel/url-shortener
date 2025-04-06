package storage

import "errors"

var ErrNotFound = errors.New("URL not found")

type Storage interface {
	Save(shortURL, originalURL string) (string, error)
	Get(shortURL string) (string, error)
}
