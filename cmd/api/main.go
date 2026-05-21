package main

import (
    "log"
    "fintech_api_golang/internal/config"
    "fintech_api_golang/internal/pkg/db"
    "fintech_api_golang/internal/server"
)

func main() {
    // Load configuration
    cfg := config.Load()
    
    log.Printf("Starting Fintech API - Environment: %s", cfg.Environment)
    log.Printf("Debug mode: %v", cfg.Debug)
    
    // Initialize database
    gormDB, err := db.InitGORM(&cfg.Database)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    
    // Run migrations (only in development and if not already run)
        if cfg.Environment == "development" {
            if err := db.RunMigrationsIfNeeded(gormDB); err != nil {
                log.Fatalf("Failed to run migrations: %v", err)
            }
        }
        
    
    // Initialize Redis (if needed)
    redisClient, err := db.InitRedis(&cfg.Redis)
    if err != nil {
        log.Fatalf("Failed to initialize Redis: %v", err)
    }
    
    // Close connections on shutdown
    defer func() {
        if err := db.CloseDB(gormDB); err != nil {
            log.Printf("Error closing database: %v", err)
        }
        log.Println("Database connection closed")
    }()
    
    // Initialize and start server
    server := server.NewServer(cfg, gormDB, redisClient)
    if err := server.Start(); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}