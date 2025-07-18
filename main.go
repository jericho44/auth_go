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
	fileRepo := repositories.NewFileRepository(database.DB)

	// Initialize services
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, authRepo)
	fileService := services.NewFileService(fileRepo, userRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)
	fileController := controllers.NewFileController(fileService)

	// Set controllers in handlers
	handlers.SetUserController(userController)
	handlers.SetAuthController(authController)
	handlers.SetFileController(fileController)
	handlers.SetFileService(fileService)

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
	log.Printf("  POST /api/files/upload - Upload single file (protected)")
	log.Printf("  POST /api/files/upload-multiple - Upload multiple files (protected)")
	log.Printf("  GET  /api/files/my - Get user's files (protected)")
	log.Printf("  GET  /api/files/{id} - Get file info")
	log.Printf("  PUT  /api/files/{id} - Update file (protected)")
	log.Printf("  DELETE /api/files/{id} - Delete file (protected)")
	log.Printf("  GET  /api/files/download/{filename} - Download file")
	log.Printf("  GET  /api/files/public - Get public files")
	log.Printf("  GET  /api/files/type/{type} - Get files by type")

	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
