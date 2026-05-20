package response

import (
    "time"
    "github.com/google/uuid"
)

type AuthResponse struct {
    AccessToken  string       `json:"access_token"`
    RefreshToken string       `json:"refresh_token"`
    ExpiresIn    int64        `json:"expires_in"`
    TokenType    string       `json:"token_type"`
    User         UserResponse `json:"user"`
}

type UserResponse struct {
    ID           uuid.UUID `json:"id"`
    PhoneNumber  string    `json:"phone_number"`
    Email        string    `json:"email,omitempty"`
    Tier         int       `json:"tier"`
    IsActive     bool      `json:"is_active"`
    IsSuspended  bool      `json:"is_suspended"`
    CreatedAt    time.Time `json:"created_at"`
}

type TierLimitResponse struct {
    Tier              int   `json:"tier"`
    DailyLimit        int64 `json:"daily_limit"`   // In kobo
    WeeklyLimit       int64 `json:"weekly_limit"`  // In kobo
    MonthlyLimit      int64 `json:"monthly_limit"` // In kobo
    SingleTxLimit     int64 `json:"single_tx_limit"` // In kobo
    DailySpent        int64 `json:"daily_spent"`   // In kobo
    WeeklySpent       int64 `json:"weekly_spent"`  // In kobo
    MonthlySpent      int64 `json:"monthly_spent"` // In kobo
    DailyRemaining    int64 `json:"daily_remaining"`
    WeeklyRemaining   int64 `json:"weekly_remaining"`
    MonthlyRemaining  int64 `json:"monthly_remaining"`
}

type KYCStatusResponse struct {
    BVNVerified   bool       `json:"bvn_verified"`
    NINVerified   bool       `json:"nin_verified"`
    FaceVerified  bool       `json:"face_verified"`
    Status        string     `json:"status"`
    VerifiedAt    *time.Time `json:"verified_at,omitempty"`
}