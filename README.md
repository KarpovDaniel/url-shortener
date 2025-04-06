Структура проекта:

```
.
├── cmd
│   └── url-shortener
│       └── main.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── handler
│   │   ├── handler_test.go
│   │   └── handler.go
│   ├── storage
│   │   ├── memory
│   │   │   ├── memory.go
│   │   │   └── memory_test.go
│   │   ├── postgres
│   │   │   ├── postgres.go
│   │   │   └── postgres.go
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

# Клонируйте проект

```
git clone https://github.com/KarpovDaniel/url-shortener.git
```

# Запуск сервера:

## Запуск тестов:
```
make tests
```

## memory: 
```
make compose-memory-up
```

## postgres: 
```
make compose-postgres-up
```

## Остановка: 
```
make down
```

# Примеры запросов:

## HTTP API:

POST:

```
curl -X POST -H "Content-Type: application/json" -d '{"original_url": "https://example.com"}' http://localhost:8080
```

Пример ответа:

```
{"short_url":"knFxGJkkDb"}
```

GET:

```
curl http://localhost:8080/knFxGJkkDb
```

Пример ответа:

```
{"original_url":"https://example.comm"}
```

Пример ответа при Invalid URL:

```
URL not found
```
