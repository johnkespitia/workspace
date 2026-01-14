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
		// Dividir el SQL en statements individuales y ejecutarlos uno por uno
		// Esto es necesario porque algunos drivers/SGBD requieren statements separados
		statements := splitSQLStatements(migration.SQL)
		
		for i, statement := range statements {
			statement = strings.TrimSpace(statement)
			if statement == "" {
				continue // Saltar líneas vacías
			}
			
			// Saltar comentarios de línea completa
			if strings.HasPrefix(statement, "--") {
				continue
			}
			
			// Ejecutar el statement
			// Si es un trigger o función de trigger y falla, solo mostrar un warning (CockroachDB puede no soportarlos)
			_, err := db.Exec(statement)
			if err != nil {
				upperStmt := strings.ToUpper(statement)
				// Detectar si es un CREATE TRIGGER o CREATE FUNCTION relacionado con triggers
				isTrigger := strings.Contains(upperStmt, "CREATE TRIGGER")
				isTriggerFunction := strings.Contains(upperStmt, "CREATE") && 
				                     strings.Contains(upperStmt, "FUNCTION") && 
				                     (strings.Contains(upperStmt, "UPDATE_UPDATED_AT") || 
				                      strings.Contains(upperStmt, "TRIGGER") ||
				                      strings.Contains(upperStmt, "RETURNS TRIGGER"))
				
				if isTrigger || isTriggerFunction {
					fmt.Printf("⚠️  Warning: Trigger/Function creation failed (may not be supported in this CockroachDB version): %v\n", err)
					fmt.Printf("   You may need to update 'updated_at' manually in your code.\n")
					fmt.Printf("   Continuing without trigger...\n")
					continue
				}
				return fmt.Errorf("failed to execute migration %s (statement %d): %w\nStatement: %s", 
					migration.Filename, i+1, err, statement)
			}
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
			SELECT 1 FROM information_schema.tables 
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

// splitSQLStatements divide un string SQL en statements individuales
// separados por punto y coma
func splitSQLStatements(sql string) []string {
	var statements []string
	var current strings.Builder
	inString := false
	stringChar := byte(0)
	inDollarQuote := false
	dollarTag := ""
	
	lines := strings.Split(sql, "\n")
	
	for _, line := range lines {
		// Saltar comentarios de línea completa
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--") {
			continue
		}
		
		for i := 0; i < len(line); i++ {
			char := line[i]
			
			// Manejar comillas simples y dobles
			if !inDollarQuote {
				if char == '\'' || char == '"' {
					if !inString {
						inString = true
						stringChar = char
					} else if char == stringChar {
						// Verificar que no sea un escape
						if i == 0 || line[i-1] != '\\' {
							inString = false
						}
					}
				}
			}
			
			// Manejar dollar quoting ($$ ... $$)
			if char == '$' && !inString {
				if i+1 < len(line) {
					// Buscar el tag del dollar quote
					j := i + 1
					for j < len(line) && line[j] != '$' {
						j++
					}
					if j < len(line) {
						tag := line[i+1 : j]
						if !inDollarQuote {
							dollarTag = tag
							inDollarQuote = true
						} else if tag == dollarTag {
							inDollarQuote = false
							dollarTag = ""
						}
					}
				}
			}
			
			current.WriteByte(char)
			
			// Si encontramos un punto y coma fuera de strings, es el final de un statement
			if char == ';' && !inString && !inDollarQuote {
				stmt := strings.TrimSpace(current.String())
				if stmt != "" {
					statements = append(statements, stmt)
				}
				current.Reset()
			}
		}
		
		// Agregar salto de línea si no estamos al final de un statement
		if current.Len() > 0 {
			current.WriteByte('\n')
		}
	}
	
	// Agregar el último statement si no termina con punto y coma
	stmt := strings.TrimSpace(current.String())
	if stmt != "" {
		statements = append(statements, stmt)
	}
	
	return statements
}
