package usecase

import (
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

func (u *userUsecase) Register(user domain.User) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hash)

	if err := u.userRepo.Create(&user); err != nil {
		return  nil, fmt.Errorf("failed to create user, error: %w", err)
	}

	return &user, nil
}

func (u *userUsecase) Login(user domain.User) (*domain.User, string, string, error) {
	if err := u.userRepo.GetByEmail(&user); err != nil {
		return nil, "", "", errors.New("email not found")
	}

	var User domain.User
	hashedPass := User.Password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(user.Password)); err != nil {
		return nil, "", "", err
	}

	accessKey := os.Getenv("JWT_ACCESS_SECRET")
	accessToken, err := jwt.GenerateToken(user.Id, accessKey, 1 * time.Hour)
	if err != nil {
		return nil, "", "", err
	}

	refreshKey := os.Getenv("JWT_REFRESH_SECRET")
	refreshToken, err := jwt.GenerateToken(user.Id, refreshKey, 24 * time.Hour)
	if err != nil {
		return nil, "", "", err
	}

	return &user, accessToken, refreshToken, nil
}

func (u *userUsecase) GetProfile() (*domain.User, error) {
	return nil, nil
}