package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/compress"
)

// Compression adds response compression
func Compression() fiber.Handler {
    return compress.New(compress.Config{
        Level: compress.LevelBestSpeed, // or compress.LevelBestCompression
    })
}