<template>
  <div class="stock-detail-container p-6">
    <div class="mb-6">
      <Button @click="goBack" variant="ghost" size="sm">
        ← Volver
      </Button>
    </div>

    <Card v-if="stock" variant="elevated" padding="lg">
      <template #header>
        <div class="flex items-center justify-between">
          <h1 class="text-3xl font-bold text-gray-900">{{ stock.ticker }}</h1>
          <span :class="getRatingClass(stock.ratingTo)" class="px-3 py-1 rounded-lg text-sm font-semibold">
            {{ stock.ratingTo || 'N/A' }}
          </span>
        </div>
      </template>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <h3 class="text-sm font-medium text-gray-500 mb-2">Compañía</h3>
          <p class="text-lg text-gray-900">{{ stock.companyName }}</p>
        </div>

        <div>
          <h3 class="text-sm font-medium text-gray-500 mb-2">Brokerage</h3>
          <p class="text-lg text-gray-900">{{ stock.brokerage || 'N/A' }}</p>
        </div>

        <div>
          <h3 class="text-sm font-medium text-gray-500 mb-2">Rating Anterior</h3>
          <p class="text-lg text-gray-900">{{ stock.ratingFrom || 'N/A' }}</p>
        </div>

        <div>
          <h3 class="text-sm font-medium text-gray-500 mb-2">Rating Actual</h3>
          <p class="text-lg text-gray-900">{{ stock.ratingTo || 'N/A' }}</p>
        </div>

        <div>
          <h3 class="text-sm font-medium text-gray-500 mb-2">Target Anterior</h3>
          <p class="text-lg font-semibold text-gray-900">
            ${{ stock.targetFrom?.toFixed(2) || 'N/A' }}
          </p>
        </div>

        <div>
          <h3 class="text-sm font-medium text-gray-500 mb-2">Target Actual</h3>
          <p class="text-lg font-semibold text-indigo-600">
            ${{ stock.targetTo?.toFixed(2) || 'N/A' }}
          </p>
        </div>

        <div>
          <h3 class="text-sm font-medium text-gray-500 mb-2">Cambio</h3>
          <p :class="getChangeClass(calculateChange())" class="text-lg font-semibold">
            {{ calculateChange() > 0 ? '+' : '' }}{{ calculateChange().toFixed(2) }}%
          </p>
        </div>

        <div>
          <h3 class="text-sm font-medium text-gray-500 mb-2">Acción</h3>
          <p class="text-lg text-gray-900">{{ stock.action || 'N/A' }}</p>
        </div>
      </div>

      <div class="mt-6 pt-6 border-t border-gray-200">
        <div class="grid grid-cols-2 gap-4 text-sm text-gray-500">
          <div>
            <span class="font-medium">Creado:</span>
            {{ formatDate(stock.createdAt) }}
          </div>
          <div>
            <span class="font-medium">Actualizado:</span>
            {{ formatDate(stock.updatedAt) }}
          </div>
        </div>
      </div>
    </Card>

    <!-- Gráfico de evolución del precio -->
    <Card v-if="stock && stock.targetFrom && stock.targetTo" variant="elevated" padding="lg" class="mt-6">
      <PriceChart :target-from="stock.targetFrom" :target-to="stock.targetTo"
        :aria-label="`Gráfico de evolución del precio objetivo para ${stock.ticker}`" />
    </Card>

    <Card v-else variant="elevated" padding="lg">
      <div class="text-center py-8">
        <p class="text-gray-500">Cargando información de la acción...</p>
      </div>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useStock } from "@/composables/useStock";
import Card from "@/design-system/components/Card/Card.vue";
import Button from "@/design-system/components/Button/Button.vue";
import PriceChart from "@/components/PriceChart.vue";
import type { Stock } from "@/utils/api";

const route = useRoute();
const router = useRouter();
const stockComposable = useStock();

const stock = ref<Stock | null>(null);

onMounted(async () => {
  const ticker = route.params.ticker as string;
  if (ticker) {
    stock.value = await stockComposable.loadStock(ticker);
  }
});

const goBack = () => {
  router.push('/stocks');
};

const calculateChange = (): number => {
  if (!stock.value?.targetFrom || !stock.value?.targetTo) return 0;
  return ((stock.value.targetTo - stock.value.targetFrom) / stock.value.targetFrom) * 100;
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

const getChangeClass = (change: number) => {
  if (change > 0) return 'text-green-600';
  if (change < 0) return 'text-red-600';
  return 'text-gray-600';
};

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('es-ES', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
};
</script>

<style scoped>
.stock-detail-container {
  max-width: 1200px;
  margin: 0 auto;
}
</style>
