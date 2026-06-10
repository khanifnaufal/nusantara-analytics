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
import CustomVizBuilder from '~/components/widgets/CustomVizBuilder.vue'

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

<<template>
  <div class="relative min-h-screen bg-[#0D1117] text-text-primary p-4 sm:p-6 md:p-8">
    <!-- Premium Decorative background glow blobs -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none z-0">
      <div class="absolute -top-40 -left-40 w-96 h-96 bg-indigo-500/10 rounded-full blur-3xl"></div>
      <div class="absolute top-1/2 -right-40 w-96 h-96 bg-violet-500/5 rounded-full blur-3xl"></div>
      <div class="absolute -bottom-40 left-1/3 w-96 h-96 bg-emerald-500/5 rounded-full blur-3xl"></div>
    </div>

    <div class="max-w-7xl mx-auto space-y-8 relative z-10">
      
      <!-- Top Header Area - Sticky Navbar style -->
      <header class="sticky top-4 z-40 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 p-4 md:p-5 rounded-2xl border border-white/5 bg-[#161B22]/80 backdrop-blur-xl shadow-lg shadow-black/10">
        <div>
          <div class="flex items-center gap-2.5">
            <span class="text-2xl" aria-hidden="true">🇮🇩</span>
            <h1 class="text-xl md:text-2xl font-extrabold tracking-tight text-transparent bg-clip-text bg-gradient-to-r from-white to-slate-400">
              Nusantara Analytics
            </h1>
          </div>
          <p class="text-[11px] md:text-xs text-text-secondary mt-1 font-medium">
            Dashboard Pemantauan Finansial, Cuaca, Komoditas Pasar, & Aktivitas Seismik Indonesia.
          </p>
        </div>

        <!-- System Navigation & Status -->
        <div class="flex flex-wrap items-center gap-3 self-start sm:self-auto">
          <!-- Link to Dedicated Builder -->
          <NuxtLink
            to="/builder"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 text-xs font-bold rounded-xl border border-white/5 bg-[#1C2128] hover:bg-[#21262D] text-text-primary transition-all shadow-sm"
          >
            <span>🎨</span> Custom Builder
          </NuxtLink>

          <!-- Live status pill -->
          <div class="inline-flex items-center gap-2 bg-[#1C2128] border border-white/5 py-1.5 px-3.5 rounded-xl text-xs font-bold text-text-secondary shadow-sm">
            <span class="relative flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-450 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            </span>
            <span class="tracking-wide">Live</span>
            <span class="text-white/10" aria-hidden="true">|</span>
            <span class="font-mono text-text-primary text-[11px]">{{ liveTime }}</span>
          </div>
        </div>
      </header>

      <!-- Skeleton Grid (Initial Loading state) -->
      <div v-if="isInitialLoading" class="grid grid-cols-1 md:grid-cols-12 gap-6 animate-pulse">
        <!-- 1. Rates Skeleton (5 col) -->
        <div class="col-span-12 md:col-span-6 lg:col-span-5 rounded-2xl border border-white/5 bg-[#161B22] p-5 flex flex-col gap-6 h-[400px]">
          <div class="flex justify-between items-start">
            <div class="space-y-2 w-1/2">
              <div class="h-4 bg-white/5 rounded"></div>
              <div class="h-3 bg-white/5 rounded w-2/3"></div>
            </div>
            <div class="h-6 bg-white/5 rounded w-16"></div>
          </div>
          <div class="space-y-4 flex-1">
            <div class="h-4 bg-white/5 rounded"></div>
            <div v-for="i in 5" :key="'rate-skel-'+i" class="flex justify-between">
              <div class="h-3 bg-white/5 rounded w-1/4"></div>
              <div class="h-3 bg-white/5 rounded w-1/5"></div>
              <div class="h-3 bg-white/5 rounded w-1/6"></div>
            </div>
          </div>
        </div>

        <!-- 2. Weather Skeleton (7 col) -->
        <div class="col-span-12 md:col-span-6 lg:col-span-7 rounded-2xl border border-white/5 bg-[#161B22] p-5 flex flex-col gap-6 h-[400px]">
          <div class="space-y-2">
            <div class="h-4 bg-white/5 rounded w-1/3"></div>
            <div class="h-3 bg-white/5 rounded w-1/4"></div>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 flex-1">
            <div v-for="i in 4" :key="'weather-skel-'+i" class="h-[140px] rounded-xl border border-white/5 bg-[#1C2128] p-4 flex flex-col justify-between">
              <div class="h-3 bg-white/5 rounded w-1/2"></div>
              <div class="h-6 bg-white/5 rounded w-1/3 mx-auto"></div>
              <div class="h-3 bg-white/5 rounded w-1/4"></div>
            </div>
          </div>
        </div>

        <!-- 3. Commodities Skeleton (4 col) -->
        <div class="col-span-12 md:col-span-12 lg:col-span-4 rounded-2xl border border-white/5 bg-[#161B22] p-5 flex flex-col gap-6 h-[380px]">
          <div class="flex justify-between">
            <div class="space-y-2 w-1/2">
              <div class="h-4 bg-white/5 rounded"></div>
              <div class="h-3 bg-white/5 rounded w-2/3"></div>
            </div>
            <div class="h-6 bg-white/5 rounded w-16"></div>
          </div>
          <div class="h-[220px] rounded-xl border border-white/5 bg-[#1C2128] p-4 flex flex-col justify-between">
            <div class="h-3 bg-white/5 rounded w-1/3"></div>
            <div class="h-10 bg-white/5 rounded w-1/2"></div>
            <div class="h-[100px] bg-white/5 rounded w-full"></div>
          </div>
        </div>

        <!-- 4. Market Skeleton (5 col) -->
        <div class="col-span-12 md:col-span-6 lg:col-span-5 rounded-2xl border border-white/5 bg-[#161B22] p-5 flex flex-col gap-6 h-[380px]">
          <div class="flex justify-between">
            <div class="space-y-2 w-1/2">
              <div class="h-4 bg-white/5 rounded"></div>
              <div class="h-3 bg-white/5 rounded w-2/3"></div>
            </div>
            <div class="h-6 bg-white/5 rounded w-16"></div>
          </div>
          <div class="h-[220px] rounded-xl border border-white/5 bg-[#1C2128] p-4 flex flex-col justify-between">
            <div class="h-3 bg-white/5 rounded w-1/3"></div>
            <div class="h-8 bg-white/5 rounded w-1/2"></div>
            <div class="h-[100px] bg-white/5 rounded w-full"></div>
          </div>
        </div>

        <!-- 5. Quakes Skeleton (3 col) -->
        <div class="col-span-12 md:col-span-6 lg:col-span-3 rounded-2xl border border-white/5 bg-[#161B22] p-5 flex flex-col gap-6 h-[380px]">
          <div class="space-y-2">
            <div class="h-4 bg-white/5 rounded w-2/3"></div>
            <div class="h-3 bg-white/5 rounded w-1/2"></div>
          </div>
          <div class="space-y-3 flex-1 overflow-hidden">
            <div v-for="i in 3" :key="'quake-skel-'+i" class="flex gap-3 p-3 rounded-xl border border-white/5 bg-[#1C2128]">
              <div class="h-8 w-8 bg-white/5 rounded-lg"></div>
              <div class="flex-1 space-y-2">
                <div class="h-3 bg-white/5 rounded w-2/3"></div>
                <div class="h-2 bg-white/5 rounded w-1/3"></div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 6. Custom Viz Builder Skeleton (12 col) -->
        <div class="col-span-12 rounded-2xl border border-white/5 bg-[#161B22] p-5 flex flex-col gap-6 h-[400px]">
          <div class="h-4 bg-white/5 rounded w-1/4"></div>
          <div class="h-12 bg-white/5 rounded w-full"></div>
          <div class="h-[200px] bg-white/5 rounded w-full"></div>
        </div>
      </div>

      <!-- Main Widget Grid Layout (Data Loaded) -->
      <div v-else class="grid grid-cols-1 md:grid-cols-12 gap-6">
        <!-- Rates (5 col on lg, 6 col on md, 12 col on mobile) -->
        <div class="col-span-12 md:col-span-6 lg:col-span-5 flex">
          <RatesWidget disable-auto-poll />
        </div>
        
        <!-- Weather (7 col on lg, 6 col on md, 12 col on mobile) -->
        <div class="col-span-12 md:col-span-6 lg:col-span-7 flex">
          <WeatherWidget disable-auto-poll />
        </div>

        <!-- Commodities (4 col on lg, 12 col on md, 12 col on mobile) -->
        <div class="col-span-12 md:col-span-12 lg:col-span-4 flex">
          <CommoditiesWidget disable-auto-poll />
        </div>

        <!-- Market (5 col on lg, 6 col on md, 12 col on mobile) -->
        <div class="col-span-12 md:col-span-6 lg:col-span-5 flex">
          <MarketWidget disable-auto-poll />
        </div>

        <!-- Quakes (3 col on lg, 6 col on md, 12 col on mobile) -->
        <div class="col-span-12 md:col-span-6 lg:col-span-3 flex">
          <QuakesWidget disable-auto-poll />
        </div>

        <!-- Custom Viz Builder (12 col always) -->
        <div class="col-span-12">
          <CustomVizBuilder />
        </div>
      </div>

    </div>

    <!-- Floating Global Last Updated Badge (Bottom Right corner) -->
    <Transition name="fade">
      <div v-if="!isInitialLoading && globalLastUpdated" class="fixed bottom-6 right-6 z-50 flex items-center gap-2 px-3.5 py-2.5 bg-[#161B22]/95 hover:bg-[#161B22] border border-white/8 hover:border-emerald-500/30 backdrop-blur-md rounded-full shadow-2xl text-[10px] font-semibold text-text-secondary transition-all duration-300 select-none group">
        <span class="relative flex h-2 w-2">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
          <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
        </span>
        <span class="tracking-wide">Terakhir Diperbarui: <span class="font-mono text-text-primary font-bold ml-0.5">{{ formattedLastUpdated }}</span></span>
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
