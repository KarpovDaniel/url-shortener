package main

import (
	"database/sql"
	"fmt"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"

	"url-shortener/internal/config"
	"url-shortener/internal/handler"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/memory"
	"url-shortener/internal/storage/postgres"
	"url-shortener/proto"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	var appStorage storage.Storage
	switch cfg.StorageType {
	case "postgres":
		log.Println("DB_HOST:", cfg.DBHost)
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
		defer db.Close() //nolint:errcheck

		if err := goose.Up(db, "migrations"); err != nil {
			log.Fatal("Failed to apply migrations:", err)
		}
		appStorage = postgres.NewPostgres(db)
	case "memory":
		appStorage = memory.NewMemory()
	default:
		log.Fatal("Unknown storage type")
	}

	// Создаём сервис, который реализует как HTTP, так и gRPC интерфейсы
	svc := service.NewService(appStorage)

	// Запуск gRPC-сервера в отдельной горутине
	go func() {
		grpcAddr := ":" + cfg.GRPCPort // добавьте поле GRPCPort в конфигурацию
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("Failed to listen on %s: %v", grpcAddr, err)
		}

		grpcServer := grpc.NewServer()
		proto.RegisterURLShortenerServer(grpcServer, svc)
		reflection.Register(grpcServer)
		log.Println("Starting gRPC server on", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Запуск HTTP-сервера
	h := handler.NewHandler(svc)
	r := h.SetupRoutes()

	log.Println("Starting HTTP server on port", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
