# SaaS Ecommerce — Backend API (LAKU)

Backend Go untuk platform ecommerce multi-tenant SaaS. Setiap "tenant" adalah
sebuah toko dengan data terisolasi (semua tabel domain membawa `tenant_id`).

## Tech Stack

- **Framework:** Gin
- **ORM:** GORM (PostgreSQL)
- **Cache:** Redis
- **Storage:** MinIO (S3-compatible)
- **Migration:** Goose (SQL-based)
- **DI:** google/wire
- **Auth:** JWT access token + refresh token rotasi (bcrypt)
- **Payment:** Mock provider (swappable ke Midtrans/Xendit)
- **Linter:** golangci-lint
- **Logger:** Logrus (JSON di production, colorized di development)

## Struktur Proyek

```
.
├── cmd/api/                  # Entry point + wire (DI)
├── constants/                # Shared constants (messages, formats)
├── docs/                     # erd.md (ERD) & api.md (API reference)
├── internal/
│   ├── entities/             # Shared DTOs & tenant context
│   ├── auth/                 # Login, refresh, logout, me
│   ├── tenant/               # Platform admin: tenants & plans
│   ├── user/                 # Users & roles (RBAC)
│   ├── category/             # Kategori produk (tree)
│   ├── product/              # Produk + images + stok (reserved/incoming)
│   ├── customer/             # Customer + addresses + activities
│   ├── order/                # Order + items + timeline + status flow
│   ├── payment/              # Payment methods + charge + webhook
│   ├── inventory/            # Warehouse, stock movements, shipments
│   ├── setting/              # Store/company/tax/shipping/notifications
│   ├── analytics/            # Aggregates: summary, revenue, sales, top
│   └── review/               # Review + moderasi (rating produk)
├── middleware/               # Auth, tenant scope, role guard
├── migration/                # SQL migration files (goose)
├── pkg/                      # Reusable packages
├── third_party/payment/      # Payment provider interface + mock client
├── scripts/                  # Git hooks & module scaffolding
├── docker-compose.yml        # postgres, redis, minio
├── Makefile
└── .env.example              # Salin ke .env
```

## Quick Start

```bash
# 1. Copy env & start infra
cp .env.example .env
docker compose up -d

# 2. Install tools
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 3. Run migrations & start server
make migrate-up
make run
```

Server berjalan di `:8000`. Seed user (password `admin123`):

| Email | Role |
|---|---|
| `platform@laku.id` | Platform admin |
| `admin@laku.id` | Store admin (tenant demo `laku-demo`) |
| `staff@laku.id` | Staff |

## Module Pattern

Setiap module: `model.go` (struct + DTOs) → `repository.go` (interface + impl,
semua query menerima `tx *gorm.DB` + `tenantID`) → `service.go` (business logic
+ interface untuk cross-service) → `handler.go` → `routes.go` (REST).

- Semua data difilter `WHERE tenant_id = ?` dari klaim JWT
- Cross-module via interface (contoh: order → `product.ServiceInterface`)
- Response envelope: `{ "code": 200, "data": {...}, "msg": "..." }`
- Error via `apperror` → HTTP status mapping
- Transaksi ditulis memakai `s.db.WithContext(ctx).Transaction(...)`

## API Endpoints

Dokumentasi lengkap: [`docs/api.md`](docs/api.md). Ringkasan:

- `POST /api/v2/auth/login|refresh|logout`, `GET /auth/me`
- `/api/v2/categories`, `/products`, `/customers`, `/orders` (CRUD REST)
- `/api/v2/payments/methods`, `/payments/charge`, `/payments/webhook`
- `/api/v2/inventory/warehouses|stock-movements|shipments`
- `/api/v2/settings/*` (store, company, tax-rates, shipping-zones, notifications)
- `/api/v2/analytics/summary|revenue|sales|top-products|top-categories`
- `/api/v2/reviews`, `/api/v2/users`, `/api/v2/roles`
- `/api/v2/platform/tenants|plans` (khusus platform admin)

## ERD

Database design (24 tabel): [`docs/erd.md`](docs/erd.md)

## Commands

```bash
make run                        # Run server
make dev                        # Run server (air hot reload)
make build                      # Build binary
make test                       # Run tests (gotestsum)
make lint                       # Run linter
make wire                       # Regenerate wire_gen.go

make migrate-up / down / status / create name=xxx / redo / reset

make new-module name=xxx [tests=true]   # Scaffold module
make docker-up / down / build           # Docker compose / image
```
