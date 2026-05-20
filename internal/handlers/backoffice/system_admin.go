package backoffice

import (
    "github.com/gofiber/fiber/v2"
    
    "fintech_api_golang/internal/core/services"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
)

type SystemAdminHandler struct {
    adminService *services.AdminService
}

func NewSystemAdminHandler(adminService *services.AdminService) *SystemAdminHandler {
    return &SystemAdminHandler{
        adminService: adminService,
    }
}

// GetSystemSettings - GET /admin/settings
func (h *SystemAdminHandler) GetSystemSettings(c *fiber.Ctx) error {
    settings, err := h.adminService.GetSystemSettings(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Settings",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    settings,
    })
}

// UpdateSystemSettings - PUT /admin/settings
func (h *SystemAdminHandler) UpdateSystemSettings(c *fiber.Ctx) error {
    var req request.SystemSettingsRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.adminService.UpdateSystemSettings(c.Context(), &req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Update Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "System settings updated successfully",
    })
}

// HealthCheck - GET /admin/health
func (h *SystemAdminHandler) HealthCheck(c *fiber.Ctx) error {
    health, err := h.adminService.HealthCheck(c.Context())
    if err != nil {
        return c.Status(fiber.StatusServiceUnavailable).JSON(response.ErrorResponse{
            Error:   "Health Check Failed",
            Message: err.Error(),
        })
    }
    
    status := fiber.StatusOK
    if health.Status != "healthy" {
        status = fiber.StatusServiceUnavailable
    }
    
    return c.Status(status).JSON(response.SuccessResponse{
        Success: health.Status == "healthy",
        Data:    health,
    })
}

// ListAuditLogs - GET /admin/audit-logs
func (h *SystemAdminHandler) ListAuditLogs(c *fiber.Ctx) error {
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 20)
    offset := (page - 1) * limit
    
    filters := make(map[string]interface{})
    if adminID := c.Query("admin_id"); adminID != "" {
        filters["admin_id"] = adminID
    }
    if action := c.Query("action"); action != "" {
        filters["action"] = action
    }
    if fromDate := c.Query("from_date"); fromDate != "" {
        filters["from_date"] = fromDate
    }
    if toDate := c.Query("to_date"); toDate != "" {
        filters["to_date"] = toDate
    }
    
    logs, total, err := h.adminService.GetAuditLogs(c.Context(), offset, limit, filters)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Audit Logs",
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

// GetAuditLog - GET /admin/audit-logs/:id
func (h *SystemAdminHandler) GetAuditLog(c *fiber.Ctx) error {
    logID := c.Params("id")
    
    log, err := h.adminService.GetAuditLog(c.Context(), logID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
            Error:   "Audit Log Not Found",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    log,
    })
}

// TriggerDatabaseBackup - POST /admin/backup/database
func (h *SystemAdminHandler) TriggerDatabaseBackup(c *fiber.Ctx) error {
    backup, err := h.adminService.TriggerDatabaseBackup(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Backup Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Database backup triggered successfully",
        Data:    backup,
    })
}

// GetSystemMetrics - GET /admin/metrics
func (h *SystemAdminHandler) GetSystemMetrics(c *fiber.Ctx) error {
    metrics, err := h.adminService.GetSystemMetrics(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Metrics",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    metrics,
    })
}