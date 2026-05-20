package postgres

import (
    "context"
    "time"
    "errors"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
    "fintech_api_golang/internal/core/entities"
    "fintech_api_golang/internal/core/interfaces"
)

type walletRepository struct {
    db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) interfaces.WalletRepository {
    return &walletRepository{db: db}
}

func (r *walletRepository) Create(ctx context.Context, wallet *entities.Wallet) error {
    return r.db.WithContext(ctx).Create(wallet).Error
}

func (r *walletRepository) Update(ctx context.Context, wallet *entities.Wallet) error {
    return r.db.WithContext(ctx).Save(wallet).Error
}

func (r *walletRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Wallet, error) {
    var wallet entities.Wallet
    err := r.db.WithContext(ctx).First(&wallet, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &wallet, nil
}

func (r *walletRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.Wallet, error) {
    var wallet entities.Wallet
    err := r.db.WithContext(ctx).First(&wallet, "user_id = ?", userID).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &wallet, nil
}

func (r *walletRepository) GetByUserIDForUpdate(ctx context.Context, userID uuid.UUID) (*entities.Wallet, error) {
    var wallet entities.Wallet
    err := r.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
        First(&wallet, "user_id = ?", userID).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &wallet, nil
}

func (r *walletRepository) Debit(ctx context.Context, walletID uuid.UUID, amount int64, reference string) error {
    return r.db.WithContext(ctx).Model(&entities.Wallet{}).
        Where("id = ? AND balance >= ?", walletID, amount).
        Update("balance", gorm.Expr("balance - ?", amount)).Error
}

func (r *walletRepository) Credit(ctx context.Context, walletID uuid.UUID, amount int64, reference string) error {
    return r.db.WithContext(ctx).Model(&entities.Wallet{}).
        Where("id = ?", walletID).
        Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *walletRepository) Lock(ctx context.Context, walletID uuid.UUID, reason string) error {
    now := time.Now()
    return r.db.WithContext(ctx).Model(&entities.Wallet{}).
        Where("id = ?", walletID).
        Updates(map[string]interface{}{
            "is_locked":  true,
            "locked_at":  now,
            "lock_reason": reason,
        }).Error
}

func (r *walletRepository) Unlock(ctx context.Context, walletID uuid.UUID) error {
    return r.db.WithContext(ctx).Model(&entities.Wallet{}).
        Where("id = ?", walletID).
        Updates(map[string]interface{}{
            "is_locked":   false,
            "locked_at":   nil,
            "lock_reason": nil,
        }).Error
}

func (r *walletRepository) UpdateSpentLimits(ctx context.Context, walletID uuid.UUID, amount int64) error {
    return r.db.WithContext(ctx).Model(&entities.Wallet{}).
        Where("id = ?", walletID).
        Updates(map[string]interface{}{
            "daily_spent":   gorm.Expr("daily_spent + ?", amount),
            "weekly_spent":  gorm.Expr("weekly_spent + ?", amount),
            "monthly_spent": gorm.Expr("monthly_spent + ?", amount),
        }).Error
}

func (r *walletRepository) ResetDailySpent(ctx context.Context) error {
    return r.db.WithContext(ctx).Model(&entities.Wallet{}).
        Where("last_daily_reset < ?", time.Now().Truncate(24*time.Hour)).
        Updates(map[string]interface{}{
            "daily_spent":    0,
            "last_daily_reset": time.Now(),
        }).Error
}

func (r *walletRepository) ResetWeeklySpent(ctx context.Context) error {
    // Get start of week (Monday)
    now := time.Now()
    weekday := int(now.Weekday())
    if weekday == 0 {
        weekday = 7
    }
    startOfWeek := now.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
    
    return r.db.WithContext(ctx).Model(&entities.Wallet{}).
        Where("last_weekly_reset < ?", startOfWeek).
        Updates(map[string]interface{}{
            "weekly_spent":    0,
            "last_weekly_reset": time.Now(),
        }).Error
}

func (r *walletRepository) ResetMonthlySpent(ctx context.Context) error {
    startOfMonth := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -time.Now().Day()+1)
    return r.db.WithContext(ctx).Model(&entities.Wallet{}).
        Where("last_monthly_reset < ?", startOfMonth).
        Updates(map[string]interface{}{
            "monthly_spent":    0,
            "last_monthly_reset": time.Now(),
        }).Error
}

func (r *walletRepository) GetTotalBalance(ctx context.Context) (int64, error) {
    var total int64
    err := r.db.WithContext(ctx).Model(&entities.Wallet{}).
        Select("COALESCE(SUM(balance), 0)").
        Scan(&total).Error
    return total, err
}

func (r *walletRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]entities.Wallet, int64, error) {
    var wallets []entities.Wallet
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entities.Wallet{})
    
    if isLocked, ok := filters["is_locked"]; ok {
        query = query.Where("is_locked = ?", isLocked)
    }
    if currency, ok := filters["currency"]; ok {
        query = query.Where("currency = ?", currency)
    }
    
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := query.Offset(offset).Limit(limit).Preload("User").Find(&wallets).Error
    return wallets, total, err
}