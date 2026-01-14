package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/john/go-react-test/api/internal/application/graphql"
	"github.com/john/go-react-test/api/internal/application/handlers"
	"github.com/john/go-react-test/api/internal/application/services"
	"github.com/john/go-react-test/api/internal/config"
	"github.com/john/go-react-test/api/internal/domain/recommendation"
	"github.com/john/go-react-test/api/internal/domain/stock"
	"github.com/john/go-react-test/api/internal/infrastructure/database"
	"github.com/john/go-react-test/api/internal/infrastructure/external"
	"github.com/john/go-react-test/api/internal/infrastructure/repository"
)

func main() {
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

	log.Println("Database connected successfully")

	// Verificar si las migraciones ya están aplicadas
	exists, err := database.CheckMigrations()
	if err != nil {
		log.Fatalf("Failed to check migrations: %v", err)
	}

	if !exists {
		log.Println("Database not initialized, running migrations...")
		if err := database.RunMigrations(); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Migrations executed successfully")
	} else {
		log.Println("Database already initialized")
	}

	// Inicializar dependencias
	stockRepo := repository.NewCockroachStockRepository()
	stockDomainSvc := stock.NewDomainService()
	stockService := services.NewStockService(stockRepo, stockDomainSvc)

	apiClient := external.NewKarenAIClient(cfg.API.BaseURL, cfg.API.APIKey)
	syncService := services.NewSyncService(apiClient, stockRepo)

	recommendationAlgorithm := recommendation.NewRecommendationAlgorithm(stockDomainSvc)
	recommendationService := services.NewRecommendationService(stockService, recommendationAlgorithm)

	// Inicializar GraphQL schema
	graphqlSchema, err := graphql.NewSchema(stockService, syncService, recommendationService)
	if err != nil {
		log.Fatalf("Failed to create GraphQL schema: %v", err)
	}

	// Crear handler GraphQL
	graphqlHandler := handlers.NewGraphQLHandler(graphqlSchema.GetSchema())

	// Configurar servidor HTTP
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// GraphQL endpoint
	mux.Handle("/query", graphqlHandler)

	// GraphQL Playground (solo en desarrollo)
	mux.Handle("/playground", handlers.PlaygroundHandler("GraphQL Playground", "/query"))

	// Documentación
	mux.HandleFunc("/docs", handlers.DocsHandler())
	mux.HandleFunc("/docs/swagger", handlers.SwaggerUIHandler())
	mux.HandleFunc("/docs/openapi.yaml", handlers.OpenAPISpecHandler())
	
	// Documentación Markdown
	mux.HandleFunc("/docs/api", handlers.MarkdownDocHandler("API_DOCUMENTATION.md"))
	mux.HandleFunc("/docs/guide", handlers.MarkdownDocHandler("USER_GUIDE.md"))
	mux.HandleFunc("/docs/graphql", func(w http.ResponseWriter, r *http.Request) {
		// Redirigir a la documentación GraphQL en docs/
		http.Redirect(w, r, "/docs/graphql-reference", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/docs/graphql-reference", handlers.MarkdownDocHandler("GRAPHQL_API_REFERENCE.md"))

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second, // Aumentado para operaciones largas como syncStocks
		IdleTimeout:  120 * time.Second,
	}

	// Iniciar servidor en goroutine
	go func() {
		log.Printf("Starting server on :%s\n", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Esperar señal de interrupción
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
