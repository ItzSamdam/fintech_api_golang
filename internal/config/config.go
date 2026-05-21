package config

import (
    "log"
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
)

type Config struct {
    // Server
    Port         string
    Environment  string
    Debug        bool
    // Database
    Database     DatabaseConfig
    // Redis
    Redis        RedisConfig
    // JWT
    JWT          JWTConfig
    // Encryption
    Encryption   EncryptionConfig
    // External APIs
    ExternalAPIs ExternalAPIConfig
    // Rate Limiting
    RateLimit    RateLimitConfig
    // Queue (RabbitMQ)
    Queue        QueueConfig
    // Email/SMS
    SMTP         SMTPConfig
    SMS          SMSConfig
    // Cloud Storage
    Cloudinary   CloudinaryConfig
    // Webhook
    Webhook      WebhookConfig
    // Face Recognition
    FaceRecognition FaceRecognitionConfig
}

type DatabaseConfig struct {
    Host        string
    Port        int
    User        string
    Password    string
    DBName      string
    SSLMode     string
    MaxOpenConns int
    MaxIdleConns int
    ConnMaxLifetime time.Duration
}

type RedisConfig struct {
    Host     string
    Port     int
    Password string
    DB       int
    PoolSize int
}

type JWTConfig struct {
    Secret        string
    AccessExpiry  time.Duration
    RefreshExpiry time.Duration
    Issuer        string
}

type EncryptionConfig struct {
    AESKey    string
}

type ExternalAPIConfig struct {
    // Bank APIs (NIP)
    NIPBaseURL      string
    NIPAPIKey       string
    NIPHook			string
    // Face Recognition
    FaceBaseURL     string
    FaceAPIKey      string
}

type RateLimitConfig struct {
    RequestsPerMinute int
    BurstSize         int
}

type QueueConfig struct {
    RabbitMQURL string
    Exchange    string
    Queue       string
}

type SMTPConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    From     string
}

type SMSConfig struct {
    TermiiAPIKey     string
}

type CloudinaryConfig struct {
    CloudName string
    APIKey    string
    APISecret string
}

type WebhookConfig struct {
    TimeoutSeconds int
    RetryCount     int
}

type FaceRecognitionConfig struct {
    Tolerance   float64
    LivenessThreshold float64
}

// Load loads configuration from environment variables
func Load() *Config {
    // Load .env file if it exists (development)
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    return &Config{
        Port:        getEnv("PORT", "3000"),
        Environment: getEnv("ENVIRONMENT", "development"),
        Debug:       getEnvAsBool("DEBUG", false),
        
        Database: DatabaseConfig{
            Host:            getEnv("DB_HOST", "localhost"),
            Port:            getEnvAsInt("DB_PORT", 5432),
            User:            getEnv("DB_USER", "postgres"),
            Password:        getEnv("DB_PASSWORD", "postgres"),
            DBName:          getEnv("DB_NAME", "fintech_db"),
            SSLMode:         getEnv("DB_SSLMODE", "disable"),
            MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
            MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
            ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 1*time.Hour),
        },
        
        Redis: RedisConfig{
            Host:     getEnv("REDIS_HOST", "localhost"),
            Port:     getEnvAsInt("REDIS_PORT", 6379),
            Password: getEnv("REDIS_PASSWORD", ""),
            DB:       getEnvAsInt("REDIS_DB", 0),
            PoolSize: getEnvAsInt("REDIS_POOL_SIZE", 10),
        },
        
        JWT: JWTConfig{
            Secret:        getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
            AccessExpiry:  getEnvAsDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
            RefreshExpiry: getEnvAsDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
            Issuer:        getEnv("JWT_ISSUER", "fintech-api"),
        },
        
        Encryption: EncryptionConfig{
            AESKey: getEnv("AES_ENCRYPTION_KEY", "32-byte-key-for-aes-encryption!!"),
        },
        
        ExternalAPIs: ExternalAPIConfig{
            NIPBaseURL:  getEnv("NIP_BASE_URL", "https://api.nip.com.ng/v1"),
            NIPAPIKey:   getEnv("NIP_API_KEY", ""),
            NIPHook:     getEnv("NIP_API_HOOK", ""),
            FaceBaseURL: getEnv("FACE_BASE_URL", "https://api.face.com/v1"),
            FaceAPIKey:  getEnv("FACE_API_KEY", ""),
        },
        
        RateLimit: RateLimitConfig{
            RequestsPerMinute: getEnvAsInt("RATE_LIMIT_REQUESTS", 60),
            BurstSize:         getEnvAsInt("RATE_LIMIT_BURST", 10),
        },
        
        Queue: QueueConfig{
            RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
            Exchange:    getEnv("QUEUE_EXCHANGE", "fintech"),
            Queue:       getEnv("QUEUE_NAME", "transactions"),
        },
        
        SMTP: SMTPConfig{
            Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
            Port:     getEnvAsInt("SMTP_PORT", 587),
            Username: getEnv("SMTP_USERNAME", ""),
            Password: getEnv("SMTP_PASSWORD", ""),
            From:     getEnv("SMTP_FROM", "noreply@fintech.com"),
        },
        
        SMS: SMSConfig{
            TermiiAPIKey:   getEnv("TERMII_API_KEY", ""),
        },
        
        Cloudinary: CloudinaryConfig{
            CloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
            APIKey:    getEnv("CLOUDINARY_API_KEY", ""),
            APISecret: getEnv("CLOUDINARY_API_SECRET", ""),
        },
        
        Webhook: WebhookConfig{
            TimeoutSeconds: getEnvAsInt("WEBHOOK_TIMEOUT", 30),
            RetryCount:     getEnvAsInt("WEBHOOK_RETRY_COUNT", 3),
        },
        
        FaceRecognition: FaceRecognitionConfig{
            Tolerance:         getEnvAsFloat("FACE_TOLERANCE", 0.6),
            LivenessThreshold: getEnvAsFloat("FACE_LIVENESS_THRESHOLD", 0.8),
        },
    }
}

// Helper functions to get environment variables with defaults
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        intVal, err := strconv.Atoi(value)
        if err == nil {
            return intVal
        }
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        boolVal, err := strconv.ParseBool(value)
        if err == nil {
            return boolVal
        }
    }
    return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        durationVal, err := time.ParseDuration(value)
        if err == nil {
            return durationVal
        }
    }
    return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
    if value := os.Getenv(key); value != "" {
        floatVal, err := strconv.ParseFloat(value, 64)
        if err == nil {
            return floatVal
        }
    }
    return defaultValue
}