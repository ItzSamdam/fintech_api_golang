package middleware

import (
    "context"
    "time"
    
    "github.com/gofiber/fiber/v2"
)

// Timeout adds request timeout
func Timeout(timeout time.Duration) fiber.Handler {
    return func(c *fiber.Ctx) error {
        ctx, cancel := context.WithTimeout(c.Context(), timeout)
        defer cancel()
        
        // Replace context with timeout context
        c.SetUserContext(ctx)
        
        // Create channel to handle request completion
        done := make(chan error, 1)
        
        go func() {
            done <- c.Next()
        }()
        
        select {
        case err := <-done:
            return err
        case <-ctx.Done():
            return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
                "error": "Request timeout",
            })
        }
    }
}