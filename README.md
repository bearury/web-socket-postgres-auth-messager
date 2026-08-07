# Запуск проекта

## Локально (Go + Postgres в Docker)

База данных поднимается в Docker, само приложение — локально через `go run`.

```bash
# 1. Поднять Postgres в Docker (один терминал)
docker compose -f docker-compose.db.yml up -d

# 2. Убедиться, что .env задан
# echo "DB_PASSWORD=qwerty" > .env

# 3. Запустить приложение
go run ./cmd/main.go
```

Сервер будет на `http://localhost:8080`.

Чтобы остановить БД:

```bash
docker compose -f docker-compose.db.yml down -v
```

## Полностью в Docker

```bash
docker compose up --build -d
```

Приложение: `http://localhost:8000`  
Postgres: `localhost:5436`

Если порт `8000` занят:

```bash
APP_PORT=8001 docker compose up --build -d
```

## Проверка

```bash
curl -X POST http://localhost:8080/auth/sign-up \
  -H "Content-Type: application/json" \
  -d '{"name":"test","username":"test","password":"test"}'
```

Должен вернуться `200`.

> **Примечание:** `service/auth.go` и `handler/auth.go` сейчас заглушки. Приложение запускается и отвечает `200`, но логика регистрации/авторизации не реализована.
