# avanpost-show

Сервис управления пользователями с REST API, PostgreSQL и NATS.

## Обзор

Обзор преимуществ testcontainers в проектах на Go/

## Быстрый старт

### 1. Запуск через Docker Compose

```bash
docker compose up --build
```

После запуска сервис будет доступен на:

- API: `http://localhost:5222`
- Swagger: `http://localhost:5222/swagger/index.html`

### 2. Локальный запуск

Скопируйте `config.yaml` в корень проекта и при необходимости скорректируйте параметры.

```bash
go mod download
go run ./cmd/main.go
```

Если используется альтернативный путь к файлу конфигурации:

```bash
API_CONFIG_PATH=./config.yaml go run ./cmd/main.go
```

## Документация

Swagger UI доступен по адресу:

```text
http://localhost:5222/swagger/index.html
```

## Миграции

SQL-скрипты находятся в `migration/postgres`.

Для применения миграций при запуске через Docker Compose используется сервис `migrate`.

## Тесты

Запустить все тесты можно командой:

```bash
go test ./...
```

## Структура проекта

- `cmd/main.go` — точка входа приложения
- `internal/app` — инициализация приложения, маршрутизация и middleware
- `internal/config` — загрузка конфигурации
- `internal/user/delivery/http` — HTTP-контроллеры
- `internal/user/usecase` — бизнес-логика
- `internal/user/repository/postgres` — работа с PostgreSQL
- `pkg/publisher` — публикация событий в NATS
- `migration/postgres` — SQL-миграции
