<script setup lang="ts">
import { useQuakesStore } from '~/stores/quakes'
import { usePolling } from '~/composables/usePolling'

const quakesStore = useQuakesStore()

// Auto-fetch and poll quakes every 60s
usePolling(() => quakesStore.fetchQuakes(), 60000)

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
const getMagnitudeStyle = (level: string) => {
  switch (level?.toLowerCase()) {
    case 'major':
      return {
        badge: 'bg-red-50 dark:bg-red-950/30 text-red-600 dark:text-red-400 border-red-200 dark:border-red-800/40',
        text: 'font-black'
      }
    case 'strong':
      return {
        badge: 'bg-orange-50 dark:bg-orange-950/30 text-orange-600 dark:text-orange-400 border-orange-200 dark:border-orange-800/40',
        text: 'font-extrabold'
      }
    case 'moderate':
      return {
        badge: 'bg-yellow-50 dark:bg-yellow-950/30 text-yellow-600 dark:text-yellow-400 border-yellow-200 dark:border-yellow-800/40',
        text: 'font-bold'
      }
    case 'minor':
    default:
      return {
        badge: 'bg-slate-100 dark:bg-slate-800/40 text-slate-600 dark:text-slate-400 border-slate-200 dark:border-slate-800/40',
        text: 'font-semibold'
      }
  }
}
</script>

<template>
  <div class="relative w-full rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900/60 p-6 shadow-md transition-all duration-300 hover:shadow-lg backdrop-blur-md">
    <!-- Header -->
    <div class="mb-6 flex justify-between items-center">
      <div>
        <h3 class="text-lg font-bold text-slate-800 dark:text-slate-100 flex items-center gap-2">
          <span>🌋</span> Gempa Bumi Terbaru
        </h3>
        <p class="text-xs text-slate-400 dark:text-slate-500 mt-1">
          Informasi gempa bumi real-time • Update terakhir: {{ quakesStore.data?.lastUpdated ? new Date(quakesStore.data.lastUpdated).toLocaleString('id-ID') : '-' }}
        </p>
      </div>
      
      <!-- Stats Badge -->
      <span v-if="quakesStore.data?.total !== undefined" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-350">
        {{ quakesStore.data.total }} Terdeteksi
      </span>
    </div>

    <!-- Error State -->
    <div v-if="quakesStore.error" class="bg-red-50 dark:bg-red-950/20 border border-red-200 dark:border-red-800/40 rounded-xl p-4 mb-4 text-sm text-red-600 dark:text-red-400">
      <div class="flex items-center gap-2 font-medium">
        <span>⚠️</span> Gagal memuat data gempa bumi
      </div>
      <p class="text-xs mt-1 opacity-90">{{ quakesStore.error }}</p>
    </div>

    <!-- Quakes List Container -->
    <div class="max-h-[500px] overflow-y-auto pr-1 select-none custom-scrollbar">
      <!-- Loading Skeleton State -->
      <div v-if="quakesStore.loading && (!quakesStore.data || !quakesStore.data.quakes)" class="space-y-3 animate-pulse">
        <div v-for="n in 5" :key="'skeleton-quake-'+n" class="flex items-center gap-4 p-4 rounded-xl border border-slate-100 dark:border-slate-800/60 bg-slate-50/50 dark:bg-slate-900/20">
          <div class="h-10 w-10 bg-slate-200 dark:bg-slate-800 rounded-lg"></div>
          <div class="flex-1 space-y-2">
            <div class="h-4 bg-slate-200 dark:bg-slate-800 rounded w-2/3"></div>
            <div class="h-3 bg-slate-200 dark:bg-slate-800 rounded w-1/3"></div>
          </div>
          <div class="h-4 bg-slate-200 dark:bg-slate-800 rounded w-12"></div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="!quakesStore.data || !quakesStore.data.quakes || quakesStore.data.quakes.length === 0" class="py-12 text-center text-slate-400 dark:text-slate-500">
        Tidak ada aktivitas gempa bumi terbaru
      </div>

      <!-- Quake List Items -->
      <div v-else class="space-y-3">
        <div
          v-for="quake in quakesStore.data.quakes"
          :key="quake.id"
          class="flex items-center gap-4 p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/30 hover:bg-slate-50/70 dark:bg-slate-900/30 dark:hover:bg-slate-800/20 transition-all duration-300 shadow-xs hover:shadow-sm"
        >
          <!-- Magnitude Badge (Left) -->
          <div 
            :class="[
              'flex flex-col items-center justify-center h-12 w-12 rounded-xl border font-mono text-center flex-shrink-0 shadow-inner',
              getMagnitudeStyle(quake.magnitudeLevel).badge
            ]"
          >
            <span class="text-xs font-semibold uppercase tracking-tighter opacity-80">Mag</span>
            <span :class="['text-lg tracking-tighter', getMagnitudeStyle(quake.magnitudeLevel).text]">
              {{ quake.magnitude.toFixed(1) }}
            </span>
          </div>

          <!-- Location, Time and Alerts (Middle) -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 flex-wrap mb-1">
              <span class="text-xs font-bold text-slate-700 dark:text-slate-350 truncate max-w-full">
                {{ quake.place }}
              </span>
              
              <!-- Tsunami Warning Badge -->
              <span 
                v-if="quake.tsunami === 1" 
                class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-extrabold uppercase bg-red-150 dark:bg-red-950/40 text-red-600 dark:text-red-400 border border-red-200 dark:border-red-800/40 animate-pulse"
              >
                ⚠️ Potensi Tsunami
              </span>
            </div>

            <!-- Details: Depth, Time -->
            <div class="flex gap-x-4 gap-y-1 flex-wrap text-[11px] font-medium text-slate-400 dark:text-slate-500">
              <span class="flex items-center gap-1">
                <span>📍</span>
                <span>Kedalaman: <strong class="text-slate-600 dark:text-slate-400 font-semibold">{{ quake.depth }} km</strong></span>
              </span>
              
              <span class="flex items-center gap-1">
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
              class="inline-flex items-center justify-center p-2 rounded-lg border border-slate-200 hover:border-slate-300 dark:border-slate-800 dark:hover:border-slate-700 bg-white hover:bg-slate-50 dark:bg-slate-900/60 dark:hover:bg-slate-800 text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 transition-all shadow-xs"
              title="Lihat Detail di USGS"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Scrollbar Customization for a premium look */
.custom-scrollbar::-webkit-scrollbar {
  width: 5px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(156, 163, 175, 0.25);
  border-radius: 9999px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(156, 163, 175, 0.45);
}
</style>
