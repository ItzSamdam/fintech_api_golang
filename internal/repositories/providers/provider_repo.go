package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"fintech_api_golang/internal/core/entities"
	"fintech_api_golang/internal/core/interfaces"
)

type providerRepository struct {
    db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) interfaces.ProviderRepository {
    return &providerRepository{db: db}
}

func (r *providerRepository) Create(ctx context.Context, provider *entities.Provider) error {
    return r.db.WithContext(ctx).Create(provider).Error
}

func (r *providerRepository) Update(ctx context.Context, provider *entities.Provider) error {
    return r.db.WithContext(ctx).Save(provider).Error
}

func (r *providerRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Provider, error) {
    var provider entities.Provider
    err := r.db.WithContext(ctx).First(&provider, "id = ?", id).Error
    return &provider, err
}

func (r *providerRepository) GetByCode(ctx context.Context, code string) (*entities.Provider, error) {
    var provider entities.Provider
    err := r.db.WithContext(ctx).First(&provider, "code = ?", code).Error
    return &provider, err
}

func (r *providerRepository) GetByType(ctx context.Context, providerType string) ([]entities.Provider, error) {
    var providers []entities.Provider
    err := r.db.WithContext(ctx).
        Where("type = ?", providerType).
        Find(&providers).Error
    return providers, err
}

func (r *providerRepository) GetActiveByType(ctx context.Context, providerType string) ([]entities.Provider, error) {
    var providers []entities.Provider
    err := r.db.WithContext(ctx).
        Where("type = ? AND is_active = ?", providerType, true).
        Order("priority ASC").
        Find(&providers).Error
    return providers, err
}

func (r *providerRepository) GetByPriority(ctx context.Context, providerType string) ([]entities.Provider, error) {
    var providers []entities.Provider
    err := r.db.WithContext(ctx).
        Where("type = ? AND is_active = ?", providerType, true).
        Order("priority ASC").
        Find(&providers).Error
    return providers, err
}

func (r *providerRepository) ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error {
    return r.db.WithContext(ctx).Model(&entities.Provider{}).
        Where("id = ?", id).
        Update("is_active", isActive).Error
}

func (r *providerRepository) UpdatePriority(ctx context.Context, id uuid.UUID, priority int) error {
    return r.db.WithContext(ctx).Model(&entities.Provider{}).
        Where("id = ?", id).
        Update("priority", priority).Error
}

func (r *providerRepository) UpdateHealthStatus(ctx context.Context, id uuid.UUID, status string, lastCheck time.Time) error {
    return r.db.WithContext(ctx).Model(&entities.Provider{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "health_status":     status,
            "last_health_check": lastCheck,
        }).Error
}

func (r *providerRepository) UpdateMargin(ctx context.Context, id uuid.UUID, marginPercent float64) error {
    return r.db.WithContext(ctx).Model(&entities.Provider{}).
        Where("id = ?", id).
        Update("margin_percent", marginPercent).Error
}

func (r *providerRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]entities.Provider, int64, error) {
    var providers []entities.Provider
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entities.Provider{})
    
    if providerType, ok := filters["type"]; ok {
        query = query.Where("type = ?", providerType)
    }
    if isActive, ok := filters["is_active"]; ok {
        query = query.Where("is_active = ?", isActive)
    }
    
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := query.Offset(offset).Limit(limit).Order("priority ASC").Find(&providers).Error
    return providers, total, err
}

type providerLogRepository struct {
    db *gorm.DB
}

func NewProviderLogRepository(db *gorm.DB) interfaces.ProviderLogRepository {
    return &providerLogRepository{db: db}
}

func (r *providerLogRepository) Create(ctx context.Context, log *entities.ProviderLog) error {
    return r.db.WithContext(ctx).Create(log).Error
}

func (r *providerLogRepository) GetByProviderID(ctx context.Context, providerID uuid.UUID, offset, limit int) ([]entities.ProviderLog, int64, error) {
    var logs []entities.ProviderLog
    var total int64
    
    query := r.db.WithContext(ctx).Model(&entities.ProviderLog{}).
        Where("provider_id = ?", providerID)
    
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
    return logs, total, err
}

func (r *providerLogRepository) GetByTransactionID(ctx context.Context, transactionID uuid.UUID) (*entities.ProviderLog, error) {
    var log entities.ProviderLog
    err := r.db.WithContext(ctx).First(&log, "transaction_id = ?", transactionID).Error
    return &log, err
}

func (r *providerLogRepository) GetErrorLogs(ctx context.Context, startDate, endDate time.Time) ([]entities.ProviderLog, error) {
    var logs []entities.ProviderLog
    err := r.db.WithContext(ctx).
        Where("is_error = ? AND created_at BETWEEN ? AND ?", true, startDate, endDate).
        Order("created_at DESC").
        Find(&logs).Error
    return logs, err
}

func (r *providerLogRepository) GetProviderStats(ctx context.Context, providerID uuid.UUID, startDate, endDate time.Time) (*interfaces.ProviderStats, error) {
    var stats interfaces.ProviderStats
    
    err := r.db.WithContext(ctx).Model(&entities.ProviderLog{}).
        Where("provider_id = ? AND created_at BETWEEN ? AND ?", providerID, startDate, endDate).
        Select("COUNT(*) as total_requests, " +
            "SUM(CASE WHEN is_error = false THEN 1 ELSE 0 END) as success_count, " +
            "SUM(CASE WHEN is_error = true THEN 1 ELSE 0 END) as failed_count, " +
            "AVG(response_time) as avg_response_time").
        Scan(&stats).Error
    
    if stats.TotalRequests > 0 {
        stats.SuccessRate = float64(stats.SuccessCount) / float64(stats.TotalRequests) * 100
    }
    
    return &stats, err
}

// ========== REDBILLER PROVIDER CLIENT ==========

type RedBillerClient struct {
    baseURL    string
    privateKey string
    httpClient *http.Client
}

type RedBillerResponse struct {
    Success bool                   `json:"success"`
    Message string                 `json:"message"`
    Data    map[string]interface{} `json:"data"`
}

func NewRedBillerClient(baseURL, privateKey string) *RedBillerClient {
    return &RedBillerClient{
        baseURL:    baseURL,
        privateKey: privateKey,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (c *RedBillerClient) doRequest(ctx context.Context, method, endpoint string, body interface{}) (*RedBillerResponse, error) {
    url := c.baseURL + endpoint
    
    var reqBody io.Reader
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return nil, err
        }
        reqBody = bytes.NewBuffer(jsonBody)
    }
    
    req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Private-Key", c.privateKey)
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return &RedBillerResponse{
            Success: false,
            Message: err.Error(),
        }, nil
    }
    defer resp.Body.Close()
    
    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    // Check if response indicates success
    success := true
    if status, ok := result["status"]; ok {
        if status == "error" || status == "failed" {
            success = false
        }
    }
    
    message := ""
    if msg, ok := result["message"]; ok {
        message = fmt.Sprintf("%v", msg)
    }
    
    return &RedBillerResponse{
        Success: success,
        Message: message,
        Data:    result,
    }, nil
}

// GetWalletBalance - GET /1.0/get/balance
func (c *RedBillerClient) GetWalletBalance(ctx context.Context) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "GET", "/1.0/get/balance", nil)
}

// FetchBanks - GET /1.0/payout/bank-transfer/banks/list
func (c *RedBillerClient) FetchBanks(ctx context.Context, bankType, countryCode, currencyCode string) (*RedBillerResponse, error) {
    endpoint := fmt.Sprintf("/1.0/payout/bank-transfer/banks/list?type=%s&country_code=%s&currency_code=%s",
        bankType, countryCode, currencyCode)
    return c.doRequest(ctx, "GET", endpoint, nil)
}

// SendMoney - POST /2.0/payout/bank-transfer/create
type SendMoneyRequest struct {
    AccountNo   string `json:"account_no"`
    BankCode    string `json:"bank_code"`
    Amount      int64  `json:"amount"`
    Narration   string `json:"narration"`
    CallbackURL string `json:"callback_url"`
    Reference   string `json:"reference"`
}

func (c *RedBillerClient) SendMoney(ctx context.Context, req *SendMoneyRequest) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/2.0/payout/bank-transfer/create", req)
}

// RetrySendMoney - POST /2.0/payout/bank-transfer/retry
type RetryRequest struct {
    Reference string `json:"reference"`
}

func (c *RedBillerClient) RetrySendMoney(ctx context.Context, reference string) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/2.0/payout/bank-transfer/retry", &RetryRequest{Reference: reference})
}

// SuggestBank - POST /1.0/payout/bank-transfer/banks/suggest
type SuggestBankRequest struct {
    AccountNo string `json:"account_no"`
}

func (c *RedBillerClient) SuggestBank(ctx context.Context, accountNo string) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/1.0/payout/bank-transfer/banks/suggest", &SuggestBankRequest{AccountNo: accountNo})
}

// VerifyAccountDetails - POST /1.0/kyc/bank-account/verify
type VerifyAccountRequest struct {
    AccountNo string `json:"account_no"`
    BankCode  string `json:"bank_code"`
}

func (c *RedBillerClient) VerifyAccountDetails(ctx context.Context, accountNo, bankCode string) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/1.0/kyc/bank-account/verify", &VerifyAccountRequest{
        AccountNo: accountNo,
        BankCode:  bankCode,
    })
}

// VerifyTransaction - POST /1.0/payout/bank-transfer/status
func (c *RedBillerClient) VerifyTransaction(ctx context.Context, reference string) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/1.0/payout/bank-transfer/status", &RetryRequest{Reference: reference})
}

// CreateVirtualAccount - POST /1.0/collections/PSA/create
type VirtualAccountRequest struct {
    Bank           string `json:"bank"`
    FirstName      string `json:"first_name"`
    Surname        string `json:"surname"`
    PhoneNo        string `json:"phone_no"`
    Email          string `json:"email"`
    BVN            string `json:"bvn"`
    DateOfBirth    string `json:"date_of_birth"`
    AutoSettlement bool   `json:"auto_settlement"`
    CallbackURL    string `json:"callback_url"`
    Reference      string `json:"reference"`
}

func (c *RedBillerClient) CreateVirtualAccount(ctx context.Context, req *VirtualAccountRequest) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/1.0/collections/PSA/create", req)
}

// VerifyVirtualAccountPayment - POST /1.0/collections/PSA/payments/verify
func (c *RedBillerClient) VerifyVirtualAccountPayment(ctx context.Context, reference string) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/1.0/collections/PSA/payments/verify", &RetryRequest{Reference: reference})
}

// VerifyDisco - POST /1.0/bills/disco/meter/verify
type VerifyDiscoRequest struct {
    Product    string `json:"product"`
    MeterNo    string `json:"meter_no"`
    MeterType  string `json:"meter_type"`
}

func (c *RedBillerClient) VerifyDisco(ctx context.Context, product, meterNo, meterType string) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/1.0/bills/disco/meter/verify", &VerifyDiscoRequest{
        Product:   product,
        MeterNo:   meterNo,
        MeterType: meterType,
    })
}

// PurchaseDisco - POST /1.1/bills/disco/purchase/create
type PurchaseDiscoRequest struct {
    Product      string `json:"product"`
    MeterNo      string `json:"meter_no"`
    CustomerName string `json:"customer_name"`
    MeterType    string `json:"meter_type"`
    PhoneNo      string `json:"phone_no"`
    Amount       int64  `json:"amount"`
    CallbackURL  string `json:"callback_url"`
    Reference    string `json:"reference"`
}

func (c *RedBillerClient) PurchaseDisco(ctx context.Context, req *PurchaseDiscoRequest) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/1.1/bills/disco/purchase/create", req)
}

// GetBetProviders - GET /1.5/bills/betting/providers/list
func (c *RedBillerClient) GetBetProviders(ctx context.Context) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "GET", "/1.5/bills/betting/providers/list", nil)
}

// CreditBetWallet - POST /1.5/bills/betting/account/payment/create
type CreditBetWalletRequest struct {
    Product     string `json:"product"`
    CustomerID  string `json:"customer_id"`
    Amount      int64  `json:"amount"`
    PhoneNo     string `json:"phone_no"`
    CallbackURL string `json:"callback_url"`
    Reference   string `json:"reference"`
}

func (c *RedBillerClient) CreditBetWallet(ctx context.Context, req *CreditBetWalletRequest) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/1.5/bills/betting/account/payment/create", req)
}

// VerifyBetWallet - POST /1.5/bills/betting/account/verify
type VerifyBetWalletRequest struct {
    Product    string `json:"product"`
    CustomerID string `json:"customer_id"`
}

func (c *RedBillerClient) VerifyBetWallet(ctx context.Context, product, customerID string) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/1.5/bills/betting/account/verify", &VerifyBetWalletRequest{
        Product:    product,
        CustomerID: customerID,
    })
}

// GetCablePlans - POST /1.0/bills/cable/plans/list
type GetCablePlansRequest struct {
    Product string `json:"product"`
}

func (c *RedBillerClient) GetCablePlans(ctx context.Context, product string) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/1.0/bills/cable/plans/list", &GetCablePlansRequest{Product: product})
}

// VerifyCable - POST /1.0/bills/cable/decoder/verify
type VerifyCableRequest struct {
    Product     string `json:"product"`
    SmartCardNo string `json:"smart_card_no"`
}

func (c *RedBillerClient) VerifyCable(ctx context.Context, product, smartCardNo string) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/1.0/bills/cable/decoder/verify", &VerifyCableRequest{
        Product:     product,
        SmartCardNo: smartCardNo,
    })
}

// PurchaseCable - POST /1.1/bills/cable/plans/purchase/create
type PurchaseCableRequest struct {
    Product      string `json:"product"`
    Code         string `json:"code"`
    SmartCardNo  string `json:"smart_card_no"`
    CustomerName string `json:"customer_name"`
    PhoneNo      string `json:"phone_no"`
    CallbackURL  string `json:"callback_url"`
    Reference    string `json:"reference"`
}

func (c *RedBillerClient) PurchaseCable(ctx context.Context, req *PurchaseCableRequest) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/1.1/bills/cable/plans/purchase/create", req)
}

// PurchaseTopUp - POST /1.0/bills/airtime/purchase/create
type PurchaseTopUpRequest struct {
    Product     string `json:"product"`
    PhoneNo     string `json:"phone_no"`
    Amount      int64  `json:"amount"`
    Ported      bool   `json:"ported"`
    CallbackURL string `json:"callback_url"`
    Reference   string `json:"reference"`
}

func (c *RedBillerClient) PurchaseTopUp(ctx context.Context, req *PurchaseTopUpRequest) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/1.0/bills/airtime/purchase/create", req)
}

// GetDataPlans - POST /1.0/bills/data/plans/list
func (c *RedBillerClient) GetDataPlans(ctx context.Context, product string) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/1.0/bills/data/plans/list", &GetCablePlansRequest{Product: product})
}

// PurchaseData - POST /1.0/bills/data/plans/purchase/create
type PurchaseDataRequest struct {
    Product     string `json:"product"`
    PhoneNo     string `json:"phone_no"`
    Code        string `json:"code"`
    Ported      bool   `json:"ported"`
    CallbackURL string `json:"callback_url"`
    Reference   string `json:"reference"`
}

func (c *RedBillerClient) PurchaseData(ctx context.Context, req *PurchaseDataRequest) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/1.0/bills/data/plans/purchase/create", req)
}

// VerifyAnyTransaction - POST to various endpoints based on transaction type
func (c *RedBillerClient) VerifyAnyTransaction(ctx context.Context, reference, txnType string) (*RedBillerResponse, error) {
    endpoints := map[string]string{
        "betting":   "/1.4/bills/betting/account/payment/status",
        "disco":     "/1.0/bills/disco/purchase/status",
        "cable":     "/1.0/bills/cable/plans/purchase/status",
        "data":      "/1.0/bills/data/plans/purchase/status",
        "airtime":   "/1.0/bills/airtime/purchase/status",
        "transfer":  "/1.0/payout/bank-transfer/status",
    }
    
    endpoint, ok := endpoints[txnType]
    if !ok {
        return &RedBillerResponse{
            Success: false,
            Message: fmt.Sprintf("Unsupported transaction type: %s", txnType),
        }, nil
    }
    
    return c.doRequest(ctx, "POST", endpoint, &RetryRequest{Reference: reference})
}

// VerifyBVN - POST /1.0/kyc/bvn/verify.3.0
type VerifyBVNRequest struct {
    BVN       string `json:"bvn"`
    Reference string `json:"reference"`
}

func (c *RedBillerClient) VerifyBVN(ctx context.Context, bvn, reference string) (*RedBillerResponse, error) {
    return c.doRequest(ctx, "POST", "/1.0/kyc/bvn/verify.3.0", &VerifyBVNRequest{
        BVN:       bvn,
        Reference: reference,
    })
}

// Provider Registry with RedBiller
type ProviderRegistry struct {
    providers   map[string]interface{}
    redBiller   *RedBillerClient
}

func NewProviderRegistry(redBiller *RedBillerClient) *ProviderRegistry {
    return &ProviderRegistry{
        providers: make(map[string]interface{}),
        redBiller: redBiller,
    }
}

func (r *ProviderRegistry) GetRedBiller() *RedBillerClient {
    return r.redBiller
}

func (r *ProviderRegistry) Register(name string, client interface{}) {
    r.providers[name] = client
}

func (r *ProviderRegistry) Get(name string) (interface{}, error) {
    client, ok := r.providers[name]
    if !ok {
        return nil, fmt.Errorf("provider not found: %s", name)
    }
    return client, nil
}