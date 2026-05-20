package services

import (
	"fmt"
	
    "context"
    "errors"
    
    "github.com/google/uuid"
    "gorm.io/gorm"
    
    "fintech_api_golang/internal/core/entities"
    "fintech_api_golang/internal/core/interfaces"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
)

type WalletService struct {
    walletRepo      interfaces.WalletRepository
    transactionRepo interfaces.TransactionRepository
    userRepo        interfaces.UserRepository
    db              *gorm.DB
}

func NewWalletService(
    walletRepo interfaces.WalletRepository,
    transactionRepo interfaces.TransactionRepository,
    userRepo interfaces.UserRepository,
    db *gorm.DB,
) *WalletService {
    return &WalletService{
        walletRepo:      walletRepo,
        transactionRepo: transactionRepo,
        userRepo:        userRepo,
        db:              db,
    }
}

func (s *WalletService) CreateWallet(ctx context.Context, userID uuid.UUID, currency string) (*response.WalletResponse, error) {
    // Check if wallet already exists
    existing, err := s.walletRepo.GetByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    if existing != nil {
        return nil, errors.New("wallet already exists for this user")
    }
    
    if currency == "" {
        currency = "NGN"
    }
    
    wallet := &entities.Wallet{
        ID:       uuid.New(),
        UserID:   userID,
        Balance:  0,
        Currency: currency,
    }
    
    if err := s.walletRepo.Create(ctx, wallet); err != nil {
        return nil, err
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

func (s *WalletService) GetBalance(ctx context.Context, userID uuid.UUID) (*response.BalanceResponse, error) {
    wallet, err := s.walletRepo.GetByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    if wallet == nil {
        return nil, errors.New("wallet not found")
    }
    
    return &response.BalanceResponse{
        Balance:      int64(wallet.Balance),
        BalanceNaira: float64(wallet.Balance) / 100,
        Currency:     wallet.Currency,
        IsLocked:     wallet.IsLocked,
    }, nil
}

func (s *WalletService) GetTransactions(ctx context.Context, userID uuid.UUID, req *request.GetTransactionsRequest) (*response.TransactionHistoryResponse, error) {
    wallet, err := s.walletRepo.GetByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    if wallet == nil {
        return nil, errors.New("wallet not found")
    }
    
    page := req.Page
    if page < 1 {
        page = 1
    }
    
    limit := req.Limit
    if limit < 1 {
        limit = 20
    }
    if limit > 100 {
        limit = 100
    }
    
    offset := (page - 1) * limit
    
    filters := make(map[string]interface{})
    if req.Category != "" {
        filters["category"] = req.Category
    }
    if req.Status != "" {
        filters["status"] = req.Status
    }
    if req.FromDate != "" {
        filters["from_date"] = req.FromDate
    }
    if req.ToDate != "" {
        filters["to_date"] = req.ToDate
    }
    
    transactions, total, err := s.transactionRepo.GetByWalletID(ctx, wallet.ID, offset, limit, filters)
    if err != nil {
        return nil, err
    }
    
    totalPages := int(total) / limit
    if int(total)%limit > 0 {
        totalPages++
    }
    
    txResponses := make([]response.TransactionResponse, len(transactions))
    for i, tx := range transactions {
        txResponses[i] = s.mapTransactionToResponse(&tx)
    }
    
    return &response.TransactionHistoryResponse{
        Transactions: txResponses,
        Total:        total,
        Page:         page,
        Limit:        limit,
        TotalPages:   totalPages,
    }, nil
}

func (s *WalletService) GetLimits(ctx context.Context, userID uuid.UUID) (*response.TierLimitResponse, error) {
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    if user == nil {
        return nil, errors.New("user not found")
    }
    
    wallet, err := s.walletRepo.GetByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    if wallet == nil {
        return nil, errors.New("wallet not found")
    }
    
    // Get tier limits from config or database
    limits := s.getTierLimits(user.Tier)
    
    return &response.TierLimitResponse{
        Tier:             user.Tier,
        DailyLimit:       limits.DailyLimit,
        WeeklyLimit:      limits.WeeklyLimit,
        MonthlyLimit:     limits.MonthlyLimit,
        SingleTxLimit:    limits.SingleTxLimit,
        DailySpent:       int64(wallet.DailySpent),
        WeeklySpent:      int64(wallet.WeeklySpent),
        MonthlySpent:     int64(wallet.MonthlySpent),
        DailyRemaining:   int64(limits.DailyLimit) - int64(wallet.DailySpent),
        WeeklyRemaining:  int64(limits.WeeklyLimit) - int64(wallet.WeeklySpent),
        MonthlyRemaining: int64(limits.MonthlyLimit) - int64(wallet.MonthlySpent),
    }, nil
}

func (s *WalletService) LockWallet(ctx context.Context, userID uuid.UUID, reason string) error {
    wallet, err := s.walletRepo.GetByUserID(ctx, userID)
    if err != nil {
        return err
    }
    
    if wallet == nil {
        return errors.New("wallet not found")
    }
    
    return s.walletRepo.Lock(ctx, wallet.ID, reason)
}

func (s *WalletService) UnlockWallet(ctx context.Context, userID uuid.UUID) error {
    wallet, err := s.walletRepo.GetByUserID(ctx, userID)
    if err != nil {
        return err
    }
    
    if wallet == nil {
        return errors.New("wallet not found")
    }
    
    return s.walletRepo.Unlock(ctx, wallet.ID)
}

func (s *WalletService) GetStatement(ctx context.Context, userID uuid.UUID, req *request.GetStatementRequest) ([]byte, string, error) {
    wallet, err := s.walletRepo.GetByUserID(ctx, userID)
    if err != nil {
        return nil, "", err
    }
    
    if wallet == nil {
        return nil, "", errors.New("wallet not found")
    }
    
    filters := make(map[string]interface{})
    filters["from_date"] = req.FromDate
    filters["to_date"] = req.ToDate
    
    transactions, _, err := s.transactionRepo.GetByWalletID(ctx, wallet.ID, 0, 10000, filters)
    if err != nil {
        return nil, "", err
    }
    
    // Generate statement (CSV for simplicity)
    if req.Format == "csv" {
        csv := s.generateCSVStatement(transactions, wallet)
        return []byte(csv), "text/csv", nil
    }
    
    // For PDF, you'd use a PDF generator library
    // Return CSV as fallback
    csv := s.generateCSVStatement(transactions, wallet)
    return []byte(csv), "text/csv", nil
}

// Private helper methods
func (s *WalletService) mapTransactionToResponse(tx *entities.Transaction) response.TransactionResponse {
    resp := response.TransactionResponse{
        ID:              tx.ID,
        Reference:       tx.Reference,
        Type:            tx.Type,
        Category:        tx.Category,
        SubCategory:     tx.SubCategory,
        Amount:          int64(tx.Amount),
        AmountNaira:     float64(tx.Amount) / 100,
        Fee:             int64(tx.Fee),
        FeeNaira:        float64(tx.Fee) / 100,
        TotalAmount:     int64(tx.TotalAmount),
        TotalAmountNaira: float64(tx.TotalAmount) / 100,
        Status:          tx.Status,
        Description:     tx.Description,
        BalanceBefore:   int64(tx.BalanceBefore),
        BalanceAfter:    int64(tx.BalanceAfter),
        CompletedAt:     tx.CompletedAt,
        CreatedAt:       tx.CreatedAt,
    }
    
    if tx.TransferDetail != nil {
        resp.TransferDetail = &response.TransferDetailResponse{
            RecipientType: tx.TransferDetail.RecipientType,
            RecipientID:   tx.TransferDetail.RecipientID,
            RecipientName: tx.TransferDetail.RecipientName,
            RecipientBank: tx.TransferDetail.RecipientBankName,
            Narration:     tx.TransferDetail.Narration,
        }
    }
    
    if tx.BillDetail != nil {
        resp.BillDetail = &response.BillDetailResponse{
            BillType:      tx.BillDetail.BillType,
            ProviderName:  tx.BillDetail.ProviderName,
            PhoneNumber:   tx.BillDetail.PhoneNumber,
            MeterNumber:   tx.BillDetail.MeterNumber,
            CustomerName:  tx.BillDetail.CustomerName,
            ElectricityToken: tx.BillDetail.ElectricityToken,
            DataPlanName:  tx.BillDetail.DataPlanName,
            DataVolume:    tx.BillDetail.DataVolume,
        }
    }
    
    return resp
}

func (s *WalletService) generateCSVStatement(transactions []entities.Transaction, wallet *entities.Wallet) string {
    csv := "Date,Reference,Type,Category,Amount (NGN),Fee (NGN),Total (NGN),Status,Balance Before (NGN),Balance After (NGN),Description\n"
    
    for _, tx := range transactions {
        csv += tx.CreatedAt.Format("2006-01-02 15:04:05") + ","
        csv += tx.Reference + ","
        csv += tx.Type + ","
        csv += tx.Category + ","
        csv += formatNaira(int64(tx.Amount)) + ","
        csv += formatNaira(int64(tx.Fee)) + ","
        csv += formatNaira(int64(tx.TotalAmount)) + ","
        csv += tx.Status + ","
        csv += formatNaira(int64(tx.BalanceBefore)) + ","
        csv += formatNaira(int64(tx.BalanceAfter)) + ","
        csv += tx.Description + "\n"
    }
    
    return csv
}

func (s *WalletService) getTierLimits(tier int) *TierLimits {
    limits := map[int]*TierLimits{
        0: {DailyLimit: 0, WeeklyLimit: 0, MonthlyLimit: 0, SingleTxLimit: 0},
        1: {DailyLimit: 5000000, WeeklyLimit: 20000000, MonthlyLimit: 50000000, SingleTxLimit: 2500000},
        2: {DailyLimit: 20000000, WeeklyLimit: 100000000, MonthlyLimit: 300000000, SingleTxLimit: 10000000},
        3: {DailyLimit: 500000000, WeeklyLimit: 2000000000, MonthlyLimit: 5000000000, SingleTxLimit: 500000000},
    }
    
    if limit, ok := limits[tier]; ok {
        return limit
    }
    
    return limits[0]
}

type TierLimits struct {
    DailyLimit    int64
    WeeklyLimit   int64
    MonthlyLimit  int64
    SingleTxLimit int64
}

func formatNaira(amount int64) string {
    naira := float64(amount) / 100
    return "₦" + formatFloat(naira)
}

func formatFloat(f float64) string {
    return fmt.Sprintf("%.2f", f)
}