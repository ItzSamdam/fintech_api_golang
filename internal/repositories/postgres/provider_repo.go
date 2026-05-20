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

type providerRepository struct {
    db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) interfaces.ProviderRepository {
    return &providerRepository{db: db}
}

func (r *providerRepository) Create(ctx context.Context, provider *entities.Provider) error {
    return r.db.WithContext(ctx).Create(provider).Error
}

func (r *providerRepository) Update(ctx context.Context, provider *entities.Provider) error {
    return r.db.WithContext(ctx).Save(provider).Error
}

func (r *providerRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Provider, error) {
    var provider entities.Provider
    err := r.db.WithContext(ctx).First(&provider, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &provider, nil
}

func (r *providerRepository) GetByCode(ctx context.Context, code string) (*entities.Provider, error) {
    var provider entities.Provider
    err := r.db.WithContext(ctx).First(&provider, "code = ?", code).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &provider, nil
}

func (r *providerRepository) GetByType(ctx context.Context, providerType string) ([]entities.Provider, error) {
    var providers []entities.Provider
    err := r.db.WithContext(ctx).Where("type = ?", providerType).Find(&providers).Error
    return providers, err
}

func (r *providerRepository) GetActiveByType(ctx context.Context, providerType string) ([]entities.Provider, error) {
    var providers []entities.Provider
    err := r.db.WithContext(ctx).
        Where("type = ? AND is_active = ?", providerType, true).
        Order("priority ASC").
        Find(&providers).Error
    return providers, err
}

func (r *providerRepository) GetByPriority(ctx context.Context, providerType string) ([]entities.Provider, error) {
    var providers []entities.Provider
    err := r.db.WithContext(ctx).
        Where("type = ? AND is_active = ?", providerType, true).
        Order("priority ASC").
        Find(&providers).Error
    return providers, err
}

func (r *providerRepository) ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error {
    return r.db.WithContext(ctx).Model(&entities.Provider{}).
        Where("id = ?", id).
        Update("is_active", isActive).Error
}

func (r *providerRepository) UpdatePriority(ctx context.Context, id uuid.UUID, priority int) error {
    return r.db.WithContext(ctx).Model(&entities.Provider{}).
        Where("id = ?", id).
        Update("priority", priority).Error
}

func (r *providerRepository) UpdateHealthStatus(ctx context.Context, id uuid.UUID, status string, lastCheck time.Time) error {
    return r.db.WithContext(ctx).Model(&entities.Provider{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "health_status":     status,
            "last_health_check": lastCheck,
        }).Error
}

func (r *providerRepository) UpdateMargin(ctx context.Context, id uuid.UUID, marginPercent float64) error {
    return r.db.WithContext(ctx).Model(&entities.Provider{}).
        Where("id = ?", id).
        Update("margin_percent", marginPercent).Error
}

func (r *providerRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]entities.Provider, int64, error) {
    var providers []entities.Provider
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entities.Provider{})
    
    if providerType, ok := filters["type"]; ok {
        query = query.Where("type = ?", providerType)
    }
    if isActive, ok := filters["is_active"]; ok {
        query = query.Where("is_active = ?", isActive)
    }
    
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := query.Offset(offset).Limit(limit).Order("priority ASC").Find(&providers).Error
    return providers, total, err
}