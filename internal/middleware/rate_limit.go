package middleware

import (
    "sync"
    "time"
    
    "github.com/gofiber/fiber/v2"
    "github.com/redis/go-redis/v9"
    "context"
    "fmt"
)

type RateLimiter struct {
    redisClient *redis.Client
    requests    int           // requests per duration
    duration    time.Duration // time window
    mu          sync.RWMutex
}

type InMemoryRateLimiter struct {
    requests   map[string][]time.Time
    maxRequests int
    window      time.Duration
    mu          sync.RWMutex
}

func NewRateLimiter(redisClient *redis.Client, requests int, duration time.Duration) *RateLimiter {
    return &RateLimiter{
        redisClient: redisClient,
        requests:    requests,
        duration:    duration,
    }
}

// RateLimit implements rate limiting using Redis
func (r *RateLimiter) RateLimit() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Get identifier (IP + user ID if authenticated)
        identifier := c.IP()
        
        // If user is authenticated, include user ID for better tracking
        userID := c.Locals("user_id")
        if userID != nil {
            identifier = fmt.Sprintf("%s:%v", identifier, userID)
        }
        
        ctx := context.Background()
        key := fmt.Sprintf("rate_limit:%s", identifier)
        
        // Get current count
        count, err := r.redisClient.Get(ctx, key).Int()
        if err != nil && err != redis.Nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Rate limit check failed",
            })
        }
        
        if count >= r.requests {
            ttl, _ := r.redisClient.TTL(ctx, key).Result()
            return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
                "error":      "Rate limit exceeded",
                "retry_after": int(ttl.Seconds()),
                "limit":      r.requests,
                "window":     r.duration.String(),
            })
        }
        
        // Increment count
        pipe := r.redisClient.Pipeline()
        pipe.Incr(ctx, key)
        pipe.Expire(ctx, key, r.duration)
        _, err = pipe.Exec(ctx)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Failed to update rate limit",
            })
        }
        
        // Add rate limit headers
        c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", r.requests))
        c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", r.requests-count-1))
        c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(r.duration).Unix()))
        
        return c.Next()
    }
}

// NewInMemoryRateLimiter creates an in-memory rate limiter (for testing or simple deployments)
func NewInMemoryRateLimiter(maxRequests int, window time.Duration) *InMemoryRateLimiter {
    // Start cleanup goroutine
    limiter := &InMemoryRateLimiter{
        requests:    make(map[string][]time.Time),
        maxRequests: maxRequests,
        window:      window,
    }
    
    go limiter.cleanup()
    
    return limiter
}

func (l *InMemoryRateLimiter) cleanup() {
    ticker := time.NewTicker(time.Minute)
    for range ticker.C {
        l.mu.Lock()
        now := time.Now()
        for key, timestamps := range l.requests {
            valid := make([]time.Time, 0)
            for _, ts := range timestamps {
                if now.Sub(ts) < l.window {
                    valid = append(valid, ts)
                }
            }
            if len(valid) == 0 {
                delete(l.requests, key)
            } else {
                l.requests[key] = valid
            }
        }
        l.mu.Unlock()
    }
}

// RateLimitInMemory implements in-memory rate limiting
func (l *InMemoryRateLimiter) RateLimitInMemory() fiber.Handler {
    return func(c *fiber.Ctx) error {
        identifier := c.IP()
        userID := c.Locals("user_id")
        if userID != nil {
            identifier = fmt.Sprintf("%s:%v", identifier, userID)
        }
        
        l.mu.Lock()
        defer l.mu.Unlock()
        
        now := time.Now()
        timestamps := l.requests[identifier]
        
        // Remove old entries
        valid := make([]time.Time, 0)
        for _, ts := range timestamps {
            if now.Sub(ts) < l.window {
                valid = append(valid, ts)
            }
        }
        
        if len(valid) >= l.maxRequests {
            oldest := valid[0]
            retryAfter := l.window - now.Sub(oldest)
            return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
                "error":       "Rate limit exceeded",
                "retry_after": int(retryAfter.Seconds()),
                "limit":       l.maxRequests,
                "window":      l.window.String(),
            })
        }
        
        valid = append(valid, now)
        l.requests[identifier] = valid
        
        c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", l.maxRequests))
        c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", l.maxRequests-len(valid)))
        c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", now.Add(l.window).Unix()))
        
        return c.Next()
    }
}