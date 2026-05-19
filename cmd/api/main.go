package main

import (
    "log"
    "fintech-api/internal/config"
    "fintech-api/internal/pkg/db"
    "fintech-api/api/routes"
    "github.com/gofiber/fiber/v2"
)

func main() {
    // Load config
    cfg := config.Load()
    
    // Connect databases
    postgres := db.ConnectPostgres(cfg)
    redis := db.ConnectRedis(cfg)
    
    // Initialize app
    app := fiber.New(fiber.Config{
        ErrorHandler: customErrorHandler,
    })
    
    // Setup routes
    routes.Setup(app, postgres, redis, cfg)
    
    // Start server
    log.Fatal(app.Listen(":" + cfg.Port))
}