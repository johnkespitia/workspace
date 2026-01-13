package stock

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Rating representa el rating de una acción
type Rating string

const (
	RatingStrongBuy      Rating = "Strong Buy"
	RatingBuy            Rating = "Buy"
	RatingSpeculativeBuy Rating = "Speculative Buy"
	RatingMarketPerform  Rating = "Market Perform"
	RatingNeutral        Rating = "Neutral"
	RatingUnderperform   Rating = "Underperform"
	RatingSell           Rating = "Sell"
	RatingStrongSell     Rating = "Strong Sell"
)

// IsValid valida si el rating es válido
func (r Rating) IsValid() bool {
	switch r {
	case RatingStrongBuy, RatingBuy, RatingSpeculativeBuy,
		RatingMarketPerform, RatingNeutral, RatingUnderperform,
		RatingSell, RatingStrongSell:
		return true
	default:
		return false
	}
}

// IsPositive retorna true si el rating es positivo para inversión
func (r Rating) IsPositive() bool {
	switch r {
	case RatingStrongBuy, RatingBuy, RatingSpeculativeBuy:
		return true
	case RatingMarketPerform:
		// Market Perform puede ser positivo si hay aumento de target
		return true
	default:
		return false
	}
}

// String retorna el string del rating
func (r Rating) String() string {
	return string(r)
}

// Price representa un precio monetario
type Price struct {
	value decimal.Decimal
}

// NewPrice crea un nuevo Price
func NewPrice(value float64) (Price, error) {
	if value < 0 {
		return Price{}, fmt.Errorf("price cannot be negative")
	}
	return Price{value: decimal.NewFromFloat(value)}, nil
}

// NewPriceFromDecimal crea un Price desde decimal.Decimal
func NewPriceFromDecimal(d decimal.Decimal) Price {
	return Price{value: d}
}

// Value retorna el valor numérico del precio
func (p Price) Value() float64 {
	f, _ := p.value.Float64()
	return f
}

// Decimal retorna el valor como decimal.Decimal
func (p Price) Decimal() decimal.Decimal {
	return p.value
}

// String retorna el string del precio
func (p Price) String() string {
	return p.value.String()
}

// IsZero retorna true si el precio es cero
func (p Price) IsZero() bool {
	return p.value.IsZero()
}
