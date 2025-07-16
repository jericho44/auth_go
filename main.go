package main

import (
	"log"
	"net/http"

	"auth-jwt/config"
	"auth-jwt/database"
	_ "auth-jwt/docs"
	"auth-jwt/handlers"
	"auth-jwt/middleware"
	"auth-jwt/utils"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Authentication API
// @version 1.0
// @description JWT Authentication API with user registration and login
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	database.Connect(cfg)
	defer database.DB.Close()

	// Initialize JWT with config
	utils.InitJWT(cfg)

	r := mux.NewRouter()

	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Public routes
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.JWTMiddleware)
	protected.HandleFunc("/profile", handlers.Profile).Methods("GET")

	log.Printf("Server starting on :%s", cfg.ServerPort)
	log.Printf("Swagger documentation available at: http://localhost:%s/swagger/", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
