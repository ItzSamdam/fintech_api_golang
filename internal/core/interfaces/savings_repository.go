package interfaces

import (
    "context"
    "time"
    "github.com/google/uuid"
    "fintech_api_golang/internal/core/entities"
)

type SavingsGoalRepository interface {
    Create(ctx context.Context, goal *entities.SavingsGoal) error
    Update(ctx context.Context, goal *entities.SavingsGoal) error
    Delete(ctx context.Context, id uuid.UUID) error
    GetByID(ctx context.Context, id uuid.UUID) (*entities.SavingsGoal, error)
    GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.SavingsGoal, error)
    GetActiveByUserID(ctx context.Context, userID uuid.UUID) ([]entities.SavingsGoal, error)
    GetByStatus(ctx context.Context, status string) ([]entities.SavingsGoal, error)
    UpdateCurrentAmount(ctx context.Context, goalID uuid.UUID, amount int64) error
    Withdraw(ctx context.Context, goalID uuid.UUID) error
    Cancel(ctx context.Context, goalID uuid.UUID) error
    GetAutoDebitGoals(ctx context.Context) ([]entities.SavingsGoal, error)
    GetCompletedGoals(ctx context.Context) ([]entities.SavingsGoal, error)
}

type SavingsContributionRepository interface {
    Create(ctx context.Context, contribution *entities.SavingsContribution) error
    GetByGoalID(ctx context.Context, goalID uuid.UUID, offset, limit int) ([]entities.SavingsContribution, int64, error)
    GetByTransactionID(ctx context.Context, transactionID uuid.UUID) (*entities.SavingsContribution, error)
    GetTotalContributions(ctx context.Context, goalID uuid.UUID) (int64, error)
    GetContributionsByDateRange(ctx context.Context, goalID uuid.UUID, startDate, endDate time.Time) ([]entities.SavingsContribution, error)
}

type AutoRoundupRepository interface {
    Create(ctx context.Context, roundup *entities.AutoRoundup) error
    Update(ctx context.Context, roundup *entities.AutoRoundup) error
    GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.AutoRoundup, error)
    GetActive(ctx context.Context) ([]entities.AutoRoundup, error)
    Deactivate(ctx context.Context, userID uuid.UUID) error
    UpdateTotalRoundup(ctx context.Context, userID uuid.UUID, amount int64) error
}