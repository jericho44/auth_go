package handlers

import (
	"encoding/json"
	"net/http"

	"auth-jwt/controllers"
	"auth-jwt/models"
	"auth-jwt/utils"
)

var userController *controllers.UserController

// SetUserController sets the user controller instance
func SetUserController(controller *controllers.UserController) {
	userController = controller
}

// @Summary Get user profile
// @Description Retrieve the authenticated user's profile information including username, email, and account details
// @Description
// @Description **Profile Information:**
// @Description - User ID and username
// @Description - Email address
// @Description - Account creation and last update timestamps
// @Description - Account status and verification details
// @Description
// @Description **Authentication Required:**
// @Description - Valid JWT token must be provided in Authorization header
// @Description - Token format: "Bearer {your-jwt-token}"
// @Tags User Management
// @Produce json
// @Success 200 {object} controllers.UserResponse "User profile retrieved successfully"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - invalid or missing JWT token"
// @Failure 404 {object} utils.ErrorResponse "User profile not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to retrieve profile"
// @Security BearerAuth
// @Router /api/user/profile [get]
func Profile(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by middleware)
	claims := r.Context().Value("user").(*utils.Claims)

	// Call controller
	response, err := userController.GetProfile(claims.UserID)
	if err != nil {
		switch err.Error() {
		case "user not found":
			utils.NotFound(w, "User profile not found")
		default:
			utils.InternalServerError(w, "Failed to retrieve user profile")
		}
		return
	}

	utils.Success(w, "User profile retrieved successfully", response.User)
}

// @Summary Update user profile
// @Description Update the authenticated user's profile information such as email address
// @Description
// @Description **Updatable Fields:**
// @Description - Email address (must be valid format and unique)
// @Description - Additional profile fields as supported
// @Description
// @Description **Update Process:**
// @Description 1. Validate new email format and uniqueness
// @Description 2. Update user record in database
// @Description 3. Return updated profile information
// @Description
// @Description **Authentication Required:**
// @Description - Valid JWT token must be provided in Authorization header
// @Tags User Management
// @Accept json
// @Produce json
// @Param profile body models.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} controllers.UserResponse "Profile updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Bad request - invalid email format or validation errors"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - invalid or missing JWT token"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 409 {object} utils.ErrorResponse "Conflict - email already exists"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to update profile"
// @Security BearerAuth
// @Router /api/user/profile [put]
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", "Failed to parse JSON request")
		return
	}

	// Call controller
	response, err := userController.UpdateProfile(claims.UserID, req)
	if err != nil {
		switch err.Error() {
		case "email is required", "invalid email format":
			utils.ValidationError(w, err.Error(), "Please provide a valid email address")
		case "user not found":
			utils.NotFound(w, "User not found")
		default:
			utils.InternalServerError(w, "Failed to update profile")
		}
		return
	}

	utils.Success(w, "Profile updated successfully", response.User)
}

// @Summary Delete user account
// @Description Permanently delete the authenticated user's account and all associated data
// @Description
// @Description **Deletion Process:**
// @Description 1. Verify user authentication and ownership
// @Description 2. Remove all user-associated data (files, records, etc.)
// @Description 3. Permanently delete user account from database
// @Description 4. Invalidate all existing JWT tokens for the user
// @Description
// @Description **Warning:**
// @Description - This action is irreversible
// @Description - All user data will be permanently lost
// @Description - All uploaded files will be deleted
// @Description
// @Description **Authentication Required:**
// @Description - Valid JWT token must be provided in Authorization header
// @Tags User Management
// @Produce json
// @Success 200 {object} controllers.UserResponse "Account deleted successfully"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - invalid or missing JWT token"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to delete account"
// @Security BearerAuth
// @Router /api/user/account [delete]
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)

	// Call controller
	response, err := userController.DeleteAccount(claims.UserID)
	if err != nil {
		switch err.Error() {
		case "user not found":
			utils.NotFound(w, "User not found")
		default:
			utils.InternalServerError(w, "Failed to delete account")
		}
		return
	}

	utils.Success(w, response.Message, nil)
}

// @Summary Change user password
// @Description Change the authenticated user's password with current password verification
// @Description
// @Description **Password Requirements:**
// @Description - Current password must be provided for verification
// @Description - New password must be at least 6 characters long
// @Description - New password must be different from current password
// @Description - New password will be securely hashed with bcrypt
// @Description
// @Description **Change Process:**
// @Description 1. Verify current password is correct
// @Description 2. Validate new password meets requirements
// @Description 3. Hash new password securely
// @Description 4. Update password in database
// @Description 5. Optionally invalidate existing JWT tokens
// @Description
// @Description **Authentication Required:**
// @Description - Valid JWT token must be provided in Authorization header
// @Tags User Management
// @Accept json
// @Produce json
// @Param password body models.ChangePasswordRequest true "Password change data"
// @Success 200 {object} controllers.UserResponse "Password changed successfully"
// @Failure 400 {object} utils.ErrorResponse "Bad request - validation errors or password requirements not met"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - invalid or missing JWT token"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - current password is incorrect"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to change password"
// @Security BearerAuth
// @Router /api/user/change-password [post]
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)

	var req models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", "Failed to parse JSON request")
		return
	}

	// Call controller
	response, err := userController.ChangePassword(claims.UserID, req)
	if err != nil {
		switch err.Error() {
		case "current password and new password are required",
			"new password must be at least 6 characters long",
			"new password must be different from current password":
			utils.ValidationError(w, err.Error(), "Please check your password requirements")
		case "user not found":
			utils.NotFound(w, "User not found")
		case "current password is incorrect":
			utils.Forbidden(w, "Current password is incorrect")
		default:
			utils.InternalServerError(w, "Failed to change password")
		}
		return
	}

	utils.Success(w, response.Message, nil)
}

// @Summary Get user statistics
// @Description Retrieve comprehensive statistics about the authenticated user's account and activity
// @Description
// @Description **Statistics Include:**
// @Description - Total number of uploaded files
// @Description - Total storage space used
// @Description - Account creation date and activity metrics
// @Description - File type breakdown (images, documents, etc.)
// @Description - Public vs private file counts
// @Description
// @Description **Use Cases:**
// @Description - Dashboard displays
// @Description - Storage usage monitoring
// @Description - Account activity tracking
// @Description - Usage analytics
// @Description
// @Description **Authentication Required:**
// @Description - Valid JWT token must be provided in Authorization header
// @Tags User Management
// @Produce json
// @Success 200 {object} controllers.UserResponse "User statistics retrieved successfully"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - invalid or missing JWT token"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to retrieve statistics"
// @Security BearerAuth
// @Router /api/user/stats [get]
func UserStats(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)

	// Call controller
	response, err := userController.GetUserStats(claims.UserID)
	if err != nil {
		switch err.Error() {
		case "user not found":
			utils.NotFound(w, "User not found")
		default:
			utils.InternalServerError(w, "Failed to retrieve user statistics")
		}
		return
	}

	utils.Success(w, "User statistics retrieved successfully", response.Stats)
}
