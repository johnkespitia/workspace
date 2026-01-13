<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-200">
    <!-- Navigation -->
    <nav class="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center">
            <router-link to="/" class="text-xl font-bold text-indigo-600 dark:text-indigo-400">
              Stock Info
            </router-link>
            <div class="ml-10 flex items-baseline space-x-4">
              <router-link
                to="/stocks"
                class="px-3 py-2 rounded-md text-sm font-medium transition-colors duration-200"
                :class="
                  $route.name === 'StockList'
                    ? 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900 dark:text-indigo-300'
                    : 'text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700'
                "
              >
                Acciones
              </router-link>
              <router-link
                to="/recommendations"
                class="px-3 py-2 rounded-md text-sm font-medium transition-colors duration-200"
                :class="
                  $route.name === 'Recommendations'
                    ? 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900 dark:text-indigo-300'
                    : 'text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700'
                "
              >
                Recomendaciones
              </router-link>
            </div>
          </div>
          <ThemeToggle />
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<script setup lang="ts">
import ThemeToggle from '@/design-system/components/ThemeToggle/ThemeToggle.vue';
</script>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
