# Этап 1: Сборка Go-приложения
FROM golang:1.23-alpine AS builder

# Установка рабочей директории
WORKDIR /app

# Кэшируем зависимости Go
COPY go.mod go.sum ./
RUN go mod tidy

# Копируем исходный код Go после установки зависимостей, чтобы избежать пересборки на изменениях в коде
COPY . .

# Сборка бинарного файла Go
RUN go build -o statements ./cmd/main.go

# Этап 2: Минимальный образ для запуска
FROM alpine:latest

# Установка необходимых пакетов для Python
RUN apk add --no-cache python3 py3-pip

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем Go-бинарник из builder этапа
COPY --from=builder /app/statements .

# Копируем файлы конфигурации и необходимые ресурсы
COPY config.yaml .
COPY ./assets /app/assets
COPY ./migrations /app/migrations
COPY ./scripts /app/scripts

# Копируем Python зависимости
COPY requirements.txt /app/requirements.txt

# Создание и настройка виртуального окружения Python
RUN python3 -m venv /app/venv \
    && /app/venv/bin/pip install --upgrade pip \
    && /app/venv/bin/pip install -r /app/requirements.txt

# Экспонируем порт для приложения
EXPOSE 9080

# Запускаем приложение
CMD ["./statements", "--config", "config.yaml"]