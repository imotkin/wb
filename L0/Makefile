.PHONY: all build run up down test cover lint tidy help

all: lint test build

build:
	mkdir -p bin
	go build -v -o app ./cmd/main.go

run:
	go run ./cmd/main.go

up:
	docker-compose up -d --build

down:
	docker-compose down

test:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

cover: test
	go tool cover -html=coverage.out

lint:
	@golangci-lint run --timeout=5m

tidy:
	go mod tidy
	go mod verify

help:
	@echo "Доступные команды:"
	@echo "  make build       - сборка проекта"
	@echo "  make run         - лольканый запуск проекта"
	@echo "  make up          - запуск проекта в Docker"
	@echo "  make down        - остановка проекта в Docker"
	@echo "  make test        - запуск тестов"
	@echo "  make cover       - запуск покрытия тестами в HTML"
	@echo "  make lint        - запустить линтера golangci-lint"
	@echo "  make tidy        - запуск проверки зависимостей"
	@echo "  make all         - проверка линтера, запуск тестов и сборка проекта"
	@echo "  make help        - вывод списка доступных команд"