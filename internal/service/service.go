package service

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"

	"url-shortener/internal/storage"
)

const (
	shortURLLength = 10
	chars          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

type Service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Create(originalURL string) (shortURL string, err error) {
	// Проверяем, существует ли originalURL
	shortURL, err = s.storage.FindByOriginalURL(originalURL)
	if err == nil {
		return shortURL, nil
	}
	if !errors.Is(err, storage.ErrNotFound) {
		return "", err
	}
	// Если не найден, генерируем новый shortURL
	for {
		shortURL, err = generateShortURL()
		if err != nil {
			return "", err
		}
		err = s.storage.Save(shortURL, originalURL)
		if err == nil {
			return shortURL, nil
		}
		// Обрабатываем возможные конфликты shortURL
		if strings.Contains(err.Error(), "short URL") || strings.Contains(err.Error(), "urls_pkey") {
			continue
		}
		// Обрабатываем случай, если originalURL был добавлен параллельно
		if strings.Contains(err.Error(), "original URL") || strings.Contains(err.Error(), "urls_original") {
			shortURL, err = s.storage.FindByOriginalURL(originalURL)
			if err == nil {
				return shortURL, nil
			}
		}
		return "", err
	}
}

func (s *Service) Get(shortURL string) (string, error) {
	return s.storage.Get(shortURL)
}

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
