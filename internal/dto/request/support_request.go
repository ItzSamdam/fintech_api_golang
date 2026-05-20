package request

type CreateTicketRequest struct {
    Category    string `json:"category" validate:"required,oneof=transaction account billing technical other"`
    Subject     string `json:"subject" validate:"required,max=255"`
    Description string `json:"description" validate:"required"`
    TransactionID string `json:"transaction_id" validate:"omitempty"`
    Priority    string `json:"priority" validate:"omitempty,oneof=low medium high urgent"`
}

type ReplyToTicketRequest struct {
    Message     string `json:"message" validate:"required"`
    AttachmentURL string `json:"attachment_url"`
}

type UpdateTicketStatusRequest struct {
    Status string `json:"status" validate:"required,oneof=open in_progress resolved closed"`
    Note   string `json:"note"`
}