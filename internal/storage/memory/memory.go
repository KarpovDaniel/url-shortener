package memory

import (
	"errors"
	"sync"

	"url-shortener/internal/storage"
)

type MemoryStorage struct {
	shortToOriginal sync.Map
	originalToShort sync.Map
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (s *MemoryStorage) Save(shortURL, originalURL string) error {
	if _, exists := s.shortToOriginal.Load(shortURL); exists {
		return errors.New("short URL already exists")
	}
	if _, exists := s.originalToShort.Load(originalURL); exists {
		return errors.New("original URL already exists")
	}
	s.shortToOriginal.Store(shortURL, originalURL)
	s.originalToShort.Store(originalURL, shortURL)
	return nil
}

func (s *MemoryStorage) Get(shortURL string) (string, error) {
	val, exists := s.shortToOriginal.Load(shortURL)
	if !exists {
		return "", storage.ErrNotFound
	}
	return val.(string), nil
}

func (s *MemoryStorage) FindByOriginalURL(originalURL string) (string, error) {
	val, exists := s.originalToShort.Load(originalURL)
	if !exists {
		return "", storage.ErrNotFound
	}
	return val.(string), nil
}
