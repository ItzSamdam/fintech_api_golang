package db

import (
    "fmt"
    "time"
    "fintech_api_golang/internal/core/entities"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func InitGORM(cfg *config.DatabaseConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Africa/Lagos",
        cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    })
    
    if err != nil {
        return nil, err
    }
    
    // Auto migrate (in development)
    err = db.AutoMigrate(
        &entities.User{},
        &entities.KYC{},
        &entities.Session{},
        &entities.TrustedDevice{},
        &entities.TwoFA{},
        &entities.OTP{},
        &entities.Wallet{},
        &entities.Transaction{},
        &entities.TransferDetail{},
        &entities.BillDetail{},
        &entities.Provider{},
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
    )
    
    return db, err
}

// GORM Hooks for auto-converting NGN to Kobo
type AmountInKobo int64

func (a *AmountInKobo) Scan(value interface{}) error {
    // Convert from DB (int64) to struct
    *a = AmountInKobo(value.(int64))
    return nil
}

// BeforeCreate hook for generating UUID if not set
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return nil
}

// AfterFind hook for formatting (if needed)
func (t *Transaction) AfterFind(tx *gorm.DB) error {
    // Any post-load transformations
    return nil
}