<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50">
    <div class="container mx-auto px-4 py-16">
      <!-- Header -->
      <header class="text-center mb-12">
        <h1 class="text-5xl font-bold text-gray-900 mb-4">
          Go + Vue 3
        </h1>
        <p class="text-xl text-gray-600">
          Stack tecnológico moderno para desarrollo full-stack
        </p>
      </header>

      <!-- Main Content Card -->
      <div class="max-w-2xl mx-auto">
        <div class="bg-white rounded-2xl shadow-xl p-8 mb-8">
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-2xl font-semibold text-gray-800">
              Conexión con API
            </h2>
            <button
              @click="demoStore.fetchMessage"
              :disabled="demoStore.isLoading"
              class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200 font-medium"
            >
              <span v-if="demoStore.isLoading">Cargando...</span>
              <span v-else>Obtener Mensaje</span>
            </button>
          </div>

          <!-- Message Display -->
          <div v-if="demoStore.hasMessage || demoStore.error" class="mt-6">
            <div
              v-if="demoStore.error"
              class="p-4 bg-red-50 border border-red-200 rounded-lg text-red-800"
            >
              <p class="font-medium">Error: {{ demoStore.error }}</p>
            </div>
            <div
              v-else-if="demoStore.message"
              class="p-4 bg-green-50 border border-green-200 rounded-lg"
            >
              <p class="text-gray-700">
                <span class="font-semibold text-green-800">Mensaje de la API:</span>
                <span class="ml-2 text-gray-900">{{ demoStore.message }}</span>
              </p>
            </div>
          </div>

          <!-- Empty State -->
          <div v-else class="text-center py-8 text-gray-500">
            <p>Haz clic en el botón para obtener un mensaje del backend</p>
          </div>
        </div>

        <!-- Tech Stack Info -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="bg-white rounded-xl shadow-lg p-6">
            <h3 class="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <span class="w-2 h-2 bg-blue-500 rounded-full mr-2"></span>
              Backend
            </h3>
            <ul class="space-y-2 text-gray-600">
              <li class="flex items-center">
                <span class="mr-2">✓</span>
                Golang
              </li>
            </ul>
          </div>

          <div class="bg-white rounded-xl shadow-lg p-6">
            <h3 class="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <span class="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
              Frontend
            </h3>
            <ul class="space-y-2 text-gray-600">
              <li class="flex items-center">
                <span class="mr-2">✓</span>
                Vue 3
              </li>
              <li class="flex items-center">
                <span class="mr-2">✓</span>
                TypeScript
              </li>
              <li class="flex items-center">
                <span class="mr-2">✓</span>
                Pinia
              </li>
              <li class="flex items-center">
                <span class="mr-2">✓</span>
                Tailwind CSS
              </li>
            </ul>
          </div>

          <div class="bg-white rounded-xl shadow-lg p-6 md:col-span-2">
            <h3 class="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <span class="w-2 h-2 bg-purple-500 rounded-full mr-2"></span>
              Base de Datos
            </h3>
            <ul class="space-y-2 text-gray-600">
              <li class="flex items-center">
                <span class="mr-2">✓</span>
                CockroachDB
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useDemoStore } from './stores/demo'

const demoStore = useDemoStore()

onMounted(() => {
  // Opcional: cargar mensaje automáticamente al montar
  // demoStore.fetchMessage()
})
</script>
