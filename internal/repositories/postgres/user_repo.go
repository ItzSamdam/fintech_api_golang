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

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&entities.User{}, "id = ?", id).Error
}

func (r *userRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&entities.User{}, "id = ?", id).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
    var user entities.User
    err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.User, error) {
    var user entities.User
    err := r.db.WithContext(ctx).First(&user, "phone_number = ?", phoneNumber).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
    var user entities.User
    err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) GetByBVN(ctx context.Context, bvn string) (*entities.User, error) {
    var user entities.User
    err := r.db.WithContext(ctx).First(&user, "bvn = ?", bvn).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) GetByNIN(ctx context.Context, nin string) (*entities.User, error) {
    var user entities.User
    err := r.db.WithContext(ctx).First(&user, "nin = ?", nin).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]entities.User, int64, error) {
    var users []entities.User
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entities.User{})
    
    // Apply filters
    if tier, ok := filters["tier"]; ok {
        query = query.Where("tier = ?", tier)
    }
    if isActive, ok := filters["is_active"]; ok {
        query = query.Where("is_active = ?", isActive)
    }
    if isSuspended, ok := filters["is_suspended"]; ok {
        query = query.Where("is_suspended = ?", isSuspended)
    }
    if fromDate, ok := filters["from_date"]; ok {
        query = query.Where("created_at >= ?", fromDate)
    }
    if toDate, ok := filters["to_date"]; ok {
        query = query.Where("created_at <= ?", toDate)
    }
    
    // Get total count
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Get paginated results
    err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
    return users, total, err
}

func (r *userRepository) UpdateTier(ctx context.Context, userID uuid.UUID, tier int) error {
    return r.db.WithContext(ctx).Model(&entities.User{}).
        Where("id = ?", userID).
        Update("tier", tier).Error
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID, ip string) error {
    now := time.Now()
    return r.db.WithContext(ctx).Model(&entities.User{}).
        Where("id = ?", userID).
        Updates(map[string]interface{}{
            "last_login_at": now,
            "last_login_ip": ip,
        }).Error
}

func (r *userRepository) Suspend(ctx context.Context, userID uuid.UUID, reason string, duration *time.Duration) error {
    now := time.Now()
    updates := map[string]interface{}{
        "is_suspended":     true,
        "suspension_reason": reason,
        "suspended_at":     now,
    }
    return r.db.WithContext(ctx).Model(&entities.User{}).
        Where("id = ?", userID).
        Updates(updates).Error
}

func (r *userRepository) Unsuspend(ctx context.Context, userID uuid.UUID) error {
    return r.db.WithContext(ctx).Model(&entities.User{}).
        Where("id = ?", userID).
        Updates(map[string]interface{}{
            "is_suspended":     false,
            "suspension_reason": nil,
            "suspended_at":     nil,
        }).Error
}

func (r *userRepository) Search(ctx context.Context, query string, offset, limit int) ([]entities.User, int64, error) {
    var users []entities.User
    var total int64
    
    searchQuery := "%" + query + "%"
    dbQuery := r.db.WithContext(ctx).Model(&entities.User{}).
        Where("phone_number ILIKE ? OR email ILIKE ?", searchQuery, searchQuery)
    
    if err := dbQuery.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := dbQuery.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
    return users, total, err
}

func (r *userRepository) CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).Model(&entities.User{}).
        Where("created_at BETWEEN ? AND ?", startDate, endDate).
        Count(&count).Error
    return count, err
}

func (r *userRepository) GetActiveUsers(ctx context.Context) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).Model(&entities.User{}).
        Where("is_active = ? AND is_suspended = ?", true, false).
        Count(&count).Error
    return count, err
}