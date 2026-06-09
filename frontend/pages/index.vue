<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRatesStore } from '~/stores/rates'
import { useWeatherStore } from '~/stores/weather'
import { useCommoditiesStore } from '~/stores/commodities'
import { useMarketStore } from '~/stores/market'
import { useQuakesStore } from '~/stores/quakes'
import { usePolling } from '~/composables/usePolling'

import RatesWidget from '~/components/widgets/RatesWidget.vue'
import WeatherWidget from '~/components/widgets/WeatherWidget.vue'
import CommoditiesWidget from '~/components/widgets/CommoditiesWidget.vue'
import MarketWidget from '~/components/widgets/MarketWidget.vue'
import QuakesWidget from '~/components/widgets/QuakesWidget.vue'

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

// Live Clock in Header
const updateClock = () => {
  liveTime.value = new Date().toLocaleString('id-ID', {
    dateStyle: 'medium',
    timeStyle: 'medium'
  })
}

// Fetch all stores in parallel on mounted using Promise.all
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

// Setup Page-Level Polling with usePolling (immediate: false so they don't fetch again on load)
// rates: 60 minutes
usePolling(() => ratesStore.fetchRates(), 60 * 60000, { immediate: false })

// weather: 30 minutes
usePolling(() => weatherStore.fetchWeather(), 30 * 60000, { immediate: false })

// commodities & market: 15 minutes
usePolling(() => commoditiesStore.fetchCommodities(), 15 * 60000, { immediate: false })
usePolling(() => marketStore.fetchMarket(), 15 * 60000, { immediate: false })

// quakes: 10 minutes
usePolling(() => quakesStore.fetchQuakes(), 10 * 60000, { immediate: false })

// Global Last Updated Timestamp (maximum date among all stores)
const globalLastUpdated = computed(() => {
  const dates = [
    ratesStore.data?.lastUpdated ? new Date(ratesStore.data.lastUpdated).getTime() : null,
    weatherStore.data?.lastUpdated ? new Date(weatherStore.data.lastUpdated).getTime() : null,
    commoditiesStore.data?.lastUpdated ? new Date(commoditiesStore.data.lastUpdated).getTime() : null,
    marketStore.data?.lastUpdated ? new Date(marketStore.data.lastUpdated).getTime() : null,
    quakesStore.data?.lastUpdated ? new Date(quakesStore.data.lastUpdated).getTime() : null
  ].filter((t): t is number => t !== null && !isNaN(t))

  if (dates.length === 0) return null
  return new Date(Math.max(...dates))
})

const formattedLastUpdated = computed(() => {
  if (!globalLastUpdated.value) return '-'
  return globalLastUpdated.value.toLocaleTimeString('id-ID', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
})
</script>

<template>
  <div class="relative min-h-screen bg-slate-950 text-slate-100 p-4 sm:p-6 md:p-8">
    <!-- Premium Decorative background glow blobs wrapped to prevent page overflow -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none z-0">
      <div class="absolute -top-40 -left-40 w-96 h-96 bg-indigo-500/10 rounded-full blur-3xl"></div>
      <div class="absolute top-1/2 -right-40 w-96 h-96 bg-violet-500/5 rounded-full blur-3xl"></div>
      <div class="absolute -bottom-40 left-1/3 w-96 h-96 bg-emerald-500/5 rounded-full blur-3xl"></div>
    </div>

    <div class="max-w-7xl mx-auto space-y-8 relative z-10">
      
      <!-- Top Header Area -->
      <header class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-6 border-b border-slate-800/60">
        <div>
          <div class="flex items-center gap-2">
            <span class="text-3xl" aria-hidden="true">🇮🇩</span>
            <h1 class="text-2xl md:text-3xl font-extrabold tracking-tight text-white bg-gradient-to-r from-white via-slate-100 to-slate-400 bg-clip-text text-transparent">
              Nusantara Analytics
            </h1>
          </div>
          <p class="text-xs md:text-sm text-slate-400 mt-1.5 font-medium">
            Dashboard Terpadu Pemantauan Finansial, Cuaca, Komoditas Pasar, & Aktivitas Seismik Indonesia.
          </p>
        </div>

        <!-- System Navigation & Status -->
        <div class="flex flex-wrap items-center gap-3 self-start sm:self-auto">
          <!-- Link to Dedicated Builder -->
          <NuxtLink
            to="/builder"
            class="inline-flex items-center gap-2 px-4 py-2 text-xs font-bold rounded-xl border border-slate-800 bg-slate-900/40 hover:bg-slate-900 hover:border-slate-700 text-slate-200 transition-all shadow-sm"
          >
            <span>🎨</span> Custom Builder
          </NuxtLink>

          <div class="inline-flex items-center gap-3 bg-slate-900/60 border border-slate-850/80 backdrop-blur-md py-2 px-4 rounded-xl text-xs font-semibold text-slate-300 shadow-sm">
            <span class="relative flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            </span>
            <span class="tracking-wide">Sistem Online</span>
            <span class="text-slate-805" aria-hidden="true">|</span>
            <span class="font-mono text-slate-200">{{ liveTime }}</span>
          </div>
        </div>
      </header>

      <!-- Skeleton Grid (Initial Loading state) -->
      <div v-if="isInitialLoading" class="grid grid-cols-1 md:grid-cols-2 gap-6 animate-pulse">
        
        <!-- 1. Rates Skeleton -->
        <div class="relative w-full rounded-2xl border border-slate-900 bg-slate-900/20 p-6 flex flex-col gap-6">
          <div class="flex justify-between items-start">
            <div>
              <div class="h-5 bg-slate-800 rounded w-40"></div>
              <div class="h-3 bg-slate-850 rounded w-56 mt-2"></div>
            </div>
            <div class="h-8 bg-slate-800 rounded w-24"></div>
          </div>
          <div class="space-y-4">
            <div class="h-4 bg-slate-850 rounded w-full"></div>
            <div v-for="i in 6" :key="'rate-skel-'+i" class="flex items-center justify-between gap-4 py-2 border-b border-slate-900/40">
              <div class="flex items-center gap-3 w-1/3">
                <div class="h-5 w-5 bg-slate-800 rounded-full"></div>
                <div class="h-4 bg-slate-800 rounded w-10"></div>
              </div>
              <div class="h-4 bg-slate-800 rounded w-16"></div>
              <div class="h-4 bg-slate-800 rounded w-20"></div>
              <div class="h-4 bg-slate-800 rounded w-12"></div>
            </div>
          </div>
        </div>

        <!-- 2. Weather Skeleton -->
        <div class="relative w-full rounded-2xl border border-slate-900 bg-slate-900/20 p-6 flex flex-col gap-6">
          <div>
            <div class="h-5 bg-slate-800 rounded w-48"></div>
            <div class="h-3 bg-slate-850 rounded w-32 mt-2"></div>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
            <div v-for="i in 6" :key="'weather-skel-'+i" class="h-[175px] rounded-xl border border-slate-850 bg-slate-900/10 p-5 flex flex-col justify-between">
              <div class="flex justify-between items-center">
                <div class="h-4 bg-slate-800 rounded w-16"></div>
                <div class="h-4 bg-slate-850 rounded w-12"></div>
              </div>
              <div class="h-8 bg-slate-800 rounded w-20 mx-auto"></div>
              <div class="flex justify-between gap-4 mt-2">
                <div class="h-3.5 bg-slate-850 rounded w-8"></div>
                <div class="h-3.5 bg-slate-850 rounded w-8"></div>
              </div>
            </div>
          </div>
        </div>

        <!-- 3. Commodities Skeleton -->
        <div class="relative w-full rounded-2xl border border-slate-900 bg-slate-900/20 p-6 flex flex-col gap-6">
          <div class="flex justify-between items-start">
            <div>
              <div class="h-5 bg-slate-800 rounded w-44"></div>
              <div class="h-3 bg-slate-850 rounded w-64 mt-2"></div>
            </div>
            <div class="h-8 bg-slate-800 rounded w-24"></div>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div v-for="i in 3" :key="'comm-skel-'+i" class="h-[360px] rounded-xl border border-slate-850 bg-slate-900/10 p-5 flex flex-col justify-between">
              <div>
                <div class="flex justify-between items-center mb-1">
                  <div class="h-4 bg-slate-800 rounded w-16"></div>
                  <div class="h-4 bg-slate-800 rounded w-12"></div>
                </div>
                <div class="h-3 bg-slate-850 rounded w-10 mt-1"></div>
                <div class="h-6 bg-slate-800 rounded w-24 mt-3"></div>
              </div>
              <div class="h-[180px] bg-slate-900/40 border border-slate-850 rounded-lg flex items-center justify-center">
                <div class="h-4 bg-slate-850 rounded w-20"></div>
              </div>
            </div>
          </div>
        </div>

        <!-- 4. Market Skeleton -->
        <div class="relative w-full rounded-2xl border border-slate-900 bg-slate-900/20 p-6 flex flex-col gap-6">
          <div class="flex justify-between items-start">
            <div>
              <div class="h-5 bg-slate-800 rounded w-52"></div>
              <div class="h-3 bg-slate-850 rounded w-80 mt-2"></div>
            </div>
            <div class="h-8 bg-slate-800 rounded w-24"></div>
          </div>
          <!-- IHSG Skeleton -->
          <div class="h-[320px] rounded-xl border border-slate-850 bg-slate-900/10 p-5 md:p-6 flex flex-col md:flex-row gap-6">
            <div class="flex flex-col justify-between md:w-[30%]">
              <div>
                <div class="h-4 bg-slate-800 rounded w-24"></div>
                <div class="h-5 bg-slate-800 rounded w-32 mt-2"></div>
              </div>
              <div class="my-4 md:my-0">
                <div class="h-8 bg-slate-800 rounded w-32"></div>
                <div class="h-6 bg-slate-800 rounded w-24 mt-2"></div>
              </div>
              <div class="h-3 bg-slate-850 rounded w-20"></div>
            </div>
            <div class="flex-1 bg-slate-900/40 border border-slate-850 rounded-lg"></div>
          </div>
          <!-- Other Stocks Grid -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div v-for="i in 3" :key="'market-skel-'+i" class="h-[280px] rounded-xl border border-slate-850 bg-slate-900/10 p-5 flex flex-col justify-between">
              <div>
                <div class="flex justify-between items-center mb-1">
                  <div class="h-4 bg-slate-800 rounded w-16"></div>
                  <div class="h-4 bg-slate-855 rounded w-12"></div>
                </div>
                <div class="h-3 bg-slate-850 rounded w-10 mt-1"></div>
                <div class="h-5 bg-slate-800 rounded w-20 mt-3"></div>
              </div>
              <div class="h-[140px] bg-slate-900/40 border border-slate-855 rounded-lg"></div>
            </div>
          </div>
        </div>

        <!-- 5. Quakes Skeleton (Full Width) -->
        <div class="relative w-full rounded-2xl border border-slate-900 bg-slate-900/20 p-6 md:col-span-2 flex flex-col gap-6">
          <div class="flex justify-between items-center">
            <div>
              <div class="h-5 bg-slate-800 rounded w-48"></div>
              <div class="h-3 bg-slate-850 rounded w-72 mt-2"></div>
            </div>
            <div class="h-5 bg-slate-800 rounded w-24"></div>
          </div>
          <div class="space-y-3">
            <div v-for="i in 5" :key="'quake-skel-'+i" class="flex items-center gap-4 p-4 rounded-xl border border-slate-85 bg-slate-900/10">
              <div class="h-12 w-12 bg-slate-800 rounded-xl"></div>
              <div class="flex-1 space-y-2">
                <div class="h-4 bg-slate-800 rounded w-1/3"></div>
                <div class="h-3 bg-slate-850 rounded w-1/4"></div>
              </div>
              <div class="h-8 w-8 bg-slate-800 rounded-lg"></div>
            </div>
          </div>
        </div>

      </div>

      <!-- Main Widget Grid Layout (Data Loaded) -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
        
        <!-- Rates and Weather Widgets (Row 1 on Desktop/Tablet: 50/50) -->
        <RatesWidget disable-auto-poll />
        <WeatherWidget disable-auto-poll />

        <!-- Commodities and Market Widgets (Row 2 on Desktop/Tablet: 50/50) -->
        <CommoditiesWidget disable-auto-poll />
        <MarketWidget disable-auto-poll />

        <!-- Quakes Widget (Row 3 on Desktop/Tablet: Full Width) -->
        <div class="md:col-span-2">
          <QuakesWidget disable-auto-poll />
        </div>

      </div>

    </div>

    <!-- Floating Global Last Updated Badge (Bottom Right corner) -->
    <Transition name="fade">
      <div v-if="!isInitialLoading && globalLastUpdated" class="fixed bottom-6 right-6 z-50 flex items-center gap-2.5 px-4 py-2.5 bg-slate-900/90 hover:bg-slate-900 border border-slate-800/80 hover:border-indigo-500/50 backdrop-blur-md rounded-full shadow-2xl text-[11px] font-medium text-slate-300 transition-all duration-300 select-none group">
        <span class="relative flex h-2 w-2">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
          <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
        </span>
        <span class="tracking-wide">Update Terakhir: <span class="font-mono text-slate-100 font-bold ml-0.5">{{ formattedLastUpdated }}</span></span>
      </div>
    </Transition>

  </div>
</template>

<style scoped>
/* Fade transition for floating badge */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
