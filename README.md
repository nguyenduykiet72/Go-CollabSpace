# GoCollab Backend

Backend service for **GoCollab** — a real-time collaborative workspace platform built with Go and Gin framework. The system supports real-time document editing through CRDT (Yjs) over WebSocket, role-based workspace access control, JWT authentication with token rotation, and file storage via AWS S3 presigned URLs.

## Overview

GoCollab is a real-time collaboration platform designed for teams, featuring:
- **Real-time collaborative editing** — CRDT-based synchronization with Yjs over WebSocket, with server-side persistence to PostgreSQL
- **Workspace RBAC** — hierarchical role model (Owner > Admin > Editor > Viewer) enforced per-route via middleware
- **Nested document tree** — recursive CTE-powered tree with cycle detection on move operations
- **Hybrid search** — Reciprocal Rank Fusion combining pgvector semantic search + tsvector full-text search
- **JWT authentication** — access/refresh token rotation with SHA-256 hashed refresh tokens, Google OAuth support, password reset via async email
- **Background workers** — Asynq-based async job queue (Redis) for email delivery and search indexing
- **Rate limiting** — per-IP Redis-backed GCRA rate limiter on all public auth endpoints
- **Observability** — Prometheus metrics (WebSocket connections, processed messages), Zap structured logging, Grafana + Loki stack
- **Graceful shutdown** — orderly drain of HTTP server, WebSocket hub, Asynq workers, Redis, and database connections on SIGTERM

### Planned (AI Features)

- **AI-assisted workflows** — bounded agents for summarization, action item suggestions, pre-fill task descriptions with human-in-the-loop confirmation
- **Semantic search** — RAG pipeline: chunk → embed → pgvector → LLM answer generation
- **Event streaming** — Kafka event bus for audit logs, decoupling AI task queue, analytics pipeline

## Tech Stack

| Category | Technology |
|---|---|
| **Language** | Go 1.25 |
| **HTTP Framework** | Gin |
| **ORM** | GORM |
| **Database** | PostgreSQL 16 + pgvector + tsvector |
| **Cache / Pub-Sub** | Redis 8 (go-redis/v9) |
| **Real-time** | gorilla/websocket, Yjs (CRDT) |
| **Auth** | golang-jwt/v5 (HS256), bcrypt, Google OAuth2 |
| **Background Jobs** | Asynq (Redis-backed) |
| **File Storage** | AWS S3 (presigned URLs via SDK v2) |
| **Email** | Resend API |
| **Rate Limiting** | go-redis/redis_rate (GCRA) |
| **Metrics** | Prometheus client_golang |
| **Logging** | Zap (JSON structured) |
| **Migration** | Goose |
| **Monitoring** | Prometheus + Grafana + Loki + Promtail |
| **Containerization** | Docker multi-stage (distroless) |
| **Config** | cleanenv + godotenv |
| **Linting** | golangci-lint, lefthook (pre-commit) |

## Architecture

```
┌────────────────────────────────────────────────────────────────┐
│                     Client (Next.js)                           │
│               Yjs Client + WebSocket Provider                  │
└────────────────────┬───────────────────────────────────────────┘
                     │
                     │ HTTP REST + WebSocket (binary)
                     │
┌────────────────────▼───────────────────────────────────────────┐
│                  Go Backend (Monolith)                          │
│                                                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │ REST API     │  │ WebSocket    │  │ Asynq        │         │
│  │ (Gin)        │  │ Hub          │  │ Workers      │         │
│  │              │  │ (Yjs Sync)   │  │ (Email,      │         │
│  │ Auth, RBAC,  │  │ Redis PubSub │  │  Indexing)   │         │
│  │ Docs, Users  │  │ Persistence  │  │              │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
│                                                                │
│  Middleware: Auth → RBAC → RateLimit → ErrorHandler → Recovery │
│  Metrics:    GET /metrics (Prometheus)                          │
└──────────┬──────────────┬──────────────┬───────────────────────┘
           │              │              │
    ┌──────▼──────┐ ┌─────▼─────┐ ┌─────▼─────┐
    │ PostgreSQL  │ │  Redis    │ │  AWS S3   │
    │ + pgvector  │ │ (Cache,   │ │ (Files)   │
    │ + tsvector  │ │  PubSub,  │ │           │
    │             │ │  Asynq,   │ │           │
    │             │ │  RateLimit)│ │           │
    └─────────────┘ └───────────┘ └───────────┘
```

### Processing Flows

**Real-time Document Editing:**
```
User Edit → Yjs Client → WebSocket → Hub → Redis Pub/Sub → Broadcast to Peers
                                         └→ SaveQueue → AppendYjsUpdate (PostgreSQL)
```

**Hybrid Search (pgvector + tsvector):**
```
Query → [text_search CTE (tsvector rank)]  ─┐
      → [semantic_search CTE (pgvector)]    ─┤→ RRF Score → Top K Results
                                              │
```

**Background Job Flow:**
```
Service → Asynq Distributor → Redis Queue → Processor → (Send Email / Index Document)
```

## Project Structure

```
Go-CollabSpace/
├── cmd/server/
│   └── main.go                    # Entry point, signal handling, graceful shutdown
├── config/
│   ├── config.go                  # Config loader (CONFIG_PATH → MODE → env fallback)
│   ├── config.development.yml     # Dev defaults (non-secret values only)
│   ├── config.stage.yml           # Stage defaults (non-secret values only)
│   └── config.production.yml      # Prod defaults (non-secret values only)
├── internal/
│   ├── common/
│   │   ├── apperror/              # Typed AppError → HTTP status mapping
│   │   ├── infrastructure/        # Redis, S3, OAuth providers, email senders
│   │   └── token/                 # JWT provider (HS256 access tokens, random refresh)
│   ├── constant/                  # Role constants (Owner=1, Admin=2, Editor=3, Viewer=4)
│   ├── controller/                # HTTP handlers (auth, user, workspace, document, storage, ws)
│   ├── dto/                       # Request/response DTOs with binding tags
│   ├── initialize/                # Database initialization
│   ├── middleware/                 # Auth, RBAC, error handler, rate limiter, auth context
│   ├── model/                     # GORM models (User, Session, Workspace, Document, etc.)
│   ├── realtime/                  # WebSocket hub, client, Yjs protocol parser
│   ├── repository/                # Data access (GORM), transactor pattern
│   ├── router/                    # Route registration, handler grouping
│   ├── server/                    # Server wiring (DI), graceful shutdown
│   ├── service/                   # Business logic (auth, workspace, document, storage)
│   ├── telemetry/                 # Prometheus gauge/counter definitions
│   └── worker/                    # Asynq distributor + processor (email, search index)
├── migrations/                    # Goose SQL migrations (8 files)
├── pkg/
│   ├── hash/                      # SHA-256 token hashing
│   ├── httpx/                     # JSON response helpers
│   └── logger/                    # Zap logger initialization
├── docker/
│   ├── docker-compose.monitoring.yml   # Prometheus + Grafana + Loki + Promtail
│   └── config/                    # Prometheus, Loki, Promtail config files
├── Dockerfile                     # Multi-stage build (Go 1.25 → distroless, ~15MB)
├── .dockerignore
├── docker-compose.yml             # PostgreSQL (pgvector) + Redis for local dev
├── Makefile                       # Dev commands (run, build, lint, test, migrations)
├── lefthook.yml                   # Git hooks (pre-commit: tidy, fmt, lint)
├── .golangci.yml                  # Linter configuration
└── .env.sample                    # Environment variable template
```

### Architecture Layers

| Layer | Responsibility |
|---|---|
| **Controller** | HTTP request parsing, input validation (`binding` tags), error delegation via `ctx.Error()` |
| **Service** | Business logic, authorization checks, transaction orchestration |
| **Repository** | GORM queries, context-scoped transactions via `transactor.getDB(ctx)` |
| **Model** | GORM entities (UUID PKs, soft delete, table name overrides) |
| **Middleware** | Auth extraction, RBAC enforcement, rate limiting, centralized error→JSON mapping |

## Prerequisites

- **Go 1.25+** — [install](https://go.dev/doc/install)
- **PostgreSQL 16+** with `pgvector` and `uuid-ossp` extensions — [pgvector](https://github.com/pgvector/pgvector)
- **Redis 7+** — [install](https://redis.io/docs/getting-started/)
- **Docker & Docker Compose** (recommended for local infra)

## Quick Start

```bash
# 1. Clone and install deps
git clone <repository-url>
cd Go-CollabSpace
go mod download

# 2. Install dev tools
make install-tools

# 3. Start infrastructure
make docker-up          # PostgreSQL (pgvector/pg16) + Redis

# 4. Configure
cp .env.sample .env     # fill in your secrets

# 5. Run migrations
make db-up

# 6. Start server
make run                # http://localhost:8080
```

## Configuration

Configuration is loaded in the following precedence (highest wins):

1. **Environment variables** — always override file values
2. **`$CONFIG_PATH`** — explicit file path (K8s configmap, CI override)
3. **`./config/config.<MODE>.yml`** — derived from `$MODE` (default `development`)
4. **Pure env-only** — when no config file found (typical for production containers)

`.env` is loaded first (best-effort via godotenv) so local devs can keep secrets out of shell history.

### Environment Variables

| Variable | Description | Default | Required |
|---|---|---|---|
| `PORT` | Server port | `8080` | No |
| `MODE` | App mode (`development` / `stage` / `production`) | `development` | No |
| `CONFIG_PATH` | Explicit config file path | — | No |
| `ALLOWED_ORIGINS` | Comma-separated CORS/WS origin whitelist | `http://localhost:3000,http://localhost:3001` | No |
| `FE_RETURN_URL_WHITELIST` | Comma-separated frontend return URL whitelist for backend-generated links, e.g. password reset | `http://localhost:3000/reset-password,http://localhost:3001/reset-password` | No |
| `DB_HOST` | PostgreSQL host | — | Yes |
| `DB_PORT` | PostgreSQL port | — | Yes |
| `DB_USER` | Database user | — | Yes |
| `DB_PASSWORD` | Database password | — | Yes |
| `DB_NAME` | Database name | — | Yes |
| `DB_TYPE` | Database type | `postgres` | No |
| `JWT_ACCESS_SECRET` | JWT access token signing key | — | Yes |
| `JWT_REFRESH_SECRET` | JWT refresh token signing key | — | Yes |
| `JWT_ACCESS_DURATION` | Access token TTL | `1h` | No |
| `JWT_REFRESH_DURATION` | Refresh token TTL | `168h` | No |
| `REDIS_HOST` | Redis host | — | Yes |
| `REDIS_PORT` | Redis port | — | Yes |
| `REDIS_PASSWORD` | Redis password | — | Yes |
| `AWS_REGION` | S3 region | — | Yes |
| `AWS_ACCESS_KEY_ID` | S3 access key | — | Yes |
| `AWS_SECRET_ACCESS_KEY` | S3 secret key | — | Yes |
| `AWS_BUCKET_NAME` | S3 bucket name | — | Yes |
| `RESEND_API_KEY` | Resend email API key | — | Yes |
| `RESEND_FROM_EMAIL` | Sender email address | — | Yes |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID | — | Yes |
| `GOOGLE_CLIENT_SECRET` | Google OAuth client secret | — | Yes |
| `GOOGLE_REDIRECT_URL` | Google OAuth redirect URL | — | Yes |

## API

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication

JWT-based. After login, include token in header:

```
Authorization: Bearer <access_token>
```

Refresh tokens are rotated: each `/auth/refresh` call revokes the old session and issues a new pair.

### Endpoints

#### Auth (Public, rate-limited)

| Method | Path | Description | Rate Limit |
|---|---|---|---|
| `POST` | `/auth/register` | Register new user | 5/min |
| `POST` | `/auth/login` | Login (email + password) | 10/min |
| `POST` | `/auth/refresh` | Rotate access/refresh tokens | 30/min |
| `POST` | `/auth/forgot-password` | Request password reset email | 5/hour |
| `POST` | `/auth/reset-password` | Reset password with token | 10/hour |
| `POST` | `/auth/oauth/:provider` | Social login (Google) | 10/min |

`POST /auth/forgot-password` requires the frontend reset page URL in the request body:

```json
{
  "email": "user@example.com",
  "returnUrl": "https://app.example.com/reset-password"
}
```

The backend only accepts `returnUrl` values present in `FE_RETURN_URL_WHITELIST`, then appends the reset token server-side.

#### Users (Protected)

| Method | Path | Description | Min Role |
|---|---|---|---|
| `GET` | `/user/all` | List all users (paginated) | Authenticated |

#### Workspace (Protected)

| Method | Path | Description | Min Role |
|---|---|---|---|
| `POST` | `/workspace` | Create workspace | Authenticated |
| `GET` | `/workspace/:workspaceId` | Get workspace details (with members) | Viewer |
| `POST` | `/workspace/:workspaceId/members` | Add members to workspace | Admin |

#### Documents (Protected)

| Method | Path | Description | Min Role |
|---|---|---|---|
| `GET` | `/workspace/:workspaceId/document` | List workspace documents | Viewer |
| `GET` | `/workspace/:workspaceId/document/:docId` | Get document detail | Viewer |
| `GET` | `/workspace/:workspaceId/document/tree` | Get nested document tree | Viewer |
| `POST` | `/workspace/:workspaceId/document` | Create document | Editor |
| `PUT` | `/workspace/:workspaceId/document/:docId/move` | Move document (cycle-safe) | Editor |
| `PUT` | `/workspace/:workspaceId/document/:docId/snapshot` | Save document snapshot | Editor |

#### WebSocket

| Path | Description |
|---|---|
| `GET` `/ws/?token=<jwt>&doc_id=<uuid>` | Yjs binary sync (origin-checked, role-enforced) |

#### Monitoring

| Path | Description |
|---|---|
| `GET` `/metrics` | Prometheus metrics |

### Response Format

**Success:**
```json
{
  "statusCode": 200,
  "message": "Documents retrieved successfully",
  "data": { ... }
}
```

**Error:**
```json
{
  "statusCode": 401,
  "message": "Unauthorized",
  "errorKey": "UNAUTHORIZED"
}
```

## Database Migrations

The project uses **Goose** for SQL migrations.

```bash
make db-create name=add_users_table   # Create new migration
make db-up                            # Apply all pending
make db-down                          # Rollback last
make db-reset                         # Down to 0, then up
make db-status                        # Show status
```

### Current Schema

| Table | Purpose |
|---|---|
| `tbl_users` | Users (UUID PK, email unique, bcrypt password, soft delete) |
| `tbl_sessions` | Refresh token sessions (hashed token, blocked flag, expiry) |
| `tbl_workspaces` | Workspaces (slug unique, soft delete) |
| `tbl_roles` | Role definitions (Owner, Admin, Editor, Viewer) |
| `tbl_workspace_members` | User-workspace-role mapping (unique composite) |
| `tbl_documents` | Documents (nested via `doc_parent_id`, soft delete) |
| `tbl_document_states` | Yjs binary state (bytea) + plaintext |
| `tbl_document_chunks` | Chunks with vector embedding (pgvector HNSW) + tsvector (GIN) |
| `tbl_doc_embeddings` | Document-level vector embeddings (HNSW index) |
| `tbl_files` | File metadata (S3 key, status enum, expiry) |
| `tbl_password_resets` | Password reset tokens (hashed, expiry, used flag) |
| `tbl_audit_logs` | Audit trail (workspace, actor, entity, action, JSONB payload) |

## Development

```bash
make fmt              # Format with goimports
make fmt-check        # Check formatting
make lint             # Run golangci-lint
make test             # Run tests with race detection
make test-coverage    # Tests + HTML coverage report
make pre-commit       # fmt + lint + test
```

### Monitoring Stack (Optional)

```bash
cd docker
docker compose -f docker-compose.monitoring.yml up -d
```

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)
- **Loki**: http://localhost:3100

## Docker Deployment

### Build

```bash
# Development build
docker build -t go-collabspace:dev .

# Production build with metadata
docker build \
  --build-arg VERSION=$(git describe --tags --always) \
  --build-arg COMMIT=$(git rev-parse --short HEAD) \
  --build-arg BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -t go-collabspace:$(git rev-parse --short HEAD) .
```

### Run

```bash
# Using env vars directly
docker run --rm -p 8080:8080 \
  -e MODE=production \
  -e DB_HOST=postgres -e DB_PORT=5432 \
  -e DB_USER=postgres -e DB_PASSWORD=secret \
  -e DB_NAME=CollabSpace -e DB_TYPE=postgres \
  -e REDIS_HOST=redis -e REDIS_PORT=6379 -e REDIS_PASSWORD=secret \
  -e JWT_ACCESS_SECRET=... -e JWT_REFRESH_SECRET=... \
  -e ALLOWED_ORIGINS=https://app.example.com \
  go-collabspace:latest

# Using .env file
docker run --rm -p 8080:8080 --env-file .env go-collabspace:latest
```

### Image Details

| Property | Value |
|---|---|
| Base image | `gcr.io/distroless/static-debian12:nonroot` |
| Binary | Static (`CGO_ENABLED=0`), stripped (`-ldflags "-s -w"`) |
| User | `nonroot` (uid 65532) |
| Approx. size | ~15 MB |
| Build cache | Docker mount cache for Go build + module cache |

## Roadmap

### Implemented
- [x] JWT auth with refresh token rotation (SHA-256 hashed)
- [x] Google OAuth social login
- [x] Password reset flow (async email via Asynq + Resend)
- [x] Workspace CRUD with RBAC (Owner/Admin/Editor/Viewer)
- [x] Nested document tree (recursive CTE) with cycle-safe move
- [x] WebSocket hub with Yjs binary relay + Redis pub/sub fan-out
- [x] Yjs state persistence to PostgreSQL (bytea append)
- [x] Hybrid search SQL (pgvector + tsvector + RRF)
- [x] S3 presigned URL file uploads
- [x] Per-IP Redis rate limiting on auth endpoints
- [x] Prometheus metrics (WS connections, messages)
- [x] Grafana + Loki + Promtail monitoring stack
- [x] Graceful shutdown (HTTP + WS Hub + Asynq + Redis + DB)
- [x] Multi-stage Dockerfile (distroless, ~15MB)
- [x] Environment-aware config loading (CONFIG_PATH / MODE / env-only)

### Planned
- [ ] AI agent workflows with LLM integration (summarize, action items, human-in-the-loop)
- [ ] RAG search pipeline (chunk → embed → pgvector → LLM answer)
- [ ] OpenTelemetry distributed tracing
- [ ] Health endpoints (`/healthz`, `/readyz`)
- [ ] OpenAPI/Swagger documentation
- [ ] GitHub Actions CI/CD pipeline
- [ ] Yjs server-side snapshot compaction
- [ ] WebSocket presence + cursor awareness
- [ ] Audit log decorator on service mutations
- [ ] Load testing report (k6 / vegeta + pprof)

## Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Run `make pre-commit` before committing
4. Push and create Pull Request

Commit messages must follow the pattern `GCS-<number>: description` (enforced by lefthook).

## Contact

Email: elliotnguyen909@gmail.com
