package memory

import (
	"errors"
	"sync"
)

// MemoryStorage представляет потокобезопасное хранилище URL.
type MemoryStorage struct {
	shortToOriginal map[string]string // Короткий URL -> Оригинальный URL
	originalToShort map[string]string // Оригинальный URL -> Короткий URL
	mu              sync.RWMutex      // Мьютекс для синхронизации
}

// NewMemoryStorage создаёт новое экземпляр хранилища.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		shortToOriginal: make(map[string]string),
		originalToShort: make(map[string]string),
	}
}

// Save сохраняет пару короткого и оригинального URL.
func (s *MemoryStorage) Save(shortURL, originalURL string) (string, error) {
	s.mu.Lock()         // Блокировка для записи
	defer s.mu.Unlock() // Разблокировка после завершения

	// Проверка на дубликаты
	if _, exists := s.shortToOriginal[shortURL]; exists {
		return "", errors.New("short URL already exists")
	}
	if actualShortURL, exists := s.originalToShort[originalURL]; exists {
		return actualShortURL, nil
	}

	// Сохранение данных
	s.shortToOriginal[shortURL] = originalURL
	s.originalToShort[originalURL] = shortURL
	return shortURL, nil
}

// Get возвращает оригинальный URL по короткому URL.
func (s *MemoryStorage) Get(shortURL string) (string, error) {
	s.mu.RLock()         // Блокировка для чтения
	defer s.mu.RUnlock() // Разблокировка после завершения

	originalURL, exists := s.shortToOriginal[shortURL]
	if !exists {
		return "", errors.New("short URL not found")
	}
	return originalURL, nil
}
