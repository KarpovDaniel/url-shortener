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

# Клонируйте проект

```
git clone https://github.com/KarpovDaniel/url-shortener.git
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
