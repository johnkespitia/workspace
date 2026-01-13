import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { Stock, StockFilter, StockSort } from '@/utils/api';

export const useStockStore = defineStore('stock', () => {
  const stocks = ref<Stock[]>([]);
  const currentStock = ref<Stock | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const totalCount = ref(0);
  const currentPage = ref(1);
  const pageSize = ref(50);
  const filter = ref<StockFilter>({});
  const sort = ref<StockSort>({ field: 'ticker', direction: 'asc' });

  const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value));

  const setStocks = (newStocks: Stock[]) => {
    stocks.value = newStocks;
  };

  const setCurrentStock = (stock: Stock | null) => {
    currentStock.value = stock;
  };

  const setLoading = (isLoading: boolean) => {
    loading.value = isLoading;
  };

  const setError = (err: string | null) => {
    error.value = err;
  };

  const setTotalCount = (count: number) => {
    totalCount.value = count;
  };

  const setCurrentPage = (page: number) => {
    currentPage.value = page;
  };

  const setFilter = (newFilter: StockFilter) => {
    filter.value = newFilter;
  };

  const setSort = (newSort: StockSort) => {
    sort.value = newSort;
  };

  return {
    stocks: computed(() => stocks.value),
    currentStock: computed(() => currentStock.value),
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    totalCount: computed(() => totalCount.value),
    currentPage: computed(() => currentPage.value),
    totalPages,
    pageSize: computed(() => pageSize.value),
    filter: computed(() => filter.value),
    sort: computed(() => sort.value),
    setStocks,
    setCurrentStock,
    setLoading,
    setError,
    setTotalCount,
    setCurrentPage,
    setFilter,
    setSort,
  };
});
