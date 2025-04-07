package service

import (
	"context"
	"crypto/rand"
	"math/big"
	"strings"

	"url-shortener/internal/storage"
	"url-shortener/proto"
)

const (
	shortURLLength = 10
	chars          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

// Service реализует интерфейс URLShortenerServer
type Service struct {
	proto.UnimplementedURLShortenerServer
	storage storage.Storage
}

// NewService создаёт новый экземпляр сервиса с переданным хранилищем
func NewService(storage storage.Storage) *Service {
	return &Service{storage: storage}
}

// CreateURL реализует gRPC-метод для создания короткой ссылки
func (s *Service) CreateURL(_ context.Context, req *proto.CreateURLRequest) (*proto.CreateURLResponse, error) {
	originalURL := req.GetOriginalUrl()
	for {
		shortURL, err := generateShortURL()
		if err != nil {
			return &proto.CreateURLResponse{
				Error: err.Error(),
			}, nil
		}
		shortURL, err = s.storage.Save(shortURL, originalURL)
		if err == nil {
			return &proto.CreateURLResponse{
				ShortUrl: shortURL,
			}, nil
		}
		if strings.Contains(err.Error(), "short URL") || strings.Contains(err.Error(), "urls_pkey") {
			continue // если короткая ссылка уже существует — сгенерировать новую
		}
		return &proto.CreateURLResponse{
			Error: err.Error(),
		}, nil
	}
}

// GetURL реализует gRPC-метод для получения оригинального URL по короткому
func (s *Service) GetURL(_ context.Context, req *proto.GetURLRequest) (*proto.GetURLResponse, error) {
	shortURL := req.GetShortUrl()
	originalURL, err := s.storage.Get(shortURL)
	if err != nil {
		return &proto.GetURLResponse{
			Error: err.Error(),
		}, nil
	}
	return &proto.GetURLResponse{
		OriginalUrl: originalURL,
	}, nil
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
