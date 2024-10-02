# Имя бинарного файла
BINARY_NAME=statements

# Папка с исходным кодом (main.go находится в папке cmd)
SRC_DIR=./cmd

# Папка для сборки
BUILD_DIR=./bin

# Путь к конфигурационному файлу
CONFIG_FILE=config.yaml

# Переменные окружения для приложения
PORT ?= 9080
GIN_MODE ?= debug

# Флаги для сборки
LDFLAGS_DEBUG=-ldflags="-s -w"
LDFLAGS_RELEASE=-ldflags="-s -w -X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(shell date)'"

# Текущая версия (можно изменять в зависимости от версии релиза)
VERSION=1.0.0

# Цель по умолчанию
.PHONY: all
all: build

# Проверка зависимостей и создание необходимых директорий
.PHONY: prepare
prepare:
	@echo "==> Подготовка окружения..."
	@mkdir -p $(BUILD_DIR)

# Сборка приложения в режиме Debug
.PHONY: build-debug
build-debug: prepare
	@echo "==> Сборка приложения в режиме DEBUG..."
	go build $(LDFLAGS_DEBUG) -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)/main.go

# Сборка приложения в режиме Release с оптимизациями
.PHONY: build-release
build-release: prepare
	@echo "==> Сборка приложения в режиме RELEASE..."
	go build $(LDFLAGS_RELEASE) -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)/main.go

# Полная сборка для Debug и Release
.PHONY: build
build: build-debug build-release

# Запуск приложения в режиме Debug
.PHONY: run-debug
run-debug: build-debug
	@echo "==> Запуск приложения в режиме DEBUG на порту $(PORT)..."
	GIN_MODE=debug $(BUILD_DIR)/$(BINARY_NAME) --config $(CONFIG_FILE)

# Запуск приложения в режиме Release
.PHONY: run-release
run-release: build-release
	@echo "==> Запуск приложения в режиме RELEASE на порту $(PORT)..."
	GIN_MODE=release $(BUILD_DIR)/$(BINARY_NAME) --config $(CONFIG_FILE)

# Линтинг кода
.PHONY: lint
lint:
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "golangci-lint не установлен. Пожалуйста, установите его командой: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi
	@echo "==> Линтинг кода..."
	golangci-lint run

# Запуск тестов
.PHONY: test
test:
	@echo "==> Запуск тестов..."
	go test ./... -v

# Очистка собранных файлов
.PHONY: clean
clean:
	@echo "==> Очистка сборки..."
	rm -rf $(BUILD_DIR)

# Установка зависимостей
.PHONY: deps
deps:
	@echo "==> Установка зависимостей..."
	go mod tidy

# Обновление зависимостей
.PHONY: update-deps
update-deps:
	@echo "==> Обновление зависимостей..."
	go get -u ./...

# Применение миграций базы данных
.PHONY: migrate
migrate:
	@echo "==> Применение миграций базы данных..."
	go run $(SRC_DIR)/main.go migrate --config $(CONFIG_FILE)

# Генерация Docker-образа для Debug
.PHONY: docker-build-debug
docker-build-debug:
	@echo "==> Генерация Docker-образа для DEBUG..."
	docker build --target=debug -t $(BINARY_NAME):debug .

# Генерация Docker-образа для Release
.PHONY: docker-build-release
docker-build-release:
	@echo "==> Генерация Docker-образа для RELEASE..."
	docker build -f ./docker/Dockerfile --target=release -t $(BINARY_NAME):release .

# Запуск приложения через Docker в режиме Debug
.PHONY: docker-run-debug
docker-run-debug: docker-build-debug
	@echo "==> Запуск Docker-контейнера в режиме DEBUG..."
	docker run -p $(PORT):$(PORT) --env GIN_MODE=debug $(BINARY_NAME):debug

# Запуск приложения через Docker в режиме Release
.PHONY: docker-run-release
docker-run-release: docker-build-release
	@echo "==> Запуск Docker-контейнера в режиме RELEASE..."
	docker run -p $(PORT):$(PORT) --env GIN_MODE=release $(BINARY_NAME):release

# Генерация документации
.PHONY: docs
docs:
	@echo "==> Генерация документации..."
	godoc -http=:6060

# Очистка и повторная установка зависимостей
.PHONY: reinstall-deps
reinstall-deps:
	@echo "==> Очистка и переустановка зависимостей..."
	rm -rf vendor/
	go mod vendor

# Полная сборка, установка зависимостей и запуск
.PHONY: full-build
full-build: deps build-release run-release