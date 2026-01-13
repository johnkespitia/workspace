import { ref, watch, onUnmounted, type Ref } from 'vue';

export function useDebounce<T>(value: Ref<T> | (() => T), delay: number = 300): Ref<T> {
  const debouncedValue = ref(
    typeof value === 'function' ? (value as () => T)() : value.value
  ) as Ref<T>;
  let timeoutId: ReturnType<typeof setTimeout> | null = null;

  const updateDebouncedValue = (newValue: T) => {
    if (timeoutId) {
      clearTimeout(timeoutId);
    }
    timeoutId = setTimeout(() => {
      debouncedValue.value = newValue;
    }, delay);
  };

  if (typeof value === 'function') {
    // Si es una función, no podemos hacer watch directo
    // Esto se manejará en el componente que lo use
  } else {
    watch(
      value,
      (newValue) => {
        updateDebouncedValue(newValue);
      },
      { immediate: true }
    );
  }

  onUnmounted(() => {
    if (timeoutId) {
      clearTimeout(timeoutId);
    }
  });

  return debouncedValue;
}
