package repository

import (
	"database/sql"

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
		return err
	}
	
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return err
	}

	return nil
}

func (m *mySQLUserRepository) GetByEmail(user *domain.User) error {
	query := "SELECT id, full_name, username, email, password, created_at, updated_at FROM users WHERE email = ?"
	row := m.db.QueryRow(query, &user.Email)
	
	if err := row.Scan(
			&user.Id, 
			&user.FullName, 
			&user.Username, 
			&user.Email, 
			&user.Password, 
			&user.CreatedAt, 
			&user.UpdatedAt,
		); err != nil {
		return err
	}
	
	return nil
}

func (m *mySQLUserRepository) GetById(id int) (*domain.User, error) {
	query := "SELECT id, full_name, username, email, password, created_at, updated_at FROM users WHERE id = ?"
	row := m.db.QueryRow(query, id)

	var user domain.User
	
	if err := row.Scan(
			&user.Id, 
			&user.FullName, 
			&user.Username, 
			&user.Email, 
			&user.Password, 
			&user.CreatedAt, 
			&user.UpdatedAt,
	); err != nil {
		return  nil, err
	}

	return &user, nil
}