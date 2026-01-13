package recommendation

import (
	"github.com/john/go-react-test/api/internal/domain/stock"
)

// Recommendation representa una recomendación de inversión
type Recommendation struct {
	Stock       *stock.Stock
	Score       float64
	PriceChange float64
	RatingScore float64
	ActionScore float64
}

// NewRecommendation crea una nueva recomendación
func NewRecommendation(
	s *stock.Stock,
	score float64,
	priceChange float64,
	ratingScore float64,
	actionScore float64,
) *Recommendation {
	return &Recommendation{
		Stock:       s,
		Score:       score,
		PriceChange: priceChange,
		RatingScore: ratingScore,
		ActionScore: actionScore,
	}
}
