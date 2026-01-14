package graphql

import (
	"context"
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/john/go-react-test/api/internal/application/services"
	"github.com/john/go-react-test/api/internal/domain/stock"
)

// StockLoaderKey es el tipo de clave para el DataLoader de stocks
type StockLoaderKey string

// String implementa la interfaz Key
func (k StockLoaderKey) String() string {
	return string(k)
}

// Key retorna el string del ticker
func (k StockLoaderKey) Key() string {
	return string(k)
}

// StockLoader es el DataLoader para cargar stocks por ticker
type StockLoader struct {
	loader *dataloader.Loader[StockLoaderKey, *stock.Stock]
}

// NewStockLoader crea un nuevo DataLoader para stocks
func NewStockLoader(stockService *services.StockService) *StockLoader {
	return &StockLoader{
		loader: dataloader.NewBatchedLoader(
			func(ctx context.Context, keys []StockLoaderKey) []*dataloader.Result[*stock.Stock] {
				// Convertir keys a strings
				tickers := make([]string, len(keys))
				for i, key := range keys {
					tickers[i] = key.String()
				}

				// Obtener stocks en batch
				stocks, err := stockService.GetStocksByTickers(ctx, tickers)
				if err != nil {
					// Si hay error, retornar error para todos
					results := make([]*dataloader.Result[*stock.Stock], len(keys))
					for i := range results {
						results[i] = &dataloader.Result[*stock.Stock]{
							Error: err,
						}
					}
					return results
				}

				// Crear mapa de ticker -> stock para lookup rápido
				stockMap := make(map[string]*stock.Stock)
				for _, s := range stocks {
					if s != nil {
						stockMap[s.Ticker] = s
					}
				}

				// Crear resultados manteniendo el orden de las keys
				results := make([]*dataloader.Result[*stock.Stock], len(keys))
				for i, key := range keys {
					if s, ok := stockMap[key.String()]; ok {
						results[i] = &dataloader.Result[*stock.Stock]{
							Data: s,
						}
					} else {
						// Stock no encontrado - retornar nil en lugar de error para permitir que el resolver maneje el caso
						results[i] = &dataloader.Result[*stock.Stock]{
							Data: nil,
						}
					}
				}

				return results
			},
			dataloader.WithBatchCapacity[StockLoaderKey, *stock.Stock](100),
			dataloader.WithWait[StockLoaderKey, *stock.Stock](16*time.Millisecond), // Esperar 16ms para batch
		),
	}
}

// Load carga un stock por ticker usando el DataLoader
func (l *StockLoader) Load(ctx context.Context, ticker string) (*stock.Stock, error) {
	key := StockLoaderKey(ticker)
	thunk := l.loader.Load(ctx, key)
	return thunk()
}

// LoadMany carga múltiples stocks por tickers usando el DataLoader
func (l *StockLoader) LoadMany(ctx context.Context, tickers []string) ([]*stock.Stock, []error) {
	keys := make([]StockLoaderKey, len(tickers))
	for i, ticker := range tickers {
		keys[i] = StockLoaderKey(ticker)
	}

	thunk := l.loader.LoadMany(ctx, keys)
	results, errs := thunk()

	stocks := make([]*stock.Stock, 0, len(results))
	errors := make([]error, 0, len(errs))

	for i, result := range results {
		if result != nil {
			stocks = append(stocks, result)
		} else if i < len(errs) && errs[i] != nil {
			errors = append(errors, errs[i])
		}
	}

	return stocks, errors
}

// Clear elimina un stock del cache del DataLoader
func (l *StockLoader) Clear(ctx context.Context, ticker string) {
	key := StockLoaderKey(ticker)
	l.loader.Clear(ctx, key)
}

// ClearAll elimina todos los stocks del cache del DataLoader
func (l *StockLoader) ClearAll() {
	l.loader.ClearAll()
}

// Prime agrega un stock al cache del DataLoader
func (l *StockLoader) Prime(ctx context.Context, ticker string, stock *stock.Stock) {
	key := StockLoaderKey(ticker)
	l.loader.Prime(ctx, key, stock)
}
