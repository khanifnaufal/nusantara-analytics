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
        class="w-full h-full flex flex-col items-center justify-center rounded-2xl bg-bg-subcard animate-pulse p-6 border border-white/5"
      >
        <div class="flex flex-col items-center gap-3">
          <!-- Spinner -->
          <div class="relative w-10 h-10">
            <div class="absolute inset-0 border-4 border-white/5 rounded-full"></div>
            <div class="absolute inset-0 border-4 border-emerald-500 border-t-transparent rounded-full animate-spin"></div>
          </div>
          <p class="text-xs text-text-secondary font-medium tracking-wide">
            Memuat data grafik...
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
  --tooltip-bg: #21262D;
  --tooltip-border: rgba(255, 255, 255, 0.1);
  --tooltip-title-color: #8B949E;
  --tooltip-text-color: #F0F6FC;
  --tooltip-val-color: #ffffff;
}
</style>
