# test-project

A Go project built with [Soliton-Go](https://github.com/soliton-go/framework) framework.

## Quick Start

```bash
# Install dependencies
GOWORK=off go mod tidy

# Generate domain modules (--wire auto-injects into main.go)
soliton-gen domain User --fields "username,email,status:enum(active|inactive)" --wire

# Run the server
GOWORK=off go run ./cmd/main.go
```

## Project Structure

```
test-project/
├── cmd/main.go              # Entry point
├── configs/                 # Configuration
├── internal/
│   ├── domain/              # Domain layer (entities, repos, events)
│   ├── application/         # Application layer (commands, queries)
│   ├── infrastructure/      # Infrastructure layer (repo implementations)
│   └── interfaces/          # Interface layer (HTTP handlers)
└── go.mod
```

## API Endpoints

After generating domains, the following endpoints are available:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /health | Health check |
| POST | /api/users | Create user |
| GET | /api/users | List users |
| GET | /api/users/:id | Get user |
| PUT | /api/users/:id | Update user |
| PATCH | /api/users/:id | Partial update user |
| DELETE | /api/users/:id | Delete user |

> **Note**: If running in a monorepo with go.work, use `GOWORK=off` prefix for go commands.
