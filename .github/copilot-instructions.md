# Copilot Instructions

## Project Overview

This repository is a **CRUD News API** built with **Go Fiber**, backed by **PostgreSQL** for persistent storage and **Redis** for caching, orchestrated with **Docker Compose**.

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
