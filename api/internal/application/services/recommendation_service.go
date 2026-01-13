package services

import (
	"context"

	"github.com/john/go-react-test/api/internal/domain/recommendation"
	"github.com/john/go-react-test/api/internal/domain/stock"
)

// RecommendationService es el servicio de aplicación para recomendaciones
type RecommendationService struct {
	stockService *StockService
	algorithm    recommendation.Algorithm
}

// NewRecommendationService crea un nuevo servicio de recomendaciones
func NewRecommendationService(
	stockService *StockService,
	algorithm recommendation.Algorithm,
) *RecommendationService {
	return &RecommendationService{
		stockService: stockService,
		algorithm:    algorithm,
	}
}

// GetRecommendations obtiene las mejores recomendaciones de inversión
func (s *RecommendationService) GetRecommendations(ctx context.Context, limit int) ([]*recommendation.Recommendation, error) {
	// Obtener todos los stocks con rating positivo
	filter := stock.Filter{
		Ratings: []stock.Rating{
			stock.RatingStrongBuy,
			stock.RatingBuy,
			stock.RatingSpeculativeBuy,
		},
	}

	stocks, err := s.stockService.GetStocks(ctx, filter, stock.Sort{})
	if err != nil {
		return nil, err
	}

	// Calcular recomendaciones usando el algoritmo
	return s.algorithm.CalculateRecommendations(ctx, stocks, limit)
}
