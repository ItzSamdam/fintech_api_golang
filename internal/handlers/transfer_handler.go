package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    
    "fintech_api_golang/internal/core/services"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
    "fintech_api_golang/internal/middleware"
)

type TransferHandler struct {
    transferService *services.TransferService
}

func NewTransferHandler(transferService *services.TransferService) *TransferHandler {
    return &TransferHandler{
        transferService: transferService,
    }
}

// SendTransfer - POST /transfers/send
func (h *TransferHandler) SendTransfer(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.SendTransferRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.transferService.SendTransfer(c.Context(), userID, &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Transfer Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Transfer initiated successfully",
        Data:    resp,
    })
}

// GetTransferStatus - GET /transfers/status/:reference
func (h *TransferHandler) GetTransferStatus(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    reference := c.Params("reference")
    if reference == "" {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Reference is required",
        })
    }
    
    resp, err := h.transferService.GetTransferStatus(c.Context(), userID, reference)
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

// RetryTransfer - POST /transfers/retry
func (h *TransferHandler) RetryTransfer(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.RetryTransferRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.transferService.RetryTransfer(c.Context(), userID, req.Reference)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Retry Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Transfer retry initiated",
        Data:    resp,
    })
}

// NameEnquiry - POST /transfers/name-enquiry
func (h *TransferHandler) NameEnquiry(c *fiber.Ctx) error {
    var req request.NameEnquiryRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.transferService.NameEnquiry(c.Context(), &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Name Enquiry Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// GetBanks - GET /banks/list
func (h *TransferHandler) GetBanks(c *fiber.Ctx) error {
    resp, err := h.transferService.GetBanks(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Banks",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// GetTransferHistory - GET /transfers/history
func (h *TransferHandler) GetTransferHistory(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    offset := c.QueryInt("offset", 0)
    limit := c.QueryInt("limit", 20)
    
    resp, err := h.transferService.GetTransferHistory(c.Context(), userID, offset, limit)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch History",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// SendToWallet - POST /transfers/send-to-wallet
func (h *TransferHandler) SendToWallet(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req struct {
        RecipientWalletID string `json:"recipient_wallet_id" validate:"required"`
        Amount            int64  `json:"amount" validate:"required,min=100"`
        Narration         string `json:"narration"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    recipientID, err := uuid.Parse(req.RecipientWalletID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid recipient wallet ID",
        })
    }
    
    resp, err := h.transferService.SendToWallet(c.Context(), userID, recipientID, req.Amount, req.Narration)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Transfer Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Transfer completed successfully",
        Data:    resp,
    })
}