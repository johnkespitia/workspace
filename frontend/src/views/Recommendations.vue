<template>
  <div class="recommendations-container p-6">
    <div class="mb-6">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">Recomendaciones de Inversión</h1>
      <p class="text-gray-600">
        Las mejores acciones para invertir basadas en nuestro algoritmo de recomendación
      </p>
    </div>

    <Card variant="elevated" padding="lg">
      <div class="mb-4 flex items-center justify-between">
        <h2 class="text-xl font-semibold">Top Recomendaciones</h2>
        <select
          v-model="selectedLimit"
          @change="handleLimitChange"
          class="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
        >
          <option :value="10">Top 10</option>
          <option :value="20">Top 20</option>
          <option :value="50">Top 50</option>
        </select>
      </div>

      <div v-if="recommendations.length > 0" class="space-y-4">
        <div
          v-for="(rec, index) in recommendations"
          :key="rec.stock.id"
          class="recommendation-item p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow duration-200 cursor-pointer"
          @click="goToStock(rec.stock.ticker)"
        >
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <div class="flex items-center gap-3 mb-2">
                <span class="text-2xl font-bold text-indigo-600">#{{ index + 1 }}</span>
                <div>
                  <h3 class="text-lg font-semibold text-gray-900">
                    {{ rec.stock.ticker }}
                  </h3>
                  <p class="text-sm text-gray-600">{{ rec.stock.companyName }}</p>
                </div>
              </div>

              <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mt-4">
                <div>
                  <p class="text-xs text-gray-500 mb-1">Score Total</p>
                  <p class="text-lg font-bold text-indigo-600">
                    {{ rec.score.toFixed(2) }}
                  </p>
                </div>
                <div>
                  <p class="text-xs text-gray-500 mb-1">Cambio %</p>
                  <p
                    :class="rec.priceChange > 0 ? 'text-green-600' : 'text-red-600'"
                    class="text-lg font-semibold"
                  >
                    {{ rec.priceChange > 0 ? '+' : '' }}{{ rec.priceChange.toFixed(2) }}%
                  </p>
                </div>
                <div>
                  <p class="text-xs text-gray-500 mb-1">Rating Score</p>
                  <p class="text-lg font-semibold text-gray-900">
                    {{ rec.ratingScore.toFixed(2) }}
                  </p>
                </div>
                <div>
                  <p class="text-xs text-gray-500 mb-1">Target</p>
                  <p class="text-lg font-semibold text-gray-900">
                    ${{ rec.stock.targetTo?.toFixed(2) || 'N/A' }}
                  </p>
                </div>
              </div>

              <div class="mt-3">
                <span
                  :class="getRatingClass(rec.stock.ratingTo)"
                  class="px-2 py-1 rounded text-xs font-semibold"
                >
                  {{ rec.stock.ratingTo || 'N/A' }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="!loading" class="text-center py-8">
        <p class="text-gray-500">No hay recomendaciones disponibles</p>
      </div>
    </Card>

    <Card v-if="recommendations.length > 0" variant="outlined" padding="md" class="mt-6">
      <h3 class="text-lg font-semibold mb-3">Sobre el Algoritmo</h3>
      <p class="text-gray-600 text-sm">
        El algoritmo de recomendación considera múltiples factores:
      </p>
      <ul class="mt-2 space-y-1 text-sm text-gray-600 list-disc list-inside">
        <li>Cambio porcentual en el precio objetivo (50% del score)</li>
        <li>Rating de la acción (30% del score)</li>
        <li>Acción realizada (20% del score)</li>
      </ul>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useRecommendations } from '@/composables/useRecommendations';
import Card from '@/design-system/components/Card/Card.vue';

const router = useRouter();
const recommendationsComposable = useRecommendations();

const recommendations = recommendationsComposable.recommendations;
const loading = recommendationsComposable.loading;
const selectedLimit = ref(10);

onMounted(async () => {
  await recommendationsComposable.loadRecommendations(selectedLimit.value);
});

const handleLimitChange = async () => {
  await recommendationsComposable.loadRecommendations(selectedLimit.value);
};

const goToStock = (ticker: string) => {
  router.push(`/stocks/${ticker}`);
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
.recommendations-container {
  max-width: 1200px;
  margin: 0 auto;
}

.recommendation-item {
  @apply bg-white;
}
</style>
