# Команда по умолчанию: сборка и запуск
.DEFAULT_GOAL := run

# Команды для сборки и запуска контейнеров с использованием Docker Compose

build and run:
	@echo "Building services..."
	docker-compose up --build

build:
	@echo "Building services..."
	docker-compose build

run:
	@echo "Running services..."
	docker-compose up

clean:
	@echo "Cleaning up..."
	docker-compose down

.PHONY: build run clean help
