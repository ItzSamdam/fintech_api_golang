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

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) interfaces.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *entities.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *sessionRepository) Update(ctx context.Context, session *entities.Session) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *sessionRepository) GetByToken(ctx context.Context, token string) (*entities.Session, error) {
	var session entities.Session
	err := r.db.WithContext(ctx).First(&session, "token = ?", token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.Session, error) {
	var sessions []entities.Session
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("last_active_at DESC").
		Find(&sessions).Error
	return sessions, err
}

func (r *sessionRepository) Invalidate(ctx context.Context, sessionID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.Session{}).
		Where("id = ?", sessionID).
		Update("is_active", false).Error
}

func (r *sessionRepository) InvalidateAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.Session{}).
		Where("user_id = ?", userID).
		Update("is_active", false).Error
}

func (r *sessionRepository) CleanupExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Model(&entities.Session{}).
		Where("expires_at < ? OR (is_active = ? AND last_active_at < ?)",
			time.Now(), true, time.Now().Add(-30*24*time.Hour)).
		Update("is_active", false).Error
}

func (r *sessionRepository) ExtendSession(ctx context.Context, sessionID uuid.UUID, newExpiry time.Time) error {
	return r.db.WithContext(ctx).Model(&entities.Session{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"expires_at":     newExpiry,
			"last_active_at": time.Now(),
		}).Error
}
