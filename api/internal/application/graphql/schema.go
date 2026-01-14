package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/john/go-react-test/api/internal/application/services"
)

// Schema contiene el schema GraphQL completo
type Schema struct {
	schema graphql.Schema
}

// NewSchema crea un nuevo schema GraphQL
func NewSchema(
	stockService *services.StockService,
	syncService *services.SyncService,
	recommendationService *services.RecommendationService,
) (*Schema, error) {
	resolver := NewResolver(stockService, syncService, recommendationService)

	schema, err := buildSchema(resolver)
	if err != nil {
		return nil, err
	}

	return &Schema{schema: schema}, nil
}

// GetSchema retorna el schema de graphql-go
func (s *Schema) GetSchema() graphql.Schema {
	return s.schema
}

// buildSchema construye el schema GraphQL
func buildSchema(resolver *Resolver) (graphql.Schema, error) {
	// Definir tipos
	stockType := defineStockType()
	recommendationType := defineRecommendationType(stockType)
	stockConnectionType := defineStockConnectionType(stockType)
	syncStocksResultType := defineSyncStocksResultType()

	// Definir inputs
	stockFilterInput := defineStockFilterInput()
	stockSortInput := defineStockSortInput()

	// Definir queries
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"stocks": &graphql.Field{
				Type: stockConnectionType,
				Args: graphql.FieldConfigArgument{
					"filter": &graphql.ArgumentConfig{
						Type: stockFilterInput,
					},
					"sort": &graphql.ArgumentConfig{
						Type: stockSortInput,
					},
					"limit": &graphql.ArgumentConfig{
						Type: graphql.Int,
						DefaultValue: 50,
					},
					"offset": &graphql.ArgumentConfig{
						Type: graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: resolver.Stocks,
			},
			"stock": &graphql.Field{
				Type: stockType,
				Args: graphql.FieldConfigArgument{
					"ticker": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: resolver.Stock,
			},
			"recommendations": &graphql.Field{
				Type: graphql.NewList(recommendationType),
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type: graphql.Int,
						DefaultValue: 10,
					},
				},
				Resolve: resolver.Recommendations,
			},
		},
	})

	// Definir mutations
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"syncStocks": &graphql.Field{
				Type: syncStocksResultType,
				Resolve: resolver.SyncStocks,
			},
		},
	})

	// Crear schema
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
}

// defineStockType define el tipo Stock
func defineStockType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Stock",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"ticker": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"companyName": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"brokerage": &graphql.Field{
				Type: graphql.String,
			},
			"action": &graphql.Field{
				Type: graphql.String,
			},
			"ratingFrom": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"ratingTo": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"targetFrom": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
			},
			"targetTo": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
			},
			"createdAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"updatedAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
		},
	})
}

// defineRecommendationType define el tipo Recommendation
func defineRecommendationType(stockType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Recommendation",
		Fields: graphql.Fields{
			"stock": &graphql.Field{
				Type: graphql.NewNonNull(stockType),
			},
			"score": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
			},
			"priceChange": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
			},
			"ratingScore": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
			},
			"actionScore": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
			},
		},
	})
}

// defineStockConnectionType define el tipo StockConnection
func defineStockConnectionType(stockType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "StockConnection",
		Fields: graphql.Fields{
			"stocks": &graphql.Field{
				Type: graphql.NewList(stockType),
			},
			"totalCount": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"pageInfo": &graphql.Field{
				Type: graphql.NewNonNull(definePageInfoType()),
			},
		},
	})
}

// definePageInfoType define el tipo PageInfo
func definePageInfoType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "PageInfo",
		Fields: graphql.Fields{
			"hasNextPage": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"hasPreviousPage": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
		},
	})
}

// defineSyncStocksResultType define el tipo SyncStocksResult
func defineSyncStocksResultType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "SyncStocksResult",
		Fields: graphql.Fields{
			"success": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"message": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"stocksSynced": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
	})
}

// defineStockFilterInput define el input StockFilter
func defineStockFilterInput() *graphql.InputObject {
	return graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "StockFilter",
		Fields: graphql.InputObjectConfigFieldMap{
			"ticker": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"companyName": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"ratings": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(graphql.NewNonNull(graphql.String)),
			},
			"action": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	})
}

// defineStockSortInput define el input StockSort
func defineStockSortInput() *graphql.InputObject {
	stockSortFieldEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "StockSortField",
		Values: graphql.EnumValueConfigMap{
			"TICKER": &graphql.EnumValueConfig{
				Value: "ticker",
			},
			"COMPANY_NAME": &graphql.EnumValueConfig{
				Value: "company_name",
			},
			"RATING_TO": &graphql.EnumValueConfig{
				Value: "rating_to",
			},
			"TARGET_TO": &graphql.EnumValueConfig{
				Value: "target_to",
			},
			"CREATED_AT": &graphql.EnumValueConfig{
				Value: "created_at",
			},
			// También aceptar valores en minúsculas directamente
			"ticker": &graphql.EnumValueConfig{
				Value: "ticker",
			},
			"company_name": &graphql.EnumValueConfig{
				Value: "company_name",
			},
			"rating_to": &graphql.EnumValueConfig{
				Value: "rating_to",
			},
			"target_to": &graphql.EnumValueConfig{
				Value: "target_to",
			},
			"created_at": &graphql.EnumValueConfig{
				Value: "created_at",
			},
		},
	})

	sortDirectionEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "SortDirection",
		Values: graphql.EnumValueConfigMap{
			"ASC": &graphql.EnumValueConfig{
				Value: "asc",
			},
			"DESC": &graphql.EnumValueConfig{
				Value: "desc",
			},
			// También aceptar valores en minúsculas directamente
			"asc": &graphql.EnumValueConfig{
				Value: "asc",
			},
			"desc": &graphql.EnumValueConfig{
				Value: "desc",
			},
		},
	})

	return graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "StockSort",
		Fields: graphql.InputObjectConfigFieldMap{
			"field": &graphql.InputObjectFieldConfig{
				Type: stockSortFieldEnum, // Opcional, no NewNonNull
			},
			"direction": &graphql.InputObjectFieldConfig{
				Type: sortDirectionEnum, // Opcional, no NewNonNull
			},
		},
	})
}
