package middleware

import (
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
    
    "fintech_api_golang/internal/core/entities"
)

type TierLimiter struct {
    db *gorm.DB
}

func NewTierLimiter(db *gorm.DB) *TierLimiter {
    return &TierLimiter{
        db: db,
    }
}

// CheckTierLimit checks if user can perform action based on tier
func (t *TierLimiter) CheckTierLimit(action string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user, err := GetUserFromContext(c)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "User not authenticated",
            })
        }
        
        // Get tier limits
        var tierLimit entities.TierLimit
        if err := t.db.Where("tier = ?", user.Tier).First(&tierLimit).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Failed to fetch tier limits",
            })
        }
        
        // Check specific action
        canPerform := false
        switch action {
        case "send_money":
            canPerform = tierLimit.CanSendMoney
        case "buy_airtime":
            canPerform = tierLimit.CanBuyAirtime
        case "buy_data":
            canPerform = tierLimit.CanBuyData
        case "pay_electricity":
            canPerform = tierLimit.CanPayElectricity
        case "fund_betting":
            canPerform = tierLimit.CanFundBetting
        case "save":
            canPerform = tierLimit.CanSave
        default:
            canPerform = true
        }
        
        if !canPerform {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error":   "Action not allowed for your tier",
                "tier":    user.Tier,
                "action":  action,
                "upgrade_required": true,
            })
        }
        
        c.Locals("tier_limit", &tierLimit)
        return c.Next()
    }
}

// CheckTransactionLimit checks if transaction amount is within tier limits
func (t *TierLimiter) CheckTransactionLimit() fiber.Handler {
    return func(c *fiber.Ctx) error {
        user, err := GetUserFromContext(c)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "User not authenticated",
            })
        }
        
        // Parse amount from request body (implement based on your DTO)
        var request struct {
            Amount entities.AmountInKobo `json:"amount"`
        }
        
        if err := c.BodyParser(&request); err != nil {
            return c.Next() // Skip if no amount field
        }
        
        // Get tier limits
        var tierLimit entities.TierLimit
        if err := t.db.Where("tier = ?", user.Tier).First(&tierLimit).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Failed to fetch tier limits",
            })
        }
        
        // Check single transaction limit
        if request.Amount > tierLimit.SingleTxLimit {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error":        "Transaction amount exceeds tier limit",
                "amount":       request.Amount,
                "max_allowed":  tierLimit.SingleTxLimit,
                "tier":         user.Tier,
                "upgrade_needed": true,
            })
        }
        
        // Get wallet to check daily/weekly/monthly spent
        var wallet entities.Wallet
        if err := t.db.Where("user_id = ?", user.ID).First(&wallet).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Failed to fetch wallet",
            })
        }
        
        // Check daily limit
        if wallet.DailySpent+request.Amount > tierLimit.DailyLimit {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error":         "Daily transaction limit exceeded",
                "daily_spent":   wallet.DailySpent,
                "daily_limit":   tierLimit.DailyLimit,
                "remaining":     tierLimit.DailyLimit - wallet.DailySpent,
            })
        }
        
        c.Locals("tier_limit", &tierLimit)
        c.Locals("wallet", &wallet)
        
        return c.Next()
    }
}