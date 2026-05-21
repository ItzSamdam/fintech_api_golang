package routes

import (
    "time"
    
    "github.com/gofiber/fiber/v2"
    // "github.com/gofiber/fiber/v2/middleware/recover"
    "go.uber.org/zap"
    "gorm.io/gorm"

    redisLib "github.com/redis/go-redis/v9"
    "fintech_api_golang/internal/repositories/redis"
    
    "fintech_api_golang/internal/config"
    "fintech_api_golang/internal/handlers"
    "fintech_api_golang/internal/middleware"
    "fintech_api_golang/internal/core/services"
    "fintech_api_golang/internal/handlers/backoffice"
    "fintech_api_golang/internal/repositories/postgres"
    "fintech_api_golang/internal/repositories/providers"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, rdb *redisLib.Client, logger *zap.Logger, cfg *config.Config) {
    // ========== INITIALIZE REPOSITORIES ==========
    // PostgreSQL Repositories
    userRepo := postgres.NewUserRepository(db)
    walletRepo := postgres.NewWalletRepository(db)
    transactionRepo := postgres.NewTransactionRepository(db)
    transferDetailRepo := postgres.NewTransferDetailRepository(db)
    kycRepo := postgres.NewKYCRepository(db)
    sessionRepo := postgres.NewSessionRepository(db)
    otpRepo := postgres.NewOTPRepository(db)
    providerRepo := postgres.NewProviderRepository(db)
    billDetailRepo := postgres.NewBillDetailRepository(db)
    savingsGoalRepo := postgres.NewSavingsGoalRepository(db)
    savingsContributionRepo := postgres.NewSavingsContributionRepository(db)
    roundupRepo := postgres.NewAutoRoundupRepository(db)
    auditLogRepo := postgres.NewAuditLogRepository(db)
    adminUserRepo := postgres.NewAdminUserRepository(db)
    roleRepo := postgres.NewRoleRepository(db)
    
    // Redis Repositories
    cacheRepo := redis.NewCacheRepository(rdb)
    _ = cacheRepo // Use as needed
    
    // Provider Clients
    redBiller := providers.NewRedBillerClient(
        cfg.ExternalAPIs.NIPBaseURL,
        cfg.ExternalAPIs.NIPAPIKey,
    )
    providerRegistry := providers.NewProviderRegistry(redBiller)
    _ = providerRegistry
    
    // ========== INITIALIZE SERVICES ==========
    authService := services.NewAuthService(
        userRepo, kycRepo, sessionRepo, otpRepo, walletRepo, cfg,
    )
    
    walletService := services.NewWalletService(
        walletRepo, transactionRepo, userRepo, db,
    )
    
    transferService := services.NewTransferService(
        walletRepo, transactionRepo, transferDetailRepo, userRepo, redBiller, db, cfg,
    )
    
    billPaymentService := services.NewBillPaymentService(
        walletRepo, transactionRepo, billDetailRepo, providerRepo, userRepo, redBiller, db,
    )
    
    savingsService := services.NewSavingsService(
        savingsGoalRepo, savingsContributionRepo, roundupRepo, walletRepo, transactionRepo, userRepo, db,
    )
    
    adminService := services.NewAdminService(
        userRepo, walletRepo, transactionRepo, kycRepo, providerRepo, auditLogRepo, adminUserRepo, roleRepo, db,
    )
    
    // ========== INITIALIZE HANDLERS ==========
    authHandler := handlers.NewAuthHandler(authService)
    walletHandler := handlers.NewWalletHandler(walletService)
    transferHandler := handlers.NewTransferHandler(transferService)
    billHandler := handlers.NewBillHandler(billPaymentService)
    savingsHandler := handlers.NewSavingsHandler(savingsService)
    
    // Admin Handlers
    userAdminHandler := backoffice.NewUserAdminHandler(adminService)
    transactionAdminHandler := backoffice.NewTransactionAdminHandler(adminService)
    walletAdminHandler := backoffice.NewWalletAdminHandler(adminService)
    kycAdminHandler := backoffice.NewKYCAdminHandler(adminService)
    providerAdminHandler := backoffice.NewProviderAdminHandler(adminService)
    reportAdminHandler := backoffice.NewReportAdminHandler(adminService)
    systemAdminHandler := backoffice.NewSystemAdminHandler(adminService)
    roleAdminHandler := backoffice.NewRoleAdminHandler(adminService)
    
    // ========== INITIALIZE MIDDLEWARES ==========
    userAuthMiddleware := middleware.NewUserAuthMiddleware(db, cfg.JWT.Secret)
    adminAuthMiddleware := middleware.NewAdminAuthMiddleware(db, cfg.JWT.Secret)
    rateLimiter := middleware.NewRateLimiter(rdb, cfg.RateLimit.RequestsPerMinute, time.Minute)
    tierLimiter := middleware.NewTierLimiter(db)
    
    // ========== GLOBAL MIDDLEWARE ==========
    app.Use(middleware.RequestID())
    app.Use(middleware.NewLoggerMiddleware(logger).LogRequests())
    app.Use(middleware.NewRecoveryMiddleware(logger).Recover())
    app.Use(middleware.CorsConfig())
    app.Use(middleware.SecurityHeaders())
    app.Use(middleware.Compression())
    app.Use(middleware.Helmet())

    
    // ========== SERVER CHECK ==========
    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "ok",
            "time":   time.Now().Unix(),
            "data": fiber.Map{
                "version":     "1.0.0",
                "environment": cfg.Environment,
                "port":        cfg.Port,
                "debug":       cfg.Debug,
                "message":     "Server is running",
            },
        })
    })

    // In SetupRoutes function, add:
    app.Static("/docs", "./api/docs")

    app.Get("/swagger", func(c *fiber.Ctx) error {
        c.Set("Content-Security-Policy", "")
        c.Set("X-Content-Type-Options", "")
        c.Set("X-Frame-Options", "")
        
        html := `<!DOCTYPE html>
    <html>
    <head>
        <title>Fintech API Documentation</title>
        <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
    </head>
    <body>
        <div id="swagger-ui"></div>
        <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
        <script>
            window.onload = function() {
                SwaggerUIBundle({
                    url: "/docs/swagger.yaml",
                    dom_id: '#swagger-ui',
                    presets: [SwaggerUIBundle.presets.apis],
                    layout: "BaseLayout",
                    deepLinking: true
                });
            }
        </script>
    </body>
    </html>`
        return c.Type("html").Send([]byte(html))
    })

    
    // ========== HEALTH CHECK ==========
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "ok",
            "time":   time.Now().Unix(),
        })
    })
    
    // ========== API v1 GROUP ==========
    v1 := app.Group("/api/v1", rateLimiter.RateLimit())
    
    // ========== PUBLIC ROUTES (No Auth) ==========
    auth := v1.Group("/auth")
    {
        auth.Post("/register/phone", authHandler.RegisterPhone)
        auth.Post("/verify/otp", authHandler.VerifyOTP)
        auth.Post("/register/bvn", authHandler.RegisterBVN)
        auth.Post("/verify/face", authHandler.VerifyFace)
        auth.Post("/login", authHandler.Login)
        auth.Post("/reset-password", authHandler.ResetPassword)
        auth.Post("/refresh", authHandler.RefreshToken)
    }
    
    // Public utility routes
    v1.Get("/banks/list", transferHandler.GetBanks)
    v1.Post("/transfers/name-enquiry", transferHandler.NameEnquiry)
    
    // ========== USER ROUTES (Customer Auth Required) ==========
    userRoutes := v1.Group("", userAuthMiddleware.UserAuthRequired())
    {
        // User Profile
        userRoutes.Get("/auth/me", authHandler.GetMe)
        userRoutes.Put("/auth/me", authHandler.UpdateMe)
        userRoutes.Post("/auth/change-password", authHandler.ChangePassword)
        userRoutes.Post("/auth/logout", authHandler.Logout)
        
        // Wallet
        wallet := userRoutes.Group("/wallets")
        {
            wallet.Post("/create", walletHandler.CreateWallet)
            wallet.Get("/balance", walletHandler.GetBalance)
            wallet.Get("/transactions", walletHandler.GetTransactions)
            wallet.Get("/limits", walletHandler.GetLimits)
            wallet.Post("/lock", walletHandler.LockWallet)
            wallet.Post("/unlock", walletHandler.UnlockWallet)
            wallet.Get("/statement", walletHandler.GetStatement)
        }
        
        // Transfers
        transfers := userRoutes.Group("/transfers", tierLimiter.CheckTierLimit("send_money"))
        {
            transfers.Post("/send", transferHandler.SendTransfer)
            transfers.Post("/send-to-wallet", transferHandler.SendToWallet)
            transfers.Get("/status/:reference", transferHandler.GetTransferStatus)
            transfers.Post("/retry", transferHandler.RetryTransfer)
            transfers.Get("/history", transferHandler.GetTransferHistory)
        }
        
        // Airtime
        airtime := userRoutes.Group("/airtime", tierLimiter.CheckTierLimit("buy_airtime"))
        {
            airtime.Get("/networks", billHandler.GetAirtimeNetworks)
            airtime.Post("/purchase", billHandler.PurchaseAirtime)
            airtime.Get("/history", billHandler.GetAirtimeHistory)
        }
        
        // Data Bundle
        data := userRoutes.Group("/data", tierLimiter.CheckTierLimit("buy_data"))
        {
            data.Get("/networks", billHandler.GetDataNetworks)
            data.Get("/plans/:network", billHandler.GetDataPlans)
            data.Post("/purchase", billHandler.PurchaseData)
            data.Get("/history", billHandler.GetDataHistory)
        }
        
        // Electricity
        electricity := userRoutes.Group("/electricity", tierLimiter.CheckTierLimit("pay_electricity"))
        {
            electricity.Get("/providers", billHandler.GetElectricityProviders)
            electricity.Post("/validate-meter", billHandler.ValidateMeter)
            electricity.Post("/pay-prepaid", billHandler.PayElectricity)
            electricity.Post("/pay-postpaid", billHandler.PayElectricity)
            electricity.Get("/token/:transaction_id", billHandler.GetElectricityToken)
            electricity.Get("/history", billHandler.GetElectricityHistory)
        }
        
        // Betting
        betting := userRoutes.Group("/betting", tierLimiter.CheckTierLimit("fund_betting"))
        {
            betting.Get("/providers", billHandler.GetBettingProviders)
            betting.Post("/validate-account", billHandler.ValidateBettingAccount)
            betting.Post("/fund", billHandler.FundBettingWallet)
            betting.Get("/history", billHandler.GetBettingHistory)
        }
        
        // Savings
        savings := userRoutes.Group("/savings", tierLimiter.CheckTierLimit("save"))
        {
            savings.Post("/goals/create", savingsHandler.CreateGoal)
            savings.Post("/goals/contribute", savingsHandler.ContributeToGoal)
            savings.Get("/goals", savingsHandler.GetGoals)
            savings.Get("/goals/:id", savingsHandler.GetGoal)
            savings.Put("/goals/:id", savingsHandler.UpdateGoal)
            savings.Delete("/goals/:id", savingsHandler.DeleteGoal)
            savings.Post("/roundup/activate", savingsHandler.ActivateRoundup)
            savings.Post("/roundup/deactivate", savingsHandler.DeactivateRoundup)
            savings.Get("/roundup/status", savingsHandler.GetRoundupStatus)
        }
        
        // Transactions
        userRoutes.Get("/transactions", walletHandler.GetTransactions)
        userRoutes.Get("/transactions/:id", walletHandler.GetTransactionByID)
        userRoutes.Get("/bills/history", billHandler.GetBillHistory)
        
        // Compliance & Security
        compliance := userRoutes.Group("/compliance")
        {
            compliance.Post("/report/suspicious", nil) // Add compliance handler
            compliance.Get("/limits/check", nil)
        }
        
        security := userRoutes.Group("/security")
        {
            security.Post("/sim-swap/check", nil)
            security.Post("/device/trust", nil)
            security.Post("/2fa/enable", nil)
            security.Post("/2fa/verify", nil)
            security.Get("/sessions", nil)
            security.Delete("/sessions/:id", nil)
        }
        
        // Notifications
        userRoutes.Get("/notifications/in-app", nil)
        userRoutes.Put("/notifications/:id/read", nil)
        
        // Support
        support := userRoutes.Group("/support")
        {
            support.Post("/tickets/create", nil)
            support.Get("/tickets", nil)
            support.Get("/tickets/:id", nil)
            support.Post("/tickets/:id/reply", nil)
            support.Put("/tickets/:id/status", nil)
        }
    }
    
    // ========== ADMIN ROUTES ==========
    adminRoutes := v1.Group("/admin", adminAuthMiddleware.AdminAuthRequired())
    {
        // User Management
        users := adminRoutes.Group("/users")
        {
            users.Get("/", userAdminHandler.ListUsers)
            users.Get("/:id", userAdminHandler.GetUserDetails)
            users.Post("/:id/tier/upgrade", userAdminHandler.UpgradeUserTier)
            users.Post("/:id/suspend", userAdminHandler.SuspendUser)
            users.Post("/:id/unsuspend", userAdminHandler.UnsuspendUser)
            users.Delete("/:id", userAdminHandler.DeleteUser)
            users.Put("/:id/limits", userAdminHandler.OverrideLimits)
            users.Get("/search", userAdminHandler.SearchUsers)
        }
        
        // Transaction Management
        transactions := adminRoutes.Group("/transactions")
        {
            transactions.Get("/", transactionAdminHandler.ListTransactions)
            transactions.Get("/:id", transactionAdminHandler.GetTransactionDetails)
            transactions.Post("/reverse", transactionAdminHandler.ReverseTransaction)
            transactions.Post("/void", transactionAdminHandler.VoidTransaction)
            transactions.Get("/summary", transactionAdminHandler.GetTransactionSummary)
        }
        
        // Wallet Management
        wallets := adminRoutes.Group("/wallets")
        {
            wallets.Get("/", walletAdminHandler.ListWallets)
            wallets.Get("/:id", walletAdminHandler.GetWalletDetails)
            wallets.Post("/credit", walletAdminHandler.CreditWallet)
            wallets.Post("/debit", walletAdminHandler.DebitWallet)
            wallets.Post("/freeze", walletAdminHandler.FreezeWallet)
            wallets.Post("/unfreeze", walletAdminHandler.UnfreezeWallet)
            wallets.Get("/balances/summary", walletAdminHandler.GetBalanceSummary)
        }
        
        // KYC Management
        kyc := adminRoutes.Group("/kyc")
        {
            kyc.Get("/pending", kycAdminHandler.ListPendingKYC)
            kyc.Get("/:id", kycAdminHandler.GetKYCDetails)
            kyc.Post("/:id/approve", kycAdminHandler.ApproveKYC)
            kyc.Post("/:id/reject", kycAdminHandler.RejectKYC)
        }
        
        // Provider Management
        providers := adminRoutes.Group("/providers")
        {
            providers.Get("/", providerAdminHandler.ListProviders)
            providers.Put("/:id/toggle", providerAdminHandler.ToggleProvider)
            providers.Put("/:id/priority", providerAdminHandler.SetProviderPriority)
            providers.Get("/:id/health", providerAdminHandler.CheckProviderHealth)
            providers.Get("/logs", providerAdminHandler.GetProviderLogs)
        }
        
        // Fee Management
        fees := adminRoutes.Group("/fees")
        {
            fees.Get("/", nil)
            fees.Put("/:bill_type", nil)
            fees.Get("/margins", nil)
            fees.Put("/margins/:provider_id", nil)
        }
        
        // Reports
        reports := adminRoutes.Group("/reports")
        {
            reports.Get("/daily", reportAdminHandler.GetDailyReport)
            reports.Get("/monthly", reportAdminHandler.GetMonthlyReport)
            reports.Get("/revenue/by-bill-type", reportAdminHandler.GetRevenueByBillType)
            reports.Get("/top-users", reportAdminHandler.GetTopUsers)
            reports.Get("/fraud-attempts", reportAdminHandler.GetFraudStats)
            reports.Get("/provider-performance", reportAdminHandler.GetProviderPerformance)
            reports.Post("/export", reportAdminHandler.ExportReport)
        }
        
        // System Settings (Super Admin only)
        settings := adminRoutes.Group("/settings", adminAuthMiddleware.RequireAdminRole("super_admin"))
        {
            settings.Get("/", systemAdminHandler.GetSystemSettings)
            settings.Put("/", systemAdminHandler.UpdateSystemSettings)
            settings.Get("/health", systemAdminHandler.HealthCheck)
            settings.Get("/audit-logs", systemAdminHandler.ListAuditLogs)
            settings.Get("/audit-logs/:id", systemAdminHandler.GetAuditLog)
            settings.Post("/backup/database", systemAdminHandler.TriggerDatabaseBackup)
            settings.Get("/metrics", systemAdminHandler.GetSystemMetrics)
        }
        
        // Role Management (Super Admin only)
        roles := adminRoutes.Group("/roles", adminAuthMiddleware.RequireAdminRole("super_admin"))
        {
            roles.Get("/", roleAdminHandler.ListRoles)
            roles.Post("/", roleAdminHandler.CreateRole)
            roles.Put("/:id", roleAdminHandler.UpdateRole)
            roles.Delete("/:id", roleAdminHandler.DeleteRole)
            roles.Get("/staff", roleAdminHandler.ListStaff)
            roles.Post("/staff/invite", roleAdminHandler.InviteStaff)
            roles.Put("/staff/:id/role", roleAdminHandler.AssignRole)
            roles.Delete("/staff/:id", roleAdminHandler.RemoveStaff)
            roles.Get("/staff/:id/audit", roleAdminHandler.GetStaffAudit)
        }
    }
    
    // ========== 404 HANDLER ==========
    app.Use(func(c *fiber.Ctx) error {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error":   "Not Found",
            "message": "Route not found",
        })
    })
}