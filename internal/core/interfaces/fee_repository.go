package interfaces

import (
    "context"
    "fintech_api_golang/internal/core/entities"
)

type FeeConfigRepository interface {
    Create(ctx context.Context, fee *entities.FeeConfig) error
    Update(ctx context.Context, fee *entities.FeeConfig) error
    GetByBillType(ctx context.Context, billType string) (*entities.FeeConfig, error)
    GetActiveFee(ctx context.Context, billType string) (*entities.FeeConfig, error)
    List(ctx context.Context, offset, limit int) ([]entities.FeeConfig, int64, error)
    CalculateFee(ctx context.Context, billType string, amount int64) (*FeeCalculation, error)
}

type FeeCalculation struct {
    Fee         int64
    VAT         int64
    TotalCharge int64
}

type TierLimitRepository interface {
    Create(ctx context.Context, limit *entities.TierLimit) error
    Update(ctx context.Context, limit *entities.TierLimit) error
    GetByTier(ctx context.Context, tier int) (*entities.TierLimit, error)
    List(ctx context.Context) ([]entities.TierLimit, error)
}