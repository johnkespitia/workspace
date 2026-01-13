package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/john/go-react-test/api/internal/domain/stock"
	"github.com/john/go-react-test/api/internal/infrastructure/database"
)

// CockroachStockRepository implementa el repositorio de stocks para CockroachDB
type CockroachStockRepository struct {
	db *sql.DB
}

// NewCockroachStockRepository crea un nuevo repositorio
func NewCockroachStockRepository() stock.Repository {
	return &CockroachStockRepository{
		db: database.GetDB(),
	}
}

// Save guarda o actualiza una acción (UPSERT)
func (r *CockroachStockRepository) Save(ctx context.Context, s *stock.Stock) error {
	query := `
		INSERT INTO stocks (
			id, ticker, company_name, brokerage, action,
			rating_from, rating_to, target_from, target_to,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (ticker) 
		DO UPDATE SET
			company_name = EXCLUDED.company_name,
			brokerage = EXCLUDED.brokerage,
			action = EXCLUDED.action,
			rating_from = EXCLUDED.rating_from,
			rating_to = EXCLUDED.rating_to,
			target_from = EXCLUDED.target_from,
			target_to = EXCLUDED.target_to,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		s.ID,
		s.Ticker,
		s.CompanyName,
		s.Brokerage,
		s.Action,
		s.RatingFrom.String(),
		s.RatingTo.String(),
		s.TargetFrom.Value(),
		s.TargetTo.Value(),
		s.CreatedAt,
		s.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save stock: %w", err)
	}

	return nil
}

// BatchUpsert guarda o actualiza múltiples acciones en batch
func (r *CockroachStockRepository) BatchUpsert(ctx context.Context, stocks []*stock.Stock) error {
	if len(stocks) == 0 {
		return nil
	}

	batchSize := 100
	for i := 0; i < len(stocks); i += batchSize {
		end := i + batchSize
		if end > len(stocks) {
			end = len(stocks)
		}

		batch := stocks[i:end]
		if err := r.upsertBatch(ctx, batch); err != nil {
			return fmt.Errorf("failed to upsert batch %d-%d: %w", i, end, err)
		}
	}

	return nil
}

// upsertBatch realiza un upsert de un batch de stocks
func (r *CockroachStockRepository) upsertBatch(ctx context.Context, stocks []*stock.Stock) error {
	if len(stocks) == 0 {
		return nil
	}

	// Construir query con múltiples valores
	valueStrings := make([]string, 0, len(stocks))
	valueArgs := make([]interface{}, 0, len(stocks)*11)

	for i, s := range stocks {
		offset := i * 11
		valueStrings = append(valueStrings, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			offset+1, offset+2, offset+3, offset+4, offset+5,
			offset+6, offset+7, offset+8, offset+9, offset+10, offset+11,
		))

		valueArgs = append(valueArgs,
			s.ID,
			s.Ticker,
			s.CompanyName,
			s.Brokerage,
			s.Action,
			s.RatingFrom.String(),
			s.RatingTo.String(),
			s.TargetFrom.Value(),
			s.TargetTo.Value(),
			s.CreatedAt,
			s.UpdatedAt,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO stocks (
			id, ticker, company_name, brokerage, action,
			rating_from, rating_to, target_from, target_to,
			created_at, updated_at
		) VALUES %s
		ON CONFLICT (ticker) 
		DO UPDATE SET
			company_name = EXCLUDED.company_name,
			brokerage = EXCLUDED.brokerage,
			action = EXCLUDED.action,
			rating_from = EXCLUDED.rating_from,
			rating_to = EXCLUDED.rating_to,
			target_from = EXCLUDED.target_from,
			target_to = EXCLUDED.target_to,
			updated_at = EXCLUDED.updated_at
	`, strings.Join(valueStrings, ","))

	_, err := r.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to upsert batch: %w", err)
	}

	return nil
}

// FindByID busca una acción por ID
func (r *CockroachStockRepository) FindByID(ctx context.Context, id uuid.UUID) (*stock.Stock, error) {
	query := `
		SELECT id, ticker, company_name, brokerage, action,
		       rating_from, rating_to, target_from, target_to,
		       created_at, updated_at
		FROM stocks
		WHERE id = $1
	`

	var s stock.Stock
	var ratingFromStr, ratingToStr string
	var targetFromVal, targetToVal float64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&s.Ticker,
		&s.CompanyName,
		&s.Brokerage,
		&s.Action,
		&ratingFromStr,
		&ratingToStr,
		&targetFromVal,
		&targetToVal,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("stock not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find stock: %w", err)
	}

	s.RatingFrom = stock.Rating(ratingFromStr)
	s.RatingTo = stock.Rating(ratingToStr)

	targetFrom, err := stock.NewPrice(targetFromVal)
	if err != nil {
		return nil, fmt.Errorf("invalid target_from: %w", err)
	}
	s.TargetFrom = targetFrom

	targetTo, err := stock.NewPrice(targetToVal)
	if err != nil {
		return nil, fmt.Errorf("invalid target_to: %w", err)
	}
	s.TargetTo = targetTo

	return &s, nil
}

// FindByTicker busca una acción por ticker
func (r *CockroachStockRepository) FindByTicker(ctx context.Context, ticker string) (*stock.Stock, error) {
	query := `
		SELECT id, ticker, company_name, brokerage, action,
		       rating_from, rating_to, target_from, target_to,
		       created_at, updated_at
		FROM stocks
		WHERE ticker = $1
	`

	var s stock.Stock
	var ratingFromStr, ratingToStr string
	var targetFromVal, targetToVal float64

	err := r.db.QueryRowContext(ctx, query, ticker).Scan(
		&s.ID,
		&s.Ticker,
		&s.CompanyName,
		&s.Brokerage,
		&s.Action,
		&ratingFromStr,
		&ratingToStr,
		&targetFromVal,
		&targetToVal,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("stock not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find stock: %w", err)
	}

	s.RatingFrom = stock.Rating(ratingFromStr)
	s.RatingTo = stock.Rating(ratingToStr)

	targetFrom, err := stock.NewPrice(targetFromVal)
	if err != nil {
		return nil, fmt.Errorf("invalid target_from: %w", err)
	}
	s.TargetFrom = targetFrom

	targetTo, err := stock.NewPrice(targetToVal)
	if err != nil {
		return nil, fmt.Errorf("invalid target_to: %w", err)
	}
	s.TargetTo = targetTo

	return &s, nil
}

// FindAll busca todas las acciones con filtros y ordenamiento
func (r *CockroachStockRepository) FindAll(ctx context.Context, filter stock.Filter, sort stock.Sort) ([]*stock.Stock, error) {
	query := "SELECT id, ticker, company_name, brokerage, action, rating_from, rating_to, target_from, target_to, created_at, updated_at FROM stocks WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	// Aplicar filtros
	if filter.Ticker != "" {
		query += fmt.Sprintf(" AND ticker = $%d", argIndex)
		args = append(args, filter.Ticker)
		argIndex++
	}

	if filter.CompanyName != "" {
		query += fmt.Sprintf(" AND company_name ILIKE $%d", argIndex)
		args = append(args, "%"+filter.CompanyName+"%")
		argIndex++
	}

	if len(filter.Ratings) > 0 {
		placeholders := make([]string, len(filter.Ratings))
		for i, rating := range filter.Ratings {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, rating.String())
			argIndex++
		}
		query += fmt.Sprintf(" AND rating_to = ANY(ARRAY[%s])", strings.Join(placeholders, ","))
	}

	if filter.Action != "" {
		query += fmt.Sprintf(" AND action = $%d", argIndex)
		args = append(args, filter.Action)
		argIndex++
	}

	// Aplicar ordenamiento
	if sort.Field != "" {
		validFields := map[string]string{
			"ticker":       "ticker",
			"company_name": "company_name",
			"rating_to":    "rating_to",
			"target_to":    "target_to",
			"created_at":   "created_at",
		}

		if field, ok := validFields[sort.Field]; ok {
			direction := "ASC"
			if sort.Direction == "desc" {
				direction = "DESC"
			}
			query += fmt.Sprintf(" ORDER BY %s %s", field, direction)
		}
	} else {
		query += " ORDER BY created_at DESC"
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query stocks: %w", err)
	}
	defer rows.Close()

	var stocks []*stock.Stock
	for rows.Next() {
		var s stock.Stock
		var ratingFromStr, ratingToStr string
		var targetFromVal, targetToVal float64

		err := rows.Scan(
			&s.ID,
			&s.Ticker,
			&s.CompanyName,
			&s.Brokerage,
			&s.Action,
			&ratingFromStr,
			&ratingToStr,
			&targetFromVal,
			&targetToVal,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock: %w", err)
		}

		s.RatingFrom = stock.Rating(ratingFromStr)
		s.RatingTo = stock.Rating(ratingToStr)

		targetFrom, err := stock.NewPrice(targetFromVal)
		if err != nil {
			continue // Skip invalid price
		}
		s.TargetFrom = targetFrom

		targetTo, err := stock.NewPrice(targetToVal)
		if err != nil {
			continue // Skip invalid price
		}
		s.TargetTo = targetTo

		stocks = append(stocks, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return stocks, nil
}

// Count cuenta el número de acciones que coinciden con el filtro
func (r *CockroachStockRepository) Count(ctx context.Context, filter stock.Filter) (int, error) {
	query := "SELECT COUNT(*) FROM stocks WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	// Aplicar los mismos filtros que FindAll
	if filter.Ticker != "" {
		query += fmt.Sprintf(" AND ticker = $%d", argIndex)
		args = append(args, filter.Ticker)
		argIndex++
	}

	if filter.CompanyName != "" {
		query += fmt.Sprintf(" AND company_name ILIKE $%d", argIndex)
		args = append(args, "%"+filter.CompanyName+"%")
		argIndex++
	}

	if len(filter.Ratings) > 0 {
		placeholders := make([]string, len(filter.Ratings))
		for i, rating := range filter.Ratings {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, rating.String())
			argIndex++
		}
		query += fmt.Sprintf(" AND rating_to = ANY(ARRAY[%s])", strings.Join(placeholders, ","))
	}

	if filter.Action != "" {
		query += fmt.Sprintf(" AND action = $%d", argIndex)
		args = append(args, filter.Action)
		argIndex++
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count stocks: %w", err)
	}

	return count, nil
}
