# test-project

A Go project built with [Soliton-Go](https://github.com/soliton-go/framework) framework.

## Quick Start

```bash
# Install dependencies
go mod tidy

# Generate domain modules
soliton-gen domain User --fields "username,email,status:enum(active|inactive)"

# Run the server
go run ./cmd/main.go
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
| DELETE | /api/users/:id | Delete user |
