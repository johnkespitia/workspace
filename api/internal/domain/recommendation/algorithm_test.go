package recommendation

import (
	"context"
	"testing"

	"github.com/john/go-react-test/api/internal/domain/stock"
	"github.com/stretchr/testify/assert"
)

func TestRecommendationAlgorithm_CalculateRecommendations(t *testing.T) {
	stockService := stock.NewDomainService()
	algorithm := NewRecommendationAlgorithm(stockService)

	tests := []struct {
		name           string
		stocks         []*stock.Stock
		limit          int
		expectedCount  int
		expectedTopTicker string
	}{
		{
			name: "filter positive ratings only",
			stocks: []*stock.Stock{
				createTestStock("AAPL", 100.0, 120.0, stock.RatingBuy, stock.RatingStrongBuy, "target raised by"),
				createTestStock("MSFT", 50.0, 60.0, stock.RatingNeutral, stock.RatingBuy, "target raised"),
				createTestStock("GOOGL", 200.0, 180.0, stock.RatingBuy, stock.RatingNeutral, "target lowered"), // Should be filtered out
			},
			limit:          10,
			expectedCount:  2, // Only AAPL and MSFT have positive ratings
			expectedTopTicker: "AAPL", // Should have higher score
		},
		{
			name: "respect limit",
			stocks: []*stock.Stock{
				createTestStock("AAPL", 100.0, 150.0, stock.RatingBuy, stock.RatingStrongBuy, "target raised by"),
				createTestStock("MSFT", 50.0, 60.0, stock.RatingBuy, stock.RatingBuy, "target raised"),
				createTestStock("GOOGL", 200.0, 220.0, stock.RatingBuy, stock.RatingBuy, "initiated coverage"),
			},
			limit:         2,
			expectedCount: 2,
		},
		{
			name: "empty stocks",
			stocks: []*stock.Stock{},
			limit: 10,
			expectedCount: 0,
		},
		{
			name: "all negative ratings filtered out",
			stocks: []*stock.Stock{
				createTestStock("BAD1", 100.0, 80.0, stock.RatingBuy, stock.RatingNeutral, "target lowered"),
				createTestStock("BAD2", 50.0, 40.0, stock.RatingBuy, stock.RatingSell, "target lowered"),
			},
			limit: 10,
			expectedCount: 0, // All filtered out
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recommendations, err := algorithm.CalculateRecommendations(context.Background(), tt.stocks, tt.limit)
			
			assert.NoError(t, err)
			assert.Len(t, recommendations, tt.expectedCount)

			if tt.expectedCount > 0 && tt.expectedTopTicker != "" {
				// Verify that recommendations are sorted by score descending
				if len(recommendations) > 0 {
					assert.Equal(t, tt.expectedTopTicker, recommendations[0].Stock.Ticker)
					
					// Verify sorting
					for i := 1; i < len(recommendations); i++ {
						assert.GreaterOrEqual(t, recommendations[i-1].Score, recommendations[i].Score,
							"recommendations should be sorted by score descending")
					}
				}
			}
		})
	}
}

// TestRecommendationAlgorithm_filterPositiveRatings tests filtering indirectly through CalculateRecommendations
func TestRecommendationAlgorithm_filterPositiveRatings(t *testing.T) {
	stockService := stock.NewDomainService()
	algorithm := NewRecommendationAlgorithm(stockService)

	stocks := []*stock.Stock{
		createTestStock("AAPL", 100.0, 120.0, stock.RatingBuy, stock.RatingStrongBuy, ""),
		createTestStock("MSFT", 50.0, 60.0, stock.RatingNeutral, stock.RatingBuy, ""),
		createTestStock("GOOGL", 200.0, 180.0, stock.RatingBuy, stock.RatingNeutral, ""),
		createTestStock("TSLA", 300.0, 350.0, stock.RatingBuy, stock.RatingSpeculativeBuy, ""),
	}

	recommendations, err := algorithm.CalculateRecommendations(context.Background(), stocks, 10)
	assert.NoError(t, err)
	
	// Should only include stocks with positive ratings (Strong Buy, Buy, Speculative Buy, Market Perform)
	assert.Len(t, recommendations, 3) // AAPL, MSFT, TSLA (GOOGL has Neutral which is not positive)
	
	tickers := make(map[string]bool)
	for _, rec := range recommendations {
		tickers[rec.Stock.Ticker] = true
	}
	
	assert.True(t, tickers["AAPL"])
	assert.True(t, tickers["MSFT"])
	assert.True(t, tickers["TSLA"])
	assert.False(t, tickers["GOOGL"])
}

// TestRecommendationAlgorithm_calculateScore tests score calculation indirectly through CalculateRecommendations
func TestRecommendationAlgorithm_calculateScore(t *testing.T) {
	stockService := stock.NewDomainService()
	algorithm := NewRecommendationAlgorithm(stockService)

	tests := []struct {
		name           string
		targetFrom     float64
		targetTo       float64
		ratingFrom     stock.Rating
		ratingTo       stock.Rating
		action         string
		expectedMin    float64
		expectedMax    float64
	}{
		{
			name:        "high score scenario",
			targetFrom:  100.0,
			targetTo:   150.0, // 50% increase
			ratingFrom: stock.RatingNeutral,
			ratingTo:   stock.RatingStrongBuy, // 5 + 2 (upgrade) = 7
			action:     "target raised by", // 3
			expectedMin: 25.0,
			expectedMax: 30.0,
		},
		{
			name:        "low score scenario",
			targetFrom:  100.0,
			targetTo:   80.0, // -20% decrease
			ratingFrom: stock.RatingBuy,
			ratingTo:   stock.RatingNeutral, // 0 (but Neutral is not positive, so filtered out)
			action:     "target lowered by", // -2
			expectedMin: -15.0,
			expectedMax: -5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := createTestStock("TEST", tt.targetFrom, tt.targetTo, tt.ratingFrom, tt.ratingTo, tt.action)
			
			// Only test if rating is positive (otherwise it will be filtered out)
			if s.RatingTo.IsPositive() {
				recommendations, err := algorithm.CalculateRecommendations(context.Background(), []*stock.Stock{s}, 1)
				assert.NoError(t, err)
				
				if len(recommendations) > 0 {
					recommendation := recommendations[0]
					assert.NotNil(t, recommendation)
					assert.Equal(t, s, recommendation.Stock)
					assert.GreaterOrEqual(t, recommendation.Score, tt.expectedMin)
					assert.LessOrEqual(t, recommendation.Score, tt.expectedMax)
				}
			}
		})
	}
}

// Helper function to create test stocks
func createTestStock(
	ticker string,
	targetFrom, targetTo float64,
	ratingFrom, ratingTo stock.Rating,
	action string,
) *stock.Stock {
	tf, _ := stock.NewPrice(targetFrom)
	tt, _ := stock.NewPrice(targetTo)
	
	s, _ := stock.NewStock(
		ticker,
		"Test Company",
		"Test Brokerage",
		action,
		ratingFrom,
		ratingTo,
		tf,
		tt,
	)
	
	return s
}
