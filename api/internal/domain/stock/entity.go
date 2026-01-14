package stock

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Errores del dominio
var (
	ErrStockNotFound = errors.New("stock not found")
)

// Stock representa una acción en el dominio
type Stock struct {
	ID          uuid.UUID
	Ticker      string
	CompanyName string
	Brokerage   string
	Action      string
	RatingFrom  Rating
	RatingTo    Rating
	TargetFrom  Price
	TargetTo    Price
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewStock crea una nueva entidad Stock con validaciones
func NewStock(
	ticker string,
	companyName string,
	brokerage string,
	action string,
	ratingFrom Rating,
	ratingTo Rating,
	targetFrom Price,
	targetTo Price,
) (*Stock, error) {
	// Validaciones
	if ticker == "" {
		return nil, fmt.Errorf("ticker cannot be empty")
	}
	if companyName == "" {
		return nil, fmt.Errorf("company name cannot be empty")
	}
	if !ratingFrom.IsValid() {
		return nil, fmt.Errorf("invalid rating_from: %s", ratingFrom)
	}
	if !ratingTo.IsValid() {
		return nil, fmt.Errorf("invalid rating_to: %s", ratingTo)
	}

	now := time.Now()
	return &Stock{
		ID:          uuid.New(),
		Ticker:      ticker,
		CompanyName: companyName,
		Brokerage:   brokerage,
		Action:      action,
		RatingFrom:  ratingFrom,
		RatingTo:    ratingTo,
		TargetFrom:  targetFrom,
		TargetTo:    targetTo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update actualiza los campos de la acción
func (s *Stock) Update(
	companyName string,
	brokerage string,
	action string,
	ratingFrom Rating,
	ratingTo Rating,
	targetFrom Price,
	targetTo Price,
) error {
	if companyName == "" {
		return fmt.Errorf("company name cannot be empty")
	}
	if !ratingFrom.IsValid() {
		return fmt.Errorf("invalid rating_from: %s", ratingFrom)
	}
	if !ratingTo.IsValid() {
		return fmt.Errorf("invalid rating_to: %s", ratingTo)
	}

	s.CompanyName = companyName
	s.Brokerage = brokerage
	s.Action = action
	s.RatingFrom = ratingFrom
	s.RatingTo = ratingTo
	s.TargetFrom = targetFrom
	s.TargetTo = targetTo
	s.UpdatedAt = time.Now()

	return nil
}

// CalculatePriceChange calcula el cambio porcentual del precio objetivo
func (s *Stock) CalculatePriceChange() float64 {
	if s.TargetFrom.IsZero() {
		return 0
	}
	change := s.TargetTo.Value() - s.TargetFrom.Value()
	percentChange := (change / s.TargetFrom.Value()) * 100
	return percentChange
}

// IsRatingUpgrade retorna true si el rating mejoró
func (s *Stock) IsRatingUpgrade() bool {
	fromScore := s.getRatingScore(s.RatingFrom)
	toScore := s.getRatingScore(s.RatingTo)
	return toScore > fromScore
}

// getRatingScore retorna un score numérico para el rating
func (s *Stock) getRatingScore(rating Rating) int {
	switch rating {
	case RatingStrongBuy:
		return 5
	case RatingBuy:
		return 3
	case RatingSpeculativeBuy:
		return 2
	case RatingMarketPerform:
		return 1
	case RatingNeutral:
		return 0
	case RatingUnderperform:
		return -1
	case RatingSell:
		return -2
	case RatingStrongSell:
		return -3
	default:
		return 0
	}
}
