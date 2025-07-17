package controllers

import (
	"errors"
	"strings"

	"auth-jwt/models"
	"auth-jwt/utils"

	"github.com/lib/pq"
)

type AuthController struct{}

type AuthResponse struct {
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty"`
	UserID  int    `json:"user_id,omitempty"`
}

// NewAuthController creates a new auth controller instance
func NewAuthController() *AuthController {
	return &AuthController{}
}

// Register handles user registration business logic
func (ac *AuthController) Register(req models.RegisterRequest) (*AuthResponse, error) {
	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("all fields are required")
	}

	// Additional validation
	if len(req.Username) < 3 {
		return nil, errors.New("username must be at least 3 characters long")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters long")
	}

	if !isValidEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Create new user
	user, err := models.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				if strings.Contains(pqErr.Detail, "username") {
					return nil, errors.New("username already exists")
				} else if strings.Contains(pqErr.Detail, "email") {
					return nil, errors.New("email already exists")
				} else {
					return nil, errors.New("user already exists")
				}
			}
		}
		return nil, errors.New("error creating user")
	}

	return &AuthResponse{
		Message: "User registered successfully",
		UserID:  user.ID,
	}, nil
}

// Login handles user authentication business logic
func (ac *AuthController) Login(req models.LoginRequest) (*AuthResponse, error) {
	// Validate input
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	// Find user
	user, err := models.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("error generating token")
	}

	return &AuthResponse{
		Token: token,
	}, nil
}

// ForgotPassword handles password reset request business logic
func (ac *AuthController) ForgotPassword(email string) (*AuthResponse, error) {
	// TODO: Implement password reset functionality
	// 1. Validate email exists
	// 2. Generate reset token
	// 3. Send email with reset link
	// 4. Store reset token with expiration

	return &AuthResponse{
		Message: "Password reset functionality coming soon",
	}, nil
}

// ResetPassword handles password reset confirmation business logic
func (ac *AuthController) ResetPassword(token, newPassword string) (*AuthResponse, error) {
	// TODO: Implement password reset confirmation
	// 1. Validate reset token
	// 2. Check token expiration
	// 3. Update user password
	// 4. Invalidate reset token

	return &AuthResponse{
		Message: "Password reset functionality coming soon",
	}, nil
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
