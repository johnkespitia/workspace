package recommendation

import (
	"context"
	"sort"

	"github.com/john/go-react-test/api/internal/domain/stock"
)

// Algorithm define el algoritmo de recomendación
type Algorithm interface {
	// CalculateRecommendations calcula las recomendaciones basadas en stocks
	CalculateRecommendations(ctx context.Context, stocks []*stock.Stock, limit int) ([]*Recommendation, error)
}

// RecommendationAlgorithm implementa el algoritmo de recomendación
type RecommendationAlgorithm struct {
	stockService stock.Service
}

// NewRecommendationAlgorithm crea un nuevo algoritmo de recomendación
func NewRecommendationAlgorithm(stockService stock.Service) Algorithm {
	return &RecommendationAlgorithm{
		stockService: stockService,
	}
}

// CalculateRecommendations calcula las recomendaciones
// Complejidad: O(n log n) donde n = número de stocks
func (a *RecommendationAlgorithm) CalculateRecommendations(
	ctx context.Context,
	stocks []*stock.Stock,
	limit int,
) ([]*Recommendation, error) {
	// Paso 1: Filtrar stocks con rating positivo
	filteredStocks := a.filterPositiveRatings(stocks)

	// Paso 2: Calcular scores para cada stock (O(n))
	recommendations := make([]*Recommendation, 0, len(filteredStocks))
	for _, s := range filteredStocks {
		score := a.calculateScore(s)
		recommendations = append(recommendations, score)
	}

	// Paso 3: Ordenar por score descendente (O(n log n))
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Paso 4: Retornar top N
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}

// filterPositiveRatings filtra stocks con rating positivo
func (a *RecommendationAlgorithm) filterPositiveRatings(stocks []*stock.Stock) []*stock.Stock {
	filtered := make([]*stock.Stock, 0)
	for _, s := range stocks {
		if s.RatingTo.IsPositive() {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// calculateScore calcula el score de recomendación para un stock
func (a *RecommendationAlgorithm) calculateScore(s *stock.Stock) *Recommendation {
	// Calcular cambio porcentual
	priceChange := s.CalculatePriceChange()

	// Calcular score de rating
	ratingScore := a.getRatingScore(s.RatingTo)
	if s.IsRatingUpgrade() {
		ratingScore += 2 // Bonus por upgrade
	}

	// Calcular score de acción
	actionScore := a.getActionScore(s.Action)

	// Score final con pesos: 50% price change, 30% rating, 20% action
	finalScore := (priceChange * 0.5) + (ratingScore * 0.3) + (actionScore * 0.2)

	return NewRecommendation(s, finalScore, priceChange, ratingScore, actionScore)
}

// getRatingScore retorna un score numérico para el rating
func (a *RecommendationAlgorithm) getRatingScore(rating stock.Rating) float64 {
	switch rating {
	case stock.RatingStrongBuy:
		return 5
	case stock.RatingBuy:
		return 3
	case stock.RatingSpeculativeBuy:
		return 2
	case stock.RatingMarketPerform:
		return 1
	default:
		return 0
	}
}

// getActionScore retorna un score numérico para la acción
func (a *RecommendationAlgorithm) getActionScore(action string) float64 {
	switch action {
	case "target raised by", "target raised":
		return 3
	case "target lowered by", "target lowered":
		return -2
	case "initiated coverage":
		return 1
	default:
		return 0
	}
}
