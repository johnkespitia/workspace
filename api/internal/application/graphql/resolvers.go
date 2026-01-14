package graphql

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/john/go-react-test/api/internal/application/services"
	"github.com/john/go-react-test/api/internal/domain/stock"
)

// Resolver contiene los resolvers de GraphQL
type Resolver struct {
	stockService         *services.StockService
	syncService         *services.SyncService
	recommendationService *services.RecommendationService
}

// NewResolver crea un nuevo resolver
func NewResolver(
	stockService *services.StockService,
	syncService *services.SyncService,
	recommendationService *services.RecommendationService,
) *Resolver {
	return &Resolver{
		stockService:         stockService,
		syncService:          syncService,
		recommendationService: recommendationService,
	}
}

// Stocks resuelve la query stocks
func (r *Resolver) Stocks(p graphql.ResolveParams) (interface{}, error) {
	ctx := p.Context

	// Obtener argumentos con manejo robusto de tipos
	limit := 50
	if l, ok := p.Args["limit"]; ok && l != nil {
		switch v := l.(type) {
		case int:
			limit = v
		case int32:
			limit = int(v)
		case int64:
			limit = int(v)
		case float64:
			limit = int(v)
		}
	}

	offset := 0
	if o, ok := p.Args["offset"]; ok && o != nil {
		switch v := o.(type) {
		case int:
			offset = v
		case int32:
			offset = int(v)
		case int64:
			offset = int(v)
		case float64:
			offset = int(v)
		}
	}

	// Convertir filtros con manejo robusto de tipos
	domainFilter := stock.Filter{}
	if filter, ok := p.Args["filter"].(map[string]interface{}); ok && filter != nil {
		// Ticker
		if tickerVal, ok := filter["ticker"]; ok && tickerVal != nil {
			if ticker, ok := tickerVal.(string); ok && ticker != "" {
				domainFilter.Ticker = ticker
			}
		}
		
		// CompanyName
		if companyNameVal, ok := filter["companyName"]; ok && companyNameVal != nil {
			if companyName, ok := companyNameVal.(string); ok && companyName != "" {
				domainFilter.CompanyName = companyName
			}
		}
		
		// Action
		if actionVal, ok := filter["action"]; ok && actionVal != nil {
			if action, ok := actionVal.(string); ok && action != "" {
				domainFilter.Action = action
			}
		}
		
		// Ratings - Solo usar "ratings" (plural)
		if ratingsVal, ok := filter["ratings"]; ok && ratingsVal != nil {
			if ratings, ok := ratingsVal.([]interface{}); ok {
				domainFilter.Ratings = make([]stock.Rating, 0, len(ratings))
				for _, r := range ratings {
					if ratingStr, ok := r.(string); ok && ratingStr != "" {
						domainFilter.Ratings = append(domainFilter.Ratings, stock.Rating(ratingStr))
					}
				}
			}
		}
	}

	// Convertir ordenamiento con valores por defecto
	domainSort := stock.Sort{
		Field:     "created_at",
		Direction: "desc",
	}
	if sort, ok := p.Args["sort"].(map[string]interface{}); ok && sort != nil {
		// El field viene como enum, necesitamos convertirlo
		if fieldVal, ok := sort["field"]; ok && fieldVal != nil {
			fieldStr := ""
			// graphql-go puede pasar el enum como string directamente (el valor del enum)
			if field, ok := fieldVal.(string); ok && field != "" {
				fieldStr = field
			} else {
				// Si viene como otro tipo, convertir a string
				fieldStr = fmt.Sprintf("%v", fieldVal)
			}
			// Mapear el enum a string de campo de BD y validar
			mappedField := mapEnumFieldToDBField(fieldStr)
			if mappedField != "" {
				domainSort.Field = mappedField
			}
		}
		
		// El direction también viene como enum
		if directionVal, ok := sort["direction"]; ok && directionVal != nil {
			directionStr := ""
			if direction, ok := directionVal.(string); ok && direction != "" {
				directionStr = direction
			} else {
				directionStr = fmt.Sprintf("%v", directionVal)
			}
			// Mapear el enum a string de dirección de BD y validar
			mappedDirection := mapEnumDirectionToDBDirection(directionStr)
			if mappedDirection != "" {
				domainSort.Direction = mappedDirection
			}
		}
	}

	// Obtener stocks
	stocks, err := r.stockService.GetStocks(ctx, domainFilter, domainSort)
	if err != nil {
		return nil, err
	}

	// Contar total
	totalCount, err := r.stockService.CountStocks(ctx, domainFilter)
	if err != nil {
		return nil, err
	}

	// Aplicar paginación
	start := offset
	end := start + limit
	if end > len(stocks) {
		end = len(stocks)
	}

	var paginatedStocks []*stock.Stock
	if start < len(stocks) {
		paginatedStocks = stocks[start:end]
	}

	// Convertir a formato GraphQL
	graphqlStocks := make([]map[string]interface{}, len(paginatedStocks))
	for i, s := range paginatedStocks {
		graphqlStocks[i] = stockToMap(s)
	}

	// Calcular información de paginación
	hasNextPage := end < totalCount
	hasPreviousPage := offset > 0

	return map[string]interface{}{
		"stocks": graphqlStocks,
		"totalCount": totalCount,
		"pageInfo": map[string]interface{}{
			"hasNextPage":     hasNextPage,
			"hasPreviousPage": hasPreviousPage,
		},
	}, nil
}

// Stock resuelve la query stock
func (r *Resolver) Stock(p graphql.ResolveParams) (interface{}, error) {
	ctx := p.Context

	ticker, ok := p.Args["ticker"].(string)
	if !ok {
		return nil, fmt.Errorf("ticker is required")
	}

	s, err := r.stockService.GetStock(ctx, ticker)
	if err != nil {
		return nil, err
	}

	return stockToMap(s), nil
}

// Recommendations resuelve la query recommendations
func (r *Resolver) Recommendations(p graphql.ResolveParams) (interface{}, error) {
	ctx := p.Context

	limit := 10
	if l, ok := p.Args["limit"].(int); ok {
		limit = l
	}

	recommendations, err := r.recommendationService.GetRecommendations(ctx, limit)
	if err != nil {
		return nil, err
	}

	// Convertir a formato GraphQL
	result := make([]map[string]interface{}, len(recommendations))
	for i, rec := range recommendations {
		result[i] = map[string]interface{}{
			"stock":       stockToMap(rec.Stock),
			"score":       rec.Score,
			"priceChange": rec.PriceChange,
			"ratingScore": rec.RatingScore,
			"actionScore": rec.ActionScore,
		}
	}

	return result, nil
}

// SyncStocks resuelve la mutation syncStocks
func (r *Resolver) SyncStocks(p graphql.ResolveParams) (interface{}, error) {
	ctx := p.Context

	count, err := r.syncService.SyncAllStocks(ctx)
	if err != nil {
		return map[string]interface{}{
			"success":      false,
			"message":      err.Error(),
			"stocksSynced": 0,
		}, nil
	}

	return map[string]interface{}{
		"success":      true,
		"message":      "Stocks synchronized successfully",
		"stocksSynced": count,
	}, nil
}

// stockToMap convierte un stock de dominio a mapa para GraphQL
func stockToMap(s *stock.Stock) map[string]interface{} {
	brokerage := s.Brokerage
	action := s.Action

	return map[string]interface{}{
		"id":          s.ID.String(),
		"ticker":      s.Ticker,
		"companyName": s.CompanyName,
		"brokerage":   brokerage,
		"action":      action,
		"ratingFrom":  s.RatingFrom.String(),
		"ratingTo":    s.RatingTo.String(),
		"targetFrom":  s.TargetFrom.Value(),
		"targetTo":    s.TargetTo.Value(),
		"createdAt":   s.CreatedAt,
		"updatedAt":   s.UpdatedAt,
	}
}

// mapEnumFieldToDBField mapea el enum de GraphQL al nombre de campo de la BD
func mapEnumFieldToDBField(enumValue string) string {
	if enumValue == "" {
		return ""
	}
	
	// Normalizar a mayúsculas para comparación
	upper := strings.ToUpper(strings.TrimSpace(enumValue))
	
	switch upper {
	case "TICKER":
		return "ticker"
	case "COMPANY_NAME":
		return "company_name"
	case "RATING_TO":
		return "rating_to"
	case "TARGET_TO":
		return "target_to"
	case "CREATED_AT":
		return "created_at"
	default:
		// Si ya viene como nombre de campo válido, retornarlo tal cual
		// Validar que sea uno de los campos permitidos
		validFields := map[string]bool{
			"ticker":       true,
			"company_name": true,
			"rating_to":    true,
			"target_to":    true,
			"created_at":   true,
		}
		lower := strings.ToLower(enumValue)
		if validFields[lower] {
			return lower
		}
		// Si no es válido, retornar campo por defecto
		return "created_at"
	}
}

// mapEnumDirectionToDBDirection mapea el enum de dirección al string de la BD
func mapEnumDirectionToDBDirection(enumValue string) string {
	if enumValue == "" {
		return ""
	}
	
	upper := strings.ToUpper(strings.TrimSpace(enumValue))
	
	switch upper {
	case "ASC":
		return "asc"
	case "DESC":
		return "desc"
	default:
		// Si ya viene como dirección válida, retornarlo tal cual
		lower := strings.ToLower(enumValue)
		if lower == "asc" || lower == "desc" {
			return lower
		}
		// Si no es válido, retornar dirección por defecto
		return "desc"
	}
}
