package main

import (
	"flag"
	"fmt"
	"kepler-auth-go/internal/config"
	"kepler-auth-go/internal/database"
	"log"
	"os"
)

func main() {
	var (
		action = flag.String("action", "up", "Migration action: up, down, status, fresh")
		force  = flag.Bool("force", false, "Force migration even if already applied")
		help   = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	cfg := config.Load()

	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	switch *action {
	case "up":
		fmt.Println("Running migrations...")
		if err := database.RunMigrations(*force); err != nil {
			log.Fatal("Migration failed:", err)
		}
		if err := database.SeedDefaultData(); err != nil {
			log.Printf("Warning: Failed to seed default data: %v", err)
		}
		fmt.Println("Migrations completed successfully!")

	case "down":
		fmt.Println("Rolling back last migration...")
		if err := database.RollbackMigration(); err != nil {
			log.Fatal("Rollback failed:", err)
		}
		fmt.Println("Rollback completed successfully!")

	case "status":
		fmt.Println("Migration status:")
		status, err := database.GetMigrationStatus()
		if err != nil {
			log.Fatal("Failed to get migration status:", err)
		}
		for _, migration := range status {
			applied := "❌ Not Applied"
			if migration["applied"].(bool) {
				applied = "✅ Applied"
			}
			fmt.Printf("  %s: %s\n", migration["id"], applied)
		}

	case "fresh":
		fmt.Println("Running fresh migrations (WARNING: This will drop all tables)...")
		fmt.Print("Are you sure? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Aborted.")
			return
		}

		// Drop all tables
		if err := database.GetDB().Migrator().DropTable(
			&database.MigrationRecord{},
			"auth_group_permissions",
			"user_groups",
		); err != nil {
			log.Printf("Warning: Failed to drop some tables: %v", err)
		}

		// Run fresh migrations
		if err := database.RunMigrations(true); err != nil {
			log.Fatal("Fresh migration failed:", err)
		}
		if err := database.SeedDefaultData(); err != nil {
			log.Printf("Warning: Failed to seed default data: %v", err)
		}
		fmt.Println("Fresh migrations completed successfully!")

	default:
		fmt.Printf("Unknown action: %s\n", *action)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Migration Tool for Kepler Auth Go")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/migrate/main.go [options]")
	fmt.Println("  OR use make commands:")
	fmt.Println("  make migrate, make migrate-status, make migrate-rollback")
	fmt.Println("")
	fmt.Println("Actions:")
	fmt.Println("  -action=up      Run pending migrations (default)")
	fmt.Println("  -action=down    Rollback last migration")
	fmt.Println("  -action=status  Show migration status")
	fmt.Println("  -action=fresh   Drop all tables and run fresh migrations")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -force          Force migration even if already applied")
	fmt.Println("  -help           Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/migrate/main.go                    # Run pending migrations")
	fmt.Println("  go run cmd/migrate/main.go -action=status     # Check migration status")
	fmt.Println("  go run cmd/migrate/main.go -action=down       # Rollback last migration")
	fmt.Println("  go run cmd/migrate/main.go -action=fresh      # Fresh install")
	fmt.Println("")
	fmt.Println("  make migrate          # Same as -action=up")
	fmt.Println("  make migrate-status   # Same as -action=status")
	fmt.Println("  make migrate-rollback # Same as -action=down")
}
