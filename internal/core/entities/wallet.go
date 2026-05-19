package entities

import (
    "time"
    "github.com/google/uuid"
)

type Wallet struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID          uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_wallet_user;not null"`
    Balance         int64          `gorm:"not null;default:0"`          // IN KOBO
    LedgerBalance   int64          `gorm:"not null;default:0"`          // IN KOBO (for reconciliation)
    Currency        string         `gorm:"size:3;default:'NGN'"`
    IsLocked        bool           `gorm:"default:false"`
    LockedAt        *time.Time
    LockReason      string         `gorm:"size:255"`
    DailySpent      int64          `gorm:"default:0"`                   // IN KOBO
    WeeklySpent     int64          `gorm:"default:0"`                   // IN KOBO
    MonthlySpent    int64          `gorm:"default:0"`                   // IN KOBO
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
    Amount          int64          `gorm:"not null"`                                  // IN KOBO
    Fee             int64          `gorm:"default:0"`                                 // IN KOBO
    VAT             int64          `gorm:"default:0"`                                 // IN KOBO
    TotalAmount     int64          `gorm:"not null"`                                  // Amount + Fee + VAT (IN KOBO)
    BalanceBefore   int64          `gorm:"not null"`                                  // IN KOBO
    BalanceAfter    int64          `gorm:"not null"`                                  // IN KOBO
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