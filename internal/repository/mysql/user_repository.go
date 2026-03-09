package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Hdeee1/go-register-login-profile/internal/domain"
)

type mySqlUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) (domain.UserRepository, error) {
	return &mySqlUserRepository{db: db}, nil
}

func (m *mySqlUserRepository) Create(user *domain.User, ctx context.Context) error {
	query := "INSERT INTO users (full_name, username, email, password) VALUES (?, ?, ?, ?)"
	res, err := m.db.Exec(query, user.FullName, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	UserID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	user.UserID = int(UserID)
	return nil
}

func (m *mySqlUserRepository) GetByEmail(user *domain.User, ctx context.Context) error {
	query := "SELECT user_id, full_name, username, email, password, created_at, updated_at FROM users WHERE email = ?"
	row := m.db.QueryRow(query, user.Email)

	if err := row.Scan(
		&user.UserID,
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

func (m *mySqlUserRepository) GetByUserID(UserID int) (*domain.User, error) {
	query := "SELECT user_id, full_name, username, email, password, created_at, updated_at FROM users WHERE user_id = ?"
	row := m.db.QueryRow(query, UserID)

	var user domain.User
	if err := row.Scan(
		&user.UserID,
		&user.FullName,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *mySqlUserRepository) FindByEmailOrUsername(email, username string) (*domain.User, error) {
	query := "SELECT user_id, full_name, username, email, password, created_at, updated_at FROM users WHERE email = ? OR username = ?"
	row := m.db.QueryRow(query, email, username)

	var user domain.User
	if err := row.Scan(
		&user.UserID,
		&user.FullName,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *mySqlUserRepository) Update(user *domain.User, ctx context.Context) error {
	fields := []string{}
	args := []any{}

	if user.Username != "" {
		fields = append(fields, "username = ?")
		args = append(args, user.Username)
	}

	if user.Password != "" {
		fields = append(fields, "password = ?")
		args = append(args, user.Password)
	}

	if len(fields) == 0 {
		return errors.New("no fields to update")
	}

	args = append(args, user.UserID)
	query := "UPDATE users SET " + strings.Join(fields, ", ") + " WHERE user_id = ?"

	_, err := m.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *mySqlUserRepository) SaveOTP(email, otp string, expiresAt time.Time, ctx context.Context) error {
	query := "INSERT INTO password_resets (email, otp, expires_at) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE otp = ?, expires_at = ?"
	_, err := m.db.Exec(query, email, otp, expiresAt, otp, expiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (m *mySqlUserRepository) FindOTP(email string, ctx context.Context) (string, time.Time, error) {
	query := "SELECT otp, expires_at FROM password_resets WHERE email = ?"
	row := m.db.QueryRow(query, email)

	var otp string
	var expires time.Time
	if err := row.Scan(
		&otp,
		&expires,
	); err != nil {
		return "", time.Time{}, err
	}

	return otp, expires, nil
}

func (m *mySqlUserRepository) DeleteOTP(email string, ctx context.Context) error {
	query := "DELETE FROM password_resets WHERE email = ?"
	_, err := m.db.Exec(query, email)
	return err
}
