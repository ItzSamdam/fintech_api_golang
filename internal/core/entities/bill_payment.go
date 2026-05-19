package entities

import (
    "time"
    "github.com/google/uuid"
)

type BillDetail struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    TransactionID   uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_bill_txn;not null"`
    BillType        string         `gorm:"size:50;index:idx_bill_type;not null"` // airtime, data, electricity, betting
    ProviderID      uuid.UUID      `gorm:"type:uuid;index:idx_bill_provider;not null"`
    ProviderName    string         `gorm:"size:100;not null"`
    PhoneNumber     string         `gorm:"size:15"`                              // For airtime/data
    MeterNumber     string         `gorm:"size:50"`                              // For electricity
    MeterType       string         `gorm:"size:20"`                              // prepaid, postpaid
    CustomerName    string         `gorm:"size:255"`
    CustomerAddress string         `gorm:"size:500"`                             // For electricity
    BettingAccount  string         `gorm:"size:100"`                             // For betting
    BettingOperator string         `gorm:"size:50"`                              // For betting
    DataPlanID      string         `gorm:"size:50"`                              // For data
    DataPlanName    string         `gorm:"size:100"`
    DataVolume      string         `gorm:"size:20"`                              // e.g., "1GB"
    DataValidity    string         `gorm:"size:50"`                              // e.g., "30 days"
    ElectricityToken string        `gorm:"size:100"`                             // Prepaid token
    ElectricityUnits int           `gorm:"default:0"`
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    
    // Relationships
    Transaction     Transaction    `gorm:"foreignKey:TransactionID"`
    Provider        Provider       `gorm:"foreignKey:ProviderID"`
}

type Provider struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Name            string         `gorm:"size:100;not null"`
    Code            string         `gorm:"uniqueIndex:idx_provider_code;size:50;not null"`
    Type            string         `gorm:"size:50;index:idx_provider_type;not null"` // airtime, data, electricity, betting, bank
    Category        string         `gorm:"size:50"`                                  // mtn, glo, ikeja_electric, etc.
    IsActive        bool           `gorm:"default:true;not null"`
    Priority        int            `gorm:"default:100"`                              // Lower = higher priority (for fallback)
    BaseURL         string         `gorm:"size:255"`
    APIKey          string         `gorm:"size:255"`                                 // Encrypted
    APISecret       string         `gorm:"size:255"`                                 // Encrypted
    TimeoutSeconds  int            `gorm:"default:30"`
    RetryCount      int            `gorm:"default:2"`
    MarginPercent   float64        `gorm:"type:decimal(5,2);default:0"`              // Profit margin
    Metadata        string         `gorm:"type:jsonb"`
    LastHealthCheck *time.Time
    HealthStatus    string         `gorm:"size:20;default:'unknown'"`                // healthy, degraded, down
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    UpdatedAt       time.Time      `gorm:"not null;default:now()"`
    
    // Relationships
    BillDetails     []BillDetail   `gorm:"foreignKey:ProviderID"`
    ProviderLogs    []ProviderLog  `gorm:"foreignKey:ProviderID"`
}

type ProviderLog struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    ProviderID      uuid.UUID      `gorm:"type:uuid;index:idx_provider_log_provider;not null"`
    TransactionID   *uuid.UUID     `gorm:"type:uuid;index:idx_provider_log_txn"`
    Endpoint        string         `gorm:"size:255;not null"`
    Request         string         `gorm:"type:text"`
    Response        string         `gorm:"type:text"`
    StatusCode      int
    ResponseTime    int            // milliseconds
    IsError         bool           `gorm:"default:false"`
    ErrorMessage    string         `gorm:"type:text"`
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    
    // Relationships
    Provider        Provider       `gorm:"foreignKey:ProviderID"`
}