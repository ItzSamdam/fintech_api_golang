package postgres

import (
    "context"
    "errors"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "fintech_api_golang/internal/core/entities"
    "fintech_api_golang/internal/core/interfaces"
)

type billDetailRepository struct {
    db *gorm.DB
}

func NewBillDetailRepository(db *gorm.DB) interfaces.BillDetailRepository {
    return &billDetailRepository{db: db}
}

func (r *billDetailRepository) Create(ctx context.Context, detail *entities.BillDetail) error {
    return r.db.WithContext(ctx).Create(detail).Error
}

func (r *billDetailRepository) GetByTransactionID(ctx context.Context, transactionID uuid.UUID) (*entities.BillDetail, error) {
    var detail entities.BillDetail
    err := r.db.WithContext(ctx).First(&detail, "transaction_id = ?", transactionID).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &detail, nil
}

func (r *billDetailRepository) GetByType(ctx context.Context, userID uuid.UUID, billType string, offset, limit int) ([]entities.BillDetail, int64, error) {
    var details []entities.BillDetail
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entities.BillDetail{}).
        Joins("JOIN transactions ON transactions.id = bill_details.transaction_id").
        Where("transactions.user_id = ? AND bill_details.bill_type = ?", userID, billType)
    
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := query.Offset(offset).Limit(limit).Order("bill_details.created_at DESC").Find(&details).Error
    return details, total, err
}

func (r *billDetailRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string, offset, limit int) ([]entities.BillDetail, int64, error) {
    var details []entities.BillDetail
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entities.BillDetail{}).
        Where("phone_number = ?", phoneNumber)
    
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&details).Error
    return details, total, err
}

func (r *billDetailRepository) GetByMeterNumber(ctx context.Context, meterNumber string, offset, limit int) ([]entities.BillDetail, int64, error) {
    var details []entities.BillDetail
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entities.BillDetail{}).
        Where("meter_number = ?", meterNumber)
    
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&details).Error
    return details, total, err
}

func (r *billDetailRepository) UpdateToken(ctx context.Context, transactionID uuid.UUID, token string, units int) error {
    return r.db.WithContext(ctx).Model(&entities.BillDetail{}).
        Where("transaction_id = ?", transactionID).
        Updates(map[string]interface{}{
            "electricity_token": token,
            "electricity_units": units,
        }).Error
}