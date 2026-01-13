package stock

import (
	"context"

	"github.com/google/uuid"
)

// Filter representa los filtros para búsqueda de stocks
type Filter struct {
	Ticker      string
	CompanyName string
	Ratings     []Rating
	Action      string
}

// Sort representa el ordenamiento para búsqueda de stocks
type Sort struct {
	Field     string // "ticker", "company_name", "rating_to", "target_to", "created_at"
	Direction string // "asc", "desc"
}

// Repository define la interfaz del repositorio de stocks
type Repository interface {
	// Save guarda o actualiza una acción
	Save(ctx context.Context, stock *Stock) error

	// BatchUpsert guarda o actualiza múltiples acciones en batch
	BatchUpsert(ctx context.Context, stocks []*Stock) error

	// FindByID busca una acción por ID
	FindByID(ctx context.Context, id uuid.UUID) (*Stock, error)

	// FindByTicker busca una acción por ticker
	FindByTicker(ctx context.Context, ticker string) (*Stock, error)

	// FindAll busca todas las acciones con filtros y ordenamiento
	FindAll(ctx context.Context, filter Filter, sort Sort) ([]*Stock, error)

	// Count cuenta el número de acciones que coinciden con el filtro
	Count(ctx context.Context, filter Filter) (int, error)
}
