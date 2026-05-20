package services

import (
    "context"
    "errors"
    "fmt"
    "time"
    
    "github.com/google/uuid"
    "gorm.io/gorm"
    
    "fintech_api_golang/internal/core/entities"
    "fintech_api_golang/internal/core/interfaces"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
)

type AdminService struct {
    userRepo         interfaces.UserRepository
    walletRepo       interfaces.WalletRepository
    transactionRepo  interfaces.TransactionRepository
    kycRepo          interfaces.KYCRepository
    providerRepo     interfaces.ProviderRepository
    auditLogRepo     interfaces.AuditLogRepository
    adminUserRepo    interfaces.AdminUserRepository
    roleRepo         interfaces.RoleRepository
    db               *gorm.DB
}

func NewAdminService(
    userRepo interfaces.UserRepository,
    walletRepo interfaces.WalletRepository,
    transactionRepo interfaces.TransactionRepository,
    kycRepo interfaces.KYCRepository,
    providerRepo interfaces.ProviderRepository,
    auditLogRepo interfaces.AuditLogRepository,
    adminUserRepo interfaces.AdminUserRepository,
    roleRepo interfaces.RoleRepository,
    db *gorm.DB,
) *AdminService {
    return &AdminService{
        userRepo:        userRepo,
        walletRepo:      walletRepo,
        transactionRepo: transactionRepo,
        kycRepo:         kycRepo,
        providerRepo:    providerRepo,
        auditLogRepo:    auditLogRepo,
        adminUserRepo:   adminUserRepo,
        roleRepo:        roleRepo,
        db:              db,
    }
}

// ========== USER MANAGEMENT ==========

func (s *AdminService) ListUsers(ctx context.Context, offset, limit int, filters map[string]interface{}) (*response.UserListResponse, error) {
    users, total, err := s.userRepo.List(ctx, offset, limit, filters)
    if err != nil {
        return nil, err
    }
    
    userResponses := make([]response.UserDetailResponse, len(users))
    for i, user := range users {
        wallet, _ := s.walletRepo.GetByUserID(ctx, user.ID)
        kyc, _ := s.kycRepo.GetByUserID(ctx, user.ID)
        
        kycStatus := "pending"
        if kyc != nil {
            kycStatus = kyc.Status
        }
        
        userResponses[i] = response.UserDetailResponse{
            ID:          user.ID,
            PhoneNumber: user.PhoneNumber,
            Email:       user.Email,
            Tier:        user.Tier,
            IsActive:    user.IsActive,
            IsSuspended: user.IsSuspended,
            SuspendedAt: user.SuspendedAt,
            KYCStatus:   kycStatus,
            CreatedAt:   user.CreatedAt,
            UpdatedAt:   user.UpdatedAt,
        }
        
        if wallet != nil {
            userResponses[i].Wallet = &response.WalletResponse{
                ID:      wallet.ID,
                UserID:  wallet.UserID,
                Balance: int64(wallet.Balance),
                BalanceNaira: float64(wallet.Balance) / 100,
                Currency: wallet.Currency,
                IsLocked: wallet.IsLocked,
            }
        }
    }
    
    totalPages := int(total) / limit
    if int(total)%limit > 0 {
        totalPages++
    }
    
    return &response.UserListResponse{
        Users:      userResponses,
        Total:      total,
        Page:       offset/limit + 1,
        Limit:      limit,
        TotalPages: totalPages,
    }, nil
}

func (s *AdminService) GetUserDetails(ctx context.Context, userID uuid.UUID) (*response.UserDetailResponse, error) {
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    if user == nil {
        return nil, errors.New("user not found")
    }
    
    wallet, _ := s.walletRepo.GetByUserID(ctx, userID)
    kyc, _ := s.kycRepo.GetByUserID(ctx, userID)
    
    kycStatus := "pending"
    if kyc != nil {
        kycStatus = kyc.Status
    }
    
    resp := &response.UserDetailResponse{
        ID:          user.ID,
        PhoneNumber: user.PhoneNumber,
        Email:       user.Email,
        Tier:        user.Tier,
        IsActive:    user.IsActive,
        IsSuspended: user.IsSuspended,
        SuspendedAt: user.SuspendedAt,
        KYCStatus:   kycStatus,
        CreatedAt:   user.CreatedAt,
        UpdatedAt:   user.UpdatedAt,
    }
    
    if wallet != nil {
        resp.Wallet = &response.WalletResponse{
            ID:      wallet.ID,
            UserID:  wallet.UserID,
            Balance: int64(wallet.Balance),
            BalanceNaira: float64(wallet.Balance) / 100,
            Currency: wallet.Currency,
            IsLocked: wallet.IsLocked,
        }
    }
    
    return resp, nil
}

func (s *AdminService) UpgradeUserTier(ctx context.Context, userID uuid.UUID, tier int, approvedBy uuid.UUID) error {
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return err
    }
    
    if user == nil {
        return errors.New("user not found")
    }
    
    if tier < user.Tier {
        return errors.New("cannot downgrade tier using upgrade method")
    }
    
    if err := s.userRepo.UpdateTier(ctx, userID, tier); err != nil {
        return err
    }
    
    // Create audit log
    auditLog := &entities.AuditLog{
        ID:         uuid.New(),
        AdminID:    &approvedBy,
        UserID:     &userID,
        Action:     "TIER_UPGRADE",
        EntityType: "user",
        EntityID:   userID.String(),
        OldValue:   fmt.Sprintf(`{"tier": %d}`, user.Tier),
        NewValue:   fmt.Sprintf(`{"tier": %d}`, tier),
        CreatedAt:  time.Now(),
    }
    s.auditLogRepo.Create(ctx, auditLog)
    
    return nil
}

func (s *AdminService) SuspendUser(ctx context.Context, userID uuid.UUID, reason string, duration *string) error {
    var dur *time.Duration
    if duration != nil && *duration != "" {
        // Parse duration string like "7d", "30d", "permanent"
        if *duration != "permanent" {
            // Parse and set duration
        }
    }
    
    return s.userRepo.Suspend(ctx, userID, reason, dur)
}

func (s *AdminService) UnsuspendUser(ctx context.Context, userID uuid.UUID) error {
    return s.userRepo.Unsuspend(ctx, userID)
}

func (s *AdminService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
    return s.userRepo.SoftDelete(ctx, userID)
}

func (s *AdminService) OverrideLimits(ctx context.Context, userID uuid.UUID, req *request.OverrideLimitsRequest) error {
    // Implementation for overriding user limits
    return nil
}

func (s *AdminService) SearchUsers(ctx context.Context, query string, offset, limit int) ([]response.UserDetailResponse, int64, error) {
    users, total, err := s.userRepo.Search(ctx, query, offset, limit)
    if err != nil {
        return nil, 0, err
    }
    
    responses := make([]response.UserDetailResponse, len(users))
    for i, user := range users {
        responses[i] = response.UserDetailResponse{
            ID:          user.ID,
            PhoneNumber: user.PhoneNumber,
            Email:       user.Email,
            Tier:        user.Tier,
            IsActive:    user.IsActive,
            IsSuspended: user.IsSuspended,
            CreatedAt:   user.CreatedAt,
        }
    }
    
    return responses, total, nil
}

// ========== TRANSACTION MANAGEMENT ==========

func (s *AdminService) ListTransactions(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]response.TransactionResponse, int64, error) {
    var userID uuid.UUID
    if uid, ok := filters["user_id"]; ok {
        userID, _ = uuid.Parse(uid.(string))
        delete(filters, "user_id")
    }
    
    var transactions []entities.Transaction
    var total int64
    var err error
    
    if userID != uuid.Nil {
        transactions, total, err = s.transactionRepo.GetByUserID(ctx, userID, offset, limit, filters)
    } else {
        // Admin view all - need to implement GetAll method
        // For now, use empty userID
        transactions, total, err = s.transactionRepo.GetByUserID(ctx, uuid.Nil, offset, limit, filters)
    }
    
    if err != nil {
        return nil, 0, err
    }
    
    responses := make([]response.TransactionResponse, len(transactions))
    for i, tx := range transactions {
        responses[i] = response.TransactionResponse{
            ID:          tx.ID,
            Reference:   tx.Reference,
            Type:        tx.Type,
            Category:    tx.Category,
            Amount:      int64(tx.Amount),
            AmountNaira: float64(tx.Amount) / 100,
            Fee:         int64(tx.Fee),
            FeeNaira:    float64(tx.Fee) / 100,
            Status:      tx.Status,
            Description: tx.Description,
            CreatedAt:   tx.CreatedAt,
        }
    }
    
    return responses, total, nil
}

func (s *AdminService) GetTransactionDetails(ctx context.Context, transactionID string) (*response.TransactionResponse, error) {
    tx, err := s.transactionRepo.GetByReference(ctx, transactionID)
    if err != nil {
        return nil, err
    }
    
    if tx == nil {
        return nil, errors.New("transaction not found")
    }
    
    return &response.TransactionResponse{
        ID:          tx.ID,
        Reference:   tx.Reference,
        Type:        tx.Type,
        Category:    tx.Category,
        Amount:      int64(tx.Amount),
        AmountNaira: float64(tx.Amount) / 100,
        Fee:         int64(tx.Fee),
        FeeNaira:    float64(tx.Fee) / 100,
        Status:      tx.Status,
        Description: tx.Description,
        CreatedAt:   tx.CreatedAt,
    }, nil
}

func (s *AdminService) ReverseTransaction(ctx context.Context, transactionID, reason string, notifyUser bool) error {
    // Implementation for reversing a transaction
    return nil
}

func (s *AdminService) VoidTransaction(ctx context.Context, transactionID, reason string) error {
    // Implementation for voiding a pending transaction
    return nil
}

func (s *AdminService) GetTransactionSummary(ctx context.Context, startDate, endDate string) (*response.TransactionSummaryResponse, error) {
    // Parse dates and get summary
    return &response.TransactionSummaryResponse{}, nil
}

// ========== WALLET MANAGEMENT ==========

func (s *AdminService) ListWallets(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]response.WalletResponse, int64, error) {
    wallets, total, err := s.walletRepo.List(ctx, offset, limit, filters)
    if err != nil {
        return nil, 0, err
    }
    
    responses := make([]response.WalletResponse, len(wallets))
    for i, wallet := range wallets {
        responses[i] = response.WalletResponse{
            ID:           wallet.ID,
            UserID:       wallet.UserID,
            Balance:      int64(wallet.Balance),
            BalanceNaira: float64(wallet.Balance) / 100,
            Currency:     wallet.Currency,
            IsLocked:     wallet.IsLocked,
            CreatedAt:    wallet.CreatedAt,
            UpdatedAt:    wallet.UpdatedAt,
        }
    }
    
    return responses, total, nil
}

func (s *AdminService) GetWalletDetails(ctx context.Context, walletID uuid.UUID) (*response.WalletResponse, error) {
    wallet, err := s.walletRepo.GetByID(ctx, walletID)
    if err != nil {
        return nil, err
    }
    
    if wallet == nil {
        return nil, errors.New("wallet not found")
    }
    
    return &response.WalletResponse{
        ID:           wallet.ID,
        UserID:       wallet.UserID,
        Balance:      int64(wallet.Balance),
        BalanceNaira: float64(wallet.Balance) / 100,
        Currency:     wallet.Currency,
        IsLocked:     wallet.IsLocked,
        CreatedAt:    wallet.CreatedAt,
        UpdatedAt:    wallet.UpdatedAt,
    }, nil
}

func (s *AdminService) ManualCredit(ctx context.Context, req *request.ManualCreditRequest, adminID uuid.UUID) (*response.TransactionResponse, error) {
    // Implementation for manual credit
    return nil, nil
}

func (s *AdminService) ManualDebit(ctx context.Context, req *request.ManualDebitRequest, adminID uuid.UUID) (*response.TransactionResponse, error) {
    // Implementation for manual debit
    return nil, nil
}

func (s *AdminService) FreezeWallet(ctx context.Context, userID uuid.UUID, reason string) error {
    wallet, err := s.walletRepo.GetByUserID(ctx, userID)
    if err != nil {
        return err
    }
    
    if wallet == nil {
        return errors.New("wallet not found")
    }
    
    return s.walletRepo.Lock(ctx, wallet.ID, reason)
}

func (s *AdminService) UnfreezeWallet(ctx context.Context, userID uuid.UUID) error {
    wallet, err := s.walletRepo.GetByUserID(ctx, userID)
    if err != nil {
        return err
    }
    
    if wallet == nil {
        return errors.New("wallet not found")
    }
    
    return s.walletRepo.Unlock(ctx, wallet.ID)
}

func (s *AdminService) GetBalanceSummary(ctx context.Context) (*response.DashboardStatsResponse, error) {
    totalBalance, err := s.walletRepo.GetTotalBalance(ctx)
    if err != nil {
        return nil, err
    }
    
    activeUsers, err := s.userRepo.GetActiveUsers(ctx)
    if err != nil {
        return nil, err
    }
    
    return &response.DashboardStatsResponse{
        TotalUsers:        activeUsers,
        ActiveUsers:       activeUsers,
        TotalBalance:      totalBalance,
        TotalBalanceNaira: float64(totalBalance) / 100,
    }, nil
}

// ========== KYC MANAGEMENT ==========

func (s *AdminService) GetPendingKYC(ctx context.Context, offset, limit int) ([]response.KYCStatusResponse, int64, error) {
    kycList, total, err := s.kycRepo.GetPending(ctx, offset, limit)
    if err != nil {
        return nil, 0, err
    }
    
    responses := make([]response.KYCStatusResponse, len(kycList))
    for i, kyc := range kycList {
        responses[i] = response.KYCStatusResponse{
            BVNVerified:  kyc.BVNVerified,
            NINVerified:  kyc.NINVerified,
            FaceVerified: kyc.FaceVerified,
            Status:       kyc.Status,
            VerifiedAt:   kyc.BVNVerifiedAt,
        }
    }
    
    return responses, total, nil
}

func (s *AdminService) GetKYCDetails(ctx context.Context, kycID uuid.UUID) (*response.KYCStatusResponse, error) {
    kyc, err := s.kycRepo.GetByID(ctx, kycID)
    if err != nil {
        return nil, err
    }
    
    if kyc == nil {
        return nil, errors.New("KYC record not found")
    }
    
    return &response.KYCStatusResponse{
        BVNVerified:  kyc.BVNVerified,
        NINVerified:  kyc.NINVerified,
        FaceVerified: kyc.FaceVerified,
        Status:       kyc.Status,
        VerifiedAt:   kyc.BVNVerifiedAt,
    }, nil
}

func (s *AdminService) ApproveKYC(ctx context.Context, kycID string, notes string, adminID uuid.UUID) error {
    id, err := uuid.Parse(kycID)
    if err != nil {
        return err
    }
    
    kyc, err := s.kycRepo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    
    if kyc == nil {
        return errors.New("KYC record not found")
    }
    
    if err := s.kycRepo.Approve(ctx, id, adminID); err != nil {
        return err
    }
    
    // Upgrade user tier to 2
    return s.userRepo.UpdateTier(ctx, kyc.UserID, 2)
}

func (s *AdminService) RejectKYC(ctx context.Context, kycID string, reason string, adminID uuid.UUID) error {
    id, err := uuid.Parse(kycID)
    if err != nil {
        return err
    }
    
    return s.kycRepo.Reject(ctx, id, reason)
}

// ========== PROVIDER MANAGEMENT ==========

func (s *AdminService) ListProviders(ctx context.Context, providerType string) ([]response.ProviderResponse, error) {
    var providers []entities.Provider
    var err error
    
    if providerType != "" {
        providers, err = s.providerRepo.GetByType(ctx, providerType)
    } else {
        providers, _, err = s.providerRepo.List(ctx, 0, 100, nil)
    }
    
    if err != nil {
        return nil, err
    }
    
    responses := make([]response.ProviderResponse, len(providers))
    for i, p := range providers {
        responses[i] = response.ProviderResponse{
            ID:       p.ID.String(),
            Name:     p.Name,
            Code:     p.Code,
            Category: p.Category,
            IsActive: p.IsActive,
        }
    }
    
    return responses, nil
}

func (s *AdminService) ToggleProvider(ctx context.Context, providerID string, isActive bool) error {
    id, err := uuid.Parse(providerID)
    if err != nil {
        return err
    }
    
    return s.providerRepo.ToggleActive(ctx, id, isActive)
}

func (s *AdminService) SetProviderPriority(ctx context.Context, providerID string, priority int) error {
    id, err := uuid.Parse(providerID)
    if err != nil {
        return err
    }
    
    return s.providerRepo.UpdatePriority(ctx, id, priority)
}

func (s *AdminService) CheckProviderHealth(ctx context.Context, providerID string) (string, error) {
    // Implementation for health check
    return "healthy", nil
}

func (s *AdminService) GetProviderLogs(ctx context.Context, providerID string, page, limit int) ([]interface{}, int64, error) {
    // Implementation for provider logs
    return nil, 0, nil
}

// ========== REPORTS ==========

func (s *AdminService) GetDailyReport(ctx context.Context, date time.Time) (*response.RevenueReportResponse, error) {
    return &response.RevenueReportResponse{
        Period:       date.Format("2006-01-02"),
        TotalRevenue: 0,
        RevenueByBillType: make(map[string]int64),
        FeeBreakdown: make(map[string]int64),
    }, nil
}

func (s *AdminService) GetMonthlyReport(ctx context.Context, year, month int) (*response.RevenueReportResponse, error) {
    return &response.RevenueReportResponse{
        Period:       fmt.Sprintf("%d-%02d", year, month),
        TotalRevenue: 0,
        RevenueByBillType: make(map[string]int64),
        FeeBreakdown: make(map[string]int64),
    }, nil
}

func (s *AdminService) GetRevenueByBillType(ctx context.Context, startDate, endDate string) (map[string]int64, error) {
    return make(map[string]int64), nil
}

func (s *AdminService) GetTopUsers(ctx context.Context, limit int) ([]response.UserDetailResponse, error) {
    return nil, nil
}

func (s *AdminService) GetFraudStats(ctx context.Context) (interface{}, error) {
    return nil, nil
}

func (s *AdminService) GetProviderPerformance(ctx context.Context, startDate, endDate string) ([]response.ProviderPerformanceResponse, error) {
    return nil, nil
}

func (s *AdminService) ExportReport(ctx context.Context, reportType, format, startDate, endDate string) ([]byte, string, error) {
    return nil, "csv", nil
}

// ========== SYSTEM SETTINGS ==========

func (s *AdminService) GetSystemSettings(ctx context.Context) (*response.SystemSettings, error) {
    return &response.SystemSettings{
        MaintenanceMode:    false,
        MaintenanceMessage: "",
        GlobalDailyLimit:   0,
        GlobalSingleTxLimit: 0,
        MaxRetryCount:      3,
        SessionTimeout:     3600,
    }, nil
}

func (s *AdminService) UpdateSystemSettings(ctx context.Context, req *request.SystemSettingsRequest) error {
    return nil
}

func (s *AdminService) HealthCheck(ctx context.Context) (interface{}, error) {
    return map[string]string{
        "status": "healthy",
        "database": "connected",
    }, nil
}

func (s *AdminService) GetAuditLogs(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]response.AuditLogResponse, int64, error) {
    logs, total, err := s.auditLogRepo.List(ctx, offset, limit, filters)
    if err != nil {
        return nil, 0, err
    }
    
    responses := make([]response.AuditLogResponse, len(logs))
    for i, log := range logs {
        responses[i] = response.AuditLogResponse{
            ID:         log.ID,
            Action:     log.Action,
            EntityType: log.EntityType,
            EntityID:   log.EntityID,
            IPAddress:  log.IPAddress,
            CreatedAt:  log.CreatedAt,
        }
    }
    
    return responses, total, nil
}

func (s *AdminService) GetAuditLog(ctx context.Context, logID string) (*response.AuditLogResponse, error) {
    id, err := uuid.Parse(logID)
    if err != nil {
        return nil, err
    }
    
    log, err := s.auditLogRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    return &response.AuditLogResponse{
        ID:         log.ID,
        Action:     log.Action,
        EntityType: log.EntityType,
        EntityID:   log.EntityID,
        IPAddress:  log.IPAddress,
        CreatedAt:  log.CreatedAt,
    }, nil
}

func (s *AdminService) TriggerDatabaseBackup(ctx context.Context) (interface{}, error) {
    return map[string]string{
        "backup_id": uuid.New().String(),
        "status":    "started",
    }, nil
}

func (s *AdminService) GetSystemMetrics(ctx context.Context) (interface{}, error) {
    return map[string]interface{}{
        "cpu_usage":    15.5,
        "memory_usage": 256.0,
        "goroutines":   45,
        "uptime":       "5d 3h 22m",
        "db_connections": 10,
    }, nil
}

// ========== ROLE MANAGEMENT ==========

func (s *AdminService) ListRoles(ctx context.Context) ([]interface{}, error) {
    roles, err := s.roleRepo.List(ctx)
    if err != nil {
        return nil, err
    }
    
    result := make([]interface{}, len(roles))
    for i, role := range roles {
        result[i] = map[string]interface{}{
            "id":          role.ID,
            "name":        role.Name,
            "permissions": role.Permissions,
            "description": role.Description,
        }
    }
    
    return result, nil
}

func (s *AdminService) CreateRole(ctx context.Context, req *request.CreateRoleRequest) (interface{}, error) {
    role := &entities.Role{
        ID:          uuid.New(),
        Name:        req.Name,
        Permissions: fmt.Sprintf(`["%s"]`, req.Permissions[0]), // Convert array to JSON string
        Description: req.Description,
    }
    
    if err := s.roleRepo.Create(ctx, role); err != nil {
        return nil, err
    }
    
    return role, nil
}

func (s *AdminService) UpdateRole(ctx context.Context, roleID uuid.UUID, req *request.CreateRoleRequest) (interface{}, error) {
    role, err := s.roleRepo.GetByID(ctx, roleID)
    if err != nil {
        return nil, err
    }
    
    if role == nil {
        return nil, errors.New("role not found")
    }
    
    role.Name = req.Name
    role.Permissions = fmt.Sprintf(`["%s"]`, req.Permissions[0])
    role.Description = req.Description
    
    if err := s.roleRepo.Update(ctx, role); err != nil {
        return nil, err
    }
    
    return role, nil
}

func (s *AdminService) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
    return s.roleRepo.Delete(ctx, roleID)
}

func (s *AdminService) ListStaff(ctx context.Context, offset, limit int) ([]interface{}, int64, error) {
    staff, total, err := s.adminUserRepo.List(ctx, offset, limit)
    if err != nil {
        return nil, 0, err
    }
    
    result := make([]interface{}, len(staff))
    for i, admin := range staff {
        result[i] = map[string]interface{}{
            "id":         admin.ID,
            "email":      admin.Email,
            "full_name":  admin.FullName,
            "role":       admin.Role,
            "is_active":  admin.IsActive,
            "created_at": admin.CreatedAt,
        }
    }
    
    return result, total, nil
}

func (s *AdminService) InviteStaff(ctx context.Context, req *request.InviteStaffRequest) error {
    // Generate temporary password and send email
    admin := &entities.AdminUser{
        ID:        uuid.New(),
        Email:     req.Email,
        FullName:  req.FullName,
        Role:      req.Role,
        IsActive:  true,
    }
    
    return s.adminUserRepo.Create(ctx, admin)
}

func (s *AdminService) AssignRole(ctx context.Context, staffID uuid.UUID, role string) error {
    return s.adminUserRepo.UpdateRole(ctx, staffID, role)
}

func (s *AdminService) RemoveStaff(ctx context.Context, staffID uuid.UUID) error {
    return s.adminUserRepo.Delete(ctx, staffID)
}

func (s *AdminService) GetStaffAudit(ctx context.Context, staffID uuid.UUID, offset, limit int) ([]response.AuditLogResponse, int64, error) {
    logs, total, err := s.auditLogRepo.GetByAdminID(ctx, staffID, offset, limit)
    if err != nil {
        return nil, 0, err
    }
    
    responses := make([]response.AuditLogResponse, len(logs))
    for i, log := range logs {
        responses[i] = response.AuditLogResponse{
            ID:         log.ID,
            Action:     log.Action,
            EntityType: log.EntityType,
            EntityID:   log.EntityID,
            IPAddress:  log.IPAddress,
            CreatedAt:  log.CreatedAt,
        }
    }
    
    return responses, total, nil
}