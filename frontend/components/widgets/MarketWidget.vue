<script setup lang="ts">
import { computed } from 'vue'
import { useMarketStore } from '~/stores/market'
import { usePolling } from '~/composables/usePolling'
import { useExportCsv } from '~/composables/useExportCsv'
import LineChart from '../charts/LineChart.vue'
import AreaChart from '../charts/AreaChart.vue'

const marketStore = useMarketStore()

// Auto-fetch and poll market data every 60s
usePolling(() => marketStore.fetchMarket(), 60000)

// Special formatter for stock index / prices
const formatMarketPrice = (value: number, symbol: string) => {
  if (symbol === '^JKSE') {
    return new Intl.NumberFormat('id-ID', { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(value)
  }
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(value)
}

// Separate IHSG and other stocks
const ihsgStock = computed(() => {
  if (!marketStore.data?.stocks) return null
  return marketStore.data.stocks.find(s => s.symbol === '^JKSE') || null
})

const otherStocks = computed(() => {
  if (!marketStore.data?.stocks) return []
  return marketStore.data.stocks.filter(s => s.symbol !== '^JKSE')
})

// Generate chart props (timestamp sorting, axis labeling)
const getChartProps = (history: { timestamp: string; close: number }[], name: string) => {
  if (!history || history.length === 0) {
    return { xAxis: [] as string[], data: [] as { name: string; data: number[] }[] }
  }

  // Sort history ascending to draw correctly left-to-right
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
  if (!marketStore.data?.stocks) return []
  return marketStore.data.stocks.map(s => ({
    Simbol: s.symbol,
    Nama: s.name,
    'Harga / Indeks': s.price,
    MataUang: s.currency,
    Perubahan: s.change,
    'Persen Perubahan (%)': s.changePercent.toFixed(2) + '%',
    'Update Terakhir': s.lastUpdated ? new Date(s.lastUpdated).toLocaleString('id-ID') : '-'
  }))
})

const { exportCsv } = useExportCsv(csvData, 'pasar_saham_analytics')
</script>

<template>
  <div class="relative w-full rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900/60 p-6 shadow-md transition-all duration-300 hover:shadow-lg backdrop-blur-md">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h3 class="text-lg font-bold text-slate-800 dark:text-slate-100 flex items-center gap-2">
          <span>🏛️</span> Bursa Saham Indonesia
        </h3>
        <p class="text-xs text-slate-400 dark:text-slate-500 mt-1">
          Indeks Harga Saham Gabungan (IHSG) & Saham Unggulan • Update terakhir: {{ marketStore.data?.lastUpdated ? new Date(marketStore.data.lastUpdated).toLocaleString('id-ID') : '-' }}
        </p>
      </div>

      <!-- Export Button -->
      <button
        @click="exportCsv()"
        :disabled="!marketStore.data?.stocks || marketStore.data.stocks.length === 0"
        class="inline-flex items-center gap-2 px-3 py-1.5 text-xs font-semibold rounded-lg border border-slate-200 dark:border-slate-800 bg-slate-50 hover:bg-slate-100 dark:bg-slate-900/80 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed transition-all"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
        </svg>
        Export CSV
      </button>
    </div>

    <!-- Error State -->
    <div v-if="marketStore.error" class="bg-red-50 dark:bg-red-950/20 border border-red-200 dark:border-red-800/40 rounded-xl p-4 mb-6 text-sm text-red-600 dark:text-red-400">
      <div class="flex items-center gap-2 font-medium">
        <span>⚠️</span> Gagal memuat data bursa saham
      </div>
      <p class="text-xs mt-1 opacity-90">{{ marketStore.error }}</p>
    </div>

    <!-- Loading Skeleton State (Full Widget) -->
    <div v-if="marketStore.loading && (!marketStore.data || !marketStore.data.stocks)" class="flex flex-col gap-6 animate-pulse">
      <!-- Prominent IHSG skeleton -->
      <div class="h-[320px] rounded-xl border border-slate-100 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/30 p-5">
        <div class="h-6 bg-slate-200 dark:bg-slate-800 rounded w-1/4 mb-4"></div>
        <div class="h-8 bg-slate-200 dark:bg-slate-800 rounded w-1/3 mb-4"></div>
        <div class="h-[180px] bg-slate-200 dark:bg-slate-800 rounded-lg w-full"></div>
      </div>
      <!-- Other stocks skeleton -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div v-for="n in 3" :key="'skeleton-market-'+n" class="h-[280px] rounded-xl border border-slate-100 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/30 p-5">
          <div class="h-5 bg-slate-200 dark:bg-slate-800 rounded w-1/3 mb-2"></div>
          <div class="h-6 bg-slate-200 dark:bg-slate-800 rounded w-1/2 mb-4"></div>
          <div class="h-[140px] bg-slate-200 dark:bg-slate-800 rounded-lg w-full"></div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="!marketStore.data || !marketStore.data.stocks || marketStore.data.stocks.length === 0" class="py-12 text-center text-slate-400 dark:text-slate-500">
      Tidak ada data bursa saham tersedia
    </div>

    <!-- Main Content -->
    <div v-else class="flex flex-col gap-6">
      
      <!-- 1. Prominent IHSG Card (Wide Layout) -->
      <div 
        v-if="ihsgStock" 
        class="rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/40 p-5 md:p-6 shadow-sm hover:shadow-md transition-all duration-300 flex flex-col md:flex-row gap-6 justify-between items-stretch"
      >
        <!-- IHSG Info (Left side on desktop) -->
        <div class="flex flex-col justify-between md:w-[30%]">
          <div>
            <div class="flex items-center gap-2 mb-1">
              <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-bold bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400">
                Indeks Utama
              </span>
              <span class="text-[10px] text-slate-400 dark:text-slate-500 font-semibold uppercase font-mono tracking-wider">
                {{ ihsgStock.symbol }}
              </span>
            </div>
            
            <h4 class="font-extrabold text-2xl text-slate-900 dark:text-slate-100 tracking-tight">
              {{ ihsgStock.name }}
            </h4>
            <p class="text-xs text-slate-400 dark:text-slate-500 mt-0.5">
              Indeks Harga Saham Gabungan
            </p>
          </div>

          <div class="my-4 md:my-0">
            <!-- Index Value -->
            <div class="text-4xl font-black text-slate-900 dark:text-slate-100 tracking-tight font-mono">
              {{ formatMarketPrice(ihsgStock.price, ihsgStock.symbol) }}
            </div>
            
            <!-- Daily Change Badge -->
            <div class="mt-2">
              <span 
                :class="[
                  'inline-flex items-center gap-1.5 px-3 py-1 rounded-lg text-sm font-bold shadow-xs',
                  ihsgStock.changePercent >= 0 
                    ? 'bg-emerald-50 dark:bg-emerald-950/20 text-emerald-600 dark:text-emerald-400' 
                    : 'bg-rose-50 dark:bg-rose-950/20 text-rose-600 dark:text-rose-400'
                ]"
              >
                <span class="text-xs">{{ ihsgStock.changePercent >= 0 ? '▲' : '▼' }}</span>
                <span>{{ ihsgStock.changePercent >= 0 ? '+' : '' }}{{ ihsgStock.changePercent.toFixed(2) }}%</span>
                <span class="opacity-80 text-xs">({{ ihsgStock.change >= 0 ? '+' : '' }}{{ ihsgStock.change.toFixed(2) }} pt)</span>
              </span>
            </div>
          </div>

          <div class="text-[10px] text-slate-400 dark:text-slate-500 font-medium">
            Diperbarui: {{ ihsgStock.lastUpdated ? new Date(ihsgStock.lastUpdated).toLocaleTimeString('id-ID') : '-' }}
          </div>
        </div>

        <!-- IHSG LineChart (Right side on desktop) -->
        <div class="flex-1 min-h-[220px]">
          <p class="text-[10px] text-slate-400 dark:text-slate-500 font-bold uppercase tracking-wider mb-2">
            Performa IHSG 7 Hari Terakhir
          </p>
          <div class="h-[220px] w-full">
            <LineChart
              v-bind="getChartProps(ihsgStock.history, ihsgStock.name)"
              :loading="marketStore.loading"
              unit="pt"
              height="100%"
            />
          </div>
        </div>
      </div>

      <!-- 2. Other Stocks Grid (BCA, BRI, Telkom) -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div
          v-for="stock in otherStocks"
          :key="stock.symbol"
          class="rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/40 p-5 flex flex-col justify-between shadow-sm hover:shadow-md transition-all duration-300"
        >
          <!-- Stock Info -->
          <div class="mb-4">
            <div class="flex items-center justify-between mb-1">
              <div>
                <h5 class="font-bold text-base text-slate-900 dark:text-slate-100 tracking-tight">
                  {{ stock.name }}
                </h5>
                <p class="text-[10px] text-slate-400 dark:text-slate-500 font-semibold uppercase font-mono tracking-wider">
                  {{ stock.symbol }}
                </p>
              </div>

              <!-- Change Badge -->
              <span 
                :class="[
                  'inline-flex items-center px-2 py-0.5 rounded text-xs font-bold',
                  stock.changePercent >= 0 
                    ? 'bg-emerald-50 dark:bg-emerald-950/20 text-emerald-600 dark:text-emerald-400' 
                    : 'bg-rose-50 dark:bg-rose-950/20 text-rose-600 dark:text-rose-400'
                ]"
              >
                {{ stock.changePercent >= 0 ? '+' : '' }}{{ stock.changePercent.toFixed(2) }}%
              </span>
            </div>

            <!-- Price -->
            <div class="text-xl font-bold text-slate-900 dark:text-slate-100 tracking-tight font-mono">
              {{ formatMarketPrice(stock.price, stock.symbol) }}
            </div>
          </div>

          <!-- Trend Chart -->
          <div class="w-full mt-2">
            <p class="text-[10px] text-slate-400 dark:text-slate-500 font-bold uppercase tracking-wider mb-2">
              Tren 7 Hari Terakhir
            </p>
            <div class="h-[140px] w-full">
              <AreaChart
                v-bind="getChartProps(stock.history, stock.name)"
                :loading="marketStore.loading"
                unit="Rp"
                height="100%"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
