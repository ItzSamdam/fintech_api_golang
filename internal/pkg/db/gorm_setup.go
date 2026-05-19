package db

import (
    "fmt"
    "log"
    "time"
    
    "fintech_api_golang/internal/config"
    
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// InitGORM initializes the GORM database connection
func InitGORM(cfg *config.DatabaseConfig) (*gorm.DB, error) {
    // Build DSN string
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Africa/Lagos",
        cfg.Host,
        cfg.User,
        cfg.Password,
        cfg.DBName,
        cfg.Port,
        cfg.SSLMode,
    )
    
    // Configure GORM logger
    gormLogger := logger.Default.LogMode(logger.Info)
    if cfg.SSLMode == "disable" {
        // In development, log all queries
        gormLogger = logger.Default.LogMode(logger.Info)
    } else {
        // In production, only log errors
        gormLogger = logger.Default.LogMode(logger.Error)
    }
    
    // Open connection
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger:                 gormLogger,
        SkipDefaultTransaction: true, // For better performance
        PrepareStmt:            true, // For prepared statements
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    })
    
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    
    // Get underlying sql.DB
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
    }
    
    // Configure connection pool
    sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
    sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
    sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
    
    // Test connection
    if err := sqlDB.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    log.Println("✓ Database connection established")
    log.Printf("  Host: %s:%d", cfg.Host, cfg.Port)
    log.Printf("  Database: %s", cfg.DBName)
    log.Printf("  Max Open Conns: %d", cfg.MaxOpenConns)
    log.Printf("  Max Idle Conns: %d", cfg.MaxIdleConns)
    
    return db, nil
}

// CloseDB closes the database connection
func CloseDB(db *gorm.DB) error {
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}

// TransactionWithRetry executes a transaction with retry logic
func TransactionWithRetry(db *gorm.DB, retries int, txFunc func(*gorm.DB) error) error {
    var err error
    
    for i := 0; i < retries; i++ {
        err = db.Transaction(txFunc)
        if err == nil {
            return nil
        }
        
        // Check if retryable error
        if isRetryableError(err) && i < retries-1 {
            log.Printf("Transaction failed, retrying (%d/%d): %v", i+1, retries, err)
            time.Sleep(time.Duration(i+1) * 100 * time.Millisecond) // Backoff
            continue
        }
        break
    }
    
    return err
}

// isRetryableError checks if the error is retryable
func isRetryableError(err error) bool {
    if err == nil {
        return false
    }
    
    // Check for common retryable errors
    errMsg := err.Error()
    retryablePatterns := []string{
        "deadlock detected",
        "lock timeout",
        "connection reset",
        "connection refused",
        "too many connections",
    }
    
    for _, pattern := range retryablePatterns {
        if contains(errMsg, pattern) {
            return true
        }
    }
    
    return false
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
    return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
        (len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
    for i := 0; i <= len(s)-len(substr); i++ {
        if s[i:i+len(substr)] == substr {
            return true
        }
    }
    return false
}