package db

import (
    "context"
    "fmt"
    "log"
    
    "github.com/redis/go-redis/v9"
    "fintech_api_golang/internal/config"
)

func InitRedis(cfg *config.RedisConfig) (*redis.Client, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
        Password: cfg.Password,
        DB:       cfg.DB,
        PoolSize: cfg.PoolSize,
    })
    
    // Test connection
    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }
    
    log.Printf("✓ Redis connection established")
    log.Printf("  Host: %s:%d", cfg.Host, cfg.Port)
    log.Printf("  DB: %d", cfg.DB)
    
    return client, nil
}

func CloseRedis(client *redis.Client) error {
    if client != nil {
        return client.Close()
    }
    return nil
}