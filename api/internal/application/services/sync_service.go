package services

import (
	"context"
	"fmt"

	"github.com/john/go-react-test/api/internal/domain/stock"
	"github.com/john/go-react-test/api/internal/infrastructure/external"
)

// SyncService maneja la sincronización de stocks desde la API externa
type SyncService struct {
	apiClient *external.KarenAIClient
	repo      stock.Repository
}

// NewSyncService crea un nuevo servicio de sincronización
func NewSyncService(apiClient *external.KarenAIClient, repo stock.Repository) *SyncService {
	return &SyncService{
		apiClient: apiClient,
		repo:      repo,
	}
}

// SyncAllStocks sincroniza todos los stocks desde la API externa
func (s *SyncService) SyncAllStocks(ctx context.Context) (int, error) {
	// Obtener todos los stocks de la API externa
	stocks, err := s.apiClient.FetchAllStocks(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch stocks from API: %w", err)
	}

	// Guardar en base de datos usando batch upsert
	if err := s.repo.BatchUpsert(ctx, stocks); err != nil {
		return 0, fmt.Errorf("failed to save stocks to database: %w", err)
	}

	return len(stocks), nil
}
