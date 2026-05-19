package entities

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    PhoneNumber     string         `gorm:"uniqueIndex:idx_user_phone;size:15;not null"`
	FirstName       string         `gorm:"size:100"`
	LastName        string         `gorm:"size:100"`
	MiddleName      string         `gorm:"size:100"`
	DateOfBirth     time.Time      `gorm:"type:date"`
	Gender          string         `gorm:"size:10"`
    Email           string         `gorm:"uniqueIndex:idx_user_email;size:255"`
    BVN             string         `gorm:"index:idx_user_bvn;size:11"`              // Encrypted
    NIN             string         `gorm:"index:idx_user_nin;size:11"`              // Encrypted
    Tier            int            `gorm:"default:0;not null"`                      // 0,1,2,3
    IsActive        bool           `gorm:"default:true;not null"`
    IsSuspended     bool           `gorm:"default:false;not null"`
    SuspendedAt     *time.Time     
    SuspensionReason string        `gorm:"size:255"`
    PasswordHash    string         `gorm:"size:255;not null"`
    FacePhotoURL    string         `gorm:"size:500"`
    FaceEmbedding   string         `gorm:"type:text"`                               // For liveness check
    DeviceID        string         `gorm:"index:idx_user_device;size:255"`
    LastLoginAt     *time.Time
    LastLoginIP     string         `gorm:"size:45"`
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    UpdatedAt       time.Time      `gorm:"not null;default:now()"`
    DeletedAt       gorm.DeletedAt `gorm:"index"`
    
    // Relationships
    // Wallets         []Wallet       `gorm:"foreignKey:UserID"`
    // KYC             KYC            `gorm:"foreignKey:UserID"`
    // Sessions        []Session      `gorm:"foreignKey:UserID"`
    // Devices         []TrustedDevice `gorm:"foreignKey:UserID"`
    // TwoFAs          []TwoFA        `gorm:"foreignKey:UserID"`
}

type KYC struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID          uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_kyc_user;not null"`
    BVNVerified     bool           `gorm:"default:false"`
    BVNVerifiedAt   *time.Time
    NINVerified     bool           `gorm:"default:false"`
    NINVerifiedAt   *time.Time
    FaceVerified    bool           `gorm:"default:false"`
    FaceVerifiedAt  *time.Time
    LivenessScore   float64        `gorm:"type:decimal(5,2)"`
    IDCardFront     string         `gorm:"size:500"`     // S3/Cloudinary URL
    IDCardBack      string         `gorm:"size:500"`
    PassportPhoto   string         `gorm:"size:500"`
    UtilityBill     string         `gorm:"size:500"`     // Address proof
    Status          string         `gorm:"size:20;default:'pending'"` // pending, approved, rejected
    RejectionReason string         `gorm:"size:255"`
    ApprovedBy      *uuid.UUID     `gorm:"type:uuid"`
    ApprovedAt      *time.Time
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    UpdatedAt       time.Time      `gorm:"not null;default:now()"`
    
    // Relationships
    User            User           `gorm:"foreignKey:UserID"`
}

type Session struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID          uuid.UUID      `gorm:"type:uuid;index:idx_session_user;not null"`
    Token           string         `gorm:"uniqueIndex:idx_session_token;size:500;not null"`
    RefreshToken    string         `gorm:"uniqueIndex:idx_session_refresh;size:500"`
    IPAddress       string         `gorm:"size:45;not null"`
    UserAgent       string         `gorm:"size:500"`
    DeviceName      string         `gorm:"size:255"`
    IsActive        bool           `gorm:"default:true;not null"`
    ExpiresAt       time.Time      `gorm:"not null"`
    LastActiveAt    time.Time      `gorm:"not null;default:now()"`
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    
    // Relationships
    User            User           `gorm:"foreignKey:UserID"`
}

type TrustedDevice struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID          uuid.UUID      `gorm:"type:uuid;index:idx_device_user;not null"`
    DeviceID        string         `gorm:"size:255;not null"`
    DeviceName      string         `gorm:"size:255"`
    DeviceType      string         `gorm:"size:50"` // ios, android, web
    IsTrusted       bool           `gorm:"default:true"`
    LastUsedAt      time.Time
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    
    // Composite unique index
    // Unique constraint: UserID + DeviceID
    
    // Relationships
    User            User           `gorm:"foreignKey:UserID"`
}

type TwoFA struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID          uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_2fa_user;not null"`
    Secret          string         `gorm:"size:255;not null"`
    BackupCodes     string         `gorm:"type:text"` // JSON array of hashed codes
    IsEnabled       bool           `gorm:"default:false"`
    VerifiedAt      *time.Time
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
    UpdatedAt       time.Time      `gorm:"not null;default:now()"`
    
    // Relationships
    User            User           `gorm:"foreignKey:UserID"`
}

type OTP struct {
    ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    PhoneNumber     string         `gorm:"index:idx_otp_phone;size:15;not null"`
    Code            string         `gorm:"size:6;not null"`
    Purpose         string         `gorm:"size:50;not null"` // registration, login, reset_password, transfer
    IsUsed          bool           `gorm:"default:false"`
    Attempts        int            `gorm:"default:0"`
    ExpiresAt       time.Time      `gorm:"not null"`
    CreatedAt       time.Time      `gorm:"not null;default:now()"`
}

// BeforeCreate GORM hook - generates UUID if not set
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return nil
}