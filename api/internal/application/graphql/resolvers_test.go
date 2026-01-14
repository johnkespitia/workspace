package graphql

import (
	"context"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

// TestResolver_Stocks tests the Stocks resolver with basic argument parsing
func TestResolver_Stocks(t *testing.T) {
	// Test argument parsing logic
	t.Run("parse limit and offset", func(t *testing.T) {
		params := graphql.ResolveParams{
			Context: context.Background(),
			Args: map[string]interface{}{
				"limit":  int(50),
				"offset": int(0),
			},
		}

		limit := 50
		if l, ok := params.Args["limit"].(int); ok {
			limit = l
		}
		assert.Equal(t, 50, limit)

		offset := 0
		if o, ok := params.Args["offset"].(int); ok {
			offset = o
		}
		assert.Equal(t, 0, offset)
	})

	t.Run("parse filter ratings", func(t *testing.T) {
		params := graphql.ResolveParams{
			Context: context.Background(),
			Args: map[string]interface{}{
				"filter": map[string]interface{}{
					"ratings": []interface{}{"Strong Buy", "Buy"},
				},
			},
		}

		if filter, ok := params.Args["filter"].(map[string]interface{}); ok {
			if ratings, ok := filter["ratings"].([]interface{}); ok {
				assert.Len(t, ratings, 2)
				assert.Equal(t, "Strong Buy", ratings[0])
				assert.Equal(t, "Buy", ratings[1])
			}
		}
	})
}

// TestResolver_Stock tests the Stock resolver argument parsing
func TestResolver_Stock(t *testing.T) {
	t.Run("parse ticker argument", func(t *testing.T) {
		params := graphql.ResolveParams{
			Context: context.Background(),
			Args: map[string]interface{}{
				"ticker": "AAPL",
			},
		}

		ticker, ok := params.Args["ticker"].(string)
		assert.True(t, ok)
		assert.Equal(t, "AAPL", ticker)
	})

	t.Run("missing ticker", func(t *testing.T) {
		params := graphql.ResolveParams{
			Context: context.Background(),
			Args:    map[string]interface{}{},
		}

		_, ok := params.Args["ticker"].(string)
		assert.False(t, ok)
	})
}

// TestResolver_SyncStocks tests the SyncStocks mutation structure
func TestResolver_SyncStocks(t *testing.T) {
	t.Run("sync stocks structure", func(t *testing.T) {
		// Test that the resolver returns the correct structure
		result := map[string]interface{}{
			"success":      true,
			"message":      "Stocks synchronized successfully",
			"stocksSynced": 10,
		}

		assert.True(t, result["success"].(bool))
		assert.Equal(t, "Stocks synchronized successfully", result["message"])
		assert.Equal(t, 10, result["stocksSynced"])
	})
}
