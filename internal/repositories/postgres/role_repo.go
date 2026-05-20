package postgres

import (
    "context"
    "errors"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "fintech_api_golang/internal/core/entities"
    "fintech_api_golang/internal/core/interfaces"
)

type roleRepository struct {
    db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) interfaces.RoleRepository {
    return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *entities.Role) error {
    return r.db.WithContext(ctx).Create(role).Error
}

func (r *roleRepository) Update(ctx context.Context, role *entities.Role) error {
    return r.db.WithContext(ctx).Save(role).Error
}

func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&entities.Role{}, "id = ?", id).Error
}

func (r *roleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Role, error) {
    var role entities.Role
    err := r.db.WithContext(ctx).First(&role, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &role, nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*entities.Role, error) {
    var role entities.Role
    err := r.db.WithContext(ctx).First(&role, "name = ?", name).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &role, nil
}

func (r *roleRepository) List(ctx context.Context) ([]entities.Role, error) {
    var roles []entities.Role
    err := r.db.WithContext(ctx).Find(&roles).Error
    return roles, err
}

func (r *roleRepository) GetPermissions(ctx context.Context, roleName string) ([]string, error) {
    var role entities.Role
    err := r.db.WithContext(ctx).First(&role, "name = ?", roleName).Error
    if err != nil {
        return nil, err
    }
    
    // Parse JSON permissions
    var permissions []string
    if role.Permissions != "" {
        // Parse JSON string to []string
        // You'll need to import encoding/json here
    }
    
    return permissions, nil
}