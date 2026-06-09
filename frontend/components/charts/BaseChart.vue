<script setup lang="ts">
import { computed } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, BarChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent
} from 'echarts/components'
import VChart from 'vue-echarts'
import { useDark } from '@vueuse/core'

// Register ECharts modules
use([
  CanvasRenderer,
  LineChart,
  BarChart,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent
])

interface Props {
  option: Record<string, any>
  height?: string
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  height: '300px',
  loading: false
})

const isDark = useDark()
const theme = computed(() => (isDark.value ? 'dark' : 'light'))
</script>

<template>
  <ClientOnly>
    <div class="relative w-full" :style="{ height }">
      <!-- Skeleton loader when loading is true -->
      <div
        v-if="loading"
        class="w-full h-full flex flex-col items-center justify-center rounded-xl border border-slate-100 dark:border-slate-800 bg-white/50 dark:bg-slate-900/50 backdrop-blur-sm animate-pulse p-6"
      >
        <div class="flex flex-col items-center gap-3">
          <!-- Spinner -->
          <div class="relative w-10 h-10">
            <div class="absolute inset-0 border-4 border-emerald-500/20 rounded-full"></div>
            <div class="absolute inset-0 border-4 border-emerald-500 border-t-transparent rounded-full animate-spin"></div>
          </div>
          <p class="text-xs text-slate-400 dark:text-slate-500 font-medium tracking-wide">
            Loading chart data...
          </p>
        </div>
      </div>

      <!-- VChart when loading is false -->
      <VChart
        v-show="!loading"
        class="w-full h-full"
        :option="option"
        :theme="theme"
        :autoresize="true"
      />
    </div>
  </ClientOnly>
</template>

<style>
/* CSS Variables for custom styled ECharts tooltip text colors matching dark/light mode */
:root {
  --tooltip-bg: rgba(255, 255, 255, 0.95);
  --tooltip-border: rgba(226, 232, 240, 0.8);
  --tooltip-title-color: #64748b;
  --tooltip-text-color: #475569;
  --tooltip-val-color: #0f172a;
}
html.dark {
  --tooltip-bg: rgba(15, 23, 42, 0.95);
  --tooltip-border: rgba(51, 65, 85, 0.8);
  --tooltip-title-color: #94a3b8;
  --tooltip-text-color: #cbd5e1;
  --tooltip-val-color: #f8fafc;
}
</style>
