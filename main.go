package main

import (
	"log"
	"net/http"

	"auth-jwt/config"
	"auth-jwt/controllers"
	"auth-jwt/database"
	_ "auth-jwt/docs"
	"auth-jwt/handlers"
	"auth-jwt/repositories"
	"auth-jwt/routes"
	"auth-jwt/services"
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

	// Get underlying sql.DB for defer close
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
	}
	defer sqlDB.Close()

	// Initialize JWT with config
	utils.InitJWT(cfg)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database.DB)
	authRepo := repositories.NewAuthRepository(database.DB)

	// Initialize services
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, authRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)

	// Set controllers in handlers
	handlers.SetUserController(userController)
	handlers.SetAuthController(authController)

	// Setup routes
	r := routes.SetupRoutes()

	log.Printf("Server starting on :%s", cfg.ServerPort)
	log.Printf("Swagger documentation available at: http://localhost:%s/swagger/", cfg.ServerPort)
	log.Printf("API endpoints:")
	log.Printf("  POST /auth/register - User registration")
	log.Printf("  POST /auth/login - User login")
	log.Printf("  POST /auth/forgot-password - Password reset request")
	log.Printf("  POST /auth/reset-password - Password reset confirmation")
	log.Printf("  GET  /api/user/profile - Get user profile (protected)")
	log.Printf("  PUT  /api/user/profile - Update user profile (protected)")
	log.Printf("  POST /api/user/change-password - Change password (protected)")
	log.Printf("  DELETE /api/user/account - Delete account (protected)")
	log.Printf("  GET  /api/user/stats - Get user statistics (protected)")

	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
