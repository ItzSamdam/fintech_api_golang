package routes

import (
    "github.com/gofiber/fiber/v2"
)

// RegisterCustomMiddleware registers custom middleware for specific routes
func RegisterCustomMiddleware(app *fiber.App) {
    // Route-specific middleware can be registered here
    // For example, IP whitelisting for admin routes
}