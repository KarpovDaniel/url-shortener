## Запуск тестов
tests:
	go test -v ./...

## Покрытие тестами
tests-coverage:
	go test -cover ./...

## Запуск с postgres-хранилищем
postgres:
	docker-compose --profile postgres up --build

## Запуск с memory-хранилищем
memory:
	docker-compose --profile memory up --build

## Остановка Docker Compose
down:
	docker-compose --profile postgres down
	docker-compose --profile memory down
