# GoCollab Backend

Backend service for **GoCollab** - an AI-augmented real-time collaborative workspace platform, built with Go and Gin framework. The system supports real-time collaboration through CRDT (Yjs), semantic search with vector embeddings, and agent-assisted workflows with human-in-the-loop confirmation.

## Overview

GoCollab is a real-time collaboration platform designed for teams and SMBs, integrating features:
- **Real-time collaborative editing**: CRDT-based synchronization with Yjs, supporting offline-first and conflict-free merging
- **AI-assisted workflows**: Bounded agents supporting summarization, action item suggestions, and pre-fill task descriptions with human confirmation
- **Semantic search**: Vector search across internal knowledge base (docs, tickets, transcripts) for quick question answering
- **Event-driven architecture**: Audit/event stream with replay capability, easy to scale into microservices
- **Privacy & governance**: Tenant isolation, AI governance hooks, rate limiting on LLM calls

## Tech Stack

### Core Framework & Language
- **Go 1.25+**: Primary programming language
- **Gin**: HTTP web framework
- **GORM**: ORM for database operations

### Database & Storage
- **PostgreSQL**: Primary relational database for persistent data (users, teams, documents metadata, transactions)
- **pgvector**: PostgreSQL extension for vector embeddings and semantic search
- **Redis**: Caching, presence tracking, ephemeral locks, rate limiting, pub/sub for low-latency notifications

### Real-time & Collaboration
- **Yjs**: CRDT library for conflict-free collaborative editing
- **WebSocket**: Real-time bidirectional communication for presence, awareness, and Yjs sync

### Migration & Tooling
- **Goose**: Database migration tool
- **golangci-lint**: Static code analysis
- **Zap**: Structured logging

### Infrastructure (Planned)
- **Kafka**: Event bus for audit logs, decoupling AI task queue, analytics pipeline
- **Elasticsearch**: Full-text + vector search for semantic queries
- **Prometheus + Grafana**: Metrics and monitoring
- **OpenTelemetry**: Distributed tracing

## Architecture

### High-level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Client (Next.js)                        │
│              Yjs Client + WebSocket Provider                  │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       │ HTTP/gRPC + WebSocket
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                    Go Backend (Monolith)                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ HTTP/gRPC    │  │ WebSocket    │  │ Background   │      │
│  │ Handlers     │  │ Server       │  │ Workers      │      │
│  │ (Auth, Docs) │  │ (Yjs Sync)   │  │ (AI, Index)  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└──────────┬──────────────┬──────────────┬────────────────────┘
           │              │              │
    ┌──────▼──────┐ ┌─────▼─────┐ ┌─────▼─────┐
    │ PostgreSQL  │ │  Redis    │ │  Kafka    │
    │ + pgvector  │ │ (Cache,   │ │ (Events)  │
    │             │ │  Pub/Sub) │ │           │
    └─────────────┘ └───────────┘ └───────────┘
```

### Main Processing Flows

**1. Real-time Document Editing:**
```
User Edit → Yjs Client → WebSocket → Go Yjs Server → CRDT Update → Broadcast to Peers
```

**2. AI-assisted Workflow:**
```
Document Save → Event (Kafka) → Worker → Index to Elasticsearch → 
AI Job Queue → LLM Call (with context) → Save Draft → Notification → 
Human Approval → Publish Summary
```

**3. Semantic Search:**
```
Query → Vector Embedding → pgvector/Elasticsearch → 
Context Retrieval → LLM Answer Generation → Response
```

## Project Structure

```
Go-CollabSpace/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── config/
│   ├── config.go                # Configuration management
│   └── config.dev.yml           # Development config template
├── internal/
│   ├── common/
│   │   ├── apperror/            # Custom error types
│   │   └── token/               # JWT token provider
│   ├── constant/                # Application constants
│   ├── controller/              # HTTP handlers (presentation layer)
│   ├── dto/                     # Data transfer objects
│   ├── initialize/              # Initialization logic (DB, etc.)
│   ├── middleware/              # HTTP middleware (auth, error handling)
│   ├── model/                   # Domain models (GORM entities)
│   ├── repository/              # Data access layer
│   ├── router/                  # Route definitions
│   ├── server/                  # Server setup and initialization
│   └── service/                 # Business logic layer
├── migrations/                  # Goose migration files
├── pkg/
│   ├── httpx/                   # HTTP utilities (response helpers)
│   └── logger/                  # Logging utilities
├── docker-compose.yml           # Local development infrastructure
├── Makefile                     # Development commands
├── go.mod                       # Go module dependencies
└── README.md                    # This file
```

### Architecture Pattern

The project follows **Clean Architecture** with the following layers:

- **Controller**: HTTP request/response handling, validation
- **Service**: Business logic, orchestration
- **Repository**: Data persistence abstraction
- **Model**: Domain entities and database schema

## Prerequisites

- **Go 1.25+**: [Installation guide](https://go.dev/doc/install)
- **PostgreSQL 16+** with `pgvector` extension: [pgvector installation](https://github.com/pgvector/pgvector)
- **Redis 7+**: [Redis installation](https://redis.io/docs/getting-started/)
- **Docker & Docker Compose** (optional, for local development)
- **Goose**: Migration tool (will be installed via Makefile)

## Installation & Setup

### 1. Clone Repository

```bash
git clone <repository-url>
cd Go-CollabSpace
```

### 2. Install Dependencies

```bash
go mod download
go mod verify
```

### 3. Install Development Tools

```bash
make install-tools
```

Or install manually:
```bash
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### 4. Setup Infrastructure with Docker Compose

```bash
# Start PostgreSQL (with pgvector) and Redis
make docker-up

# Check running containers
docker ps
```

### 5. Database Setup

```bash
# Apply migrations
make db-up

# Check migration status
make db-status
```

### 6. Configuration

Create `.env` file from template or copy `config/config.dev.yml` and adjust:

```bash
cp config/config.dev.yml config/config.local.yml
```

Or use environment variables:

```bash
# Server
export PORT=8080
export MODE=development

# Database
export DB_HOST=localhost
export DB_PORT=5445
export DB_USER=postgres
export DB_PASSWORD=123456789
export DB_NAME=CollabSpace
export DB_TYPE=postgres

# JWT
export JWT_ACCESS_SECRET=your-access-secret-key
export JWT_REFRESH_SECRET=your-refresh-secret-key
export JWT_ACCESS_DURATION=1h
export JWT_REFRESH_DURATION=168h
```

### 7. Run Application

```bash
# Development mode
make run

# Or directly
go run cmd/server/main.go
```

Server will run at `http://localhost:8080`

## Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | No |
| `MODE` | Gin mode (development/release/test) | `development` | No |
| `DB_HOST` | PostgreSQL host | - | Yes |
| `DB_PORT` | PostgreSQL port | - | Yes |
| `DB_USER` | Database user | - | Yes |
| `DB_PASSWORD` | Database password | - | Yes |
| `DB_NAME` | Database name | - | Yes |
| `DB_TYPE` | Database type (postgres/mysql) | `postgres` | No |
| `JWT_ACCESS_SECRET` | JWT access token secret | - | Yes |
| `JWT_REFRESH_SECRET` | JWT refresh token secret | - | Yes |
| `JWT_ACCESS_DURATION` | Access token duration | `1h` | No |
| `JWT_REFRESH_DURATION` | Refresh token duration | `168h` | No |

### Config File Structure

Configuration is loaded in the following priority order:
1. Environment variables (highest priority)
2. `config/config.dev.yml` (or specified file)
3. Default values

## Database Migrations

The project uses **Goose** for database migrations.

### Create New Migration

```bash
make db-create name=add_users_table
```

Migration file will be created at `migrations/YYYYMMDDHHMMSS_add_users_table.sql`

### Apply Migrations

```bash
# Apply all pending migrations
make db-up

# Rollback last migration
make db-down

# Reset database (rollback all then reapply)
make db-reset

# View migration status
make db-status
```

### Migration Best Practices

- Each migration must be **reversible** (has both `up` and `down`)
- Do not modify migrations that have been applied to production
- Test migrations on staging before applying to production
- Use transactions when possible

## Development

### Code Quality

```bash
# Format code
make fmt

# Check formatting
make fmt-check

# Run linter
make lint

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run all checks (pre-commit)
make pre-commit
```

### Project Structure Guidelines

- **Controllers**: Only handle HTTP request/response, validation, error handling
- **Services**: Contains business logic, not dependent on HTTP layer
- **Repositories**: Data access layer, abstract database operations
- **Models**: Domain entities, GORM models
- **DTOs**: Data transfer objects for API contracts

### Logging

The project uses **Zap** for structured logging:

```go
logger.Log.Info("User created", zap.String("userID", userID))
logger.Log.Error("Database error", zap.Error(err))
```

Log levels: `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`

## API Documentation

### Base URL

```
Development: http://localhost:8080/api/v1
```

### Authentication

API uses JWT-based authentication. After login, include token in header:

```
Authorization: Bearer <access_token>
```

### Endpoints

#### Authentication

**POST** `/api/v1/auth/register`
- Register new user
- Request body: `{ "email": "...", "password": "...", "name": "..." }`
- Response: `{ "accessToken": "...", "refreshToken": "..." }`

**POST** `/api/v1/auth/login`
- User login
- Request body: `{ "email": "...", "password": "..." }`
- Response: `{ "accessToken": "...", "refreshToken": "..." }`

#### Workspace (Protected)

**POST** `/api/v1/workspace`
- Create new workspace
- Headers: `Authorization: Bearer <token>`
- Request body: `{ "name": "...", "description": "..." }`

**GET** `/api/v1/workspace/:workspaceId`
- Get workspace information
- Headers: `Authorization: Bearer <token>`

#### Document (Protected)

**POST** `/api/v1/document`
- Create new document
- Headers: `Authorization: Bearer <token>`
- Request body: `{ "workspaceId": "...", "title": "...", "content": "..." }`

**GET** `/api/v1/document/:docId`
- Get document detail
- Headers: `Authorization: Bearer <token>`

**GET** `/api/v1/document/doc/:workspaceId`
- Get all documents in workspace
- Headers: `Authorization: Bearer <token>`

### Error Response Format

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {}
  }
}
```

## Testing

### Unit Tests

```bash
# Run all tests
make test

# Run tests with race detection
go test -race ./...

# Run tests for specific package
go test ./internal/service/...
```

### Integration Tests

Integration tests will test with test database. Setup test database before running:

```bash
# Create test database
createdb collabspace_test

# Run integration tests
go test -tags=integration ./...
```

### Test Coverage

```bash
make test-coverage
```

Coverage report will be generated at `coverage.html`

## Deployment

### Build Binary

```bash
make build
```

Binary will be created at `bin/Go-CollabSpace`

### Production Considerations

1. **Environment Variables**: Use secret management (AWS Secrets Manager, HashiCorp Vault, etc.)
2. **Database**: Setup connection pooling, read replicas for scaling
3. **Redis**: Configure persistence and high availability
4. **Monitoring**: Setup Prometheus metrics, Grafana dashboards
5. **Logging**: Centralized logging (ELK stack, Loki, etc.)
6. **Security**: 
   - Enable HTTPS/TLS
   - Rate limiting
   - CORS configuration
   - Input validation and sanitization

### Docker Deployment

```dockerfile
# Example Dockerfile (needs to be created)
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o bin/Go-CollabSpace cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/Go-CollabSpace .
CMD ["./Go-CollabSpace"]
```

## Roadmap & Future Enhancements

### Planned Features

- [ ] WebSocket server for Yjs synchronization
- [ ] Vector search integration with pgvector
- [ ] AI agent workflows with LLM integration
- [ ] Event streaming with Kafka
- [ ] Elasticsearch integration for full-text + vector search
- [ ] Background workers for async tasks
- [ ] Rate limiting middleware
- [ ] Tenant isolation and multi-tenancy support
- [ ] OpenTelemetry distributed tracing
- [ ] Prometheus metrics export

### Architecture Evolution

Currently a **monolithic architecture**, but designed to easily evolve into **microservices**:

- Event-driven design with Kafka
- Clear separation of concerns (Clean Architecture)
- Repository pattern for database abstraction
- Service layer can be extracted into independent services

## Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Create Pull Request

### Code Style

- Follow Go conventions and best practices
- Run `make pre-commit` before committing
- Write tests for new features
- Update documentation when needed

## Support & Contact

Email: elliotnguyen909@gmail.com


