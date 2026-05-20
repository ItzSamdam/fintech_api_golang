package postgres

import (
    "context"
    "time"
    "errors"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "fintech_api_golang/internal/core/entities"
    "fintech_api_golang/internal/core/interfaces"
)

type transactionRepository struct {
    db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) interfaces.TransactionRepository {
    return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *entities.Transaction) error {
    return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *transactionRepository) Update(ctx context.Context, transaction *entities.Transaction) error {
    return r.db.WithContext(ctx).Save(transaction).Error
}

func (r *transactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Transaction, error) {
    var transaction entities.Transaction
    err := r.db.WithContext(ctx).
        Preload("TransferDetail").
        Preload("BillDetail").
        First(&transaction, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &transaction, nil
}

func (r *transactionRepository) GetByReference(ctx context.Context, reference string) (*entities.Transaction, error) {
    var transaction entities.Transaction
    err := r.db.WithContext(ctx).
        Preload("TransferDetail").
        Preload("BillDetail").
        First(&transaction, "reference = ?", reference).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &transaction, nil
}

func (r *transactionRepository) GetByWalletID(ctx context.Context, walletID uuid.UUID, offset, limit int, filters map[string]interface{}) ([]entities.Transaction, int64, error) {
    var transactions []entities.Transaction
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entities.Transaction{}).
        Where("wallet_id = ?", walletID)
    
    if category, ok := filters["category"]; ok {
        query = query.Where("category = ?", category)
    }
    if status, ok := filters["status"]; ok {
        query = query.Where("status = ?", status)
    }
    if fromDate, ok := filters["from_date"]; ok {
        query = query.Where("created_at >= ?", fromDate)
    }
    if toDate, ok := filters["to_date"]; ok {
        query = query.Where("created_at <= ?", toDate)
    }
    
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := query.Offset(offset).Limit(limit).
        Order("created_at DESC").
        Preload("TransferDetail").
        Preload("BillDetail").
        Find(&transactions).Error
    return transactions, total, err
}

func (r *transactionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int, filters map[string]interface{}) ([]entities.Transaction, int64, error) {
    var transactions []entities.Transaction
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entities.Transaction{}).
        Where("user_id = ?", userID)
    
    if category, ok := filters["category"]; ok {
        query = query.Where("category = ?", category)
    }
    if status, ok := filters["status"]; ok {
        query = query.Where("status = ?", status)
    }
    if fromDate, ok := filters["from_date"]; ok {
        query = query.Where("created_at >= ?", fromDate)
    }
    if toDate, ok := filters["to_date"]; ok {
        query = query.Where("created_at <= ?", toDate)
    }
    
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := query.Offset(offset).Limit(limit).
        Order("created_at DESC").
        Preload("TransferDetail").
        Preload("BillDetail").
        Find(&transactions).Error
    return transactions, total, err
}

func (r *transactionRepository) GetByCategory(ctx context.Context, userID uuid.UUID, category string, offset, limit int) ([]entities.Transaction, int64, error) {
    var transactions []entities.Transaction
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entities.Transaction{}).
        Where("user_id = ? AND category = ?", userID, category)
    
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := query.Offset(offset).Limit(limit).
        Order("created_at DESC").
        Find(&transactions).Error
    return transactions, total, err
}

func (r *transactionRepository) UpdateStatus(ctx context.Context, reference string, status string, completedAt *time.Time) error {
    updates := map[string]interface{}{
        "status": status,
    }
    if completedAt != nil {
        updates["completed_at"] = completedAt
    }
    return r.db.WithContext(ctx).Model(&entities.Transaction{}).
        Where("reference = ?", reference).
        Updates(updates).Error
}

func (r *transactionRepository) Reverse(ctx context.Context, reference string, reversedTxnID uuid.UUID) error {
    return r.db.WithContext(ctx).Model(&entities.Transaction{}).
        Where("reference = ?", reference).
        Updates(map[string]interface{}{
            "is_reversed":    true,
            "reversed_txn_id": reversedTxnID,
            "status":         "reversed",
        }).Error
}

func (r *transactionRepository) GetSummaryByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*interfaces.TransactionSummary, error) {
    var summary interfaces.TransactionSummary
    
    err := r.db.WithContext(ctx).Model(&entities.Transaction{}).
        Where("user_id = ? AND created_at BETWEEN ? AND ? AND status = ?", userID, startDate, endDate, "success").
        Select("COUNT(*) as total_count, COALESCE(SUM(amount), 0) as total_amount, COALESCE(SUM(fee), 0) as total_fee, COALESCE(SUM(vat), 0) as total_vat").
        Scan(&summary).Error
    
    if err != nil {
        return nil, err
    }
    
    // Get counts by status
    r.db.WithContext(ctx).Model(&entities.Transaction{}).
        Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
        Select("status, COUNT(*) as count").
        Group("status").
        Scan(&summary)
    
    return &summary, nil
}

func (r *transactionRepository) GetAdminSummary(ctx context.Context, startDate, endDate time.Time) (*interfaces.AdminTransactionSummary, error) {
    var summary interfaces.AdminTransactionSummary
    
    err := r.db.WithContext(ctx).Model(&entities.Transaction{}).
        Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, "success").
        Select("COALESCE(SUM(total_amount), 0) as total_volume, COALESCE(SUM(fee + vat), 0) as total_revenue").
        Scan(&summary).Error
    
    return &summary, err
}

func (r *transactionRepository) GetDailyVolume(ctx context.Context, date time.Time) (int64, error) {
    var volume int64
    startOfDay := date.Truncate(24 * time.Hour)
    endOfDay := startOfDay.Add(24 * time.Hour)
    
    err := r.db.WithContext(ctx).Model(&entities.Transaction{}).
        Where("created_at BETWEEN ? AND ? AND status = ?", startOfDay, endOfDay, "success").
        Select("COALESCE(SUM(amount), 0)").
        Scan(&volume).Error
    return volume, err
}

func (r *transactionRepository) GetPendingTransactions(ctx context.Context) ([]entities.Transaction, error) {
    var transactions []entities.Transaction
    err := r.db.WithContext(ctx).
        Where("status = ? AND created_at < ?", "pending", time.Now().Add(-5*time.Minute)).
        Find(&transactions).Error
    return transactions, err
}

func (r *transactionRepository) MarkAsFailed(ctx context.Context, reference string, response string) error {
    return r.db.WithContext(ctx).Model(&entities.Transaction{}).
        Where("reference = ?", reference).
        Updates(map[string]interface{}{
            "status":           "failed",
            "provider_response": response,
        }).Error
}