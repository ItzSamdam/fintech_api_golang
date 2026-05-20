package middleware

import (
    "encoding/json"
    "time"
    
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "gorm.io/gorm"
    
    "fintech_api_golang/internal/core/entities"
)

type AuditMiddleware struct {
    db *gorm.DB
}

func NewAuditMiddleware(db *gorm.DB) *AuditMiddleware {
    return &AuditMiddleware{
        db: db,
    }
}

// LogAdminAction logs actions performed by admin users
func (a *AuditMiddleware) LogAdminAction(entityType string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Get admin user from context
        adminUser := c.Locals("admin_user")
        if adminUser == nil {
            return c.Next()
        }
        
        admin, ok := adminUser.(*entities.AdminUser)
        if !ok {
            return c.Next()
        }
        
        // Capture request body for old/new values
        var oldValue, newValue interface{}
        
        // For GET requests, don't capture body
        if c.Method() != "GET" {
            body := c.Body()
            if len(body) > 0 {
                json.Unmarshal(body, &newValue)
            }
        }
        
        // Process request
        err := c.Next()
        
        // Create audit log asynchronously (goroutine to not block response)
        go func() {
            auditLog := &entities.AuditLog{
                ID:         uuid.New(),
                AdminID:    &admin.ID,
                Action:     c.Method() + " " + c.Path(),
                EntityType: entityType,
                EntityID:   c.Params("id"),
				OldValue:   a.toJSON(oldValue),
                NewValue:   a.toJSON(newValue),
                IPAddress:  c.IP(),
                UserAgent:  c.Get("User-Agent"),
                CreatedAt:  time.Now(),
            }
            
            // Try to get old value if this is an update
            if c.Method() == "PUT" || c.Method() == "PATCH" {
                // Fetch old entity based on ID (implement as needed)
                // auditLog.OldValue = a.getOldEntity(entityType, c.Params("id"))
            }
            
            a.db.Create(auditLog)
        }()
        
        return err
    }
}

// LogUserAction logs important user actions (login, password change, etc.)
func (a *AuditMiddleware) LogUserAction(action string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID, err := GetUserIDFromContext(c)
        if err != nil {
            return c.Next()
        }
        
        err = c.Next()
        
        go func() {
            auditLog := &entities.AuditLog{
                ID:         uuid.New(),
                UserID:     &userID,
                Action:     action,
                EntityType: "user",
                IPAddress:  c.IP(),
                UserAgent:  c.Get("User-Agent"),
                Metadata:   a.toJSON(map[string]string{"path": c.Path()}),
                CreatedAt:  time.Now(),
            }
            
            a.db.Create(auditLog)
        }()
        
        return err
    }
}

func (a *AuditMiddleware) toJSON(v interface{}) string {
    if v == nil {
        return ""
    }
    bytes, err := json.Marshal(v)
    if err != nil {
        return ""
    }
    return string(bytes)
}