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

type otpRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) interfaces.OTPRepository {
	return &otpRepository{db: db}
}

func (r *otpRepository) Create(ctx context.Context, otp *entities.OTP) error {
	return r.db.WithContext(ctx).Create(otp).Error
}

func (r *otpRepository) GetValidOTP(ctx context.Context, phoneNumber, code, purpose string) (*entities.OTP, error) {
	var otp entities.OTP
	err := r.db.WithContext(ctx).
		Where("phone_number = ? AND code = ? AND purpose = ? AND is_used = ? AND expires_at > ?",
			phoneNumber, code, purpose, false, time.Now()).
		First(&otp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &otp, nil
}

func (r *otpRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.OTP{}).
		Where("id = ?", id).
		Update("is_used", true).Error
}

func (r *otpRepository) IncrementAttempts(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.OTP{}).
		Where("id = ?", id).
		Update("attempts", gorm.Expr("attempts + ?", 1)).Error
}

func (r *otpRepository) InvalidateByPhoneNumber(ctx context.Context, phoneNumber, purpose string) error {
	return r.db.WithContext(ctx).Model(&entities.OTP{}).
		Where("phone_number = ? AND purpose = ? AND is_used = ?", phoneNumber, purpose, false).
		Update("is_used", true).Error
}

func (r *otpRepository) CleanupExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&entities.OTP{}).Error
}
