package request

type ReportSuspiciousRequest struct {
    TransactionID string `json:"transaction_id" validate:"required"`
    Reason        string `json:"reason" validate:"required"`
    Description   string `json:"description"`
}

type CheckLimitsRequest struct {
    TransactionType string `json:"transaction_type" validate:"required"`
    Amount          int64  `json:"amount" validate:"required"`
}

type VerifySIMSwapRequest struct {
    PhoneNumber string `json:"phone_number" validate:"required,len=11,numeric"`
}

type TrustDeviceRequest struct {
    DeviceID   string `json:"device_id" validate:"required"`
    DeviceName string `json:"device_name" validate:"required"`
    DeviceType string `json:"device_type" validate:"required,oneof=ios android web"`
}

type Enable2FARequest struct {
    Password string `json:"password" validate:"required"`
}

type Verify2FARequest struct {
    Code string `json:"code" validate:"required,len=6"`
}