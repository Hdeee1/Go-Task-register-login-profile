package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Hdeee1/go-register-login-profile/internal/domain"
	"github.com/Hdeee1/go-register-login-profile/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)


type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(r domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo: r}
}

func (u *userUsecase) Register(input domain.RegisterRequest, ctx context.Context) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	input.Password = string(hash)

	var user domain.User
	user.FullName = input.FullName
	user.Username = input.Username
	user.Email = input.Email
	user.Password = input.Password

	if err := u.userRepo.Create(&user, ctx); err != nil {
		return  nil, fmt.Errorf("failed to create user, error: %w", err)
	}

	return &user, nil
}

func (u *userUsecase) Login(input domain.LoginRequest, ctx context.Context) (*domain.User, string, string, error) {
	password := input.Password

	var user domain.User
	user.Email = input.Email
	user.Password = input.Password
	
	if err := u.userRepo.GetByEmail(&user, ctx); err != nil {
		return nil, "", "", errors.New("wrong email or password")
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", "", errors.New("wrong email or password")
	}

	accessKey := os.Getenv("JWT_ACCESS_SECRET")
	accessToken, err := jwt.GenerateToken(user.Id, accessKey, 1 * time.Hour)
	if err != nil {
		return nil, "", "", errors.New("wrong email or password")
	}

	refreshKey := os.Getenv("JWT_REFRESH_SECRET")
	refreshToken, err := jwt.GenerateToken(user.Id, refreshKey, 24 * time.Hour)
	if err != nil {
		return nil, "", "", errors.New("wrong email or password")
	}

	return &user, accessToken, refreshToken, nil
}

func (u *userUsecase) GetProfile(userId int, ctx context.Context) (*domain.User, error) {
	user, err := u.userRepo.GetById(userId)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}