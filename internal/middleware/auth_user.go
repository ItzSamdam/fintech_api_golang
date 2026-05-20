package middleware

import (
    "strings"
    "fmt"
    "time"
    
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"fintech_api_golang/internal/core/entities"
)

type UserAuthMiddleware struct {
    db     *gorm.DB
    secret string
}

func NewUserAuthMiddleware(db *gorm.DB, secret string) *UserAuthMiddleware {
    return &UserAuthMiddleware{
        db:     db,
        secret: secret,
    }
}

// UserAuthRequired validates JWT token for customers
func (m *UserAuthMiddleware) UserAuthRequired() fiber.Handler {
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
        
        // Verify this is a USER token (not admin)
        tokenType, ok := claims["token_type"].(string)
        if !ok || tokenType != "user" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Invalid token type",
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
        
        // Get user ID from claims
        userIDStr, ok := claims["user_id"].(string)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Invalid user ID in token",
            })
        }
        
        userID, err := uuid.Parse(userIDStr)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "Invalid user ID format",
            })
        }
        
        // Check if session exists and is active
        var session entities.Session
        err = m.db.Where("user_id = ? AND token = ? AND is_active = ? AND expires_at > ?", 
            userID, tokenString, true, time.Now()).First(&session).Error
        
        if err != nil {
            if err == gorm.ErrRecordNotFound {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "error":   "Unauthorized",
                    "message": "Session not found or expired",
                })
            }
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error":   "Internal Server Error",
                "message": "Database error",
            })
        }
        
        // Get user details
        var user entities.User
        if err := m.db.First(&user, "id = ?", userID).Error; err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error":   "Unauthorized",
                "message": "User not found",
            })
        }
        
        // Check if user is suspended
        if user.IsSuspended {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error":   "Forbidden",
                "message": "Account is suspended",
                "reason":  user.SuspensionReason,
            })
        }
        
        // Check if user is active
        if !user.IsActive {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error":   "Forbidden",
                "message": "Account is deactivated",
            })
        }
        
        // Add user and session to context
        c.Locals("user", &user)
        c.Locals("user_id", userID)
        c.Locals("session_id", session.ID)
        c.Locals("token", tokenString)
        
        return c.Next()
    }
}

// OptionalUserAuth validates token if present but doesn't require it
func (m *UserAuthMiddleware) OptionalUserAuth() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Next()
        }
        
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return c.Next()
        }
        
        tokenString := parts[1]
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(m.secret), nil
        })
        
        if err == nil && token.Valid {
            claims, ok := token.Claims.(jwt.MapClaims)
            if ok {
                tokenType, ok := claims["token_type"].(string)
                if ok && tokenType == "user" {
                    userIDStr, ok := claims["user_id"].(string)
                    if ok {
                        userID, err := uuid.Parse(userIDStr)
                        if err == nil {
                            var user entities.User
                            if m.db.First(&user, "id = ?", userID).Error == nil {
                                c.Locals("user", &user)
                                c.Locals("user_id", userID)
                            }
                        }
                    }
                }
            }
        }
        
        return c.Next()
    }
}

// GetUserFromContext retrieves user from context
func GetUserFromContext(c *fiber.Ctx) (*entities.User, error) {
    user := c.Locals("user")
    if user == nil {
        return nil, fmt.Errorf("user not found in context")
    }
    
    u, ok := user.(*entities.User)
    if !ok {
        return nil, fmt.Errorf("invalid user type in context")
    }
    
    return u, nil
}

// GetUserIDFromContext retrieves user ID from context
func GetUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
    userID := c.Locals("user_id")
    if userID == nil {
        return uuid.Nil, fmt.Errorf("user ID not found in context")
    }
    
    uid, ok := userID.(uuid.UUID)
    if !ok {
        return uuid.Nil, fmt.Errorf("invalid user ID type in context")
    }
    
    return uid, nil
}