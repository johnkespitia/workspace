package stock

// Service define los servicios de dominio para stocks
type Service interface {
	// CalculatePriceChange calcula el cambio porcentual en el precio objetivo
	CalculatePriceChange(stock *Stock) float64

	// IsRatingUpgrade determina si el rating mejoró
	IsRatingUpgrade(stock *Stock) bool

	// CalculateRecommendationScore calcula un score de recomendación
	CalculateRecommendationScore(stock *Stock) float64
}

// DomainService implementa los servicios de dominio
type DomainService struct{}

// NewDomainService crea un nuevo servicio de dominio
func NewDomainService() Service {
	return &DomainService{}
}

// CalculatePriceChange calcula el cambio porcentual en el precio objetivo
func (s *DomainService) CalculatePriceChange(stock *Stock) float64 {
	return stock.CalculatePriceChange()
}

// IsRatingUpgrade determina si el rating mejoró
func (s *DomainService) IsRatingUpgrade(stock *Stock) bool {
	return stock.IsRatingUpgrade()
}

// CalculateRecommendationScore calcula un score de recomendación
// Score = (priceChange * 0.5) + (ratingScore * 0.3) + (actionScore * 0.2)
func (s *DomainService) CalculateRecommendationScore(stock *Stock) float64 {
	// Calcular cambio porcentual
	priceChange := stock.CalculatePriceChange()

	// Calcular score de rating
	ratingScore := s.getRatingScore(stock.RatingTo)
	if stock.IsRatingUpgrade() {
		ratingScore += 2 // Bonus por upgrade
	}

	// Calcular score de acción
	actionScore := s.getActionScore(stock.Action)

	// Score final con pesos
	finalScore := (priceChange * 0.5) + (float64(ratingScore) * 0.3) + (float64(actionScore) * 0.2)

	return finalScore
}

// getRatingScore retorna un score numérico para el rating
func (s *DomainService) getRatingScore(rating Rating) float64 {
	switch rating {
	case RatingStrongBuy:
		return 5
	case RatingBuy:
		return 3
	case RatingSpeculativeBuy:
		return 2
	case RatingMarketPerform:
		return 1
	default:
		return 0
	}
}

// getActionScore retorna un score numérico para la acción
func (s *DomainService) getActionScore(action string) float64 {
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
