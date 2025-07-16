package routes

import (
	"auth-jwt/middleware"
	"encoding/json"
	"net/http"
	"time"

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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	}).Methods("GET")
}
