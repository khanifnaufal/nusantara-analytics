<script setup lang="ts">
import { computed } from 'vue'
import { useMarketStore } from '~/stores/market'
import { usePolling } from '~/composables/usePolling'
import { useExportCsv } from '~/composables/useExportCsv'
import LineChart from '../charts/LineChart.vue'
import AreaChart from '../charts/AreaChart.vue'

const props = withDefaults(defineProps<{
  disableAutoPoll?: boolean
}>(), {
  disableAutoPoll: false
})

const marketStore = useMarketStore()

// Auto-fetch and poll market data every 60s
if (!props.disableAutoPoll) {
  usePolling(() => marketStore.fetchMarket(), 60000)
}

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

// IHSG dynamic Y-axis min/max scaling (tight range to highlight small changes)
const ihsgYMinFn = (value: any) => value.min * 0.998;
const ihsgYMaxFn = (value: any) => value.max * 1.002;

// Other stocks dynamic Y-axis min/max scaling with 2% padding (auto scaled, not starting from zero)
const stockYMinFn = (value: any) => {
  const range = value.max - value.min;
  const minVal = value.min - range * 0.02;
  return value.min >= 0 && minVal < 0 ? 0 : minVal;
};

const stockYMaxFn = (value: any) => {
  const range = value.max - value.min;
  return value.max + range * 0.02;
};
</script>

<template>
  <div class="widget-card">

    <!-- Header -->
    <div class="flex items-center justify-between mb-6 mt-1">
      <div>
        <h3 class="text-base font-bold text-text-primary flex items-center gap-2">
          <span class="text-violet-400">🏛️</span> Bursa Saham Indonesia
        </h3>
        <p class="text-xs text-text-tertiary mt-1">
          Indeks Harga Saham Gabungan (IHSG) & Saham Unggulan • Update terakhir: {{ marketStore.data?.lastUpdated ? new Date(marketStore.data.lastUpdated).toLocaleTimeString('id-ID') : '-' }}
        </p>
      </div>

      <!-- Export Button -->
      <button
        @click="exportCsv()"
        :disabled="!marketStore.data?.stocks || marketStore.data.stocks.length === 0"
        class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-lg bg-transparent border border-white/10 hover:border-white/20 text-zinc-400 hover:text-white cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed transition-all"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
        </svg>
        Export CSV
      </button>
    </div>

    <!-- Error State -->
    <div v-if="marketStore.error" class="bg-rose-500/10 border border-rose-500/20 rounded-xl p-4 mb-6 text-xs text-rose-400">
      <div class="flex items-center gap-2 font-semibold">
        <span>⚠️</span> Gagal memuat data bursa saham
      </div>
      <p class="mt-1 opacity-90 font-mono">{{ marketStore.error }}</p>
    </div>

    <!-- Loading Skeleton State (Full Widget) -->
    <div v-if="marketStore.loading && (!marketStore.data || !marketStore.data.stocks)" class="flex flex-col gap-6 animate-pulse">
      <!-- Prominent IHSG skeleton -->
      <div class="h-[230px] rounded-xl border border-white/5 bg-bg-subcard p-4 flex flex-col md:flex-row gap-4 justify-between">
        <div class="md:w-1/3 flex flex-col justify-between">
          <div>
            <div class="h-4 bg-white/5 rounded w-1/2 mb-2"></div>
            <div class="h-6 bg-white/5 rounded w-2/3"></div>
          </div>
          <div class="h-5 bg-white/5 rounded w-1/3 mt-4"></div>
        </div>
        <div class="flex-1 bg-white/5 rounded-lg h-full"></div>
      </div>
      <!-- Other stocks skeleton -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div v-for="n in 3" :key="'skeleton-market-'+n" class="h-[190px] rounded-xl border border-white/5 bg-bg-subcard p-4 flex flex-col justify-between">
          <div>
            <div class="h-4 bg-white/5 rounded w-1/3 mb-2"></div>
            <div class="h-5 bg-white/5 rounded w-1/2"></div>
          </div>
          <div class="h-[90px] bg-white/5 rounded-lg w-full"></div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="!marketStore.data || !marketStore.data.stocks || marketStore.data.stocks.length === 0" class="py-12 text-center text-text-tertiary">
      Tidak ada data bursa saham tersedia
    </div>

    <!-- Main Content -->
    <div v-else class="flex flex-col gap-6">
      
      <!-- 1. Prominent IHSG Card (Wide Layout) -->
      <div 
        v-if="ihsgStock" 
        class="rounded-xl border border-white/5 bg-[#161616] p-4 md:p-5 flex flex-col md:flex-row gap-6 justify-between items-center transition-all duration-300 hover:border-blue-500/20"
      >
        <!-- IHSG Info (Left side on desktop) -->
        <div class="flex flex-col justify-between w-full md:w-[35%] self-stretch">
          <div>
            <div class="flex items-center gap-2 mb-1">
              <span class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-bold bg-white/5 text-violet-400 border border-violet-500/10">
                Indeks Utama
              </span>
              <span class="text-[9px] text-text-tertiary font-semibold uppercase font-mono tracking-wider">
                {{ ihsgStock.symbol }}
              </span>
            </div>
            
            <h4 class="font-bold text-lg text-text-primary tracking-tight">
              {{ ihsgStock.name }}
            </h4>
          </div>

          <div class="my-3 md:my-1">
            <!-- Index Value -->
            <div class="text-2xl font-extrabold text-text-primary tracking-tight font-mono">
              {{ formatMarketPrice(ihsgStock.price, ihsgStock.symbol) }}
            </div>
            
             <!-- Daily Change Badge -->
            <div class="mt-1.5">
              <span 
                :class="[
                  'inline-flex items-center gap-1.5 px-2 py-0.5 rounded text-[11px] font-bold font-mono',
                  ihsgStock.changePercent > 0 
                    ? 'bg-green-500/15 text-green-400' 
                    : ihsgStock.changePercent < 0
                      ? 'bg-red-500/15 text-red-400'
                      : 'bg-zinc-800 text-zinc-500'
                ]"
              >
                <span>{{ ihsgStock.changePercent >= 0 ? '▲' : '▼' }}</span>
                <span>{{ ihsgStock.changePercent > 0 ? '+' : '' }}{{ ihsgStock.changePercent.toFixed(2) }}%</span>
                <span class="opacity-80 font-normal">({{ ihsgStock.change >= 0 ? '+' : '' }}{{ ihsgStock.change.toFixed(2) }} pt)</span>
              </span>
            </div>
          </div>

          <div class="text-[9px] text-text-tertiary">
            Diperbarui: {{ ihsgStock.lastUpdated ? new Date(ihsgStock.lastUpdated).toLocaleTimeString('id-ID') : '-' }}
          </div>
        </div>

        <!-- IHSG LineChart (Right side on desktop) -->
        <div class="w-full md:flex-1 h-[150px] self-center">
          <p class="text-[9px] text-text-tertiary font-bold uppercase tracking-wider mb-2">
            Performa IHSG 7 Hari Terakhir
          </p>
          <div class="h-[130px] w-full">
            <LineChart
              v-bind="getChartProps(ihsgStock.history, ihsgStock.name)"
              :loading="marketStore.loading"
              unit="pt"
              accentColor="#3B82F6"
              height="100%"
              :yAxisMin="ihsgYMinFn"
              :yAxisMax="ihsgYMaxFn"
            />
          </div>
        </div>
      </div>

      <!-- 2. Other Stocks Grid (BCA, BRI, Telkom) -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div
          v-for="stock in otherStocks"
          :key="stock.symbol"
          class="rounded-xl border border-white/5 bg-[#161616] p-4 flex flex-col justify-between transition-all duration-300 hover:border-blue-500/20"
        >
          <!-- Stock Info -->
          <div class="mb-2">
            <div class="flex items-center justify-between mb-1">
              <div>
                <h5 class="font-bold text-sm text-text-primary tracking-tight">
                  {{ stock.name }}
                </h5>
                <p class="text-[9px] text-text-tertiary font-semibold uppercase font-mono tracking-wider">
                  {{ stock.symbol }}
                </p>
              </div>

              <!-- Change Badge -->
              <span 
                :class="[
                  'inline-flex items-center px-2 py-0.5 rounded text-[10px] font-bold font-mono',
                  stock.changePercent > 0 
                    ? 'bg-green-500/15 text-green-400' 
                    : stock.changePercent < 0
                      ? 'bg-red-500/15 text-red-400'
                      : 'bg-zinc-800 text-zinc-500'
                ]"
              >
                {{ stock.changePercent > 0 ? '+' : '' }}{{ stock.changePercent.toFixed(2) }}%
              </span>
            </div>

            <!-- Price -->
            <div class="text-lg font-extrabold text-text-primary tracking-tight font-mono mt-1">
              {{ formatMarketPrice(stock.price, stock.symbol) }}
            </div>
          </div>

          <!-- Trend Chart -->
          <div class="w-full mt-2 border-t border-white/5 pt-3">
            <p class="text-[9px] text-text-tertiary font-bold uppercase tracking-wider mb-2">
              Tren 7 Hari Terakhir
            </p>
            <div class="h-[90px] w-full">
              <AreaChart
                v-bind="getChartProps(stock.history, stock.name)"
                :loading="marketStore.loading"
                unit="Rp"
                accentColor="#3B82F6"
                height="100%"
                :yAxisMin="stockYMinFn"
                :yAxisMax="stockYMaxFn"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
