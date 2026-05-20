package middleware

import (
    "runtime/debug"
    
    "github.com/gofiber/fiber/v2"
    "go.uber.org/zap"
)

type RecoveryMiddleware struct {
    logger *zap.Logger
}

func NewRecoveryMiddleware(logger *zap.Logger) *RecoveryMiddleware {
    return &RecoveryMiddleware{
        logger: logger,
    }
}

// Recover panics and returns 500 error
func (r *RecoveryMiddleware) Recover() fiber.Handler {
    return func(c *fiber.Ctx) error {
        defer func() {
            if err := recover(); err != nil {
                // Log the panic with stack trace
                r.logger.Error("Panic recovered",
                    zap.Any("error", err),
                    zap.String("stack", string(debug.Stack())),
                    zap.String("path", c.Path()),
                    zap.String("method", c.Method()),
                )
                
                // Return 500 error to client
                c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                    "error":   "Internal server error",
                    "message": "An unexpected error occurred",
                })
            }
        }()
        
        return c.Next()
    }
}