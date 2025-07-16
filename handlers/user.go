package handlers

import (
	"encoding/json"
	"net/http"

	"auth-jwt/models"
	"auth-jwt/utils"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by middleware)
	claims := r.Context().Value("user").(*utils.Claims)

	user, err := models.GetUserByID(claims.UserID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update user profile
	user, err := models.UpdateUserProfile(claims.UserID, req.Email)
	if err != nil {
		http.Error(w, "Error updating profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)

	err := models.DeleteUser(claims.UserID)
	if err != nil {
		http.Error(w, "Error deleting account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Account deleted successfully",
	})
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)

	var req models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get current user
	user, err := models.GetUserByID(claims.UserID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify current password
	if !user.CheckPassword(req.CurrentPassword) {
		http.Error(w, "Current password is incorrect", http.StatusForbidden)
		return
	}

	// Update password
	err = models.UpdateUserPassword(claims.UserID, req.NewPassword)
	if err != nil {
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password changed successfully",
	})
}
