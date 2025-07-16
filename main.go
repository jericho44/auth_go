package main

import (
	"log"
	"net/http"

	"auth-jwt/config"
	"auth-jwt/database"
	"auth-jwt/handlers"
	"auth-jwt/middleware"
	"auth-jwt/utils"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	database.Connect(cfg)
	defer database.DB.Close()

	// Initialize JWT with config
	utils.InitJWT(cfg)

	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.JWTMiddleware)
	protected.HandleFunc("/profile", handlers.Profile).Methods("GET")

	log.Printf("Server starting on :%s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
