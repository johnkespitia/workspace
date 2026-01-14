<template>
    <div class="price-chart-container">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Evolución del Precio Objetivo</h3>
        <div class="chart-wrapper" :style="{ height: height + 'px' }">
            <Line v-if="chartData" :data="chartData" :options="chartOptions" :aria-label="ariaLabel" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { Line } from "vue-chartjs";
// Importar configuración de Chart.js (registra componentes automáticamente)
import "@/utils/chartConfig";

interface Props {
    targetFrom?: number;
    targetTo?: number;
    height?: number;
    ariaLabel?: string;
}

const props = withDefaults(defineProps<Props>(), {
    height: 200,
    ariaLabel: "Gráfico de evolución del precio objetivo",
});

const chartData = computed(() => {
    if (!props.targetFrom || !props.targetTo) return null;

    const change = ((props.targetTo - props.targetFrom) / props.targetFrom) * 100;
    const isPositive = change >= 0;

    return {
        labels: ["Target Anterior", "Target Actual"],
        datasets: [
            {
                label: "Precio Objetivo ($)",
                data: [props.targetFrom, props.targetTo],
                borderColor: isPositive ? "#22c55e" : "#ef4444",
                backgroundColor: isPositive
                    ? "rgba(34, 197, 94, 0.1)"
                    : "rgba(239, 68, 68, 0.1)",
                borderWidth: 3,
                fill: true,
                tension: 0.4,
                pointRadius: 6,
                pointHoverRadius: 8,
                pointBackgroundColor: isPositive ? "#22c55e" : "#ef4444",
                pointBorderColor: "#ffffff",
                pointBorderWidth: 2,
            },
        ],
    };
});

const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: {
            display: true,
            position: "top" as const,
        },
        tooltip: {
            callbacks: {
                label: (context: any) => {
                    return `$${Number(context.parsed.y).toFixed(2)}`;
                },
            },
        },
    },
    scales: {
        y: {
            beginAtZero: false,
            ticks: {
                callback: (value: string | number) => {
                    return `$${Number(value).toFixed(2)}`;
                },
            },
            grid: {
                color: "rgba(0, 0, 0, 0.05)",
            },
        },
        x: {
            grid: {
                display: false,
            },
        },
    },
};
</script>

<style scoped>
.price-chart-container {
    @apply w-full;
}

.chart-wrapper {
    @apply relative w-full;
}
</style>
