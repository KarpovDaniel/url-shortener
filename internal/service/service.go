package service

import (
	"crypto/rand"
	"math/big"
	"strings"

	"url-shortener/internal/storage"
)

const (
	shortURLLength = 10
	chars          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

// Service реализует бизнес-логику сервиса сокращения URL
type Service struct {
	storage storage.Storage
}

// NewService создает новый экземпляр сервиса с заданным хранилищем
func NewService(storage storage.Storage) *Service {
	return &Service{storage: storage}
}

// Create генерирует и сохраняет короткий URL для оригинального URL
func (s *Service) Create(originalURL string) (shortURL string, err error) {
	for {
		shortURL, err = generateShortURL()
		if err != nil {
			return "", err
		}
		shortURL, err = s.storage.Save(shortURL, originalURL)
		if err == nil {
			return shortURL, nil
		}
		if strings.Contains(err.Error(), "short URL") || strings.Contains(err.Error(), "urls_pkey") {
			continue
		}
		return "", err
	}
}

// Get возвращает оригинальный URL по короткому URL
func (s *Service) Get(shortURL string) (string, error) {
	return s.storage.Get(shortURL)
}

// generateShortURL генерирует случайный короткий URL заданной длины
func generateShortURL() (string, error) {
	var shortURL string
	for i := 0; i < shortURLLength; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		shortURL += string(chars[idx.Int64()])
	}
	return shortURL, nil
}
