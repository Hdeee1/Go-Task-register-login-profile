package usecase

import (
	"fmt"

	"github.com/Hdeee1/go-register-login-profile/internal/domain"
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

func (u *userUsecase) Login(user domain.User) (*domain.User, error) {
	return nil, nil
}

func (u *userUsecase) GetProfile() (*domain.User, error) {
	return nil, nil
}