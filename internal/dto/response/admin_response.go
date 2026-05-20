package response

import (
    "time"
    "github.com/google/uuid"
)

type UserListResponse struct {
    Users      []UserDetailResponse `json:"users"`
    Total      int64                `json:"total"`
    Page       int                  `json:"page"`
    Limit      int                  `json:"limit"`
    TotalPages int                  `json:"total_pages"`
}

type UserDetailResponse struct {
    ID           uuid.UUID `json:"id"`
    PhoneNumber  string    `json:"phone_number"`
    Email        string    `json:"email"`
    Tier         int       `json:"tier"`
    IsActive     bool      `json:"is_active"`
    IsSuspended  bool      `json:"is_suspended"`
    SuspendedAt  *time.Time `json:"suspended_at,omitempty"`
    KYCStatus    string    `json:"kyc_status"`
    Wallet       *WalletResponse `json:"wallet,omitempty"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type TransactionSummaryResponse struct {
    TotalTransactions  int64   `json:"total_transactions"`
    TotalVolume        int64   `json:"total_volume"`     // In kobo
    TotalVolumeNaira   float64 `json:"total_volume_naira"`
    SuccessfulCount    int64   `json:"successful_count"`
    FailedCount        int64   `json:"failed_count"`
    PendingCount       int64   `json:"pending_count"`
    AverageTransaction int64   `json:"average_transaction"`
}

type DashboardStatsResponse struct {
    TotalUsers        int64   `json:"total_users"`
    ActiveUsers       int64   `json:"active_users"`
    NewUsersToday     int64   `json:"new_users_today"`
    TotalWallets      int64   `json:"total_wallets"`
    TotalBalance      int64   `json:"total_balance"`      // In kobo
    TotalBalanceNaira float64 `json:"total_balance_naira"`
    TodayVolume       int64   `json:"today_volume"`       // In kobo
    TodayRevenue      int64   `json:"today_revenue"`      // In kobo
}

type RevenueReportResponse struct {
    Period          string                     `json:"period"`
    TotalRevenue    int64                      `json:"total_revenue"`
    RevenueByBillType map[string]int64         `json:"revenue_by_bill_type"`
    FeeBreakdown    map[string]int64           `json:"fee_breakdown"`
    Chart           []RevenuePointResponse     `json:"chart"`
}

type RevenuePointResponse struct {
    Date    string `json:"date"`
    Amount  int64  `json:"amount"`
}

type ProviderPerformanceResponse struct {
    ProviderID      string  `json:"provider_id"`
    ProviderName    string  `json:"provider_name"`
    TotalRequests   int64   `json:"total_requests"`
    SuccessCount    int64   `json:"success_count"`
    FailedCount     int64   `json:"failed_count"`
    SuccessRate     float64 `json:"success_rate"`
    AvgResponseTime int     `json:"avg_response_time"` // milliseconds
}

type AuditLogResponse struct {
    ID          uuid.UUID `json:"id"`
    AdminEmail  string    `json:"admin_email"`
    Action      string    `json:"action"`
    EntityType  string    `json:"entity_type"`
    EntityID    string    `json:"entity_id"`
    IPAddress   string    `json:"ip_address"`
    CreatedAt   time.Time `json:"created_at"`
}