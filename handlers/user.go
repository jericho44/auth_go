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
