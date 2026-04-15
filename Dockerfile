# Используем официальный образ Go для сборки
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Копируем файлы модулей и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o meowcut ./cmd

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder
COPY --from=builder /app/meowcut .
# Копируем конфигурационные файлы
COPY configs ./configs
COPY migrations ./migrations

# Создаём директорию для логов
RUN mkdir -p logs

# Экспонируем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./meowcut"]