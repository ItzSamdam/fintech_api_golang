package postgres

import (
	"context"
	"errors"
	"fintech_api_golang/internal/core/entities"
	"fintech_api_golang/internal/core/interfaces"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type kycRepository struct {
	db *gorm.DB
}

func NewKYCRepository(db *gorm.DB) interfaces.KYCRepository {
	return &kycRepository{db: db}
}

func (r *kycRepository) Create(ctx context.Context, kyc *entities.KYC) error {
	return r.db.WithContext(ctx).Create(kyc).Error
}

func (r *kycRepository) Update(ctx context.Context, kyc *entities.KYC) error {
	return r.db.WithContext(ctx).Save(kyc).Error
}

func (r *kycRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.KYC, error) {
	var kyc entities.KYC
	err := r.db.WithContext(ctx).First(&kyc, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &kyc, nil
}

func (r *kycRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.KYC, error) {
	var kyc entities.KYC
	err := r.db.WithContext(ctx).First(&kyc, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &kyc, nil
}

func (r *kycRepository) GetPending(ctx context.Context, offset, limit int) ([]entities.KYC, int64, error) {
	var kycList []entities.KYC
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.KYC{}).
		Where("status = ?", "pending")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).
		Preload("User").
		Order("created_at ASC").
		Find(&kycList).Error
	return kycList, total, err
}

func (r *kycRepository) Approve(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entities.KYC{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      "approved",
			"approved_by": approvedBy,
			"approved_at": now,
		}).Error
}

func (r *kycRepository) Reject(ctx context.Context, id uuid.UUID, reason string) error {
	return r.db.WithContext(ctx).Model(&entities.KYC{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":           "rejected",
			"rejection_reason": reason,
		}).Error
}

func (r *kycRepository) UpdateBVNVerification(ctx context.Context, userID uuid.UUID, verified bool) error {
	now := time.Now()
	updates := map[string]interface{}{
		"bvn_verified": verified,
	}
	if verified {
		updates["bvn_verified_at"] = now
	}
	return r.db.WithContext(ctx).Model(&entities.KYC{}).
		Where("user_id = ?", userID).
		Updates(updates).Error
}

func (r *kycRepository) UpdateNINVerification(ctx context.Context, userID uuid.UUID, verified bool) error {
	now := time.Now()
	updates := map[string]interface{}{
		"nin_verified": verified,
	}
	if verified {
		updates["nin_verified_at"] = now
	}
	return r.db.WithContext(ctx).Model(&entities.KYC{}).
		Where("user_id = ?", userID).
		Updates(updates).Error
}

func (r *kycRepository) UpdateFaceVerification(ctx context.Context, userID uuid.UUID, verified bool, score float64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entities.KYC{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"face_verified":    verified,
			"face_verified_at": now,
			"liveness_score":   score,
		}).Error
}
