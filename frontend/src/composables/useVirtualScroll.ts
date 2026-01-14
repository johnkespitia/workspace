import { ref, computed, onMounted, onUnmounted } from "vue";

interface UseVirtualScrollOptions {
  itemHeight: number;
  containerHeight: number;
  overscan?: number; // Número de items a renderizar fuera del viewport
}

/**
 * Composable para virtual scrolling - solo renderiza items visibles
 * Útil para listas grandes con mejor rendimiento
 */
export function useVirtualScroll<T>(
  items: T[],
  options: UseVirtualScrollOptions
) {
  const { itemHeight, containerHeight, overscan = 3 } = options;

  const scrollTop = ref(0);
  const containerRef = ref<HTMLElement | null>(null);

  const totalHeight = computed(() => items.length * itemHeight);

  const visibleCount = computed(() =>
    Math.ceil(containerHeight / itemHeight)
  );

  const startIndex = computed(() => {
    const index = Math.floor(scrollTop.value / itemHeight);
    return Math.max(0, index - overscan);
  });

  const endIndex = computed(() => {
    const index = startIndex.value + visibleCount.value + overscan * 2;
    return Math.min(items.length, index);
  });

  const visibleItems = computed(() => {
    return items.slice(startIndex.value, endIndex.value);
  });

  const offsetY = computed(() => startIndex.value * itemHeight);

  const handleScroll = (event: Event) => {
    const target = event.target as HTMLElement;
    scrollTop.value = target.scrollTop;
  };

  onMounted(() => {
    if (containerRef.value) {
      containerRef.value.addEventListener("scroll", handleScroll);
    }
  });

  onUnmounted(() => {
    if (containerRef.value) {
      containerRef.value.removeEventListener("scroll", handleScroll);
    }
  });

  return {
    containerRef,
    visibleItems,
    totalHeight,
    offsetY,
    startIndex,
    endIndex,
  };
}
