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

// @Summary User registration
// @Description Register a new user account with username, email, and password
// @Description
// @Description **Registration Requirements:**
// @Description - Username: minimum 3 characters, alphanumeric only
// @Description - Email: valid email format, must be unique in system
// @Description - Password: minimum 6 characters, will be securely hashed
// @Description
// @Description **Registration Process:**
// @Description 1. Input validation (format, length, uniqueness checks)
// @Description 2. Password hashing with bcrypt for security
// @Description 3. User account creation in database
// @Description 4. Automatic JWT token generation for immediate login
// @Description
// @Description **Response includes:**
// @Description - JWT access token for authentication
// @Description - User profile information (without password)
// @Description - Token expiration information
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User registration data"
// @Success 201 {object} controllers.AuthResponse "User registered successfully with JWT token"
// @Failure 400 {object} utils.ErrorResponse "Bad request - validation errors"
// @Failure 409 {object} utils.ErrorResponse "Conflict - username or email already exists"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to create user account"
// @Router /auth/register [post]
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

// @Summary User login
// @Description Authenticate user with username/email and password to receive JWT token
// @Description
// @Description **Login Options:**
// @Description - Login with username OR email
// @Description - Password must match the registered password
// @Description
// @Description **Security Features:**
// @Description - Rate limiting for failed attempts
// @Description - Password verification with bcrypt
// @Description - JWT token generation with configurable expiration
// @Description
// @Description **Response includes:**
// @Description - JWT access token for API authentication
// @Description - User profile information
// @Description - Token expiration timestamp
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "User login credentials"
// @Success 200 {object} controllers.AuthResponse "Login successful"
// @Failure 400 {object} utils.ErrorResponse "Bad request - missing credentials"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - invalid credentials"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auth/login [post]
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
// @Summary Request password reset
// @Description Send password reset token to user's email address
// @Description
// @Description **Process:**
// @Description 1. Validate email exists in system
// @Description 2. Generate secure reset token
// @Description 3. Store token with expiration (typically 1 hour)
// @Description 4. Send reset email with token
// @Description
// @Description **Security:**
// @Description - Tokens expire after 1 hour
// @Description - One-time use tokens
// @Description - Rate limiting to prevent abuse
// @Tags Authentication
// @Accept json
// @Produce json
// @Param email body models.ForgotPasswordRequest true "Email address for password reset"
// @Success 200 {object} controllers.AuthResponse "Password reset email sent"
// @Failure 400 {object} utils.ErrorResponse "Bad request - invalid email format"
// @Failure 404 {object} utils.ErrorResponse "Email not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to send email"
// @Router /auth/forgot-password [post]
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
// @Summary Reset password with token
// @Description Reset user password using the token received via email
// @Description
// @Description **Process:**
// @Description 1. Validate reset token (exists, not expired, not used)
// @Description 2. Validate new password requirements
// @Description 3. Hash new password with bcrypt
// @Description 4. Update user password
// @Description 5. Invalidate reset token
// @Description
// @Description **Requirements:**
// @Description - Valid reset token from forgot-password request
// @Description - New password minimum 6 characters
// @Description - Token must not be expired or already used
// @Tags Authentication
// @Accept json
// @Produce json
// @Param reset body models.ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} controllers.AuthResponse "Password reset successful"
// @Failure 400 {object} utils.ErrorResponse "Bad request - invalid token or password requirements"
// @Failure 404 {object} utils.ErrorResponse "Invalid or expired reset token"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auth/reset-password [post]
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
