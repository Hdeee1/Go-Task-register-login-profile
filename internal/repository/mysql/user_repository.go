package repository

import (
	"database/sql"
	"log"

	"github.com/Hdeee1/go-register-login-profile/internal/domain"
)

type mySQLUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) (domain.UserRepository, error) {
	return &mySQLUserRepository{db: db}, nil
}

func (m *mySQLUserRepository) Create(user *domain.User) error {
	query := "INSERT INTO users (full_name, username, email, password) VALUES (?, ?, ?, ?)"
	res, err := m.db.Exec(query, user.FullName, user.Username, user.Email, user.Password)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer m.db.Close()
	
	rows, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err.Error())
	}

	if rows != 1 {
		log.Fatal(err.Error())
	}

	return nil
}

func (m *mySQLUserRepository) GetByEmail(user *domain.User) error {
	return nil
}