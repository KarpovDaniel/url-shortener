package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"url-shortener/internal/storage"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"url-shortener/internal/config"
	"url-shortener/internal/handler"
	"url-shortener/internal/service"
	"url-shortener/internal/storage/memory"
	"url-shortener/internal/storage/postgres"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	var app_storage storage.Storage
	if cfg.StorageType == "postgres" {
		log.Println("DB_HOST:", cfg.DBHost)
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		// Применение миграций
		if err := goose.Up(db, "migrations"); err != nil {
			log.Fatal("Failed to apply migrations:", err)
		}
		app_storage = postgres.NewPostgresStorage(db)
	} else {
		app_storage = memory.NewMemoryStorage()
	}

	svc := service.NewService(app_storage)
	h := handler.NewHandler(svc)
	r := h.SetupRoutes()

	log.Println("Starting server on port", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
