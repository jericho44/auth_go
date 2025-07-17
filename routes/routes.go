package routes

import (
	"auth-jwt/handlers"
	"auth-jwt/middleware"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupRoutes configures all application routes
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Setup route groups
	setupSwaggerRoutes(r)
	setupAuthRoutes(r)
	setupAPIRoutes(r)

	return r
}

// setupSwaggerRoutes configures Swagger documentation routes
func setupSwaggerRoutes(r *mux.Router) {
	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}

// setupAuthRoutes configures public authentication routes
func setupAuthRoutes(r *mux.Router) {
	// Public authentication routes
	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", handlers.Register).Methods("POST")
	auth.HandleFunc("/login", handlers.Login).Methods("POST")

	// Legacy routes for backward compatibility
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")
}

// setupAPIRoutes configures protected API routes
func setupAPIRoutes(r *mux.Router) {
	// Protected API routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware)
	api.Use(middleware.CORSMiddleware)
	api.Use(middleware.LoggingMiddleware)

	// User routes
	setupUserRoutes(api)
}

// setupUserRoutes configures user-related routes
func setupUserRoutes(api *mux.Router) {
	user := api.PathPrefix("/user").Subrouter()
	user.HandleFunc("/profile", handlers.Profile).Methods("GET")
	user.HandleFunc("/profile", handlers.UpdateProfile).Methods("PUT")
	user.HandleFunc("/change-password", handlers.ChangePassword).Methods("POST")
	user.HandleFunc("/account", handlers.DeleteAccount).Methods("DELETE")
	user.HandleFunc("/stats", handlers.UserStats).Methods("GET")

	// Legacy route for backward compatibility
	api.HandleFunc("/profile", handlers.Profile).Methods("GET")
}
