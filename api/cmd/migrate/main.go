package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/john/go-react-test/api/internal/config"
	"github.com/john/go-react-test/api/internal/infrastructure/database"
)

func main() {
	var (
		check = flag.Bool("check", false, "Check if migrations have been applied")
		reset = flag.Bool("reset", false, "Reset database (drop all tables and reapply migrations)")
		up    = flag.Bool("up", false, "Run all pending migrations")
		help  = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Conectar a la base de datos
	if err := database.Connect(cfg.DatabaseDSN()); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Ejecutar comando según flag
	if *check {
		exists, err := database.CheckMigrations()
		if err != nil {
			log.Fatalf("Error checking migrations: %v", err)
		}

		if exists {
			fmt.Println("✓ Database is initialized: stocks table exists")
			os.Exit(0)
		} else {
			fmt.Println("✗ Database is not initialized: stocks table does not exist")
			os.Exit(1)
		}
	}

	if *reset {
		fmt.Println("⚠️  WARNING: This will delete all data in the database!")
		fmt.Print("Are you sure you want to continue? (yes/no): ")
		var confirmation string
		fmt.Scanln(&confirmation)

		if confirmation != "yes" {
			fmt.Println("Operation cancelled")
			return
		}

		if err := database.ResetDatabase(); err != nil {
			log.Fatalf("Failed to reset database: %v", err)
		}
		fmt.Println("✓ Database reset completed successfully")
		return
	}

	if *up {
		if err := database.RunMigrations(); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("✓ All migrations completed successfully")
		return
	}

	// Si no se especificó ningún comando, mostrar ayuda
	showHelp()
}

func showHelp() {
	fmt.Println("Database Migration Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd/migrate [command]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  -check    Check if migrations have been applied")
	fmt.Println("  -up       Run all pending migrations")
	fmt.Println("  -reset    Reset database (drop all tables and reapply migrations)")
	fmt.Println("  -help     Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run ./cmd/migrate -check")
	fmt.Println("  go run ./cmd/migrate -up")
	fmt.Println("  go run ./cmd/migrate -reset")
}
