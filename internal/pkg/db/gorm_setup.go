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
// RunMigrationsIfNeeded checks if migrations have already been run
// 
// // Check if users table exists (indicating migrations were already run)
func RunMigrationsIfNeeded(db *gorm.DB) error {
    if db.Migrator().HasTable("users") {
        log.Println("Tables already exist. Skipping migrations.")
        return nil
    }
    
    log.Println("No tables found. Running migrations...")
    return RunMigrations(db)
}
    
// func RunMigrationsIfNeeded(db *gorm.DB) error {
//     // Check if migration tracking table exists
//     if !db.Migrator().HasTable("schema_migrations") {
//         log.Println("No migration tracking found, running all migrations...")
//         return RunMigrations(db)
//     }
    
//     // Get last migration version
//     var lastMigration struct {
//         Version uint
//     }
    
//     result := db.Table("schema_migrations").Select("MAX(version) as version").Scan(&lastMigration)
//     if result.Error != nil {
//         log.Printf("Failed to get last migration version: %v", result.Error)
//         log.Println("Running migrations anyway...")
//         return RunMigrations(db)
//     }
    
//     // Define expected migrations (update this as you add migrations)
//     expectedMigrations := []string{
//         "000_init_extensions",
//         "001_create_users_table",
//         "002_create_kyc_table",
//         "003_create_sessions_table",
//         "004_create_trusted_devices_table",
//         "005_create_2fa_table",
//         "006_create_otp_table",
//         "007_create_wallets_table",
//         "008_create_transactions_table",
//         "009_create_transfer_details_table",
//         "010_create_providers_table",
//         "011_create_bill_details_table",
//         "012_create_provider_logs_table",
//         "013_create_savings_goals_table",
//         "014_create_savings_contributions_table",
//         "015_create_auto_roundup_table",
//         "016_create_fee_configs_table",
//         "017_create_tier_limits_table",
//         "018_create_admin_users_table",
//         "019_create_roles_table",
//         "020_create_audit_logs_table",
//         "021_create_support_tickets_table",
//         "022_create_ticket_messages_table",
//         "023_seed_default_data",
//     }
    
//     if lastMigration.Version >= uint(len(expectedMigrations)) {
//         log.Printf("Migrations already completed (version %d). Skipping...", lastMigration.Version)
//         return nil
//     }
    
//     log.Printf("Migrations incomplete. Last version: %d, Expected: %d. Running missing migrations...", 
//         lastMigration.Version, len(expectedMigrations))
//     return RunMigrations(db)
// }

// // CreateMigrationTable creates the schema_migrations table if it doesn't exist
// func CreateMigrationTable(db *gorm.DB) error {
//     if !db.Migrator().HasTable("schema_migrations") {
//         return db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
//             version BIGINT PRIMARY KEY,
//             dirty BOOLEAN NOT NULL DEFAULT FALSE,
//             executed_at TIMESTAMP NOT NULL DEFAULT NOW()
//         )`).Error
//     }
//     return nil
// }

// // RecordMigration records a successful migration
// func RecordMigration(db *gorm.DB, version uint) error {
//     return db.Exec(`INSERT INTO schema_migrations (version, dirty, executed_at) 
//         VALUES (?, FALSE, NOW()) 
//         ON CONFLICT (version) DO UPDATE SET dirty = FALSE, executed_at = NOW()`, version).Error
// }