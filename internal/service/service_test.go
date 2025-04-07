package service

import (
	"context"
	"errors"
	"testing"
	"url-shortener/internal/storage"
	"url-shortener/proto"

	"github.com/stretchr/testify/assert"
)

// FakeStorage — поддельное хранилище для тестов
type FakeStorage struct {
	storage map[string]string // Короткий URL -> Оригинальный URL
	err     error
}

func NewFakeStorage() *FakeStorage {
	return &FakeStorage{
		storage: make(map[string]string),
		err:     nil,
	}
}

func (f *FakeStorage) Save(shortURL, originalURL string) (string, error) {
	if f.err != nil {
		defer func() { f.err = nil }()
		return "", f.err
	}
	// Если URL уже сохранён — вернуть существующий короткий URL.
	for k, v := range f.storage {
		if v == originalURL {
			return k, nil
		}
	}
	if _, exists := f.storage[shortURL]; exists {
		return "", errors.New("short URL already exists")
	}
	f.storage[shortURL] = originalURL
	return shortURL, nil
}

func (f *FakeStorage) Get(shortURL string) (string, error) {
	originalURL, exists := f.storage[shortURL]
	if !exists {
		return "", storage.ErrNotFound
	}
	return originalURL, nil
}

func TestService_CreateURL(t *testing.T) {
	tests := []struct {
		name          string
		originalURL   string
		setup         func(*FakeStorage)
		expectedShort string // если не пустая — ожидается именно она
		expectedErr   error
	}{
		{
			name:        "Создание нового URL",
			originalURL: "https://example.com",
			setup:       func(f *FakeStorage) {},
			// Ожидается, что будет сгенерирован новый короткий URL случайной длины
			expectedShort: "",
			expectedErr:   nil,
		},
		{
			name:        "Повторное использование существующего URL",
			originalURL: "https://example.com",
			setup: func(f *FakeStorage) {
				// Имитируем, что URL уже сохранён с коротким значением "abc123"
				f.storage["abc123"] = "https://example.com"
			},
			expectedShort: "abc123",
			expectedErr:   nil,
		},
		{
			name:        "Конфликт короткого URL (повторная генерация)",
			originalURL: "https://newexample.com",
			setup: func(f *FakeStorage) {
				// При первом вызове Save вернётся ошибка о конфликте,
				// затем ошибка сбрасывается, и операция должна пройти успешно.
				f.err = errors.New("short URL already exists")
			},
			expectedShort: "",
			expectedErr:   nil,
		},
		{
			name:        "Непредвиденная ошибка",
			originalURL: "https://newexample.com",
			setup: func(f *FakeStorage) {
				f.err = errors.New("unexpected error")
			},
			expectedShort: "",
			expectedErr:   errors.New("unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeStorage := NewFakeStorage()
			tt.setup(fakeStorage)

			s := NewService(fakeStorage)
			req := &proto.CreateURLRequest{
				OriginalUrl: tt.originalURL,
			}
			resp, err := s.CreateURL(context.Background(), req)
			// Метод gRPC возвращает ошибку только при критических сбоях, остальные ошибки передаются в поле Error.
			if tt.expectedErr != nil {
				// Если ожидается ошибка, то err должен быть nil, а сообщение об ошибке в ответе должно совпадать.
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedErr.Error(), resp.Error)
				assert.Empty(t, resp.ShortUrl)
			} else {
				assert.NoError(t, err)
				assert.Empty(t, resp.Error)
				if tt.expectedShort != "" {
					assert.Equal(t, tt.expectedShort, resp.ShortUrl)
				} else {
					// Проверяем, что сгенерированный короткий URL имеет нужную длину и соответствует шаблону.
					assert.Len(t, resp.ShortUrl, shortURLLength)
					assert.Regexp(t, "^[a-zA-Z0-9_]+$", resp.ShortUrl)
				}
			}
		})
	}
}

func TestService_GetURL(t *testing.T) {
	tests := []struct {
		name        string
		shortURL    string
		setup       func(*FakeStorage)
		expectedURL string
		expectedErr error
	}{
		{
			name:     "Успешное получение URL",
			shortURL: "abc123",
			setup: func(f *FakeStorage) {
				f.storage["abc123"] = "https://example.com"
			},
			expectedURL: "https://example.com",
			expectedErr: nil,
		},
		{
			name:        "URL не найден",
			shortURL:    "xyz789",
			setup:       func(f *FakeStorage) {},
			expectedURL: "",
			expectedErr: storage.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeStorage := NewFakeStorage()
			tt.setup(fakeStorage)

			s := NewService(fakeStorage)
			req := &proto.GetURLRequest{
				ShortUrl: tt.shortURL,
			}
			resp, err := s.GetURL(context.Background(), req)
			// Как и в CreateURL, ошибка в gRPC-методе передаётся через поле Error.
			if tt.expectedErr != nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedErr.Error(), resp.Error)
				assert.Empty(t, resp.OriginalUrl)
			} else {
				assert.NoError(t, err)
				assert.Empty(t, resp.Error)
				assert.Equal(t, tt.expectedURL, resp.OriginalUrl)
			}
		})
	}
}
