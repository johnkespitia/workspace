package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/john/go-react-test/api/internal/domain/stock"
	"golang.org/x/time/rate"
)

// Cache interface para almacenamiento en caché
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
}

// InMemoryCache implementa un cache en memoria simple
type InMemoryCache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewInMemoryCache crea un nuevo cache en memoria
func NewInMemoryCache() *InMemoryCache {
	c := &InMemoryCache{
		items: make(map[string]cacheItem),
	}
	// Limpiar items expirados cada minuto
	go c.cleanup()
	return c
}

func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	
	if time.Now().After(item.expiration) {
		delete(c.items, key)
		return nil, false
	}
	
	return item.value, true
}

func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.items[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

func (c *InMemoryCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.expiration) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// KarenAIClient es el cliente para la API externa de KarenAI
type KarenAIClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	rateLimiter *rate.Limiter
	cache      Cache
	maxRetries int
	retryDelay time.Duration
}

// APIResponse representa la respuesta de la API externa
type APIResponse struct {
	Items    []StockDTO `json:"items"`
	NextPage string     `json:"next_page,omitempty"`
}

// StockDTO representa un stock en la respuesta de la API
type StockDTO struct {
	Ticker     string `json:"ticker"`
	Company    string `json:"company"`
	Brokerage  string `json:"brokerage"`
	Action     string `json:"action"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	TargetFrom string `json:"target_from"` // Viene como string con $, ej: "$3.00"
	TargetTo   string `json:"target_to"`   // Viene como string con $, ej: "$2.50"
	Time       string `json:"time,omitempty"`
}

// NewKarenAIClient crea un nuevo cliente para la API externa
func NewKarenAIClient(baseURL, apiKey string) *KarenAIClient {
	return &KarenAIClient{
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // Aumentado para operaciones largas
		},
		baseURL:     baseURL,
		apiKey:      apiKey,
		rateLimiter: rate.NewLimiter(rate.Limit(10), 1), // 10 requests por segundo
		cache:       NewInMemoryCache(),
		maxRetries:  3,
		retryDelay:  1 * time.Second,
	}
}

// NewKarenAIClientWithOptions crea un cliente con opciones personalizadas
func NewKarenAIClientWithOptions(baseURL, apiKey string, requestsPerSecond float64, maxRetries int, cache Cache) *KarenAIClient {
	if cache == nil {
		cache = NewInMemoryCache()
	}
	return &KarenAIClient{
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		baseURL:     baseURL,
		apiKey:      apiKey,
		rateLimiter: rate.NewLimiter(rate.Limit(requestsPerSecond), 1),
		cache:       cache,
		maxRetries:  maxRetries,
		retryDelay:  1 * time.Second,
	}
}

// FetchStocks obtiene stocks de la API externa con paginación
// Incluye rate limiting, retry logic y caching
func (c *KarenAIClient) FetchStocks(ctx context.Context, nextPage string) (*APIResponse, error) {
	// Verificar cache primero
	cacheKey := fmt.Sprintf("stocks:%s", nextPage)
	if cached, ok := c.cache.Get(cacheKey); ok {
		return cached.(*APIResponse), nil
	}

	// Usar retry logic
	response, err := c.fetchWithRetry(ctx, nextPage)
	if err != nil {
		return nil, err
	}

	// Guardar en cache (TTL de 5 minutos)
	c.cache.Set(cacheKey, response, 5*time.Minute)

	return response, nil
}

// fetchWithRetry implementa retry logic con exponential backoff
func (c *KarenAIClient) fetchWithRetry(ctx context.Context, nextPage string) (*APIResponse, error) {
	var lastErr error
	
	for attempt := 0; attempt < c.maxRetries; attempt++ {
		// Rate limiting: esperar antes de hacer la request
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter error: %w", err)
		}

		// Intentar hacer la request
		response, err := c.fetchPage(ctx, nextPage)
		if err == nil {
			return response, nil
		}

		lastErr = err

		// Si no es el último intento, esperar antes de reintentar
		if attempt < c.maxRetries-1 {
			// Exponential backoff: 1s, 2s, 4s...
			backoff := time.Duration(1<<uint(attempt)) * c.retryDelay
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
				// Continuar con el siguiente intento
			}
		}
	}

	return nil, fmt.Errorf("failed after %d retries: %w", c.maxRetries, lastErr)
}

// fetchPage hace una request HTTP a la API
func (c *KarenAIClient) fetchPage(ctx context.Context, nextPage string) (*APIResponse, error) {
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
// Incluye caching para evitar requests innecesarias
func (c *KarenAIClient) FetchAllStocks(ctx context.Context) ([]*stock.Stock, error) {
	// Verificar cache para todos los stocks
	cacheKey := "stocks:all"
	if cached, ok := c.cache.Get(cacheKey); ok {
		return cached.([]*stock.Stock), nil
	}

	var allStocks []*stock.Stock
	nextPage := ""
	pageCount := 0

	for {
		response, err := c.FetchStocks(ctx, nextPage)
		if err != nil {
			return nil, fmt.Errorf("error fetching page %d: %w", pageCount, err)
		}

		// Convertir DTOs a entidades de dominio
		convertedCount := 0
		for _, dto := range response.Items {
			s, err := c.convertToDomainEntity(dto)
			if err != nil {
				// Log error pero continuar con otros stocks
				// Podríamos agregar logging aquí si fuera necesario
				continue
			}
			allStocks = append(allStocks, s)
			convertedCount++
		}

		// Si no se convirtió ningún stock en esta página, puede ser un problema
		if len(response.Items) > 0 && convertedCount == 0 {
			return nil, fmt.Errorf("failed to convert any stocks from page %d (received %d stocks)", pageCount, len(response.Items))
		}

		// Verificar si hay más páginas
		if response.NextPage == "" {
			break
		}
		nextPage = response.NextPage
		pageCount++
	}

	if len(allStocks) == 0 {
		return nil, fmt.Errorf("no stocks found after fetching all pages")
	}

	// Guardar en cache (TTL de 10 minutos para todos los stocks)
	c.cache.Set(cacheKey, allStocks, 10*time.Minute)

	return allStocks, nil
}

// convertToDomainEntity convierte un DTO a una entidad de dominio
func (c *KarenAIClient) convertToDomainEntity(dto StockDTO) (*stock.Stock, error) {
	ratingFrom := stock.Rating(dto.RatingFrom)
	ratingTo := stock.Rating(dto.RatingTo)

	// Parsear precios que vienen como strings con $, ej: "$3.00"
	targetFromFloat, err := parsePriceString(dto.TargetFrom)
	if err != nil {
		return nil, fmt.Errorf("invalid target_from '%s': %w", dto.TargetFrom, err)
	}

	targetToFloat, err := parsePriceString(dto.TargetTo)
	if err != nil {
		return nil, fmt.Errorf("invalid target_to '%s': %w", dto.TargetTo, err)
	}

	targetFrom, err := stock.NewPrice(targetFromFloat)
	if err != nil {
		return nil, fmt.Errorf("invalid target_from: %w", err)
	}

	targetTo, err := stock.NewPrice(targetToFloat)
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

// parsePriceString parsea un string de precio como "$3.00" a float64
func parsePriceString(priceStr string) (float64, error) {
	// Remover el símbolo $ y espacios
	cleaned := strings.TrimSpace(priceStr)
	cleaned = strings.TrimPrefix(cleaned, "$")
	cleaned = strings.TrimSpace(cleaned)

	// Convertir a float64
	price, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse price: %w", err)
	}

	return price, nil
}
