# Этап сборки
FROM golang:1.23.1-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем модульные файлы и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

RUN go build -o file_upload_service ./cmd/main.go

# Этап выполнения
FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk add --no-cache ca-certificates

# Создаём пользователя для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем скомпилированный бинарный файл из этапа сборки
COPY --from=builder /app/file_upload_service .

# Копируем файлы окружения и публичный ключ
COPY .env .
COPY public_key.pem .

RUN chown -R appuser:appgroup /app

# Переключаемся на созданного пользователя
USER appuser

EXPOSE 8081

# Команда для запуска приложения
CMD ["./file_upload_service", "--config=.env"]
