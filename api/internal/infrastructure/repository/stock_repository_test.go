package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/john/go-react-test/api/internal/domain/stock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCockroachStockRepository_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &CockroachStockRepository{db: db}

	targetFrom, _ := stock.NewPrice(100.0)
	targetTo, _ := stock.NewPrice(120.0)
	
	s := &stock.Stock{
		ID:          uuid.New(),
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		Brokerage:   "Test Brokerage",
		Action:      "target raised by",
		RatingFrom:  stock.RatingBuy,
		RatingTo:    stock.RatingStrongBuy,
		TargetFrom:  targetFrom,
		TargetTo:    targetTo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test INSERT (new stock)
	t.Run("insert new stock", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO stocks`).
			WithArgs(
				s.ID, s.Ticker, s.CompanyName, s.Brokerage, s.Action,
				s.RatingFrom.String(), s.RatingTo.String(),
				s.TargetFrom.Value(), s.TargetTo.Value(),
				s.CreatedAt, s.UpdatedAt,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(context.Background(), s)
		assert.NoError(t, err)
	})

	// Test UPDATE (existing stock)
	t.Run("update existing stock", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO stocks`).
			WithArgs(
				s.ID, s.Ticker, s.CompanyName, s.Brokerage, s.Action,
				s.RatingFrom.String(), s.RatingTo.String(),
				s.TargetFrom.Value(), s.TargetTo.Value(),
				s.CreatedAt, s.UpdatedAt,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(context.Background(), s)
		assert.NoError(t, err)
	})

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCockroachStockRepository_FindByTicker(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &CockroachStockRepository{db: db}

	t.Run("find existing stock", func(t *testing.T) {
		ticker := "AAPL"
		stockID := uuid.New()
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "ticker", "company_name", "brokerage", "action",
			"rating_from", "rating_to", "target_from", "target_to",
			"created_at", "updated_at",
		}).AddRow(
			stockID, ticker, "Apple Inc.", "Test Brokerage", "target raised by",
			"Buy", "Strong Buy", 100.0, 120.0,
			now, now,
		)

		mock.ExpectQuery(`SELECT .+ FROM stocks WHERE ticker = \$1`).
			WithArgs(ticker).
			WillReturnRows(rows)

		result, err := repo.FindByTicker(context.Background(), ticker)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ticker, result.Ticker)
		assert.Equal(t, "Apple Inc.", result.CompanyName)
	})

	t.Run("stock not found", func(t *testing.T) {
		ticker := "NONEXISTENT"

		mock.ExpectQuery(`SELECT .+ FROM stocks WHERE ticker = \$1`).
			WithArgs(ticker).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByTicker(context.Background(), ticker)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "stock not found")
	})

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCockroachStockRepository_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &CockroachStockRepository{db: db}

	t.Run("find all without filters", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"id", "ticker", "company_name", "brokerage", "action",
			"rating_from", "rating_to", "target_from", "target_to",
			"created_at", "updated_at",
		}).
			AddRow(
				uuid.New(), "AAPL", "Apple Inc.", "Brokerage1", "target raised by",
				"Buy", "Strong Buy", 100.0, 120.0, now, now,
			).
			AddRow(
				uuid.New(), "MSFT", "Microsoft Corp.", "Brokerage2", "target raised",
				"Neutral", "Buy", 50.0, 60.0, now, now,
			)

		mock.ExpectQuery(`SELECT .+ FROM stocks WHERE 1=1 ORDER BY created_at DESC`).
			WillReturnRows(rows)

		filter := stock.Filter{}
		sort := stock.Sort{Field: "", Direction: ""}

		result, err := repo.FindAll(context.Background(), filter, sort)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("find with ticker filter", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"id", "ticker", "company_name", "brokerage", "action",
			"rating_from", "rating_to", "target_from", "target_to",
			"created_at", "updated_at",
		}).
			AddRow(
				uuid.New(), "AAPL", "Apple Inc.", "Brokerage1", "target raised by",
				"Buy", "Strong Buy", 100.0, 120.0, now, now,
			)

		mock.ExpectQuery(`SELECT .+ FROM stocks WHERE 1=1 AND ticker = \$1`).
			WithArgs("AAPL").
			WillReturnRows(rows)

		filter := stock.Filter{Ticker: "AAPL"}
		sort := stock.Sort{Field: "", Direction: ""}

		result, err := repo.FindAll(context.Background(), filter, sort)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "AAPL", result[0].Ticker)
	})

	t.Run("find with ratings filter", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"id", "ticker", "company_name", "brokerage", "action",
			"rating_from", "rating_to", "target_from", "target_to",
			"created_at", "updated_at",
		}).
			AddRow(
				uuid.New(), "AAPL", "Apple Inc.", "Brokerage1", "target raised by",
				"Buy", "Strong Buy", 100.0, 120.0, now, now,
			)

		mock.ExpectQuery(`SELECT .+ FROM stocks WHERE 1=1 AND rating_to = ANY`).
			WithArgs("Strong Buy").
			WillReturnRows(rows)

		filter := stock.Filter{Ratings: []stock.Rating{stock.RatingStrongBuy}}
		sort := stock.Sort{Field: "", Direction: ""}

		result, err := repo.FindAll(context.Background(), filter, sort)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("find with sort", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"id", "ticker", "company_name", "brokerage", "action",
			"rating_from", "rating_to", "target_from", "target_to",
			"created_at", "updated_at",
		}).
			AddRow(
				uuid.New(), "AAPL", "Apple Inc.", "Brokerage1", "target raised by",
				"Buy", "Strong Buy", 100.0, 120.0, now, now,
			).
			AddRow(
				uuid.New(), "MSFT", "Microsoft Corp.", "Brokerage2", "target raised",
				"Neutral", "Buy", 50.0, 60.0, now, now,
			)

		mock.ExpectQuery(`SELECT .+ FROM stocks WHERE 1=1 ORDER BY ticker ASC`).
			WillReturnRows(rows)

		filter := stock.Filter{}
		sort := stock.Sort{Field: "ticker", Direction: "asc"}

		result, err := repo.FindAll(context.Background(), filter, sort)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
	})

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCockroachStockRepository_Count(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &CockroachStockRepository{db: db}

	t.Run("count all", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(10)

		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM stocks WHERE 1=1`).
			WillReturnRows(rows)

		filter := stock.Filter{}
		count, err := repo.Count(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, 10, count)
	})

	t.Run("count with filter", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(5)

		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM stocks WHERE 1=1 AND ticker = \$1`).
			WithArgs("AAPL").
			WillReturnRows(rows)

		filter := stock.Filter{Ticker: "AAPL"}
		count, err := repo.Count(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, 5, count)
	})

	require.NoError(t, mock.ExpectationsWereMet())
}
