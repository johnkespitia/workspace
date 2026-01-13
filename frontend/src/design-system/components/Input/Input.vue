<template>
  <div class="input-wrapper">
    <label v-if="label" :for="inputId" class="block text-sm font-medium mb-1">
      {{ label }}
      <span v-if="required" class="text-red-500" aria-label="requerido">*</span>
    </label>
    <div class="relative">
      <input
        :id="inputId"
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :required="required"
        :aria-label="ariaLabel || label"
        :aria-invalid="hasError"
        :aria-describedby="hasError ? `${inputId}-error` : undefined"
        :class="inputClasses"
        @input="handleInput"
        @blur="handleBlur"
        @focus="handleFocus"
      />
      <span v-if="icon" class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400" aria-hidden="true">
        {{ icon }}
      </span>
    </div>
    <p v-if="error" :id="`${inputId}-error`" class="mt-1 text-sm text-red-600" role="alert">
      {{ error }}
    </p>
    <p v-else-if="hint" class="mt-1 text-sm text-gray-500">
      {{ hint }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';

interface Props {
  modelValue: string;
  type?: 'text' | 'email' | 'password' | 'number' | 'search';
  label?: string;
  placeholder?: string;
  disabled?: boolean;
  required?: boolean;
  error?: string;
  hint?: string;
  ariaLabel?: string;
  icon?: string;
  id?: string;
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  disabled: false,
  required: false,
});

const emit = defineEmits<{
  'update:modelValue': [value: string];
  blur: [event: FocusEvent];
  focus: [event: FocusEvent];
}>();

const inputId = computed(() => props.id || `input-${Math.random().toString(36).substr(2, 9)}`);
const hasError = computed(() => !!props.error);

const inputClasses = computed(() => {
  const base = 'w-full px-4 py-2 border rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-offset-1 disabled:bg-gray-100 disabled:cursor-not-allowed';
  
  if (hasError.value) {
    return `${base} border-red-500 focus:ring-red-500 focus:border-red-500`;
  }
  
  return `${base} border-gray-300 focus:ring-indigo-500 focus:border-indigo-500`;
});

const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement;
  emit('update:modelValue', target.value);
};

const handleBlur = (event: FocusEvent) => {
  emit('blur', event);
};

const handleFocus = (event: FocusEvent) => {
  emit('focus', event);
};
</script>

<style scoped>
.input-wrapper {
  @apply w-full;
}
</style>
