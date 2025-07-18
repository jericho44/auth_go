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
	setupApiPublicRoutes(r)

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
	auth.HandleFunc("/forgot-password", handlers.ForgotPassword).Methods("POST")
	auth.HandleFunc("/reset-password", handlers.ResetPassword).Methods("POST")

	// Legacy routes for backward compatibility
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")
}

// setupAPIRoutes configures API routes (both public and protected)
func setupAPIRoutes(r *mux.Router) {
	// Add global middleware
	r.Use(middleware.RequestIDMiddleware)
	r.Use(middleware.ResponseHeadersMiddleware)
	r.Use(middleware.CORSMiddleware)
	r.Use(middleware.LoggingMiddleware)

	// Create API router
	api := r.PathPrefix("/api").Subrouter()

	// Setup public file routes directly on the API router (no authentication required)
	// These routes will be accessible without authentication
	// api.HandleFunc("/files/public", handlers.GetPublicFiles).Methods("GET")
	// api.HandleFunc("/files/download/{filename}", handlers.DownloadFile).Methods("GET")
	// api.HandleFunc("/files/type/{type}", handlers.GetFilesByType).Methods("GET")
	// api.HandleFunc("/files/{id:[0-9]+}", handlers.GetFile).Methods("GET")

	// Create protected subrouter for routes that require authentication
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.JWTMiddleware)

	// User routes
	setupUserRoutes(protected)

	// Protected file routes
	files := protected.PathPrefix("/files").Subrouter()
	files.HandleFunc("/upload", handlers.UploadSingleFile).Methods("POST")
	files.HandleFunc("/upload-multiple", handlers.UploadMultipleFiles).Methods("POST")
	files.HandleFunc("/my", handlers.GetUserFiles).Methods("GET")
	files.HandleFunc("/{id:[0-9]+}", handlers.UpdateFile).Methods("PUT")
	files.HandleFunc("/{id:[0-9]+}", handlers.DeleteFile).Methods("DELETE")
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

func setupApiPublicRoutes(r *mux.Router) {

	// Public authentication routes
	api := r.PathPrefix("/").Subrouter()
	api.HandleFunc("/files/public", handlers.GetPublicFiles).Methods("GET")
	api.HandleFunc("/files/download/{filename}", handlers.DownloadFile).Methods("GET")
	api.HandleFunc("/files/type/{type}", handlers.GetFilesByType).Methods("GET")
	api.HandleFunc("/files/{id:[0-9]+}", handlers.GetFile).Methods("GET")
}
