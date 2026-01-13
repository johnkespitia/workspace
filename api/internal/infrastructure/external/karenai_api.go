package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/john/go-react-test/api/internal/domain/stock"
)

// KarenAIClient es el cliente para la API externa de KarenAI
type KarenAIClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

// APIResponse representa la respuesta de la API externa
type APIResponse struct {
	Stocks   []StockDTO `json:"stocks"`
	NextPage string     `json:"next_page,omitempty"`
}

// StockDTO representa un stock en la respuesta de la API
type StockDTO struct {
	Ticker     string  `json:"TICKER"`
	Company    string  `json:"COMPANY"`
	Brokerage  string  `json:"BROKERAGE"`
	Action     string  `json:"ACTION"`
	RatingFrom string  `json:"RATING FROM"`
	RatingTo   string  `json:"RATING TO"`
	TargetFrom float64 `json:"TARGET FROM"`
	TargetTo   float64 `json:"TARGET TO"`
}

// NewKarenAIClient crea un nuevo cliente para la API externa
func NewKarenAIClient(baseURL, apiKey string) *KarenAIClient {
	return &KarenAIClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

// FetchStocks obtiene stocks de la API externa con paginación
func (c *KarenAIClient) FetchStocks(ctx context.Context, nextPage string) (*APIResponse, error) {
	url := fmt.Sprintf("%s/swechallenge/list", c.baseURL)
	if nextPage != "" {
		url = fmt.Sprintf("%s?next_page=%s", url, nextPage)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResp, nil
}

// FetchAllStocks obtiene todas las páginas de stocks
func (c *KarenAIClient) FetchAllStocks(ctx context.Context) ([]*stock.Stock, error) {
	var allStocks []*stock.Stock
	nextPage := ""
	pageCount := 0

	for {
		response, err := c.FetchStocks(ctx, nextPage)
		if err != nil {
			return nil, fmt.Errorf("error fetching page %d: %w", pageCount, err)
		}

		// Convertir DTOs a entidades de dominio
		for _, dto := range response.Stocks {
			s, err := c.convertToDomainEntity(dto)
			if err != nil {
				// Log error pero continuar con otros stocks
				continue
			}
			allStocks = append(allStocks, s)
		}

		// Verificar si hay más páginas
		if response.NextPage == "" {
			break
		}
		nextPage = response.NextPage
		pageCount++
	}

	return allStocks, nil
}

// convertToDomainEntity convierte un DTO a una entidad de dominio
func (c *KarenAIClient) convertToDomainEntity(dto StockDTO) (*stock.Stock, error) {
	ratingFrom := stock.Rating(dto.RatingFrom)
	ratingTo := stock.Rating(dto.RatingTo)

	targetFrom, err := stock.NewPrice(dto.TargetFrom)
	if err != nil {
		return nil, fmt.Errorf("invalid target_from: %w", err)
	}

	targetTo, err := stock.NewPrice(dto.TargetTo)
	if err != nil {
		return nil, fmt.Errorf("invalid target_to: %w", err)
	}

	return stock.NewStock(
		dto.Ticker,
		dto.Company,
		dto.Brokerage,
		dto.Action,
		ratingFrom,
		ratingTo,
		targetFrom,
		targetTo,
	)
}
