package interfaces

import (
    "context"
    
    "github.com/google/uuid"
    "fintech_api_golang/internal/core/entities"
)

type SupportTicketRepository interface {
    Create(ctx context.Context, ticket *entities.SupportTicket) error
    Update(ctx context.Context, ticket *entities.SupportTicket) error
    GetByID(ctx context.Context, id uuid.UUID) (*entities.SupportTicket, error)
    GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]entities.SupportTicket, int64, error)
    GetByStatus(ctx context.Context, status string, offset, limit int) ([]entities.SupportTicket, int64, error)
    GetByPriority(ctx context.Context, priority string, offset, limit int) ([]entities.SupportTicket, int64, error)
    GetAssignedTo(ctx context.Context, adminID uuid.UUID, offset, limit int) ([]entities.SupportTicket, int64, error)
    UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
    AssignTo(ctx context.Context, id uuid.UUID, adminID uuid.UUID) error
    Resolve(ctx context.Context, id uuid.UUID, rating int, feedback string) error
    List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]entities.SupportTicket, int64, error)
    GetStats(ctx context.Context) (*TicketStats, error)
}

type TicketMessageRepository interface {
    Create(ctx context.Context, message *entities.TicketMessage) error
    GetByTicketID(ctx context.Context, ticketID uuid.UUID, offset, limit int) ([]entities.TicketMessage, int64, error)
    MarkAsRead(ctx context.Context, messageID uuid.UUID) error
    MarkTicketMessagesAsRead(ctx context.Context, ticketID uuid.UUID, senderType string) error
    GetUnreadCount(ctx context.Context, ticketID uuid.UUID, senderType string) (int64, error)
}

type TicketStats struct {
    OpenCount      int64
    InProgressCount int64
    ResolvedCount  int64
    ClosedCount    int64
    HighPriorityCount int64
    AvgResolutionTime float64
}