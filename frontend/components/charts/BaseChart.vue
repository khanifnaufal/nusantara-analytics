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

// Helper for deep merging options recursively
function deepMerge(target: any, source: any): any {
  if (!source) return target
  const output = { ...target }
  for (const key in source) {
    if (Object.prototype.hasOwnProperty.call(source, key)) {
      if (source[key] && typeof source[key] === 'object' && !Array.isArray(source[key])) {
        output[key] = deepMerge(target[key] || {}, source[key])
      } else {
        output[key] = source[key]
      }
    }
  }
  return output
}

// Global default ECharts config
const defaultOption = {
  grid: {
    top: 8,
    right: 8,
    bottom: 24,
    left: 48
  },
  xAxis: {
    axisLabel: {
      fontSize: 10,
      color: '#52525B'
    }
  },
  yAxis: {
    axisLabel: {
      fontSize: 10,
      color: '#52525B'
    },
    splitLine: {
      lineStyle: {
        color: 'rgba(255,255,255,0.04)'
      }
    }
  },
  tooltip: {
    backgroundColor: '#21262D',
    borderColor: 'rgba(255,255,255,0.1)',
    textStyle: {
      color: '#F0F6FC',
      fontSize: 12
    }
  }
}

const mergedOption = computed(() => {
  return deepMerge(defaultOption, props.option)
})
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
            <div class="absolute inset-0 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
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
        :option="mergedOption"
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
  --tooltip-bg: #111111;
  --tooltip-border: rgba(255, 255, 255, 0.08);
  --tooltip-title-color: #52525B;
  --tooltip-text-color: #A1A1AA;
  --tooltip-val-color: #ffffff;
}
</style>
