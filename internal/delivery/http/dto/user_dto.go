package dto

// Request
type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Username string `json:"username" binding:"required,min=3"`
	Phone 	 string `json:"phone" biding:"required,min=11"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Identifier	string	`json:"identifier" gorm:"column:target;not null;uniqueIndex"`
	Password	string	`json:"password" binding:"required,min=8"`
}

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordRequest struct {
	Identifier	string	`json:"identifier" gorm:"column:target;not null;uniqueIndex"`
}

type ResetPasswordRequest struct {
	Identifier	string	`json:"identifier" `
	OTP         string	`json:"otp_code" binding:"required"`
	NewPassword string 	`json:"new_password" binding:"required,min=8"`
}

// Responses
type RegisterResponse struct {
	UserID		int	   `json:"user_id"`
	FullName	string `json:"full_name"`
	Username	string `json:"username"`
	Phone		string `json:"phone"`
	Email		string `json:"email"`
}

type LoginResponse struct {
	Username	 string `json:"username"`
	Phone		 string `json:"phone"`
	Email		 string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}