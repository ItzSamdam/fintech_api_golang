// package response

// import (
//     "time"
//     "github.com/google/uuid"
// )

// type AuthResponse struct {
//     AccessToken  string       `json:"access_token"`
//     RefreshToken string       `json:"refresh_token"`
//     ExpiresIn    int64        `json:"expires_in"`
//     TokenType    string       `json:"token_type"`
//     User         UserResponse `json:"user"`
// }

// type UserResponse struct {
//     ID           uuid.UUID `json:"id"`
//     PhoneNumber  string    `json:"phone_number"`
//     Email        string    `json:"email,omitempty"`
//     Tier         int       `json:"tier"`
//     IsActive     bool      `json:"is_active"`
//     IsSuspended  bool      `json:"is_suspended"`
//     CreatedAt    time.Time `json:"created_at"`
// }

// type TierLimitResponse struct {
//     Tier              int   `json:"tier"`
//     DailyLimit        int64 `json:"daily_limit"`   // In kobo
//     WeeklyLimit       int64 `json:"weekly_limit"`  // In kobo
//     MonthlyLimit      int64 `json:"monthly_limit"` // In kobo
//     SingleTxLimit     int64 `json:"single_tx_limit"` // In kobo
//     DailySpent        int64 `json:"daily_spent"`   // In kobo
//     WeeklySpent       int64 `json:"weekly_spent"`  // In kobo
//     MonthlySpent      int64 `json:"monthly_spent"` // In kobo
//     DailyRemaining    int64 `json:"daily_remaining"`
//     WeeklyRemaining   int64 `json:"weekly_remaining"`
//     MonthlyRemaining  int64 `json:"monthly_remaining"`
// }

// type KYCStatusResponse struct {
//     BVNVerified   bool       `json:"bvn_verified"`
//     NINVerified   bool       `json:"nin_verified"`
//     FaceVerified  bool       `json:"face_verified"`
//     Status        string     `json:"status"`
//     VerifiedAt    *time.Time `json:"verified_at,omitempty"`
// }
// 
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
    FirstName    string    `json:"first_name,omitempty"`
    LastName     string    `json:"last_name,omitempty"`
    Tier         int       `json:"tier"`
    IsActive     bool      `json:"is_active"`
    IsSuspended  bool      `json:"is_suspended"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
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

// Session Management Responses
type SessionResponse struct {
    ID           string    `json:"id"`
    DeviceName   string    `json:"device_name"`
    DeviceType   string    `json:"device_type"`
    IPAddress    string    `json:"ip_address"`
    Location     string    `json:"location,omitempty"`
    IsCurrent    bool      `json:"is_current"`
    LastActiveAt time.Time `json:"last_active_at"`
    CreatedAt    time.Time `json:"created_at"`
    ExpiresAt    time.Time `json:"expires_at"`
}

// 2FA Responses
type TwoFASetupResponse struct {
    Secret       string   `json:"secret"`
    QRCodeURL    string   `json:"qr_code_url"`
    BackupCodes  []string `json:"backup_codes"`
}

type TwoFAVerifyResponse struct {
    IsVerified   bool   `json:"is_verified"`
    Message      string `json:"message"`
}

// Device Management Responses
type DeviceTrustResponse struct {
    DeviceID     string `json:"device_id"`
    DeviceName   string `json:"device_name"`
    IsTrusted    bool   `json:"is_trusted"`
    Message      string `json:"message"`
}

// SIM Swap Response
type SIMSwapResponse struct {
    IsSwapped     bool      `json:"is_swapped"`
    SwappedAt     *time.Time `json:"swapped_at,omitempty"`
    PreviousSIM   string     `json:"previous_sim,omitempty"`
    CurrentSIM    string     `json:"current_sim,omitempty"`
    Message       string     `json:"message"`
}

// Limit Check Response
type LimitCheckResponse struct {
    IsAllowed     bool   `json:"is_allowed"`
    CurrentAmount int64  `json:"current_amount"`
    LimitAmount   int64  `json:"limit_amount"`
    LimitType     string `json:"limit_type"` // daily, weekly, monthly, single
    Remaining     int64  `json:"remaining"`
    ResetsAt      string `json:"resets_at,omitempty"`
}

// Suspicious Report Response
type SuspiciousReportResponse struct {
    ReportID      string `json:"report_id"`
    Status        string `json:"status"`
    Message       string `json:"message"`
    Reference     string `json:"reference"`
}

// Password Reset Response
type PasswordResetResponse struct {
    Message      string `json:"message"`
    Reference    string `json:"reference"`
    ExpiresIn    int    `json:"expires_in"` // seconds
}

// Logout Response
type LogoutResponse struct {
    Message      string `json:"message"`
    Success      bool   `json:"success"`
}