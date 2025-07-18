package routes

import (
	"auth-jwt/middleware"
	"net/http"
	"time"

	"auth-jwt/utils"

	"github.com/gorilla/mux"
)

// APIRoutes defines API route groups
type APIRoutes struct {
	userRoutes *UserRoutes
}

// NewAPIRoutes creates a new API routes instance
func NewAPIRoutes() *APIRoutes {
	return &APIRoutes{
		userRoutes: NewUserRoutes(),
	}
}

// RegisterRoutes registers all API routes with middleware
func (ar *APIRoutes) RegisterRoutes(r *mux.Router) {
	// API v1 routes
	v1 := r.PathPrefix("/api/v1").Subrouter()
	v1.Use(middleware.JWTMiddleware)
	v1.Use(middleware.CORSMiddleware)
	v1.Use(middleware.LoggingMiddleware)

	// Register user routes under v1
	ar.userRoutes.RegisterRoutes(v1)

	// Health check endpoint (no auth required)
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		healthData := map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
			"uptime":    time.Since(time.Now()).String(), // This would be calculated from app start time
		}
		utils.Success(w, "Service is healthy", healthData)
	}).Methods("GET")
}
