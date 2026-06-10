<script setup lang="ts">
import { useWeatherStore } from '~/stores/weather'
import { usePolling } from '~/composables/usePolling'

const props = withDefaults(defineProps<{
  disableAutoPoll?: boolean
}>(), {
  disableAutoPoll: false
})

const weatherStore = useWeatherStore()

// Auto-fetch and poll weather every 60s
if (!props.disableAutoPoll) {
  usePolling(() => weatherStore.fetchWeather(), 60000)
}

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

// Helper to get card styles based on temperature
const getTempColorClass = (temp: number) => {
  if (temp < 24) {
    // Biru (<24)
    return {
      border: 'hover:border-[#38BDF8]/30 hover:shadow-[0_0_15px_rgba(56,189,248,0.08)]',
      tempText: 'text-[#38BDF8]',
      badge: 'bg-[#38BDF8]/10 text-[#38BDF8] border-[#38BDF8]/20'
    }
  } else if (temp <= 30) {
    // Hijau (24-30)
    return {
      border: 'hover:border-[#10B981]/30 hover:shadow-[0_0_15px_rgba(16,185,129,0.08)]',
      tempText: 'text-[#10B981]',
      badge: 'bg-[#10B981]/10 text-[#10B981] border-[#10B981]/20'
    }
  } else {
    // Oranye (>30)
    return {
      border: 'hover:border-[#F59E0B]/30 hover:shadow-[0_0_15px_rgba(245,158,11,0.08)]',
      tempText: 'text-[#F59E0B]',
      badge: 'bg-[#F59E0B]/10 text-[#F59E0B] border-[#F59E0B]/20'
    }
  }
}
</script>

<template>
  <div class="relative w-full overflow-hidden rounded-2xl border border-white/6 bg-[#161B22] p-5 md:p-6 transition-all duration-300 hover:border-sky-400/30 hover:shadow-[0_0_20px_rgba(56,189,248,0.08)]">
    <!-- Gradient Accent Top Border -->
    <div class="absolute top-0 left-0 h-1 w-full bg-gradient-to-r from-sky-400 to-blue-500"></div>

    <!-- Header -->
    <div class="mb-6 mt-1">
      <h3 class="text-base font-bold text-text-primary flex items-center gap-2">
        <span class="text-sky-400">🌤️</span> Cuaca Kota Indonesia
      </h3>
      <p class="text-xs text-text-tertiary mt-1">
        Update terakhir: {{ weatherStore.data?.lastUpdated ? new Date(weatherStore.data.lastUpdated).toLocaleTimeString('id-ID') : '-' }}
      </p>
    </div>

    <!-- Error State -->
    <div v-if="weatherStore.error" class="bg-rose-500/10 border border-rose-500/20 rounded-xl p-4 mb-4 text-xs text-rose-400">
      <div class="flex items-center gap-2 font-semibold">
        <span>⚠️</span> Gagal memuat data cuaca
      </div>
      <p class="mt-1 opacity-90 font-mono">{{ weatherStore.error }}</p>
    </div>

    <!-- Grid Weather Card Layout (2 columns x 3 rows on large screen) -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <!-- Loading Skeleton State -->
      <template v-if="weatherStore.loading && (!weatherStore.data || !weatherStore.data.cities)">
        <div 
          v-for="n in 6" 
          :key="'skeleton-weather-'+n" 
          class="h-[140px] rounded-xl border border-white/5 bg-[#1C2128] animate-pulse p-4 flex flex-col justify-between"
        >
          <div class="flex items-start justify-between">
            <div class="h-4 bg-white/5 rounded w-1/2"></div>
            <div class="h-4 bg-white/5 rounded w-1/4"></div>
          </div>
          <div class="h-8 bg-white/5 rounded w-1/3 my-2 mx-auto"></div>
          <div class="flex justify-between gap-4 mt-2">
            <div class="h-3 bg-white/5 rounded w-1/3"></div>
            <div class="h-3 bg-white/5 rounded w-1/3"></div>
          </div>
        </div>
      </template>

      <!-- Empty State -->
      <div v-else-if="!weatherStore.data || !weatherStore.data.cities || weatherStore.data.cities.length === 0" class="col-span-full py-12 text-center text-text-tertiary">
        Tidak ada data cuaca kota tersedia
      </div>

      <!-- Real Cards -->
      <template v-else>
        <div
          v-for="city in weatherStore.data.cities"
          :key="city.city"
          :class="[
            'relative rounded-xl border border-white/5 bg-[#1C2128] p-4 transition-all duration-300 flex flex-col justify-between h-[140px]',
            getTempColorClass(city.temperature).border
          ]"
        >
          <!-- Top Area: City Name & Weather Description Badge -->
          <div class="flex items-start justify-between gap-2">
            <h4 class="font-bold text-sm text-text-primary tracking-tight">
              {{ city.city }}
            </h4>
            
            <span :class="['inline-flex items-center px-2 py-0.5 rounded-full text-[9px] font-bold uppercase tracking-wider border', getTempColorClass(city.temperature).badge]">
              {{ city.weatherDesc }}
            </span>
          </div>

          <!-- Middle Area: Temperature + Emoji Icon -->
          <div class="flex items-center justify-center gap-3 my-1">
            <span class="text-3xl animate-bounce-subtle" aria-hidden="true">
              {{ getWeatherEmoji(city.weatherCode, city.weatherDesc) }}
            </span>
            <span class="text-3xl font-extrabold tracking-tighter" :class="getTempColorClass(city.temperature).tempText">
              {{ Math.round(city.temperature) }}°C
            </span>
          </div>

          <!-- Bottom Area: Humidity & Wind Speed -->
          <div class="grid grid-cols-2 border-t border-white/5 pt-2.5 text-[10px] font-medium text-text-secondary">
            <div class="flex items-center gap-1.5 justify-start">
              <span>💧</span>
              <span class="text-text-primary font-bold">{{ city.humidity }}%</span>
            </div>
            
            <div class="flex items-center gap-1.5 justify-end border-l border-white/5 pl-2">
              <span>💨</span>
              <span class="text-text-primary font-bold">{{ city.windSpeed }} km/h</span>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
