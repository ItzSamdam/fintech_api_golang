package interfaces

import (
    "context"
    "time"
    "github.com/google/uuid"
    "fintech_api_golang/internal/core/entities"
)

type BillDetailRepository interface {
    Create(ctx context.Context, detail *entities.BillDetail) error
    GetByTransactionID(ctx context.Context, transactionID uuid.UUID) (*entities.BillDetail, error)
    GetByType(ctx context.Context, userID uuid.UUID, billType string, offset, limit int) ([]entities.BillDetail, int64, error)
    GetByPhoneNumber(ctx context.Context, phoneNumber string, offset, limit int) ([]entities.BillDetail, int64, error)
    GetByMeterNumber(ctx context.Context, meterNumber string, offset, limit int) ([]entities.BillDetail, int64, error)
    UpdateToken(ctx context.Context, transactionID uuid.UUID, token string, units int) error
}

type ProviderRepository interface {
    Create(ctx context.Context, provider *entities.Provider) error
    Update(ctx context.Context, provider *entities.Provider) error
    GetByID(ctx context.Context, id uuid.UUID) (*entities.Provider, error)
    GetByCode(ctx context.Context, code string) (*entities.Provider, error)
    GetByType(ctx context.Context, providerType string) ([]entities.Provider, error)
    GetActiveByType(ctx context.Context, providerType string) ([]entities.Provider, error)
    GetByPriority(ctx context.Context, providerType string) ([]entities.Provider, error)
    ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error
    UpdatePriority(ctx context.Context, id uuid.UUID, priority int) error
    UpdateHealthStatus(ctx context.Context, id uuid.UUID, status string, lastCheck time.Time) error
    UpdateMargin(ctx context.Context, id uuid.UUID, marginPercent float64) error
    List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]entities.Provider, int64, error)
}

type ProviderLogRepository interface {
    Create(ctx context.Context, log *entities.ProviderLog) error
    GetByProviderID(ctx context.Context, providerID uuid.UUID, offset, limit int) ([]entities.ProviderLog, int64, error)
    GetByTransactionID(ctx context.Context, transactionID uuid.UUID) (*entities.ProviderLog, error)
    GetErrorLogs(ctx context.Context, startDate, endDate time.Time) ([]entities.ProviderLog, error)
    GetProviderStats(ctx context.Context, providerID uuid.UUID, startDate, endDate time.Time) (*ProviderStats, error)
}

type ProviderStats struct {
    TotalRequests   int64
    SuccessCount    int64
    FailedCount     int64
    SuccessRate     float64
    AvgResponseTime int64
    TotalRevenue    int64
}