import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useDemoStore = defineStore('demo', () => {
  const message = ref<string>('')
  const isLoading = ref<boolean>(false)
  const error = ref<string | null>(null)

  const hasMessage = computed(() => message.value.length > 0)

  async function fetchMessage() {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await fetch('/hello')
      if (!response.ok) {
        throw new Error('Error al obtener el mensaje')
      }
      const data = await response.json()
      message.value = data.message
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Error desconocido'
      message.value = ''
    } finally {
      isLoading.value = false
    }
  }

  function clearMessage() {
    message.value = ''
    error.value = null
  }

  return {
    message,
    isLoading,
    error,
    hasMessage,
    fetchMessage,
    clearMessage
  }
})
