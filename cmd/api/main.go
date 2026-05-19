package main

import (
    "log"
    "fintech_api_golang/internal/config"
    "fintech_api_golang/internal/pkg/db"
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
    
    // Run migrations (only in development)
    if cfg.Environment == "development" {
        if err := db.RunMigrations(gormDB); err != nil {
            log.Fatalf("Failed to run migrations: %v", err)
        }
    }
    
    // Initialize Redis (if needed)
    // redisClient := db.InitRedis(&cfg.Redis)
    
    // Close connections on shutdown
    defer func() {
        if err := db.CloseDB(gormDB); err != nil {
            log.Printf("Error closing database: %v", err)
        }
        log.Println("Database connection closed")
    }()
    
    // Initialize and start server
    // server := server.NewServer(cfg, gormDB, redisClient)
    // if err := server.Start(); err != nil {
    //     log.Fatalf("Failed to start server: %v", err)
    // }
}