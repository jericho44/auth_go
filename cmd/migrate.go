package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"auth-jwt/config"
	"auth-jwt/database"
)

func main() {
	var action = flag.String("action", "up", "Migration action: up, down, version")
	flag.Parse()

	cfg := config.Load()

	switch *action {
	case "up":
		if err := database.RunMigrations(cfg); err != nil {
			log.Fatal("Migration failed:", err)
		}
		fmt.Println("Migrations completed successfully")

	case "down":
		if err := database.RollbackMigrations(cfg); err != nil {
			log.Fatal("Rollback failed:", err)
		}
		fmt.Println("Rollback completed successfully")

	case "version":
		version, dirty, err := database.GetMigrationVersion(cfg)
		if err != nil {
			log.Fatal("Failed to get version:", err)
		}
		fmt.Printf("Current migration version: %d\n", version)
		if dirty {
			fmt.Println("Warning: Database is in dirty state")
		}

	default:
		fmt.Println("Usage: go run cmd/migrate.go -action=[up|down|version]")
		os.Exit(1)
	}
}
