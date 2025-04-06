package postgres

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"url-shortener/internal/storage"
)

// convertArgs преобразует []interface{} в []driver.Value
func convertArgs(args []interface{}) []driver.Value {
	driverArgs := make([]driver.Value, len(args))
	for i, arg := range args {
		driverArgs[i] = arg
	}
	return driverArgs
}

func TestPostgres_Save(t *testing.T) {
	tests := []struct {
		name          string
		shortURL      string
		originalURL   string
		setup         func(sqlmock.Sqlmock)
		expectedShort string
		expectedErr   error
	}{
		{
			name:        "Успешное сохранение нового URL",
			shortURL:    "abc123",
			originalURL: "https://example.com",
			setup: func(mock sqlmock.Sqlmock) {
				query, args, _ := squirrel.Insert("urls").
					Columns("short_url", "original_url").
					Values("abc123", "https://example.com").ToSql()
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(convertArgs(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedShort: "abc123",
			expectedErr:   nil,
		},
		{
			name:        "Повторное использование существующего originalURL",
			shortURL:    "xyz789",
			originalURL: "https://example.com",
			setup: func(mock sqlmock.Sqlmock) {
				query, args, _ := squirrel.Insert("urls").
					Columns("short_url", "original_url").
					Values("xyz789", "https://example.com").ToSql()
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(convertArgs(args)...).
					WillReturnError(errors.New("urls_original_url_key"))

				selectQuery := squirrel.Select("short_url").
					From("urls").
					Where(squirrel.Eq{"original_url": "https://example.com"})
				selectSQL, selectArgs, _ := selectQuery.ToSql()
				mock.ExpectQuery(selectSQL).
					WithArgs(convertArgs(selectArgs)...).
					WillReturnRows(sqlmock.NewRows([]string{"short_url"}).AddRow("abc123"))
			},
			expectedShort: "abc123",
			expectedErr:   nil,
		},
		{
			name:        "Ошибка при сохранении",
			shortURL:    "def456",
			originalURL: "https://newexample.com",
			setup: func(mock sqlmock.Sqlmock) {
				query, args, _ := squirrel.Insert("urls").
					Columns("short_url", "original_url").
					Values("def456", "https://newexample.com").ToSql()
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(convertArgs(args)...).
					WillReturnError(errors.New("database error"))
			},
			expectedShort: "",
			expectedErr:   errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close() //nolint:errcheck

			tt.setup(mock)

			pg := NewPostgres(db)
			shortURL, err := pg.Save(tt.shortURL, tt.originalURL)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				assert.Empty(t, shortURL)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedShort, shortURL)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgres_Get(t *testing.T) {
	tests := []struct {
		name        string
		shortURL    string
		setup       func(sqlmock.Sqlmock)
		expectedURL string
		expectedErr error
	}{
		{
			name:     "Успешное получение URL",
			shortURL: "abc123",
			setup: func(mock sqlmock.Sqlmock) {
				query, args, _ := squirrel.Select("original_url").
					From("urls").
					Where(squirrel.Eq{"short_url": "abc123"}).ToSql()
				mock.ExpectQuery(query).
					WithArgs(convertArgs(args)...).
					WillReturnRows(sqlmock.NewRows([]string{"original_url"}).AddRow("https://example.com"))
			},
			expectedURL: "https://example.com",
			expectedErr: nil,
		},
		{
			name:     "URL не найден",
			shortURL: "xyz789",
			setup: func(mock sqlmock.Sqlmock) {
				query, args, _ := squirrel.Select("original_url").
					From("urls").
					Where(squirrel.Eq{"short_url": "xyz789"}).ToSql()
				mock.ExpectQuery(query).
					WithArgs(convertArgs(args)...).
					WillReturnError(sql.ErrNoRows)
			},
			expectedURL: "",
			expectedErr: storage.ErrNotFound,
		},
		{
			name:     "Ошибка базы данных",
			shortURL: "def456",
			setup: func(mock sqlmock.Sqlmock) {
				query, args, _ := squirrel.Select("original_url").
					From("urls").
					Where(squirrel.Eq{"short_url": "def456"}).ToSql()
				mock.ExpectQuery(query).
					WithArgs(convertArgs(args)...).
					WillReturnError(errors.New("database error"))
			},
			expectedURL: "",
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close() //nolint:errcheck

			tt.setup(mock)

			pg := NewPostgres(db)
			originalURL, err := pg.Get(tt.shortURL)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				assert.Empty(t, originalURL)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedURL, originalURL)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
