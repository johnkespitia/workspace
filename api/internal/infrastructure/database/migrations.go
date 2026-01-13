package database

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// MigrationInfo contiene información sobre una migración
type MigrationInfo struct {
	Filename string
	SQL      string
}

// RunMigrations ejecuta todas las migraciones en orden
func RunMigrations() error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration.SQL); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migration.Filename, err)
		}
		fmt.Printf("✓ Migration %s executed successfully\n", migration.Filename)
	}

	return nil
}

// CheckMigrations verifica si las tablas existen en la base de datos
func CheckMigrations() (bool, error) {
	db := GetDB()
	if db == nil {
		return false, fmt.Errorf("database connection is not initialized")
	}

	// Verificar si la tabla stocks existe
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'stocks'
		)
	`

	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check migrations: %w", err)
	}

	return exists, nil
}

// ResetDatabase elimina todas las tablas y las recrea
func ResetDatabase() error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	// Eliminar tablas en orden inverso (por si hay dependencias)
	dropQueries := []string{
		"DROP TRIGGER IF EXISTS update_stocks_updated_at ON stocks",
		"DROP FUNCTION IF EXISTS update_updated_at_column()",
		"DROP TABLE IF EXISTS stocks CASCADE",
	}

	for _, query := range dropQueries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to drop: %w", err)
		}
	}

	fmt.Println("✓ Database reset: all tables dropped")

	// Ejecutar migraciones nuevamente
	return RunMigrations()
}

// loadMigrations carga todos los archivos SQL de la carpeta migrations
func loadMigrations() ([]MigrationInfo, error) {
	// Buscar el directorio de migraciones relativo a este archivo
	// migrations.go está en internal/infrastructure/database/
	// Las migraciones están en internal/infrastructure/database/migrations/

	// Obtener el directorio donde está este archivo
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	migrationsPath := filepath.Join(baseDir, "migrations")

	var migrations []MigrationInfo

	err := filepath.WalkDir(migrationsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Solo procesar archivos .sql
		if !d.IsDir() && strings.HasSuffix(path, ".sql") {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", path, err)
			}

			filename := filepath.Base(path)
			migrations = append(migrations, MigrationInfo{
				Filename: filename,
				SQL:      string(content),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Ordenar migraciones por nombre (para ejecutarlas en orden)
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Filename < migrations[j].Filename
	})

	return migrations, nil
}
