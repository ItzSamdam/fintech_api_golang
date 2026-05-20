package response

type NetworkListResponse struct {
    Networks []Network `json:"networks"`
}

type Network struct {
    Code        string `json:"code"`
    Name        string `json:"name"`
    LogoURL     string `json:"logo_url"`
    IsActive    bool   `json:"is_active"`
}

type DataPlanResponse struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Volume      string  `json:"volume"`      // e.g., "1GB"
    Price       int64   `json:"price"`       // In kobo
    PriceNaira  float64 `json:"price_naira"` // In naira
    Validity    string  `json:"validity"`    // e.g., "30 days"
    Network     string  `json:"network"`
}

type AirtimeDenominationResponse struct {
    Amount      int64   `json:"amount"`       // In kobo
    AmountNaira float64 `json:"amount_naira"` // In naira
}

type ProviderResponse struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Code        string `json:"code"`
    Category    string `json:"category"`
    IsActive    bool   `json:"is_active"`
}

type MeterValidationResponse struct {
    CustomerName    string `json:"customer_name"`
    CustomerAddress string `json:"customer_address"`
    MeterNumber     string `json:"meter_number"`
    MeterType       string `json:"meter_type"`
    ProviderID      string `json:"provider_id"`
    ProviderName    string `json:"provider_name"`
}

type ElectricityPaymentResponse struct {
    TransactionID   string `json:"transaction_id"`
    Reference       string `json:"reference"`
    MeterNumber     string `json:"meter_number"`
    CustomerName    string `json:"customer_name"`
    Amount          int64  `json:"amount"`        // In kobo
    AmountNaira     float64 `json:"amount_naira"` // In naira
    Token           string `json:"token,omitempty"` // For prepaid
    Units           int    `json:"units,omitempty"`
    Status          string `json:"status"`
}

type BettingAccountResponse struct {
    ProviderID      string `json:"provider_id"`
    ProviderName    string `json:"provider_name"`
    AccountID       string `json:"account_id"`
    AccountName     string `json:"account_name"`
    IsValid         bool   `json:"is_valid"`
    Balance         int64  `json:"balance,omitempty"` // In kobo
}

type BillHistoryResponse struct {
    ID              string    `json:"id"`
    BillType        string    `json:"bill_type"`
    Provider        string    `json:"provider"`
    Amount          int64     `json:"amount"` // In kobo
    AmountNaira     float64   `json:"amount_naira"`
    Status          string    `json:"status"`
    Recipient       string    `json:"recipient"` // Phone number, meter number, or account ID
    CreatedAt       string    `json:"created_at"`
}