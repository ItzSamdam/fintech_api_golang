package entities

import (
    "time"
    "github.com/google/uuid"
)

type FeeConfig struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    BillType        string         `gorm:"size:50;uniqueIndex:idx_fee_type;not null"` // transfer, airtime, data, electricity, betting
    FeeType         string         `gorm:"size:20;not null"`                          // percentage, fixed
    FeeValue        float64        `gorm:"type:decimal(10,2);not null"`               // If percentage: 1.5 = 1.5%
    // CapAmount       int64          `gorm:"default:0"`                                 // Maximum fee in kobo (0 = no cap)
    CapAmount       AmountInKobo   `gorm:"type:bigint;default:0"`
    // MinAmount       int64          `gorm:"default:0"`                                 // Minimum transaction amount in kobo
    MinAmount       AmountInKobo   `gorm:"type:bigint;default:0"`
    // MaxAmount       int64          `gorm:"default:0"`                                 // Maximum transaction amount in kobo
    MaxAmount       AmountInKobo   `gorm:"type:bigint;default:0"`
    VATRate         float64        `gorm:"type:decimal(5,2);default:7.5"`             // VAT percentage
    IsActive        bool           `gorm:"default:true"`
    EffectiveFrom   time.Time      `gorm:"not null;default:now()"`
    EffectiveTo     *time.Time
    CreatedBy       *uuid.UUID     `gorm:"type:uuid"`
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    UpdatedAt       time.Time      `gorm:"not null;default:now()"`
}

type TierLimit struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Tier            int            `gorm:"uniqueIndex:idx_tier_limit;not null"`
    // DailyLimit      int64          `gorm:"not null"`          // IN KOBO
    DailyLimit      AmountInKobo   `gorm:"type:bigint;not null"`
    // WeeklyLimit     int64          `gorm:"not null"`          // IN KOBO
    WeeklyLimit     AmountInKobo   `gorm:"type:bigint;not null"`
    // MonthlyLimit    int64          `gorm:"not null"`          // IN KOBO
    MonthlyLimit    AmountInKobo   `gorm:"type:bigint;not null"`
    // SingleTxLimit   int64          `gorm:"not null"`          // IN KOBO
    SingleTxLimit   AmountInKobo   `gorm:"type:bigint;not null"`
    CanSendMoney    bool           `gorm:"default:true"`
    CanBuyAirtime   bool           `gorm:"default:true"`
    CanBuyData      bool           `gorm:"default:true"`
    CanPayElectricity bool         `gorm:"default:true"`
    CanFundBetting  bool           `gorm:"default:true"`
    CanSave         bool           `gorm:"default:true"`
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    UpdatedAt       time.Time      `gorm:"not null;default:now()"`
}