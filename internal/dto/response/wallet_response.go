package response

import (
    "time"
    "github.com/google/uuid"
)

type WalletResponse struct {
    ID            uuid.UUID `json:"id"`
    UserID        uuid.UUID `json:"user_id"`
    Balance       int64     `json:"balance"`        // In kobo
    BalanceNaira  float64   `json:"balance_naira"`  // In naira
    Currency      string    `json:"currency"`
    IsLocked      bool      `json:"is_locked"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type TransactionResponse struct {
    ID              uuid.UUID `json:"id"`
    Reference       string    `json:"reference"`
    Type            string    `json:"type"`
    Category        string    `json:"category"`
    SubCategory     string    `json:"sub_category,omitempty"`
    Amount          int64     `json:"amount"`         // In kobo
    AmountNaira     float64   `json:"amount_naira"`   // In naira
    Fee             int64     `json:"fee"`            // In kobo
    FeeNaira        float64   `json:"fee_naira"`      // In naira
    TotalAmount     int64     `json:"total_amount"`   // In kobo
    TotalAmountNaira float64  `json:"total_amount_naira"` // In naira
    Status          string    `json:"status"`
    Description     string    `json:"description"`
    BalanceBefore   int64     `json:"balance_before"`
    BalanceAfter    int64     `json:"balance_after"`
    CompletedAt     *time.Time `json:"completed_at,omitempty"`
    CreatedAt       time.Time `json:"created_at"`
    
    // Additional details based on category
    TransferDetail  *TransferDetailResponse  `json:"transfer_detail,omitempty"`
    BillDetail      *BillDetailResponse      `json:"bill_detail,omitempty"`
}

type TransferDetailResponse struct {
    RecipientType   string `json:"recipient_type"`
    RecipientID     string `json:"recipient_id"`
    RecipientName   string `json:"recipient_name"`
    RecipientBank   string `json:"recipient_bank,omitempty"`
    Narration       string `json:"narration"`
}

type BillDetailResponse struct {
    BillType        string `json:"bill_type"`
    ProviderName    string `json:"provider_name"`
    PhoneNumber     string `json:"phone_number,omitempty"`
    MeterNumber     string `json:"meter_number,omitempty"`
    CustomerName    string `json:"customer_name"`
    ElectricityToken string `json:"electricity_token,omitempty"`
    DataPlanName    string `json:"data_plan_name,omitempty"`
    DataVolume      string `json:"data_volume,omitempty"`
}

type TransactionHistoryResponse struct {
    Transactions []TransactionResponse `json:"transactions"`
    Total        int64                 `json:"total"`
    Page         int                   `json:"page"`
    Limit        int                   `json:"limit"`
    TotalPages   int                   `json:"total_pages"`
}