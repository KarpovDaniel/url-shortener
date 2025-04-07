package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config содержит конфигурационные параметры приложения
type Config struct {
	StorageType string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	ServerPort  string
}

// LoadConfig загружает конфигурацию из .env файла и переменных окружения
func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{
		StorageType: os.Getenv("STORAGE_TYPE"),
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		ServerPort:  os.Getenv("SERVER_PORT"),
	}, nil
}
