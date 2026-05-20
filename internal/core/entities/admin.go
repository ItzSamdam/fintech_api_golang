package entities

import (
    "time"
    "github.com/google/uuid"
)

type AdminUser struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Email           string         `gorm:"uniqueIndex:idx_admin_email;size:255;not null"`
    PasswordHash    string         `gorm:"size:255;not null"`
    FullName        string         `gorm:"size:255;not null"`
    Role            string         `gorm:"size:50;not null;default:'viewer'"` // super_admin, admin, viewer, support
    IsActive        bool           `gorm:"default:true"`
    LastLoginAt     *time.Time
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    UpdatedAt       time.Time      `gorm:"not null;default:now()"`
    
    // Relationships
    AuditLogs       []AuditLog     `gorm:"foreignKey:AdminID"`
}

type Role struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Name            string         `gorm:"uniqueIndex:idx_role_name;size:50;not null"`
    Permissions     string         `gorm:"type:jsonb;not null"`       // JSON array of permissions
    Description     string         `gorm:"size:255"`
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    UpdatedAt       time.Time      `gorm:"not null;default:now()"`
}

type AuditLog struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    AdminID         *uuid.UUID     `gorm:"type:uuid;index:idx_audit_admin"`
    UserID          *uuid.UUID     `gorm:"type:uuid;index:idx_audit_user"`
    Action          string         `gorm:"size:100;not null;index:idx_audit_action"`
    EntityType      string         `gorm:"size:50"`                  // user, wallet, transaction, provider
    EntityID        string         `gorm:"size:100"`
    OldValue        string         `gorm:"type:jsonb"`
    NewValue        string         `gorm:"type:jsonb"`
    IPAddress       string         `gorm:"size:45;not null"`
    UserAgent       string         `gorm:"size:500"`
    Metadata        string         `gorm:"type:jsonb"`
    CreatedAt       time.Time      `gorm:"not null;default:now();index:idx_audit_created"`
    
    // Relationships
    Admin           *AdminUser     `gorm:"foreignKey:AdminID"`
    User            *User          `gorm:"foreignKey:UserID"`
}