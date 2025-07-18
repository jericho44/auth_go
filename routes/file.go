package routes

import (
	"auth-jwt/handlers"
	"auth-jwt/middleware"

	"github.com/gorilla/mux"
)

func SetupFileRoutes(r *mux.Router) {
	// Public routes
	r.HandleFunc("/files/{id:[0-9]+}", handlers.GetFile).Methods("GET")
	r.HandleFunc("/files/public", handlers.GetPublicFiles).Methods("GET")
	r.HandleFunc("/files/download/{filename}", handlers.DownloadFile).Methods("GET")
	r.HandleFunc("/files/type/{type}", handlers.GetFilesByType).Methods("GET")

	// Protected routes
	protected := r.PathPrefix("/files").Subrouter()
	protected.Use(middleware.JWTMiddleware)

	// Single file upload
	protected.HandleFunc("/upload", handlers.UploadSingleFile).Methods("POST")

	// Multiple files upload
	protected.HandleFunc("/upload-multiple", handlers.UploadMultipleFiles).Methods("POST")

	// User's files
	protected.HandleFunc("/my", handlers.GetUserFiles).Methods("GET")

	// File management
	protected.HandleFunc("/{id:[0-9]+}", handlers.UpdateFile).Methods("PUT")
	protected.HandleFunc("/{id:[0-9]+}", handlers.DeleteFile).Methods("DELETE")
}

// SetupPublicFileRoutes configures public file routes (no authentication required)
func SetupPublicFileRoutes(r *mux.Router) {
	r.HandleFunc("/files/{id:[0-9]+}", handlers.GetFile).Methods("GET")
	r.HandleFunc("/files/public", handlers.GetPublicFiles).Methods("GET")
	r.HandleFunc("/files/download/{filename}", handlers.DownloadFile).Methods("GET")
	r.HandleFunc("/files/type/{type}", handlers.GetFilesByType).Methods("GET")
}

// SetupProtectedFileRoutes configures protected file routes (authentication required)
func SetupProtectedFileRoutes(r *mux.Router) {
	files := r.PathPrefix("/files").Subrouter()

	// Single file upload
	files.HandleFunc("/upload", handlers.UploadSingleFile).Methods("POST")

	// Multiple files upload
	files.HandleFunc("/upload-multiple", handlers.UploadMultipleFiles).Methods("POST")

	// User's files
	files.HandleFunc("/my", handlers.GetUserFiles).Methods("GET")

	// File management
	files.HandleFunc("/{id:[0-9]+}", handlers.UpdateFile).Methods("PUT")
	files.HandleFunc("/{id:[0-9]+}", handlers.DeleteFile).Methods("DELETE")
}
