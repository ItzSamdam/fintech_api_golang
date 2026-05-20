package interfaces

import (
    "context"
    "time"
    "github.com/google/uuid"
    "fintech_api_golang/internal/core/entities"
)

type WalletRepository interface {
    Create(ctx context.Context, wallet *entities.Wallet) error
    Update(ctx context.Context, wallet *entities.Wallet) error
    GetByID(ctx context.Context, id uuid.UUID) (*entities.Wallet, error)
    GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.Wallet, error)
    GetByUserIDForUpdate(ctx context.Context, userID uuid.UUID) (*entities.Wallet, error)
    Debit(ctx context.Context, walletID uuid.UUID, amount int64, reference string) error
    Credit(ctx context.Context, walletID uuid.UUID, amount int64, reference string) error
    Lock(ctx context.Context, walletID uuid.UUID, reason string) error
    Unlock(ctx context.Context, walletID uuid.UUID) error
    UpdateSpentLimits(ctx context.Context, walletID uuid.UUID, amount int64) error
    ResetDailySpent(ctx context.Context) error
    ResetWeeklySpent(ctx context.Context) error
    ResetMonthlySpent(ctx context.Context) error
    GetTotalBalance(ctx context.Context) (int64, error)
    List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]entities.Wallet, int64, error)
}

type TransactionRepository interface {
    Create(ctx context.Context, transaction *entities.Transaction) error
    Update(ctx context.Context, transaction *entities.Transaction) error
    GetByID(ctx context.Context, id uuid.UUID) (*entities.Transaction, error)
    GetByReference(ctx context.Context, reference string) (*entities.Transaction, error)
    GetByWalletID(ctx context.Context, walletID uuid.UUID, offset, limit int, filters map[string]interface{}) ([]entities.Transaction, int64, error)
    GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int, filters map[string]interface{}) ([]entities.Transaction, int64, error)
    GetByCategory(ctx context.Context, userID uuid.UUID, category string, offset, limit int) ([]entities.Transaction, int64, error)
    UpdateStatus(ctx context.Context, reference string, status string, completedAt *time.Time) error
    Reverse(ctx context.Context, reference string, reversedTxnID uuid.UUID) error
    GetSummaryByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*TransactionSummary, error)
    GetAdminSummary(ctx context.Context, startDate, endDate time.Time) (*AdminTransactionSummary, error)
    GetDailyVolume(ctx context.Context, date time.Time) (int64, error)
    GetPendingTransactions(ctx context.Context) ([]entities.Transaction, error)
    MarkAsFailed(ctx context.Context, reference string, response string) error
}

type TransactionSummary struct {
    TotalCount      int64
    TotalAmount     int64
    TotalFee        int64
    TotalVAT        int64
    SuccessCount    int64
    FailedCount     int64
    PendingCount    int64
}

type AdminTransactionSummary struct {
    TotalVolume     int64
    TotalRevenue    int64
    TotalFee        int64
    TotalVAT        int64
    SuccessRate     float64
}

type TransferDetailRepository interface {
    Create(ctx context.Context, detail *entities.TransferDetail) error
    GetByTransactionID(ctx context.Context, transactionID uuid.UUID) (*entities.TransferDetail, error)
    GetByRecipientID(ctx context.Context, recipientID string, offset, limit int) ([]entities.TransferDetail, int64, error)
}