import { ref, computed } from 'vue';
import { graphqlClient, GET_STOCKS_QUERY, GET_STOCK_QUERY, GET_RECOMMENDATIONS_QUERY, SYNC_STOCKS_MUTATION, type Stock, type StockFilter, type StockSort, type StockConnection, type Recommendation } from '@/utils/api';

interface CacheEntry<T> {
  data: T;
  timestamp: number;
  ttl: number;
}

class ApiCache {
  private cache = new Map<string, CacheEntry<any>>();

  get<T>(key: string): T | null {
    const entry = this.cache.get(key);
    if (!entry) return null;

    if (Date.now() - entry.timestamp > entry.ttl) {
      this.cache.delete(key);
      return null;
    }

    return entry.data as T;
  }

  set<T>(key: string, data: T, ttl: number = 60000): void {
    this.cache.set(key, {
      data,
      timestamp: Date.now(),
      ttl,
    });
  }

  clear(): void {
    this.cache.clear();
  }
}

const cache = new ApiCache();
const pendingRequests = new Map<string, Promise<any>>();

export function useApi() {
  const loading = ref(false);
  const error = ref<string | null>(null);

  const fetchStocks = async (
    filter?: StockFilter,
    sort?: StockSort,
    limit: number = 50,
    offset: number = 0
  ): Promise<StockConnection> => {
    const cacheKey = JSON.stringify({ filter, sort, limit, offset });
    
    // Verificar cache
    const cached = cache.get<StockConnection>(cacheKey);
    if (cached) {
      return cached;
    }

    // Request deduplication
    if (pendingRequests.has(cacheKey)) {
      return pendingRequests.get(cacheKey)!;
    }

    loading.value = true;
    error.value = null;

    try {
      const promise = graphqlClient
        .query(GET_STOCKS_QUERY, { filter, sort, limit, offset })
        .toPromise()
        .then((result) => {
          if (result.error) {
            throw new Error(result.error.message);
          }
          const data = result.data?.stocks as StockConnection;
          cache.set(cacheKey, data, 60000); // Cache por 1 minuto
          return data;
        });

      pendingRequests.set(cacheKey, promise);
      const data = await promise;
      pendingRequests.delete(cacheKey);
      return data;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Error desconocido';
      pendingRequests.delete(cacheKey);
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const fetchStock = async (ticker: string): Promise<Stock | null> => {
    const cacheKey = `stock:${ticker}`;
    
    const cached = cache.get<Stock>(cacheKey);
    if (cached) {
      return cached;
    }

    loading.value = true;
    error.value = null;

    try {
      const result = await graphqlClient
        .query(GET_STOCK_QUERY, { ticker })
        .toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      const stock = result.data?.stock as Stock;
      if (stock) {
        cache.set(cacheKey, stock, 60000);
      }
      return stock || null;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Error desconocido';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const fetchRecommendations = async (limit: number = 10): Promise<Recommendation[]> => {
    const cacheKey = `recommendations:${limit}`;
    
    const cached = cache.get<Recommendation[]>(cacheKey);
    if (cached) {
      return cached;
    }

    loading.value = true;
    error.value = null;

    try {
      const result = await graphqlClient
        .query(GET_RECOMMENDATIONS_QUERY, { limit })
        .toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      const recommendations = (result.data?.recommendations || []) as Recommendation[];
      cache.set(cacheKey, recommendations, 300000); // Cache por 5 minutos
      return recommendations;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Error desconocido';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const syncStocks = async (): Promise<{ success: boolean; message: string; stocksSynced: number }> => {
    loading.value = true;
    error.value = null;

    try {
      const result = await graphqlClient
        .mutation(SYNC_STOCKS_MUTATION)
        .toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      // Limpiar cache después de sincronizar
      cache.clear();

      return result.data?.syncStocks || { success: false, message: 'Error desconocido', stocksSynced: 0 };
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Error desconocido';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return {
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    fetchStocks,
    fetchStock,
    fetchRecommendations,
    syncStocks,
    clearCache: () => cache.clear(),
  };
}
