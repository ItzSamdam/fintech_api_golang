package backoffice

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    
    "fintech_api_golang/internal/core/services"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
)

type RoleAdminHandler struct {
    adminService *services.AdminService
}

func NewRoleAdminHandler(adminService *services.AdminService) *RoleAdminHandler {
    return &RoleAdminHandler{
        adminService: adminService,
    }
}

// ListRoles - GET /admin/roles
func (h *RoleAdminHandler) ListRoles(c *fiber.Ctx) error {
    roles, err := h.adminService.ListRoles(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Roles",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    roles,
    })
}

// CreateRole - POST /admin/roles
func (h *RoleAdminHandler) CreateRole(c *fiber.Ctx) error {
    var req request.CreateRoleRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    role, err := h.adminService.CreateRole(c.Context(), &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Role Creation Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse{
        Success: true,
        Message: "Role created successfully",
        Data:    role,
    })
}

// UpdateRole - PUT /admin/roles/:id
func (h *RoleAdminHandler) UpdateRole(c *fiber.Ctx) error {
    roleID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid role ID",
        })
    }
    
    var req request.CreateRoleRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    role, err := h.adminService.UpdateRole(c.Context(), roleID, &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Role Update Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Role updated successfully",
        Data:    role,
    })
}

// DeleteRole - DELETE /admin/roles/:id
func (h *RoleAdminHandler) DeleteRole(c *fiber.Ctx) error {
    roleID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid role ID",
        })
    }
    
    if err := h.adminService.DeleteRole(c.Context(), roleID); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Role Deletion Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Role deleted successfully",
    })
}

// ListStaff - GET /admin/staff
func (h *RoleAdminHandler) ListStaff(c *fiber.Ctx) error {
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 20)
    offset := (page - 1) * limit
    
    staff, total, err := h.adminService.ListStaff(c.Context(), offset, limit)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Staff",
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
            Data:       staff,
            Total:      total,
            Page:       page,
            Limit:      limit,
            TotalPages: totalPages,
            HasNext:    page < totalPages,
            HasPrev:    page > 1,
        },
    })
}

// InviteStaff - POST /admin/staff/invite
func (h *RoleAdminHandler) InviteStaff(c *fiber.Ctx) error {
    var req request.InviteStaffRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.adminService.InviteStaff(c.Context(), &req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Invitation Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Staff invitation sent successfully",
    })
}

// AssignRole - PUT /admin/staff/:id/role
func (h *RoleAdminHandler) AssignRole(c *fiber.Ctx) error {
    staffID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid staff ID",
        })
    }
    
    var req struct {
        Role string `json:"role"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.adminService.AssignRole(c.Context(), staffID, req.Role); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Role Assignment Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Role assigned successfully",
    })
}

// RemoveStaff - DELETE /admin/staff/:id
func (h *RoleAdminHandler) RemoveStaff(c *fiber.Ctx) error {
    staffID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid staff ID",
        })
    }
    
    if err := h.adminService.RemoveStaff(c.Context(), staffID); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Staff Removal Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Staff removed successfully",
    })
}

// GetStaffAudit - GET /admin/staff/:id/audit
func (h *RoleAdminHandler) GetStaffAudit(c *fiber.Ctx) error {
    staffID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid staff ID",
        })
    }
    
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 20)
    offset := (page - 1) * limit
    
    logs, total, err := h.adminService.GetStaffAudit(c.Context(), staffID, offset, limit)
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