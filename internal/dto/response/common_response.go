package response

type SuccessResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type BankListResponse struct {
    Banks []Bank `json:"banks"`
}

type Bank struct {
    Code string `json:"code"`
    Name string `json:"name"`
}

type NameEnquiryResponse struct {
    AccountNumber string `json:"account_number"`
    AccountName   string `json:"account_name"`
    BankCode      string `json:"bank_code"`
    BankName      string `json:"bank_name"`
}

type OTPResponse struct {
    Reference string `json:"reference"`
    ExpiresIn int    `json:"expires_in"` // seconds
}

type BalanceResponse struct {
    Balance      int64   `json:"balance"`       // In kobo
    BalanceNaira float64 `json:"balance_naira"` // In naira
    Currency     string  `json:"currency"`
    IsLocked     bool    `json:"is_locked"`
}