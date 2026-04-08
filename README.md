# Сервис бронирования переговорок

## Технологии

| Компонент | Версия / пакет |
|-----------|----------------|
| **Go** | 1.25.3 |
| **HTTP** | [gin-gonic/gin](https://github.com/gin-gonic/gin) |
| **JWT** | [golang-jwt/jwt/v4](https://github.com/golang-jwt/jwt) |
| **UUID** | [google/uuid](https://github.com/google/uuid) |
| **PostgreSQL (драйвер)** | [jackc/pgx/v5](https://github.com/jackc/pgx) |
| **Конфиг из env** | [ilyakaznacheev/cleanenv](https://github.com/ilyakaznacheev/cleanenv) |
| **Миграции** | [golang-migrate/migrate](https://github.com/golang-migrate/migrate) |
| **Пароли** | `golang.org/x/crypto` (bcrypt) |
| **Тесты** | [stretchr/testify](https://github.com/stretchr/testify) |

## Требования

- Docker и Docker Compose
- Go 1.25+

## Быстрый старт

1. Скопируйте переменные окружения:

   ```bash
   cp .env.example .env
   ```

2. Поднимите сервисы:

   ```bash
   docker compose up -d --build
   ```

3. API по умолчанию: `http://127.0.0.1:8080`  
   Проверка: `GET /_info` → `{"status":"ok"}`

Переменные окружения (см. `.env.example`): `JWT_SECRET`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`. В Docker `DB_HOST` обычно `postgres`.

## Основные эндпоинты

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/dummyLogin` | Тестовый JWT по роли (`admin` / `user`) |
| POST | `/register`, `/login` | Регистрация и вход |
| GET | `/rooms/list`, POST `/rooms/create` | Список / создание переговорки (create — admin) |
| POST | `/rooms/:roomId/schedule/create` | Расписание (admin) |
| GET | `/rooms/:roomId/slots/list?date=YYYY-MM-DD` | Доступные слоты на дату |
| POST | `/bookings/create`, GET `/bookings/my`, … | Бронирования |

Защищённые маршруты требуют заголовок `Authorization: Bearer <token>`.

## Тесты

**Юнит-тесты** (без поднятого API):

```bash
go test ./... -count=1
```

**Интеграционные E2E**:

```bash
go test -tags=integration ./e2e/...
```

Переменная `E2E_BASE_URL` (по умолчанию `http://127.0.0.1:8080`).

## Нагрузочное тестирование (wrk)

Пример прогона с Lua-скриптом `slots.lua` (запросы к нужному URL с JWT и параметрами слотов), базовый URL сервера:

```bash
wrk -t4 -c100 -d30s --latency -s slots.lua http://127.0.0.1:8080
```
`slots.lua`:
```bash
wrk.method = "GET"

local token = "<Bearer token>"

request = function()
  wrk.headers["Authorization"] = "Bearer " .. token
  return wrk.format(
    nil,
    "/rooms/<roomID>/slots/list?date=yyyy-MM-dd"
  )
end
```

### Результаты нагрузки самого высоконагруженного endpoint

Условия: 4 потока, 100 соединений, 30 с, сценарий из `slots.lua`.
Endpoint: `/rooms/<roomID>/slots/list?date=yyyy-MM-dd`

```
Running 30s test @ http://127.0.0.1:8080
  4 threads and 100 connections

  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    15.54ms    3.43ms 119.12ms   81.48%
    Req/Sec     1.61k   238.49     3.59k    79.92%

  Latency Distribution
     50%   15.04ms
     75%   16.91ms
     90%   19.16ms
     99%   25.62ms

  193442 requests in 32.85s, 578.35MB read
  Socket errors: connect 0, read 0, write 0, timeout 200

Requests/sec:   5888.91
Transfer/sec:     17.61MB
```
**Итоговые показатели после нагрузки**

- **Пропускная способность:** ~5900 запросов в секунду (RPS) при 4 потоках и 100 соединениях.
- **Задержка:** p50 ≈ 15 ms, p90 ≈ 19 ms, p99 ≈ 26 ms.
- **Успешность:** ~99.9% запросов без ошибок (таймауты ≈ 0.1%).
