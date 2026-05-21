package server

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/gofiber/fiber/v2"
    "github.com/redis/go-redis/v9"
    "go.uber.org/zap"
    "gorm.io/gorm"
    
    "fintech_api_golang/internal/config"
    "fintech_api_golang/api/routes"
)

type Server struct {
    app    *fiber.App
    cfg    *config.Config
    db     *gorm.DB
    rdb    *redis.Client
    logger *zap.Logger
}

func NewServer(cfg *config.Config, db *gorm.DB, rdb *redis.Client) *Server {
    var logger *zap.Logger
    var err error
    
    if cfg.Environment == "production" {
        logger, err = zap.NewProduction()
    } else {
        logger, err = zap.NewDevelopment()
    }
    
    if err != nil {
        log.Fatalf("Failed to create logger: %v", err)
    }
    
    app := fiber.New(fiber.Config{
        AppName:      "Fintech API",
        Prefork:      false,
        ServerHeader: "Fintech",
        ErrorHandler: customErrorHandler,
    })
    
    return &Server{
        app:    app,
        cfg:    cfg,
        db:     db,
        rdb:    rdb,
        logger: logger,
    }
}

func (s *Server) Start() error {
    // Setup routes
    routes.SetupRoutes(s.app, s.db, s.rdb, s.logger, s.cfg)
    
    // Start server in a goroutine
    addr := fmt.Sprintf(":%s", s.cfg.Port)
    
    // Handle graceful shutdown
    go func() {
        log.Printf("Server starting on %s", addr)
        log.Printf("Environment: %s", s.cfg.Environment)
        log.Printf("Press Ctrl+C to stop")
        
        if err := s.app.Listen(addr); err != nil {
            log.Fatalf("Server failed to start: %v", err)
        }
    }()
    
    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    return s.app.Shutdown()
}

func customErrorHandler(c *fiber.Ctx, err error) error {
    code := fiber.StatusInternalServerError
    
    if e, ok := err.(*fiber.Error); ok {
        code = e.Code
    }
    
    return c.Status(code).JSON(fiber.Map{
        "error":   "Server Error",
        "message": err.Error(),
    })
}