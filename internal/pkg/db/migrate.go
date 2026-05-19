package db

import (
    "fmt"
    "log"
    "fintech_api_golang/internal/core/entities"
    "gorm.io/gorm"
)


// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB) error {
    log.Println("Running database migrations...")
    
    // Order matters for foreign keys
    // Child tables should come after parent tables
    models := []interface{}{
        // Parent tables (no foreign keys)
        &entities.User{},
        &entities.Provider{},
        &entities.AdminUser{},
        &entities.Role{},
        
        // Child tables (with foreign keys)
        &entities.KYC{},
        &entities.Session{},
        &entities.TrustedDevice{},
        &entities.TwoFA{},
        &entities.OTP{},
        &entities.Wallet{},
        &entities.FeeConfig{},
        &entities.TierLimit{},
        &entities.SavingsGoal{},
        &entities.AutoRoundup{},
        
        // Grandchild tables
        &entities.Transaction{},
        &entities.TransferDetail{},
        &entities.BillDetail{},
        &entities.ProviderLog{},
        &entities.SavingsContribution{},
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
    
    // Create additional indexes for better performance
    if err := createIndexes(db); err != nil {
        log.Printf("Warning: Some indexes could not be created: %v", err)
    }
    
    log.Println("All migrations completed successfully!")
    return nil
}

// createIndexes creates additional composite indexes
func createIndexes(db *gorm.DB) error {
    queries := []string{
        // Transaction indexes
        "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_transactions_user_status ON transactions(user_id, status, created_at DESC);",
        "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_transactions_category_status ON transactions(category, status, created_at DESC);",
        "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_transactions_created_date ON transactions(date(created_at));",
        
        // Wallet indexes
        "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_wallets_user_balance ON wallets(user_id, balance);",
        
        // Savings indexes
        "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_savings_user_status ON savings_goals(user_id, status);",
        "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_savings_target_date ON savings_goals(target_date) WHERE status = 'active';",
        
        // Provider indexes
        "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_provider_type_active ON providers(type, is_active, priority);",
        
        // Bill details indexes
        "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_bill_details_phone ON bill_details(phone_number) WHERE phone_number IS NOT NULL;",
        "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_bill_details_meter ON bill_details(meter_number) WHERE meter_number IS NOT NULL;",
        
        // Audit log indexes
        "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_audit_logs_action_date ON audit_logs(action, created_at DESC);",
        
        // Support ticket indexes
        "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tickets_user_status ON support_tickets(user_id, status);",
        "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tickets_assigned_status ON support_tickets(assigned_to, status);",
    }
    
    for _, query := range queries {
        if err := db.Exec(query).Error; err != nil {
            log.Printf("Warning: Could not create index: %v - %s", err, query)
            // Don't fail the entire migration, just log the warning
        }
    }
    
    return nil
}
