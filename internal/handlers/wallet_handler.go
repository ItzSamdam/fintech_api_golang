package handlers

import (
    "github.com/gofiber/fiber/v2"
    
    "fintech_api_golang/internal/core/services"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
    "fintech_api_golang/internal/middleware"
)

type WalletHandler struct {
    walletService *services.WalletService
}

func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
    return &WalletHandler{
        walletService: walletService,
    }
}

// CreateWallet - POST /wallets/create
func (h *WalletHandler) CreateWallet(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.CreateWalletRequest
    c.BodyParser(&req)
    
    resp, err := h.walletService.CreateWallet(c.Context(), userID, req.Currency)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Wallet Creation Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Wallet created successfully",
        Data:    resp,
    })
}

// GetBalance - GET /wallets/balance
func (h *WalletHandler) GetBalance(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    resp, err := h.walletService.GetBalance(c.Context(), userID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
            Error:   "Not Found",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// GetTransactions - GET /wallets/transactions
func (h *WalletHandler) GetTransactions(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.GetTransactionsRequest
    if err := c.QueryParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid query parameters",
        })
    }
    
    resp, err := h.walletService.GetTransactions(c.Context(), userID, &req)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Transactions",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// GetTransactionByID - GET /wallets/transactions/:id
func (h *WalletHandler) GetTransactionByID(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    transactionID := c.Params("id")
    if transactionID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "ID is required",
        })
    }
    
    resp, err := h.walletService.GetTransactionByID(c.Context(), userID, transactionID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
            Error:   "Token Not Found",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:  resp,
    })
}


// GetLimits - GET /wallets/limits
func (h *WalletHandler) GetLimits(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    resp, err := h.walletService.GetLimits(c.Context(), userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Limits",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// LockWallet - POST /wallets/lock
func (h *WalletHandler) LockWallet(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.WalletActionRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.walletService.LockWallet(c.Context(), userID, req.Reason); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Lock Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Wallet locked successfully",
    })
}

// UnlockWallet - POST /wallets/unlock
func (h *WalletHandler) UnlockWallet(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    if err := h.walletService.UnlockWallet(c.Context(), userID); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Unlock Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Wallet unlocked successfully",
    })
}

// GetStatement - GET /wallets/statement
func (h *WalletHandler) GetStatement(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.GetStatementRequest
    if err := c.QueryParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid query parameters",
        })
    }
    
    data, format, err := h.walletService.GetStatement(c.Context(), userID, &req)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Statement Generation Failed",
            Message: err.Error(),
        })
    }
    
    c.Set("Content-Type", format)
    c.Set("Content-Disposition", "attachment; filename=statement."+format)
    
    return c.Send(data)
}