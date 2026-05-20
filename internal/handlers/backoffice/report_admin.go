package backoffice

import (
    "time"
    
    "github.com/gofiber/fiber/v2"
    
    "fintech_api_golang/internal/core/services"
    "fintech_api_golang/internal/dto/response"
)

type ReportAdminHandler struct {
    adminService *services.AdminService
}

func NewReportAdminHandler(adminService *services.AdminService) *ReportAdminHandler {
    return &ReportAdminHandler{
        adminService: adminService,
    }
}

// GetDailyReport - GET /admin/reports/daily
func (h *ReportAdminHandler) GetDailyReport(c *fiber.Ctx) error {
    dateStr := c.Query("date")
    var date time.Time
    var err error
    
    if dateStr == "" {
        date = time.Now()
    } else {
        date, err = time.Parse("2006-01-02", dateStr)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
                Error:   "Validation Error",
                Message: "Invalid date format. Use YYYY-MM-DD",
            })
        }
    }
    
    report, err := h.adminService.GetDailyReport(c.Context(), date)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Generate Report",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    report,
    })
}

// GetMonthlyReport - GET /admin/reports/monthly
func (h *ReportAdminHandler) GetMonthlyReport(c *fiber.Ctx) error {
    year := c.QueryInt("year")
    month := c.QueryInt("month")
    
    if year == 0 {
        year = time.Now().Year()
    }
    if month == 0 {
        month = int(time.Now().Month())
    }
    
    report, err := h.adminService.GetMonthlyReport(c.Context(), year, month)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Generate Report",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    report,
    })
}

// GetRevenueByBillType - GET /admin/reports/revenue/by-bill-type
func (h *ReportAdminHandler) GetRevenueByBillType(c *fiber.Ctx) error {
    startDate := c.Query("start_date")
    endDate := c.Query("end_date")
    
    revenue, err := h.adminService.GetRevenueByBillType(c.Context(), startDate, endDate)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Revenue",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    revenue,
    })
}

// GetTopUsers - GET /admin/reports/top-users
func (h *ReportAdminHandler) GetTopUsers(c *fiber.Ctx) error {
    limit := c.QueryInt("limit", 10)
    
    users, err := h.adminService.GetTopUsers(c.Context(), limit)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Top Users",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    users,
    })
}

// GetFraudStats - GET /admin/reports/fraud-attempts (Super Admin only)
func (h *ReportAdminHandler) GetFraudStats(c *fiber.Ctx) error {
    stats, err := h.adminService.GetFraudStats(c.Context())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Fraud Stats",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    stats,
    })
}

// GetProviderPerformance - GET /admin/reports/provider-performance
func (h *ReportAdminHandler) GetProviderPerformance(c *fiber.Ctx) error {
    startDate := c.Query("start_date")
    endDate := c.Query("end_date")
    
    performance, err := h.adminService.GetProviderPerformance(c.Context(), startDate, endDate)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Failed to Fetch Performance",
            Message: err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
        Success: true,
        Data:    performance,
    })
}

// ExportReport - POST /admin/reports/export (Super Admin only)
func (h *ReportAdminHandler) ExportReport(c *fiber.Ctx) error {
    var req struct {
        ReportType string `json:"report_type"`
        Format     string `json:"format"`
        StartDate  string `json:"start_date"`
        EndDate    string `json:"end_date"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error:   "Validation Error",
            Message: "Invalid request body",
        })
    }
    
    data, format, err := h.adminService.ExportReport(c.Context(), req.ReportType, req.Format, req.StartDate, req.EndDate)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error:   "Export Failed",
            Message: err.Error(),
        })
    }
    
    contentType := "text/csv"
    if format == "excel" {
        contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
    }
    
    c.Set("Content-Type", contentType)
    c.Set("Content-Disposition", "attachment; filename=report."+format)
    
    return c.Send(data)
}