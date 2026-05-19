package db

import (
    "fmt"
    "log"
    "fintech-api/internal/core/entities"
    "gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
    log.Println("Running database migrations...")
    
    // Order matters for foreign keys
    models := []interface{}{
        &entities.User{},
        &entities.KYC{},
        &entities.Session{},
        &entities.TrustedDevice{},
        &entities.TwoFA{},
        &entities.OTP{},
        &entities.Wallet{},
        &entities.Transaction{},
        &entities.TransferDetail{},
        &entities.Provider{},
        &entities.BillDetail{},
        &entities.ProviderLog{},
        &entities.SavingsGoal{},
        &entities.SavingsContribution{},
        &entities.AutoRoundup{},
        &entities.FeeConfig{},
        &entities.TierLimit{},
        &entities.AdminUser{},
        &entities.Role{},
        &entities.AuditLog{},
        &entities.SupportTicket{},
        &entities.TicketMessage{},
    }
    
    for _, model := range models {
        if err := db.AutoMigrate(model); err != nil {
            return fmt.Errorf("failed to migrate %T: %w", model, err)
        }
        log.Printf("✓ Migrated %T", model)
    }
    
    // Create indexes for better performance
    if err := createIndexes(db); err != nil {
        return err
    }
    
    log.Println("All migrations completed successfully!")
    return nil
}

func createIndexes(db *gorm.DB) error {
    // Composite indexes
    queries := []string{
        "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_transactions_user_status ON transactions(user_id, status, created_at DESC);",
        "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_wallets_user_limit ON wallets(user_id, is_locked);",
        "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_savings_user_status ON savings_goals(user_id, status);",
        "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_provider_type_active ON providers(type, is_active, priority);",
    }
    
    for _, query := range queries {
        if err := db.Exec(query).Error; err != nil {
            log.Printf("Warning: Could not create index: %v", err)
        }
    }
    
    return nil
}