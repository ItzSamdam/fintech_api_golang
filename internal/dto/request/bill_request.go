package request

// Airtime Requests
type PurchaseAirtimeRequest struct {
    PhoneNumber string `json:"phone_number" validate:"required,len=11,numeric"`
    Network     string `json:"network" validate:"required,oneof=MTN GLO AIRTEL 9MOBILE"`
    Amount      int64  `json:"amount" validate:"required,min=5000"` // In kobo (min 50 NGN)
    SaveBeneficiary bool `json:"save_beneficiary"`
}

// Data Requests
type PurchaseDataRequest struct {
    PhoneNumber string `json:"phone_number" validate:"required,len=11,numeric"`
    Network     string `json:"network" validate:"required,oneof=MTN GLO AIRTEL 9MOBILE"`
    PlanID      string `json:"plan_id" validate:"required"`
}

// Electricity Requests
type ValidateMeterRequest struct {
    MeterNumber string `json:"meter_number" validate:"required"`
    ProviderID  string `json:"provider_id" validate:"required"`
    MeterType   string `json:"meter_type" validate:"required,oneof=prepaid postpaid"`
}

type PayElectricityRequest struct {
    MeterNumber string `json:"meter_number" validate:"required"`
    ProviderID  string `json:"provider_id" validate:"required"`
    MeterType   string `json:"meter_type" validate:"required,oneof=prepaid postpaid"`
    Amount      int64  `json:"amount" validate:"required,min=10000"` // In kobo (min 100 NGN)
    CustomerName string `json:"customer_name"` // From validation response
    Address     string `json:"address"`
    Email       string `json:"email" validate:"omitempty,email"`
}

// Betting Requests
type ValidateBettingAccountRequest struct {
    ProviderID string `json:"provider_id" validate:"required"`
    AccountID  string `json:"account_id" validate:"required"`
}

type FundBettingRequest struct {
    ProviderID string `json:"provider_id" validate:"required"`
    AccountID  string `json:"account_id" validate:"required"`
    Amount     int64  `json:"amount" validate:"required,min=10000"` // In kobo
}

type WithdrawBettingRequest struct {
    ProviderID string `json:"provider_id" validate:"required"`
    AccountID  string `json:"account_id" validate:"required"`
    Amount     int64  `json:"amount" validate:"required,min=10000"`
}