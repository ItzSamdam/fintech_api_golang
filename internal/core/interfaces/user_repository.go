package interfaces

import (
    "context"
    "time"
    "github.com/google/uuid"
    "fintech_api_golang/internal/core/entities"
)

type UserRepository interface {
    Create(ctx context.Context, user *entities.User) error
    Update(ctx context.Context, user *entities.User) error
    Delete(ctx context.Context, id uuid.UUID) error
    SoftDelete(ctx context.Context, id uuid.UUID) error
    GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
    GetByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.User, error)
    GetByEmail(ctx context.Context, email string) (*entities.User, error)
    GetByBVN(ctx context.Context, bvn string) (*entities.User, error)
    GetByNIN(ctx context.Context, nin string) (*entities.User, error)
    List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]entities.User, int64, error)
    UpdateTier(ctx context.Context, userID uuid.UUID, tier int) error
    UpdateLastLogin(ctx context.Context, userID uuid.UUID, ip string) error
    Suspend(ctx context.Context, userID uuid.UUID, reason string, duration *time.Duration) error
    Unsuspend(ctx context.Context, userID uuid.UUID) error
    Search(ctx context.Context, query string, offset, limit int) ([]entities.User, int64, error)
    CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)
    GetActiveUsers(ctx context.Context) (int64, error)
}

type KYCRepository interface {
    Create(ctx context.Context, kyc *entities.KYC) error
    Update(ctx context.Context, kyc *entities.KYC) error
    GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.KYC, error)
    GetByID(ctx context.Context, id uuid.UUID) (*entities.KYC, error)
    GetPending(ctx context.Context, offset, limit int) ([]entities.KYC, int64, error)
    Approve(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) error
    Reject(ctx context.Context, id uuid.UUID, reason string) error
    UpdateBVNVerification(ctx context.Context, userID uuid.UUID, verified bool) error
    UpdateNINVerification(ctx context.Context, userID uuid.UUID, verified bool) error
    UpdateFaceVerification(ctx context.Context, userID uuid.UUID, verified bool, score float64) error
}

type SessionRepository interface {
    Create(ctx context.Context, session *entities.Session) error
    Update(ctx context.Context, session *entities.Session) error
    GetByToken(ctx context.Context, token string) (*entities.Session, error)
    GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.Session, error)
    Invalidate(ctx context.Context, sessionID uuid.UUID) error
    InvalidateAllUserSessions(ctx context.Context, userID uuid.UUID) error
    CleanupExpired(ctx context.Context) error
    ExtendSession(ctx context.Context, sessionID uuid.UUID, newExpiry time.Time) error
}

type OTPRepository interface {
    Create(ctx context.Context, otp *entities.OTP) error
    GetValidOTP(ctx context.Context, phoneNumber, code, purpose string) (*entities.OTP, error)
    MarkAsUsed(ctx context.Context, id uuid.UUID) error
    IncrementAttempts(ctx context.Context, id uuid.UUID) error
    InvalidateByPhoneNumber(ctx context.Context, phoneNumber, purpose string) error
    CleanupExpired(ctx context.Context) error
}