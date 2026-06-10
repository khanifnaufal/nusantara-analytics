<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from '#app'
import { useDark, useClipboard } from '@vueuse/core'
import { useRatesStore } from '~/stores/rates'
import { useWeatherStore } from '~/stores/weather'
import { useCommoditiesStore } from '~/stores/commodities'
import { useMarketStore } from '~/stores/market'
import { useExportCsv } from '~/composables/useExportCsv'
import BaseChart from '../charts/BaseChart.vue'

// Initialize stores
const ratesStore = useRatesStore()
const weatherStore = useWeatherStore()
const commoditiesStore = useCommoditiesStore()
const marketStore = useMarketStore()

// Clipboard utility from VueUse
const { copy, copied } = useClipboard()

// Form selections state (draft state)
const selectedMetric = ref('rates')
const selectedItem = ref('USD')
const selectedWeatherMetric = ref('temperature')
const selectedChartType = ref('line')
const selectedRange = ref(7)

// Applied state (used to render the actual chart)
const appliedConfig = ref({
  metric: 'rates',
  item: 'USD',
  weatherMetric: 'temperature',
  chartType: 'line',
  range: 7
})

// Options metadata mappings
const METRIC_CATEGORIES = [
  { value: 'rates', label: '💵 Kurs Mata Uang (IDR)' },
  { value: 'commodities', label: '📈 Komoditas Global' },
  { value: 'market', label: '🏛️ Bursa Saham' },
  { value: 'weather', label: '🌤️ Cuaca Kota' }
]

const CURRENCIES = [
  { value: 'USD', name: 'US Dollar', flag: '🇺🇸' },
  { value: 'EUR', name: 'Euro', flag: '🇪🇺' },
  { value: 'SGD', name: 'Singapore Dollar', flag: '🇸🇬' },
  { value: 'JPY', name: 'Japanese Yen', flag: '🇯🇵' },
  { value: 'MYR', name: 'Malaysian Ringgit', flag: '🇲🇾' },
  { value: 'SAR', name: 'Saudi Riyal', flag: '🇸🇦' },
  { value: 'AUD', name: 'Australian Dollar', flag: '🇦🇺' }
]

const COMMODITIES = [
  { value: 'GC=F', name: 'Emas (Gold)', flag: '🪙' },
  { value: 'CL=F', name: 'Minyak Mentah (Crude Oil)', flag: '🛢️' },
  { value: 'PALM.KL', name: 'CPO (Palm Oil)', flag: '🌴' }
]

const MARKET_STOCKS = [
  { value: '^JKSE', name: 'IHSG (Indeks Gabungan)', flag: '🇮🇩' },
  { value: 'BBCA.JK', name: 'Bank Central Asia', flag: '🏦' },
  { value: 'BBRI.JK', name: 'Bank Rakyat Indonesia', flag: '🏦' },
  { value: 'TLKM.JK', name: 'Telkom Indonesia', flag: '📞' }
]

const CITIES = [
  { value: 'Jakarta', name: 'Jakarta', flag: '🏙️' },
  { value: 'Surabaya', name: 'Surabaya', flag: '⚓' },
  { value: 'Bandung', name: 'Bandung', flag: '⛰️' },
  { value: 'Medan', name: 'Medan', flag: '🏢' },
  { value: 'Semarang', name: 'Semarang', flag: '🌅' },
  { value: 'Makassar', name: 'Makassar', flag: '🏖️' }
]

const WEATHER_METRICS = [
  { value: 'temperature', label: 'Suhu', unit: '°C' },
  { value: 'humidity', label: 'Kelembapan', unit: '%' },
  { value: 'windSpeed', label: 'Kecepatan Angin', unit: 'km/h' }
]

// Computed dropdown options for the secondary item selector
const itemOptions = computed(() => {
  switch (selectedMetric.value) {
    case 'rates':
      return CURRENCIES.map(c => ({ value: c.value, label: `${c.flag} ${c.value} - ${c.name}` }))
    case 'commodities':
      return COMMODITIES.map(c => ({ value: c.value, label: `${c.flag} ${c.name}` }))
    case 'market':
      return MARKET_STOCKS.map(m => ({ value: m.value, label: `${m.flag} ${m.name}` }))
    case 'weather':
      return CITIES.map(c => ({ value: c.value, label: `${c.flag} ${c.name}` }))
    default:
      return []
  }
})

// Triggered when metric category dropdown changes
const onMetricChange = () => {
  if (selectedMetric.value === 'rates') {
    selectedItem.value = 'USD'
  } else if (selectedMetric.value === 'commodities') {
    selectedItem.value = 'GC=F'
  } else if (selectedMetric.value === 'market') {
    selectedItem.value = '^JKSE'
  } else if (selectedMetric.value === 'weather') {
    selectedItem.value = 'Jakarta'
  }
}

// Deterministic hashing helper
const getDeterministicHash = (str: string): number => {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = (hash << 5) - hash + str.charCodeAt(i)
    hash |= 0
  }
  return hash
}

// Generate daily changes based on currency & date
const getDailyChangePercent = (currency: string, dateStr: string): number => {
  const hash = getDeterministicHash(currency + dateStr)
  return (hash % 100) / 150 // between -0.66% and +0.66%
}

// Rates history generation
const generateRatesHistory = (currency: string, currentValInIdr: number, rangeDays: number) => {
  const tempHistory = []
  const today = new Date()
  let val = currentValInIdr
  
  for (let i = 0; i < rangeDays; i++) {
    const d = new Date()
    d.setDate(today.getDate() - i)
    const dateStr = d.toISOString().split('T')[0] || ''
    
    tempHistory.push({
      date: dateStr,
      value: val
    })
    
    const changePercent = getDailyChangePercent(currency, dateStr)
    const changeRatio = changePercent / 100
    val = val / (1 + changeRatio)
  }
  
  return tempHistory.reverse()
}

// Weather history generation
const generateWeatherHistory = (
  city: string, 
  metric: 'temperature' | 'humidity' | 'windSpeed', 
  currentVal: number, 
  rangeDays: number
) => {
  const tempHistory = []
  const today = new Date()
  let val = currentVal
  
  for (let i = 0; i < rangeDays; i++) {
    const d = new Date()
    d.setDate(today.getDate() - i)
    const dateStr = d.toISOString().split('T')[0] || ''
    
    tempHistory.push({
      date: dateStr,
      value: Number(val.toFixed(1))
    })
    
    const hash = getDeterministicHash(city + metric + dateStr)
    if (metric === 'temperature') {
      const change = (hash % 150) / 100 // -1.5 to +1.5
      val = val - change
    } else if (metric === 'humidity') {
      const change = hash % 6 // -5 to +5
      val = Math.max(30, Math.min(100, val - change))
    } else { // windSpeed
      const change = (hash % 300) / 100 // -3 to +3
      val = Math.max(0, val - change)
    }
  }
  
  return tempHistory.reverse()
}

// Finance fallback history generation
const generateDeterministicFinanceHistory = (symbol: string, currentPrice: number, rangeDays: number) => {
  const tempHistory = []
  const today = new Date()
  let val = currentPrice
  
  for (let i = 0; i < rangeDays; i++) {
    const d = new Date()
    d.setDate(today.getDate() - i)
    const dateStr = d.toISOString().split('T')[0] || ''
    
    tempHistory.push({
      date: dateStr,
      value: Number(val.toFixed(2))
    })
    
    const hash = getDeterministicHash(symbol + dateStr)
    const changePercent = (hash % 120) / 100 // -1.2% to +1.2%
    const changeRatio = changePercent / 100
    val = val / (1 + changeRatio)
  }
  
  return tempHistory.reverse()
}

// Extract historical quote data
const getQuoteHistory = (quote: any, rangeDays: number) => {
  if (quote && quote.history && quote.history.length > 0) {
    const sortedHistory = [...quote.history].sort(
      (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
    )
    const sliceStart = Math.max(0, sortedHistory.length - rangeDays)
    return sortedHistory.slice(sliceStart).map((h: any) => {
      const d = new Date(h.timestamp)
      return {
        date: d.toISOString().split('T')[0] || '',
        value: parseFloat(h.close.toFixed(2))
      }
    })
  } else if (quote) {
    return generateDeterministicFinanceHistory(quote.symbol, quote.price, rangeDays)
  }
  return []
}

// Reactive store checking loading state
const isLoading = computed(() => {
  const config = appliedConfig.value
  if (config.metric === 'rates') return ratesStore.loading || !ratesStore.data
  if (config.metric === 'commodities') return commoditiesStore.loading || !commoditiesStore.data
  if (config.metric === 'market') return marketStore.loading || !marketStore.data
  if (config.metric === 'weather') return weatherStore.loading || !weatherStore.data
  return false
})

// Primary computed chart data binding
const chartData = computed(() => {
  const config = appliedConfig.value
  const rangeDays = config.range
  
  let xAxis: string[] = []
  let dataPoints: number[] = []
  let seriesName = ''
  let unit = ''
  let color = '#10b981'
  
  if (config.metric === 'rates') {
    seriesName = `Kurs ${config.item}`
    unit = 'IDR'
    color = '#10b981'
    
    const rawVal = ratesStore.data?.rates[config.item]
    if (rawVal !== undefined && rawVal !== null) {
      const inIdr = 1 / rawVal
      const generated = generateRatesHistory(config.item, inIdr, rangeDays)
      xAxis = generated.map(g => {
        const d = new Date(g.date)
        return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short' })
      })
      dataPoints = generated.map(g => parseFloat(g.value.toFixed(2)))
    }
  } else if (config.metric === 'commodities') {
    const item = commoditiesStore.data?.commodities.find(c => c.symbol === config.item)
    if (item) {
      seriesName = item.name
      unit = item.currency
      color = '#f59e0b'
      
      const history = getQuoteHistory(item, rangeDays)
      xAxis = history.map(h => {
        const d = new Date(h.date)
        return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short' })
      })
      dataPoints = history.map(h => h.value)
    }
  } else if (config.metric === 'market') {
    const item = marketStore.data?.stocks.find(s => s.symbol === config.item)
    if (item) {
      seriesName = item.name
      unit = item.symbol === '^JKSE' ? 'pt' : 'Rp'
      color = '#6366f1'
      
      const history = getQuoteHistory(item, rangeDays)
      xAxis = history.map(h => {
        const d = new Date(h.date)
        return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short' })
      })
      dataPoints = history.map(h => h.value)
    }
  } else if (config.metric === 'weather') {
    const item = weatherStore.data?.cities.find(c => c.city === config.item)
    if (item) {
      const metricObj = WEATHER_METRICS.find(m => m.value === config.weatherMetric)
      seriesName = `${item.city} - ${metricObj?.label}`
      unit = metricObj?.unit || ''
      color = '#06b6d4'
      
      const currentVal = config.weatherMetric === 'temperature'
        ? item.temperature
        : config.weatherMetric === 'humidity'
          ? item.humidity
          : item.windSpeed
          
      const generated = generateWeatherHistory(config.item, config.weatherMetric as 'temperature' | 'humidity' | 'windSpeed', currentVal, rangeDays)
      xAxis = generated.map(g => {
        const d = new Date(g.date)
        return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short' })
      })
      dataPoints = generated.map(g => g.value)
    }
  }
  
  return { xAxis, dataPoints, seriesName, unit, color }
})

// Dynamic title formatting for the visual display
const activeTitle = computed(() => {
  const config = appliedConfig.value
  let itemName = config.item
  
  if (config.metric === 'rates') {
    const c = CURRENCIES.find(curr => curr.value === config.item)
    itemName = c ? `${c.flag} ${c.value} ke IDR` : config.item
  } else if (config.metric === 'commodities') {
    const c = COMMODITIES.find(comm => comm.value === config.item)
    itemName = c ? `${c.flag} ${c.name}` : config.item
  } else if (config.metric === 'market') {
    const m = MARKET_STOCKS.find(st => st.value === config.item)
    itemName = m ? `${m.flag} ${m.name}` : config.item
  } else if (config.metric === 'weather') {
    const c = CITIES.find(ct => ct.value === config.item)
    const m = WEATHER_METRICS.find(met => met.value === config.weatherMetric)
    itemName = `${c ? c.flag : ''} ${config.item} (${m ? m.label : ''})`
  }
  
  return `${itemName} — Tren ${config.range} Hari`
})

// ECharts reactive option configurations
const chartOption = computed(() => {
  const chart = chartData.value
  const config = appliedConfig.value
  const unitSuffix = chart.unit ? ' ' + chart.unit : ''
  
  return {
    backgroundColor: 'transparent',
    color: [chart.color],
    tooltip: {
      trigger: 'axis',
      confine: true,
      backgroundColor: 'var(--tooltip-bg)',
      borderColor: 'var(--tooltip-border)',
      borderWidth: 1,
      shadowColor: 'rgba(0, 0, 0, 0.05)',
      shadowBlur: 10,
      formatter: (params: any) => {
        if (!params || params.length === 0) return ''
        const title = params[0].axisValueLabel || params[0].name
        let result = `<div style="font-family: Inter, sans-serif; font-size: 12px; line-height: 1.5; padding: 4px;">`
        result += `<div style="font-weight: 600; margin-bottom: 6px; color: var(--tooltip-title-color, #8B949E);">${title}</div>`
        params.forEach((item: any) => {
          let val = '-'
          if (item.value !== undefined && item.value !== null) {
            if (chart.unit === 'Rp') {
              val = new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(item.value)
            } else if (chart.unit === 'USD') {
              val = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 2 }).format(item.value)
            } else {
              val = new Intl.NumberFormat('id-ID').format(item.value) + unitSuffix
            }
          }
          result += `<div style="display: flex; align-items: center; justify-content: space-between; gap: 16px; margin-top: 4px;">
            <div style="display: flex; align-items: center; gap: 6px;">
              ${item.marker}
              <span style="color: var(--tooltip-text-color, #8B949E); font-weight: 500;">${item.seriesName}</span>
            </div>
            <span style="font-weight: 600; color: var(--tooltip-val-color, #ffffff);">${val}</span>
          </div>`
        })
        result += `</div>`
        return result
      }
    },
    grid: {
      top: 25,
      left: '2%',
      right: '2%',
      bottom: 15,
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: config.chartType === 'bar',
      data: chart.xAxis,
      axisLine: {
        show: true,
        lineStyle: {
          color: 'rgba(255, 255, 255, 0.08)'
        }
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        fontFamily: 'Inter, sans-serif',
        fontSize: 10,
        color: '#8B949E'
      }
    },
    yAxis: {
      type: 'value',
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        fontFamily: 'Inter, sans-serif',
        fontSize: 10,
        color: '#8B949E',
        formatter: (value: number) => {
          if (chart.unit === 'Rp') {
            return new Intl.NumberFormat('id-ID', { notation: 'compact', style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(value)
          } else if (chart.unit === 'USD') {
            return new Intl.NumberFormat('en-US', { notation: 'compact', style: 'currency', currency: 'USD', maximumFractionDigits: 1 }).format(value)
          }
          return new Intl.NumberFormat(undefined, { notation: 'compact' }).format(value) + (chart.unit ? ' ' + chart.unit : '')
        }
      },
      splitLine: {
        lineStyle: {
          color: 'rgba(255, 255, 255, 0.04)',
          type: 'dashed'
        }
      }
    },
    series: [
      {
        name: chart.seriesName,
        type: config.chartType === 'bar' ? 'bar' : 'line',
        data: chart.dataPoints,
        smooth: config.chartType !== 'bar',
        showSymbol: false,
        symbolSize: 6,
        barMaxWidth: config.chartType === 'bar' ? 24 : undefined,
        itemStyle: config.chartType === 'bar' ? {
          borderRadius: [4, 4, 0, 0]
        } : undefined,
        areaStyle: config.chartType === 'area' ? {
          opacity: 0.15,
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: chart.color },
              { offset: 1, color: 'transparent' }
            ],
            global: false
          }
        } : undefined,
        emphasis: {
          focus: 'series',
          itemStyle: {
            borderWidth: 2
          }
        },
        lineStyle: config.chartType !== 'bar' ? {
          width: config.chartType === 'area' ? 2 : 2.5
        } : undefined
      }
    ]
  }
})

// Apply Selection values to applied config & sync query params
const applySettings = () => {
  appliedConfig.value = {
    metric: selectedMetric.value,
    item: selectedItem.value,
    weatherMetric: selectedWeatherMetric.value,
    chartType: selectedChartType.value,
    range: Number(selectedRange.value)
  }
  
  // Update query params in URL
  const route = useRoute()
  const router = useRouter()
  const query: Record<string, string> = {
    metric: selectedMetric.value,
    type: selectedChartType.value,
    range: String(selectedRange.value)
  }
  
  if (selectedMetric.value === 'rates') {
    query.currency = selectedItem.value
  } else if (selectedMetric.value === 'commodities' || selectedMetric.value === 'market') {
    query.symbol = selectedItem.value
  } else if (selectedMetric.value === 'weather') {
    query.city = selectedItem.value
    query.weatherMetric = selectedWeatherMetric.value
  }
  
  router.replace({ query })
}

// Copy URL link handler
const copyLink = () => {
  if (typeof window !== 'undefined') {
    copy(window.location.href)
  }
}

// Setup CSV Exporter
const csvData = computed(() => {
  const chart = chartData.value
  if (!chart || chart.xAxis.length === 0) return []
  
  const headerValue = `${chart.seriesName} (${chart.unit})`
  return chart.xAxis.map((date, idx) => ({
    'Tanggal': date,
    [headerValue]: chart.dataPoints[idx]
  }))
})

const exportFilename = computed(() => {
  const config = appliedConfig.value
  let name = `custom_viz_${config.metric}`
  if (config.metric === 'rates') {
    name += `_${config.item}`
  } else if (config.metric === 'commodities' || config.metric === 'market') {
    name += `_${config.item.replace('^', '')}`
  } else if (config.metric === 'weather') {
    name += `_${config.item}_${config.weatherMetric}`
  }
  return name + `_${config.range}d`
})

const { exportCsv } = useExportCsv(csvData, exportFilename)

// Read URL query params on mount
onMounted(async () => {
  const route = useRoute()
  const q = route.query
  
  if (q.metric) {
    selectedMetric.value = String(q.metric)
    if (q.metric === 'rates' && q.currency) {
      selectedItem.value = String(q.currency)
    } else if ((q.metric === 'commodities' || q.metric === 'market') && q.symbol) {
      selectedItem.value = String(q.symbol)
    } else if (q.metric === 'weather' && q.city) {
      selectedItem.value = String(q.city)
      if (q.weatherMetric) {
        selectedWeatherMetric.value = String(q.weatherMetric)
      }
    }
  }
  
  if (q.type) {
    selectedChartType.value = String(q.type)
  }
  
  if (q.range) {
    selectedRange.value = Number(q.range)
  }
  
  // Proactively fetch stores if they are empty
  try {
    const promises = []
    if (!ratesStore.data) promises.push(ratesStore.fetchRates())
    if (!weatherStore.data) promises.push(weatherStore.fetchWeather())
    if (!commoditiesStore.data) promises.push(commoditiesStore.fetchCommodities())
    if (!marketStore.data) promises.push(marketStore.fetchMarket())
    if (promises.length > 0) {
      await Promise.all(promises)
    }
  } catch (err) {
    console.error('Error fetching data during builder mount:', err)
  }
  
  // Finalize options mapping
  applySettings()
})
</script>

<template>
  <div class="relative w-full overflow-hidden rounded-2xl border border-white/6 bg-[#161B22] p-5 md:p-6 transition-all duration-300 hover:border-emerald-500/30 hover:shadow-[0_0_20px_rgba(16,185,129,0.08)]">
    <!-- Gradient Accent Top Border -->
    <div class="absolute top-0 left-0 h-1 w-full bg-gradient-to-r from-emerald-500 to-indigo-500"></div>

    <!-- Header area -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6 mt-1">
      <div>
        <h3 class="text-base font-bold text-text-primary flex items-center gap-2">
          <span>🎨</span> Custom Visualisasi Builder
        </h3>
        <p class="text-xs text-text-tertiary mt-1">
          Pilih metrik, atur tipe grafik dan rentang waktu sesuai keinginan Anda.
        </p>
      </div>

      <!-- Action buttons -->
      <div class="flex items-center gap-2 self-start sm:self-auto">
        <button
          @click="copyLink()"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-lg border border-white/5 bg-[#1C2128] hover:bg-[#21262D] text-text-secondary hover:text-text-primary cursor-pointer transition-all"
        >
          <span v-if="copied" class="flex items-center gap-1.5 text-emerald-400 font-bold">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
            </svg>
            Tersalin!
          </span>
          <span v-else class="flex items-center gap-1.5">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M8.684 10.742l-1.922 1.922a4.5 4.5 0 106.364 6.364l3.182-3.182a4.5 4.5 0 00-6.364-6.364L8.682 10.74m.002.002h.008" />
            </svg>
            Salin Link
          </span>
        </button>

        <button
          @click="exportCsv()"
          :disabled="isLoading || !chartData || chartData.dataPoints.length === 0"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-lg border border-white/5 bg-[#1C2128] hover:bg-[#21262D] text-text-secondary hover:text-text-primary cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed transition-all"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
          Export CSV
        </button>
      </div>
    </div>

    <!-- Controls Selector Grid -->
    <div class="grid grid-cols-1 md:grid-cols-12 gap-4 p-4 mb-6 rounded-xl bg-[#1C2128] border border-white/5 transition-all duration-300">
      <!-- 1. Metric Category -->
      <div class="col-span-1 md:col-span-3">
        <label class="block text-[10px] font-bold text-text-tertiary uppercase tracking-wider mb-1.5">
          Kategori Metrik
        </label>
        <select
          v-model="selectedMetric"
          @change="onMetricChange()"
          class="w-full text-xs rounded-lg border border-white/5 bg-[#161B22] px-3 py-2 text-text-primary focus:border-emerald-500 focus:outline-none transition-colors cursor-pointer"
        >
          <option v-for="m in METRIC_CATEGORIES" :key="m.value" :value="m.value">
            {{ m.label }}
          </option>
        </select>
      </div>

      <!-- 2. Sub Option (Item) -->
      <div class="col-span-1" :class="selectedMetric === 'weather' ? 'md:col-span-3' : 'md:col-span-5'">
        <label class="block text-[10px] font-bold text-text-tertiary uppercase tracking-wider mb-1.5">
          Pilih Item
        </label>
        <select
          v-model="selectedItem"
          class="w-full text-xs rounded-lg border border-white/5 bg-[#161B22] px-3 py-2 text-text-primary focus:border-emerald-500 focus:outline-none transition-colors cursor-pointer"
        >
          <option v-for="item in itemOptions" :key="item.value" :value="item.value">
            {{ item.label }}
          </option>
        </select>
      </div>

      <!-- 3. Weather Metric Parameter (Cuaca only) -->
      <div v-if="selectedMetric === 'weather'" class="col-span-1 md:col-span-2">
        <label class="block text-[10px] font-bold text-text-tertiary uppercase tracking-wider mb-1.5">
          Parameter Cuaca
        </label>
        <select
          v-model="selectedWeatherMetric"
          class="w-full text-xs rounded-lg border border-white/5 bg-[#161B22] px-3 py-2 text-text-primary focus:border-emerald-500 focus:outline-none transition-colors cursor-pointer"
        >
          <option v-for="w in WEATHER_METRICS" :key="w.value" :value="w.value">
            {{ w.label }} ({{ w.unit }})
          </option>
        </select>
      </div>

      <!-- 4. Chart Type Selection -->
      <div class="col-span-1 md:col-span-2">
        <label class="block text-[10px] font-bold text-text-tertiary uppercase tracking-wider mb-1.5">
          Tipe Grafik
        </label>
        <select
          v-model="selectedChartType"
          class="w-full text-xs rounded-lg border border-white/5 bg-[#161B22] px-3 py-2 text-text-primary focus:border-emerald-500 focus:outline-none transition-colors cursor-pointer"
        >
          <option value="line">Line Chart</option>
          <option value="bar">Bar Chart</option>
          <option value="area">Area Chart</option>
        </select>
      </div>

      <!-- 5. Time Range Selection -->
      <div class="col-span-1 md:col-span-1.5" :class="selectedMetric === 'weather' ? 'md:col-span-1.5' : 'md:col-span-2'">
        <label class="block text-[10px] font-bold text-text-tertiary uppercase tracking-wider mb-1.5">
          Waktu
        </label>
        <select
          v-model="selectedRange"
          class="w-full text-xs rounded-lg border border-white/5 bg-[#161B22] px-3 py-2 text-text-primary focus:border-emerald-500 focus:outline-none transition-colors cursor-pointer"
        >
          <option :value="3">3 Hari</option>
          <option :value="7">7 Hari</option>
        </select>
      </div>

      <!-- 6. Apply Button -->
      <div class="col-span-1 md:col-span-2 flex items-end">
        <button
          @click="applySettings()"
          class="w-full h-[34px] inline-flex items-center justify-center gap-1.5 px-4 py-2 text-xs font-bold rounded-lg bg-gradient-to-r from-emerald-500 to-teal-600 hover:from-emerald-600 hover:to-teal-700 text-white shadow-md shadow-emerald-500/10 hover:shadow-emerald-500/20 active:scale-[0.98] transition-all cursor-pointer"
        >
          <span>⚡</span> Terapkan
        </button>
      </div>
    </div>

    <!-- Chart Panel -->
    <div class="mt-6">
      <div class="flex items-center justify-between mb-4 border-b border-white/5 pb-3">
        <h4 class="text-xs font-extrabold text-text-secondary uppercase tracking-wide">
          {{ activeTitle }}
        </h4>
        <div class="text-[9px] text-text-tertiary font-semibold uppercase tracking-wider font-mono">
          Y-Axis Unit: <span class="text-text-primary font-bold">{{ chartData.unit || 'n/a' }}</span>
        </div>
      </div>

      <div class="relative w-full h-[320px]">
        <!-- BaseChart rendering with reactive configuration -->
        <BaseChart
          :option="chartOption"
          :loading="isLoading"
          height="100%"
        />
      </div>

      <!-- Compact numeric data preview table -->
      <div v-if="!isLoading && chartData.dataPoints.length > 0" class="mt-6 overflow-x-auto border border-white/5 rounded-xl bg-[#1C2128]/50 scrollbar-thin">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="bg-[#1C2128] border-b border-white/5 text-text-tertiary font-bold tracking-wider uppercase">
              <th class="py-2 px-4">Tanggal</th>
              <th class="py-2 px-4 text-right">Nilai ({{ chartData.unit }})</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-white/5">
            <tr v-for="(val, idx) in chartData.dataPoints" :key="'preview-'+idx" class="hover:bg-white/[0.02] text-text-secondary">
              <td class="py-2 px-4 font-medium">{{ chartData.xAxis[idx] }}</td>
              <td class="py-2 px-4 text-right font-mono font-semibold text-text-primary">{{ val }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Focus border accents */
select:focus {
  border-color: #10B981;
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.1);
}
</style>
