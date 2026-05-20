package backoffice

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    
    "fintech_api_golang/internal/core/services"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
    "fintech_api_golang/internal/middleware"
)

type KYCAdminHandler struct {
    adminService *services.AdminService
}

func NewKYCAdminHandler(adminService *services.AdminService) *KYCAdminHandler {
    return &KYCAdminHandler{
        adminService: adminService,
    }
}

// ListPendingKYC - GET /admin/kyc/pending
func (h *KYCAdminHandler) ListPendingKYC(c *fiber.Ctx) error {
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 20)
    offset := (page - 1) * limit
    
    kycList, total, err := h.adminService.GetPendingKYC(c.Context(), offset, limit)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Pending KYC",
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
            Data:       kycList,
            Total:      total,
            Page:       page,
            Limit:      limit,
            TotalPages: totalPages,
            HasNext:    page < totalPages,
            HasPrev:    page > 1,
        },
    })
}

// GetKYCDetails - GET /admin/kyc/:id
func (h *KYCAdminHandler) GetKYCDetails(c *fiber.Ctx) error {
    kycID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid KYC ID",
        })
    }
    
    kyc, err := h.adminService.GetKYCDetails(c.Context(), kycID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
            Error:   "KYC Not Found",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    kyc,
    })
}

// ApproveKYC - POST /admin/kyc/:id/approve
func (h *KYCAdminHandler) ApproveKYC(c *fiber.Ctx) error {
    adminID, err := middleware.GetAdminIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "Admin not authenticated",
        })
    }
    
    kycID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid KYC ID",
        })
    }
    
    var req request.ApproveKYCRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.adminService.ApproveKYC(c.Context(), kycID.String(), req.Notes, adminID); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Approval Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "KYC approved successfully",
    })
}

// RejectKYC - POST /admin/kyc/:id/reject
func (h *KYCAdminHandler) RejectKYC(c *fiber.Ctx) error {
    adminID, err := middleware.GetAdminIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "Admin not authenticated",
        })
    }
    
    kycID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid KYC ID",
        })
    }
    
    var req request.RejectKYCRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.adminService.RejectKYC(c.Context(), kycID.String(), req.Reason, adminID); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Rejection Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "KYC rejected successfully",
    })
}