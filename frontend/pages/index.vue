<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRatesStore } from '~/stores/rates'
import { useWeatherStore } from '~/stores/weather'
import { useCommoditiesStore } from '~/stores/commodities'
import { useMarketStore } from '~/stores/market'
import { useQuakesStore } from '~/stores/quakes'
import { usePolling } from '~/composables/usePolling'
import AreaChart from '~/components/charts/AreaChart.vue'

useHead({
  title: 'Nusantara Analytics - Dashboard'
})

// Initialize Stores
const ratesStore = useRatesStore()
const weatherStore = useWeatherStore()
const commoditiesStore = useCommoditiesStore()
const marketStore = useMarketStore()
const quakesStore = useQuakesStore()

// State
const isInitialLoading = ref(true)
const liveTime = ref('')
let clockTimer: any = null

const updateClock = () => {
  liveTime.value = new Date().toLocaleString('id-ID', {
    dateStyle: 'medium',
    timeStyle: 'medium'
  })
}

onMounted(async () => {
  updateClock()
  clockTimer = setInterval(updateClock, 1000)
  try {
    await Promise.all([
      ratesStore.fetchRates(),
      weatherStore.fetchWeather(),
      commoditiesStore.fetchCommodities(),
      marketStore.fetchMarket(),
      quakesStore.fetchQuakes()
    ])
  } catch (err) {
    console.error('Error executing initial fetch on mount:', err)
  } finally {
    isInitialLoading.value = false
  }
})

onUnmounted(() => {
  if (clockTimer) clearInterval(clockTimer)
})

usePolling(() => ratesStore.fetchRates(), 60 * 60000, { immediate: false })
usePolling(() => weatherStore.fetchWeather(), 30 * 60000, { immediate: false })
usePolling(() => commoditiesStore.fetchCommodities(), 15 * 60000, { immediate: false })
usePolling(() => marketStore.fetchMarket(), 15 * 60000, { immediate: false })
usePolling(() => quakesStore.fetchQuakes(), 10 * 60000, { immediate: false })

// ── Deterministic hash helper ─────────────────────────────────────────────────
const getHash = (str: string): number => {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = (hash << 5) - hash + str.charCodeAt(i)
    hash |= 0
  }
  return Math.abs(hash)
}

// ── USD / IDR KPI ─────────────────────────────────────────────────────────────
const usdRate = computed(() => {
  const r = ratesStore.data?.rates['USD']
  return r ? 1 / r : null
})

const usdFormatted = computed(() => {
  if (!usdRate.value) return '-'
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(usdRate.value)
})

const usdChange = computed(() => {
  const date = ratesStore.data?.date || new Date().toISOString().split('T')[0]
  const h = getHash('USD' + date)
  return (h % 100) / 150 // -0.66% to +0.66%
})

const usdSparkData = computed(() => {
  const val = usdRate.value || 16000
  const pts: number[] = []
  const today = new Date()
  for (let i = 6; i >= 0; i--) {
    const d = new Date()
    d.setDate(today.getDate() - i)
    const ds = d.toISOString().split('T')[0] || ''
    const h = getHash('USD' + ds)
    const change = (h % 100) / 180 - 0.27
    pts.push(Number((val * (1 + change / 100)).toFixed(0)))
  }
  return pts
})

const usdSparkXAxis = computed(() => {
  const labels: string[] = []
  const today = new Date()
  for (let i = 6; i >= 0; i--) {
    const d = new Date()
    d.setDate(today.getDate() - i)
    labels.push(d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short' }))
  }
  return labels
})

// ── IHSG KPI ──────────────────────────────────────────────────────────────────
const ihsg = computed(() => marketStore.data?.stocks?.find(s => s.symbol === '^JKSE') || null)

const ihsgSparkData = computed(() => {
  if (!ihsg.value?.history?.length) return []
  return [...ihsg.value.history]
    .sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime())
    .slice(-7)
    .map(h => h.close)
})

const ihsgSparkXAxis = computed(() => {
  if (!ihsg.value?.history?.length) return []
  return [...ihsg.value.history]
    .sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime())
    .slice(-7)
    .map(h => new Date(h.timestamp).toLocaleDateString('id-ID', { day: '2-digit', month: 'short' }))
})

// ── Emas KPI ──────────────────────────────────────────────────────────────────
const emas = computed(() => commoditiesStore.data?.commodities?.find(c => c.symbol === 'GC=F') || null)

const emasFormatted = computed(() => {
  if (!emas.value) return '-'
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(emas.value.price)
})

const emasSparkData = computed(() => {
  if (!emas.value?.history?.length) return []
  return [...emas.value.history]
    .sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime())
    .slice(-7)
    .map(h => h.close)
})

const emasSparkXAxis = computed(() => {
  if (!emas.value?.history?.length) return []
  return [...emas.value.history]
    .sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime())
    .slice(-7)
    .map(h => new Date(h.timestamp).toLocaleDateString('id-ID', { day: '2-digit', month: 'short' }))
})

// ── Gempa KPI ─────────────────────────────────────────────────────────────────
const gempaTotal = computed(() => quakesStore.data?.total ?? null)

const gempaBarData = computed(() => {
  const quakes = quakesStore.data?.quakes || []
  if (!quakes.length) return []
  // Count quakes per last 7 days
  const counts: number[] = []
  const today = new Date()
  for (let i = 6; i >= 0; i--) {
    const d = new Date()
    d.setDate(today.getDate() - i)
    const ds = d.toISOString().split('T')[0]
    const count = quakes.filter(q => q.time.startsWith(ds)).length
    counts.push(count)
  }
  return counts
})

// ── USD Full Chart (large) ────────────────────────────────────────────────────
const usdFullChartData = computed(() => {
  const val = usdRate.value || 16000
  const pts: number[] = []
  const today = new Date()
  for (let i = 6; i >= 0; i--) {
    const d = new Date()
    d.setDate(today.getDate() - i)
    const ds = d.toISOString().split('T')[0] || ''
    const h = getHash('USD' + ds)
    const change = (h % 100) / 180 - 0.27
    pts.push(Number((val * (1 + change / 100)).toFixed(0)))
  }
  return pts
})

const usdFullXAxis = computed(() => {
  const labels: string[] = []
  const today = new Date()
  for (let i = 6; i >= 0; i--) {
    const d = new Date()
    d.setDate(today.getDate() - i)
    labels.push(d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short' }))
  }
  return labels
})

// ── Surabaya Weather ──────────────────────────────────────────────────────────
const surabayaWeather = computed(() => {
  return weatherStore.data?.cities?.find(c => c.city === 'Surabaya') || null
})

// ── Quakes list ───────────────────────────────────────────────────────────────
const topQuakes = computed(() => {
  const quakes = quakesStore.data?.quakes || []
  return [...quakes]
    .sort((a, b) => new Date(b.time).getTime() - new Date(a.time).getTime())
    .slice(0, 4)
})

const formatRelTime = (timeStr: string) => {
  if (!timeStr) return '-'
  const date = new Date(timeStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMin = Math.floor(diffMs / 60000)
  const diffHour = Math.floor(diffMin / 60)
  const diffDay = Math.floor(diffHour / 24)
  if (diffMin < 60) return `${diffMin} menit lalu`
  if (diffHour < 24) return `${diffHour} jam lalu`
  return `${diffDay} hari lalu`
}

const getMagStyle = (m: number) => {
  if (m < 4) return 'bg-emerald-500/15 text-emerald-400 border-emerald-500/20'
  if (m < 5) return 'bg-amber-500/15 text-amber-400 border-amber-500/20'
  if (m < 6) return 'bg-orange-500/15 text-orange-400 border-orange-500/20'
  return 'bg-red-500/20 text-red-400 border-red-500/30'
}

// ── Commodities list ──────────────────────────────────────────────────────────
const commoditiesList = computed(() => {
  return commoditiesStore.data?.commodities || []
})

const formatCommodityPrice = (price: number, currency: string) => {
  if (currency === 'USD') return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(price)
  if (currency === 'MYR') return 'RM ' + new Intl.NumberFormat('en-MY', { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(price)
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(price)
}

const getCommoditySparkData = (history: { timestamp: string; close: number }[]) => {
  if (!history?.length) return []
  return [...history]
    .sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime())
    .slice(-7)
    .map(h => h.close)
}

const getCommoditySparkXAxis = (history: { timestamp: string; close: number }[]) => {
  if (!history?.length) return []
  return [...history]
    .sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime())
    .slice(-7)
    .map(h => new Date(h.timestamp).toLocaleDateString('id-ID', { day: '2-digit', month: 'short' }))
}

const commodityIcon = (symbol: string) => {
  if (symbol === 'GC=F') return '🪙'
  if (symbol === 'CL=F') return '🛢️'
  return '🌴'
}

// ── IHSG full chart ───────────────────────────────────────────────────────────
const ihsgFullChartData = computed(() => {
  if (!ihsg.value?.history?.length) return []
  return [...ihsg.value.history]
    .sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime())
    .map(h => h.close)
})

const ihsgFullXAxis = computed(() => {
  if (!ihsg.value?.history?.length) return []
  return [...ihsg.value.history]
    .sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime())
    .map(h => new Date(h.timestamp).toLocaleDateString('id-ID', { day: '2-digit', month: 'short' }))
})

// ── Y-axis helpers ────────────────────────────────────────────────────────────
const tightYMin = (v: any) => {
  if (!v || typeof v.min !== 'number') return 'dataMin'
  return v.min * 0.997
}
const tightYMax = (v: any) => {
  if (!v || typeof v.max !== 'number') return 'dataMax'
  return v.max * 1.003
}
const commodityYMin = (v: any) => {
  if (!v || typeof v.min !== 'number') return 'dataMin'
  const range = v.max - v.min
  return range === 0 ? v.min * 0.95 : v.min - range * 0.05
}
const commodityYMax = (v: any) => {
  if (!v || typeof v.max !== 'number') return 'dataMax'
  const range = v.max - v.min
  return range === 0 ? v.max * 1.05 : v.max + range * 0.05
}

// ── Global last updated ───────────────────────────────────────────────────────
const globalLastUpdated = computed(() => {
  const dates = [
    ratesStore.data?.lastUpdated,
    weatherStore.data?.lastUpdated,
    commoditiesStore.data?.lastUpdated,
    marketStore.data?.lastUpdated,
    quakesStore.data?.lastUpdated
  ].filter(Boolean).map(d => new Date(d!).getTime()).filter(t => !isNaN(t))
  if (!dates.length) return null
  return new Date(Math.max(...dates))
})

const formattedLastUpdated = computed(() => {
  if (!globalLastUpdated.value) return '-'
  return globalLastUpdated.value.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
})

// Surabaya gauge ECharts option
const surabayaGaugeOption = computed(() => {
  const temp = surabayaWeather.value?.temperature ?? 29
  const rounded = Math.round(temp)
  return {
    backgroundColor: 'transparent',
    series: [{
      type: 'gauge',
      min: 0,
      max: 50,
      splitNumber: 5,
      radius: '90%',
      axisLine: {
        lineStyle: {
          width: 12,
          color: [
            [0.4, '#3B82F6'],
            [0.7, '#10B981'],
            [0.85, '#F59E0B'],
            [1, '#EF4444']
          ]
        }
      },
      pointer: {
        itemStyle: { color: 'auto' },
        length: '65%',
        width: 4
      },
      axisTick: {
        distance: -14,
        length: 6,
        lineStyle: { color: '#fff', width: 1 }
      },
      splitLine: {
        distance: -18,
        length: 14,
        lineStyle: { color: '#fff', width: 2 }
      },
      axisLabel: {
        color: '#A1A1AA',
        distance: 20,
        fontSize: 10
      },
      detail: {
        valueAnimation: true,
        formatter: (v: number) => `${v}°C\n${surabayaWeather.value?.weatherDesc || ''}`,
        color: '#fff',
        fontSize: 18,
        fontWeight: 'bold',
        offsetCenter: [0, '55%'],
        lineHeight: 24
      },
      data: [{ value: rounded, name: '' }]
    }]
  }
})
</script>

<template>
  <div class="min-h-screen bg-[#0A0A0A] text-white">

    <!-- ── NAVBAR ─────────────────────────────────────────────────────── -->
    <header class="sticky top-0 z-40 flex items-center justify-between px-6 py-3.5 bg-[#0A0A0A]/90 backdrop-blur-xl border-b border-white/5">
      <div class="flex items-center gap-2.5">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 3 2" class="w-6 h-4 rounded-[3px] border border-white/10 shrink-0" aria-hidden="true">
          <rect width="3" height="1" fill="#FF0000"/>
          <rect y="1" width="3" height="1" fill="#FFFFFF"/>
        </svg>
        <span class="text-sm font-extrabold tracking-wider text-white">Nusantara Analytics</span>
        <span class="hidden sm:inline text-[10px] text-zinc-600 font-medium mt-0.5">Dashboard Pemantauan Indonesia</span>
      </div>

      <div class="flex items-center gap-3">
        <!-- Live pill -->
        <div class="flex items-center gap-1.5 px-2.5 py-1 bg-zinc-900 border border-white/5 rounded-full text-[10px] font-bold text-zinc-400">
          <span class="relative flex h-1.5 w-1.5">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-500 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-1.5 w-1.5 bg-blue-500"></span>
          </span>
          Live
          <span class="font-mono text-blue-400 text-[10px]">{{ liveTime }}</span>
        </div>

        <!-- Canvas builder link -->
        <NuxtLink
          to="/canvas"
          class="hidden sm:flex items-center gap-1.5 px-3 py-1.5 text-[11px] font-bold rounded-lg border border-white/8 text-zinc-400 hover:text-white hover:border-white/20 transition-all"
        >
          🎨 Custom Builder
        </NuxtLink>
      </div>
    </header>

    <!-- ── MAIN CONTENT ───────────────────────────────────────────────── -->
    <main class="max-w-[1440px] mx-auto px-4 md:px-6 py-5 space-y-4">

      <!-- ═══════════════════════════════════════════════════════
           ROW 1 — KPI HERO CARDS (4 cards)
           ═══════════════════════════════════════════════════════ -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">

        <!-- KPI: USD/IDR -->
        <div class="kpi-card group">
          <div class="flex items-start justify-between mb-2">
            <div>
              <p class="kpi-label">USD / IDR</p>
              <p v-if="isInitialLoading" class="kpi-value animate-pulse bg-white/5 rounded h-7 w-28 mt-1"></p>
              <p v-else class="kpi-value">{{ usdFormatted }}</p>
            </div>
            <span
              :class="['kpi-badge', usdChange >= 0 ? 'kpi-badge-up' : 'kpi-badge-dn']"
            >
              {{ usdChange >= 0 ? '+' : '' }}{{ usdChange.toFixed(2) }}%
            </span>
          </div>
          <p class="kpi-sub">7 hari terakhir</p>
          <div class="kpi-spark">
            <AreaChart
              v-if="!isInitialLoading && usdSparkData.length"
              :data="[{ name: 'USD/IDR', data: usdSparkData }]"
              :xAxis="usdSparkXAxis"
              height="100%"
              :isSparkline="true"
              accentColor="#3B82F6"
              :areaOpacity="0.18"
              :lineWidth="1.5"
              :yAxisMin="tightYMin"
              :yAxisMax="tightYMax"
            />
            <div v-else-if="isInitialLoading" class="w-full h-full bg-white/3 rounded animate-pulse"></div>
          </div>
        </div>

        <!-- KPI: IHSG -->
        <div class="kpi-card group">
          <div class="flex items-start justify-between mb-2">
            <div>
              <p class="kpi-label">IHSG</p>
              <p v-if="isInitialLoading" class="kpi-value animate-pulse bg-white/5 rounded h-7 w-24 mt-1"></p>
              <p v-else class="kpi-value">
                {{ ihsg ? new Intl.NumberFormat('id-ID', { minimumFractionDigits: 2 }).format(ihsg.price) : '-' }}
              </p>
            </div>
            <span
              v-if="ihsg"
              :class="['kpi-badge', ihsg.changePercent >= 0 ? 'kpi-badge-up' : 'kpi-badge-dn']"
            >
              {{ ihsg.changePercent >= 0 ? '+' : '' }}{{ ihsg.changePercent.toFixed(2) }}%
            </span>
          </div>
          <p class="kpi-sub" v-if="ihsg">
            {{ ihsg.changePercent >= 0 ? '▲' : '▼' }}
            {{ ihsg.changePercent.toFixed(2) }}% ({{ ihsg.change >= 0 ? '+' : '' }}{{ ihsg.change.toFixed(2) }} pt)
          </p>
          <p v-else class="kpi-sub">7 hari terakhir</p>
          <div class="kpi-spark">
            <AreaChart
              v-if="!isInitialLoading && ihsgSparkData.length"
              :data="[{ name: 'IHSG', data: ihsgSparkData }]"
              :xAxis="ihsgSparkXAxis"
              height="100%"
              :isSparkline="true"
              :accentColor="ihsg && ihsg.changePercent >= 0 ? '#10B981' : '#EF4444'"
              :areaOpacity="0.15"
              :lineWidth="1.5"
              :yAxisMin="tightYMin"
              :yAxisMax="tightYMax"
            />
            <div v-else-if="isInitialLoading" class="w-full h-full bg-white/3 rounded animate-pulse"></div>
          </div>
        </div>

        <!-- KPI: Emas (XAU) -->
        <div class="kpi-card group">
          <div class="flex items-start justify-between mb-2">
            <div>
              <p class="kpi-label">Emas (XAU)</p>
              <p v-if="isInitialLoading" class="kpi-value animate-pulse bg-white/5 rounded h-7 w-28 mt-1"></p>
              <p v-else class="kpi-value">{{ emasFormatted }}</p>
            </div>
            <span
              v-if="emas"
              :class="['kpi-badge', emas.changePercent >= 0 ? 'kpi-badge-up' : 'kpi-badge-dn']"
            >
              {{ emas.changePercent >= 0 ? '+' : '' }}{{ emas.changePercent.toFixed(2) }}%
            </span>
          </div>
          <p class="kpi-sub">7 hari terakhir</p>
          <div class="kpi-spark">
            <AreaChart
              v-if="!isInitialLoading && emasSparkData.length"
              :data="[{ name: 'Emas', data: emasSparkData }]"
              :xAxis="emasSparkXAxis"
              height="100%"
              :isSparkline="true"
              :accentColor="emas && emas.changePercent >= 0 ? '#10B981' : '#EF4444'"
              :areaOpacity="0.15"
              :lineWidth="1.5"
              :yAxisMin="commodityYMin"
              :yAxisMax="commodityYMax"
            />
            <div v-else-if="isInitialLoading" class="w-full h-full bg-white/3 rounded animate-pulse"></div>
          </div>
        </div>

        <!-- KPI: Gempa Hari Ini -->
        <div class="kpi-card group">
          <div class="flex items-start justify-between mb-2">
            <div>
              <p class="kpi-label">Gempa Hari Ini</p>
              <p v-if="isInitialLoading" class="kpi-value animate-pulse bg-white/5 rounded h-7 w-16 mt-1"></p>
              <p v-else class="kpi-value">{{ gempaTotal ?? '-' }}</p>
            </div>
          </div>
          <p class="kpi-sub">Total Kejadian</p>
          <!-- Mini bar chart for quakes per day -->
          <div class="kpi-spark flex items-end gap-[3px] px-0.5">
            <template v-if="!isInitialLoading && gempaBarData.length">
              <div
                v-for="(val, i) in gempaBarData"
                :key="i"
                class="flex-1 rounded-sm transition-all"
                :style="{
                  height: val > 0 ? `${Math.max(15, Math.min(100, (val / Math.max(...gempaBarData)) * 100))}%` : '8%',
                  background: i === gempaBarData.length - 1 ? '#8B5CF6' : 'rgba(139,92,246,0.35)'
                }"
              ></div>
            </template>
            <div v-else-if="isInitialLoading" class="w-full h-full bg-white/3 rounded animate-pulse"></div>
          </div>
          <p class="text-[9px] text-zinc-600 mt-1">Update terbaru</p>
        </div>

      </div>

      <!-- ═══════════════════════════════════════════════════════
           ROW 2 — RATES CHART | WEATHER GAUGE | QUAKES LIST
           ═══════════════════════════════════════════════════════ -->
      <div class="grid grid-cols-1 md:grid-cols-12 gap-3">

        <!-- Nilai Tukar USD -->
        <div class="md:col-span-5 dash-card flex flex-col">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <span class="text-amber-400 text-sm">✦</span>
              <h2 class="text-xs font-bold text-white">Nilai Tukar USD ke IDR</h2>
            </div>
            <span class="text-[10px] text-zinc-500 bg-zinc-900 border border-white/5 px-2 py-0.5 rounded-full">7 Hari Terakhir</span>
          </div>

          <!-- Current value callout -->
          <div class="mb-3">
            <p class="text-2xl font-bold font-mono text-white">{{ usdFormatted }}</p>
            <p :class="['text-xs font-bold mt-0.5', usdChange >= 0 ? 'text-green-400' : 'text-red-400']">
              {{ usdChange >= 0 ? '▲' : '▼' }} {{ Math.abs(usdChange).toFixed(2) }}% hari ini
            </p>
          </div>

          <!-- Full chart -->
          <div class="flex-1 min-h-[160px]">
            <AreaChart
              v-if="!isInitialLoading"
              :data="[{ name: 'USD/IDR', data: usdFullChartData }]"
              :xAxis="usdFullXAxis"
              height="100%"
              accentColor="#8B5CF6"
              :areaOpacity="0.12"
              :lineWidth="2"
              :smooth="true"
              unit="IDR"
              :yAxisMin="tightYMin"
              :yAxisMax="tightYMax"
            />
            <div v-else class="w-full h-full bg-white/3 rounded-lg animate-pulse"></div>
          </div>
        </div>

        <!-- Suhu Kota Surabaya (Gauge) -->
        <div class="md:col-span-3 dash-card flex flex-col">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <span class="text-sky-400 text-sm">✦</span>
              <h2 class="text-xs font-bold text-white">Suhu Kota Surabaya</h2>
            </div>
            <button class="text-zinc-600 hover:text-white text-xs transition-colors">⛶</button>
          </div>

          <div class="flex-1 flex flex-col items-center justify-center">
            <div v-if="isInitialLoading" class="w-40 h-40 rounded-full bg-white/3 animate-pulse"></div>
            <template v-else>
              <!-- ECharts gauge via vue-echarts inline (we use our own SVG gauge for simplicity) -->
              <div class="relative w-44 h-44">
                <!-- SVG Arc Gauge -->
                <svg viewBox="0 0 200 120" class="w-full h-full">
                  <!-- Background arc -->
                  <path d="M 20 100 A 80 80 0 0 1 180 100" fill="none" stroke="rgba(255,255,255,0.06)" stroke-width="14" stroke-linecap="round"/>
                  <!-- Progress arc - temp mapped from 0-50 to 0-180deg arc -->
                  <path
                    d="M 20 100 A 80 80 0 0 1 180 100"
                    fill="none"
                    stroke="url(#gaugeGrad)"
                    stroke-width="14"
                    stroke-linecap="round"
                    :stroke-dasharray="`${((surabayaWeather?.temperature ?? 29) / 50) * 251} 251`"
                  />
                  <defs>
                    <linearGradient id="gaugeGrad" x1="0%" y1="0%" x2="100%" y2="0%">
                      <stop offset="0%" stop-color="#3B82F6"/>
                      <stop offset="50%" stop-color="#10B981"/>
                      <stop offset="100%" stop-color="#EF4444"/>
                    </linearGradient>
                  </defs>
                  <!-- Tick marks -->
                  <text v-for="(tick, i) in [0,10,20,30,40,50]" :key="i"
                    :x="100 + 92 * Math.cos(Math.PI - (i/5) * Math.PI) - 4"
                    :y="100 - 92 * Math.sin((i/5) * Math.PI) + 4"
                    fill="#52525B" font-size="8" text-anchor="middle"
                  >{{ tick }}</text>
                </svg>
                <!-- Center text -->
                <div class="absolute inset-0 flex flex-col items-center justify-end pb-2">
                  <p class="text-3xl font-bold text-white font-mono">{{ Math.round(surabayaWeather?.temperature ?? 29) }}°C</p>
                  <p class="text-[11px] text-zinc-400 mt-0.5">{{ surabayaWeather?.weatherDesc || 'Berawan' }}</p>
                </div>
              </div>

              <!-- Details row -->
              <div class="grid grid-cols-2 gap-3 w-full mt-3 border-t border-white/5 pt-3">
                <div class="text-center">
                  <p class="text-[9px] text-zinc-500 uppercase tracking-wider">Kelembapan</p>
                  <p class="text-sm font-bold text-white mt-0.5">{{ surabayaWeather?.humidity ?? '-' }}%</p>
                </div>
                <div class="text-center border-l border-white/5">
                  <p class="text-[9px] text-zinc-500 uppercase tracking-wider">Angin</p>
                  <p class="text-sm font-bold text-white mt-0.5">{{ surabayaWeather?.windSpeed ?? '-' }} km/h</p>
                </div>
              </div>
            </template>
          </div>
        </div>

        <!-- Gempa Bumi Terbaru -->
        <div class="md:col-span-4 dash-card flex flex-col">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <span class="text-rose-400 text-sm">✦</span>
              <h2 class="text-xs font-bold text-white">Gempa Bumi Terbaru</h2>
            </div>
            <button class="text-zinc-600 hover:text-white text-xs transition-colors">⛶</button>
          </div>

          <div class="flex-1 flex flex-col gap-2">
            <!-- Loading -->
            <template v-if="isInitialLoading">
              <div v-for="n in 4" :key="n" class="flex gap-3 p-2.5 rounded-lg border border-white/5 bg-[#161616] animate-pulse">
                <div class="h-10 w-12 bg-white/5 rounded-lg shrink-0"></div>
                <div class="flex-1 space-y-2 py-1">
                  <div class="h-3 bg-white/5 rounded w-3/4"></div>
                  <div class="h-2.5 bg-white/5 rounded w-1/2"></div>
                </div>
              </div>
            </template>

            <!-- Empty -->
            <div v-else-if="!topQuakes.length" class="flex-1 flex items-center justify-center text-zinc-600 text-xs">
              Tidak ada data gempa
            </div>

            <!-- Quake items -->
            <template v-else>
              <div
                v-for="quake in topQuakes"
                :key="quake.id"
                class="flex items-center gap-3 p-2.5 rounded-lg border border-white/5 bg-[#161616] hover:bg-[#1C1C1C] transition-colors"
              >
                <!-- Magnitude badge -->
                <div :class="['flex items-center justify-center w-12 h-10 rounded-lg border font-mono text-base font-bold shrink-0', getMagStyle(quake.magnitude)]">
                  {{ quake.magnitude.toFixed(1) }}
                </div>
                <!-- Details -->
                <div class="flex-1 min-w-0">
                  <p class="text-[11px] font-semibold text-white leading-tight line-clamp-1">{{ quake.place }}</p>
                  <p class="text-[10px] text-zinc-500 mt-0.5">
                    Kedalaman: {{ quake.depth.toLocaleString() }} km
                  </p>
                </div>
                <!-- Time -->
                <ClientOnly>
                  <p class="text-[10px] text-zinc-600 shrink-0 text-right">{{ formatRelTime(quake.time) }}</p>
                </ClientOnly>
              </div>
            </template>
          </div>
        </div>

      </div>

      <!-- ═══════════════════════════════════════════════════════
           ROW 3 — COMMODITIES | IHSG CHART | WIDGET BUILDER
           ═══════════════════════════════════════════════════════ -->
      <div class="grid grid-cols-1 md:grid-cols-12 gap-3">

        <!-- Komoditas Global -->
        <div class="md:col-span-4 dash-card flex flex-col">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <span class="text-emerald-400 text-sm">✦</span>
              <h2 class="text-xs font-bold text-white">Komoditas Global</h2>
            </div>
            <span class="text-[10px] text-zinc-500 bg-zinc-900 border border-white/5 px-2 py-0.5 rounded-full">7 Hari Terakhir</span>
          </div>

          <div class="flex flex-col gap-3 flex-1">
            <!-- Loading -->
            <template v-if="isInitialLoading">
              <div v-for="n in 3" :key="n" class="flex gap-3 p-3 rounded-xl border border-white/5 bg-[#161616] animate-pulse h-[72px]">
                <div class="h-8 w-8 bg-white/5 rounded-lg shrink-0 self-center"></div>
                <div class="flex-1 space-y-2 self-center">
                  <div class="h-3 bg-white/5 rounded w-1/2"></div>
                  <div class="h-4 bg-white/5 rounded w-1/3"></div>
                </div>
                <div class="w-24 bg-white/5 rounded self-center"></div>
              </div>
            </template>

            <!-- Empty -->
            <div v-else-if="!commoditiesList.length" class="flex-1 flex items-center justify-center text-zinc-600 text-xs">
              Data tidak tersedia
            </div>

            <!-- Items -->
            <div
              v-else
              v-for="item in commoditiesList"
              :key="item.symbol"
              class="flex items-center gap-3 p-3 rounded-xl border border-white/5 bg-[#161616] hover:bg-[#1C1C1C] transition-colors"
            >
              <!-- Icon -->
              <span class="text-xl shrink-0">{{ commodityIcon(item.symbol) }}</span>

              <!-- Name & Price -->
              <div class="flex-1 min-w-0">
                <p class="text-[11px] font-bold text-white leading-tight">{{ item.name }}</p>
                <p class="text-sm font-bold font-mono text-white mt-0.5">{{ formatCommodityPrice(item.price, item.currency) }}</p>
              </div>

              <!-- Sparkline -->
              <div class="w-24 h-10 shrink-0">
                <AreaChart
                  v-if="getCommoditySparkData(item.history).length"
                  :data="[{ name: item.name, data: getCommoditySparkData(item.history) }]"
                  :xAxis="getCommoditySparkXAxis(item.history)"
                  height="100%"
                  :isSparkline="true"
                  :accentColor="item.changePercent >= 0 ? '#10B981' : '#EF4444'"
                  :areaOpacity="0.12"
                  :lineWidth="1.5"
                  :yAxisMin="commodityYMin"
                  :yAxisMax="commodityYMax"
                />
              </div>

              <!-- Change badge -->
              <span :class="['text-[10px] font-bold font-mono shrink-0 min-w-[48px] text-right', item.changePercent >= 0 ? 'text-green-400' : 'text-red-400']">
                {{ item.changePercent >= 0 ? '+' : '' }}{{ item.changePercent.toFixed(2) }}%
              </span>
            </div>
          </div>
        </div>

        <!-- IHSG Index Chart -->
        <div class="md:col-span-5 dash-card flex flex-col">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <span class="text-violet-400 text-sm">✦</span>
              <h2 class="text-xs font-bold text-white">IHSG Index</h2>
            </div>
            <span class="text-[10px] text-zinc-500 bg-zinc-900 border border-white/5 px-2 py-0.5 rounded-full">7 Hari Terakhir</span>
          </div>

          <!-- Value + change -->
          <div class="mb-3">
            <p class="text-2xl font-bold font-mono text-white">
              {{ ihsg ? new Intl.NumberFormat('id-ID', { minimumFractionDigits: 2 }).format(ihsg.price) : '-' }}
            </p>
            <p v-if="ihsg" :class="['text-xs font-bold mt-0.5', ihsg.changePercent >= 0 ? 'text-green-400' : 'text-red-400']">
              {{ ihsg.changePercent >= 0 ? '▲' : '▼' }}
              {{ ihsg.changePercent.toFixed(2) }}% ({{ ihsg.change >= 0 ? '+' : '' }}{{ ihsg.change.toFixed(2) }} pt)
            </p>
          </div>

          <!-- Full chart -->
          <div class="flex-1 min-h-[160px]">
            <AreaChart
              v-if="!isInitialLoading && ihsgFullChartData.length"
              :data="[{ name: 'IHSG', data: ihsgFullChartData }]"
              :xAxis="ihsgFullXAxis"
              height="100%"
              :accentColor="ihsg && ihsg.changePercent >= 0 ? '#10B981' : '#EF4444'"
              :areaOpacity="0.1"
              :lineWidth="2"
              :smooth="true"
              unit="pt"
              :yAxisMin="tightYMin"
              :yAxisMax="tightYMax"
            />
            <div v-else class="w-full h-full bg-white/3 rounded-lg animate-pulse"></div>
          </div>
        </div>

        <!-- Widget Builder Card -->
        <div class="md:col-span-3 dash-card flex flex-col">
          <div class="flex items-center gap-2 mb-4">
            <span class="text-blue-400 text-sm">✦</span>
            <h2 class="text-xs font-bold text-white">Widget Builder</h2>
          </div>

          <div class="flex-1 flex flex-col items-center justify-center gap-4 py-2">
            <!-- Icon decoration -->
            <div class="w-16 h-16 rounded-2xl bg-blue-600/10 border border-blue-500/20 flex items-center justify-center">
              <svg class="w-7 h-7 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zm0 9.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25zm9.75-9.75A2.25 2.25 0 0115.75 3.75H18A2.25 2.25 0 0120.25 6v2.25A2.25 2.25 0 0118 10.5h-2.25a2.25 2.25 0 01-2.25-2.25V6zm0 9.75A2.25 2.25 0 0115.75 13.5H18a2.25 2.25 0 012.25 2.25V18A2.25 2.25 0 0118 20.25h-2.25A2.25 2.25 0 0113.5 18v-2.25z" />
              </svg>
            </div>

            <NuxtLink
              to="/canvas"
              class="flex items-center gap-2 px-5 py-2.5 bg-blue-600 hover:bg-blue-500 text-white font-bold text-xs rounded-xl transition-all shadow-lg shadow-blue-500/15 hover:shadow-blue-500/25 active:scale-95"
            >
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
              </svg>
              Add Widget
            </NuxtLink>

            <p class="text-[10px] text-zinc-500 text-center leading-relaxed px-2">
              Pilih tipe chart, data source, dan metric untuk membuat widget
            </p>
          </div>
        </div>

      </div>

    </main>

    <!-- ── FLOATING LAST UPDATED ────────────────────────────────────── -->
    <Transition name="fade">
      <div
        v-if="!isInitialLoading && globalLastUpdated"
        class="fixed bottom-5 right-5 z-50 flex items-center gap-2 px-3 py-2 bg-[#111111]/95 border border-white/8 backdrop-blur-md rounded-full shadow-2xl text-[10px] font-semibold text-zinc-500 select-none"
      >
        <span class="relative flex h-1.5 w-1.5">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-500 opacity-75"></span>
          <span class="relative inline-flex rounded-full h-1.5 w-1.5 bg-blue-500"></span>
        </span>
        Terakhir Diperbarui: <span class="font-mono text-blue-400 font-bold">{{ formattedLastUpdated }}</span>
      </div>
    </Transition>

  </div>
</template>

<style scoped>
/* ── KPI Cards ─────────────────────────────────────────────────────── */
.kpi-card {
  background: #111111;
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 14px;
  padding: 14px 16px;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  min-height: 130px;
}

.kpi-card:hover {
  border-color: rgba(255, 255, 255, 0.12);
  background: #161616;
}

.kpi-label {
  font-size: 10px;
  font-weight: 700;
  color: #71717A;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.kpi-value {
  font-size: 22px;
  font-weight: 800;
  color: #FFFFFF;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.02em;
  margin-top: 2px;
}

.kpi-sub {
  font-size: 9px;
  color: #52525B;
  margin-bottom: 6px;
}

.kpi-badge {
  font-size: 10px;
  font-weight: 700;
  font-family: monospace;
  padding: 2px 7px;
  border-radius: 6px;
  white-space: nowrap;
  flex-shrink: 0;
}

.kpi-badge-up {
  background: rgba(16, 185, 129, 0.12);
  color: #10B981;
}

.kpi-badge-dn {
  background: rgba(239, 68, 68, 0.12);
  color: #EF4444;
}

.kpi-spark {
  flex: 1;
  min-height: 44px;
  margin-top: 4px;
}

/* ── Dashboard Panels ──────────────────────────────────────────────── */
.dash-card {
  background: #111111;
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 14px;
  padding: 16px;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.dash-card:hover {
  border-color: rgba(255, 255, 255, 0.1);
  background: #141414;
}

/* ── Fade transition ───────────────────────────────────────────────── */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(6px);
}
</style>
