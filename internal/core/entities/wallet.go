package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallet struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID          uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_wallet_user;not null"`
    // Balance         int64          `gorm:"not null;default:0"`          // IN KOBO
    Balance AmountInKobo `gorm:"type:bigint;not null;default:0"`
    // LedgerBalance   int64          `gorm:"not null;default:0"`          // IN KOBO (for reconciliation)
    LedgerBalance AmountInKobo `gorm:"type:bigint;not null;default:0"`
    Currency        string         `gorm:"size:3;default:'NGN'"`
    IsLocked        bool           `gorm:"default:false"`
    LockedAt        *time.Time
    LockReason      string         `gorm:"size:255"`
    // DailySpent      int64          `gorm:"default:0"`                   // IN KOBO
    DailySpent AmountInKobo `gorm:"type:bigint;not null;default:0"`
    // WeeklySpent     int64          `gorm:"default:0"`                   // IN KOBO
    WeeklySpent AmountInKobo `gorm:"type:bigint;not null;default:0"`
    // MonthlySpent    int64          `gorm:"default:0"`                   // IN KOBO
    MonthlySpent AmountInKobo `gorm:"type:bigint;not null;default:0"`
    LastDailyReset  time.Time      `gorm:"not null;default:now()"`
    LastWeeklyReset time.Time      `gorm:"not null;default:now()"`
    LastMonthlyReset time.Time     `gorm:"not null;default:now()"`
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    UpdatedAt       time.Time      `gorm:"not null;default:now()"`
    
    // Relationships
    User            User           `gorm:"foreignKey:UserID"`
    Transactions    []Transaction  `gorm:"foreignKey:WalletID"`
}

type Transaction struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Reference       string         `gorm:"uniqueIndex:idx_txn_ref;size:100;not null"` // Unique transaction reference
    WalletID        uuid.UUID      `gorm:"type:uuid;index:idx_txn_wallet;not null"`
    UserID          uuid.UUID      `gorm:"type:uuid;index:idx_txn_user;not null"`
    Type            string         `gorm:"size:50;index:idx_txn_type;not null"`       // credit, debit
    Category        string         `gorm:"size:50;index:idx_txn_category;not null"`   // transfer, airtime, data, electricity, betting, savings
    SubCategory     string         `gorm:"size:50"`                                    // mtn, glo, ikeja_electric, bet9ja
    // Amount          int64          `gorm:"not null"`                                  // IN KOBO
    Amount AmountInKobo `gorm:"type:bigint;not null"`
    // Fee             int64          `gorm:"default:0"`                                 // IN KOBO
    Fee AmountInKobo `gorm:"type:bigint;default:0"`
    // VAT             int64          `gorm:"default:0"`                                 // IN KOBO
    VAT AmountInKobo `gorm:"type:bigint;default:0"`
    // TotalAmount     int64          `gorm:"not null"`                                  // Amount + Fee + VAT (IN KOBO)
    TotalAmount AmountInKobo `gorm:"type:bigint;not null"`
    // BalanceBefore   int64          `gorm:"not null"`                                  // IN KOBO
    BalanceBefore AmountInKobo `gorm:"type:bigint;not null"`
    // BalanceAfter    int64          `gorm:"not null"`                                  // IN KOBO
    BalanceAfter AmountInKobo `gorm:"type:bigint;not null"`
    Status          string         `gorm:"size:20;index:idx_txn_status;default:'pending'"` // pending, success, failed, reversed
    Description     string         `gorm:"size:255"`
    Metadata        string         `gorm:"type:jsonb"`                                 // JSON metadata
    ProviderReference string       `gorm:"size:100"`                                  // Third-party reference
    ProviderResponse string        `gorm:"type:text"`                                 // Raw provider response
    RetryCount      int            `gorm:"default:0"`
    IsReversed      bool           `gorm:"default:false"`
    ReversedTxnID   *uuid.UUID     `gorm:"type:uuid"`
    IPAddress       string         `gorm:"size:45"`
    DeviceID        string         `gorm:"size:255"`
    CompletedAt     *time.Time
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    UpdatedAt       time.Time      `gorm:"not null;default:now()"`
    
    // Relationships
    Wallet          Wallet         `gorm:"foreignKey:WalletID"`
    User            User           `gorm:"foreignKey:UserID"`
    TransferDetail  *TransferDetail `gorm:"foreignKey:TransactionID"`
    BillDetail      *BillDetail     `gorm:"foreignKey:TransactionID"`
}

type TransferDetail struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    TransactionID   uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_transfer_txn;not null"`
    RecipientType   string         `gorm:"size:20;not null"` // bank, wallet
    RecipientID     string         `gorm:"size:100;not null"` // account number or wallet ID
    RecipientName   string         `gorm:"size:255"`
    RecipientBankCode string        `gorm:"size:10"`
    RecipientBankName string        `gorm:"size:100"`
    NIPSessionID    string         `gorm:"size:100"`         // NIP session reference
    Narration       string         `gorm:"size:255"`
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    
    // Relationships
    Transaction     Transaction    `gorm:"foreignKey:TransactionID"`
}

// AfterFind GORM hook - called after loading from database
func (t *Transaction) AfterFind(tx *gorm.DB) error {
    // Example: You can add post-load transformations here
    // For instance, you might want to populate a cache or log access
    
    // If you need to convert amounts to a different format for API response,
    // do it in the DTO layer, not here
    
    return nil
}

// BeforeCreate GORM hook for Transaction
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
    if t.ID == uuid.Nil {
        t.ID = uuid.New()
    }
    
    // Generate reference if not set
    if t.Reference == "" {
        t.Reference = generateTransactionReference()
    }
    
    return nil
}

// Helper function to generate transaction reference
func generateTransactionReference() string {
    return "TXN" + time.Now().Format("20060102150405") + uuid.New().String()[:8]
}

// BeforeCreate hook for Wallet
func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
    if w.ID == uuid.Nil {
        w.ID = uuid.New()
    }
    return nil
}

// BeforeUpdate hook for Wallet - update timestamps
func (w *Wallet) BeforeUpdate(tx *gorm.DB) error {
    w.UpdatedAt = time.Now()
    return nil
}