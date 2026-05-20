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

type auditLogRepository struct {
    db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) interfaces.AuditLogRepository {
    return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, log *entities.AuditLog) error {
    return r.db.WithContext(ctx).Create(log).Error
}

func (r *auditLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.AuditLog, error) {
    var log entities.AuditLog
    err := r.db.WithContext(ctx).First(&log, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &log, nil
}

func (r *auditLogRepository) GetByAdminID(ctx context.Context, adminID uuid.UUID, offset, limit int) ([]entities.AuditLog, int64, error) {
    var logs []entities.AuditLog
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entities.AuditLog{}).
        Where("admin_id = ?", adminID)
    
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
    return logs, total, err
}

func (r *auditLogRepository) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]entities.AuditLog, int64, error) {
    var logs []entities.AuditLog
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entities.AuditLog{}).
        Where("user_id = ?", userID)
    
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
    return logs, total, err
}

func (r *auditLogRepository) GetByAction(ctx context.Context, action string, startDate, endDate time.Time) ([]entities.AuditLog, error) {
    var logs []entities.AuditLog
    err := r.db.WithContext(ctx).
        Where("action = ? AND created_at BETWEEN ? AND ?", action, startDate, endDate).
        Order("created_at DESC").
        Find(&logs).Error
    return logs, err
}

func (r *auditLogRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]entities.AuditLog, int64, error) {
    var logs []entities.AuditLog
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entities.AuditLog{})
    
    if adminID, ok := filters["admin_id"]; ok {
        query = query.Where("admin_id = ?", adminID)
    }
    if action, ok := filters["action"]; ok {
        query = query.Where("action = ?", action)
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
    
    err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
    return logs, total, err
}

func (r *auditLogRepository) GetAdminActions(ctx context.Context, adminID uuid.UUID, startDate, endDate time.Time) ([]entities.AuditLog, error) {
    var logs []entities.AuditLog
    err := r.db.WithContext(ctx).
        Where("admin_id = ? AND created_at BETWEEN ? AND ?", adminID, startDate, endDate).
        Order("created_at DESC").
        Find(&logs).Error
    return logs, err
}