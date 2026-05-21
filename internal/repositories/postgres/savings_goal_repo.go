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

type savingsGoalRepository struct {
    db *gorm.DB
}

func NewSavingsGoalRepository(db *gorm.DB) interfaces.SavingsGoalRepository {
    return &savingsGoalRepository{db: db}
}

func (r *savingsGoalRepository) Create(ctx context.Context, goal *entities.SavingsGoal) error {
    return r.db.WithContext(ctx).Create(goal).Error
}

func (r *savingsGoalRepository) Update(ctx context.Context, goal *entities.SavingsGoal) error {
    return r.db.WithContext(ctx).Save(goal).Error
}

func (r *savingsGoalRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&entities.SavingsGoal{}, "id = ?", id).Error
}

func (r *savingsGoalRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.SavingsGoal, error) {
    var goal entities.SavingsGoal
    err := r.db.WithContext(ctx).First(&goal, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &goal, nil
}

func (r *savingsGoalRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.SavingsGoal, error) {
    var goals []entities.SavingsGoal
    err := r.db.WithContext(ctx).
        Where("user_id = ?", userID).
        Order("created_at DESC").
        Find(&goals).Error
    return goals, err
}

func (r *savingsGoalRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) ([]entities.SavingsGoal, error) {
    var goals []entities.SavingsGoal
    err := r.db.WithContext(ctx).
        Where("user_id = ? AND status = ?", userID, "active").
        Order("target_date ASC").
        Find(&goals).Error
    return goals, err
}

func (r *savingsGoalRepository) GetByStatus(ctx context.Context, status string) ([]entities.SavingsGoal, error) {
    var goals []entities.SavingsGoal
    err := r.db.WithContext(ctx).
        Where("status = ?", status).
        Find(&goals).Error
    return goals, err
}

func (r *savingsGoalRepository) UpdateCurrentAmount(ctx context.Context, goalID uuid.UUID, amount int64) error {
    return r.db.WithContext(ctx).Model(&entities.SavingsGoal{}).
        Where("id = ?", goalID).
        Update("current_amount", entities.AmountInKobo(amount)).Error
}

func (r *savingsGoalRepository) Withdraw(ctx context.Context, goalID uuid.UUID) error {
    now := time.Now()
    return r.db.WithContext(ctx).Model(&entities.SavingsGoal{}).
        Where("id = ?", goalID).
        Updates(map[string]interface{}{
            "status":       "withdrawn",
            "withdrawn_at": now,
        }).Error
}

func (r *savingsGoalRepository) Cancel(ctx context.Context, goalID uuid.UUID) error {
    return r.db.WithContext(ctx).Model(&entities.SavingsGoal{}).
        Where("id = ?", goalID).
        Update("status", "cancelled").Error
}

func (r *savingsGoalRepository) GetAutoDebitGoals(ctx context.Context) ([]entities.SavingsGoal, error) {
    var goals []entities.SavingsGoal
    err := r.db.WithContext(ctx).
        Where("is_auto_debit = ? AND status = ?", true, "active").
        Find(&goals).Error
    return goals, err
}

func (r *savingsGoalRepository) GetCompletedGoals(ctx context.Context) ([]entities.SavingsGoal, error) {
    var goals []entities.SavingsGoal
    err := r.db.WithContext(ctx).
        Where("status = ? OR current_amount >= target_amount", "active").
        Find(&goals).Error
    return goals, err
}