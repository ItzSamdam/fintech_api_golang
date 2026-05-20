package services

import (
    "context"
    "errors"
    "fmt"
    "time"
    
    "github.com/google/uuid"
    "gorm.io/gorm"
    
    "fintech_api_golang/internal/config"
    "fintech_api_golang/internal/core/entities"
    "fintech_api_golang/internal/core/interfaces"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
    "fintech_api_golang/internal/repositories/providers"
)

type TransferService struct {
    walletRepo      interfaces.WalletRepository
    transactionRepo interfaces.TransactionRepository
    transferDetailRepo interfaces.TransferDetailRepository
    userRepo        interfaces.UserRepository
    redBiller       *providers.RedBillerClient
    db              *gorm.DB
    config          *config.Config
}

func NewTransferService(
    walletRepo interfaces.WalletRepository,
    transactionRepo interfaces.TransactionRepository,
    transferDetailRepo interfaces.TransferDetailRepository,
    userRepo interfaces.UserRepository,
    redBiller *providers.RedBillerClient,
    db *gorm.DB,
    cfg *config.Config,
) *TransferService {
    return &TransferService{
        walletRepo:        walletRepo,
        transactionRepo:   transactionRepo,
        transferDetailRepo: transferDetailRepo,
        userRepo:          userRepo,
        redBiller:         redBiller,
        db:                db,
        config:            cfg,
    }
}

// SendTransfer - Handles both in-app and external transfers
func (s *TransferService) SendTransfer(ctx context.Context, userID uuid.UUID, req *request.SendTransferRequest) (*response.TransactionResponse, error) {
    // Get user's wallet
    wallet, err := s.walletRepo.GetByUserIDForUpdate(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    if wallet == nil {
        return nil, errors.New("wallet not found")
    }
    
    if wallet.IsLocked {
        return nil, errors.New("wallet is locked")
    }
    
    var recipientName string
    var recipientWallet *entities.Wallet
    
    // Handle different recipient types
    if req.RecipientType == "bank" {
        // External transfer - verify bank account first
        nameEnquiry, err := s.NameEnquiry(ctx, &request.NameEnquiryRequest{
            AccountNumber: req.RecipientID,
            BankCode:      req.RecipientBankCode,
        })
        if err != nil {
            return nil, errors.New("account verification failed: " + err.Error())
        }
        recipientName = nameEnquiry.AccountName
    } else if req.RecipientType == "wallet" {
        // In-app transfer - verify recipient wallet exists
        recipientWallet, err = s.walletRepo.GetByUserID(ctx, uuid.MustParse(req.RecipientID))
        if err != nil || recipientWallet == nil {
            return nil, errors.New("recipient wallet not found")
        }
        
        // Get recipient user for name
        recipientUser, err := s.userRepo.GetByID(ctx, recipientWallet.UserID)
        if err == nil && recipientUser != nil {
            recipientName = recipientUser.PhoneNumber
        }
        
        // Prevent self-transfer
        if recipientWallet.UserID == userID {
            return nil, errors.New("cannot transfer to yourself")
        }
    }
    
    // Check sufficient balance
    if int64(wallet.Balance) < req.Amount {
        return nil, errors.New("insufficient balance")
    }
    
    // Calculate fee
    fee := s.calculateTransferFee(req.Amount)
    vat := int64(float64(fee) * 0.075)
    totalAmount := req.Amount + fee + vat
    
    if int64(wallet.Balance) < totalAmount {
        return nil, errors.New("insufficient balance including fees")
    }
    
    reference := s.generateReference("TRF")
    
    var transaction *entities.Transaction
    
    err = s.db.Transaction(func(tx *gorm.DB) error {
        // Debit sender's wallet
        if err := s.walletRepo.Debit(ctx, wallet.ID, totalAmount, reference); err != nil {
            return err
        }
        
        // Update spent limits
        if err := s.walletRepo.UpdateSpentLimits(ctx, wallet.ID, req.Amount); err != nil {
            return err
        }
        
        // For in-app transfer, credit recipient immediately
        if req.RecipientType == "wallet" && recipientWallet != nil {
            if err := s.walletRepo.Credit(ctx, recipientWallet.ID, req.Amount, reference); err != nil {
                return err
            }
        }
        
        // Create transaction record
        transaction = &entities.Transaction{
            ID:            uuid.New(),
            Reference:     reference,
            WalletID:      wallet.ID,
            UserID:        userID,
            Type:          "debit",
            Category:      "transfer",
            Amount:        entities.AmountInKobo(req.Amount),
            Fee:           entities.AmountInKobo(fee),
            VAT:           entities.AmountInKobo(vat),
            TotalAmount:   entities.AmountInKobo(totalAmount),
            BalanceBefore: wallet.Balance,
            BalanceAfter:  entities.AmountInKobo(int64(wallet.Balance) - totalAmount),
            Status:        "processing",
            Description:   req.Narration,
            CreatedAt:     time.Now(),
        }
        
        if err := s.transactionRepo.Create(ctx, transaction); err != nil {
            return err
        }
        
        // Create transfer detail
        transferDetail := &entities.TransferDetail{
            ID:            uuid.New(),
            TransactionID: transaction.ID,
            RecipientType: req.RecipientType,
            RecipientID:   req.RecipientID,
            RecipientName: recipientName,
            Narration:     req.Narration,
        }
        
        if req.RecipientType == "bank" {
            transferDetail.RecipientBankCode = req.RecipientBankCode
        }
        
        if err := s.transferDetailRepo.Create(ctx, transferDetail); err != nil {
            return err
        }
        
        // If in-app transfer, mark as success immediately
        if req.RecipientType == "wallet" {
            now := time.Now()
            transaction.Status = "success"
            transaction.CompletedAt = &now
            return s.transactionRepo.Update(ctx, transaction)
        }
        
        return nil
    })
    
    if err != nil {
        return nil, err
    }
    
    // Only process external transfer via provider
    if req.RecipientType == "bank" {
        go s.processExternalTransfer(context.Background(), transaction, req)
    } else {
        // For in-app, return success immediately
        return s.mapTransactionToResponse(transaction), nil
    }
    
    return s.mapTransactionToResponse(transaction), nil
}

// SendToWallet - Direct wallet to wallet transfer
func (s *TransferService) SendToWallet(ctx context.Context, userID uuid.UUID, recipientWalletID uuid.UUID, amount int64, narration string) (*response.TransactionResponse, error) {
    // Get sender's wallet
    senderWallet, err := s.walletRepo.GetByUserIDForUpdate(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    if senderWallet == nil {
        return nil, errors.New("wallet not found")
    }
    
    // Get recipient's wallet
    recipientWallet, err := s.walletRepo.GetByID(ctx, recipientWalletID)
    if err != nil || recipientWallet == nil {
        return nil, errors.New("recipient wallet not found")
    }
    
    // Check sufficient balance
    if int64(senderWallet.Balance) < amount {
        return nil, errors.New("insufficient balance")
    }
    
    reference := s.generateReference("WTR")
    
    var transaction *entities.Transaction
    
    err = s.db.Transaction(func(tx *gorm.DB) error {
        // Debit sender
        if err := s.walletRepo.Debit(ctx, senderWallet.ID, amount, reference); err != nil {
            return err
        }
        
        // Credit recipient
        if err := s.walletRepo.Credit(ctx, recipientWallet.ID, amount, reference); err != nil {
            return err
        }
        
        // Create sender transaction
        transaction = &entities.Transaction{
            ID:            uuid.New(),
            Reference:     reference,
            WalletID:      senderWallet.ID,
            UserID:        userID,
            Type:          "debit",
            Category:      "transfer",
            Amount:        entities.AmountInKobo(amount),
            Fee:           0,
            VAT:           0,
            TotalAmount:   entities.AmountInKobo(amount),
            BalanceBefore: senderWallet.Balance,
            BalanceAfter:  entities.AmountInKobo(int64(senderWallet.Balance) - amount),
            Status:        "success",
            Description:   narration,
            CompletedAt:   timeNow(),
            CreatedAt:     time.Now(),
        }
        
        if err := s.transactionRepo.Create(ctx, transaction); err != nil {
            return err
        }
        
        // Create transfer detail
        transferDetail := &entities.TransferDetail{
            ID:            uuid.New(),
            TransactionID: transaction.ID,
            RecipientType: "wallet",
            RecipientID:   recipientWallet.UserID.String(),
            Narration:     narration,
        }
        
        return s.transferDetailRepo.Create(ctx, transferDetail)
    })
    
    if err != nil {
        return nil, err
    }
    
    return s.mapTransactionToResponse(transaction), nil
}


// processExternalTransfer - Only called for bank transfers
func (s *TransferService) processExternalTransfer(ctx context.Context, transaction *entities.Transaction, req *request.SendTransferRequest) {
    // Verify account details again before sending (fresh validation)
    nameEnquiry, err := s.NameEnquiry(ctx, &request.NameEnquiryRequest{
        AccountNumber: req.RecipientID,
        BankCode:      req.RecipientBankCode,
    })
    
    if err != nil || nameEnquiry.AccountName == "" {
        errorMsg := "account verification failed before transfer"
        if err != nil {
            errorMsg = err.Error()
        }
        s.transactionRepo.MarkAsFailed(ctx, transaction.Reference, errorMsg)
        return
    }
    
    // Call RedBiller API to send money
    resp, err := s.redBiller.SendMoney(ctx, &providers.SendMoneyRequest{
        AccountNo:   req.RecipientID,
        BankCode:    req.RecipientBankCode,
        Amount:      req.Amount,
        Narration:   req.Narration,
        CallbackURL: "https://domain.com/api/callback", //s.config.Webhook.CallbackURL,
        Reference:   transaction.Reference,
    })
    
    if err != nil || (resp != nil && !resp.Success) {
        errorMsg := "transfer failed"
        if resp != nil {
            errorMsg = resp.Message
        }
        s.transactionRepo.MarkAsFailed(ctx, transaction.Reference, errorMsg)
        return
    }
    
    now := time.Now()
    s.transactionRepo.UpdateStatus(ctx, transaction.Reference, "success", &now)
}

func (s *TransferService) GetTransferStatus(ctx context.Context, userID uuid.UUID, reference string) (*response.TransactionResponse, error) {
    transaction, err := s.transactionRepo.GetByReference(ctx, reference)
    if err != nil {
        return nil, err
    }
    
    if transaction == nil {
        return nil, errors.New("transaction not found")
    }
    
    if transaction.UserID != userID {
        return nil, errors.New("unauthorized to view this transaction")
    }
    
    return s.mapTransactionToResponse(transaction), nil
}

func (s *TransferService) RetryTransfer(ctx context.Context, userID uuid.UUID, reference string) (*response.TransactionResponse, error) {
    transaction, err := s.transactionRepo.GetByReference(ctx, reference)
    if err != nil {
        return nil, err    }
    
    if transaction == nil {
        return nil, errors.New("transaction not found")
    }
    
    if transaction.UserID != userID {
        return nil, errors.New("unauthorized to retry this transaction")
    }
    
    if transaction.Status != "failed" {
        return nil, errors.New("only failed transactions can be retried")
    }
    
    if transaction.RetryCount >= 3 {
        return nil, errors.New("maximum retry attempts reached")
    }
    
    // Increment retry count and reset status
    transaction.RetryCount++
    transaction.Status = "processing"
    if err := s.transactionRepo.Update(ctx, transaction); err != nil {
        return nil, err
    }
    
    // Retry the transfer
    go s.processExternalTransferRetry(context.Background(), transaction)
    
    return s.mapTransactionToResponse(transaction), nil
}

func (s *TransferService) processExternalTransferRetry(ctx context.Context, transaction *entities.Transaction) {
    resp, err := s.redBiller.RetrySendMoney(ctx, transaction.Reference)
    
    if err != nil || (resp != nil && !resp.Success) {
        errorMsg := "retry failed"
        if resp != nil {
            errorMsg = resp.Message
        }
        s.transactionRepo.MarkAsFailed(ctx, transaction.Reference, errorMsg)
        return
    }
    
    now := time.Now()
    s.transactionRepo.UpdateStatus(ctx, transaction.Reference, "success", &now)
}

func (s *TransferService) NameEnquiry(ctx context.Context, req *request.NameEnquiryRequest) (*response.NameEnquiryResponse, error) {
    resp, err := s.redBiller.VerifyAccountDetails(ctx, req.AccountNumber, req.BankCode)
    if err != nil {
        return nil, err
    }
    
    if !resp.Success {
        return nil, errors.New("account verification failed")
    }
    
    // Extract account name from response
    accountName := ""
    if data, ok := resp.Data["data"].(map[string]interface{}); ok {
        if name, ok := data["account_name"].(string); ok {
            accountName = name
        }
    }
    
    return &response.NameEnquiryResponse{
        AccountNumber: req.AccountNumber,
        AccountName:   accountName,
        BankCode:      req.BankCode,
    }, nil
}

func (s *TransferService) GetBanks(ctx context.Context) (*response.BankListResponse, error) {
    resp, err := s.redBiller.FetchBanks(ctx, "", "NG", "NGN")
    if err != nil {
        return nil, err
    }
    
    banks := []response.Bank{}
    if data, ok := resp.Data["data"].([]interface{}); ok {
        for _, bank := range data {
            if bankMap, ok := bank.(map[string]interface{}); ok {
                banks = append(banks, response.Bank{
                    Code: fmt.Sprintf("%v", bankMap["code"]),
                    Name: fmt.Sprintf("%v", bankMap["name"]),
                })
            }
        }
    }
    
    return &response.BankListResponse{Banks: banks}, nil
}

func (s *TransferService) GetTransferHistory(ctx context.Context, userID uuid.UUID, offset, limit int) (*response.TransactionHistoryResponse, error) {
    filters := map[string]interface{}{
        "category": "transfer",
    }
    
    transactions, total, err := s.transactionRepo.GetByUserID(ctx, userID, offset, limit, filters)
    if err != nil {
        return nil, err
    }
    
    totalPages := int(total) / limit
    if int(total)%limit > 0 {
        totalPages++
    }
    
    txResponses := make([]response.TransactionResponse, len(transactions))
    for i, tx := range transactions {
        txResponses[i] = *s.mapTransactionToResponse(&tx)
    }
    
    return &response.TransactionHistoryResponse{
        Transactions: txResponses,
        Total:        total,
        Page:         offset/limit + 1,
        Limit:        limit,
        TotalPages:   totalPages,
    }, nil
}

func (s *TransferService) generateReference(prefix string) string {
    return fmt.Sprintf("%s%s%d", prefix, time.Now().Format("20060102150405"), time.Now().UnixNano()%10000)
}

func (s *TransferService) calculateTransferFee(amount int64) int64 {
    // 0.5% fee, max 5000 NGN (500000 kobo)
    fee := int64(float64(amount) * 0.005)
    maxFee := int64(500000)
    if fee > maxFee {
        fee = maxFee
    }
    return fee
}

func (s *TransferService) mapTransactionToResponse(tx *entities.Transaction) *response.TransactionResponse {
    resp := &response.TransactionResponse{
        ID:               tx.ID,
        Reference:        tx.Reference,
        Type:             tx.Type,
        Category:         tx.Category,
        SubCategory:      tx.SubCategory,
        Amount:           int64(tx.Amount),
        AmountNaira:      float64(tx.Amount) / 100,
        Fee:              int64(tx.Fee),
        FeeNaira:         float64(tx.Fee) / 100,
        TotalAmount:      int64(tx.TotalAmount),
        TotalAmountNaira: float64(tx.TotalAmount) / 100,
        Status:           tx.Status,
        Description:      tx.Description,
        BalanceBefore:    int64(tx.BalanceBefore),
        BalanceAfter:     int64(tx.BalanceAfter),
        CompletedAt:      tx.CompletedAt,
        CreatedAt:        tx.CreatedAt,
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
    
    return resp
}