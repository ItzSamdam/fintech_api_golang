package entities

import (
    "time"
    "github.com/google/uuid"
)

type SavingsGoal struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID          uuid.UUID      `gorm:"type:uuid;index:idx_savings_user;not null"`
    Name            string         `gorm:"size:100;not null"`
    TargetAmount    int64          `gorm:"not null"`                      // IN KOBO
    CurrentAmount   int64          `gorm:"default:0"`                     // IN KOBO
    InterestRate    float64        `gorm:"type:decimal(5,2);default:0"`
    DurationDays    int            `gorm:"not null"`                      // Goal period in days
    StartDate       time.Time      `gorm:"not null;default:now()"`
    TargetDate      time.Time      `gorm:"not null"`
    IsAutoDebit     bool           `gorm:"default:false"`
    AutoDebitAmount int64          `gorm:"default:0"`                     // IN KOBO
    AutoDebitDay    int            `gorm:"default:1"`                     // Day of month
    Status          string         `gorm:"size:20;default:'active'"`      // active, completed, withdrawn, cancelled
    WithdrawnAt     *time.Time
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    UpdatedAt       time.Time      `gorm:"not null;default:now()"`
    
    // Relationships
    User            User           `gorm:"foreignKey:UserID"`
    Contributions   []SavingsContribution `gorm:"foreignKey:SavingsGoalID"`
    Transactions    []Transaction  `gorm:"foreignKey:ID"`
}

type SavingsContribution struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    SavingsGoalID   uuid.UUID      `gorm:"type:uuid;index:idx_contribution_goal;not null"`
    TransactionID   uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_contribution_txn;not null"`
    Amount          int64          `gorm:"not null"`                      // IN KOBO
    InterestEarned  int64          `gorm:"default:0"`                     // IN KOBO
    ContributionDate time.Time     `gorm:"not null;default:now()"`
    IsAutoDebit     bool           `gorm:"default:false"`
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    
    // Relationships
    SavingsGoal     SavingsGoal    `gorm:"foreignKey:SavingsGoalID"`
    Transaction     Transaction    `gorm:"foreignKey:TransactionID"`
}

type AutoRoundup struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID          uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_roundup_user;not null"`
    SavingsGoalID   uuid.UUID      `gorm:"type:uuid;not null"`
    IsActive        bool           `gorm:"default:true"`
    Multiplier      int            `gorm:"default:1"`                     // Round up to nearest X naira
    MaxDailyAmount  int64          `gorm:"default:1000"`                  // IN KOBO (10 naira default)
    TotalRoundup    int64          `gorm:"default:0"`                     // IN KOBO
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    UpdatedAt       time.Time      `gorm:"not null;default:now()"`
    
    // Relationships
    User            User           `gorm:"foreignKey:UserID"`
    SavingsGoal     SavingsGoal    `gorm:"foreignKey:SavingsGoalID"`
}