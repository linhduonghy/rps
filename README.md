# Resource Policy Service

A centralized resource-based authorization service (inspired by AWS Resource Policy)

The service provides CRUD APIs for managing resource policies and integrate with Open Policy Agent (OPA) for policy evaluation.

---

## Features

- Resource Policy CRUD
- OPA Integration
- PostgreSQL with sqlc
- Logging (slog)
- Prometheus Metrics
- OpenTelemetry Tracing
- Configuration via YAML & Environment Variables

---

## Tech Stack

- Go
- Chi Router
- PostgreSQL
- pgx
- sqlc
- OpenTelemetry
- Prometheus
- OPA (planned)

---

## Project Structure

```
cmd/
config/
internal/
migrations/
sql/
```

---

## Prerequisites

- Go 1.25+
- PostgreSQL 16+
- sqlc
- golang-migrate

---

## Run Locally

### 1. Generate sqlc code

```bash
sqlc generate
```

### 2. Run database migrations

```bash
go run ./cmd/migrate
```

### 3. Start the server

```bash
go run ./cmd/server
```

The server starts on:

```
http://localhost:8088
```

---

## Observability

- Structured logging with slog
- Prometheus metrics
- OpenTelemetry tracing

---