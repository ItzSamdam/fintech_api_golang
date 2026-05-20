package interfaces

import (
    "context"
    "time"
    "github.com/google/uuid"
)

type NotificationRepository interface {
    CreateInApp(ctx context.Context, userID uuid.UUID, title, message, notificationType string, metadata map[string]interface{}) error
    GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int, unreadOnly bool) ([]InAppNotification, int64, error)
    MarkAsRead(ctx context.Context, notificationID uuid.UUID) error
    MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
    Delete(ctx context.Context, notificationID uuid.UUID) error
    GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error)
}

type InAppNotification struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    Title     string
    Message   string
    Type      string
    IsRead    bool
    Metadata  string
    CreatedAt time.Time
}