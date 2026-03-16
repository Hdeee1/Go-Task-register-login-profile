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

func (m *mySqlUserRepository) GetByIdentifier(identifier string, ctx context.Context) (*domain.User, error) {
	var user domain.User
	if err := m.db.WithContext(ctx).Where("phone = ? OR email = ?", identifier, identifier).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			m.logger.Error("failed to get user by phone or email", zap.Error(err))
		}
		return nil, err
	}
	return &user, nil
}

func (m *mySqlUserRepository) GetByUserID(UserID int, ctx context.Context) (*domain.User, error) {
	var user domain.User
	if err := m.db.WithContext(ctx).First(&user, UserID).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			m.logger.Error("failed to get user by user_id", zap.Error(err))
		}
		return nil, err
	}

	return &user, nil
}

func (m *mySqlUserRepository) FindByPhoneOrEmailOrUsername(phone, email, username string, ctx context.Context) (*domain.User, error) {
	var user domain.User
	if err := m.db.WithContext(ctx).Where("phone = ? OR email = ? OR username = ?", phone, email, username).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			m.logger.Error("failed to find by phone, email or username", zap.Error(err))
		}
		return nil, err
	}

	return &user, nil
}

func (m *mySqlUserRepository) Update(user *domain.User, ctx context.Context) error {
	return m.db.WithContext(ctx).Model(user).Updates(user).Error
}

func (m *mySqlUserRepository) SaveOTP(identifier, otp string, expiresAt time.Time, ctx context.Context) error {
	dataOTP := domain.PasswordReset{
		Identifier: identifier,
		OTPCode: otp,
		ExpiresAt: expiresAt,
	}
	return m.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "identifier"}},
		DoUpdates: clause.AssignmentColumns([]string{"otp_code", "expires_at"}),
	}).Create(&dataOTP).Error
}

func (m *mySqlUserRepository) FindOTP(identifier string, ctx context.Context) (string, time.Time, error) {
	var reset domain.PasswordReset
	if err := m.db.WithContext(ctx).Where("identifier = ?", identifier).First(&reset).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			m.logger.Error("failed to find opt", zap.Error(err))
		}
		return "", time.Time{}, err
	}
	return reset.OTPCode, reset.ExpiresAt, nil
}

func (m *mySqlUserRepository) DeleteOTP(identifier string, ctx context.Context) error {
	return m.db.WithContext(ctx).Where("identifier = ?", identifier).Delete(&domain.PasswordReset{}).Error
}
