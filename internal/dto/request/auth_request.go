package request

type RegisterPhoneRequest struct {
    PhoneNumber string `json:"phone_number" validate:"required,len=11,numeric"`
}

type SendOTPRequest struct {
    PhoneNumber string `json:"phone_number" validate:"required,len=11,numeric"`
    Purpose     string `json:"purpose" validate:"required"`
}

type VerifyOTPRequest struct {
    PhoneNumber string `json:"phone_number" validate:"required,len=11,numeric"`
    Code        string `json:"code" validate:"required,len=6,numeric"`
    DeviceID    string `json:"device_id" validate:"required"`
    DeviceName  string `json:"device_name"`
}

type RegisterBVNRequest struct {
    BVN            string `json:"bvn" validate:"required,len=11,numeric"`
    NIN            string `json:"nin" validate:"required,len=11,numeric"`
    DateOfBirth    string `json:"date_of_birth" validate:"required"`
    FirstName      string `json:"first_name" validate:"required"`
    LastName       string `json:"last_name" validate:"required"`
    MiddleName     string `json:"middle_name"`
}

type VerifyFaceRequest struct {
    FacePhoto     string `json:"face_photo" validate:"required"` // Base64 encoded image
    LivenessVideo string `json:"liveness_video"`                 // Optional for liveness check
}

type LoginRequest struct {
    PhoneNumber string `json:"phone_number" validate:"required,len=11,numeric"`
    Password    string `json:"password" validate:"required,min=6"`
    DeviceID    string `json:"device_id" validate:"required"`
}

type ChangePasswordRequest struct {
    OldPassword string `json:"old_password" validate:"required"`
    NewPassword string `json:"new_password" validate:"required,min=6"`
}

type ResetPasswordRequest struct {
    PhoneNumber string `json:"phone_number" validate:"required,len=11,numeric"`
    Code        string `json:"code" validate:"required,len=6,numeric"`
    NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
    Email       string `json:"email" validate:"omitempty,email"`
    FirstName   string `json:"first_name"`
    LastName    string `json:"last_name"`
    Address     string `json:"address"`
    DateOfBirth string `json:"date_of_birth"`
}

type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutRequest struct {
    AllDevices bool `json:"all_devices"` // Logout from all devices
}