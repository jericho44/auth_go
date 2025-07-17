package models

import (
	"time"

	"gorm.io/gorm"
)

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	User      User           `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Token     string         `json:"token" gorm:"uniqueIndex;size:255;not null"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// LoginAttempt represents a login attempt record
type LoginAttempt struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Username    string         `json:"username" gorm:"size:50;not null;index"`
	IPAddress   string         `json:"ip_address" gorm:"size:45"`
	Success     bool           `json:"success" gorm:"not null;default:false"`
	AttemptedAt time.Time      `json:"attempted_at" gorm:"not null;index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserStats represents user statistics
type UserStats struct {
	UserID         uint       `json:"user_id"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	AccountAgeDays int        `json:"account_age_days"`
	LoginCount     int        `json:"login_count,omitempty"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
}

// TableName specifies the table name for GORM
func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// TableName specifies the table name for GORM
func (LoginAttempt) TableName() string {
	return "login_attempts"
}
