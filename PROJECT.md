# Project Plan: Go HTTP/gRPC + Postgres + ClickHouse + Mongo (+ Python gRPC)

---

# 1. Repository & Project Structure

- [x] **1.1** Создать структуру репозитория  
  - `go-service/`  
  - `py-service/`  
  - `proto/`

- [x] **1.2** Добавить базовый `README` с описанием архитектуры

- [x] **1.3** Описать роли БД:
  - Postgres (OLTP)
  - Mongo (raw events / документы)
  - ClickHouse (аналитика)

- [x] **1.4** Настроить gitignore

- [x] **1.5** Добавить Makefile

---

# 2. Docker Infrastructure

- [x] **2.1** Создать `docker-compose.yml`

- [x] **2.2** Добавить сервис `postgres`

- [x] **2.3** Добавить сервис `clickhouse`

- [ ] **2.4** Добавить сервис `mongo`

- [x] **2.5** Добавить сервис `go-service`

- [ ] **2.6** Добавить сервис `py-service`

- [x] **2.7** Добавить docker volumes

- [x] **2.8** Добавить docker networks

- [x] **2.9** Добавить healthchecks

- [ ] **2.10** Добавить `.env.example`

---

# 3. Protobuf Contracts

- [ ] **3.1** Создать `proto/common.proto`

- [ ] **3.2** Создать `proto/orders.proto`

- [ ] **3.3** Создать `proto/events.proto`

- [ ] **3.4** Настроить генерацию Go protobuf

- [ ] **3.5** Настроить генерацию Python protobuf

- [ ] **3.6** Настроить Makefile команду `make proto`

---

# 4. Go Service Skeleton

- [x] **4.1** Создать структуру Go проекта

- [x] **4.2** Добавить конфигурацию приложения

- [x] **4.3** Добавить structured logging

- [x] **4.4** Реализовать graceful shutdown

- [x] **4.5** Реализовать dependency wiring

- [x] **4.6** Middleware: request-id

- [x] **4.7** Middleware: panic recovery

- [x] **4.8** Middleware: logging

---

# 5. Postgres (Transactional Data)

## 5.1 Миграции

- [x] **5.1.1** Подключить инструмент миграций

- [ ] **5.1.2** Создать таблицу `users`

- [ ] **5.1.3** Создать таблицу `orders`

- [ ] **5.1.4** Создать таблицу `order_items`

- [ ] **5.1.5** Создать таблицу `outbox`

- [ ] **5.1.6** Добавить индексы

---

## 5.2 Репозитории

- [ ] **5.2.1** Реализовать Users repository

- [ ] **5.2.2** Реализовать Orders repository

- [ ] **5.2.3** Реализовать Outbox repository

- [ ] **5.2.4** Реализовать транзакционный слой

---

# 6. MongoDB (Raw Events)

- [ ] **6.1** Создать коллекцию `raw_events`

- [ ] **6.2** Добавить индекс `event_id`

- [ ] **6.3** Добавить индекс `order_id + created_at`

- [ ] **6.4** Добавить TTL индекс

- [ ] **6.5** Описать структуру документа события

---

# 7. ClickHouse (Analytics)

- [ ] **7.1** Создать таблицу `events_fact`

- [ ] **7.2** Добавить партиционирование

- [ ] **7.3** Добавить ORDER BY ключ

- [ ] **7.4** Добавить аналитические индексы

- [ ] **7.5** Реализовать repository для аналитики

---

# 8. Python gRPC Service (Mongo/Event Service)

## 8.1 Каркас сервиса

- [ ] **8.1.1** Создать gRPC сервер

- [ ] **8.1.2** Подключить Mongo

- [ ] **8.1.3** Добавить конфигурацию

- [ ] **8.1.4** Добавить structured logging

---

## 8.2 gRPC методы

- [ ] **8.2.1** Реализовать `StoreRawEvent`

- [ ] **8.2.2** Реализовать `GetOrderTimeline`

- [ ] **8.2.3** Реализовать `NormalizeEvent`

---

## 8.3 Тестирование

- [ ] **8.3.1** Unit tests

- [ ] **8.3.2** Integration tests

---

# 9. Go Business Logic

## 9.1 Orders Service

- [ ] **9.1.1** Use-case: CreateOrder

- [ ] **9.1.2** Use-case: UpdateOrderStatus

- [ ] **9.1.3** Use-case: GetOrder

- [ ] **9.1.4** Use-case: ListOrders

- [x] **9.1.5** Registration and logining

---

## 9.2 Outbox Worker

- [ ] **9.2.1** Реализовать polling outbox

- [ ] **9.2.2** Добавить batch processing

- [ ] **9.2.3** Вызов Python gRPC

- [ ] **9.2.4** Запись событий в ClickHouse

- [ ] **9.2.5** Retry policy

---

# 10. HTTP API

- [ ] **10.1** `POST /orders`

- [ ] **10.2** `POST /orders/{id}/status`

- [ ] **10.3** `GET /orders/{id}`

- [ ] **10.4** `GET /orders/{id}/timeline`

- [ ] **10.5** `GET /analytics/orders`

- [x] **10.6** `GET /healthcheck`

- [x] **10.7** Валидация входных данных

- [x] **10.8** Унифицированный формат ошибок

---

# 11. gRPC API (Go)

- [ ] **11.1** Реализовать `CreateOrder`

- [ ] **11.2** Реализовать `UpdateOrderStatus`

- [ ] **11.3** Реализовать `GetOrder`

- [ ] **11.4** Реализовать `GetOrderTimeline`

- [ ] **11.5** Реализовать `GetAnalytics`

- [ ] **11.6** gRPC interceptors (logging)

---

# 12. ClickHouse Analytics Queries

- [ ] **12.1** Orders per day

- [ ] **12.2** Funnel conversion

- [ ] **12.3** Average order value

- [ ] **12.4** Top users

---

# 13. Testing

- [ ] **13.1** Unit tests (Go)

- [ ] **13.2** Integration tests

- [ ] **13.3** gRPC contract tests

- [ ] **13.4** Load testing

---

# 14. CI/CD

- [ ] **14.1** Настроить CI pipeline

- [ ] **14.2** Запуск тестов

- [ ] **14.3** Запуск линтеров

- [ ] **14.4** Build docker images

---

# 15. Documentation

- [x] **15.1** README: запуск проекта

- [x] **15.2** README: архитектура

- [ ] **15.4** Примеры HTTP запросов

- [ ] **15.5** Примеры gRPC запросов

---

# 16. Optional Improvements

- [ ] **16.1** Swagger/OpenAPI

- [ ] **16.2** gRPC reflection

- [ ] **16.3** Dead-letter queue

- [ ] **16.4** Materialized views в ClickHouse

- [ ] **16.5** Версионирование событий

---