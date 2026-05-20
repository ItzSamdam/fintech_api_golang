package backoffice

import (
    "github.com/gofiber/fiber/v2"
    
    "fintech_api_golang/internal/core/services"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
)

type TransactionAdminHandler struct {
    adminService *services.AdminService
}

func NewTransactionAdminHandler(adminService *services.AdminService) *TransactionAdminHandler {
    return &TransactionAdminHandler{
        adminService: adminService,
    }
}

// ListTransactions - GET /admin/transactions
func (h *TransactionAdminHandler) ListTransactions(c *fiber.Ctx) error {
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 20)
    offset := (page - 1) * limit
    
    filters := make(map[string]interface{})
    if category := c.Query("category"); category != "" {
        filters["category"] = category
    }
    if status := c.Query("status"); status != "" {
        filters["status"] = status
    }
    if fromDate := c.Query("from_date"); fromDate != "" {
        filters["from_date"] = fromDate
    }
    if toDate := c.Query("to_date"); toDate != "" {
        filters["to_date"] = toDate
    }
    if userID := c.Query("user_id"); userID != "" {
        filters["user_id"] = userID
    }
    
    transactions, total, err := h.adminService.ListTransactions(c.Context(), offset, limit, filters)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Transactions",
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
            Data:       transactions,
            Total:      total,
            Page:       page,
            Limit:      limit,
            TotalPages: totalPages,
            HasNext:    page < totalPages,
            HasPrev:    page > 1,
        },
    })
}

// GetTransactionDetails - GET /admin/transactions/:id
func (h *TransactionAdminHandler) GetTransactionDetails(c *fiber.Ctx) error {
    transactionID := c.Params("id")
    
    transaction, err := h.adminService.GetTransactionDetails(c.Context(), transactionID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
            Error:   "Transaction Not Found",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    transaction,
    })
}

// ReverseTransaction - POST /admin/transactions/reverse
func (h *TransactionAdminHandler) ReverseTransaction(c *fiber.Ctx) error {
    var req request.ReverseTransactionRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.adminService.ReverseTransaction(c.Context(), req.TransactionID, req.Reason, req.NotifyUser); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Reversal Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Transaction reversed successfully",
    })
}

// GetTransactionSummary - GET /admin/transactions/summary
func (h *TransactionAdminHandler) GetTransactionSummary(c *fiber.Ctx) error {
    startDate := c.Query("start_date")
    endDate := c.Query("end_date")
    
    summary, err := h.adminService.GetTransactionSummary(c.Context(), startDate, endDate)
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

// VoidTransaction - POST /admin/transactions/void
func (h *TransactionAdminHandler) VoidTransaction(c *fiber.Ctx) error {
    var req struct {
        TransactionID string `json:"transaction_id"`
        Reason        string `json:"reason"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.adminService.VoidTransaction(c.Context(), req.TransactionID, req.Reason); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Void Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Transaction voided successfully",
    })
}