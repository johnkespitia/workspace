package stock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainService_CalculatePriceChange(t *testing.T) {
	service := NewDomainService()

	tests := []struct {
		name           string
		targetFrom     float64
		targetTo       float64
		expectedChange float64
	}{
		{
			name:           "price increase",
			targetFrom:     100.0,
			targetTo:       120.0,
			expectedChange: 20.0, // 20% increase
		},
		{
			name:           "price decrease",
			targetFrom:     100.0,
			targetTo:       80.0,
			expectedChange: -20.0, // 20% decrease
		},
		{
			name:           "no change",
			targetFrom:     100.0,
			targetTo:       100.0,
			expectedChange: 0.0,
		},
		{
			name:           "zero target from",
			targetFrom:     0.0,
			targetTo:       100.0,
			expectedChange: 0.0, // Should return 0 if targetFrom is zero
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetFrom, _ := NewPrice(tt.targetFrom)
			targetTo, _ := NewPrice(tt.targetTo)
			
			stock := &Stock{
				TargetFrom: targetFrom,
				TargetTo:   targetTo,
			}

			result := service.CalculatePriceChange(stock)
			assert.InDelta(t, tt.expectedChange, result, 0.01, "price change should match expected")
		})
	}
}

func TestDomainService_IsRatingUpgrade(t *testing.T) {
	service := NewDomainService()

	tests := []struct {
		name         string
		ratingFrom   Rating
		ratingTo     Rating
		expectedUpgrade bool
	}{
		{
			name:         "upgrade from Neutral to Buy",
			ratingFrom:   RatingNeutral,
			ratingTo:     RatingBuy,
			expectedUpgrade: true,
		},
		{
			name:         "upgrade from Buy to Strong Buy",
			ratingFrom:   RatingBuy,
			ratingTo:     RatingStrongBuy,
			expectedUpgrade: true,
		},
		{
			name:         "downgrade from Buy to Neutral",
			ratingFrom:   RatingBuy,
			ratingTo:     RatingNeutral,
			expectedUpgrade: false,
		},
		{
			name:         "no change",
			ratingFrom:   RatingBuy,
			ratingTo:     RatingBuy,
			expectedUpgrade: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stock := &Stock{
				RatingFrom: tt.ratingFrom,
				RatingTo:   tt.ratingTo,
			}

			result := service.IsRatingUpgrade(stock)
			assert.Equal(t, tt.expectedUpgrade, result, "rating upgrade should match expected")
		})
	}
}

func TestDomainService_CalculateRecommendationScore(t *testing.T) {
	service := NewDomainService()

	tests := []struct {
		name           string
		targetFrom     float64
		targetTo       float64
		ratingFrom     Rating
		ratingTo       Rating
		action         string
		expectedMin    float64 // Minimum expected score
		expectedMax    float64 // Maximum expected score
	}{
		{
			name:        "high score: price increase + strong buy + upgrade + target raised",
			targetFrom:  100.0,
			targetTo:   150.0, // 50% increase
			ratingFrom: RatingNeutral,
			ratingTo:   RatingStrongBuy, // 5 + 2 (upgrade bonus) = 7
			action:     "target raised by", // 3
			expectedMin: 25.0, // (50 * 0.5) + (7 * 0.3) + (3 * 0.2) = 25 + 2.1 + 0.6 = 27.7
			expectedMax: 30.0,
		},
		{
			name:        "low score: price decrease + neutral rating + target lowered",
			targetFrom:  100.0,
			targetTo:   80.0, // -20% decrease
			ratingFrom: RatingBuy,
			ratingTo:   RatingNeutral, // 0, no upgrade
			action:     "target lowered by", // -2
			expectedMin: -15.0,
			expectedMax: -5.0,
		},
		{
			name:        "medium score: small increase + buy rating",
			targetFrom:  100.0,
			targetTo:   110.0, // 10% increase
			ratingFrom: RatingBuy,
			ratingTo:   RatingBuy, // 3, no upgrade
			action:     "initiated coverage", // 1
			expectedMin: 5.0,
			expectedMax: 10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetFrom, _ := NewPrice(tt.targetFrom)
			targetTo, _ := NewPrice(tt.targetTo)
			
			stock := &Stock{
				TargetFrom:  targetFrom,
				TargetTo:    targetTo,
				RatingFrom:  tt.ratingFrom,
				RatingTo:    tt.ratingTo,
				Action:      tt.action,
			}

			result := service.CalculateRecommendationScore(stock)
			assert.GreaterOrEqual(t, result, tt.expectedMin, "score should be at least minimum")
			assert.LessOrEqual(t, result, tt.expectedMax, "score should be at most maximum")
		})
	}
}

// TestDomainService_getRatingScore tests rating scores indirectly through CalculateRecommendationScore
func TestDomainService_getRatingScore(t *testing.T) {
	service := NewDomainService()

	tests := []struct {
		name     string
		rating   Rating
		expected float64
	}{
		{"Strong Buy", RatingStrongBuy, 5.0},
		{"Buy", RatingBuy, 3.0},
		{"Speculative Buy", RatingSpeculativeBuy, 2.0},
		{"Market Perform", RatingMarketPerform, 1.0},
		{"Neutral", RatingNeutral, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test indirectly through CalculateRecommendationScore
			// Create a stock with no price change and no action to isolate rating score
			targetFrom, _ := NewPrice(100.0)
			targetTo, _ := NewPrice(100.0) // No change
			
			stock := &Stock{
				TargetFrom: targetFrom,
				TargetTo:   targetTo,
				RatingFrom: tt.rating,
				RatingTo:   tt.rating, // No upgrade
				Action:     "", // No action
			}

			score := service.CalculateRecommendationScore(stock)
			// Score should be approximately (ratingScore * 0.3) since priceChange=0 and action=0
			expectedScore := tt.expected * 0.3
			assert.InDelta(t, expectedScore, score, 0.1, "rating score contribution should match")
		})
	}
}

// TestDomainService_getActionScore tests action scores indirectly through CalculateRecommendationScore
func TestDomainService_getActionScore(t *testing.T) {
	service := NewDomainService()

	tests := []struct {
		name     string
		action   string
		expected float64
	}{
		{"target raised by", "target raised by", 3.0},
		{"target raised", "target raised", 3.0},
		{"target lowered by", "target lowered by", -2.0},
		{"target lowered", "target lowered", -2.0},
		{"initiated coverage", "initiated coverage", 1.0},
		{"unknown action", "unknown", 0.0},
		{"empty action", "", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test indirectly through CalculateRecommendationScore
			// Create a stock with no price change and neutral rating to isolate action score
			targetFrom, _ := NewPrice(100.0)
			targetTo, _ := NewPrice(100.0) // No change
			
			stock := &Stock{
				TargetFrom: targetFrom,
				TargetTo:   targetTo,
				RatingFrom: RatingNeutral,
				RatingTo:   RatingNeutral, // No upgrade, neutral rating = 0
				Action:     tt.action,
			}

			score := service.CalculateRecommendationScore(stock)
			// Score should be approximately (actionScore * 0.2) since priceChange=0 and ratingScore=0
			expectedScore := tt.expected * 0.2
			assert.InDelta(t, expectedScore, score, 0.1, "action score contribution should match")
		})
	}
}
