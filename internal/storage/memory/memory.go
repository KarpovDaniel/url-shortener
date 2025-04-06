package memory

import (
	"errors"
	"sync"
)

// Memory представляет потокобезопасное хранилище URL.
type Memory struct {
	shortToOriginal map[string]string
	originalToShort map[string]string
	mu              sync.RWMutex
}

// NewMemory создаёт новое экземпляр хранилища.
func NewMemory() *Memory {
	return &Memory{
		shortToOriginal: make(map[string]string),
		originalToShort: make(map[string]string),
	}
}

// Save сохраняет пару короткого и оригинального URL.
func (s *Memory) Save(shortURL, originalURL string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.shortToOriginal[shortURL]; exists {
		return "", errors.New("short URL already exists")
	}
	if actualShortURL, exists := s.originalToShort[originalURL]; exists {
		return actualShortURL, nil
	}

	s.shortToOriginal[shortURL] = originalURL
	s.originalToShort[originalURL] = shortURL
	return shortURL, nil
}

// Get возвращает оригинальный URL по короткому URL.
func (s *Memory) Get(shortURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	originalURL, exists := s.shortToOriginal[shortURL]
	if !exists {
		return "", errors.New("short URL not found")
	}
	return originalURL, nil
}
