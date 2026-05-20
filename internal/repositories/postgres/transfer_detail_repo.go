package postgres

import (
    "context"
    "errors"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "fintech_api_golang/internal/core/entities"
    "fintech_api_golang/internal/core/interfaces"
)

type transferDetailRepository struct {
    db *gorm.DB
}

func NewTransferDetailRepository(db *gorm.DB) interfaces.TransferDetailRepository {
    return &transferDetailRepository{db: db}
}

func (r *transferDetailRepository) Create(ctx context.Context, detail *entities.TransferDetail) error {
    return r.db.WithContext(ctx).Create(detail).Error
}

func (r *transferDetailRepository) GetByTransactionID(ctx context.Context, transactionID uuid.UUID) (*entities.TransferDetail, error) {
    var detail entities.TransferDetail
    err := r.db.WithContext(ctx).First(&detail, "transaction_id = ?", transactionID).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &detail, nil
}

func (r *transferDetailRepository) GetByRecipientID(ctx context.Context, recipientID string, offset, limit int) ([]entities.TransferDetail, int64, error) {
    var details []entities.TransferDetail
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entities.TransferDetail{}).
        Where("recipient_id = ?", recipientID)
    
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&details).Error
    return details, total, err
}