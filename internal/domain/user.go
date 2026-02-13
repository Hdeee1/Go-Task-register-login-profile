package domain

import (
	"context"
	"time"
)

type User struct {
	Id			int			`json:"id" `
	FullName	string		`json:"full_name"`
	Username	string		`json:"username"`
	Email		string		`json:"email"`
	Password	string		`json:"password"`
	CreatedAt	time.Time	`json:"created_at" `
	UpdatedAt	time.Time	`json:"updated_at" `
}

type RegisterRequest struct {
	FullName	string `json:"full_name" binding:"required"`
	Username	string `json:"username" binding:"required"`
	Email		string `json:"email" binding:"required,email"`
	Password 	string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email		 string `json:"email" binding:"required,email"`
	Password	 string `json:"password" binding:"required,min=6"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserRepository interface {
	Create(user *User, ctx context.Context) error
	GetByEmail(user *User, ctx context.Context) error
	GetById(id int) (*User, error)
}

type UserUsecase interface {
	Register(user RegisterRequest, ctx context.Context) (*User, error)
	Login(user LoginRequest, ctx context.Context) (*User, string, string, error)
	GetProfile(userId int, ctx context.Context) (*User, error)
}