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

// GetStocksByTickers obtiene múltiples stocks por sus tickers (para DataLoader)
func (s *StockService) GetStocksByTickers(ctx context.Context, tickers []string) ([]*stock.Stock, error) {
	if len(tickers) == 0 {
		return []*stock.Stock{}, nil
	}

	// Crear un mapa para resultados
	resultMap := make(map[string]*stock.Stock)
	
	// Obtener cada stock por ticker
	for _, ticker := range tickers {
		stock, err := s.repo.FindByTicker(ctx, ticker)
		if err != nil {
			// Si no se encuentra, continuar con el siguiente
			continue
		}
		resultMap[ticker] = stock
	}

	// Convertir mapa a slice manteniendo el orden de los tickers
	result := make([]*stock.Stock, 0, len(tickers))
	for _, ticker := range tickers {
		if stock, ok := resultMap[ticker]; ok {
			result = append(result, stock)
		}
	}

	return result, nil
}
