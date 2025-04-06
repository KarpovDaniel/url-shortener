-- +goose Up
CREATE TABLE urls (
                      short_url VARCHAR(10) PRIMARY KEY,
                      original_url TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE urls;