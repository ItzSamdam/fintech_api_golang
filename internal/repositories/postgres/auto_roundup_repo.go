package postgres

import (
    "context"
    "errors"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "fintech_api_golang/internal/core/entities"
    "fintech_api_golang/internal/core/interfaces"
)

type autoRoundupRepository struct {
    db *gorm.DB
}

func NewAutoRoundupRepository(db *gorm.DB) interfaces.AutoRoundupRepository {
    return &autoRoundupRepository{db: db}
}

func (r *autoRoundupRepository) Create(ctx context.Context, roundup *entities.AutoRoundup) error {
    return r.db.WithContext(ctx).Create(roundup).Error
}

func (r *autoRoundupRepository) Update(ctx context.Context, roundup *entities.AutoRoundup) error {
    return r.db.WithContext(ctx).Save(roundup).Error
}

func (r *autoRoundupRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.AutoRoundup, error) {
    var roundup entities.AutoRoundup
    err := r.db.WithContext(ctx).First(&roundup, "user_id = ?", userID).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &roundup, nil
}

func (r *autoRoundupRepository) GetActive(ctx context.Context) ([]entities.AutoRoundup, error) {
    var roundups []entities.AutoRoundup
    err := r.db.WithContext(ctx).
        Where("is_active = ?", true).
        Find(&roundups).Error
    return roundups, err
}

func (r *autoRoundupRepository) Deactivate(ctx context.Context, userID uuid.UUID) error {
    return r.db.WithContext(ctx).Model(&entities.AutoRoundup{}).
        Where("user_id = ?", userID).
        Update("is_active", false).Error
}

func (r *autoRoundupRepository) UpdateTotalRoundup(ctx context.Context, userID uuid.UUID, amount int64) error {
    return r.db.WithContext(ctx).Model(&entities.AutoRoundup{}).
        Where("user_id = ?", userID).
        Update("total_roundup", gorm.Expr("total_roundup + ?", amount)).Error
}