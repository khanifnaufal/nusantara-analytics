<script setup lang="ts">
import { useQuakesStore } from '~/stores/quakes'
import { usePolling } from '~/composables/usePolling'

const props = withDefaults(defineProps<{
  disableAutoPoll?: boolean
}>(), {
  disableAutoPoll: false
})

const quakesStore = useQuakesStore()

// Auto-fetch and poll quakes every 60s
if (!props.disableAutoPoll) {
  usePolling(() => quakesStore.fetchQuakes(), 60000)
}

// Helper to determine relative time in Indonesian
const formatRelativeTime = (timeStr: string) => {
  if (!timeStr) return '-'
  const date = new Date(timeStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffSec = Math.floor(diffMs / 1000)
  const diffMin = Math.floor(diffSec / 60)
  const diffHour = Math.floor(diffMin / 60)
  const diffDay = Math.floor(diffHour / 24)

  if (diffSec < 0) {
    return 'Baru saja' // Time sync fallback
  }

  if (diffSec < 60) {
    return 'Baru saja'
  } else if (diffMin < 60) {
    return `${diffMin} menit lalu`
  } else if (diffHour < 24) {
    return `${diffHour} jam lalu`
  } else {
    return `${diffDay} hari lalu`
  }
}

// Get magnitude color classes
const getMagnitudeStyle = (magnitude: number) => {
  if (magnitude < 4.0) {
    // 3-4 (minor)
    return {
      badge: 'bg-zinc-800 text-zinc-400 border-white/5',
      text: 'font-semibold',
      pulse: false
    }
  } else if (magnitude < 5.0) {
    // 4-5 (moderate)
    return {
      badge: 'bg-zinc-700 text-zinc-300 border-white/5',
      text: 'font-bold',
      pulse: false
    }
  } else if (magnitude < 6.0) {
    // 5-6 (strong)
    return {
      badge: 'bg-zinc-600 text-white border-white/5',
      text: 'font-extrabold',
      pulse: false
    }
  } else {
    // 6+ (major)
    return {
      badge: 'bg-red-500/20 text-red-400 border-red-500/30',
      text: 'font-black',
      pulse: true
    }
  }
}
</script>

<template>
  <div class="widget-card">

    <!-- Header -->
    <div class="mb-6 mt-1 flex justify-between items-center">
      <div>
        <h3 class="text-base font-bold text-text-primary flex items-center gap-2">
          <span class="text-rose-400">🌋</span> Gempa Bumi Terbaru
        </h3>
        <p class="text-xs text-text-tertiary mt-1">
          Informasi gempa bumi real-time • Update terakhir: {{ quakesStore.data?.lastUpdated ? new Date(quakesStore.data.lastUpdated).toLocaleTimeString('id-ID') : '-' }}
        </p>
      </div>
      
      <!-- Stats Badge -->
      <span v-if="quakesStore.data?.total !== undefined" class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-bold bg-white/5 text-text-secondary border border-white/5">
        {{ quakesStore.data.total }} Data
      </span>
    </div>

    <!-- Error State -->
    <div v-if="quakesStore.error" class="bg-rose-500/10 border border-rose-500/20 rounded-xl p-4 mb-4 text-xs text-rose-400">
      <div class="flex items-center gap-2 font-semibold">
        <span>⚠️</span> Gagal memuat data gempa bumi
      </div>
      <p class="mt-1 opacity-90 font-mono">{{ quakesStore.error }}</p>
    </div>

    <!-- Quakes List Container -->
    <div class="pr-1 select-none">
      <!-- Loading Skeleton State -->
      <div v-if="quakesStore.loading && (!quakesStore.data || !quakesStore.data.quakes)" class="space-y-3 animate-pulse">
        <div v-for="n in 5" :key="'skeleton-quake-'+n" class="flex items-center gap-4 p-4 rounded-xl border border-white/5 bg-[#111111]">
          <div class="h-10 w-10 bg-white/5 rounded-lg"></div>
          <div class="flex-1 space-y-2">
            <div class="h-4 bg-white/5 rounded w-2/3"></div>
            <div class="h-3 bg-white/5 rounded w-1/3"></div>
          </div>
          <div class="h-4 bg-white/5 rounded w-12"></div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="!quakesStore.data || !quakesStore.data.quakes || quakesStore.data.quakes.length === 0" class="py-12 text-center text-text-tertiary">
        Tidak ada aktivitas gempa bumi terbaru
      </div>

      <!-- Quake List Items (scrollable viewport max-h-[400px]) -->
      <div v-else class="max-h-[400px] overflow-y-auto scrollbar-thin pr-1.5 space-y-3">
        <div
          v-for="quake in quakesStore.data.quakes"
          :key="quake.id"
          class="flex items-center gap-3.5 p-3 md:p-3.5 rounded-xl border border-white/5 bg-[#161616] hover:bg-white/[0.02] transition-all duration-300"
        >
          <!-- Magnitude Badge (Left) -->
          <div 
            :class="[
              'flex flex-col items-center justify-center h-12 w-12 rounded-xl border font-mono text-center flex-shrink-0 transition-transform duration-300 hover:scale-105',
              getMagnitudeStyle(quake.magnitude).badge,
              getMagnitudeStyle(quake.magnitude).pulse ? 'animate-major-pulse' : ''
            ]"
          >
            <span class="text-[9px] font-semibold uppercase tracking-tighter opacity-80">Mag</span>
            <span :class="['text-base tracking-tighter leading-tight', getMagnitudeStyle(quake.magnitude).text]">
              {{ quake.magnitude.toFixed(1) }}
            </span>
          </div>

          <!-- Location, Time and Alerts (Middle) -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 flex-wrap mb-1">
              <span class="text-xs font-bold text-text-primary truncate max-w-full">
                {{ quake.place }}
              </span>
              
              <!-- Tsunami Warning Badge -->
              <span 
                v-if="quake.tsunami === 1" 
                class="inline-flex items-center px-1.5 py-0.5 rounded text-[9px] font-extrabold uppercase bg-rose-500/10 text-rose-450 border border-rose-500/20 animate-pulse"
              >
                ⚠️ Tsunami
              </span>
            </div>

            <!-- Details: Depth, Time -->
            <div class="flex gap-x-4 gap-y-1 flex-wrap text-[10px] font-medium text-text-secondary">
              <span class="flex items-center gap-1.5">
                <span>📍</span>
                <span>Kedalaman: <strong class="text-text-primary font-bold">{{ quake.depth }} km</strong></span>
              </span>
              
              <span class="flex items-center gap-1.5">
                <span>⏱️</span>
                <span>{{ formatRelativeTime(quake.time) }}</span>
              </span>
            </div>
          </div>

          <!-- USGS Link (Right) -->
          <div class="flex-shrink-0">
            <a 
              :href="quake.url" 
              target="_blank" 
              rel="noopener noreferrer"
              class="inline-flex items-center justify-center p-2 rounded-lg bg-transparent border border-white/10 hover:border-white/20 text-zinc-400 hover:text-white transition-all shadow-xs"
              title="Lihat Detail di USGS"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
