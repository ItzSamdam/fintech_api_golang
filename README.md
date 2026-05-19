# Fintech API (Golang)

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
