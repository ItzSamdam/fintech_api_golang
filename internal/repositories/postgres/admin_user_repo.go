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

type adminUserRepository struct {
    db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) interfaces.AdminUserRepository {
    return &adminUserRepository{db: db}
}

func (r *adminUserRepository) Create(ctx context.Context, admin *entities.AdminUser) error {
    return r.db.WithContext(ctx).Create(admin).Error
}

func (r *adminUserRepository) Update(ctx context.Context, admin *entities.AdminUser) error {
    return r.db.WithContext(ctx).Save(admin).Error
}

func (r *adminUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&entities.AdminUser{}, "id = ?", id).Error
}

func (r *adminUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.AdminUser, error) {
    var admin entities.AdminUser
    err := r.db.WithContext(ctx).First(&admin, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &admin, nil
}

func (r *adminUserRepository) GetByEmail(ctx context.Context, email string) (*entities.AdminUser, error) {
    var admin entities.AdminUser
    err := r.db.WithContext(ctx).First(&admin, "email = ?", email).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &admin, nil
}

func (r *adminUserRepository) List(ctx context.Context, offset, limit int) ([]entities.AdminUser, int64, error) {
    var admins []entities.AdminUser
    var total int64
    
    if err := r.db.WithContext(ctx).Model(&entities.AdminUser{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Order("created_at DESC").Find(&admins).Error
    return admins, total, err
}

func (r *adminUserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID, ip string) error {
    now := time.Now()
    return r.db.WithContext(ctx).Model(&entities.AdminUser{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "last_login_at": now,
        }).Error
}

func (r *adminUserRepository) UpdateRole(ctx context.Context, id uuid.UUID, role string) error {
    return r.db.WithContext(ctx).Model(&entities.AdminUser{}).
        Where("id = ?", id).
        Update("role", role).Error
}

func (r *adminUserRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Model(&entities.AdminUser{}).
        Where("id = ?", id).
        Update("is_active", false).Error
}

func (r *adminUserRepository) Activate(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Model(&entities.AdminUser{}).
        Where("id = ?", id).
        Update("is_active", true).Error
}