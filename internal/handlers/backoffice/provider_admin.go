package backoffice

import (
    "github.com/gofiber/fiber/v2"
    
    "fintech_api_golang/internal/core/services"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
)

type ProviderAdminHandler struct {
    adminService *services.AdminService
}

func NewProviderAdminHandler(adminService *services.AdminService) *ProviderAdminHandler {
    return &ProviderAdminHandler{
        adminService: adminService,
    }
}

// ListProviders - GET /admin/providers
func (h *ProviderAdminHandler) ListProviders(c *fiber.Ctx) error {
    providerType := c.Query("type")
    
    providers, err := h.adminService.ListProviders(c.Context(), providerType)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Providers",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    providers,
    })
}

// ToggleProvider - PUT /admin/providers/:id/toggle
func (h *ProviderAdminHandler) ToggleProvider(c *fiber.Ctx) error {
    providerID := c.Params("id")
    
    var req request.ToggleProviderRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.adminService.ToggleProvider(c.Context(), providerID, req.IsActive); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Toggle Failed",
            Message: err.Error(),
        })
    }
    
    status := "disabled"
    if req.IsActive {
        status = "enabled"
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Provider " + status + " successfully",
    })
}

// SetProviderPriority - PUT /admin/providers/:id/priority
func (h *ProviderAdminHandler) SetProviderPriority(c *fiber.Ctx) error {
    providerID := c.Params("id")
    
    var req request.SetProviderPriorityRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.adminService.SetProviderPriority(c.Context(), providerID, req.Priority); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Priority Update Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Provider priority updated successfully",
    })
}

// CheckProviderHealth - GET /admin/providers/:id/health
func (h *ProviderAdminHandler) CheckProviderHealth(c *fiber.Ctx) error {
    providerID := c.Params("id")
    
    status, err := h.adminService.CheckProviderHealth(c.Context(), providerID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Health Check Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data: map[string]string{
            "status": status,
        },
    })
}

// GetProviderLogs - GET /admin/providers/logs
func (h *ProviderAdminHandler) GetProviderLogs(c *fiber.Ctx) error {
    providerID := c.Query("provider_id")
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 50)
    
    logs, total, err := h.adminService.GetProviderLogs(c.Context(), providerID, page, limit)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Logs",
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
            Data:       logs,
            Total:      total,
            Page:       page,
            Limit:      limit,
            TotalPages: totalPages,
            HasNext:    page < totalPages,
            HasPrev:    page > 1,
        },
    })
}