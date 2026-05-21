package services

import (
    "context"
    "crypto/rand"
    "errors"
    "fmt"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    
    "fintech_api_golang/internal/config"
    "fintech_api_golang/internal/core/entities"
    "fintech_api_golang/internal/core/interfaces"
    "fintech_api_golang/internal/dto/request"
    "fintech_api_golang/internal/dto/response"
)

type AuthService struct {
    userRepo      interfaces.UserRepository
    kycRepo       interfaces.KYCRepository
    sessionRepo   interfaces.SessionRepository
    otpRepo       interfaces.OTPRepository
    walletRepo    interfaces.WalletRepository
    config        *config.Config
}

func NewAuthService(
    userRepo interfaces.UserRepository,
    kycRepo interfaces.KYCRepository,
    sessionRepo interfaces.SessionRepository,
    otpRepo interfaces.OTPRepository,
    walletRepo interfaces.WalletRepository,
    cfg *config.Config,
) *AuthService {
    return &AuthService{
        userRepo:    userRepo,
        kycRepo:     kycRepo,
        sessionRepo: sessionRepo,
        otpRepo:     otpRepo,
        walletRepo:  walletRepo,
        config:      cfg,
    }
}

func (s *AuthService) RegisterPhone(ctx context.Context, req *request.RegisterPhoneRequest) (*response.OTPResponse, error) {
    // Check if user already exists
    existingUser, err := s.userRepo.GetByPhoneNumber(ctx, req.PhoneNumber)
    if err != nil {
        return nil, err
    }
    
    if existingUser != nil {
        return nil, errors.New("user already exists with this phone number")
    }
    
    // Generate OTP
    otpCode := generateOTP()
    
    // Create OTP record
    otp := &entities.OTP{
        ID:          uuid.New(),
        PhoneNumber: req.PhoneNumber,
        Code:        otpCode,
        Purpose:     "registration",
        ExpiresAt:   time.Now().Add(10 * time.Minute),
    }
    
    if err := s.otpRepo.Create(ctx, otp); err != nil {
        return nil, err
    }
    
    // TODO: Send SMS via notification service
    // For now, just return the OTP in response (remove in production)
    
    return &response.OTPResponse{
        Reference: otp.ID.String(),
        ExpiresIn: 600, // 10 minutes in seconds
    }, nil
}

func (s *AuthService) VerifyOTP(ctx context.Context, req *request.VerifyOTPRequest) (*response.AuthResponse, error) {
    // Validate OTP
    otp, err := s.otpRepo.GetValidOTP(ctx, req.PhoneNumber, req.Code, "registration")
    if err != nil {
        return nil, err
    }
    
    if otp == nil {
        return nil, errors.New("invalid or expired OTP")
    }
    
    // Mark OTP as used
    if err := s.otpRepo.MarkAsUsed(ctx, otp.ID); err != nil {
        return nil, err
    }
    
    // Create user (Tier 0 - unverified)
    user := &entities.User{
        ID:          uuid.New(),
        PhoneNumber: req.PhoneNumber,
        PasswordHash: "", // Will be set later
        FirstName: "", // Will be set later
        LastName: "", // Will be set later
        MiddleName: "", // Will be set later
        Tier:        0,
        IsActive:    true,
        DeviceID:    req.DeviceID,
    }
    
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    // Create KYC record
    kyc := &entities.KYC{
        ID:     uuid.New(),
        UserID: user.ID,
        Status: "pending",
    }
    
    if err := s.kycRepo.Create(ctx, kyc); err != nil {
        return nil, err
    }
    
    // Create wallet
    wallet := &entities.Wallet{
        ID:     uuid.New(),
        UserID: user.ID,
        Balance: 0,
        Currency: "NGN",
    }
    
    if err := s.walletRepo.Create(ctx, wallet); err != nil {
        return nil, err
    }
    
    // Create session
    accessToken, refreshToken, err := s.createUserSession(ctx, user, req.DeviceID, req.DeviceName, req.PhoneNumber)
    if err != nil {
        return nil, err
    }
    
    return &response.AuthResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    int64(s.config.JWT.AccessExpiry.Seconds()),
        TokenType:    "Bearer",
        User: response.UserResponse{
            ID:          user.ID,
            PhoneNumber: user.PhoneNumber,
            Email:       user.Email,
            Tier:        user.Tier,
            IsActive:    user.IsActive,
            IsSuspended: user.IsSuspended,
            CreatedAt:   user.CreatedAt,
        },
    }, nil
}

func (s *AuthService) RegisterBVN(ctx context.Context, userID uuid.UUID, req *request.RegisterBVNRequest) error {
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return err
    }
    
    if user == nil {
        return errors.New("user not found")
    }
    
    // TODO: Verify BVN with external service
    // For now, just store the encrypted values
    
    user.BVN = req.BVN
    user.NIN = req.NIN
    user.Tier = 1 // Upgrade to Tier 1 after BVV verification
    
    if err := s.userRepo.Update(ctx, user); err != nil {
        return err
    }
    
    // Update KYC record
    kyc, err := s.kycRepo.GetByUserID(ctx, userID)
    if err != nil {
        return err
    }
    
    if kyc != nil {
        kyc.BVNVerified = true
        kyc.BVNVerifiedAt = timeNow()
        if err := s.kycRepo.Update(ctx, kyc); err != nil {
            return err
        }
    }
    
    return nil
}

func (s *AuthService) VerifyFace(ctx context.Context, userID uuid.UUID, req *request.VerifyFaceRequest) error {
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return err
    }
    
    if user == nil {
        return errors.New("user not found")
    }
    
    // TODO: Call face recognition API
    // For now, simulate successful verification
    
    user.FacePhotoURL = req.FacePhoto
    user.Tier = 2 // Upgrade to Tier 2 after face verification
    
    if err := s.userRepo.Update(ctx, user); err != nil {
        return err
    }
    
    // Update KYC record
    kyc, err := s.kycRepo.GetByUserID(ctx, userID)
    if err != nil {
        return err
    }
    
    if kyc != nil {
        kyc.FaceVerified = true
        kyc.FaceVerifiedAt = timeNow()
        kyc.LivenessScore = 0.95 // Mock score
        if err := s.kycRepo.Update(ctx, kyc); err != nil {
            return err
        }
    }
    
    return nil
}

func (s *AuthService) Login(ctx context.Context, req *request.LoginRequest) (*response.AuthResponse, error) {
    user, err := s.userRepo.GetByPhoneNumber(ctx, req.PhoneNumber)
    if err != nil {
        return nil, err
    }
    
    if user == nil {
        return nil, errors.New("invalid credentials")
    }
    
    // Check password (if set)
    if user.PasswordHash != "" {
        if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
            return nil, errors.New("invalid credentials")
        }
    }
    
    // Check if user is suspended
    if user.IsSuspended {
        return nil, errors.New("account is suspended")
    }
    
    // Update last login
    if err := s.userRepo.UpdateLastLogin(ctx, user.ID, req.PhoneNumber); err != nil {
        // Log but don't fail
    }
    
    // Create session
    accessToken, refreshToken, err := s.createUserSession(ctx, user, req.DeviceID, "", req.PhoneNumber)
    if err != nil {
        return nil, err
    }
    
    return &response.AuthResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    int64(s.config.JWT.AccessExpiry.Seconds()),
        TokenType:    "Bearer",
        User: response.UserResponse{
            ID:          user.ID,
            PhoneNumber: user.PhoneNumber,
            Email:       user.Email,
            Tier:        user.Tier,
            IsActive:    user.IsActive,
            IsSuspended: user.IsSuspended,
            CreatedAt:   user.CreatedAt,
        },
    }, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, req *request.ChangePasswordRequest) error {
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return err
    }
    
    if user == nil {
        return errors.New("user not found")
    }
    
    // Verify old password
    if user.PasswordHash != "" {
        if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
            return errors.New("invalid old password")
        }
    }
    
    // Hash new password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    
    user.PasswordHash = string(hashedPassword)
    
    return s.userRepo.Update(ctx, user)
}

func (s *AuthService) ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) error {
    // Validate OTP
    otp, err := s.otpRepo.GetValidOTP(ctx, req.PhoneNumber, req.Code, "reset_password")
    if err != nil {
        return err
    }
    
    if otp == nil {
        return errors.New("invalid or expired OTP")
    }
    
    // Mark OTP as used
    if err := s.otpRepo.MarkAsUsed(ctx, otp.ID); err != nil {
        return err
    }
    
    // Get user
    user, err := s.userRepo.GetByPhoneNumber(ctx, req.PhoneNumber)
    if err != nil {
        return err
    }
    
    if user == nil {
        return errors.New("user not found")
    }
    
    // Hash new password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    
    user.PasswordHash = string(hashedPassword)
    
    // Invalidate all sessions
    if err := s.sessionRepo.InvalidateAllUserSessions(ctx, user.ID); err != nil {
        return err
    }
    
    return s.userRepo.Update(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID, token string, allDevices bool) error {
    if allDevices {
        return s.sessionRepo.InvalidateAllUserSessions(ctx, userID)
    }
    
    // Find session by token and invalidate
    session, err := s.sessionRepo.GetByToken(ctx, token)
    if err != nil {
        return err
    }
    
    if session != nil {
        return s.sessionRepo.Invalidate(ctx, session.ID)
    }
    
    return nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*response.AuthResponse, error) {
    // Parse and validate refresh token
    token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.config.JWT.Secret), nil
    })
    
    if err != nil || !token.Valid {
        return nil, errors.New("invalid refresh token")
    }
    
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, errors.New("invalid token claims")
    }
    
    // Check token type
    tokenType, ok := claims["token_type"].(string)
    if !ok || tokenType != "refresh" {
        return nil, errors.New("invalid token type")
    }
    
    userIDStr, ok := claims["user_id"].(string)
    if !ok {
        return nil, errors.New("invalid user ID")
    }
    
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        return nil, err
    }
    
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    if user == nil {
        return nil, errors.New("user not found")
    }
    
    // Create new session
    newAccessToken, newRefreshToken, err := s.createUserSession(ctx, user, "", "", "")
    if err != nil {
        return nil, err
    }
    
    return &response.AuthResponse{
        AccessToken:  newAccessToken,
        RefreshToken: newRefreshToken,
        ExpiresIn:    int64(s.config.JWT.AccessExpiry.Seconds()),
        TokenType:    "Bearer",
        User: response.UserResponse{
            ID:          user.ID,
            PhoneNumber: user.PhoneNumber,
            Email:       user.Email,
            Tier:        user.Tier,
            IsActive:    user.IsActive,
            IsSuspended: user.IsSuspended,
            CreatedAt:   user.CreatedAt,
        },
    }, nil
}

func (s *AuthService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*response.UserResponse, error) {
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    if user == nil {
        return nil, errors.New("user not found")
    }
    
    return &response.UserResponse{
        ID:          user.ID,
        PhoneNumber: user.PhoneNumber,
        Email:       user.Email,
        Tier:        user.Tier,
        IsActive:    user.IsActive,
        IsSuspended: user.IsSuspended,
        CreatedAt:   user.CreatedAt,
        UpdatedAt:   user.UpdatedAt,
    }, nil
}

func (s *AuthService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *request.UpdateUserRequest) error {
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return err
    }
    
    if user == nil {
        return errors.New("user not found")
    }
    
    if req.Email != "" {
        user.Email = req.Email
    }
    
    return s.userRepo.Update(ctx, user)
}

// Private helper methods
func (s *AuthService) createUserSession(ctx context.Context, user *entities.User, deviceID, deviceName, ipAddress string) (string, string, error) {
    // Generate access token
    accessToken, err := s.generateToken(user, "access")
    if err != nil {
        return "", "", err
    }
    
    // Generate refresh token
    refreshToken, err := s.generateToken(user, "refresh")
    if err != nil {
        return "", "", err
    }
    
    // Create session record
    session := &entities.Session{
        ID:           uuid.New(),
        UserID:       user.ID,
        Token:        accessToken,
        RefreshToken: refreshToken,
        IPAddress:    ipAddress,
        DeviceName:   deviceName,
        IsActive:     true,
        ExpiresAt:    time.Now().Add(s.config.JWT.AccessExpiry),
        LastActiveAt: time.Now(),
    }
    
    if err := s.sessionRepo.Create(ctx, session); err != nil {
        return "", "", err
    }
    
    return accessToken, refreshToken, nil
}

func (s *AuthService) generateToken(user *entities.User, tokenType string) (string, error) {
    var expiry time.Duration
    if tokenType == "access" {
        expiry = s.config.JWT.AccessExpiry
    } else {
        expiry = s.config.JWT.RefreshExpiry
    }
    
    claims := jwt.MapClaims{
        "user_id":    user.ID.String(),
        "phone":      user.PhoneNumber,
        "email":      user.Email,
        "tier":       user.Tier,
        "token_type": tokenType,
        "exp":        time.Now().Add(expiry).Unix(),
        "iat":        time.Now().Unix(),
        "iss":        s.config.JWT.Issuer,
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.config.JWT.Secret))
}

// func generateOTP() string {
//     // Generate 6-digit OTP
//     b := make([]byte, 3)
//     rand.Read(b)
//     return fmt.Sprintf("%06d", int(b[0])<<16|int(b[1])<<8|int(b[2])%1000000)
// }
func generateOTP() string {
    b := make([]byte, 3)
    _, _ = rand.Read(b) // always check error in production
    n := (int(b[0]) << 16) | (int(b[1]) << 8) | int(b[2])
    return fmt.Sprintf("%06d", n % 1000000)
}

func timeNow() *time.Time {
    now := time.Now()
    return &now
}