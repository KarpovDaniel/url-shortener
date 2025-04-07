# Структура проекта:

```
.
├── cmd
│   └── url-shortener
│       └── main.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── handler
│   │   ├── handler.go
│   │   └── handler_test.go
│   ├── storage
│   │   ├── memory
│   │   │   ├── memory.go
│   │   │   └── memory_test.go
│   │   ├── postgres
│   │   │   ├── postgres.go
│   │   │   └── postgres_test.go
│   │   └── storage.go
│   └── service
│       ├── service.go
│       └── service_test.go
├── migrations
│   └── 00001_create_urls_table.sql
├── .env
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

# Запуск сервера:

## Запуск тестов:
```
make tests
```

## Покрытие тестами
```
make tests-coverage
```

## memory: 
```
make memory
```

## postgres: 
```
make postgres
```

## Остановка: 
```
make down
```

# Примеры запросов:

## gRPC:

CreateLink:

```
grpcurl -plaintext -d '{"original_url": "https://example.com"}' localhost:50051 proto.URLShortener/CreateURL
```

Пример ответа:

```
{
  "shortUrl": "_shortURL_"
}
```

GetLink :

```
grpcurl -plaintext -d '{"short_url": "_shortURL_"}' localhost:50051 proto.URLShortener/GetURL
```

Пример ответа:

```
{
  "originalUrl": "https://example.com"
}
```

Пример ответа когда не надено:

```
{
  "error": "short URL not found"
}
```

## HTTP API:

POST:

```
curl -X POST -d "url=https://example.com" http://localhost:8080
```

Пример ответа:

```
_shortURL_
```

GET:

```
curl http://localhost:8080/_shortURL_
```

Пример ответа:

```
https://example.com
```

Пример ответа при Invalid URL:

```
URL not found
```
