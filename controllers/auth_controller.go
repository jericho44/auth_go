package controllers

import (
	"errors"
	"strings"

	"auth-jwt/models"
	"auth-jwt/services"
	"auth-jwt/utils"
)

type AuthController struct {
	authService services.AuthServiceInterface
}

type AuthResponse struct {
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty"`
	UserID  int    `json:"user_id,omitempty"`
}

// NewAuthController creates a new auth controller instance
func NewAuthController(authService services.AuthServiceInterface) *AuthController {
	return &AuthController{
		authService: authService,
	}
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

	if !utils.IsValidEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Register user through service
	user, err := ac.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			if strings.Contains(strings.ToLower(err.Error()), "username") {
				return nil, errors.New("username already exists")
			} else if strings.Contains(strings.ToLower(err.Error()), "email") {
				return nil, errors.New("email already exists")
			}
			return nil, errors.New("user already exists")
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

	// TODO: Add IP address from request context for rate limiting
	// if err := ac.authService.ValidateLoginAttempts(req.Username, ipAddress); err != nil {
	//     return nil, err
	// }

	// Login through service
	token, err := ac.authService.Login(req.Username, req.Password)
	if err != nil {
		// TODO: Record failed login attempt
		// ac.authService.RecordLoginAttempt(req.Username, ipAddress, false)
		return nil, err
	}

	// TODO: Record successful login attempt
	// ac.authService.RecordLoginAttempt(req.Username, ipAddress, true)

	return &AuthResponse{
		Token: token,
	}, nil
}

// ForgotPassword handles password reset request business logic
func (ac *AuthController) ForgotPassword(email string) (*AuthResponse, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	if !utils.IsValidEmail(email) {
		return nil, errors.New("invalid email format")
	}

	// Generate reset token through service
	token, err := ac.authService.GeneratePasswordResetToken(email)
	if err != nil {
		if err.Error() == "user not found" {
			// Don't reveal if email exists or not for security
			return &AuthResponse{
				Message: "If the email exists, a password reset link has been sent",
			}, nil
		}
		return nil, errors.New("error generating reset token")
	}

	// TODO: Send email with reset token
	// In a real application, you would send an email here

	return &AuthResponse{
		Message: "Password reset token generated: " + token, // Remove this in production
	}, nil
}

// ResetPassword handles password reset confirmation business logic
func (ac *AuthController) ResetPassword(token, newPassword string) (*AuthResponse, error) {
	if token == "" || newPassword == "" {
		return nil, errors.New("token and new password are required")
	}

	if len(newPassword) < 6 {
		return nil, errors.New("new password must be at least 6 characters long")
	}

	// Reset password through service
	if err := ac.authService.ResetPassword(token, newPassword); err != nil {
		return nil, err
	}

	return &AuthResponse{
		Message: "Password reset successfully",
	}, nil
}
