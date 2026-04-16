# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this project is

**mkit** is a Go-based microservices toolkit — a collection of reusable packages for building HTTP (Gin) and gRPC services with built-in support for PostgreSQL, Redis, RabbitMQ, MinIO, OpenTelemetry, and cron jobs.

## Common commands

```bash
# Run examples
go run ./example/ginapp/cmd/main.go
go run ./example/grpcapp/cmd/main.go

# Start local infrastructure (PostgreSQL, Redis)
cd example && docker-compose up -d

# Proto code generation (requires buf CLI v1.57.0)
cd example/grpcapp && buf generate

# Database migrations (requires atlas CLI v1.1.0)
atlas migrate diff <migration_name> --env gorm
atlas migrate hash --dir "file://db/migrations"

# Dependency management
go mod tidy
```

There are no tests or linting configs in the codebase yet.

## Architecture

### Core packages (`pkg/`)

- **`server/`** — Central orchestrator. Boots Gin HTTP and/or gRPC servers simultaneously, handles graceful shutdown, and wires all dependencies together via functional options (`server.Postgres(db)`, `server.Logger(logger)`, etc.).
- **`config/`** — Viper-based YAML config. Schema defined in `pkg/config/app.go` covering Postgres, Redis, HTTP, gRPC, Tracing, Swagger, MinIO, and RabbitMQ.
- **`middleware/`** — Gin middleware (logger, recovery, CORS) and gRPC interceptors (proto validation, logging, panic recovery).
- **`error/serviceerror/`** — Custom error types that map to HTTP status codes (400/401/403/404/409/500). Middleware reads these and writes the response automatically.
- **`error/repoerror/`** — Repository-level error types.
- **`tracing/`** — OpenTelemetry setup: OTLP gRPC exporter for traces, metrics, and logs.
- **`pubsub/`** — RabbitMQ topic-based pub/sub. Two modes: `Broadcast` (one instance) and `BroadcastAll` (all instances). Uses proto serialization.
- **`cron/`** — Distributed cron via Redsync (Redis-backed) to prevent duplicate execution across instances.
- **`postgres/`**, **`cache/redis/`**, **`cache/redsync/`**, **`rabbitmq/`**, **`minio/`** — Thin wrappers with OpenTelemetry instrumentation.
- **`httpclient/`** — HTTP client wrapper with retry logic.

### Layered example apps (`example/`)

Each example app follows this structure:
```
cmd/main.go              → entry point
internal/server/serve.go → bootstrap: builds deps, calls server.NewServer(...)
internal/controller/     → HTTP/gRPC handlers
internal/service/        → business logic
internal/repository/     → GORM data access
internal/model/          → domain models
config/config.yaml       → environment config
```

### Key design decisions

- **Functional options for DI**: All dependencies are optional. Pass only what your service needs to `server.NewServer(...)`.
- **gRPC + JSON transcoding**: gRPC services can be called as REST via Vanguard (h2c). Configure in `pkg/server/grpc/server.go`.
- **Proto validation**: Uses `protovalidate` library. Validated automatically by the gRPC unary interceptor.
- **Code-first migrations**: Write GORM models → Atlas generates SQL migrations from them.
- **Bob query builder**: Used alongside GORM for complex queries. Config in `example/bob.gen.yaml`.

## Required external tools

| Tool | Version | Purpose |
|------|---------|---------|
| `buf` | v1.57.0 | Proto code generation |
| `protoc-gen-go` | v1.25.0 | Go protobuf stubs |
| `protoc-gen-go-grpc` | v1.5.1 | gRPC stubs |
| `atlas` | v1.1.0 | DB migrations from GORM models |

## Configuration

Config files are YAML, loaded by Viper. Environment-specific overrides: `config-dev.yaml`, `config-prod.yaml`. The full config schema is in `pkg/config/app.go`.

Local dev services run via `example/docker-compose.yml`:
- PostgreSQL 18.1 on port 5433
- Redis 7 on port 6379
