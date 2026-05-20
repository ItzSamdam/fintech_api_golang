package request

type CreateWalletRequest struct {
    Currency string `json:"currency" validate:"omitempty,len=3"` // Default NGN
}

type WalletActionRequest struct {
    WalletID string `json:"wallet_id" validate:"required"`
    Reason   string `json:"reason"`
}

type GetTransactionsRequest struct {
    Page     int    `query:"page" validate:"omitempty,min=1"`
    Limit    int    `query:"limit" validate:"omitempty,min=1,max=100"`
    Category string `query:"category"`
    Status   string `query:"status"`
    FromDate string `query:"from_date"`
    ToDate   string `query:"to_date"`
}

type GetStatementRequest struct {
    FromDate string `query:"from_date" validate:"required"`
    ToDate   string `query:"to_date" validate:"required"`
    Format   string `query:"format" validate:"omitempty,oneof=pdf csv"`
}