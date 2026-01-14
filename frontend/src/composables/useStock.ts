import { ref, computed } from "vue";
import { useApi } from "./useApi";
import type { Stock, StockFilter, StockSort } from "@/utils/api";

export function useStock() {
  const api = useApi();
  const stocks = ref<Stock[]>([]);
  const currentStock = ref<Stock | null>(null);
  const totalCount = ref(0);
  const currentPage = ref(1);
  const pageSize = ref(50);
  const filter = ref<StockFilter>({});
  const sort = ref<StockSort>({ field: "TICKER", direction: "ASC" });

  const totalPages = computed(() =>
    Math.ceil(totalCount.value / pageSize.value)
  );

  const loadStocks = async (
    newFilter?: StockFilter,
    newSort?: StockSort,
    page: number = 1
  ) => {
    try {
      if (newFilter) filter.value = newFilter;
      if (newSort) sort.value = newSort;
      currentPage.value = page;

      const offset = (page - 1) * pageSize.value;
      const result = await api.fetchStocks(
        filter.value,
        sort.value,
        pageSize.value,
        offset
      );

      stocks.value = result.stocks;
      totalCount.value = result.totalCount;
    } catch (error) {
      console.error("Error loading stocks:", error);
      // No fallar completamente, solo mostrar lista vacía
      stocks.value = [];
      totalCount.value = 0;
    }
  };

  const loadStock = async (ticker: string) => {
    try {
      currentStock.value = await api.fetchStock(ticker);
      return currentStock.value;
    } catch (error) {
      console.error("Error loading stock:", error);
      currentStock.value = null;
      return null;
    }
  };

  const searchStocks = async (query: string) => {
    // Para búsqueda, usar companyName (búsqueda parcial) en lugar de ticker (búsqueda exacta)
    // Si el query parece un ticker (corto y en mayúsculas), usar ticker, sino companyName
    const isLikelyTicker = query.length <= 10 && query === query.toUpperCase();
    const searchFilter: StockFilter = {
      ...filter.value,
      ...(isLikelyTicker 
        ? { ticker: query || undefined }
        : { companyName: query || undefined }
      ),
    };
    await loadStocks(searchFilter, sort.value, 1);
  };

  const setFilter = async (newFilter: StockFilter) => {
    await loadStocks(newFilter, sort.value, 1);
  };

  const setSort = async (newSort: StockSort) => {
    await loadStocks(filter.value, newSort, 1);
  };

  const goToPage = async (page: number) => {
    await loadStocks(filter.value, sort.value, page);
  };

  return {
    stocks: computed(() => stocks.value),
    currentStock: computed(() => currentStock.value),
    loading: api.loading,
    error: api.error,
    totalCount: computed(() => totalCount.value),
    currentPage: computed(() => currentPage.value),
    totalPages,
    pageSize: computed(() => pageSize.value),
    filter: computed(() => filter.value),
    sort: computed(() => sort.value),
    loadStocks,
    loadStock,
    searchStocks,
    setFilter,
    setSort,
    goToPage,
  };
}
