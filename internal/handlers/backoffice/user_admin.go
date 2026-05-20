package backoffice

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    
    "fintech_api_golang/internal/core/services"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
    "fintech_api_golang/internal/middleware"
)

type UserAdminHandler struct {
    adminService *services.AdminService
}

func NewUserAdminHandler(adminService *services.AdminService) *UserAdminHandler {
    return &UserAdminHandler{
        adminService: adminService,
    }
}

// ListUsers - GET /admin/users
func (h *UserAdminHandler) ListUsers(c *fiber.Ctx) error {
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 20)
    offset := (page - 1) * limit
    
    filters := make(map[string]interface{})
    if tier := c.Query("tier"); tier != "" {
        filters["tier"] = tier
    }
    if status := c.Query("status"); status != "" {
        if status == "active" {
            filters["is_active"] = true
        } else if status == "suspended" {
            filters["is_suspended"] = true
        }
    }
    if fromDate := c.Query("from_date"); fromDate != "" {
        filters["from_date"] = fromDate
    }
    if toDate := c.Query("to_date"); toDate != "" {
        filters["to_date"] = toDate
    }
    
    resp, err := h.adminService.ListUsers(c.Context(), offset, limit, filters)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Users",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// GetUserDetails - GET /admin/users/:id
func (h *UserAdminHandler) GetUserDetails(c *fiber.Ctx) error {
    userID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid user ID",
        })
    }
    
    resp, err := h.adminService.GetUserDetails(c.Context(), userID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
            Error:   "User Not Found",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// UpgradeUserTier - POST /admin/users/:id/tier/upgrade
func (h *UserAdminHandler) UpgradeUserTier(c *fiber.Ctx) error {
    adminID, err := middleware.GetAdminIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "Admin not authenticated",
        })
    }
    
    userID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid user ID",
        })
    }
    
    var req request.UpgradeTierRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.adminService.UpgradeUserTier(c.Context(), userID, req.Tier, adminID); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Tier Upgrade Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "User tier upgraded successfully",
    })
}

// SuspendUser - POST /admin/users/:id/suspend
func (h *UserAdminHandler) SuspendUser(c *fiber.Ctx) error {
    userID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid user ID",
        })
    }
    
    var req request.SuspendUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.adminService.SuspendUser(c.Context(), userID, req.Reason, &req.Duration); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Suspension Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "User suspended successfully",
    })
}

// UnsuspendUser - POST /admin/users/:id/unsuspend
func (h *UserAdminHandler) UnsuspendUser(c *fiber.Ctx) error {
    userID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid user ID",
        })
    }
    
    if err := h.adminService.UnsuspendUser(c.Context(), userID); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Unsuspension Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "User unsuspended successfully",
    })
}

// DeleteUser - DELETE /admin/users/:id
func (h *UserAdminHandler) DeleteUser(c *fiber.Ctx) error {
    userID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid user ID",
        })
    }
    
    if err := h.adminService.DeleteUser(c.Context(), userID); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Deletion Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "User deleted successfully",
    })
}

// OverrideLimits - PUT /admin/users/:id/limits
func (h *UserAdminHandler) OverrideLimits(c *fiber.Ctx) error {
    userID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid user ID",
        })
    }
    
    var req request.OverrideLimitsRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.adminService.OverrideLimits(c.Context(), userID, &req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Override Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "User limits overridden successfully",
    })
}

// SearchUsers - GET /admin/users/search
func (h *UserAdminHandler) SearchUsers(c *fiber.Ctx) error {
    query := c.Query("q")
    if query == "" {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Search query is required",
        })
    }
    
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 20)
    offset := (page - 1) * limit
    
    users, total, err := h.adminService.SearchUsers(c.Context(), query, offset, limit)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Search Failed",
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
            Data:       users,
            Total:      total,
            Page:       page,
            Limit:      limit,
            TotalPages: totalPages,
            HasNext:    page < totalPages,
            HasPrev:    page > 1,
        },
    })
}