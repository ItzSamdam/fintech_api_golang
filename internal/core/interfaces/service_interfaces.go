package interfaces

import (
    "context"
    "time"
    "github.com/google/uuid"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
)

// Auth Service Interface
type AuthService interface {
    RegisterPhone(ctx context.Context, req *request.RegisterPhoneRequest) (*response.OTPResponse, error)
    VerifyOTP(ctx context.Context, req *request.VerifyOTPRequest) (*response.AuthResponse, error)
    RegisterBVN(ctx context.Context, userID uuid.UUID, req *request.RegisterBVNRequest) error
    VerifyFace(ctx context.Context, userID uuid.UUID, req *request.VerifyFaceRequest) error
    Login(ctx context.Context, req *request.LoginRequest) (*response.AuthResponse, error)
    ChangePassword(ctx context.Context, userID uuid.UUID, req *request.ChangePasswordRequest) error
    ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) error
    Logout(ctx context.Context, userID uuid.UUID, token string, allDevices bool) error
    RefreshToken(ctx context.Context, refreshToken string) (*response.AuthResponse, error)
    GetUserProfile(ctx context.Context, userID uuid.UUID) (*response.UserResponse, error)
    UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *request.UpdateUserRequest) error
}

// Wallet Service Interface
type WalletService interface {
    CreateWallet(ctx context.Context, userID uuid.UUID, currency string) (*response.WalletResponse, error)
    GetBalance(ctx context.Context, userID uuid.UUID) (*response.BalanceResponse, error)
    GetTransactions(ctx context.Context, userID uuid.UUID, req *request.GetTransactionsRequest) (*response.TransactionHistoryResponse, error)
    GetLimits(ctx context.Context, userID uuid.UUID) (*response.TierLimitResponse, error)
    LockWallet(ctx context.Context, userID uuid.UUID, reason string) error
    UnlockWallet(ctx context.Context, userID uuid.UUID) error
    GetStatement(ctx context.Context, userID uuid.UUID, req *request.GetStatementRequest) ([]byte, string, error)
}

// Transfer Service Interface
type TransferService interface {
    SendTransfer(ctx context.Context, userID uuid.UUID, req *request.SendTransferRequest) (*response.TransactionResponse, error)
    GetTransferStatus(ctx context.Context, userID uuid.UUID, reference string) (*response.TransactionResponse, error)
    RetryTransfer(ctx context.Context, userID uuid.UUID, reference string) (*response.TransactionResponse, error)
    NameEnquiry(ctx context.Context, req *request.NameEnquiryRequest) (*response.NameEnquiryResponse, error)
    GetBanks(ctx context.Context) (*response.BankListResponse, error)
    GetTransferHistory(ctx context.Context, userID uuid.UUID, offset, limit int) (*response.TransactionHistoryResponse, error)
}

// Bill Payment Service Interface
type BillPaymentService interface {
    // Airtime
    GetAirtimeNetworks(ctx context.Context) (*response.NetworkListResponse, error)
    PurchaseAirtime(ctx context.Context, userID uuid.UUID, req *request.PurchaseAirtimeRequest) (*response.TransactionResponse, error)
    
    // Data
    GetDataNetworks(ctx context.Context) (*response.NetworkListResponse, error)
    GetDataPlans(ctx context.Context, network string) ([]response.DataPlanResponse, error)
    PurchaseData(ctx context.Context, userID uuid.UUID, req *request.PurchaseDataRequest) (*response.TransactionResponse, error)
    
    // Electricity
    GetElectricityProviders(ctx context.Context) ([]response.ProviderResponse, error)
    ValidateMeter(ctx context.Context, req *request.ValidateMeterRequest) (*response.MeterValidationResponse, error)
    PayElectricity(ctx context.Context, userID uuid.UUID, req *request.PayElectricityRequest) (*response.ElectricityPaymentResponse, error)
    GetElectricityToken(ctx context.Context, userID uuid.UUID, transactionID string) (string, error)
    
    // Betting
    GetBettingProviders(ctx context.Context) ([]response.ProviderResponse, error)
    ValidateBettingAccount(ctx context.Context, req *request.ValidateBettingAccountRequest) (*response.BettingAccountResponse, error)
    FundBettingWallet(ctx context.Context, userID uuid.UUID, req *request.FundBettingRequest) (*response.TransactionResponse, error)
    GetBettingHistory(ctx context.Context, userID uuid.UUID, offset, limit int) (*response.TransactionHistoryResponse, error)
}

// Savings Service Interface
type SavingsService interface {
    CreateGoal(ctx context.Context, userID uuid.UUID, req *request.CreateSavingsGoalRequest) (*response.SavingsGoalResponse, error)
    ContributeToGoal(ctx context.Context, userID uuid.UUID, req *request.ContributeToGoalRequest) (*response.SavingsContributionResponse, error)
    GetGoals(ctx context.Context, userID uuid.UUID) ([]response.SavingsGoalResponse, error)
    GetGoal(ctx context.Context, userID uuid.UUID, goalID string) (*response.SavingsGoalResponse, error)
    UpdateGoal(ctx context.Context, userID uuid.UUID, goalID string, req *request.UpdateSavingsGoalRequest) error
    DeleteGoal(ctx context.Context, userID uuid.UUID, goalID string) error
    ActivateRoundup(ctx context.Context, userID uuid.UUID, req *request.ActivateRoundupRequest) error
    DeactivateRoundup(ctx context.Context, userID uuid.UUID) error
    GetRoundupStatus(ctx context.Context, userID uuid.UUID) (*response.RoundupStatusResponse, error)
}

// Compliance Service Interface
type ComplianceService interface {
    ReportSuspicious(ctx context.Context, userID uuid.UUID, req *request.ReportSuspiciousRequest) error
    CheckLimits(ctx context.Context, userID uuid.UUID, req *request.CheckLimitsRequest) (bool, error)
    VerifySIMSwap(ctx context.Context, phoneNumber string) (bool, error)
    TrustDevice(ctx context.Context, userID uuid.UUID, req *request.TrustDeviceRequest) error
    Enable2FA(ctx context.Context, userID uuid.UUID, password string) (string, []string, error)
    Verify2FA(ctx context.Context, userID uuid.UUID, code string) (bool, error)
    Disable2FA(ctx context.Context, userID uuid.UUID, password string) error
    GetActiveSessions(ctx context.Context, userID uuid.UUID) ([]response.SessionResponse, error)
    TerminateSession(ctx context.Context, userID uuid.UUID, sessionID string) error
}

// Admin Service Interface
type AdminService interface {
    // User Management
    ListUsers(ctx context.Context, offset, limit int, filters map[string]interface{}) (*response.UserListResponse, error)
    GetUserDetails(ctx context.Context, userID uuid.UUID) (*response.UserDetailResponse, error)
    UpgradeUserTier(ctx context.Context, userID uuid.UUID, tier int, approvedBy uuid.UUID) error
    SuspendUser(ctx context.Context, userID uuid.UUID, reason string, duration *string) error
    UnsuspendUser(ctx context.Context, userID uuid.UUID) error
    DeleteUser(ctx context.Context, userID uuid.UUID) error
    OverrideLimits(ctx context.Context, userID uuid.UUID, req *request.OverrideLimitsRequest) error
    
    // Transaction Management
    ListTransactions(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]response.TransactionResponse, int64, error)
    ReverseTransaction(ctx context.Context, transactionID string, reason string, notifyUser bool) error
    GetTransactionSummary(ctx context.Context, startDate, endDate time.Time) (*response.TransactionSummaryResponse, error)
    
    // Wallet Management
    ManualCredit(ctx context.Context, req *request.ManualCreditRequest, adminID uuid.UUID) (*response.TransactionResponse, error)
    ManualDebit(ctx context.Context, req *request.ManualDebitRequest, adminID uuid.UUID) (*response.TransactionResponse, error)
    FreezeWallet(ctx context.Context, userID uuid.UUID, reason string) error
    UnfreezeWallet(ctx context.Context, userID uuid.UUID) error
    GetBalanceSummary(ctx context.Context) (*response.DashboardStatsResponse, error)
    
    // KYC Management
    GetPendingKYC(ctx context.Context, offset, limit int) ([]response.KYCStatusResponse, int64, error)
    ApproveKYC(ctx context.Context, kycID string, notes string, adminID uuid.UUID) error
    RejectKYC(ctx context.Context, kycID string, reason string, adminID uuid.UUID) error
    
    // Provider Management
    ListProviders(ctx context.Context, providerType string) ([]response.ProviderResponse, error)
    ToggleProvider(ctx context.Context, providerID string, isActive bool) error
    SetProviderPriority(ctx context.Context, providerID string, priority int) error
    CheckProviderHealth(ctx context.Context, providerID string) (string, error)
    
    // Fee Management
    UpdateFee(ctx context.Context, req *request.UpdateFeeRequest) error
    UpdateMargin(ctx context.Context, req *request.UpdateMarginRequest) error
    
    // Reports
    GetDailyReport(ctx context.Context, date time.Time) (*response.RevenueReportResponse, error)
    GetMonthlyReport(ctx context.Context, year, month int) (*response.RevenueReportResponse, error)
    GetRevenueByBillType(ctx context.Context, startDate, endDate time.Time) (map[string]int64, error)
    GetTopUsers(ctx context.Context, limit int) ([]response.UserDetailResponse, error)
    ExportReport(ctx context.Context, reportType string, format string, startDate, endDate time.Time) ([]byte, string, error)
    
    // System Settings
    GetSystemSettings(ctx context.Context) (*SystemSettings, error)
    UpdateSystemSettings(ctx context.Context, req *request.SystemSettingsRequest) error
    GetAuditLogs(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]response.AuditLogResponse, int64, error)
    GetSystemMetrics(ctx context.Context) (*SystemMetrics, error)
}

type SystemSettings struct {
    MaintenanceMode      bool
    MaintenanceMessage   string
    GlobalDailyLimit     int64
    GlobalSingleTxLimit  int64
    MaxRetryCount        int
    SessionTimeout       int
}

type SystemMetrics struct {
    CPUUsage    float64
    MemoryUsage float64
    Goroutines  int
    Uptime      string
    DBConnections int
    RedisConnections int
}

// Notification Service Interface
type NotificationService interface {
    SendSMS(ctx context.Context, phoneNumber, message string) error
    SendEmail(ctx context.Context, to, subject, body string) error
    SendPushNotification(ctx context.Context, userID uuid.UUID, title, body string, data map[string]interface{}) error
    GetInAppNotifications(ctx context.Context, userID uuid.UUID, offset, limit int, unreadOnly bool) ([]response.InAppNotificationResponse, int64, error)
    MarkNotificationAsRead(ctx context.Context, userID uuid.UUID, notificationID string) error
}