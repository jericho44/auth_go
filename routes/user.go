package routes

import (
	"auth-jwt/handlers"
	"auth-jwt/middleware"

	"github.com/gorilla/mux"
)

// UserRoutes defines user-related routes
type UserRoutes struct {
	router *mux.Router
}

// NewUserRoutes creates a new user routes instance
func NewUserRoutes() *UserRoutes {
	return &UserRoutes{
		router: mux.NewRouter(),
	}
}

// RegisterRoutes registers all user routes (protected)
func (ur *UserRoutes) RegisterRoutes(r *mux.Router) {
	// Protected user routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware)

	// User management routes
	user := api.PathPrefix("/user").Subrouter()
	
	// @Summary Get user profile
	// @Description Get current user profile information
	// @Tags User
	// @Produce json
	// @Security BearerAuth
	// @Success 200 {object} models.User "User profile"
	// @Failure 401 {string} string "Unauthorized - Invalid or missing token"
	// @Failure 404 {string} string "User not found"
	// @Router /api/user/profile [get]
	user.HandleFunc("/profile", handlers.Profile).Methods("GET")
	
	// @Summary Update user profile
	// @Description Update current user profile information
	// @Tags User
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param user body models.UpdateProfileRequest true "Profile update data"
	// @Success 200 {object} models.User "Updated user profile"
	// @Failure 400 {string} string "Invalid request body"
	// @Failure 401 {string} string "Unauthorized - Invalid or missing token"
	// @Failure 404 {string} string "User not found"
	// @Router /api/user/profile [put]
	user.HandleFunc("/profile", handlers.UpdateProfile).Methods("PUT")
	
	// @Summary Delete user account
	// @Description Delete current user account
	// @Tags User
	// @Security BearerAuth
	// @Success 200 {object} map[string]string "Account deleted successfully"
	// @Failure 401 {string} string "Unauthorized - Invalid or missing token"
	// @Failure 404 {string} string "User not found"
	// @Router /api/user/account [delete]
	user.HandleFunc("/account", handlers.DeleteAccount).Methods("DELETE")

	// @Summary Change password
	// @Description Change user password
	// @Tags User
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param passwords body models.ChangePasswordRequest true "Password change data"
	// @Success 200 {object} map[string]string "Password changed successfully"
	// @Failure 400 {string} string "Invalid request body"
	// @Failure 401 {string} string "Unauthorized - Invalid or missing token"
	// @Failure 403 {string} string "Current password is incorrect"
	// @Router /api/user/change-password [post]
	user.HandleFunc("/change-password", handlers.ChangePassword).Methods("POST")
}