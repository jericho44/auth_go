package main

import (
	"log"
	"net/http"

	"auth-jwt/config"
	"auth-jwt/database"
	_ "auth-jwt/docs"
	"auth-jwt/routes"
	"auth-jwt/utils"
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

	// Setup routes
	r := routes.SetupRoutes()

	log.Printf("Server starting on :%s", cfg.ServerPort)
	log.Printf("Swagger documentation available at: http://localhost:%s/swagger/", cfg.ServerPort)
	log.Printf("API endpoints:")
	log.Printf("  POST /auth/register - User registration")
	log.Printf("  POST /auth/login - User login")
	log.Printf("  GET  /api/user/profile - Get user profile (protected)")
	log.Printf("  PUT  /api/user/profile - Update user profile (protected)")
	log.Printf("  POST /api/user/change-password - Change password (protected)")
	log.Printf("  DELETE /api/user/account - Delete account (protected)")

	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
