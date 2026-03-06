# Copilot Instructions

## Project Overview

This repository is a **CRUD News API** built with **Go Fiber**, following **Hexagonal Architecture (Ports & Adapters)**, backed by **PostgreSQL** for persistent storage and **Redis** for caching, orchestrated with **Docker Compose**.

---

## Tech Stack

| Layer         | Technology                  |
|---------------|-----------------------------|
| Language      | Go (1.25+)                  |
| Web Framework | [Fiber v2](https://github.com/gofiber/fiber) |
| Database      | PostgreSQL 16               |
| Cache         | Redis 7                     |
| ORM           | GORM                        |
| Hot Reload    | Air                         |
| Container     | Docker & Docker Compose     |

---

## Architecture — Hexagonal (Ports & Adapters)

```
┌─────────────────────────────────────────────┐
│               Primary Adapters              │
│         (Driving — HTTP / Fiber)            │
│   internal/adapters/http/                   │
│     handler/news_handler.go                 │
│     router/router.go                        │
└──────────────────┬──────────────────────────┘
                   │ uses input port
┌──────────────────▼──────────────────────────┐
│           Application (Use Cases)           │
│  internal/app/                              │
│    port/news_service.go    ← input port     │
│    port/news_repository.go ← output port    │
│    port/news_cache.go      ← output port    │
│    news_service.go         ← implementation │
│                                             │
│  internal/domain/                           │
│    news.go                 ← entity         │
└──────────────────┬──────────────────────────┘
                   │ uses output ports
┌──────────────────▼──────────────────────────┐
│              Secondary Adapters             │
│        (Driven — Postgres / Redis)          │
│   internal/adapters/                        │
│     postgres/news_repository.go             │
│     redis/news_cache.go                     │
└─────────────────────────────────────────────┘
```

---

## Project Structure

```
revise-redis/
├── cmd/
│   └── main.go                        # Composition root — wires all layers
├── config/
│   └── config.go                      # Env/config loader
├── infrastructure/
│   ├── db/postgres.go                 # GORM + PostgreSQL connection
│   └── cache/redis.go                 # go-redis client setup
├── internal/
│   ├── domain/
│   │   └── news.go                    # Pure domain entity (no framework deps)
│   ├── app/
│   │   ├── port/
│   │   │   ├── news_service.go        # Input port interface
│   │   │   ├── news_repository.go     # Output port interface (persistence)
│   │   │   └── news_cache.go          # Output port interface (cache)
│   │   └── news_service.go            # Use case implementation
│   └── adapters/
│       ├── http/
│       │   ├── handler/
│       │   │   └── news_handler.go    # Fiber HTTP handlers (driving)
│       │   └── router/
│       │       └── router.go          # Route registration
│       ├── postgres/
│       │   └── news_repository.go     # GORM implementation of output port
│       └── redis/
│           └── news_cache.go          # Redis implementation of output port
├── docker-compose.yml
├── Dockerfile
├── .air.toml
├── .env
└── go.mod / go.sum
```

---

## Domain Entity

```go
// internal/core/domain/news.go
type News struct {
    ID        uint      `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Author    string    `json:"author"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

The domain entity has **no GORM tags** and **no framework dependencies**. The Postgres adapter maps it to a GORM model internally.

---

## Ports

### Input port
```go
// internal/app/port/news_service.go
type NewsService interface {
    GetAll() ([]domain.News, error)
    GetByID(id uint) (*domain.News, error)
    Create(news *domain.News) error
    Update(id uint, news *domain.News) error
    Delete(id uint) error
}
```

### Output ports
```go
// internal/app/port/news_repository.go
type NewsRepository interface { ... }

// internal/app/port/news_cache.go
type NewsCache interface { ... }
```

---

## API Endpoints

| Method | Path            | Description           | Cached |
|--------|-----------------|-----------------------|--------|
| GET    | `/api/news`     | List all news         | Yes    |
| GET    | `/api/news/:id` | Get single news by ID | Yes    |
| POST   | `/api/news`     | Create a news article | No     |
| PUT    | `/api/news/:id` | Update a news article | No     |
| DELETE | `/api/news/:id` | Delete a news article | No     |

---

## Caching Strategy

- Cache key format: `news:<id>` (single), `news:all` (list)
- TTL: **5 minutes** by default (configurable via `REDIS_TTL`)
- Cache invalidated on `POST`, `PUT`, `DELETE`
- Cache logic lives **only in the core service** — never in adapters or handlers
- Redis values stored as **JSON strings**

---

## Response Format

Success:
```json
{ "success": true, "data": ..., "message": "..." }
```
Error:
```json
{ "success": false, "message": "error description" }
```

---

## Environment Variables (`.env`)

```env
APP_PORT=3000

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=newsdb
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_TTL=300
```

---

## Docker Compose Services

| Service       | Image                     | Port |
|---------------|---------------------------|------|
| app           | built from Dockerfile     | 3000 |
| postgres      | postgres:16-alpine        | 5432 |
| redis         | redis:7-alpine            | 6379 |
| redisinsight  | redis/redisinsight:latest | 5540 |

- **app** mounts source as volume — Air hot reload works inside container
- **postgres** data persisted via named volume
- **redis** runs cache-only (no persistence)
- **redisinsight** available at http://localhost:5540 (connect to host: `redis`, port: `6379`)

---

## Coding Conventions

- Domain entity (`internal/domain/`) must have **zero** framework or driver imports
- Input port interface lives in `internal/app/port/`
- Output port interfaces live in `internal/app/port/` (same package)
- Core service (`internal/app/`) depends **only on port interfaces**, never on concrete adapters
- Adapters import the domain but the domain never imports adapters
- GORM models are private and live **inside** the postgres adapter package
- Cache read/write/invalidation logic lives **only** in the core service (`internal/app/`)
- All handlers return `{ "success": bool, "data": ..., "message": "..." }`
- Use `context.Background()` in adapter methods for Redis calls
- Composition root (`cmd/main.go`) is the only place that imports all layers

---

## Key Dependencies (`go.mod`)

```
github.com/gofiber/fiber/v2
github.com/redis/go-redis/v9
gorm.io/gorm
gorm.io/driver/postgres
github.com/joho/godotenv
```

---

## Running the Project

```bash
# Start all services (with hot reload via Air)
docker compose up --build

# Run locally (requires external Postgres & Redis)
air
```

---

## Tech Stack

| Layer         | Technology                  |
|---------------|-----------------------------|
| Language      | Go (1.22+)                  |
| Web Framework | [Fiber v2](https://github.com/gofiber/fiber) |
| Database      | PostgreSQL 16               |
| Cache         | Redis 7                     |
| ORM           | GORM                        |
| Container     | Docker & Docker Compose     |

---

## Project Structure

```
revise-redis/
├── .github/
│   └── copilot-instructions.md
├── cmd/
│   └── main.go                 # Entry point
├── config/
│   └── config.go               # Env/config loader
├── db/
│   └── postgres.go             # PostgreSQL connection (GORM)
├── cache/
│   └── redis.go                # Redis client setup
├── models/
│   └── news.go                 # News GORM model
├── repository/
│   └── news_repository.go      # DB queries for news
├── service/
│   └── news_service.go         # Business logic + cache layer
├── handler/
│   └── news_handler.go         # Fiber route handlers
├── router/
│   └── router.go               # Route registration
├── docker-compose.yml          # Postgres + Redis + App services
├── Dockerfile                  # Multi-stage Go build
├── .env                        # Environment variables
└── go.mod / go.sum
```

---

## Data Model — `News`

```go
type News struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Author    string    `json:"author"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

---

## API Endpoints

| Method | Path            | Description              | Cached |
|--------|-----------------|--------------------------|--------|
| GET    | `/api/news`     | List all news            | Yes    |
| GET    | `/api/news/:id` | Get single news by ID    | Yes    |
| POST   | `/api/news`     | Create a news article    | No     |
| PUT    | `/api/news/:id` | Update a news article    | No     |
| DELETE | `/api/news/:id` | Delete a news article    | No     |

---

## Caching Strategy

- Cache key format:
  - Single item: `news:<id>`
  - All items: `news:all`
- TTL: **5 minutes** by default
- Cache is **invalidated** on `POST`, `PUT`, and `DELETE` operations
- On cache miss, data is fetched from PostgreSQL and written back to Redis
- Redis values are stored as **JSON strings**

---

## Environment Variables (`.env`)

```env
APP_PORT=3000

# PostgreSQL
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=newsdb
DB_SSLMODE=disable

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_TTL=300
```

---

## Docker Compose Services

```yaml
services:
  app:       # Go Fiber application (port 3000)
  postgres:  # PostgreSQL 16 (port 5432)
  redis:     # Redis 7 (port 6379)
```

- The `app` service depends on both `postgres` and `redis`
- PostgreSQL data is persisted via a named volume
- Redis runs without persistence (cache-only mode)

---

## Coding Conventions

- Use **GORM** for all database interactions; avoid raw SQL unless necessary
- Use **`go-redis/v9`** (`github.com/redis/go-redis/v9`) as the Redis client
- All handlers must return consistent JSON responses:
  ```json
  { "success": true, "data": ..., "message": "..." }
  ```
- Error responses:
  ```json
  { "success": false, "message": "error description" }
  ```
- Use **service layer** for business logic; handlers should only parse input and call services
- Cache read/write/invalidation must live in the **service layer**, not handlers or repository
- Use `context.Background()` or pass Fiber's context-derived context to Redis calls
- Auto-migrate the `News` model on startup with `db.AutoMigrate(&models.News{})`

---

## Key Dependencies (`go.mod`)

```
github.com/gofiber/fiber/v2
github.com/redis/go-redis/v9
gorm.io/gorm
gorm.io/driver/postgres
github.com/joho/godotenv
```

---

## Running the Project

```bash
# Start all services
docker compose up --build

# Run locally (requires external Postgres & Redis)
go run cmd/main.go
```
