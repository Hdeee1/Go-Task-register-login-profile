package dto

// Request
type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=8"`
}

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	OTP         string `json:"otp_code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// Responses
type RegisterResponse struct {
	UserID		int	   `json:"user_id"`
	FullName	string `json:"full_name"`
	Username	string `json:"username"`
	Email		string `json:"email"`
}

type LoginResponse struct {
	Username	 string `json:"username"`
	Email		 string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}