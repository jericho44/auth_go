package models

import "time"

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginAttempt represents a login attempt record
type LoginAttempt struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	IPAddress   string    `json:"ip_address"`
	Success     bool      `json:"success"`
	AttemptedAt time.Time `json:"attempted_at"`
}

// UserStats represents user statistics
type UserStats struct {
	UserID         int       `json:"user_id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	AccountAgeDays int       `json:"account_age_days"`
	LoginCount     int       `json:"login_count,omitempty"`
	LastLoginAt    time.Time `json:"last_login_at,omitempty"`
}
