## Запуск тестов
tests:
	# test ./...

## Запуск с postgres-хранилищем
compose-postgres-up:
	docker-compose --profile postgres up --build

## Запуск с memory-хранилищем
compose-memory-up:
	docker-compose --profile memory up --build

## Остановка Docker Compose
compose-down:
	docker-compose --profile postgres down
	docker-compose --profile memory down
