<template>
  <div class="stock-list-container p-6">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900">Lista de Acciones</h1>
      <Button @click="handleSync" :loading="syncing" variant="primary">
        Sincronizar Stocks
      </Button>
    </div>

    <Card>
      <template #header>
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-semibold">Acciones</h2>
          <div class="flex items-center gap-4">
            <select
              v-model="selectedRating"
              @change="handleFilterChange"
              class="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
            >
              <option value="">Todos los ratings</option>
              <option value="Buy">Buy</option>
              <option value="Strong Buy">Strong Buy</option>
              <option value="Speculative Buy">Speculative Buy</option>
              <option value="Market Perform">Market Perform</option>
            </select>
          </div>
        </div>
      </template>

      <Table
        :headers="tableHeaders"
        :data="displayedStocks"
        :sortable="true"
        :hoverable="true"
        :clickable="true"
        @sort="handleSort"
        @row-click="handleRowClick"
        aria-label="Tabla de acciones"
      >
        <template #cell-ratingTo="{ value }">
          <span
            :class="getRatingClass(value)"
            class="px-2 py-1 rounded text-xs font-semibold"
          >
            {{ value || 'N/A' }}
          </span>
        </template>
        <template #cell-targetTo="{ value }">
          <span class="font-medium">${{ value?.toFixed(2) || 'N/A' }}</span>
        </template>
        <template #cell-action="{ value }">
          <span class="text-sm text-gray-600">{{ value || 'N/A' }}</span>
        </template>
        <template #empty>
          <div class="text-center py-8">
            <p class="text-gray-500">No hay acciones disponibles</p>
          </div>
        </template>
      </Table>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useStock } from '@/composables/useStock';
import { useApi } from '@/composables/useApi';
import { withLoading } from '@/hoc/withLoading';
import { withError } from '@/hoc/withError';
import { withPagination } from '@/hoc/withPagination';
import { withSearch } from '@/hoc/withSearch';
import Table from '@/design-system/components/Table/Table.vue';
import Card from '@/design-system/components/Card/Card.vue';
import Button from '@/design-system/components/Button/Button.vue';
import type { StockSort, StockSortField, SortDirection } from '@/utils/api';

const router = useRouter();
const stockComposable = useStock();
const api = useApi();

const selectedRating = ref('');
const searchQuery = ref('');

const tableHeaders = [
  { key: 'ticker', label: 'Ticker', sortable: true },
  { key: 'companyName', label: 'Compañía', sortable: true },
  { key: 'brokerage', label: 'Brokerage', sortable: false },
  { key: 'ratingTo', label: 'Rating', sortable: true },
  { key: 'targetTo', label: 'Target', sortable: true },
  { key: 'action', label: 'Acción', sortable: false },
];

const displayedStocks = computed(() => {
  let stocks = stockComposable.stocks.value;

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    stocks = stocks.filter(
      (s) =>
        s.ticker.toLowerCase().includes(query) ||
        s.companyName.toLowerCase().includes(query)
    );
  }

  return stocks;
});

const syncing = ref(false);

onMounted(async () => {
  await stockComposable.loadStocks();
});

const handleSync = async () => {
  syncing.value = true;
  try {
    await api.syncStocks();
    await stockComposable.loadStocks();
  } catch (error) {
    console.error('Error sincronizando stocks:', error);
  } finally {
    syncing.value = false;
  }
};

// Mapeo de nombres de campos de la tabla a valores del enum GraphQL
const mapFieldToEnum = (field: string): StockSortField => {
  const fieldMap: Record<string, StockSortField> = {
    ticker: "TICKER",
    companyName: "COMPANY_NAME",
    ratingTo: "RATING_TO",
    targetTo: "TARGET_TO",
    createdAt: "CREATED_AT",
  };
  return fieldMap[field] || "TICKER";
};

const mapDirectionToEnum = (direction: 'asc' | 'desc'): SortDirection => {
  return direction === 'asc' ? 'ASC' : 'DESC';
};

const handleSort = (key: string, direction: 'asc' | 'desc') => {
  const sort: StockSort = {
    field: mapFieldToEnum(key),
    direction: mapDirectionToEnum(direction),
  };
  stockComposable.setSort(sort);
  stockComposable.loadStocks(stockComposable.filter.value, sort, 1);
};

const handleFilterChange = () => {
  const filter = {
    ...stockComposable.filter.value,
    ratings: selectedRating.value ? [selectedRating.value] : undefined, // ✅ Usar ratings (plural)
  };
  stockComposable.setFilter(filter);
  stockComposable.loadStocks(filter, stockComposable.sort.value, 1);
};

const handleSearch = (query: string) => {
  searchQuery.value = query;
};

const handleRowClick = (row: any) => {
  router.push(`/stocks/${row.ticker}`);
};

const getRatingClass = (rating: string | null | undefined) => {
  if (!rating) return 'bg-gray-100 text-gray-700';
  
  const ratingLower = rating.toLowerCase();
  if (ratingLower.includes('strong buy')) {
    return 'bg-green-100 text-green-800';
  }
  if (ratingLower.includes('buy')) {
    return 'bg-blue-100 text-blue-800';
  }
  if (ratingLower.includes('speculative')) {
    return 'bg-yellow-100 text-yellow-800';
  }
  return 'bg-gray-100 text-gray-700';
};
</script>

<style scoped>
.stock-list-container {
  max-width: 1400px;
  margin: 0 auto;
}
</style>
