# Go / Python Services — Monorepo

Короткое описание: сервисная платформа на Go (REST + PostgreSQL) с будущим Python gRPC-компонентом. Подробный план развития — в [PROJECT.md](PROJECT.md).

## Структура репозитория

- [docker-compose.yaml](docker-compose.yaml) — описание окружения контейнеров (Postgres, go-service).  
- [Dockerfile](Dockerfile) — образ сборки сервиса Go.
- [.env](.env) — окружение (файл в .gitignore).
- certs/ — параметры и конфиги сертификатов: [openssl.cnf](certs/openssl.cnf), [cert.conf](certs/cert.conf), [.srl](certs/.srl).

- go-service/ — основной Go-сервис:
  - [go-service/go.mod](go-service/go.mod) — зависимости модуля.
  - [cmd/app/main.go](go-service/cmd/app/main.go) — точка входа, wiring сервера и graceful shutdown (`main.main`) — [go-service/cmd/app/main.go](go-service/cmd/app/main.go).
  - [cmd/internal/http/server.go](go-service/cmd/internal/http/server.go) — сборка HTTP сервера (`http.GetHttpServer`, `http.ServeServer`) — [`http.GetHttpServer`](go-service/cmd/internal/http/server.go), [`http.ServeServer`](go-service/cmd/internal/http/server.go).
  - [cmd/internal/router/router.go](go-service/cmd/internal/router/router.go) — маршрутизация и регистрация публичных/API-эндпоинтов (`router.GetRouter`) — [`router.GetRouter`](go-service/cmd/internal/router/router.go).
  - [cmd/internal/middleware/]
    - [db.go](go-service/cmd/internal/middleware/db.go) — middleware для подключения к БД и транзакций (`middleware.DbMiddleware`) — [`middleware.DbMiddleware`](go-service/cmd/internal/middleware/db.go).
    - [auth.go](go-service/cmd/internal/middleware/auth.go) — JWT-аутентификация и генерация токенов (`middleware.AuthMiddleware`, `middleware.GetToken`) — [`middleware.AuthMiddleware`](go-service/cmd/internal/middleware/auth.go), [`middleware.GetToken`](go-service/cmd/internal/middleware/auth.go).
  - [cmd/internal/public/public.go](go-service/cmd/internal/public/public.go) — публичные endpoint'ы: `/health`, `/register`, `/login` (`public.GetPublicHandler`) — [`public.GetPublicHandler`](go-service/cmd/internal/public/public.go).
  - [cmd/internal/api/api.go](go-service/cmd/internal/api/api.go) — приватный API под префиксом `/api` (`api.GetApiHandler`) — [`api.GetApiHandler`](go-service/cmd/internal/api/api.go).
  - [cmd/internal/errorHandler/errorHandler.go](go-service/cmd/internal/errorHandler/errorHandler.go) — обёртка ошибок для handler’ов (`errorHandler.Wrap`) — [`errorHandler.Wrap`](go-service/cmd/internal/errorHandler/errorHandler.go).
  - [cmd/internal/errors/erros.go](go-service/cmd/internal/errors/erros.go) — типы и фабрики HTTP-ошибок (`errors.HTTPError`, `errors.BadRequest`, `errors.InternalServerError`, ... ) — [`errors.HTTPError`](go-service/cmd/internal/errors/erros.go).

- proto/ — protobuf контракты (планы в [PROJECT.md](PROJECT.md)).

- py-service/ — задел для Python gRPC сервиса (описание задач в [PROJECT.md](PROJECT.md)).

## Что планируется (кратко, из PROJECT.md)
Ссылка на полный план: [PROJECT.md](PROJECT.md). Ключевые направления:
- Docker: добавить ClickHouse, Mongo, доработать docker-compose.
- Protobuf и генерация для Go/Python.
- Полный каркас Go-сервиса: миграции Postgres, репозитории, outbox, бизнес-логика orders.
- Python gRPC для хранения raw events в Mongo.
- Observability (metrics, tracing), тестирование и CI/CD.

## Быстрый старт (локально)
1. Подготовить .env (см. [docker-compose.yaml](docker-compose.yaml)).
2. Запустить Postgres через docker-compose:  
   docker-compose up -d --build
3. Сконфигурировать TLS файлы, задать переменную окружения TLS_PATH и запустить сервис (локально или через Dockerfile).
