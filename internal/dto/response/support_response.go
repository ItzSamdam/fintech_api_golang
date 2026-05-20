package response

import (
    "time"
    "github.com/google/uuid"
)

type SupportTicketResponse struct {
    ID              uuid.UUID            `json:"id"`
    UserID          uuid.UUID            `json:"user_id"`
    TransactionID   *uuid.UUID           `json:"transaction_id,omitempty"`
    Category        string               `json:"category"`
    Priority        string               `json:"priority"`
    Subject         string               `json:"subject"`
    Description     string               `json:"description"`
    Status          string               `json:"status"`
    AssignedTo      *uuid.UUID           `json:"assigned_to,omitempty"`
    AssignedToName  string               `json:"assigned_to_name,omitempty"`
    Messages        []TicketMessageResponse `json:"messages"`
    CreatedAt       time.Time            `json:"created_at"`
    UpdatedAt       time.Time            `json:"updated_at"`
    ResolvedAt      *time.Time           `json:"resolved_at,omitempty"`
}

type TicketMessageResponse struct {
    ID          uuid.UUID `json:"id"`
    SenderType  string    `json:"sender_type"`
    SenderID    uuid.UUID `json:"sender_id"`
    SenderName  string    `json:"sender_name"`
    Message     string    `json:"message"`
    AttachmentURL string  `json:"attachment_url,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}

type CreateTicketResponse struct {
    TicketID    uuid.UUID `json:"ticket_id"`
    Reference   string    `json:"reference"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
}