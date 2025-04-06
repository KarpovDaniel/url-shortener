package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"url-shortener/internal/storage"
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
	if _, exists := f.storage[shortURL]; exists {
		return "", errors.New("short URL already exists")
	}
	for k, v := range f.storage {
		if v == originalURL {
			return k, nil
		}
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

func TestService_Create(t *testing.T) {
	tests := []struct {
		name          string
		originalURL   string
		setup         func(*FakeStorage)
		expectedShort string
		expectedErr   error
	}{
		{
			name:          "Создание нового URL",
			originalURL:   "https://example.com",
			setup:         func(f *FakeStorage) {},
			expectedShort: "",
			expectedErr:   nil,
		},
		{
			name:        "Повторное использование существующего URL",
			originalURL: "https://example.com",
			setup: func(f *FakeStorage) {
				f.storage["abc123"] = "https://example.com"
			},
			expectedShort: "abc123",
			expectedErr:   nil,
		}, {
			name:        "Конфликт короткого URL",
			originalURL: "https://newexample.com",
			setup: func(f *FakeStorage) {
				f.err = errors.New("short URL already exists")
			},
			expectedShort: "",
			expectedErr:   nil,
		}, {
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
			shortURL, err := s.Create(tt.originalURL)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				assert.Empty(t, shortURL)
			} else {
				assert.NoError(t, err)
				if tt.expectedShort != "" {
					assert.Equal(t, tt.expectedShort, shortURL)
				} else {
					assert.Len(t, shortURL, shortURLLength)
					assert.Regexp(t, "^[a-zA-Z0-9_]+$", shortURL)
				}
			}
		})
	}
}

func TestService_Get(t *testing.T) {
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
			originalURL, err := s.Get(tt.shortURL)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				assert.Empty(t, originalURL)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedURL, originalURL)
			}
		})
	}
}
