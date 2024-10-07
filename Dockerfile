# Этап 1: Сборка приложения
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o statements ./cmd/main.go

# Этап 2: Минимальный образ для запуска
FROM alpine:latest
WORKDIR /root/

# Установка Python3 и pip
RUN apk add --no-cache python3 py3-pip

# Создание виртуального окружения Python
RUN python3 -m venv /app/venv

# Копирование проекта
COPY --from=builder /app/statements .
COPY config.yaml .
COPY ./assets /app/assets
COPY ./migrations /app/migrations
COPY ./scripts /app/scripts
COPY requirements.txt /app/requirements.txt

# Установка Python зависимостей
RUN /app/venv/bin/pip install -r /app/requirements.txt

# Экспонирование порта
EXPOSE 9080

# Запуск приложения
CMD ["./statements", "--config", "config.yaml"]