package backoffice

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    
    "fintech_api_golang/internal/core/services"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
    "fintech_api_golang/internal/middleware"
)

type WalletAdminHandler struct {
    adminService *services.AdminService
}

func NewWalletAdminHandler(adminService *services.AdminService) *WalletAdminHandler {
    return &WalletAdminHandler{
        adminService: adminService,
    }
}

// ListWallets - GET /admin/wallets
func (h *WalletAdminHandler) ListWallets(c *fiber.Ctx) error {
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 20)
    offset := (page - 1) * limit
    
    filters := make(map[string]interface{})
    if isLocked := c.Query("is_locked"); isLocked != "" {
        filters["is_locked"] = isLocked == "true"
    }
    if currency := c.Query("currency"); currency != "" {
        filters["currency"] = currency
    }
    
    wallets, total, err := h.adminService.ListWallets(c.Context(), offset, limit, filters)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Wallets",
            Message: err.Error(),
        })
    }
    
    totalPages := int(total) / limit
    if int(total)%limit > 0 {
        totalPages++
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data: response.PaginatedResponse{
            Data:       wallets,
            Total:      total,
            Page:       page,
            Limit:      limit,
            TotalPages: totalPages,
            HasNext:    page < totalPages,
            HasPrev:    page > 1,
        },
    })
}

// GetWalletDetails - GET /admin/wallets/:id
func (h *WalletAdminHandler) GetWalletDetails(c *fiber.Ctx) error {
    walletID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid wallet ID",
        })
    }
    
    wallet, err := h.adminService.GetWalletDetails(c.Context(), walletID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
            Error:   "Wallet Not Found",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    wallet,
    })
}

// CreditWallet - POST /admin/wallets/credit
func (h *WalletAdminHandler) CreditWallet(c *fiber.Ctx) error {
    adminID, err := middleware.GetAdminIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "Admin not authenticated",
        })
    }
    
    var req request.ManualCreditRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.adminService.ManualCredit(c.Context(), &req, adminID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Credit Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Wallet credited successfully",
        Data:    resp,
    })
}

// DebitWallet - POST /admin/wallets/debit
func (h *WalletAdminHandler) DebitWallet(c *fiber.Ctx) error {
    adminID, err := middleware.GetAdminIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "Admin not authenticated",
        })
    }
    
    var req request.ManualDebitRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.adminService.ManualDebit(c.Context(), &req, adminID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Debit Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Wallet debited successfully",
        Data:    resp,
    })
}

// FreezeWallet - POST /admin/wallets/freeze
func (h *WalletAdminHandler) FreezeWallet(c *fiber.Ctx) error {
    var req struct {
        UserID string `json:"user_id"`
        Reason string `json:"reason"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    userID, err := uuid.Parse(req.UserID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid user ID",
        })
    }
    
    if err := h.adminService.FreezeWallet(c.Context(), userID, req.Reason); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Freeze Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Wallet frozen successfully",
    })
}

// UnfreezeWallet - POST /admin/wallets/unfreeze
func (h *WalletAdminHandler) UnfreezeWallet(c *fiber.Ctx) error {
    var req struct {
        UserID string `json:"user_id"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    userID, err := uuid.Parse(req.UserID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid user ID",
        })
    }
    
    if err := h.adminService.UnfreezeWallet(c.Context(), userID); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Unfreeze Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Wallet unfrozen successfully",
    })
}

// GetBalanceSummary - GET /admin/wallets/balances/summary
func (h *WalletAdminHandler) GetBalanceSummary(c *fiber.Ctx) error {
    summary, err := h.adminService.GetBalanceSummary(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Summary",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    summary,
    })
}