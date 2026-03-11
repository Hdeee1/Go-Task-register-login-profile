package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Hdeee1/go-register-login-profile/internal/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type mySqlUserRepository struct {
	db *gorm.DB
	logger *zap.Logger
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger) (domain.UserRepository, error) {
	return &mySqlUserRepository{db: db, logger: logger}, nil
}

func (m *mySqlUserRepository) Create(user *domain.User, ctx context.Context) error {
	return m.db.WithContext(ctx).Create(user).Error
}

func (m *mySqlUserRepository) GetByEmail(user *domain.User, ctx context.Context) error {
	return m.db.WithContext(ctx).Where("email = ?", user.Email).First(user).Error
}

func (m *mySqlUserRepository) GetByUserID(UserID int) (*domain.User, error) {
	var user domain.User
	if err := m.db.First(&user, UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

func (m *mySqlUserRepository) FindByEmailOrUsername(email, username string) (*domain.User, error) {
	var user domain.User
	if err := m.db.Where("email = ? OR username = ?", user.Email, user.Username).First(&user).Error; err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	return &user, nil
}

func (m *mySqlUserRepository) Update(user *domain.User, ctx context.Context) error {
	return m.db.WithContext(ctx).Model(user).Updates(user).Error
}

func (m *mySqlUserRepository) SaveOTP(email, otp string, expiresAt time.Time, ctx context.Context) error {
	dataOTP := domain.PasswordReset{
		Email: email,
		OTPCode: otp,
		ExpiresAt: expiresAt,
	}
	return m.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"otp_code", "expires_at"}),
	}).Create(&dataOTP).Error
}

func (m *mySqlUserRepository) FindOTP(email string, ctx context.Context) (string, time.Time, error) {
	var reset domain.PasswordReset
	if err := m.db.WithContext(ctx).Where("email = ?", email).First(&reset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", time.Time{}, err
		}
		return "", time.Time{}, err
	}
	return reset.OTPCode, reset.ExpiresAt, nil
}

func (m *mySqlUserRepository) DeleteOTP(email string, ctx context.Context) error {
	return m.db.WithContext(ctx).Where("email = ?", email).Delete(&domain.PasswordReset{}).Error
}
