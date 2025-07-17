package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"auth-jwt/models"
	"auth-jwt/repositories"
	"auth-jwt/utils"
)

// AuthServiceInterface defines the contract for authentication business operations
type AuthServiceInterface interface {
	Register(username, email, password string) (*models.User, error)
	Login(username, password string) (string, error)
	GeneratePasswordResetToken(email string) (string, error)
	ResetPassword(token, newPassword string) error
	ValidateLoginAttempts(username, ipAddress string) error
	RecordLoginAttempt(username, ipAddress string, success bool) error
}

// AuthService implements AuthServiceInterface
type AuthService struct {
	userRepo repositories.UserRepositoryInterface
	authRepo repositories.AuthRepositoryInterface
}

// NewAuthService creates a new auth service instance
func NewAuthService(userRepo repositories.UserRepositoryInterface, authRepo repositories.AuthRepositoryInterface) AuthServiceInterface {
	return &AuthService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

// Register handles user registration
func (as *AuthService) Register(username, email, password string) (*models.User, error) {
	// Check if user already exists
	exists, err := as.userRepo.Exists(username, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	// Create user (password will be hashed by GORM hook)
	user := &models.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	// Save to database
	if err := as.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login handles user authentication
func (as *AuthService) Login(username, password string) (string, error) {
	// Get user by username
	user, err := as.userRepo.GetByUsername(username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	// Check password
	if !user.CheckPassword(password) {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(int(user.ID), user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GeneratePasswordResetToken creates a password reset token
func (as *AuthService) GeneratePasswordResetToken(email string) (string, error) {
	// Find user by email
	user, err := as.userRepo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	// Set expiration time (1 hour from now)
	expiresAt := time.Now().Add(time.Hour)

	// Save token to database
	if err := as.authRepo.CreatePasswordResetToken(user.ID, token, expiresAt); err != nil {
		return "", err
	}

	return token, nil
}

// ResetPassword resets user password using token
func (as *AuthService) ResetPassword(token, newPassword string) error {
	// Get reset token
	resetToken, err := as.authRepo.GetPasswordResetToken(token)
	if err != nil {
		return err
	}
	if resetToken == nil {
		return errors.New("invalid or expired token")
	}

	// Hash new password
	tempUser := &models.User{Password: newPassword}
	if err := tempUser.HashPassword(); err != nil {
		return err
	}

	// Update user password
	if err := as.userRepo.UpdatePassword(resetToken.UserID, tempUser.Password); err != nil {
		return err
	}

	// Delete used token
	if err := as.authRepo.DeletePasswordResetToken(token); err != nil {
		return err
	}

	return nil
}

// ValidateLoginAttempts checks for too many failed login attempts
func (as *AuthService) ValidateLoginAttempts(username, ipAddress string) error {
	// Check attempts in the last 15 minutes
	since := time.Now().Add(-15 * time.Minute)
	attempts, err := as.authRepo.GetRecentLoginAttempts(username, since)
	if err != nil {
		return err
	}

	// Allow maximum 5 failed attempts in 15 minutes
	if attempts >= 5 {
		return errors.New("too many failed login attempts, please try again later")
	}

	return nil
}

// RecordLoginAttempt logs a login attempt
func (as *AuthService) RecordLoginAttempt(username, ipAddress string, success bool) error {
	return as.authRepo.CreateLoginAttempt(username, ipAddress, success)
}
