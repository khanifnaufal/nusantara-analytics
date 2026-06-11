<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useCanvasStore } from '~/stores/canvas'

interface Props {
  isOpen: boolean
}

defineProps<Props>()

const canvasStore = useCanvasStore()

// Get the active widget object reactively
const widget = computed(() => {
  return canvasStore.widgets.find(w => w.id === canvasStore.selectedWidgetId) || null
})

// Close panel handler
const closePanel = () => {
  canvasStore.selectedWidgetId = null
}

// Data Sources descriptive labels
const DATA_SOURCE_INFO = {
  rates: { title: 'Kurs Mata Uang', icon: '💱' },
  weather: { title: 'Cuaca Kota', icon: '🌤' },
  commodities: { title: 'Komoditas', icon: '📊' },
  market: { title: 'Bursa Saham', icon: '📈' },
  quakes: { title: 'Gempa Bumi', icon: '🌋' }
}

// Chart type presets
const CHART_TYPES = [
  { id: 'line', name: 'Line Chart', icon: '📈' },
  { id: 'bar', name: 'Bar Chart', icon: '📊' },
  { id: 'area', name: 'Area Chart', icon: '🔷' },
  { id: 'donut', name: 'Donut Chart', icon: '🍩' },
  { id: 'scatter', name: 'Scatter Plot', icon: '⚫' },
  { id: 'gauge', name: 'Gauge Chart', icon: '🌡' },
  { id: 'treemap', name: 'Treemap', icon: '🌳' },
  { id: 'radar', name: 'Radar Chart', icon: '📻' },
  { id: 'candlestick', name: 'Candlestick', icon: '🕯' },
  { id: 'heatmap', name: 'Heatmap', icon: '🔥' }
]

// Preset colors list
const PRESET_COLORS = [
  { name: 'blue', hex: '#3B82F6', class: 'bg-[#3B82F6]' },
  { name: 'green', hex: '#10B981', class: 'bg-[#10B981]' },
  { name: 'red', hex: '#EF4444', class: 'bg-[#EF4444]' },
  { name: 'amber', hex: '#F59E0B', class: 'bg-[#F59E0B]' },
  { name: 'violet', hex: '#8B5CF6', class: 'bg-[#8B5CF6]' },
  { name: 'cyan', hex: '#06B6D4', class: 'bg-[#06B6D4]' }
]

// Available metrics lists corresponding to datasource
const METRICS_BY_SOURCE: Record<string, { value: string; name: string }[]> = {
  rates: [
    { value: 'USD', name: 'US Dollar (USD)' },
    { value: 'EUR', name: 'Euro (EUR)' },
    { value: 'SGD', name: 'Singapore Dollar (SGD)' },
    { value: 'JPY', name: 'Japanese Yen (JPY)' },
    { value: 'MYR', name: 'Malaysian Ringgit (MYR)' },
    { value: 'CNY', name: 'Chinese Yuan (CNY)' },
    { value: 'AUD', name: 'Australian Dollar (AUD)' }
  ],
  weather: [
    { value: 'Jakarta', name: 'Jakarta' },
    { value: 'Surabaya', name: 'Surabaya' },
    { value: 'Bandung', name: 'Bandung' },
    { value: 'Medan', name: 'Medan' },
    { value: 'Semarang', name: 'Semarang' },
    { value: 'Makassar', name: 'Makassar' }
  ],
  commodities: [
    { value: 'GC=F', name: 'Emas (Gold)' },
    { value: 'CL=F', name: 'Minyak Mentah' },
    { value: 'PALM.KL', name: 'Minyak Sawit' }
  ],
  market: [
    { value: '^JKSE', name: 'IHSG (Indeks)' },
    { value: 'BBCA.JK', name: 'BCA (BBCA)' },
    { value: 'BBRI.JK', name: 'BRI (BBRI)' },
    { value: 'TLKM.JK', name: 'Telkom (TLKM)' }
  ],
  quakes: [
    { value: 'all', name: 'Semua Wilayah' }
  ]
}

// Metric dropdown state
const isMetricDropdownOpen = ref(false)

const currentMetricLabel = computed(() => {
  const currentWidget = widget.value
  if (!currentWidget) return ''
  const list = METRICS_BY_SOURCE[currentWidget.dataSource] || []
  return list.find(m => m.value === currentWidget.metric)?.name || currentWidget.metric
})

const selectMetric = (val: string) => {
  if (widget.value) {
    widget.value.metric = val
  }
  isMetricDropdownOpen.value = false
}

const toggleMetricDropdown = () => {
  isMetricDropdownOpen.value = !isMetricDropdownOpen.value
}

const closeMetricDropdown = () => {
  isMetricDropdownOpen.value = false
}

// Axis labels and Legend visibility computed setters
const showLabels = computed({
  get: () => widget.value?.config?.showLabels !== false,
  set: (val) => {
    if (widget.value) {
      if (!widget.value.config) widget.value.config = {}
      widget.value.config.showLabels = val
    }
  }
})

const showLegend = computed({
  get: () => widget.value?.config?.showLegend !== false,
  set: (val) => {
    if (widget.value) {
      if (!widget.value.config) widget.value.config = {}
      widget.value.config.showLegend = val
    }
  }
})

// Presets size handlers
const setPresetSize = (w: number, h: number) => {
  if (widget.value) {
    widget.value.size.w = w
    widget.value.size.h = h
  }
}

// Accent color setter
const setAccentColor = (colorName: string) => {
  if (widget.value) {
    if (!widget.value.config) widget.value.config = {}
    widget.value.config.accentColor = colorName
  }
}

// Window click to close dropdowns
const handleWindowClick = () => {
  closeMetricDropdown()
}

onMounted(() => {
  if (typeof window !== 'undefined') {
    window.addEventListener('click', handleWindowClick)
  }
})

onUnmounted(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('click', handleWindowClick)
  }
})
</script>

<template>
  <Transition name="slide">
    <div
      v-if="isOpen && widget"
      class="fixed right-0 top-0 h-screen w-[280px] bg-bg-card border-l border-white/8 z-50 flex flex-col shadow-2xl select-none"
    >
      <!-- Panel Header -->
      <div class="flex items-center justify-between px-4 py-4 border-b border-white/5 bg-[#141414]">
        <div class="flex items-center gap-2">
          <button
            @click="closePanel"
            class="p-1.5 hover:bg-white/5 text-zinc-500 hover:text-white rounded-lg transition-colors cursor-pointer"
            aria-label="Tutup Panel"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
          <span class="text-xs font-bold text-white uppercase tracking-wider">Config Widget</span>
        </div>
      </div>

      <!-- Panel Content -->
      <div class="flex-1 overflow-y-auto p-4 space-y-6 scrollbar-thin">
        
        <!-- SECTION: JUDUL WIDGET -->
        <div class="space-y-1.5">
          <label class="text-[10px] font-bold text-zinc-500 uppercase tracking-wider">Judul Widget</label>
          <input
            type="text"
            v-model="widget.title"
            class="w-full bg-[#1A1A1A] border border-white/10 rounded-xl px-3 py-2 text-xs focus:ring-1 focus:ring-blue-500/50 focus:border-blue-500/50 focus:outline-none transition-all text-white placeholder-zinc-500"
            placeholder="Edit nama widget..."
          />
        </div>

        <!-- SECTION: DATA & METRIK -->
        <div class="space-y-3">
          <div class="space-y-1.5">
            <label class="text-[10px] font-bold text-zinc-500 uppercase tracking-wider">Sumber Data</label>
            <div class="w-full flex items-center gap-2 bg-[#1A1A1A] border border-white/5 rounded-xl px-3 py-2 text-xs text-zinc-400">
              <span class="text-base leading-none">
                {{ DATA_SOURCE_INFO[widget.dataSource]?.icon || '📊' }}
              </span>
              <span>
                {{ DATA_SOURCE_INFO[widget.dataSource]?.title || widget.dataSource }}
              </span>
            </div>
          </div>

          <div class="relative space-y-1.5">
            <label class="text-[10px] font-bold text-zinc-500 uppercase tracking-wider">Metrik Data</label>
            <button
              type="button"
              @click.stop="toggleMetricDropdown"
              class="w-full flex items-center justify-between text-xs bg-[#1A1A1A] border border-white/10 rounded-xl px-3 py-2 text-white hover:border-white/20 transition-colors cursor-pointer text-left focus:outline-none"
            >
              <span>{{ currentMetricLabel }}</span>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 text-zinc-400 transition-transform duration-200" :class="{ 'rotate-180': isMetricDropdownOpen }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </button>
            
            <div
              v-if="isMetricDropdownOpen"
              class="absolute left-0 right-0 mt-1.5 z-50 max-h-40 overflow-y-auto rounded-xl border border-white/10 bg-bg-subcard shadow-xl py-1 scrollbar-thin animate-in fade-in slide-in-from-top-1 duration-150"
            >
              <button
                v-for="m in (METRICS_BY_SOURCE[widget.dataSource] || [])"
                :key="m.value"
                type="button"
                @click="selectMetric(m.value)"
                class="w-full text-left px-3 py-2 text-xs text-zinc-300 hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
              >
                <span>{{ m.name }}</span>
                <span v-if="widget.metric === m.value" class="text-blue-500 font-bold">✓</span>
              </button>
            </div>
          </div>
        </div>

        <!-- SECTION: TIPE CHART -->
        <div class="space-y-2">
          <label class="text-[10px] font-bold text-zinc-500 uppercase tracking-wider">Tipe Chart</label>
          <div class="grid grid-cols-2 gap-1.5">
            <button
              v-for="chart in CHART_TYPES"
              :key="chart.id"
              type="button"
              @click="widget.chartType = chart.id"
              class="flex items-center gap-2 p-2 rounded-xl border transition-all text-left cursor-pointer active:scale-[0.97]"
              :class="[
                widget.chartType === chart.id
                  ? 'border-blue-500/50 bg-blue-500/10 text-white shadow-md shadow-blue-500/5'
                  : 'border-white/5 bg-[#1A1A1A] hover:bg-zinc-800 text-zinc-450 hover:text-zinc-200'
              ]"
            >
              <span class="text-base leading-none">{{ chart.icon }}</span>
              <span class="text-[10px] font-semibold truncate">{{ chart.name.split(' ')[0] }}</span>
            </button>
          </div>
        </div>

        <!-- SECTION: TAMPILAN -->
        <div class="space-y-3 border-t border-white/5 pt-4">
          <label class="text-[10px] font-bold text-zinc-500 uppercase tracking-wider block">Tampilan</label>
          
          <div class="flex items-center justify-between py-0.5">
            <span class="text-xs text-zinc-300">Tampilkan Label</span>
            <button
              type="button"
              @click="showLabels = !showLabels"
              class="w-8 h-4.5 rounded-full relative transition-colors duration-200 focus:outline-none cursor-pointer"
              :class="showLabels ? 'bg-blue-600' : 'bg-zinc-800'"
            >
              <span
                class="absolute top-0.5 left-0.5 w-3.5 h-3.5 rounded-full bg-white transition-transform duration-200"
                :class="showLabels ? 'transform translate-x-3.5' : ''"
              ></span>
            </button>
          </div>

          <div class="flex items-center justify-between py-0.5">
            <span class="text-xs text-zinc-300">Tampilkan Legend</span>
            <button
              type="button"
              @click="showLegend = !showLegend"
              class="w-8 h-4.5 rounded-full relative transition-colors duration-200 focus:outline-none cursor-pointer"
              :class="showLegend ? 'bg-blue-600' : 'bg-zinc-800'"
            >
              <span
                class="absolute top-0.5 left-0.5 w-3.5 h-3.5 rounded-full bg-white transition-transform duration-200"
                :class="showLegend ? 'transform translate-x-3.5' : ''"
              ></span>
            </button>
          </div>

          <!-- Color Accent Picker -->
          <div class="space-y-1.5 pt-1">
            <span class="text-[10px] text-zinc-400 block font-semibold">Warna Aksen</span>
            <div class="flex items-center gap-2">
              <button
                v-for="color in PRESET_COLORS"
                :key="color.name"
                type="button"
                @click="setAccentColor(color.name)"
                class="w-6 h-6 rounded-full border border-black/20 relative group transition-all hover:scale-110 cursor-pointer shadow-md focus:outline-none"
                :class="[
                  color.class,
                  (widget.config?.accentColor === color.name || (!widget.config?.accentColor && color.name === 'blue'))
                    ? 'ring-2 ring-white/60 scale-105'
                    : 'opacity-70 hover:opacity-100'
                ]"
                :title="color.name"
              >
                <span
                  v-if="widget.config?.accentColor === color.name || (!widget.config?.accentColor && color.name === 'blue')"
                  class="absolute inset-0 flex items-center justify-center text-white text-[9px] font-bold"
                >
                  ✓
                </span>
              </button>
            </div>
          </div>
        </div>

        <!-- SECTION: UKURAN -->
        <div class="space-y-3.5 border-t border-white/5 pt-4">
          <label class="text-[10px] font-bold text-zinc-500 uppercase tracking-wider block">Ukuran (px)</label>
          
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1">
              <label class="text-[9px] font-bold text-zinc-500 uppercase tracking-wider">Lebar</label>
              <input
                type="number"
                min="280"
                v-model.number="widget.size.w"
                class="w-full bg-[#1A1A1A] border border-white/10 rounded-xl px-2.5 py-1.5 text-xs text-white focus:outline-none focus:border-blue-500/50"
              />
            </div>
            <div class="space-y-1">
              <label class="text-[9px] font-bold text-zinc-500 uppercase tracking-wider">Tinggi</label>
              <input
                type="number"
                min="200"
                v-model.number="widget.size.h"
                class="w-full bg-[#1A1A1A] border border-white/10 rounded-xl px-2.5 py-1.5 text-xs text-white focus:outline-none focus:border-blue-500/50"
              />
            </div>
          </div>

          <div class="space-y-1.5">
            <span class="text-[9px] font-bold text-zinc-500 uppercase tracking-wider">Preset Ukuran</span>
            <div class="grid grid-cols-3 gap-1.5">
              <button
                type="button"
                @click="setPresetSize(300, 200)"
                class="py-1.5 text-[10px] font-bold rounded-lg bg-[#1A1A1A] border border-white/5 hover:border-white/10 text-zinc-300 hover:text-white transition-colors cursor-pointer"
              >
                S
              </button>
              <button
                type="button"
                @click="setPresetSize(400, 300)"
                class="py-1.5 text-[10px] font-bold rounded-lg bg-[#1A1A1A] border border-white/5 hover:border-white/10 text-zinc-300 hover:text-white transition-colors cursor-pointer"
              >
                M
              </button>
              <button
                type="button"
                @click="setPresetSize(600, 400)"
                class="py-1.5 text-[10px] font-bold rounded-lg bg-[#1A1A1A] border border-white/5 hover:border-white/10 text-zinc-300 hover:text-white transition-colors cursor-pointer"
              >
                L
              </button>
            </div>
          </div>
        </div>

      </div>
    </div>
  </Transition>
</template>

<style scoped>
.slide-enter-active,
.slide-leave-active {
  transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}
.slide-enter-from,
.slide-leave-to {
  transform: translateX(100%);
}
</style>
