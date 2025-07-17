package repositories

import (
	"database/sql"
	"time"

	"auth-jwt/models"
)

// AuthRepositoryInterface defines the contract for authentication data operations
type AuthRepositoryInterface interface {
	CreatePasswordResetToken(userID int, token string, expiresAt time.Time) error
	GetPasswordResetToken(token string) (*models.PasswordResetToken, error)
	DeletePasswordResetToken(token string) error
	CleanupExpiredTokens() error
	CreateLoginAttempt(username, ipAddress string, success bool) error
	GetRecentLoginAttempts(username string, since time.Time) (int, error)
}

// AuthRepository implements AuthRepositoryInterface
type AuthRepository struct {
	db *sql.DB
}

// NewAuthRepository creates a new auth repository instance
func NewAuthRepository(db *sql.DB) AuthRepositoryInterface {
	return &AuthRepository{
		db: db,
	}
}

// CreatePasswordResetToken stores a password reset token
func (ar *AuthRepository) CreatePasswordResetToken(userID int, token string, expiresAt time.Time) error {
	query := `INSERT INTO password_reset_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := ar.db.Exec(query, userID, token, expiresAt)
	return err
}

// GetPasswordResetToken retrieves a password reset token
func (ar *AuthRepository) GetPasswordResetToken(token string) (*models.PasswordResetToken, error) {
	resetToken := &models.PasswordResetToken{}
	query := `SELECT id, user_id, token, expires_at, created_at FROM password_reset_tokens WHERE token = $1 AND expires_at > NOW()`
	err := ar.db.QueryRow(query, token).Scan(&resetToken.ID, &resetToken.UserID, &resetToken.Token, &resetToken.ExpiresAt, &resetToken.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return resetToken, nil
}

// DeletePasswordResetToken removes a password reset token
func (ar *AuthRepository) DeletePasswordResetToken(token string) error {
	query := `DELETE FROM password_reset_tokens WHERE token = $1`
	_, err := ar.db.Exec(query, token)
	return err
}

// CleanupExpiredTokens removes expired password reset tokens
func (ar *AuthRepository) CleanupExpiredTokens() error {
	query := `DELETE FROM password_reset_tokens WHERE expires_at <= NOW()`
	_, err := ar.db.Exec(query)
	return err
}

// CreateLoginAttempt logs a login attempt
func (ar *AuthRepository) CreateLoginAttempt(username, ipAddress string, success bool) error {
	query := `INSERT INTO login_attempts (username, ip_address, success, attempted_at) VALUES ($1, $2, $3, NOW())`
	_, err := ar.db.Exec(query, username, ipAddress, success)
	return err
}

// GetRecentLoginAttempts gets the number of recent login attempts for a username
func (ar *AuthRepository) GetRecentLoginAttempts(username string, since time.Time) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM login_attempts WHERE username = $1 AND attempted_at >= $2 AND success = false`
	err := ar.db.QueryRow(query, username, since).Scan(&count)
	return count, err
}
