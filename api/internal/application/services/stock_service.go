package services

import (
	"context"

	"github.com/john/go-react-test/api/internal/domain/stock"
)

// StockService es el servicio de aplicación para stocks
type StockService struct {
	repo       stock.Repository
	domainSvc  stock.Service
}

// NewStockService crea un nuevo servicio de stocks
func NewStockService(repo stock.Repository, domainSvc stock.Service) *StockService {
	return &StockService{
		repo:      repo,
		domainSvc: domainSvc,
	}
}

// GetStocks obtiene stocks con filtros y ordenamiento
func (s *StockService) GetStocks(ctx context.Context, filter stock.Filter, sort stock.Sort) ([]*stock.Stock, error) {
	return s.repo.FindAll(ctx, filter, sort)
}

// GetStock obtiene un stock por ticker
func (s *StockService) GetStock(ctx context.Context, ticker string) (*stock.Stock, error) {
	return s.repo.FindByTicker(ctx, ticker)
}

// CountStocks cuenta el número de stocks que coinciden con el filtro
func (s *StockService) CountStocks(ctx context.Context, filter stock.Filter) (int, error) {
	return s.repo.Count(ctx, filter)
}
