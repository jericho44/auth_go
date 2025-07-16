package database

import (
	"database/sql"
	"fmt"
	"log"

	"auth-jwt/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to database successfully")

	// Run migrations instead of creating tables manually
	if err := RunMigrations(cfg); err != nil {
		log.Printf("Migration warning: %v", err)
	}
}
