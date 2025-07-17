package handlers

import (
	"encoding/json"
	"net/http"

	"auth-jwt/controllers"
	"auth-jwt/models"
	"auth-jwt/utils"
)

var authController *controllers.AuthController

// SetAuthController sets the auth controller instance
func SetAuthController(controller *controllers.AuthController) {
	authController = controller
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", "Failed to parse JSON request")
		return
	}

	// Call controller
	response, err := authController.Register(req)
	if err != nil {
		switch err.Error() {
		case "all fields are required", "username must be at least 3 characters long",
			"password must be at least 6 characters long", "invalid email format":
			utils.ValidationError(w, err.Error(), "Please check your input and try again")
		case "username already exists", "email already exists", "user already exists":
			utils.Conflict(w, err.Error())
		default:
			utils.InternalServerError(w, "Registration failed")
		}
		return
	}

	utils.Created(w, "User registered successfully", response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", "Failed to parse JSON request")
		return
	}

	// Call controller
	response, err := authController.Login(req)
	if err != nil {
		switch err.Error() {
		case "username and password are required":
			utils.ValidationError(w, err.Error(), "Both username and password must be provided")
		case "invalid credentials":
			utils.Unauthorized(w, "Invalid username or password")
		default:
			utils.InternalServerError(w, "Login failed")
		}
		return
	}

	utils.Success(w, "Login successful", response)
}

// ForgotPassword handles password reset requests
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", "Failed to parse JSON request")
		return
	}

	response, err := authController.ForgotPassword(req.Email)
	if err != nil {
		switch err.Error() {
		case "email is required", "invalid email format":
			utils.ValidationError(w, err.Error(), "Please provide a valid email address")
		default:
			utils.InternalServerError(w, "Password reset request failed")
		}
		return
	}

	utils.Success(w, "Password reset request processed", response)
}

// ResetPassword handles password reset confirmation
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", "Failed to parse JSON request")
		return
	}

	response, err := authController.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		switch err.Error() {
		case "token and new password are required", "new password must be at least 6 characters long":
			utils.ValidationError(w, err.Error(), "Please check your input and try again")
		case "invalid or expired token":
			utils.BadRequest(w, "Invalid or expired reset token", "Please request a new password reset")
		default:
			utils.InternalServerError(w, "Password reset failed")
		}
		return
	}

	utils.Success(w, "Password reset successful", response)
}
