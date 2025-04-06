package memory

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Тест для метода Save
func TestMemory_Save(t *testing.T) {
	tests := []struct {
		name          string
		shortURL      string
		originalURL   string
		setup         func(*Memory)
		expectedShort string
		expectedErr   error
	}{
		{
			name:          "Успешное сохранение нового URL",
			shortURL:      "abc123",
			originalURL:   "https://example.com",
			setup:         func(m *Memory) {},
			expectedShort: "abc123",
			expectedErr:   nil,
		},
		{
			name:        "Повторное использование существующего originalURL",
			shortURL:    "xyz789",
			originalURL: "https://example.com",
			setup: func(m *Memory) {
				m.Save("abc123", "https://example.com") //nolint:errcheck
			},
			expectedShort: "abc123",
			expectedErr:   nil,
		},
		{
			name:        "Ошибка из-за дублирования shortURL",
			shortURL:    "abc123",
			originalURL: "https://newexample.com",
			setup: func(m *Memory) {
				m.Save("abc123", "https://example.com") //nolint:errcheck
			},
			expectedShort: "",
			expectedErr:   errors.New("short URL already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewMemory()
			tt.setup(mem)

			shortURL, err := mem.Save(tt.shortURL, tt.originalURL)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				assert.Empty(t, shortURL)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedShort, shortURL)

				// Проверяем, что данные действительно сохранены
				originalURL, getErr := mem.Get(shortURL)
				assert.NoError(t, getErr)
				assert.Equal(t, tt.originalURL, originalURL)
			}
		})
	}
}

// Тест для метода Get
func TestMemory_Get(t *testing.T) {
	tests := []struct {
		name        string
		shortURL    string
		setup       func(*Memory)
		expectedURL string
		expectedErr error
	}{
		{
			name:     "Успешное получение URL",
			shortURL: "abc123",
			setup: func(m *Memory) {
				m.Save("abc123", "https://example.com") //nolint:errcheck
			},
			expectedURL: "https://example.com",
			expectedErr: nil,
		},
		{
			name:        "URL не найден",
			shortURL:    "xyz789",
			setup:       func(m *Memory) {},
			expectedURL: "",
			expectedErr: errors.New("short URL not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewMemory()
			tt.setup(mem)

			originalURL, err := mem.Get(tt.shortURL)

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
