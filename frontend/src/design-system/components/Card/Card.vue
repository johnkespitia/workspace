<template>
  <div :class="cardClasses" :role="role">
    <div v-if="$slots.header || title" class="card-header">
      <slot name="header">
        <h3 v-if="title" class="text-lg font-semibold">{{ title }}</h3>
      </slot>
    </div>
    <div class="card-body">
      <slot></slot>
    </div>
    <div v-if="$slots.footer" class="card-footer">
      <slot name="footer"></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  title?: string;
  variant?: 'default' | 'elevated' | 'outlined';
  padding?: 'none' | 'sm' | 'md' | 'lg';
  role?: string;
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  padding: 'md',
  role: 'article',
});

const cardClasses = computed(() => {
  const base = 'rounded-lg border transition-shadow duration-200';
  
  const variants = {
    default: 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700',
    elevated: 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 shadow-lg',
    outlined: 'bg-transparent dark:bg-transparent border-gray-300 dark:border-gray-600',
  };
  
  const paddings = {
    none: '',
    sm: 'p-4',
    md: 'p-6',
    lg: 'p-8',
  };
  
  return `${base} ${variants[props.variant]} ${paddings[props.padding]}`;
});
</script>

<style scoped>
.card-header {
  @apply border-b border-gray-200 dark:border-gray-700 pb-4 mb-4;
}

.card-body {
  @apply w-full;
}

.card-footer {
  @apply border-t border-gray-200 dark:border-gray-700 pt-4 mt-4;
}
</style>
