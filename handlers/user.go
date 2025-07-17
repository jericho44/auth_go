package handlers

import (
	"encoding/json"
	"net/http"

	"auth-jwt/controllers"
	"auth-jwt/models"
	"auth-jwt/utils"
)

var userController = controllers.NewUserController()

func Profile(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by middleware)
	claims := r.Context().Value("user").(*utils.Claims)

	// Call controller
	response, err := userController.GetProfile(claims.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.User)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call controller
	response, err := userController.UpdateProfile(claims.UserID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "email is required" ||
			err.Error() == "invalid email format" {
			status = http.StatusBadRequest
		} else if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.User)
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)

	// Call controller
	response, err := userController.DeleteAccount(claims.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": response.Message,
	})
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)

	var req models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call controller
	response, err := userController.ChangePassword(claims.UserID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "current password and new password are required" ||
			err.Error() == "new password must be at least 6 characters long" ||
			err.Error() == "new password must be different from current password" {
			status = http.StatusBadRequest
		} else if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "current password is incorrect" {
			status = http.StatusForbidden
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": response.Message,
	})
}

func UserStats(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)

	// Call controller
	response, err := userController.GetUserStats(claims.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.User)
}
