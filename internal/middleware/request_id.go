package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
)

// RequestID adds a unique request ID to each request
func RequestID() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Check if request ID already exists in header
        requestID := c.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        // Set request ID in response header
        c.Set("X-Request-ID", requestID)
        
        // Store in context locals for use in handlers
        c.Locals("request_id", requestID)
        
        return c.Next()
    }
}

// GetRequestID retrieves request ID from context
func GetRequestID(c *fiber.Ctx) string {
    requestID := c.Locals("request_id")
    if requestID == nil {
        return ""
    }
    return requestID.(string)
}