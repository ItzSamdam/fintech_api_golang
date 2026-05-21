package middleware

import (
    "strings"
    
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/helmet"
)

// SecurityHeaders adds security-related headers to responses
// func SecurityHeaders() fiber.Handler {
//     return func(c *fiber.Ctx) error {
//         // Set security headers
//         c.Set("X-Content-Type-Options", "nosniff")
//         c.Set("X-Frame-Options", "DENY")
//         c.Set("X-XSS-Protection", "1; mode=block")
//         c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
//         c.Set("Content-Security-Policy", "default-src 'self'")
//         c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
        
//         return c.Next()
//     }
// }
// 
// // SecurityHeaders adds security-related headers to responses
func SecurityHeaders() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Only relax CSP for Swagger UI endpoint
        path := c.Path()
        if path == "/swagger" || path == "/swagger/" || path == "/docs" {
            c.Set("Content-Security-Policy", "default-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline' https://unpkg.com; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://unpkg.com; img-src 'self' data: https:; font-src 'self' data: https:; connect-src 'self' https:;")
        } else {
            c.Set("Content-Security-Policy", "default-src 'self'")
        }
        
        c.Set("X-Content-Type-Options", "nosniff")
        c.Set("X-Frame-Options", "DENY")
        c.Set("X-XSS-Protection", "1; mode=block")
        c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
        
        return c.Next()
    }
}

// Helmet middleware for additional security
func Helmet() fiber.Handler {
    return helmet.New(helmet.Config{
        CrossOriginEmbedderPolicy: "require-corp",
        CrossOriginOpenerPolicy:   "same-origin",
        CrossOriginResourcePolicy: "same-origin",
        OriginAgentCluster:        "?1",
        ReferrerPolicy:            "strict-origin-when-cross-origin",
        // XContentTypeOptions:       "nosniff",
        XDNSPrefetchControl:       "off",
        XFrameOptions:             "DENY",
        // XPoweredBy:                "off",
        // XXSSProtection:            "1; mode=block",
    })
}

// ValidateIP whitelists specific IPs (for admin endpoints)
func ValidateIP(allowedIPs []string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        clientIP := c.IP()
        
        // Check if IP is allowed
        for _, allowedIP := range allowedIPs {
            if clientIP == allowedIP {
                return c.Next()
            }
        }
        
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Access denied from this IP address",
        })
    }
}

// SanitizeInput sanitizes request inputs to prevent XSS
func SanitizeInput() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Parse body
        var body map[string]interface{}
        if err := c.BodyParser(&body); err == nil {
            sanitized := sanitizeMap(body)
            c.Locals("sanitized_body", sanitized)
        }
        
        // Sanitize query params
        for key, value := range c.Queries() {
            c.Request().URI().QueryArgs().Set(key, sanitizeString(value))
        }
        
        return c.Next()
    }
}

func sanitizeString(s string) string {
    // Remove dangerous characters
    replacer := strings.NewReplacer(
        "<", "&lt;",
        ">", "&gt;",
        "'", "&#39;",
        "\"", "&quot;",
        "&", "&amp;",
    )
    return replacer.Replace(s)
}

func sanitizeMap(data map[string]interface{}) map[string]interface{} {
    sanitized := make(map[string]interface{})
    for key, value := range data {
        switch v := value.(type) {
        case string:
            sanitized[key] = sanitizeString(v)
        case map[string]interface{}:
            sanitized[key] = sanitizeMap(v)
        default:
            sanitized[key] = v
        }
    }
    return sanitized
}