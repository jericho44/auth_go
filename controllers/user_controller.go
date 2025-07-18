package controllers

import (
	"errors"

	"auth-jwt/models"
	"auth-jwt/services"
	"auth-jwt/utils"
)

type UserController struct {
	userService services.UserServiceInterface
}

type UserResponse struct {
	Message string            `json:"message,omitempty"`
	User    *models.User      `json:"user,omitempty"`
	Stats   *models.UserStats `json:"stats,omitempty"`
}

// NewUserController creates a new user controller instance
func NewUserController(userService services.UserServiceInterface) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetProfile retrieves user profile business logic
func (uc *UserController) GetProfile(userID int) (*UserResponse, error) {
	user, err := uc.userService.GetUserByID(uint(userID))
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

	if !utils.IsValidEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Update user profile through service
	user, err := uc.userService.UpdateUserProfile(uint(userID), req.Email)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, errors.New("user not found")
		}
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

	// Get current user to verify password
	user, err := uc.userService.GetUserByID(uint(userID))
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

	// Update password through service
	err = uc.userService.UpdateUserPassword(uint(userID), req.NewPassword)
	if err != nil {
		return nil, errors.New("error updating password")
	}

	return &UserResponse{
		Message: "Password changed successfully",
	}, nil
}

// DeleteAccount handles account deletion business logic
func (uc *UserController) DeleteAccount(userID int) (*UserResponse, error) {
	// Delete user account through service
	err := uc.userService.DeleteUser(uint(userID))
	if err != nil {
		if err.Error() == "user not found" {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("error deleting account")
	}

	return &UserResponse{
		Message: "Account deleted successfully",
	}, nil
}

// GetUserStats retrieves user statistics
func (uc *UserController) GetUserStats(userID int) (*UserResponse, error) {
	stats, err := uc.userService.GetUserStats(uint(userID))
	if err != nil {
		return nil, errors.New("error retrieving user stats")
	}

	return &UserResponse{
		Stats: stats,
	}, nil
}
