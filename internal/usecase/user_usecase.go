package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Hdeee1/go-register-login-profile/internal/delivery/http/dto"
	"github.com/Hdeee1/go-register-login-profile/internal/domain"
	"github.com/Hdeee1/go-register-login-profile/pkg/jwt"
	"github.com/Hdeee1/go-register-login-profile/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo domain.UserRepository
	logger *zap.Logger
}

func NewUserUseCase(r domain.UserRepository, logger *zap.Logger) domain.UserUseCase {
	return &userUseCase{userRepo: r, logger: logger}
}

func (u *userUseCase) Register(input dto.RegisterRequest, ctx context.Context) (*domain.User, error) {
	data, err := u.userRepo.FindByEmailOrUsername(input.Email, input.Username, ctx)
	if err == nil && data != nil {
		if data.Email == input.Email {
			return nil, errors.New("email already registered")
		}
		if data.Username == input.Username {
			return nil, errors.New("username already taken")
		}
		u.logger.Warn("email or username already used")
	}

	if err := utils.ValidatePassword(input.Password); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error("failed to hash password")
		return nil, err
	}

	input.Password = string(hash)

	var user domain.User
	user.FullName = input.FullName
	user.Username = input.Username
	user.Email = input.Email
	user.Password = input.Password

	if err := u.userRepo.Create(&user, ctx); err != nil {
		return nil, fmt.Errorf("failed to create user, error: %w", err)
	}

	return &user, nil
}

func (u *userUseCase) Login(input dto.LoginRequest, ctx context.Context) (*domain.User, string, string, error) {
	password := input.Password

	var user domain.User
	user.Email = input.Email
	user.Password = input.Password

	if err := u.userRepo.GetByEmail(&user, ctx); err != nil {
		u.logger.Warn("email not found")
		return nil, "", "", errors.New("wrong email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		u.logger.Warn("wrong password")
		return nil, "", "", errors.New("wrong email or password")
	}

	accessKey := os.Getenv("JWT_ACCESS_SECRET")
	accessToken, err := jwt.GenerateToken(user.UserID, accessKey, 1*time.Hour)
	if err != nil {
		u.logger.Error("failed to generate access token")
		return nil, "", "", errors.New("failed to generate token")
	}

	refreshKey := os.Getenv("JWT_REFRESH_SECRET")
	refreshToken, err := jwt.GenerateToken(user.UserID, refreshKey, 5 * time.Minute)
	if err != nil {
		u.logger.Error("failed to generate refresh token")
		return nil, "", "", errors.New("failed to generate token")
	}

	return &user, accessToken, refreshToken, nil
}

func (u *userUseCase) Refresh(input dto.RefreshTokenRequest, ctx context.Context) (string, error) {
	refreshToken := input.RefreshToken

	refreshKey := os.Getenv("JWT_REFRESH_SECRET")
	claims, err := jwt.ValidateToken(refreshToken, refreshKey)
	if err != nil {
		return "", errors.New("invalid token")
	}

	accessKey := os.Getenv("JWT_ACCESS_SECRET")
	tokenString, err := jwt.GenerateToken(claims.UserId, accessKey, time.Hour)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u *userUseCase) GetProfile(userId int, ctx context.Context) (*domain.User, error) {
	user, err := u.userRepo.GetByUserID(userId, ctx)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) UpdateProfile(userId int, input dto.UpdateProfileRequest, ctx context.Context) (*domain.User, error) {
	if input.Password == "" && input.Username == "" {
		return nil, errors.New("no field to update")
	}

	if input.Password != "" {
		if err := utils.ValidatePassword(input.Password); err != nil {
			return nil, err
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		input.Password = string(hash)
	}

	var user domain.User

	user.UserID = userId
	user.Password = input.Password
	user.Username = input.Username

	if err := u.userRepo.Update(&user, ctx); err != nil {
		u.logger.Error(err.Error())
		return nil, fmt.Errorf("failed to update user, error: %w", err)
	}

	updateUser, err := u.GetProfile(userId, ctx)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, err
	}

	return updateUser, nil
}

func (u *userUseCase) ForgotPassword(input dto.ForgotPasswordRequest, ctx context.Context) error {
	var user domain.User
	user.Email = input.Email

	if err := u.userRepo.GetByEmail(&user, ctx); err != nil {
		u.logger.Warn("user not found")
		return errors.New("user not found")
	}

	randNum := rand.Intn(1000000)
	otp := fmt.Sprintf("%06d", randNum)
	exp := time.Now().Add(5 * time.Minute)

	if err := u.userRepo.SaveOTP(input.Email, otp, exp, ctx); err != nil {
		u.logger.Error(err.Error())
		return err
	}

	fmt.Println("The OTP code for", input.Email, "is", otp)
	return nil
}

func (u *userUseCase) ResetPassword(input dto.ResetPasswordRequest, ctx context.Context) error {
	otp, exp, err := u.userRepo.FindOTP(input.Email, ctx)
	if err != nil {
		u.logger.Error(err.Error())
		return errors.New("Wrong email")
	}
	if otp != input.OTP {
		u.logger.Warn("invalid otp")
		return errors.New("The OTP code is invalid")
	}
	if time.Now().After(exp) {
		u.logger.Warn("otp expired")
		return errors.New("The OTP has been expired")
	}

	var user domain.User
	user.Email = input.Email
	if err := u.userRepo.GetByEmail(&user, ctx); err != nil {
		u.logger.Error(err.Error())
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.NewPassword)); err == nil {
		u.logger.Warn("new password is same as the old password")
		return errors.New("The new password cannot be the same as the old password")
	}
	if err := utils.ValidatePassword(input.NewPassword); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	user.Password = string(hash)
	if err := u.userRepo.Update(&user, ctx); err != nil {
		u.logger.Error(err.Error())
		return err
	}
	u.userRepo.DeleteOTP(input.Email, ctx)

	return nil
}
