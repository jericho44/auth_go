package database

import (
	"fmt"
	"log"

	"auth-jwt/config"
	"auth-jwt/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Test connection
	if err = sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to database successfully")

	// Auto-migrate models
	if err := AutoMigrate(); err != nil {
		log.Printf("Auto-migration warning: %v", err)
	}
}

// AutoMigrate runs GORM auto-migration for all models
func AutoMigrate() error {
	// log.Println("Running GORM auto-migration...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.PasswordResetToken{},
		&models.LoginAttempt{},
		&models.File{},
	)

	if err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	// log.Println("GORM auto-migration completed successfully")
	return nil
}

// GetDB returns the GORM database instance
func GetDB() *gorm.DB {
	return DB
}
