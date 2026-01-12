# application

A Go project built with [Soliton-Go](https://github.com/soliton-go/framework) framework.

## Quick Start

```bash
# Install dependencies
GOWORK=off go mod tidy

# Generate domain modules (--wire auto-injects into main.go)
soliton-gen domain User --fields "username,email,status:enum(active|inactive)" --wire

# Enable soft delete (optional)
soliton-gen domain User --fields "username,email" --soft-delete --wire

# Run the server
GOWORK=off go run ./cmd/main.go
```

## Project Structure

```
application/
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

Domain CRUD endpoints are generated per module:

| Module | Base Endpoint |
|--------|----------------|
| User | /api/users |
| Product | /api/products |
| Order | /api/orders |
| Inventory | /api/inventories |
| Payment | /api/payments |
| Shipping | /api/shippings |
| Promotion | /api/promotions |
| Review | /api/reviews |

Service action endpoints:

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/payments/authorize | Authorize payment |
| POST | /api/payments/:id/capture | Capture payment |
| POST | /api/payments/:id/refund | Refund payment |
| POST | /api/payments/:id/cancel | Cancel payment |
| POST | /api/inventories/:id/adjust | Adjust stock |
| POST | /api/inventories/:id/reserve | Reserve stock |
| POST | /api/inventories/:id/release | Release stock |
| POST | /api/inventories/:id/stock-in | Stock in |
| POST | /api/inventories/:id/stock-out | Stock out |
| POST | /api/shippings/shipments | Create shipment |
| POST | /api/shippings/:id/tracking | Update tracking |
| POST | /api/shippings/:id/deliver | Mark delivered |
| POST | /api/shippings/:id/cancel | Cancel shipment |
| POST | /api/promotions/validate | Validate promotion |
| POST | /api/promotions/apply | Apply promotion |
| POST | /api/promotions/revoke | Revoke promotion |
| POST | /api/reviews/submit | Submit review |
| POST | /api/reviews/:id/moderate | Moderate review |
| POST | /api/reviews/:id/reply | Reply review |

### Pagination

List endpoints support pagination:

```bash
curl "http://localhost:8080/api/users?page=1&page_size=20"
```

Response:
```json
{
  "items": [...],
  "total": 100,
  "page": 1,
  "page_size": 20,
  "total_pages": 5
}
```

> **Note**: If running in a monorepo with go.work, use `GOWORK=off` prefix for go commands.
