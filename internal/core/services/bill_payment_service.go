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
	"fintech_api_golang/internal/repositories/providers"
)

type BillPaymentService struct {
	walletRepo      interfaces.WalletRepository
	transactionRepo interfaces.TransactionRepository
	billDetailRepo  interfaces.BillDetailRepository
	providerRepo    interfaces.ProviderRepository
	userRepo        interfaces.UserRepository
	redBiller       *providers.RedBillerClient
	db              *gorm.DB
}

func NewBillPaymentService(
	walletRepo interfaces.WalletRepository,
	transactionRepo interfaces.TransactionRepository,
	billDetailRepo interfaces.BillDetailRepository,
	providerRepo interfaces.ProviderRepository,
	userRepo interfaces.UserRepository,
	redBiller *providers.RedBillerClient,
	db *gorm.DB,
) *BillPaymentService {
	return &BillPaymentService{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		billDetailRepo:  billDetailRepo,
		providerRepo:    providerRepo,
		userRepo:        userRepo,
		redBiller:       redBiller,
		db:              db,
	}
}

// Airtime
func (s *BillPaymentService) GetAirtimeNetworks(ctx context.Context) (*response.NetworkListResponse, error) {
	networks := []response.Network{
		{Code: "MTN", Name: "MTN Nigeria", IsActive: true},
		{Code: "GLO", Name: "Glo Nigeria", IsActive: true},
		{Code: "AIRTEL", Name: "Airtel Nigeria", IsActive: true},
		{Code: "9MOBILE", Name: "9mobile", IsActive: true},
	}

	return &response.NetworkListResponse{Networks: networks}, nil
}

func (s *BillPaymentService) PurchaseAirtime(ctx context.Context, userID uuid.UUID, req *request.PurchaseAirtimeRequest) (*response.TransactionResponse, error) {
	// Get wallet
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

	// Calculate fee
	fee := s.calculateBillFee("airtime", req.Amount)
	vat := int64(float64(fee) * 0.075)
	totalAmount := req.Amount + fee + vat

	if wallet.Balance < entities.AmountInKobo(totalAmount) {
		return nil, errors.New("insufficient balance")
	}

	reference := s.generateReference("AIR")

	var transaction *entities.Transaction

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Debit wallet
		if err := s.walletRepo.Debit(ctx, wallet.ID, totalAmount, reference); err != nil {
			return err
		}

		// Create transaction
		transaction = &entities.Transaction{
			ID:            uuid.New(),
			Reference:     reference,
			WalletID:      wallet.ID,
			UserID:        userID,
			Type:          "debit",
			Category:      "airtime",
			SubCategory:   req.Network,
			Amount:        entities.AmountInKobo(req.Amount),
			Fee:           entities.AmountInKobo(fee),
			VAT:           entities.AmountInKobo(vat),
			TotalAmount:   entities.AmountInKobo(totalAmount),
			BalanceBefore: wallet.Balance,
			BalanceAfter:  entities.AmountInKobo(wallet.Balance) - entities.AmountInKobo(totalAmount),
			Status:        "processing",
			Description:   fmt.Sprintf("Airtime purchase for %s", req.PhoneNumber),
			CreatedAt:     time.Now(),
		}

		if err := s.transactionRepo.Create(ctx, transaction); err != nil {
			return err
		}

		// Create bill detail
		billDetail := &entities.BillDetail{
			ID:            uuid.New(),
			TransactionID: transaction.ID,
			BillType:      "airtime",
			PhoneNumber:   req.PhoneNumber,
			ProviderName:  req.Network,
		}

		if err := s.billDetailRepo.Create(ctx, billDetail); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Process external purchase
	go s.processAirtimePurchase(ctx, transaction, req)

	return s.mapTransactionToResponse(transaction), nil
}

func (s *BillPaymentService) processAirtimePurchase(ctx context.Context, transaction *entities.Transaction, req *request.PurchaseAirtimeRequest) {
	productMap := map[string]string{
		"MTN":     "MTN",
		"GLO":     "GLO",
		"AIRTEL":  "AIRTEL",
		"9MOBILE": "9MOBILE",
	}

	product := productMap[req.Network]
	if product == "" {
		product = req.Network
	}

	resp, err := s.redBiller.PurchaseTopUp(ctx, &providers.PurchaseTopUpRequest{
		Product:     product,
		PhoneNo:     req.PhoneNumber,
		Amount:      req.Amount,
		Ported:      false,
		CallbackURL: "https://your-callback.com/webhook",
		Reference:   transaction.Reference,
	})

	if err != nil || (resp != nil && !resp.Success) {
		errorMsg := "airtime purchase failed"
		if resp != nil {
			errorMsg = resp.Message
		}
		s.transactionRepo.MarkAsFailed(ctx, transaction.Reference, errorMsg)
		return
	}

	now := time.Now()
	s.transactionRepo.UpdateStatus(ctx, transaction.Reference, "success", &now)
}

// Data
func (s *BillPaymentService) GetDataNetworks(ctx context.Context) (*response.NetworkListResponse, error) {
	networks := []response.Network{
		{Code: "MTN", Name: "MTN Data", IsActive: true},
		{Code: "GLO", Name: "Glo Data", IsActive: true},
		{Code: "AIRTEL", Name: "Airtel Data", IsActive: true},
		{Code: "9MOBILE", Name: "9mobile Data", IsActive: true},
	}

	return &response.NetworkListResponse{Networks: networks}, nil
}

func (s *BillPaymentService) GetDataPlans(ctx context.Context, network string) ([]response.DataPlanResponse, error) {
	resp, err := s.redBiller.GetDataPlans(ctx, network)
	if err != nil {
		return nil, err
	}

	plans := []response.DataPlanResponse{}
	if data, ok := resp.Data["data"].([]interface{}); ok {
		for _, plan := range data {
			if planMap, ok := plan.(map[string]interface{}); ok {
				plans = append(plans, response.DataPlanResponse{
					ID:         fmt.Sprintf("%v", planMap["code"]),
					Name:       fmt.Sprintf("%v", planMap["name"]),
					Volume:     fmt.Sprintf("%v", planMap["size"]),
					Price:      int64(planMap["amount"].(float64)),
					PriceNaira: planMap["amount"].(float64),
					Validity:   fmt.Sprintf("%v", planMap["validity"]),
					Network:    network,
				})
			}
		}
	}

	return plans, nil
}

func (s *BillPaymentService) PurchaseData(ctx context.Context, userID uuid.UUID, req *request.PurchaseDataRequest) (*response.TransactionResponse, error) {
	// Similar to airtime purchase
	// ... (implementation similar to PurchaseAirtime)
	return nil, nil
}

// Electricity
func (s *BillPaymentService) GetElectricityProviders(ctx context.Context) ([]response.ProviderResponse, error) {
	providers, err := s.providerRepo.GetActiveByType(ctx, "electricity")
	if err != nil {
		return nil, err
	}

	result := make([]response.ProviderResponse, len(providers))
	for i, p := range providers {
		result[i] = response.ProviderResponse{
			ID:       p.ID.String(),
			Name:     p.Name,
			Code:     p.Code,
			Category: p.Category,
			IsActive: p.IsActive,
		}
	}

	return result, nil
}

func (s *BillPaymentService) ValidateMeter(ctx context.Context, req *request.ValidateMeterRequest) (*response.MeterValidationResponse, error) {
	resp, err := s.redBiller.VerifyDisco(ctx, req.ProviderID, req.MeterNumber, req.MeterType)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New("meter validation failed")
	}

	customerName := ""
	customerAddress := ""
	if data, ok := resp.Data["data"].(map[string]interface{}); ok {
		if name, ok := data["customer_name"].(string); ok {
			customerName = name
		}
		if address, ok := data["address"].(string); ok {
			customerAddress = address
		}
	}

	return &response.MeterValidationResponse{
		CustomerName:    customerName,
		CustomerAddress: customerAddress,
		MeterNumber:     req.MeterNumber,
		MeterType:       req.MeterType,
		ProviderID:      req.ProviderID,
	}, nil
}

func (s *BillPaymentService) PayElectricity(ctx context.Context, userID uuid.UUID, req *request.PayElectricityRequest) (*response.ElectricityPaymentResponse, error) {
	// Similar to airtime purchase but with electricity-specific logic
	// ... (implementation similar to PurchaseAirtime)
	return nil, nil
}

func (s *BillPaymentService) GetElectricityToken(ctx context.Context, userID uuid.UUID, transactionID string) (string, error) {
	txnID, err := uuid.Parse(transactionID)
	if err != nil {
		return "", err
	}

	billDetail, err := s.billDetailRepo.GetByTransactionID(ctx, txnID)
	if err != nil {
		return "", err
	}

	if billDetail == nil {
		return "", errors.New("transaction not found")
	}

	return billDetail.ElectricityToken, nil
}

// Betting
func (s *BillPaymentService) GetBettingProviders(ctx context.Context) ([]response.ProviderResponse, error) {
	resp, err := s.redBiller.GetBetProviders(ctx)
	if err != nil {
		return nil, err
	}

	providers := []response.ProviderResponse{}
	if data, ok := resp.Data["data"].([]interface{}); ok {
		for _, provider := range data {
			if providerMap, ok := provider.(map[string]interface{}); ok {
				providers = append(providers, response.ProviderResponse{
					ID:   fmt.Sprintf("%v", providerMap["code"]),
					Name: fmt.Sprintf("%v", providerMap["name"]),
					Code: fmt.Sprintf("%v", providerMap["code"]),
				})
			}
		}
	}

	return providers, nil
}

func (s *BillPaymentService) ValidateBettingAccount(ctx context.Context, req *request.ValidateBettingAccountRequest) (*response.BettingAccountResponse, error) {
	resp, err := s.redBiller.VerifyBetWallet(ctx, req.ProviderID, req.AccountID)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New("account validation failed")
	}

	accountName := ""
	if data, ok := resp.Data["data"].(map[string]interface{}); ok {
		if name, ok := data["customer_name"].(string); ok {
			accountName = name
		}
	}

	return &response.BettingAccountResponse{
		ProviderID:  req.ProviderID,
		AccountID:   req.AccountID,
		AccountName: accountName,
		IsValid:     true,
	}, nil
}

func (s *BillPaymentService) FundBettingWallet(ctx context.Context, userID uuid.UUID, req *request.FundBettingRequest) (*response.TransactionResponse, error) {
	// Similar to airtime purchase
	// ... (implementation similar to PurchaseAirtime)
	return nil, nil
}

func (s *BillPaymentService) GetBettingHistory(ctx context.Context, userID uuid.UUID, offset, limit int) (*response.TransactionHistoryResponse, error) {
	filters := map[string]interface{}{
		"category": "betting",
	}

	transactions, total, err := s.transactionRepo.GetByUserID(ctx, userID, offset, limit, filters)
	if err != nil {
		return nil, err
	}

	txResponses := make([]response.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		txResponses[i] = *s.mapTransactionToResponse(&tx)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &response.TransactionHistoryResponse{
		Transactions: txResponses,
		Total:        total,
		Page:         offset/limit + 1,
		Limit:        limit,
		TotalPages:   totalPages,
	}, nil
}

// Helper methods
func (s *BillPaymentService) generateReference(prefix string) string {
	return fmt.Sprintf("%s%s%d", prefix, time.Now().Format("20060102150405"), time.Now().UnixNano()%10000)
}

func (s *BillPaymentService) calculateBillFee(billType string, amount int64) int64 {
	feePercentages := map[string]float64{
		"airtime":     0.01,
		"data":        0.01,
		"electricity": 100, // Fixed 100 NGN fee
		"betting":     0.015,
	}

	percentage, ok := feePercentages[billType]
	if !ok {
		return 0
	}

	if billType == "electricity" {
		return 10000 // 100 NGN in kobo
	}

	fee := int64(float64(amount) * percentage)
	maxFee := int64(100000) // Max 1000 NGN
	if fee > maxFee {
		fee = maxFee
	}

	return fee
}


func (s *BillPaymentService) GetDataHistory(ctx context.Context, userID uuid.UUID, offset, limit int) (*response.TransactionHistoryResponse, error) {
	filters := map[string]interface{}{
		"category": "data",
	}

	transactions, total, err := s.transactionRepo.GetByUserID(ctx, userID, offset, limit, filters)
	if err != nil {
		return nil, err
	}

	txResponses := make([]response.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		txResponses[i] = *s.mapTransactionToResponse(&tx)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &response.TransactionHistoryResponse{
		Transactions: txResponses,
		Total:        total,
		Page:         offset/limit + 1,
		Limit:        limit,
		TotalPages:   totalPages,
	}, nil
}


func (s *BillPaymentService) GetElectricityHistory(ctx context.Context, userID uuid.UUID, offset, limit int) (*response.TransactionHistoryResponse, error) {
	filters := map[string]interface{}{
		"category": "electricity",
	}

	transactions, total, err := s.transactionRepo.GetByUserID(ctx, userID, offset, limit, filters)
	if err != nil {
		return nil, err
	}

	txResponses := make([]response.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		txResponses[i] = *s.mapTransactionToResponse(&tx)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &response.TransactionHistoryResponse{
		Transactions: txResponses,
		Total:        total,
		Page:         offset/limit + 1,
		Limit:        limit,
		TotalPages:   totalPages,
	}, nil
}


func (s *BillPaymentService) GetBillHistory(ctx context.Context, userID uuid.UUID, offset, limit int) (*response.TransactionHistoryResponse, error) {
	filters := map[string]interface{}{
		"category": []string{
            "airtime", 
            "data", 
            "electricity", 
            "betting",
        },
	}

	transactions, total, err := s.transactionRepo.GetByUserID(ctx, userID, offset, limit, filters)
	if err != nil {
		return nil, err
	}

	txResponses := make([]response.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		txResponses[i] = *s.mapTransactionToResponse(&tx)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &response.TransactionHistoryResponse{
		Transactions: txResponses,
		Total:        total,
		Page:         offset/limit + 1,
		Limit:        limit,
		TotalPages:   totalPages,
	}, nil
}

func (s *BillPaymentService) mapTransactionToResponse(tx *entities.Transaction) *response.TransactionResponse {
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

	if tx.BillDetail != nil {
		resp.BillDetail = &response.BillDetailResponse{
			BillType:         tx.BillDetail.BillType,
			ProviderName:     tx.BillDetail.ProviderName,
			PhoneNumber:      tx.BillDetail.PhoneNumber,
			MeterNumber:      tx.BillDetail.MeterNumber,
			CustomerName:     tx.BillDetail.CustomerName,
			ElectricityToken: tx.BillDetail.ElectricityToken,
			DataPlanName:     tx.BillDetail.DataPlanName,
			DataVolume:       tx.BillDetail.DataVolume,
		}
	}

	return resp
}
