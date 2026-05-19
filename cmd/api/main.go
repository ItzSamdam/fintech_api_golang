package main

import (
    "log"
    "fintech_api_golang/internal/config"
    "fintech_api_golang/internal/pkg/db"
    "fintech_api_golang/api/routes"
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