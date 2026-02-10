package domain

import "time"

type User struct {
	Id			int			`json:"id" `
	FullName	string		`json:"full_name" `
	Username	string		`json:"username" `
	Email		string		`json:"email" `
	Password	string		`json:"password" `
	CreatedAt	time.Time	`json:"created_at" `
	UpdatedAt	time.Time	`json:"updated_at" `
}

type UserRepository interface {
	Create()
	GetByEmail(email string)
}

type UserUsecase interface {
	Register(fullName, username, email, password string)
	Login(email, password string)
	GetProfile()
}