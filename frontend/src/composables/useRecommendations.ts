import { ref, computed } from 'vue';
import { useApi } from './useApi';
import type { Recommendation } from '@/utils/api';

export function useRecommendations() {
  const api = useApi();
  const recommendations = ref<Recommendation[]>([]);
  const limit = ref(10);

  const loadRecommendations = async (newLimit?: number) => {
    if (newLimit) limit.value = newLimit;
    recommendations.value = await api.fetchRecommendations(limit.value);
  };

  return {
    recommendations: computed(() => recommendations.value),
    loading: api.loading,
    error: api.error,
    limit: computed(() => limit.value),
    loadRecommendations,
  };
}
