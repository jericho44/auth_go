package handlers

import (
	"encoding/json"
	"net/http"

	"auth-jwt/controllers"
	"auth-jwt/models"
)

var authController = controllers.NewAuthController()

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call controller
	response, err := authController.Register(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "all fields are required" ||
			err.Error() == "username must be at least 3 characters long" ||
			err.Error() == "password must be at least 6 characters long" ||
			err.Error() == "invalid email format" {
			status = http.StatusBadRequest
		} else if err.Error() == "username already exists" ||
			err.Error() == "email already exists" ||
			err.Error() == "user already exists" {
			status = http.StatusConflict
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call controller
	response, err := authController.Login(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "username and password are required" {
			status = http.StatusBadRequest
		} else if err.Error() == "invalid credentials" {
			status = http.StatusUnauthorized
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ForgotPassword handles password reset requests
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := authController.ForgotPassword(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ResetPassword handles password reset confirmation
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := authController.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
