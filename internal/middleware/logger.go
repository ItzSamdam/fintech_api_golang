package middleware

import (
    "time"
    
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "go.uber.org/zap"
)

type LoggerMiddleware struct {
    logger *zap.Logger
}

func NewLoggerMiddleware(logger *zap.Logger) *LoggerMiddleware {
    return &LoggerMiddleware{
        logger: logger,
    }
}

// LogRequests logs all HTTP requests
func (l *LoggerMiddleware) LogRequests() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        requestID := c.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
            c.Set("X-Request-ID", requestID)
        }
        
        // Log request
        l.logger.Info("Incoming request",
            zap.String("request_id", requestID),
            zap.String("method", c.Method()),
            zap.String("path", c.Path()),
            zap.String("ip", c.IP()),
            zap.String("user_agent", c.Get("User-Agent")),
        )
        
        // Process request
        err := c.Next()
        
        // Log response
        duration := time.Since(start)
        status := c.Response().StatusCode()
        
        logFields := []zap.Field{
            zap.String("request_id", requestID),
            zap.String("method", c.Method()),
            zap.String("path", c.Path()),
            zap.Int("status", status),
            zap.Duration("duration", duration),
            zap.String("ip", c.IP()),
        }
        
        if err != nil {
            logFields = append(logFields, zap.Error(err))
            l.logger.Error("Request completed with error", logFields...)
        } else if status >= 500 {
            l.logger.Error("Request completed with server error", logFields...)
        } else if status >= 400 {
            l.logger.Warn("Request completed with client error", logFields...)
        } else {
            l.logger.Info("Request completed successfully", logFields...)
        }
        
        return err
    }
}