package request

// type SendTransferRequest struct {
//     RecipientType   string `json:"recipient_type" validate:"required,oneof=bank wallet"`
//     RecipientID     string `json:"recipient_id" validate:"required"`
//     Amount          int64  `json:"amount" validate:"required,min=100"` // In kobo
//     Narration       string `json:"narration" validate:"max=255"`
//     RecipientBankCode string `json:"recipient_bank_code" validate:"required_if=RecipientType bank"`
//     SaveBeneficiary bool   `json:"save_beneficiary"`
// }

type NameEnquiryRequest struct {
    AccountNumber string `json:"account_number" validate:"required,len=10,numeric"`
    BankCode      string `json:"bank_code" validate:"required,len=3,numeric"`
}

type RetryTransferRequest struct {
    Reference string `json:"reference" validate:"required"`
}

type TransferStatusRequest struct {
    Reference string `json:"reference" validate:"required"`
}

// In internal/dto/request/transfer_request.go, add validation tag:
type SendTransferRequest struct {
    RecipientType   string `json:"recipient_type" validate:"required,oneof=bank wallet"`
    RecipientID     string `json:"recipient_id" validate:"required"`
    Amount          int64  `json:"amount" validate:"required,min=100"`
    Narration       string `json:"narration" validate:"max=255"`
    RecipientBankCode string `json:"recipient_bank_code" validate:"required_if=RecipientType bank"`
    SaveBeneficiary bool   `json:"save_beneficiary"`
}