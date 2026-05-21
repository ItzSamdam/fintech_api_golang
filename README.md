# Fintech API (Golang)
# PS: This is only a challenge building a fintech api coming from a nodejs/laravel background. pron to errors. so please check well before using as a template for production apps
A modular, hexagonal-architecture based fintech API built with Go.  
This project provides wallet services, transfers, bill payments, savings goals, compliance checks, and admin features.

---

## 📂 Project Structure

```text
fintech_api_golang/
├── cmd/
│   └── api/                # Application entry point
│       └── main.go
│
├── internal/
│   ├── core/               # Domain models & business logic
│   │   ├── entities/       # Core entities
│   │   │   ├── user.go
│   │   │   ├── wallet.go
│   │   │   ├── transaction.go
│   │   │   ├── bill_payment.go
│   │   │   ├── savings_goal.go
│   │   │   └── audit_log.go
│   │   ├── interfaces/     # Repository interfaces (hexagonal ports)
│   │   ├── services/       # Business logic services
│   │
│   ├── handlers/           # HTTP handlers (controllers)
│   │   ├── auth_handler.go
│   │   ├── wallet_handler.go
│   │   ├── transfer_handler.go
│   │   ├── airtime_handler.go
│   │   ├── data_handler.go
│   │   ├── electricity_handler.go
│   │   ├── betting_handler.go
│   │   ├── savings_handler.go
│   │   ├── notification_handler.go
│   │   ├── compliance_handler.go
│   │   ├── support_handler.go
│   │   └── admin/          # Admin subpackage
│   │       ├── user_admin.go
│   │       ├── transaction_admin.go
│   │       ├── wallet_admin.go
│   │       ├── kyc_admin.go
│   │       ├── provider_admin.go
│   │       ├── fee_admin.go
│   │       ├── savings_admin.go
│   │       ├── report_admin.go
│   │       └── role_admin.go
│   │
│   ├── repository/         # Data layer implementations
│   │   ├── postgres/       # Postgres repositories
│   │   ├── redis/          # Redis caches
│   │   └── provider/       # External provider clients
│   │
│   ├── middleware/         # HTTP middleware (auth, rate limit, logging, etc.)
│   ├── dto/                # Request/response DTOs
│   ├── pkg/                # Shared internal packages (db, logger, utils, queue, cache, webhook)
│   └── config/             # Configuration management
│
├── api/
│   ├── routes/             # Router setup
│   │   ├── routes.go
│   │   ├── v1/             # Versioned routes
│   │   └── middleware_routes.go
│   └── docs/               # API documentation (Swagger, Postman)
│
├── migrations/             # SQL migrations
├── scripts/                # Build/test/migration scripts
├── tests/                  # Unit, integration, contract, e2e tests
├── deployments/            # Docker, Kubernetes, Terraform configs
├── .env.example            # Example environment variables
├── go.mod                  # Go modules definition
├── go.sum
├── Makefile                # Build/test commands
└── README.md
```

## 🚀 Features

- Authentication & Authorization: JWT-based auth, role-based access.
- Wallet Management: Create, fund, transfer between wallets.
- Transfers: Bank transfers via NIP, name enquiry.
- Bill Payments: Airtime, data, electricity, betting.
- Savings Goals: Create and contribute to savings goals.
- Notifications: Push/email notifications.
- Compliance: KYC, tier limits, audit logging.
- Admin Panel: Manage users, wallets, transactions, providers, fees, roles.
- Middleware: Logging, recovery, CORS, request ID, rate limiting.

## 🛠️ Tech Stack
- Language: Go
- Database: PostgreSQL, Redis
- Queue: RabbitMQ
- Logging: Zap
- Architecture: Hexagonal (ports & adapters)
- Deployment: Docker, Kubernetes, Terraform
- Docs: Swagger, Postman

## 📦 Setup
Clone the repository:

```bash
git clone https://github.com/ItzSamdam/fintech_api_golang.git
cd fintech_api_golang
```
Install dependencies:
```bash
go mod tidy
```
Configure environment variables:
```bash
cp .env.example .env
```
Run migrations:
```bash
./scripts/migrate-up.sh
```
Install Air for live reloading:
```bash
go install github.com/air-verse/air@latest
```
Start the API with Air:
```bash
air
```
## 🧪 Testing
- Unit tests: go test ./tests/unit/...
- Integration tests: go test ./tests/integration/...
- Contract tests: go test ./tests/contract/...
- End-to-end tests: go test ./tests/e2e/...

## 📖 Documentation
- Swagger: http://localhost:[port]/swagger
- Postman collection: api/docs/postman_collection.json (not updated)

## 📜 License
This project is licensed under the MIT License.
