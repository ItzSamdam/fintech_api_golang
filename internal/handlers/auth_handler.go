package handlers

import (
    "github.com/gofiber/fiber/v2"
    
    "fintech_api_golang/internal/core/services"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
    "fintech_api_golang/internal/middleware"
)

type AuthHandler struct {
    authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
    return &AuthHandler{
        authService: authService,
    }
}

// RegisterPhone - POST /auth/register/phone
func (h *AuthHandler) RegisterPhone(c *fiber.Ctx) error {
    var req request.RegisterPhoneRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.authService.RegisterPhone(c.Context(), &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Registration Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "OTP sent successfully",
        Data:    resp,
    })
}

// VerifyOTP - POST /auth/verify/otp
func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
    var req request.VerifyOTPRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.authService.VerifyOTP(c.Context(), &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Verification Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "OTP verified successfully",
        Data:    resp,
    })
}

// RegisterBVN - POST /auth/register/bvn
func (h *AuthHandler) RegisterBVN(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.RegisterBVNRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.authService.RegisterBVN(c.Context(), userID, &req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "BVN Registration Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "BVN verified successfully. Tier upgraded.",
    })
}

// VerifyFace - POST /auth/verify/face
func (h *AuthHandler) VerifyFace(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.VerifyFaceRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.authService.VerifyFace(c.Context(), userID, &req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Face Verification Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Face verified successfully. Tier upgraded.",
    })
}

// Login - POST /auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
    var req request.LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.authService.Login(c.Context(), &req)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Login Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Login successful",
        Data:    resp,
    })
}

// GetMe - GET /auth/me
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    resp, err := h.authService.GetUserProfile(c.Context(), userID)
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

// UpdateMe - PUT /auth/me
func (h *AuthHandler) UpdateMe(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.UpdateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.authService.UpdateUserProfile(c.Context(), userID, &req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Update Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Profile updated successfully",
    })
}

// ChangePassword - POST /auth/change-password
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    var req request.ChangePasswordRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.authService.ChangePassword(c.Context(), userID, &req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Password Change Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Password changed successfully",
    })
}

// ResetPassword - POST /auth/reset-password
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
    var req request.ResetPasswordRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    if err := h.authService.ResetPassword(c.Context(), &req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Password Reset Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Password reset successfully",
    })
}

// Logout - POST /auth/logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User not authenticated",
        })
    }
    
    token := c.Get("Authorization")
    if len(token) > 7 && token[:7] == "Bearer " {
        token = token[7:]
    }
    
    var req request.LogoutRequest
    c.BodyParser(&req)
    
    if err := h.authService.Logout(c.Context(), userID, token, req.AllDevices); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Logout Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Message: "Logged out successfully",
    })
}

// RefreshToken - POST /auth/refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
    var req request.RefreshTokenRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    resp, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
            Error:   "Refresh Failed",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    resp,
    })
}