package routes

import (
	"auth-jwt/handlers"

	"github.com/gorilla/mux"
)

// AuthRoutes defines authentication-related routes
type AuthRoutes struct {
	router *mux.Router
}

// NewAuthRoutes creates a new auth routes instance
func NewAuthRoutes() *AuthRoutes {
	return &AuthRoutes{
		router: mux.NewRouter(),
	}
}

// RegisterRoutes registers all authentication routes
func (ar *AuthRoutes) RegisterRoutes(r *mux.Router) {
	// Authentication routes group
	auth := r.PathPrefix("/auth").Subrouter()

	// @Summary User registration
	// @Description Create a new user account
	// @Tags Authentication
	// @Accept json
	// @Produce json
	// @Param user body models.RegisterRequest true "User registration data"
	// @Success 200 {object} map[string]interface{} "Registration successful"
	// @Failure 400 {string} string "Invalid request body or missing fields"
	// @Failure 409 {string} string "User already exists"
	// @Router /auth/register [post]
	auth.HandleFunc("/register", handlers.Register).Methods("POST")

	// @Summary User login
	// @Description Authenticate user and return JWT token
	// @Tags Authentication
	// @Accept json
	// @Produce json
	// @Param credentials body models.LoginRequest true "Login credentials"
	// @Success 200 {object} map[string]string "Login successful"
	// @Failure 400 {string} string "Invalid request body"
	// @Failure 401 {string} string "Invalid credentials"
	// @Router /auth/login [post]
	auth.HandleFunc("/login", handlers.Login).Methods("POST")

	// Password reset routes (future implementation)
	auth.HandleFunc("/forgot-password", handlers.ForgotPassword).Methods("POST")
	auth.HandleFunc("/reset-password", handlers.ResetPassword).Methods("POST")
}
