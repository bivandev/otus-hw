# Используем официальный образ Golang
FROM golang:1.22 as builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем модульные файлы и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Устанавливаем Ginkgo CLI
RUN go install github.com/onsi/ginkgo/v2/ginkgo@latest

# Сборка интеграционных тестов
RUN ginkgo build ./integration_tests

# Используем минимальный образ для запуска тестов
FROM debian:bookworm-slim

# Устанавливаем необходимые зависимости
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Устанавливаем рабочую директорию
WORKDIR /tests

# Копируем собранные тесты из builder-этапа
COPY --from=builder /app/integration_tests/integration_tests.test .

# Устанавливаем переменные окружения для тестов
ENV TEST_ENV=true

# Команда по умолчанию для запуска тестов
ENTRYPOINT ["./integration_tests.test"]
