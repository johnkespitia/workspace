<template>
  <div class="table-container">
    <div class="overflow-x-auto">
      <table :class="tableClasses" role="table" :aria-label="ariaLabel">
        <thead v-if="headers && headers.length > 0">
          <tr>
            <th
              v-for="(header, index) in headers"
              :key="index"
              :class="headerClasses"
              :scope="'col'"
              :aria-sort="getSortDirection(header.key)"
            >
              <div class="flex items-center gap-2">
                <span>{{ header.label }}</span>
                <button
                  v-if="header.sortable"
                  @click="handleSort(header.key)"
                  :class="sortButtonClasses"
                  :aria-label="`Ordenar por ${header.label}`"
                >
                  <span v-if="sortKey === header.key">
                    {{ sortDirection === 'asc' ? '↑' : '↓' }}
                  </span>
                  <span v-else>⇅</span>
                </button>
              </div>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(row, rowIndex) in data"
            :key="rowIndex"
            :class="rowClasses(rowIndex)"
            :tabindex="clickable ? 0 : -1"
            :aria-selected="clickable && focusedRowIndex === rowIndex"
            :role="clickable ? 'button' : undefined"
            @click="handleRowClick(row)"
            @keydown="handleRowKeyDown($event, row, rowIndex)"
          >
            <td
              v-for="(header, colIndex) in headers"
              :key="colIndex"
              :class="cellClasses"
            >
              <slot
                :name="`cell-${header.key}`"
                :row="row"
                :value="row[header.key]"
                :header="header"
              >
                {{ row[header.key] }}
              </slot>
            </td>
          </tr>
          <tr v-if="data.length === 0">
            <td :colspan="headers.length" class="text-center py-8 text-gray-500">
              <slot name="empty">
                No hay datos disponibles
              </slot>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue';
import { handleKeyboardNavigation } from '@/utils/accessibility';

export interface TableHeader {
  key: string;
  label: string;
  sortable?: boolean;
}

interface Props {
  headers: TableHeader[];
  data: Record<string, any>[];
  sortable?: boolean;
  striped?: boolean;
  hoverable?: boolean;
  clickable?: boolean;
  ariaLabel?: string;
}

const props = withDefaults(defineProps<Props>(), {
  sortable: true,
  striped: true,
  hoverable: true,
  clickable: false,
});

const emit = defineEmits<{
  sort: [key: string, direction: 'asc' | 'desc'];
  rowClick: [row: Record<string, any>];
}>();

const sortKey = ref<string | null>(null);
const sortDirection = ref<'asc' | 'desc'>('asc');
const focusedRowIndex = ref<number | null>(null);

const tableClasses = computed(() => {
  return 'min-w-full divide-y divide-gray-200 border border-gray-200 rounded-lg overflow-hidden';
});

const headerClasses = computed(() => {
  return 'px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-700 uppercase tracking-wider';
});

const cellClasses = computed(() => {
  return 'px-6 py-4 whitespace-nowrap text-sm text-gray-900';
});

const sortButtonClasses = computed(() => {
  return 'inline-flex items-center text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-indigo-500 rounded p-1';
});

const rowClasses = (index: number) => {
  const base = props.hoverable ? 'hover:bg-gray-50 transition-colors duration-150' : '';
  const striped = props.striped && index % 2 === 1 ? 'bg-gray-50' : 'bg-white';
  const clickable = props.clickable ? 'cursor-pointer' : '';
  return `${base} ${striped} ${clickable}`;
};

const getSortDirection = (key: string): string | undefined => {
  if (sortKey.value !== key) return undefined;
  return sortDirection.value === 'asc' ? 'ascending' : 'descending';
};

const handleSort = (key: string) => {
  if (!props.sortable) return;
  
  if (sortKey.value === key) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc';
  } else {
    sortKey.value = key;
    sortDirection.value = 'asc';
  }
  
  emit('sort', key, sortDirection.value);
};

const handleRowClick = (row: Record<string, any>) => {
  if (props.clickable) {
    emit('rowClick', row);
  }
};

const handleRowKeyDown = (event: KeyboardEvent, row: Record<string, any>, index: number) => {
  if (!props.clickable) return;

  const newIndex = handleKeyboardNavigation(
    event,
    props.data,
    focusedRowIndex.value ?? index,
    (idx) => {
      focusedRowIndex.value = idx;
      // Focus en la fila
      const rowElement = event.currentTarget as HTMLElement;
      rowElement.focus();
    }
  );

  if (event.key === 'Enter' || event.key === ' ') {
    handleRowClick(row);
  }

  if (newIndex !== null) {
    focusedRowIndex.value = newIndex;
  }
};

// Resetear focusedRowIndex cuando cambian los datos
onMounted(() => {
  focusedRowIndex.value = null;
});
</script>

<style scoped>
.table-container {
  @apply w-full;
}
</style>
