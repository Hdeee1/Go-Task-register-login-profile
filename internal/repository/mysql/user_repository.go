package mysql

import (
	"database/sql"

	"github.com/Hdeee1/go-register-login-profile/internal/domain"
)

type mySQLUserRepository struct {
	*sql.DB
}

func NewUserRepository(db *sql.DB) (domain.UserRepository, error) {
	return &mySQLUserRepository{db}, nil
}

func (m *mySQLUserRepository) Create(fullName, username, email, password string) error {
	return nil
}

func (m *mySQLUserRepository) GetByEmail(email string) error {
	return nil
}