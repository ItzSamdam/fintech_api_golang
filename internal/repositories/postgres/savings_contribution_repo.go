package postgres

import (
    "context"
    "errors"
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "fintech_api_golang/internal/core/entities"
    "fintech_api_golang/internal/core/interfaces"
)

type savingsContributionRepository struct {
    db *gorm.DB
}

func NewSavingsContributionRepository(db *gorm.DB) interfaces.SavingsContributionRepository {
    return &savingsContributionRepository{db: db}
}

func (r *savingsContributionRepository) Create(ctx context.Context, contribution *entities.SavingsContribution) error {
    return r.db.WithContext(ctx).Create(contribution).Error
}

func (r *savingsContributionRepository) GetByGoalID(ctx context.Context, goalID uuid.UUID, offset, limit int) ([]entities.SavingsContribution, int64, error) {
    var contributions []entities.SavingsContribution
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entities.SavingsContribution{}).
        Where("savings_goal_id = ?", goalID)
    
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := query.Offset(offset).Limit(limit).Order("contribution_date DESC").Find(&contributions).Error
    return contributions, total, err
}

func (r *savingsContributionRepository) GetByTransactionID(ctx context.Context, transactionID uuid.UUID) (*entities.SavingsContribution, error) {
    var contribution entities.SavingsContribution
    err := r.db.WithContext(ctx).First(&contribution, "transaction_id = ?", transactionID).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &contribution, nil
}

func (r *savingsContributionRepository) GetTotalContributions(ctx context.Context, goalID uuid.UUID) (int64, error) {
    var total int64
    err := r.db.WithContext(ctx).Model(&entities.SavingsContribution{}).
        Where("savings_goal_id = ?", goalID).
        Select("COALESCE(SUM(amount), 0)").
        Scan(&total).Error
    return total, err
}

func (r *savingsContributionRepository) GetContributionsByDateRange(ctx context.Context, goalID uuid.UUID, startDate, endDate time.Time) ([]entities.SavingsContribution, error) {
    var contributions []entities.SavingsContribution
    err := r.db.WithContext(ctx).
        Where("savings_goal_id = ? AND contribution_date BETWEEN ? AND ?", goalID, startDate, endDate).
        Order("contribution_date ASC").
        Find(&contributions).Error
    return contributions, err
}