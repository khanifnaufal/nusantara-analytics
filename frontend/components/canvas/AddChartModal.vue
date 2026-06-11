<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRatesStore } from '~/stores/rates'
import { useWeatherStore } from '~/stores/weather'
import { useCommoditiesStore } from '~/stores/commodities'
import { useMarketStore } from '~/stores/market'
import { useQuakesStore } from '~/stores/quakes'

interface Props {
  isOpen: boolean
}

const props = defineProps<Props>()
const emit = defineEmits(['close', 'add'])

const currentStep = ref(1)

// Form states
const selectedSource = ref<'rates' | 'weather' | 'commodities' | 'market' | 'quakes' | null>(null)
const selectedMetricItem = ref('') // e.g. USD, Jakarta, GC=F, etc.
const selectedMetricField = ref('') // e.g. 'nilai' / 'history' / 'temperature' / 'magnitude distribution'
const selectedChartType = ref('')
const widgetTitle = ref('')

// Custom dropdowns states
const activeCustomDropdown = ref<string | null>(null)
const toggleCustomDropdown = (name: string) => {
  activeCustomDropdown.value = activeCustomDropdown.value === name ? null : name
}
const closeCustomDropdowns = () => {
  activeCustomDropdown.value = null
}

onMounted(() => {
  if (typeof window !== 'undefined') {
    window.addEventListener('click', closeCustomDropdowns)
  }
})

onUnmounted(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('click', closeCustomDropdowns)
  }
})

// Initialize stores to prefetch if needed or get list
const ratesStore = useRatesStore()
const weatherStore = useWeatherStore()
const commoditiesStore = useCommoditiesStore()
const marketStore = useMarketStore()
const quakesStore = useQuakesStore()

// Step metadata lists
const DATA_SOURCES = [
  { id: 'rates', title: 'Kurs Mata Uang', icon: '💱', desc: 'Nilai tukar mata uang asing ke Rupiah' },
  { id: 'weather', title: 'Cuaca Kota', icon: '🌤', desc: 'Suhu, kelembapan, dan kecepatan angin kota utama' },
  { id: 'commodities', title: 'Komoditas', icon: '📊', desc: 'Harga emas, minyak mentah, dan CPO global' },
  { id: 'market', title: 'Bursa Saham', icon: '📈', desc: 'Indeks IHSG dan pergerakan saham blue chip' },
  { id: 'quakes', title: 'Gempa Bumi', icon: '🌋', desc: 'Aktivitas gempa bumi BMKG terbaru' }
]

const RATES_CURRENCIES = [
  { value: 'USD', name: 'US Dollar (USD)' },
  { value: 'EUR', name: 'Euro (EUR)' },
  { value: 'SGD', name: 'Singapore Dollar (SGD)' },
  { value: 'JPY', name: 'Japanese Yen (JPY)' },
  { value: 'MYR', name: 'Malaysian Ringgit (MYR)' },
  { value: 'CNY', name: 'Chinese Yuan (CNY)' },
  { value: 'AUD', name: 'Australian Dollar (AUD)' }
]

const WEATHER_CITIES = [
  { value: 'Jakarta', name: 'Jakarta' },
  { value: 'Surabaya', name: 'Surabaya' },
  { value: 'Bandung', name: 'Bandung' },
  { value: 'Medan', name: 'Medan' },
  { value: 'Semarang', name: 'Semarang' },
  { value: 'Makassar', name: 'Makassar' }
]

const COMMODITIES_LIST = [
  { value: 'GC=F', name: 'Emas (Gold)' },
  { value: 'CL=F', name: 'Minyak Mentah (Crude Oil)' },
  { value: 'PALM.KL', name: 'Minyak Sawit (CPO)' }
]

const MARKET_STOCKS = [
  { value: '^JKSE', name: 'IHSG (Indeks Gabungan)' },
  { value: 'BBCA.JK', name: 'Bank Central Asia (BBCA)' },
  { value: 'BBRI.JK', name: 'Bank Rakyat Indonesia (BBRI)' },
  { value: 'TLKM.JK', name: 'Telkom Indonesia (TLKM)' }
]

const QUAKES_VIZ_OPTIONS = [
  { value: 'magnitude distribution', name: 'Distribusi Magnitudo' },
  { value: 'depth distribution', name: 'Distribusi Kedalaman (Scatter)' },
  { value: 'timeline', name: 'Timeline Aktivitas Gempa' }
]

const CHART_TYPES = [
  { id: 'line', name: 'Line Chart', icon: '📈', desc: 'Tren data kontinu (contoh: riwayat harga saham BBCA, nilai kurs USD)' },
  { id: 'bar', name: 'Bar Chart', icon: '📊', desc: 'Perbandingan data diskrit (contoh: suhu udara antar kota, jumlah gempa per magnitudo)' },
  { id: 'area', name: 'Area Chart', icon: '🔷', desc: 'Tren akumulasi atau fluktuasi (contoh: harga emas global, tren suhu Jakarta)' },
  { id: 'donut', name: 'Donut Chart', icon: '🍩', desc: 'Proporsi kontribusi total (contoh: persentase jumlah gempa rendah vs tinggi)' },
  { id: 'scatter', name: 'Scatter Plot', icon: '⚫', desc: 'Hubungan 2 variabel numerik (contoh: kedalaman vs magnitudo gempa bumi)' },
  { id: 'gauge', name: 'Gauge Chart', icon: '🌡', desc: 'Menampilkan nilai tunggal / meteran batas (contoh: suhu saat ini, magnitudo gempa terbaru)' },
  { id: 'treemap', name: 'Treemap', icon: '🌳', desc: 'Struktur hierarki & pembagian area (contoh: segmentasi tingkat gempa berdasarkan wilayah)' },
  { id: 'radar', name: 'Radar Chart', icon: '📻', desc: 'Perbandingan multivariat (contoh: profil suhu vs kelembapan vs angin suatu kota)' },
  { id: 'candlestick', name: 'Candlestick', icon: '🕯', desc: 'Rentang harga pembukaan-penutupan (contoh: OHLC saham BBRI/IHSG historis)' },
  { id: 'heatmap', name: 'Heatmap', icon: '🔥', desc: 'Kepadatan intensitas grid (contoh: matriks suhu kota vs hari dalam seminggu)' }
]

// Fetch data stores on mount / open
watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    currentStep.value = 1
    selectedSource.value = null
    selectedMetricItem.value = ''
    selectedMetricField.value = ''
    selectedChartType.value = ''
    widgetTitle.value = ''
    closeCustomDropdowns()
    
    // Fetch if needed
    if (!ratesStore.data) ratesStore.fetchRates()
    if (!weatherStore.data) weatherStore.fetchWeather()
    if (!commoditiesStore.data) commoditiesStore.fetchCommodities()
    if (!marketStore.data) marketStore.fetchMarket()
    if (!quakesStore.data) quakesStore.fetchQuakes()
  }
})

// Auto-select defaults when source changes
const handleSourceSelect = (source: 'rates' | 'weather' | 'commodities' | 'market' | 'quakes') => {
  selectedSource.value = source
  
  if (source === 'rates') {
    selectedMetricItem.value = 'USD'
    selectedMetricField.value = 'history'
  } else if (source === 'weather') {
    selectedMetricItem.value = 'Jakarta'
    selectedMetricField.value = 'temperature'
  } else if (source === 'commodities') {
    selectedMetricItem.value = 'GC=F'
    selectedMetricField.value = 'history'
  } else if (source === 'market') {
    selectedMetricItem.value = '^JKSE'
    selectedMetricField.value = 'history'
  } else if (source === 'quakes') {
    selectedMetricItem.value = 'all'
    selectedMetricField.value = 'timeline'
  }
  
  // Set default recommended chart types
  if (source === 'rates' || source === 'commodities' || source === 'market') {
    selectedChartType.value = 'line'
  } else if (source === 'weather') {
    selectedChartType.value = 'gauge'
  } else if (source === 'quakes') {
    selectedChartType.value = 'bar'
  }
  
  goToStep(2)
}

// Generate pre-filled title based on choices
const generatedTitle = computed(() => {
  if (!selectedSource.value) return ''
  
  let sourceName = ''
  let itemName = ''
  let fieldName = ''
  
  if (selectedSource.value === 'rates') {
    sourceName = 'Kurs'
    itemName = selectedMetricItem.value
    fieldName = selectedMetricField.value === 'history' ? 'Tren' : 'Nilai'
    return `${fieldName} ${sourceName} ${itemName}`
  }
  
  if (selectedSource.value === 'weather') {
    itemName = selectedMetricItem.value
    const fObj = selectedMetricField.value === 'temperature' ? 'Suhu' : selectedMetricField.value === 'humidity' ? 'Kelembapan' : 'Kecepatan Angin'
    return `${fObj} ${itemName}`
  }
  
  if (selectedSource.value === 'commodities') {
    const cObj = COMMODITIES_LIST.find(c => c.value === selectedMetricItem.value)
    itemName = (cObj ? cObj.name.split(' ')[0] : 'Komoditas') || 'Komoditas'
    fieldName = selectedMetricField.value === 'history' ? 'Tren' : 'Harga'
    return `${fieldName} ${itemName}`
  }
  
  if (selectedSource.value === 'market') {
    const sObj = MARKET_STOCKS.find(s => s.value === selectedMetricItem.value)
    itemName = (sObj ? sObj.name.split(' ')[0] : 'Saham') || 'Saham'
    fieldName = selectedMetricField.value === 'history' ? 'Tren' : 'Harga'
    return `${fieldName} ${itemName}`
  }
  
  if (selectedSource.value === 'quakes') {
    if (selectedMetricField.value === 'magnitude distribution') return 'Distribusi Magnitudo Gempa'
    if (selectedMetricField.value === 'depth distribution') return 'Scatter Plot Kedalaman Gempa'
    return 'Timeline Aktivitas Gempa Bumi'
  }
  
  return ''
})

// Set widget title whenever step reaches 4
watch(currentStep, (newStep) => {
  if (newStep === 4) {
    widgetTitle.value = generatedTitle.value
  }
})

const goToStep = (step: number) => {
  if (step > currentStep.value) {
    // Validate current step
    if (currentStep.value === 1 && !selectedSource.value) return
    if (currentStep.value === 2 && (!selectedMetricItem.value || !selectedMetricField.value)) return
    if (currentStep.value === 3 && !selectedChartType.value) return
  }
  currentStep.value = step
}

const addWidget = () => {
  if (!widgetTitle.value.trim()) {
    widgetTitle.value = generatedTitle.value || 'Custom Widget'
  }
  
  const widgetData = {
    title: widgetTitle.value,
    dataSource: selectedSource.value!,
    metric: selectedMetricItem.value,
    chartType: selectedChartType.value,
    position: {
      x: 100 + Math.floor(Math.random() * 60),
      y: 100 + Math.floor(Math.random() * 60)
    },
    size: {
      w: selectedChartType.value === 'heatmap' || selectedChartType.value === 'treemap' ? 520 : 400,
      h: 300
    },
    config: {
      field: selectedMetricField.value
    }
  }
  
  emit('add', widgetData)
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <!-- Backdrop -->
    <div class="absolute inset-0 bg-black/70 backdrop-blur-md transition-opacity" @click="emit('close')"></div>
    
    <!-- Modal container -->
    <div class="relative w-full max-w-lg rounded-2xl border border-white/10 bg-[#111111] p-6 shadow-2xl z-10 flex flex-col max-h-[85vh] overflow-hidden transition-all duration-300">
      
      <!-- Modal Header -->
      <div class="flex items-center justify-between border-b border-white/5 pb-4 mb-4">
        <div>
          <h3 class="text-base font-extrabold text-white flex items-center gap-2">
            <span>⚙️</span> Tambah Widget Baru
          </h3>
          <p class="text-xs text-text-secondary mt-0.5">Custom canvas analytics builder</p>
        </div>
        <button 
          @click="emit('close')" 
          class="text-zinc-500 hover:text-white transition-colors p-1.5 hover:bg-white/5 rounded-lg text-lg cursor-pointer"
        >
          ×
        </button>
      </div>

      <!-- Step Indicator -->
      <div class="flex items-center justify-between px-1 mb-5">
        <div 
          v-for="s in 4" 
          :key="'step-'+s"
          class="flex items-center flex-1"
        >
          <button
            @click="s < currentStep || (selectedSource && s <= 3) ? goToStep(s) : null"
            :disabled="s > currentStep && (!selectedSource || s === 4)"
            class="w-6 h-6 rounded-full flex items-center justify-center text-[10px] font-extrabold transition-all"
            :class="[
              currentStep === s 
                ? 'bg-blue-500 text-white ring-4 ring-blue-500/25' 
                : s < currentStep 
                  ? 'bg-blue-900/40 text-blue-400 border border-blue-500/30' 
                  : 'bg-zinc-800 text-zinc-500 border border-zinc-700/50'
            ]"
          >
            {{ s }}
          </button>
          <div 
            v-if="s < 4" 
            class="h-0.5 flex-1 mx-2 transition-colors duration-300"
            :class="s < currentStep ? 'bg-blue-500/50' : 'bg-zinc-800'"
          ></div>
        </div>
      </div>

      <!-- Scrollable content area -->
      <div class="flex-1 overflow-y-auto pr-1 scrollbar-thin">
        
        <!-- STEP 1: PILIH DATA SOURCE -->
        <div v-if="currentStep === 1" class="space-y-3">
          <label class="block text-xs font-bold text-zinc-400 uppercase tracking-wider mb-2">
            Langkah 1: Pilih Sumber Data
          </label>
          <div class="grid grid-cols-1 gap-2.5">
            <button
              v-for="source in DATA_SOURCES"
              :key="source.id"
              @click="handleSourceSelect(source.id as any)"
              class="flex items-center gap-4 text-left p-3.5 rounded-xl border border-white/5 bg-[#161616] hover:bg-zinc-800/80 hover:border-white/10 active:scale-[0.99] transition-all cursor-pointer group"
              :class="{ 'border-blue-500/50 bg-blue-950/10 shadow-lg shadow-blue-500/5': selectedSource === source.id }"
            >
              <span class="text-2xl p-2 rounded-lg bg-zinc-850 group-hover:scale-110 transition-transform">{{ source.icon }}</span>
              <div class="flex-1 min-w-0">
                <h4 class="text-xs font-extrabold text-white">{{ source.title }}</h4>
                <p class="text-[10px] text-text-secondary mt-0.5 line-clamp-1">{{ source.desc }}</p>
              </div>
              <span 
                class="w-4 h-4 rounded-full border flex items-center justify-center text-[10px] font-bold"
                :class="selectedSource === source.id ? 'border-blue-500 bg-blue-500 text-white' : 'border-zinc-700 text-transparent'"
              >
                ✓
              </span>
            </button>
          </div>
        </div>

        <!-- STEP 2: PILIH METRIK -->
        <div v-if="currentStep === 2" class="space-y-4">
          <label class="block text-xs font-bold text-zinc-400 uppercase tracking-wider mb-2">
            Langkah 2: Konfigurasi Parameter & Metrik
          </label>

          <!-- Rates Settings -->
          <div v-if="selectedSource === 'rates'" class="space-y-3.5">
            <div class="relative">
              <label class="block text-[10px] font-bold text-text-secondary uppercase tracking-wider mb-1.5">Mata Uang</label>
              <button
                type="button"
                @click.stop="toggleCustomDropdown('ratesCurrency')"
                class="w-full flex items-center justify-between text-xs rounded-lg border border-white/8 bg-[#161616] px-3 py-2 text-white focus:border-blue-500/50 focus:outline-none cursor-pointer"
              >
                <span>{{ RATES_CURRENCIES.find(c => c.value === selectedMetricItem)?.name || 'Pilih Mata Uang' }}</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-zinc-400 transition-transform duration-200" :class="{ 'rotate-180': activeCustomDropdown === 'ratesCurrency' }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
              <div
                v-if="activeCustomDropdown === 'ratesCurrency'"
                class="absolute left-0 right-0 mt-1.5 z-50 max-h-48 overflow-y-auto rounded-lg border border-white/10 bg-[#161616] shadow-xl py-1 scrollbar-thin animate-in fade-in slide-in-from-top-1 duration-200"
              >
                <button
                  v-for="c in RATES_CURRENCIES"
                  :key="c.value"
                  type="button"
                  @click="selectedMetricItem = c.value; closeCustomDropdowns()"
                  class="w-full text-left px-3 py-2 text-xs text-zinc-300 hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
                >
                  <span>{{ c.name }}</span>
                  <span v-if="selectedMetricItem === c.value" class="text-blue-500 font-bold">✓</span>
                </button>
              </div>
            </div>
            <div>
              <label class="block text-[10px] font-bold text-text-secondary uppercase tracking-wider mb-1.5">Tipe Tampilan Data</label>
              <div class="grid grid-cols-2 gap-2">
                <button 
                  @click="selectedMetricField = 'nilai'"
                  class="p-2.5 rounded-lg border text-center text-xs font-bold cursor-pointer transition-all"
                  :class="selectedMetricField === 'nilai' ? 'bg-blue-600 border-blue-500 text-white shadow-md shadow-blue-600/10' : 'bg-[#161616] border-white/5 text-zinc-400 hover:bg-zinc-800'"
                >
                  💵 Nilai Terkini
                </button>
                <button 
                  @click="selectedMetricField = 'history'"
                  class="p-2.5 rounded-lg border text-center text-xs font-bold cursor-pointer transition-all"
                  :class="selectedMetricField === 'history' ? 'bg-blue-600 border-blue-500 text-white shadow-md shadow-blue-600/10' : 'bg-[#161616] border-white/5 text-zinc-400 hover:bg-zinc-800'"
                >
                  📈 Riwayat Historis
                </button>
              </div>
            </div>
          </div>

          <!-- Weather Settings -->
          <div v-if="selectedSource === 'weather'" class="space-y-3.5">
            <div class="relative">
              <label class="block text-[10px] font-bold text-text-secondary uppercase tracking-wider mb-1.5">Pilih Kota</label>
              <button
                type="button"
                @click.stop="toggleCustomDropdown('weatherCity')"
                class="w-full flex items-center justify-between text-xs rounded-lg border border-white/8 bg-[#161616] px-3 py-2 text-white focus:border-blue-500/50 focus:outline-none cursor-pointer"
              >
                <span>🏙️ {{ WEATHER_CITIES.find(c => c.value === selectedMetricItem)?.name || 'Pilih Kota' }}</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-zinc-400 transition-transform duration-200" :class="{ 'rotate-180': activeCustomDropdown === 'weatherCity' }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
              <div
                v-if="activeCustomDropdown === 'weatherCity'"
                class="absolute left-0 right-0 mt-1.5 z-50 max-h-48 overflow-y-auto rounded-lg border border-white/10 bg-[#161616] shadow-xl py-1 scrollbar-thin"
              >
                <button
                  v-for="c in WEATHER_CITIES"
                  :key="c.value"
                  type="button"
                  @click="selectedMetricItem = c.value; closeCustomDropdowns()"
                  class="w-full text-left px-3 py-2 text-xs text-zinc-300 hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
                >
                  <span>{{ c.name }}</span>
                  <span v-if="selectedMetricItem === c.value" class="text-blue-500 font-bold">✓</span>
                </button>
              </div>
            </div>
            <div class="relative">
              <label class="block text-[10px] font-bold text-text-secondary uppercase tracking-wider mb-1.5">Parameter Cuaca</label>
              <button
                type="button"
                @click.stop="toggleCustomDropdown('weatherField')"
                class="w-full flex items-center justify-between text-xs rounded-lg border border-white/8 bg-[#161616] px-3 py-2 text-white focus:border-blue-500/50 focus:outline-none cursor-pointer"
              >
                <span>{{ selectedMetricField === 'temperature' ? '🌡 Suhu Udara (°C)' : selectedMetricField === 'humidity' ? '💧 Kelembapan (%)' : '💨 Kecepatan Angin (km/h)' }}</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-zinc-400 transition-transform duration-200" :class="{ 'rotate-180': activeCustomDropdown === 'weatherField' }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
              <div
                v-if="activeCustomDropdown === 'weatherField'"
                class="absolute left-0 right-0 mt-1.5 z-50 max-h-48 overflow-y-auto rounded-lg border border-white/10 bg-[#161616] shadow-xl py-1 scrollbar-thin"
              >
                <button
                  type="button"
                  @click="selectedMetricField = 'temperature'; closeCustomDropdowns()"
                  class="w-full text-left px-3 py-2 text-xs text-zinc-300 hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
                >
                  <span>🌡 Suhu Udara (°C)</span>
                  <span v-if="selectedMetricField === 'temperature'" class="text-blue-500 font-bold">✓</span>
                </button>
                <button
                  type="button"
                  @click="selectedMetricField = 'humidity'; closeCustomDropdowns()"
                  class="w-full text-left px-3 py-2 text-xs text-zinc-300 hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
                >
                  <span>💧 Kelembapan (%)</span>
                  <span v-if="selectedMetricField === 'humidity'" class="text-blue-500 font-bold">✓</span>
                </button>
                <button
                  type="button"
                  @click="selectedMetricField = 'windSpeed'; closeCustomDropdowns()"
                  class="w-full text-left px-3 py-2 text-xs text-zinc-300 hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
                >
                  <span>💨 Kecepatan Angin (km/h)</span>
                  <span v-if="selectedMetricField === 'windSpeed'" class="text-blue-500 font-bold">✓</span>
                </button>
              </div>
            </div>
          </div>

          <!-- Commodities Settings -->
          <div v-if="selectedSource === 'commodities'" class="space-y-3.5">
            <div class="relative">
              <label class="block text-[10px] font-bold text-text-secondary uppercase tracking-wider mb-1.5">Komoditas</label>
              <button
                type="button"
                @click.stop="toggleCustomDropdown('commoditiesItem')"
                class="w-full flex items-center justify-between text-xs rounded-lg border border-white/8 bg-[#161616] px-3 py-2 text-white focus:border-blue-500/50 focus:outline-none cursor-pointer"
              >
                <span>📊 {{ COMMODITIES_LIST.find(c => c.value === selectedMetricItem)?.name || 'Pilih Komoditas' }}</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-zinc-400 transition-transform duration-200" :class="{ 'rotate-180': activeCustomDropdown === 'commoditiesItem' }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
              <div
                v-if="activeCustomDropdown === 'commoditiesItem'"
                class="absolute left-0 right-0 mt-1.5 z-50 max-h-48 overflow-y-auto rounded-lg border border-white/10 bg-[#161616] shadow-xl py-1 scrollbar-thin"
              >
                <button
                  v-for="c in COMMODITIES_LIST"
                  :key="c.value"
                  type="button"
                  @click="selectedMetricItem = c.value; closeCustomDropdowns()"
                  class="w-full text-left px-3 py-2 text-xs text-zinc-300 hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
                >
                  <span>{{ c.name }}</span>
                  <span v-if="selectedMetricItem === c.value" class="text-blue-500 font-bold">✓</span>
                </button>
              </div>
            </div>
            <div>
              <label class="block text-[10px] font-bold text-text-secondary uppercase tracking-wider mb-1.5">Tipe Tampilan</label>
              <div class="grid grid-cols-2 gap-2">
                <button 
                  @click="selectedMetricField = 'nilai'"
                  class="p-2.5 rounded-lg border text-center text-xs font-bold cursor-pointer transition-all"
                  :class="selectedMetricField === 'nilai' ? 'bg-blue-600 border-blue-500 text-white' : 'bg-[#161616] border-white/5 text-zinc-400 hover:bg-zinc-800'"
                >
                  💰 Harga Terkini
                </button>
                <button 
                  @click="selectedMetricField = 'history'"
                  class="p-2.5 rounded-lg border text-center text-xs font-bold cursor-pointer transition-all"
                  :class="selectedMetricField === 'history' ? 'bg-blue-600 border-blue-500 text-white' : 'bg-[#161616] border-white/5 text-zinc-400 hover:bg-zinc-800'"
                >
                  📈 Riwayat Historis
                </button>
              </div>
            </div>
          </div>

          <!-- Market Settings -->
          <div v-if="selectedSource === 'market'" class="space-y-3.5">
            <div class="relative">
              <label class="block text-[10px] font-bold text-text-secondary uppercase tracking-wider mb-1.5">Pilih Indeks / Saham</label>
              <button
                type="button"
                @click.stop="toggleCustomDropdown('marketItem')"
                class="w-full flex items-center justify-between text-xs rounded-lg border border-white/8 bg-[#161616] px-3 py-2 text-white focus:border-blue-500/50 focus:outline-none cursor-pointer"
              >
                <span>📈 {{ MARKET_STOCKS.find(s => s.value === selectedMetricItem)?.name || 'Pilih Saham' }}</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-zinc-400 transition-transform duration-200" :class="{ 'rotate-180': activeCustomDropdown === 'marketItem' }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
              <div
                v-if="activeCustomDropdown === 'marketItem'"
                class="absolute left-0 right-0 mt-1.5 z-50 max-h-48 overflow-y-auto rounded-lg border border-white/10 bg-[#161616] shadow-xl py-1 scrollbar-thin"
              >
                <button
                  v-for="s in MARKET_STOCKS"
                  :key="s.value"
                  type="button"
                  @click="selectedMetricItem = s.value; closeCustomDropdowns()"
                  class="w-full text-left px-3 py-2 text-xs text-zinc-300 hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
                >
                  <span>{{ s.name }}</span>
                  <span v-if="selectedMetricItem === s.value" class="text-blue-500 font-bold">✓</span>
                </button>
              </div>
            </div>
            <div>
              <label class="block text-[10px] font-bold text-text-secondary uppercase tracking-wider mb-1.5">Tipe Tampilan</label>
              <div class="grid grid-cols-2 gap-2">
                <button 
                  @click="selectedMetricField = 'nilai'"
                  class="p-2.5 rounded-lg border text-center text-xs font-bold cursor-pointer transition-all"
                  :class="selectedMetricField === 'nilai' ? 'bg-blue-600 border-blue-500 text-white' : 'bg-[#161616] border-white/5 text-zinc-400 hover:bg-zinc-800'"
                >
                  🏛 Harga / Indeks Terkini
                </button>
                <button 
                  @click="selectedMetricField = 'history'"
                  class="p-2.5 rounded-lg border text-center text-xs font-bold cursor-pointer transition-all"
                  :class="selectedMetricField === 'history' ? 'bg-blue-600 border-blue-500 text-white' : 'bg-[#161616] border-white/5 text-zinc-400 hover:bg-zinc-800'"
                >
                  📈 Riwayat Tren
                </button>
              </div>
            </div>
          </div>

          <!-- Quakes Settings -->
          <div v-if="selectedSource === 'quakes'" class="space-y-3.5">
            <div>
              <label class="block text-[10px] font-bold text-text-secondary uppercase tracking-wider mb-1.5">Visualisasi Data Gempa Bumi</label>
              <div class="grid grid-cols-1 gap-2">
                <button 
                  v-for="opt in QUAKES_VIZ_OPTIONS"
                  :key="opt.value"
                  @click="selectedMetricField = opt.value"
                  class="p-3 rounded-lg border text-left text-xs font-bold cursor-pointer transition-all flex items-center justify-between"
                  :class="selectedMetricField === opt.value ? 'bg-blue-600 border-blue-500 text-white shadow-md shadow-blue-600/10' : 'bg-[#161616] border-white/5 text-zinc-400 hover:bg-zinc-800'"
                >
                  <span>{{ opt.name }}</span>
                  <span v-if="selectedMetricField === opt.value" class="text-xs">✓</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- STEP 3: PILIH CHART TYPE -->
        <div v-if="currentStep === 3" class="space-y-3">
          <label class="block text-xs font-bold text-zinc-400 uppercase tracking-wider mb-2">
            Langkah 3: Pilih Tipe Grafik (ECharts)
          </label>
          <div class="grid grid-cols-2 gap-2.5">
            <button
              v-for="type in CHART_TYPES"
              :key="type.id"
              @click="selectedChartType = type.id"
              class="flex flex-col items-center text-center p-3.5 rounded-xl border transition-all cursor-pointer hover:scale-[1.01] active:scale-[0.98]"
              :class="[
                selectedChartType === type.id 
                  ? 'border-blue-400 bg-blue-600 text-white shadow-lg shadow-blue-500/25 ring-2 ring-blue-400/20 scale-[1.02]' 
                  : 'border-white/5 bg-[#161616] hover:bg-zinc-800 text-zinc-400'
              ]"
            >
              <span class="text-2xl mb-1.5">{{ type.icon }}</span>
              <h4 class="text-xs font-extrabold" :class="selectedChartType === type.id ? 'text-white' : 'text-zinc-200'">{{ type.name }}</h4>
              <p class="text-[11px] mt-1.5 px-1 leading-relaxed" :class="selectedChartType === type.id ? 'text-blue-100' : 'text-zinc-400'">{{ type.desc }}</p>
            </button>
          </div>
        </div>

        <!-- STEP 4: EDITABLE TITLE & ADD -->
        <div v-if="currentStep === 4" class="space-y-4">
          <label class="block text-xs font-bold text-zinc-400 uppercase tracking-wider mb-2">
            Langkah 4: Konfirmasi Widget
          </label>
          
          <div class="p-4 rounded-xl border border-white/5 bg-[#161616] space-y-3">
            <h4 class="text-xs font-bold text-zinc-400">Ringkasan Konfigurasi:</h4>
            <div class="grid grid-cols-2 gap-y-2 gap-x-4 text-xs">
              <span class="text-text-secondary">Sumber Data:</span>
              <span class="text-white font-semibold capitalize">{{ selectedSource }}</span>
              
              <span class="text-text-secondary">Item/Kota:</span>
              <span class="text-white font-semibold">{{ selectedMetricItem }}</span>
              
              <span class="text-text-secondary">Tipe Grafik:</span>
              <span class="text-white font-semibold capitalize">{{ selectedChartType }} Chart</span>
            </div>
          </div>

          <div>
            <label class="block text-[10px] font-bold text-text-secondary uppercase tracking-wider mb-1.5">Nama Widget (Editable)</label>
            <input 
              type="text" 
              v-model="widgetTitle" 
              placeholder="Masukkan nama widget..." 
              class="w-full text-xs rounded-lg border border-white/8 bg-[#161616] px-3.5 py-2.5 text-white focus:border-blue-500 focus:outline-none transition-colors"
            />
          </div>
        </div>

      </div>

      <!-- Modal Footer -->
      <div class="flex items-center justify-between border-t border-white/5 pt-4 mt-4 bg-[#111111]">
        <button
          v-if="currentStep > 1"
          @click="goToStep(currentStep - 1)"
          class="px-4 py-2 rounded-lg border border-white/10 hover:border-white/20 text-zinc-400 hover:text-white transition-all text-xs font-bold cursor-pointer"
        >
          ← Kembali
        </button>
        <div v-else></div> <!-- Placeholder for layout -->

        <button
          v-if="currentStep < 4"
          @click="goToStep(currentStep + 1)"
          :disabled="currentStep === 1 ? !selectedSource : currentStep === 2 ? (!selectedMetricItem || !selectedMetricField) : !selectedChartType"
          class="px-5 py-2 rounded-lg bg-blue-500 hover:bg-blue-400 text-white font-bold transition-all text-xs cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed shadow-md shadow-blue-500/10 hover:shadow-blue-500/20"
        >
          Lanjut →
        </button>
        <button
          v-else
          @click="addWidget"
          class="px-5 py-2 rounded-lg bg-green-600 hover:bg-green-500 text-white font-bold transition-all text-xs cursor-pointer shadow-md shadow-green-600/10 hover:shadow-green-600/20"
        >
          ➕ Tambahkan ke Canvas
        </button>
      </div>

    </div>
  </div>
</template>
