package interfaces

import (
    "context"
    "time"
    "github.com/google/uuid"
    "fintech_api_golang/internal/core/entities"
)

type AdminUserRepository interface {
    Create(ctx context.Context, admin *entities.AdminUser) error
    Update(ctx context.Context, admin *entities.AdminUser) error
    Delete(ctx context.Context, id uuid.UUID) error
    GetByID(ctx context.Context, id uuid.UUID) (*entities.AdminUser, error)
    GetByEmail(ctx context.Context, email string) (*entities.AdminUser, error)
    List(ctx context.Context, offset, limit int) ([]entities.AdminUser, int64, error)
    UpdateLastLogin(ctx context.Context, id uuid.UUID, ip string) error
    UpdateRole(ctx context.Context, id uuid.UUID, role string) error
    Deactivate(ctx context.Context, id uuid.UUID) error
    Activate(ctx context.Context, id uuid.UUID) error
}

type RoleRepository interface {
    Create(ctx context.Context, role *entities.Role) error
    Update(ctx context.Context, role *entities.Role) error
    Delete(ctx context.Context, id uuid.UUID) error
    GetByID(ctx context.Context, id uuid.UUID) (*entities.Role, error)
    GetByName(ctx context.Context, name string) (*entities.Role, error)
    List(ctx context.Context) ([]entities.Role, error)
    GetPermissions(ctx context.Context, roleName string) ([]string, error)
}

type AuditLogRepository interface {
    Create(ctx context.Context, log *entities.AuditLog) error
    GetByID(ctx context.Context, id uuid.UUID) (*entities.AuditLog, error)
    GetByAdminID(ctx context.Context, adminID uuid.UUID, offset, limit int) ([]entities.AuditLog, int64, error)
    GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]entities.AuditLog, int64, error)
    GetByAction(ctx context.Context, action string, startDate, endDate time.Time) ([]entities.AuditLog, error)
    List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]entities.AuditLog, int64, error)
    GetAdminActions(ctx context.Context, adminID uuid.UUID, startDate, endDate time.Time) ([]entities.AuditLog, error)
}