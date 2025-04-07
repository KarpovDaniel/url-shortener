package storage

import "errors"

// ErrNotFound возвращается когда URL не найден
var ErrNotFound = errors.New("URL not found")

// Storage определяет интерфейс для работы с хранилищем URL
type Storage interface {
	// Save сохраняет пару URL, возвращает существующий shortURL если originalURL уже есть
	Save(shortURL, originalURL string) (string, error)

	// Get возвращает оригинальный URL по его короткой версии
	Get(shortURL string) (string, error)
}
