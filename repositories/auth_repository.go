package repositories

import (
	"errors"
	"time"

	"auth-jwt/models"

	"gorm.io/gorm"
)

// AuthRepositoryInterface defines the contract for authentication data operations
type AuthRepositoryInterface interface {
	CreatePasswordResetToken(userID uint, token string, expiresAt time.Time) error
	GetPasswordResetToken(token string) (*models.PasswordResetToken, error)
	DeletePasswordResetToken(token string) error
	CleanupExpiredTokens() error
	CreateLoginAttempt(username, ipAddress string, success bool) error
	GetRecentLoginAttempts(username string, since time.Time) (int, error)
}

// AuthRepository implements AuthRepositoryInterface using GORM
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new auth repository instance
func NewAuthRepository(db *gorm.DB) AuthRepositoryInterface {
	return &AuthRepository{
		db: db,
	}
}

// CreatePasswordResetToken stores a password reset token
func (ar *AuthRepository) CreatePasswordResetToken(userID uint, token string, expiresAt time.Time) error {
	resetToken := &models.PasswordResetToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	result := ar.db.Create(resetToken)
	return result.Error
}

// GetPasswordResetToken retrieves a password reset token
func (ar *AuthRepository) GetPasswordResetToken(token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	result := ar.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&resetToken)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &resetToken, result.Error
}

// DeletePasswordResetToken removes a password reset token
func (ar *AuthRepository) DeletePasswordResetToken(token string) error {
	result := ar.db.Where("token = ?", token).Delete(&models.PasswordResetToken{})
	return result.Error
}

// CleanupExpiredTokens removes expired password reset tokens
func (ar *AuthRepository) CleanupExpiredTokens() error {
	result := ar.db.Where("expires_at <= ?", time.Now()).Delete(&models.PasswordResetToken{})
	return result.Error
}

// CreateLoginAttempt logs a login attempt
func (ar *AuthRepository) CreateLoginAttempt(username, ipAddress string, success bool) error {
	attempt := &models.LoginAttempt{
		Username:    username,
		IPAddress:   ipAddress,
		Success:     success,
		AttemptedAt: time.Now(),
	}
	result := ar.db.Create(attempt)
	return result.Error
}

// GetRecentLoginAttempts gets the number of recent login attempts for a username
func (ar *AuthRepository) GetRecentLoginAttempts(username string, since time.Time) (int, error) {
	var count int64
	result := ar.db.Model(&models.LoginAttempt{}).
		Where("username = ? AND attempted_at >= ? AND success = ?", username, since, false).
		Count(&count)
	return int(count), result.Error
}
