<script setup lang="ts">
import { useWeatherStore } from '~/stores/weather'
import { usePolling } from '~/composables/usePolling'

const weatherStore = useWeatherStore()

// Auto-fetch and poll weather every 60s
usePolling(() => weatherStore.fetchWeather(), 60000)

// Helper to determine weather emoji
const getWeatherEmoji = (code: number, desc: string) => {
  const lowercaseDesc = desc.toLowerCase()
  if (code === 0 || code === 1 || lowercaseDesc.includes('sunny') || lowercaseDesc.includes('cerah') || lowercaseDesc.includes('clear')) return '☀️'
  if (code === 2 || code === 3 || lowercaseDesc.includes('cloudy') || lowercaseDesc.includes('berawan') || lowercaseDesc.includes('cloud')) return '⛅'
  if ((code >= 51 && code <= 67) || (code >= 80 && code <= 82) || lowercaseDesc.includes('rain') || lowercaseDesc.includes('hujan') || lowercaseDesc.includes('gerimis') || lowercaseDesc.includes('drizzle')) return '🌧️'
  if (code >= 95 && code <= 99 || lowercaseDesc.includes('thunder') || lowercaseDesc.includes('storm') || lowercaseDesc.includes('petir') || lowercaseDesc.includes('kilat')) return '⛈️'
  if (code >= 71 && code <= 86 || lowercaseDesc.includes('snow') || lowercaseDesc.includes('salju')) return '❄️'
  if (code >= 45 && code <= 48 || lowercaseDesc.includes('fog') || lowercaseDesc.includes('kabut')) return '🌫️'
  return '🌡️'
}

// Helper to get card color class based on temperature
const getTempColorClass = (temp: number) => {
  if (temp < 24) {
    // Biru (<24)
    return {
      card: 'bg-gradient-to-br from-blue-500/10 to-blue-600/5 dark:from-blue-950/20 dark:to-slate-900 border-blue-200 dark:border-blue-900/40 text-blue-950 dark:text-blue-100 hover:border-blue-400 dark:hover:border-blue-700/60 shadow-blue-500/5',
      badge: 'bg-blue-50 dark:bg-blue-950/40 text-blue-600 dark:text-blue-400 border border-blue-100 dark:border-blue-800/30'
    }
  } else if (temp <= 30) {
    // Hijau (24-30)
    return {
      card: 'bg-gradient-to-br from-emerald-500/10 to-emerald-600/5 dark:from-emerald-950/20 dark:to-slate-900 border-emerald-200 dark:border-emerald-900/40 text-emerald-950 dark:text-emerald-100 hover:border-emerald-400 dark:hover:border-emerald-700/60 shadow-emerald-500/5',
      badge: 'bg-emerald-50 dark:bg-emerald-950/40 text-emerald-600 dark:text-emerald-400 border border-emerald-100 dark:border-emerald-800/30'
    }
  } else {
    // Oranye (>30)
    return {
      card: 'bg-gradient-to-br from-amber-500/10 to-amber-600/5 dark:from-amber-950/20 dark:to-slate-900 border-amber-200 dark:border-amber-900/40 text-amber-950 dark:text-amber-100 hover:border-amber-400 dark:hover:border-amber-700/60 shadow-amber-500/5',
      badge: 'bg-amber-50 dark:bg-amber-950/40 text-amber-600 dark:text-amber-400 border border-amber-100 dark:border-amber-800/30'
    }
  }
}
</script>

<template>
  <div class="relative w-full rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900/60 p-6 shadow-md transition-all duration-300 hover:shadow-lg backdrop-blur-md">
    <!-- Header -->
    <div class="mb-6">
      <h3 class="text-lg font-bold text-slate-800 dark:text-slate-100 flex items-center gap-2">
        <span>🌤️</span> Cuaca Kota Indonesia
      </h3>
      <p class="text-xs text-slate-400 dark:text-slate-500 mt-1">
        Update terakhir: {{ weatherStore.data?.lastUpdated ? new Date(weatherStore.data.lastUpdated).toLocaleString('id-ID') : '-' }}
      </p>
    </div>

    <!-- Error State -->
    <div v-if="weatherStore.error" class="bg-red-50 dark:bg-red-950/20 border border-red-200 dark:border-red-800/40 rounded-xl p-4 mb-4 text-sm text-red-600 dark:text-red-400">
      <div class="flex items-center gap-2 font-medium">
        <span>⚠️</span> Gagal memuat data cuaca
      </div>
      <p class="text-xs mt-1 opacity-90">{{ weatherStore.error }}</p>
    </div>

    <!-- Grid Weather Card Layout -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
      <!-- Loading Skeleton State -->
      <template v-if="weatherStore.loading && (!weatherStore.data || !weatherStore.data.cities)">
        <div 
          v-for="n in 6" 
          :key="'skeleton-weather-'+n" 
          class="h-[175px] rounded-xl border border-slate-100 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/30 animate-pulse p-5 flex flex-col justify-between"
        >
          <div class="flex items-start justify-between">
            <div class="h-5 bg-slate-200 dark:bg-slate-800 rounded w-1/2"></div>
            <div class="h-5 bg-slate-200 dark:bg-slate-800 rounded w-1/4"></div>
          </div>
          <div class="h-10 bg-slate-200 dark:bg-slate-800 rounded w-1/3 my-2 mx-auto"></div>
          <div class="flex justify-between gap-4 mt-2">
            <div class="h-4 bg-slate-200 dark:bg-slate-800 rounded w-1/3"></div>
            <div class="h-4 bg-slate-200 dark:bg-slate-800 rounded w-1/3"></div>
          </div>
        </div>
      </template>

      <!-- Empty State -->
      <div v-else-if="!weatherStore.data || !weatherStore.data.cities || weatherStore.data.cities.length === 0" class="col-span-full py-12 text-center text-slate-400 dark:text-slate-500">
        Tidak ada data cuaca kota tersedia
      </div>

      <!-- Real Cards -->
      <template v-else>
        <div
          v-for="city in weatherStore.data.cities"
          :key="city.city"
          :class="[
            'relative rounded-xl border p-5 shadow-sm hover:shadow-md transition-all duration-300 flex flex-col justify-between h-[175px]',
            getTempColorClass(city.temperature).card
          ]"
        >
          <!-- Top Area: City Name & Weather Description Badge -->
          <div class="flex items-start justify-between gap-2">
            <h4 class="font-bold text-base text-slate-800 dark:text-slate-100 tracking-tight truncate">
              {{ city.city }}
            </h4>
            
            <span :class="['inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider', getTempColorClass(city.temperature).badge]">
              {{ city.weatherDesc }}
            </span>
          </div>

          <!-- Middle Area: Huge Temperature + Emoji Icon -->
          <div class="flex items-center justify-center gap-2 my-2">
            <span class="text-4xl" aria-hidden="true">
              {{ getWeatherEmoji(city.weatherCode, city.weatherDesc) }}
            </span>
            <span class="text-4xl font-extrabold tracking-tighter">
              {{ Math.round(city.temperature) }}°C
            </span>
          </div>

          <!-- Bottom Area: Humidity & Wind Speed -->
          <div class="grid grid-cols-2 border-t border-slate-200/50 dark:border-slate-800/40 pt-3 text-[11px] font-medium opacity-85">
            <div class="flex items-center gap-1 text-slate-500 dark:text-slate-400 justify-start">
              <span>💧</span>
              <span class="text-slate-700 dark:text-slate-300 font-semibold">{{ city.humidity }}%</span>
            </div>
            
            <div class="flex items-center gap-1 text-slate-500 dark:text-slate-400 justify-end border-l border-slate-200/50 dark:border-slate-800/40 pl-2">
              <span>💨</span>
              <span class="text-slate-700 dark:text-slate-300 font-semibold">{{ city.windSpeed }} km/h</span>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
