package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Driver PostgreSQL para CockroachDB
)

// DB es la conexión global a la base de datos
var DB *sql.DB

// Connect establece la conexión a CockroachDB
func Connect(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verificar la conexión
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	return nil
}

// Close cierra la conexión a la base de datos
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB retorna la conexión a la base de datos
func GetDB() *sql.DB {
	return DB
}
