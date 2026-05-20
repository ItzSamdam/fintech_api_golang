package response

import (
    "time"
    "github.com/google/uuid"
)

type InAppNotificationResponse struct {
    ID        uuid.UUID `json:"id"`
    Title     string    `json:"title"`
    Message   string    `json:"message"`
    Type      string    `json:"type"` // info, success, warning, error
    IsRead    bool      `json:"is_read"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
    CreatedAt time.Time `json:"created_at"`
}