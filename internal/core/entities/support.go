package entities

import (
    "time"
    "github.com/google/uuid"
)

type SupportTicket struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID          uuid.UUID      `gorm:"type:uuid;index:idx_ticket_user;not null"`
    TransactionID   *uuid.UUID     `gorm:"type:uuid;index:idx_ticket_txn"`
    Category        string         `gorm:"size:50;not null"`          // transaction, account, billing, technical, other
    Priority        string         `gorm:"size:20;default:'medium'"`   // low, medium, high, urgent
    Subject         string         `gorm:"size:255;not null"`
    Description     string         `gorm:"type:text;not null"`
    Status          string         `gorm:"size:20;default:'open';index:idx_ticket_status"` // open, in_progress, resolved, closed
    AssignedTo      *uuid.UUID     `gorm:"type:uuid"`
    ResolvedAt      *time.Time
    ClosedAt        *time.Time
    Rating          int            // 1-5
    Feedback        string         `gorm:"type:text"`
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    UpdatedAt       time.Time      `gorm:"not null;default:now()"`
    
    // Relationships
    User            User           `gorm:"foreignKey:UserID"`
    Messages        []TicketMessage `gorm:"foreignKey:TicketID"`
    AssignedAdmin   *AdminUser     `gorm:"foreignKey:AssignedTo"`
}

type TicketMessage struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    TicketID        uuid.UUID      `gorm:"type:uuid;index:idx_message_ticket;not null"`
    SenderType      string         `gorm:"size:20;not null"`          // user, admin
    SenderID        uuid.UUID      `gorm:"type:uuid;not null"`
    Message         string         `gorm:"type:text;not null"`
    AttachmentURL   string         `gorm:"size:500"`
    IsRead          bool           `gorm:"default:false"`
    ReadAt          *time.Time
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    
    // Relationships
    Ticket          SupportTicket  `gorm:"foreignKey:TicketID"`
}