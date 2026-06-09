<script setup lang="ts">
import { computed } from 'vue'
import { useCommoditiesStore } from '~/stores/commodities'
import { usePolling } from '~/composables/usePolling'
import { useExportCsv } from '~/composables/useExportCsv'
import AreaChart from '../charts/AreaChart.vue'

const commoditiesStore = useCommoditiesStore()

// Auto-fetch and poll commodities every 60s
usePolling(() => commoditiesStore.fetchCommodities(), 60000)

// Helper to format prices according to currency
const formatPrice = (value: number, currency: string) => {
  if (currency === 'USD') {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(value)
  } else if (currency === 'MYR') {
    return new Intl.NumberFormat('en-MY', { style: 'currency', currency: 'MYR' }).format(value).replace('MYR', 'RM')
  } else if (currency === 'IDR') {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(value)
  }
  return `${currency} ${new Intl.NumberFormat(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(value)}`
}

// Generate chart data for a commodity
const getChartProps = (history: { timestamp: string; close: number }[], name: string) => {
  if (!history || history.length === 0) {
    return { xAxis: [] as string[], data: [] as { name: string; data: number[] }[] }
  }

  // Sort history by timestamp ascending to make sure chart draws correctly left-to-right
  const sortedHistory = [...history].sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime())

  const xAxis = sortedHistory.map(h => {
    const d = new Date(h.timestamp)
    return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short' })
  })

  const data = [
    {
      name,
      data: sortedHistory.map(h => parseFloat(h.close.toFixed(2)))
    }
  ]

  return { xAxis, data }
}

// Export CSV Setup
const csvData = computed(() => {
  if (!commoditiesStore.data?.commodities) return []
  return commoditiesStore.data.commodities.map(c => ({
    Simbol: c.symbol,
    Nama: c.name,
    Harga: c.price,
    MataUang: c.currency,
    Perubahan: c.change,
    'Persen Perubahan (%)': c.changePercent.toFixed(2) + '%',
    'Update Terakhir': c.lastUpdated ? new Date(c.lastUpdated).toLocaleString('id-ID') : '-'
  }))
})

const { exportCsv } = useExportCsv(csvData, 'komoditas_harga_analytics')
</script>

<template>
  <div class="relative w-full rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900/60 p-6 shadow-md transition-all duration-300 hover:shadow-lg backdrop-blur-md">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h3 class="text-lg font-bold text-slate-800 dark:text-slate-100 flex items-center gap-2">
          <span>📈</span> Komoditas Global
        </h3>
        <p class="text-xs text-slate-400 dark:text-slate-500 mt-1">
          Harga pasar komoditas dunia real-time • Update terakhir: {{ commoditiesStore.data?.lastUpdated ? new Date(commoditiesStore.data.lastUpdated).toLocaleString('id-ID') : '-' }}
        </p>
      </div>

      <!-- Export Button -->
      <button
        @click="exportCsv()"
        :disabled="!commoditiesStore.data?.commodities || commoditiesStore.data.commodities.length === 0"
        class="inline-flex items-center gap-2 px-3 py-1.5 text-xs font-semibold rounded-lg border border-slate-200 dark:border-slate-800 bg-slate-50 hover:bg-slate-100 dark:bg-slate-900/80 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed transition-all"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
        </svg>
        Export CSV
      </button>
    </div>

    <!-- Error State -->
    <div v-if="commoditiesStore.error" class="bg-red-50 dark:bg-red-950/20 border border-red-200 dark:border-red-800/40 rounded-xl p-4 mb-4 text-sm text-red-600 dark:text-red-400">
      <div class="flex items-center gap-2 font-medium">
        <span>⚠️</span> Gagal memuat data komoditas
      </div>
      <p class="text-xs mt-1 opacity-90">{{ commoditiesStore.error }}</p>
    </div>

    <!-- Commodities Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <!-- Loading Skeleton State -->
      <template v-if="commoditiesStore.loading && (!commoditiesStore.data || !commoditiesStore.data.commodities)">
        <div 
          v-for="n in 3" 
          :key="'skeleton-comm-'+n" 
          class="h-[360px] rounded-xl border border-slate-100 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/30 animate-pulse p-5 flex flex-col justify-between"
        >
          <div>
            <div class="flex justify-between items-center mb-4">
              <div class="h-5 bg-slate-200 dark:bg-slate-800 rounded w-1/3"></div>
              <div class="h-5 bg-slate-200 dark:bg-slate-800 rounded w-1/4"></div>
            </div>
            <div class="h-8 bg-slate-200 dark:bg-slate-800 rounded w-1/2 mb-4"></div>
          </div>
          <!-- Chart Placeholder -->
          <div class="h-[180px] bg-slate-200 dark:bg-slate-800 rounded-lg w-full flex items-center justify-center">
            <span class="text-xs text-slate-400">Loading chart...</span>
          </div>
        </div>
      </template>

      <!-- Empty State -->
      <div v-else-if="!commoditiesStore.data || !commoditiesStore.data.commodities || commoditiesStore.data.commodities.length === 0" class="col-span-full py-12 text-center text-slate-400 dark:text-slate-500">
        Tidak ada data komoditas tersedia
      </div>

      <!-- Real Cards -->
      <template v-else>
        <div
          v-for="item in commoditiesStore.data.commodities"
          :key="item.symbol"
          class="rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/40 p-5 flex flex-col justify-between shadow-sm hover:shadow-md transition-all duration-300"
        >
          <!-- Price and Badge -->
          <div class="mb-4">
            <div class="flex items-center justify-between mb-1">
              <div>
                <h4 class="font-bold text-base text-slate-900 dark:text-slate-100 tracking-tight">
                  {{ item.name }}
                </h4>
                <p class="text-[10px] text-slate-400 dark:text-slate-500 font-semibold uppercase font-mono tracking-wider">
                  {{ item.symbol }}
                </p>
              </div>

              <!-- Change Badge -->
              <span 
                :class="[
                  'inline-flex items-center px-2 py-0.5 rounded text-xs font-bold',
                  item.changePercent >= 0 
                    ? 'bg-emerald-50 dark:bg-emerald-950/20 text-emerald-600 dark:text-emerald-400' 
                    : 'bg-rose-50 dark:bg-rose-950/20 text-rose-600 dark:text-rose-400'
                ]"
              >
                {{ item.changePercent >= 0 ? '+' : '' }}{{ item.changePercent.toFixed(2) }}%
              </span>
            </div>

            <!-- Price -->
            <div class="text-2xl font-black text-slate-900 dark:text-slate-100 tracking-tight font-mono">
              {{ formatPrice(item.price, item.currency) }}
            </div>
          </div>

          <!-- History Chart (Last 7 Days) -->
          <div class="w-full mt-2">
            <p class="text-[10px] text-slate-400 dark:text-slate-500 font-bold uppercase tracking-wider mb-2">
              Tren 7 Hari Terakhir
            </p>
            <div class="h-[180px] w-full">
              <AreaChart
                v-bind="getChartProps(item.history, item.name)"
                :loading="commoditiesStore.loading"
                :unit="item.currency"
                height="100%"
              />
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
