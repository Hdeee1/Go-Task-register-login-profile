package domain

import (
	"context"
	"time"

	"github.com/Hdeee1/go-register-login-profile/internal/delivery/http/dto"
	"gorm.io/gorm"
)

type User struct {
	UserID    int       	 `json:"user_id" gorm:"primaryKey;column:user_id"`
	FullName  string    	 `json:"full_name" gorm:"column:full_name"`
	Username  string    	 `json:"username" gorm:"column:username"`
	Phone	  string		 `json:"phone" gorm:"column:phone"`
	Email     string    	 `json:"email" gorm:"column:email"`
	Password  string    	 `json:"password" gorm:"column:password"`
	CreatedAt time.Time 	 `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time 	 `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index"`
}

type PasswordReset struct {
	ID			int			`gorm:"PrimaryKey;column:id"`
	Identifier	string		`gorm:"column:target;not null;uniqueIndex"`
	OTPCode 	string		`gorm:"column:otp_code"`
	ExpiresAt	time.Time	`gorm:"expires_at"`
}

type UserRepository interface {
	Create(user *User, ctx context.Context) error
	Update(user *User, ctx context.Context) error
	GetByUserID(id int, ctx context.Context) (*User, error)
	GetByIdentifier(identifier string, ctx context.Context) error
	FindByPhoneOrEmailOrUsername(phone, email, username string, ctx context.Context) (*User, error)
	SaveOTP(identifier, otp string, expiresAt time.Time, ctx context.Context) error
	FindOTP(identifier string, ctx context.Context) (string, time.Time, error)
	DeleteOTP(identifier string, ctx context.Context) error
}

type UserUseCase interface {
	Login(input dto.LoginRequest, ctx context.Context) (*User, string, string, error)
	Refresh(input dto.RefreshTokenRequest, ctx context.Context) (string, error)
	Register(input dto.RegisterRequest, ctx context.Context) (*User, error)
	GetProfile(userId int, ctx context.Context) (*User, error)
	UpdateProfile(userId int, input dto.UpdateProfileRequest, ctx context.Context) (*User, error)
	ResetPassword(input dto.ResetPasswordRequest, ctx context.Context) error
	ForgotPassword(input dto.ForgotPasswordRequest, ctx context.Context) error
}
