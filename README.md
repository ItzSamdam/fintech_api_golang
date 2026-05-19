fintech_api_golang/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
│
├── internal/
│   ├── core/                       # Domain models & business logic
│   │   ├── entities/
│   │   │   ├── user.go
│   │   │   ├── wallet.go
│   │   │   ├── transaction.go
│   │   │   ├── bill_payment.go
│   │   │   ├── savings_goal.go
│   │   │   └── audit_log.go
│   │   │
│   │   ├── interfaces/             # Repository interfaces (hexagonal)
│   │   │   ├── user_repository.go
│   │   │   ├── wallet_repository.go
│   │   │   ├── transaction_repository.go
│   │   │   ├── savings_repository.go
│   │   │   └── audit_repository.go
│   │   │
│   │   └── services/               # Business logic layer
│   │       ├── auth_service.go
│   │       ├── wallet_service.go
│   │       ├── transfer_service.go
│   │       ├── airtime_service.go
│   │       ├── data_service.go
│   │       ├── electricity_service.go
│   │       ├── betting_service.go
│   │       ├── savings_service.go
│   │       ├── notification_service.go
│   │       ├── compliance_service.go
│   │       └── admin_service.go
│   │
│   ├── handlers/                   # HTTP handlers (controllers)
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
│   │   └── admin/                  # Admin subpackage
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
│   ├── repository/                 # Data layer implementations
│   │   ├── postgres/
│   │   │   ├── user_repo.go
│   │   │   ├── wallet_repo.go
│   │   │   ├── transaction_repo.go
│   │   │   ├── savings_repo.go
│   │   │   └── audit_repo.go
│   │   ├── redis/
│   │   │   ├── session_cache.go
│   │   │   └── rate_limiter.go
│   │   └── provider/               # External provider clients
│   │       ├── provider_interface.go
│   │       ├── airtime/
│   │       │   ├── mtn.go
│   │       │   ├── glo.go
│   │       │   └── airtel.go
│   │       ├── data/
│   │       │   ├── mtn_data.go
│   │       │   └── glo_data.go
│   │       ├── electricity/
│   │       │   ├── ikeja.go
│   │       │   ├── eko.go
│   │       │   └── abuja.go
│   │       ├── betting/
│   │       │   ├── bet9ja.go
│   │       │   ├── sportybet.go
│   │       │   └── onexbet.go
│   │       ├── bank/
│   │       │   ├── nip_client.go
│   │       │   └── name_enquiry.go
│   │       └── registry.go        # Provider registry with toggles
│   │
│   ├── middleware/                 # HTTP middleware
│   │   ├── auth.go                # JWT verification
│   │   ├── role.go                # Role-based access
│   │   ├── rate_limit.go
│   │   ├── logger.go
│   │   ├── recovery.go
│   │   ├── cors.go
│   │   ├── request_id.go
│   │   ├── tier_limiter.go        # Check tier limits
│   │   └── audit.go               # Log admin actions
│   │
│   ├── dto/                        # Data Transfer Objects
│   │   ├── request/
│   │   │   ├── auth_request.go
│   │   │   ├── transfer_request.go
│   │   │   ├── bill_request.go
│   │   │   ├── savings_request.go
│   │   │   └── admin_request.go
│   │   └── response/
│   │       ├── auth_response.go
│   │       ├── wallet_response.go
│   │       ├── transaction_response.go
│   │       ├── bill_response.go
│   │       └── admin_response.go
│   │
│   ├── pkg/                        # Internal packages (shared)
│   │   ├── db/
│   │   │   ├── postgres.go
│   │   │   └── redis.go
│   │   ├── logger/
│   │   │   └── zap_logger.go
│   │   ├── errors/
│   │   │   ├── app_error.go
│   │   │   ├── error_codes.go
│   │   │   └── error_handler.go
│   │   ├── utils/
│   │   │   ├── otp.go
│   │   │   ├── encryption.go
│   │   │   ├── pagination.go
│   │   │   ├── reference_gen.go
│   │   │   └── face_match.go
│   │   ├── queue/
│   │   │   ├── rabbitmq.go
│   │   │   └── worker.go
│   │   ├── cache/
│   │   │   └── redis_cache.go
│   │   └── webhook/
│   │       └── notifier.go
│   │
│   └── config/
│       └── config.go               # Configuration management
│
├── api/
│   ├── routes/
│   │   ├── routes.go               # Main router setup
│   │   ├── v1/
│   │   │   ├── auth_routes.go
│   │   │   ├── wallet_routes.go
│   │   │   ├── transfer_routes.go
│   │   │   ├── airtime_routes.go
│   │   │   ├── data_routes.go
│   │   │   ├── electricity_routes.go
│   │   │   ├── betting_routes.go
│   │   │   ├── savings_routes.go
│   │   │   ├── notification_routes.go
│   │   │   ├── compliance_routes.go
│   │   │   ├── support_routes.go
│   │   │   └── admin_routes.go
│   │   └── middleware_routes.go
│   │
│   └── docs/
│       ├── swagger.yaml
│       └── postman_collection.json
│
├── migrations/
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   ├── 002_create_wallets_table.up.sql
│   ├── 003_create_transactions_table.up.sql
│   ├── 004_create_bill_payments_table.up.sql
│   ├── 005_create_savings_goals_table.up.sql
│   ├── 006_create_savings_contributions_table.up.sql
│   ├── 007_create_audit_logs_table.up.sql
│   ├── 008_create_support_tickets_table.up.sql
│   ├── 009_create_providers_table.up.sql
│   ├── 010_create_fees_table.up.sql
│   ├── 011_create_kyc_records_table.up.sql
│   ├── 012_create_roles_table.up.sql
│   └── seed.sql
│
├── scripts/
│   ├── build.sh
│   ├── test.sh
│   ├── migrate-up.sh
│   └── migrate-down.sh
│
├── tests/
│   ├── unit/
│   │   ├── services/
│   │   │   ├── auth_service_test.go
│   │   │   ├── wallet_service_test.go
│   │   │   └── transfer_service_test.go
│   │   └── entities/
│   │       └── user_test.go
│   ├── integration/
│   │   ├── auth_test.go
│   │   ├── transfer_test.go
│   │   └── admin_test.go
│   ├── contract/                   # Provider contract tests
│   │   ├── airtime_provider_test.go
│   │   └── bank_nip_test.go
│   └── e2e/
│       └── full_flow_test.go
│
├── deployments/
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   ├── kubernetes/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   ├── configmap.yaml
│   │   └── secrets.yaml
│   └── terraform/
│       └── aws/
│
├── scripts/
│   ├── monitor.sh
│   └── backup.sh
│
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md