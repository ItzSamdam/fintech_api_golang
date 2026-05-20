package handlers

import (
    "github.com/gofiber/fiber/v2"
    
    "fintech_api_golang/internal/core/services"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
    "fintech_api_golang/internal/middleware"
)

type BillHandler struct {
    billService *services.BillPaymentService
}

func NewBillHandler(billService *services.BillPaymentService) *BillHandler {
    return &BillHandler{
        billService: billService,
    }
}

// ========== AIRTIME ==========

// GetAirtimeNetworks - GET /airtime/networks
func (h *BillHandler) GetAirtimeNetworks(c *fiber.Ctx) error {
    resp, err := h.billService.GetAirtimeNetworks(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Networks",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// PurchaseAirtime - POST /airtime/purchase
func (h *BillHandler) PurchaseAirtime(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.PurchaseAirtimeRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.billService.PurchaseAirtime(c.Context(), userID, &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Airtime Purchase Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Airtime purchased successfully",
        Data:    resp,
    })
}

// GetAirtimeHistory - GET /airtime/history
func (h *BillHandler) GetAirtimeHistory(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    offset := c.QueryInt("offset", 0)
    limit := c.QueryInt("limit", 20)
    
    resp, err := h.billService.GetBettingHistory(c.Context(), userID, offset, limit)
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

// ========== DATA ==========

// GetDataNetworks - GET /data/networks
func (h *BillHandler) GetDataNetworks(c *fiber.Ctx) error {
    resp, err := h.billService.GetDataNetworks(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Networks",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// GetDataPlans - GET /data/plans/:network
func (h *BillHandler) GetDataPlans(c *fiber.Ctx) error {
    network := c.Params("network")
    if network == "" {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Network is required",
        })
    }
    
    resp, err := h.billService.GetDataPlans(c.Context(), network)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Data Plans",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// PurchaseData - POST /data/purchase
func (h *BillHandler) PurchaseData(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.PurchaseDataRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.billService.PurchaseData(c.Context(), userID, &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Data Purchase Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Data purchased successfully",
        Data:    resp,
    })
}

// ========== ELECTRICITY ==========

// GetElectricityProviders - GET /electricity/providers
func (h *BillHandler) GetElectricityProviders(c *fiber.Ctx) error {
    resp, err := h.billService.GetElectricityProviders(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Providers",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// ValidateMeter - POST /electricity/validate-meter
func (h *BillHandler) ValidateMeter(c *fiber.Ctx) error {
    var req request.ValidateMeterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.billService.ValidateMeter(c.Context(), &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Meter Validation Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// PayElectricity - POST /electricity/pay
func (h *BillHandler) PayElectricity(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.PayElectricityRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.billService.PayElectricity(c.Context(), userID, &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Electricity Payment Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Electricity bill paid successfully",
        Data:    resp,
    })
}

// GetElectricityToken - GET /electricity/token/:transaction_id
func (h *BillHandler) GetElectricityToken(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    transactionID := c.Params("transaction_id")
    if transactionID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Transaction ID is required",
        })
    }
    
    token, err := h.billService.GetElectricityToken(c.Context(), userID, transactionID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
            Error:   "Token Not Found",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data: map[string]string{
            "token": token,
        },
    })
}

// ========== BETTING ==========

// GetBettingProviders - GET /betting/providers
func (h *BillHandler) GetBettingProviders(c *fiber.Ctx) error {
    resp, err := h.billService.GetBettingProviders(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Providers",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// ValidateBettingAccount - POST /betting/validate-account
func (h *BillHandler) ValidateBettingAccount(c *fiber.Ctx) error {
    var req request.ValidateBettingAccountRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.billService.ValidateBettingAccount(c.Context(), &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Account Validation Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// FundBettingWallet - POST /betting/fund
func (h *BillHandler) FundBettingWallet(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.FundBettingRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.billService.FundBettingWallet(c.Context(), userID, &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Betting Funding Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Betting wallet funded successfully",
        Data:    resp,
    })
}

// GetBettingHistory - GET /betting/history
func (h *BillHandler) GetBettingHistory(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    offset := c.QueryInt("offset", 0)
    limit := c.QueryInt("limit", 20)
    
    resp, err := h.billService.GetBettingHistory(c.Context(), userID, offset, limit)
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