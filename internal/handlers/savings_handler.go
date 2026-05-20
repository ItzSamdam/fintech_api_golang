package handlers

import (
    "github.com/gofiber/fiber/v2"
    
    "fintech_api_golang/internal/core/services"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
    "fintech_api_golang/internal/middleware"
)

type SavingsHandler struct {
    savingsService *services.SavingsService
}

func NewSavingsHandler(savingsService *services.SavingsService) *SavingsHandler {
    return &SavingsHandler{
        savingsService: savingsService,
    }
}

// CreateGoal - POST /savings/goals/create
func (h *SavingsHandler) CreateGoal(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.CreateSavingsGoalRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.savingsService.CreateGoal(c.Context(), userID, &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Failed to Create Goal",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse{
        Success: true,
        Message: "Savings goal created successfully",
        Data:    resp,
    })
}

// ContributeToGoal - POST /savings/goals/contribute
func (h *SavingsHandler) ContributeToGoal(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.ContributeToGoalRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.savingsService.ContributeToGoal(c.Context(), userID, &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Contribution Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Contribution added successfully",
        Data:    resp,
    })
}

// GetGoals - GET /savings/goals
func (h *SavingsHandler) GetGoals(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    resp, err := h.savingsService.GetGoals(c.Context(), userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Goals",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// GetGoal - GET /savings/goals/:id
func (h *SavingsHandler) GetGoal(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    goalID := c.Params("id")
    if goalID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Goal ID is required",
        })
    }
    
    resp, err := h.savingsService.GetGoal(c.Context(), userID, goalID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
            Error:   "Goal Not Found",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}

// UpdateGoal - PUT /savings/goals/:id
func (h *SavingsHandler) UpdateGoal(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    goalID := c.Params("id")
    if goalID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Goal ID is required",
        })
    }
    
    var req request.UpdateSavingsGoalRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.savingsService.UpdateGoal(c.Context(), userID, goalID, &req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Update Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Savings goal updated successfully",
    })
}

// DeleteGoal - DELETE /savings/goals/:id
func (h *SavingsHandler) DeleteGoal(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    goalID := c.Params("id")
    if goalID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Goal ID is required",
        })
    }
    
    if err := h.savingsService.DeleteGoal(c.Context(), userID, goalID); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Deletion Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Savings goal deleted successfully",
    })
}

// ActivateRoundup - POST /savings/roundup/activate
func (h *SavingsHandler) ActivateRoundup(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.ActivateRoundupRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.savingsService.ActivateRoundup(c.Context(), userID, &req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Activation Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Auto roundup activated successfully",
    })
}

// DeactivateRoundup - POST /savings/roundup/deactivate
func (h *SavingsHandler) DeactivateRoundup(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    if err := h.savingsService.DeactivateRoundup(c.Context(), userID); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Deactivation Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Auto roundup deactivated successfully",
    })
}

// GetRoundupStatus - GET /savings/roundup/status
func (h *SavingsHandler) GetRoundupStatus(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    resp, err := h.savingsService.GetRoundupStatus(c.Context(), userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Status",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}