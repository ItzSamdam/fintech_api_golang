package request

type UpgradeTierRequest struct {
    Tier      int    `json:"tier" validate:"required,min=0,max=3"`
    Reason    string `json:"reason"`
    ApprovedBy string `json:"approved_by" validate:"required"`
}

type SuspendUserRequest struct {
    Reason      string `json:"reason" validate:"required"`
    Duration    string `json:"duration"` // e.g., "7d", "30d", "permanent"
}

type OverrideLimitsRequest struct {
    DailyLimit    *int64 `json:"daily_limit"`
    WeeklyLimit   *int64 `json:"weekly_limit"`
    MonthlyLimit  *int64 `json:"monthly_limit"`
    SingleTxLimit *int64 `json:"single_tx_limit"`
    Reason        string `json:"reason" validate:"required"`
}

type ReverseTransactionRequest struct {
    TransactionID string `json:"transaction_id" validate:"required"`
    Reason        string `json:"reason" validate:"required"`
    NotifyUser    bool   `json:"notify_user"`
}

type ManualCreditRequest struct {
    UserID      string `json:"user_id" validate:"required"`
    Amount      int64  `json:"amount" validate:"required,min=100"`
    Description string `json:"description" validate:"required"`
    Reference   string `json:"reference" validate:"required"`
}

type ManualDebitRequest struct {
    UserID      string `json:"user_id" validate:"required"`
    Amount      int64  `json:"amount" validate:"required,min=100"`
    Description string `json:"description" validate:"required"`
    Reference   string `json:"reference" validate:"required"`
}

type ApproveKYCRequest struct {
    KYCID        string `json:"kyc_id" validate:"required"`
    Notes        string `json:"notes"`
}

type RejectKYCRequest struct {
    KYCID        string `json:"kyc_id" validate:"required"`
    Reason       string `json:"reason" validate:"required"`
}

type ToggleProviderRequest struct {
    ProviderID string `json:"provider_id" validate:"required"`
    IsActive   bool   `json:"is_active"`
}

type SetProviderPriorityRequest struct {
    ProviderID string `json:"provider_id" validate:"required"`
    Priority   int    `json:"priority" validate:"required,min=1"`
}

type UpdateFeeRequest struct {
    BillType   string  `json:"bill_type" validate:"required"`
    FeeType    string  `json:"fee_type" validate:"required,oneof=percentage fixed"`
    FeeValue   float64 `json:"fee_value" validate:"required,min=0"`
    CapAmount  int64   `json:"cap_amount"`
    MinAmount  int64   `json:"min_amount"`
    MaxAmount  int64   `json:"max_amount"`
}

type UpdateMarginRequest struct {
    ProviderID    string  `json:"provider_id" validate:"required"`
    MarginPercent float64 `json:"margin_percent" validate:"required,min=0,max=100"`
}

type SystemSettingsRequest struct {
    MaintenanceMode bool   `json:"maintenance_mode"`
    MaintenanceMessage string `json:"maintenance_message"`
    GlobalDailyLimit   int64  `json:"global_daily_limit"`
    GlobalSingleTxLimit int64 `json:"global_single_tx_limit"`
}

type CreateRoleRequest struct {
    Name        string   `json:"name" validate:"required"`
    Permissions []string `json:"permissions" validate:"required"`
    Description string   `json:"description"`
}

type InviteStaffRequest struct {
    Email    string `json:"email" validate:"required,email"`
    FullName string `json:"full_name" validate:"required"`
    Role     string `json:"role" validate:"required"`
}