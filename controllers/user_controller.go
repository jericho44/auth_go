package controllers

import (
	"errors"

	"auth-jwt/models"
)

type UserController struct{}

type UserResponse struct {
	Message string       `json:"message,omitempty"`
	User    *models.User `json:"user,omitempty"`
}

// NewUserController creates a new user controller instance
func NewUserController() *UserController {
	return &UserController{}
}

// GetProfile retrieves user profile business logic
func (uc *UserController) GetProfile(userID int) (*UserResponse, error) {
	user, err := models.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return &UserResponse{
		User: user,
	}, nil
}

// UpdateProfile handles profile update business logic
func (uc *UserController) UpdateProfile(userID int, req models.UpdateProfileRequest) (*UserResponse, error) {
	// Validate input
	if req.Email == "" {
		return nil, errors.New("email is required")
	}

	if !isValidEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Check if user exists
	existingUser, err := models.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	// Update user profile
	user, err := models.UpdateUserProfile(userID, req.Email)
	if err != nil {
		return nil, errors.New("error updating profile")
	}

	return &UserResponse{
		User: user,
	}, nil
}

// ChangePassword handles password change business logic
func (uc *UserController) ChangePassword(userID int, req models.ChangePasswordRequest) (*UserResponse, error) {
	// Validate input
	if req.CurrentPassword == "" || req.NewPassword == "" {
		return nil, errors.New("current password and new password are required")
	}

	if len(req.NewPassword) < 6 {
		return nil, errors.New("new password must be at least 6 characters long")
	}

	if req.CurrentPassword == req.NewPassword {
		return nil, errors.New("new password must be different from current password")
	}

	// Get current user
	user, err := models.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Verify current password
	if !user.CheckPassword(req.CurrentPassword) {
		return nil, errors.New("current password is incorrect")
	}

	// Update password
	err = models.UpdateUserPassword(userID, req.NewPassword)
	if err != nil {
		return nil, errors.New("error updating password")
	}

	return &UserResponse{
		Message: "Password changed successfully",
	}, nil
}

// DeleteAccount handles account deletion business logic
func (uc *UserController) DeleteAccount(userID int) (*UserResponse, error) {
	// Check if user exists
	user, err := models.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Delete user account
	err = models.DeleteUser(userID)
	if err != nil {
		return nil, errors.New("error deleting account")
	}

	return &UserResponse{
		Message: "Account deleted successfully",
	}, nil
}

// GetUserStats retrieves user statistics (example of additional business logic)
func (uc *UserController) GetUserStats(userID int) (*UserResponse, error) {
	user, err := models.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Add additional stats logic here
	// For example: login count, last login, account age, etc.

	return &UserResponse{
		User: user,
	}, nil
}
