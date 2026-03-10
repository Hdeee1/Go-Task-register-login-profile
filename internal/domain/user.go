package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID    int       	 `json:"user_id" gorm:"primaryKey;column:user_id"`
	FullName  string    	 `json:"full_name" gorm:"column:full_name"`
	Username  string    	 `json:"username" gorm:"column:username"`
	Email     string    	 `json:"email" gorm:"column:email"`
	Password  string    	 `json:"password" gorm:"column:password"`
	CreatedAt time.Time 	 `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time 	 `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index"`
}

type PasswordReset struct {
	ID			int			`json:"id" gorm:"PrimaryKey;column:id"`
	Email 		string		`json:"email" binding:"unique"`
	OTPCode 	string		`json:"otp_code" gorm:"column:otp_code"`
	ExpiresAt	time.Time	`json:"expires_at" gorm:"expires_at"`
}

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

type UserRepository interface {
	Create(user *User, ctx context.Context) error
	Update(user *User, ctx context.Context) error
	GetByUserID(id int) (*User, error)
	GetByEmail(user *User, ctx context.Context) error
	FindByEmailOrUsername(email, username string) (*User, error)
	SaveOTP(email, otp string, expiresAt time.Time, ctx context.Context) error
	FindOTP(email string, ctx context.Context) (string, time.Time, error)
	DeleteOTP(email string, ctx context.Context) error
}

type UserUseCase interface {
	Login(input LoginRequest, ctx context.Context) (*User, string, string, error)
	Refresh(input RefreshTokenRequest, ctx context.Context) (string, error)
	Register(input RegisterRequest, ctx context.Context) (*User, error)
	GetProfile(userId int, ctx context.Context) (*User, error)
	UpdateProfile(userId int, input UpdateProfileRequest, ctx context.Context) (*User, error)
	ResetPassword(input ResetPasswordRequest, ctx context.Context) error
	ForgotPassword(input ForgotPasswordRequest, ctx context.Context) error
}
