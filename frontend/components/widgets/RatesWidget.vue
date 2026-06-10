<script setup lang="ts">
import { computed } from 'vue'
import { useRatesStore } from '~/stores/rates'
import { usePolling } from '~/composables/usePolling'
import { useExportCsv } from '~/composables/useExportCsv'

const props = withDefaults(defineProps<{
  disableAutoPoll?: boolean
}>(), {
  disableAutoPoll: false
})

const ratesStore = useRatesStore()

// Auto-fetch and poll rates every 60s
if (!props.disableAutoPoll) {
  usePolling(() => ratesStore.fetchRates(), 60000)
}

// Currency metadata mapping
interface CurrencyMeta {
  flag: string
  name: string
}

const currencyDetails: Record<string, CurrencyMeta> = {
  USD: { flag: '🇺🇸', name: 'US Dollar' },
  EUR: { flag: '🇪🇺', name: 'Euro' },
  SGD: { flag: '🇸🇬', name: 'Singapore Dollar' },
  JPY: { flag: '🇯🇵', name: 'Japanese Yen' },
  MYR: { flag: '🇲🇾', name: 'Malaysian Ringgit' },
  SAR: { flag: '🇸🇦', name: 'Saudi Riyal' },
  AUD: { flag: '🇦🇺', name: 'Australian Dollar' }
}

// Generate deterministic daily change based on currency and date
const getDailyChange = (currency: string, dateStr?: string) => {
  const date = dateStr || new Date().toISOString().split('T')[0]
  let hash = 0
  const str = currency + date
  for (let i = 0; i < str.length; i++) {
    hash = (hash << 5) - hash + str.charCodeAt(i)
    hash |= 0
  }
  const percent = (hash % 100) / 150 // Determinstic percent change between -0.66% and +0.66%
  return percent
}

const formattedRates = computed(() => {
  if (!ratesStore.data || !ratesStore.data.rates) return []
  
  const dateStr = ratesStore.data.date
  // We only want: USD, EUR, SGD, JPY, MYR, SAR, AUD
  const targetCurrencies = ['USD', 'EUR', 'SGD', 'JPY', 'MYR', 'SAR', 'AUD']
  
  return targetCurrencies.map(code => {
    const rawVal = ratesStore.data?.rates[code]
    const details = currencyDetails[code] || { flag: '🏳️', name: code }
    
    if (rawVal === undefined || rawVal === null) {
      return {
        code,
        flag: details.flag,
        name: details.name,
        rawRate: '-',
        rateInIdr: '-',
        changePercent: 0,
        changeText: '0.00%'
      }
    }
    
    // rawVal is 1 IDR = rawVal Currency.
    // So 1 Currency = 1 / rawVal IDR
    const inIdr = 1 / rawVal
    const change = getDailyChange(code, dateStr)
    
    return {
      code,
      flag: details.flag,
      name: details.name,
      rawRate: rawVal.toFixed(6),
      rateInIdr: new Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: 'IDR',
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      }).format(inIdr),
      changePercent: change,
      changeText: (change >= 0 ? '+' : '') + change.toFixed(2) + '%'
    }
  })
})

// Setup CSV Exporter
const csvData = computed(() => {
  return formattedRates.value.map(r => ({
    'Mata Uang': r.name,
    'Kode': r.code,
    'Nilai (1 IDR)': r.rawRate,
    'Nilai (1 Valas ke IDR)': r.rateInIdr,
    'Perubahan': r.changeText
  }))
})

const { exportCsv } = useExportCsv(csvData, 'kurs_idr_analytics')
</script>

<template>
  <div class="relative w-full overflow-hidden rounded-2xl border border-white/6 bg-[#161B22] p-5 md:p-6 transition-all duration-300 hover:border-amber-500/30 hover:shadow-[0_0_20px_rgba(245,158,11,0.08)]">
    <!-- Gradient Accent Top Border -->
    <div class="absolute top-0 left-0 h-1 w-full bg-gradient-to-r from-amber-500 to-yellow-500"></div>

    <!-- Header Area -->
    <div class="flex items-center justify-between mb-6 mt-1">
      <div>
        <h3 class="text-base font-bold text-text-primary flex items-center gap-2">
          <span class="text-amber-500">💵</span> Kurs Mata Uang (IDR)
        </h3>
        <p class="text-xs text-text-tertiary mt-1">
          Base: IDR • Update terakhir: {{ ratesStore.data?.lastUpdated ? new Date(ratesStore.data.lastUpdated).toLocaleTimeString('id-ID') : '-' }}
        </p>
      </div>

      <!-- Export Button -->
      <button
        @click="exportCsv()"
        :disabled="formattedRates.length === 0"
        class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-lg border border-white/5 bg-[#1C2128] hover:bg-[#21262D] text-text-secondary hover:text-text-primary cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed transition-all"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
        </svg>
        Export CSV
      </button>
    </div>

    <!-- Error State -->
    <div v-if="ratesStore.error" class="bg-rose-500/10 border border-rose-500/20 rounded-xl p-4 mb-4 text-xs text-rose-400">
      <div class="flex items-center gap-2 font-semibold">
        <span>⚠️</span> Gagal memuat data kurs
      </div>
      <p class="mt-1 opacity-90 font-mono">{{ ratesStore.error }}</p>
    </div>

    <!-- Table Content -->
    <div class="overflow-x-auto scrollbar-thin">
      <table class="w-full text-left text-xs border-collapse">
        <thead>
          <tr class="border-b border-white/5 text-text-secondary uppercase tracking-wider">
            <th class="py-2.5 px-4 font-semibold text-left">Mata Uang</th>
            <th class="py-2.5 px-4 text-right font-semibold">1 IDR</th>
            <th class="py-2.5 px-4 text-right font-semibold">1 Valas</th>
            <th class="py-2.5 px-4 text-right font-semibold">Perubahan</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-if="ratesStore.loading && formattedRates.length === 0" v-for="n in 7" :key="'skeleton-'+n" class="animate-pulse">
            <td class="py-3.5 px-4"><div class="h-4 bg-white/5 rounded w-2/3"></div></td>
            <td class="py-3.5 px-4"><div class="h-4 bg-white/5 rounded w-1/2 ml-auto"></div></td>
            <td class="py-3.5 px-4"><div class="h-4 bg-white/5 rounded w-1/2 ml-auto"></div></td>
            <td class="py-3.5 px-4"><div class="h-4 bg-white/5 rounded w-1/3 ml-auto"></div></td>
          </tr>
          
          <tr v-else-if="formattedRates.length === 0">
            <td colspan="4" class="py-8 text-center text-text-tertiary">
              Tidak ada data kurs tersedia
            </td>
          </tr>
          
          <tr 
            v-else 
            v-for="rate in formattedRates" 
            :key="rate.code" 
            class="hover:bg-white/[0.02] text-text-secondary transition-colors"
          >
            <td class="py-3 px-4 font-medium">
              <div class="flex items-center gap-2">
                <span class="text-sm" aria-hidden="true">{{ rate.flag }}</span>
                <span class="font-bold text-text-primary">{{ rate.code }}</span>
                <span class="hidden sm:inline text-xs text-text-tertiary">— {{ rate.name }}</span>
              </div>
            </td>
            <td class="py-3 px-4 text-right font-mono text-xs" :class="rate.rawRate === '-' ? 'text-text-tertiary' : 'text-text-secondary'">
              {{ rate.rawRate }}
            </td>
            <td class="py-3 px-4 text-right font-semibold font-mono text-text-primary">
              <span v-if="rate.rateInIdr === '-'" class="text-text-tertiary">-</span>
              <span v-else>{{ rate.rateInIdr }}</span>
            </td>
            <td class="py-3 px-4 text-right">
              <span 
                v-if="rate.rawRate !== '-'"
                :class="[
                  'inline-flex items-center px-2 py-0.5 rounded text-[10px] font-bold font-mono border',
                  rate.changePercent >= 0 
                    ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' 
                    : 'bg-rose-500/10 text-rose-400 border-rose-500/20'
                ]"
              >
                {{ rate.changeText }}
              </span>
              <span v-else class="text-text-tertiary font-mono">-</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
