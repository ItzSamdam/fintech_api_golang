package middleware

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"fintech_api_golang/internal/core/entities"
)

type AdminAuthMiddleware struct {
    db     *gorm.DB
    secret string
}

func NewAdminAuthMiddleware(db *gorm.DB, secret string) *AdminAuthMiddleware {
    return &AdminAuthMiddleware{
        db:     db,
        secret: secret,
    }
}

// AdminAuthRequired validates JWT token for admin users
func (m *AdminAuthMiddleware) AdminAuthRequired() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Get token from header
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Missing authorization header",
            })
        }
        
        // Extract token (Bearer <token>)
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Invalid authorization format. Use: Bearer <token>",
            })
        }
        
        tokenString := parts[1]
        
        // Parse and validate token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(m.secret), nil
        })
        
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Invalid or expired token",
            })
        }
        
        if !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Invalid token",
            })
        }
        
        // Extract claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Invalid token claims",
            })
        }
        
        // Verify this is an ADMIN token (not user)
        tokenType, ok := claims["token_type"].(string)
        if !ok || tokenType != "admin" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Invalid token type - admin access required",
            })
        }
        
        // Check if token is expired
        exp, ok := claims["exp"].(float64)
        if ok && time.Now().Unix() > int64(exp) {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Token has expired",
            })
        }
        
        // Get admin ID from claims
        adminIDStr, ok := claims["admin_id"].(string)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Invalid admin ID in token",
            })
        }
        
        adminID, err := uuid.Parse(adminIDStr)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Invalid admin ID format",
            })
        }
        
        // Get admin user details
        var adminUser entities.AdminUser
        if err := m.db.First(&adminUser, "id = ?", adminID).Error; err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Admin user not found",
            })
        }
        
        // Check if admin is active
        if !adminUser.IsActive {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error":   "Forbidden",
                "message": "Admin account is deactivated",
            })
        }
        
        // Add admin user to context
        c.Locals("admin_user", &adminUser)
        c.Locals("admin_id", adminID)
        c.Locals("admin_token", tokenString)
        
        return c.Next()
    }
}

// RequireAdminRole checks if admin has specific role
func (m *AdminAuthMiddleware) RequireAdminRole(allowedRoles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        adminUser := c.Locals("admin_user")
        if adminUser == nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Admin authentication required",
            })
        }
        
        admin, ok := adminUser.(*entities.AdminUser)
        if !ok {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error":   "Internal Server Error",
                "message": "Invalid admin user context",
            })
        }
        
        // Check if admin role is allowed
        allowed := false
        for _, role := range allowedRoles {
            if admin.Role == role {
                allowed = true
                break
            }
        }
        
        if !allowed {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error":      "Forbidden",
                "message":    "Insufficient permissions",
                "required":   allowedRoles,
                "current":    admin.Role,
            })
        }
        
        return c.Next()
    }
}

// RequireAdminPermission checks if admin has specific permission
func (m *AdminAuthMiddleware) RequireAdminPermission(permission string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        adminUser := c.Locals("admin_user")
        if adminUser == nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Admin authentication required",
            })
        }
        
        admin, ok := adminUser.(*entities.AdminUser)
        if !ok {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error":   "Internal Server Error",
                "message": "Invalid admin user context",
            })
        }
        
        // Super admin has all permissions
        if admin.Role == "super_admin" {
            return c.Next()
        }
        
        // Get role permissions from database
        var role entities.Role
        if err := m.db.Where("name = ?", admin.Role).First(&role).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error":   "Internal Server Error",
                "message": "Failed to fetch role permissions",
            })
        }
        
        // Parse permissions and check
        var permissions []string
        if err := json.Unmarshal([]byte(role.Permissions), &permissions); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error":   "Internal Server Error",
                "message": "Invalid permissions format",
            })
        }
        
        // Check if permission exists
        hasPermission := false
        for _, p := range permissions {
            if p == permission || p == "*" {
                hasPermission = true
                break
            }
        }
        
        if !hasPermission {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error":      "Forbidden",
                "message":    "Insufficient permissions",
                "permission": permission,
                "role":       admin.Role,
            })
        }
        
        return c.Next()
    }
}

// Convenience middleware functions
func AdminOnly() fiber.Handler {
    return func(c *fiber.Ctx) error {
        adminUser := c.Locals("admin_user")
        if adminUser == nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Admin authentication required",
            })
        }
        return c.Next()
    }
}

func SuperAdminOnly() fiber.Handler {
    return func(c *fiber.Ctx) error {
        adminUser := c.Locals("admin_user")
        if adminUser == nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Admin authentication required",
            })
        }
        
        admin, ok := adminUser.(*entities.AdminUser)
        if !ok {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error":   "Internal Server Error",
                "message": "Invalid admin user context",
            })
        }
        
        if admin.Role != "super_admin" {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error":   "Forbidden",
                "message": "Super admin access required",
            })
        }
        
        return c.Next()
    }
}

// GetAdminFromContext retrieves admin user from context
func GetAdminFromContext(c *fiber.Ctx) (*entities.AdminUser, error) {
    admin := c.Locals("admin_user")
    if admin == nil {
        return nil, fmt.Errorf("admin user not found in context")
    }
    
    a, ok := admin.(*entities.AdminUser)
    if !ok {
        return nil, fmt.Errorf("invalid admin user type in context")
    }
    
    return a, nil
}

// GetAdminIDFromContext retrieves admin ID from context
func GetAdminIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
    adminID := c.Locals("admin_id")
    if adminID == nil {
        return uuid.Nil, fmt.Errorf("admin ID not found in context")
    }
    
    aid, ok := adminID.(uuid.UUID)
    if !ok {
        return uuid.Nil, fmt.Errorf("invalid admin ID type in context")
    }
    
    return aid, nil
}