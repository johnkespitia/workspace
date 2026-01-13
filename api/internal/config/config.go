package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config contiene la configuración de la aplicación
type Config struct {
	Database DatabaseConfig
	API      APIConfig
	Server   ServerConfig
}

// DatabaseConfig configuración de base de datos
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// APIConfig configuración de API externa
type APIConfig struct {
	BaseURL string
	APIKey  string
}

// ServerConfig configuración del servidor
type ServerConfig struct {
	Port string
}

// Load carga la configuración desde variables de entorno
func Load() (*Config, error) {
	// Intentar cargar archivos .env si existen
	// Buscar desde el directorio actual hacia arriba
	loadEnvFiles()

	cfg := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "26257"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "stocks"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		API: APIConfig{
			BaseURL: getEnv("API_BASE_URL", "https://api.karenai.click"),
			APIKey:  getEnv("API_KEY", ""),
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
	}

	if cfg.API.APIKey == "" {
		return nil, fmt.Errorf("API_KEY environment variable is required")
	}

	return cfg, nil
}

// DatabaseDSN retorna el Data Source Name para la conexión a la base de datos
func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// loadEnvFiles intenta cargar archivos .env desde el directorio del proyecto
func loadEnvFiles() {
	// Buscar el directorio api/ desde el directorio actual
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	// Buscar archivos .env en orden de prioridad
	envFiles := []string{
		filepath.Join(wd, ".env.development"),
		filepath.Join(wd, ".env.local"),
		filepath.Join(wd, ".env"),
	}

	// También buscar en el directorio padre (si estamos en cmd/)
	parentDir := filepath.Dir(wd)
	envFiles = append(envFiles,
		filepath.Join(parentDir, ".env.development"),
		filepath.Join(parentDir, ".env.local"),
		filepath.Join(parentDir, ".env"),
	)

	// Cargar el primer archivo que exista
	for _, envFile := range envFiles {
		if _, err := os.Stat(envFile); err == nil {
			_ = godotenv.Load(envFile)
			return
		}
	}
}
