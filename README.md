# Sliceflow Backend (Go + Clean Arch + Postgres)

README Manufacturing Execution System (MES) for 3D printing farms (FDM/SLS) written in Go following Clean Architecture + GORM + PostgreSQL + Redis. It includes intelligent background workers and role-based security.

# Summary

This service centralizes production control. It manages the lifecycle of production orders, material inventory control, and machine monitoring. It exposes a secure REST API and features a "Brain" (Worker) that automatically prioritizes tasks based on deadlines.

> 📢 **Includes Swagger Documentation & Bruno Collection for API testing.**

# Typical Flow:

1. The operator or admin authenticates and obtains a JWT token.
2. The admin or user creates production orders.
3. The "Worker" (background process) monitors deadlines and upgrades priority (P1/P2) if the deadline is approaching.
4. The Dashboard endpoint shows real-time data (censoring prices for non-admins).

# Features

✅ **Smart Worker:** Auto-prioritization of orders (<24hs urgent) and detection of "Zombie" (stalled) orders.

✅ **Data Security:** Automatic censorship of financial fields (`Revenue`, `Price`) for standard operators.

✅ **Stock Management:** Material validation before prioritizing an order.

✅ **High Performance Caching:** Redis implementation for heavy read endpoints (Dashboard).

✅ **Swagger UI:** Full API documentation available at `/swagger/index.html`.

✅ **Bruno Support:** API Collection included for easy testing.

✅ **Docker Compose:** Ready-to-use environment with DB and Cache.

# Tech Stack

  Language: Go ≥ 1.21.

  Architecture: Hexagonal (Clean Architecture).

  GORM: GORM (PostgreSQL).

  Cache: Redis.

  Database: PostgreSQL.

  Containers: Docker / Docker Compose.

# Setup

  Requirements
    Go ≥ 1.21.
    Docker Desktop (with WSL2 enabled on Windows).
    Git.

# Environment Variables

Copy `.env.example` to `.env.dev` and complete:

PORT=8181.

JWT_SECRET=<your_secure_secret>.
TTL=1440 (Token life in minutes).

REDIS_ADDR=redis:6379.

DB (if using Docker Compose):
POSTGRES_HOST=db, POSTGRES_USER=postgres, POSTGRES_PASSWORD=postgres, POSTGRES_DB=sliceflow_db, POSTGRES_PORT=5432.

Local without Docker: POSTGRES_HOST=localhost - POSTGRES_PORT=5432.

# Quick Start (Docker)

1) Prepare the environment

  Copy the environment variables.
  Ensure ports 8080 and 5432 are free.

2) Start everything
  `docker-compose up --build -d`

  The API runs inside the container at :8080.

3) Useful commands
  `docker-compose down` # stop everything.
  `docker-compose logs -f api` # view logs (including Smart Worker logs).

# Run Local (without Docker)
  `go mod tidy`
  `go run cmd/api/main.go`
  http://localhost:8080

# API Documentation & Tools

### 📘 Swagger UI
Once the server is running, visit:
`http://localhost:8080/swagger/index.html`

### 🐶 Bruno Collection
You will find a folder named `bruno_collection` (or `docs/bruno`) in the root of this repository. Open **Bruno**, click "Open Collection", and select that folder to load all pre-configured requests.

# Auth Flow

Login with username -> JWT.

Send `Authorization: Bearer <token>` to protected endpoints.

Roles:
- `admin`: Sees everything (including revenue/prices).
- `user`: Sees orders and machines (prices hidden/zeroed).

## Main Endpoints

**Base URL:** `http://localhost:8080`
> Use `{{token}}` (Bearer) for protected endpoints.

### Login
- **URL:** `/auth/login`
- **Method:** `POST`
- **Body:**
```json
{
  "username": "admin",
  "password": "admin"
}
