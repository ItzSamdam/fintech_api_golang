package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CorsConfig returns CORS middleware configuration
func CorsConfig() fiber.Handler {
    return cors.New(cors.Config{
        AllowOrigins:     "*", // Configure based on your needs
        AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-ID",
        ExposeHeaders:    "Content-Length, X-Request-ID",
        AllowCredentials: true,
        MaxAge:           300, // 5 minutes
    })
}

// StrictCorsConfig for production with specific origins
func StrictCorsConfig(allowedOrigins []string) fiber.Handler {
    return cors.New(cors.Config{
        AllowOrigins:     strings.Join(allowedOrigins, ", "),
        AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-ID",
        ExposeHeaders:    "Content-Length, X-Request-ID",
        AllowCredentials: true,
        MaxAge:           300,
    })
}