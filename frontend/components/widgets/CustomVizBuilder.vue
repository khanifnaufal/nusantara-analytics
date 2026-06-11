<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from '#app'
import { useDark, useClipboard } from '@vueuse/core'
import { useRatesStore } from '~/stores/rates'
import { useWeatherStore } from '~/stores/weather'
import { useCommoditiesStore } from '~/stores/commodities'
import { useMarketStore } from '~/stores/market'
import { useExportCsv } from '~/composables/useExportCsv'
import BaseChart from '../charts/BaseChart.vue'

// Custom Dropdowns active state
const activeDropdown = ref<string | null>(null)

const toggleDropdown = (name: string) => {
  if (activeDropdown.value === name) {
    activeDropdown.value = null
  } else {
    activeDropdown.value = name
  }
}

const closeAllDropdowns = () => {
  activeDropdown.value = null
}

const selectMetric = (val: string) => {
  selectedMetric.value = val
  onMetricChange()
  activeDropdown.value = null
}

const selectItem = (val: string) => {
  selectedItem.value = val
  activeDropdown.value = null
}

const selectWeatherMetric = (val: string) => {
  selectedWeatherMetric.value = val
  activeDropdown.value = null
}

const selectChartType = (val: string) => {
  selectedChartType.value = val
  activeDropdown.value = null
}

const selectRange = (val: number) => {
  selectedRange.value = val
  activeDropdown.value = null
}

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
  { value: 'USD', name: 'US Dollar', flagCode: 'us' },
  { value: 'EUR', name: 'Euro', flagCode: 'eu' },
  { value: 'SGD', name: 'Singapore Dollar', flagCode: 'sg' },
  { value: 'JPY', name: 'Japanese Yen', flagCode: 'jp' },
  { value: 'MYR', name: 'Malaysian Ringgit', flagCode: 'my' },
  { value: 'SAR', name: 'Saudi Riyal', flagCode: 'sa' },
  { value: 'AUD', name: 'Australian Dollar', flagCode: 'au' }
]

const COMMODITIES = [
  { value: 'GC=F', name: 'Emas (Gold)', icon: '🪙' },
  { value: 'CL=F', name: 'Minyak Mentah (Crude Oil)', icon: '🛢️' },
  { value: 'PALM.KL', name: 'CPO (Palm Oil)', icon: '🌴' }
]

const MARKET_STOCKS = [
  { value: '^JKSE', name: 'IHSG (Indeks Gabungan)', isIndo: true },
  { value: 'BBCA.JK', name: 'Bank Central Asia', isIndo: true },
  { value: 'BBRI.JK', name: 'Bank Rakyat Indonesia', isIndo: true },
  { value: 'TLKM.JK', name: 'Telkom Indonesia', isIndo: true }
]

const CITIES = [
  { value: 'Jakarta', name: 'Jakarta', icon: '🏙️' },
  { value: 'Surabaya', name: 'Surabaya', icon: '⚓' },
  { value: 'Bandung', name: 'Bandung', icon: '⛰️' },
  { value: 'Medan', name: 'Medan', icon: '🏢' },
  { value: 'Semarang', name: 'Semarang', icon: '🌅' },
  { value: 'Makassar', name: 'Makassar', icon: '🏖️' }
]

const WEATHER_METRICS = [
  { value: 'temperature', label: 'Suhu', unit: '°C' },
  { value: 'humidity', label: 'Kelembapan', unit: '%' },
  { value: 'windSpeed', label: 'Kecepatan Angin', unit: 'km/h' }
]

interface ItemDetail {
  value: string
  name: string
  flagCode?: string
  icon?: string
  isIndo?: boolean
}

// Resolved details of the currently selected item
const selectedItemDetails = computed((): ItemDetail => {
  const itemVal = selectedItem.value
  switch (selectedMetric.value) {
    case 'rates':
      return CURRENCIES.find(c => c.value === itemVal) || { value: itemVal, name: itemVal, flagCode: 'un' }
    case 'commodities':
      return COMMODITIES.find(c => c.value === itemVal) || { value: itemVal, name: itemVal, icon: '📈' }
    case 'market':
      return MARKET_STOCKS.find(c => c.value === itemVal) || { value: itemVal, name: itemVal, isIndo: true }
    case 'weather':
      return CITIES.find(c => c.value === itemVal) || { value: itemVal, name: itemVal, icon: '🌤️' }
    default:
      return { value: itemVal, name: itemVal }
  }
})

// Resolved details of the applied item
const appliedItemDetails = computed((): ItemDetail => {
  const itemVal = appliedConfig.value.item
  switch (appliedConfig.value.metric) {
    case 'rates':
      return CURRENCIES.find(c => c.value === itemVal) || { value: itemVal, name: itemVal, flagCode: 'un' }
    case 'commodities':
      return COMMODITIES.find(c => c.value === itemVal) || { value: itemVal, name: itemVal, icon: '📈' }
    case 'market':
      return MARKET_STOCKS.find(c => c.value === itemVal) || { value: itemVal, name: itemVal, isIndo: true }
    case 'weather':
      return CITIES.find(c => c.value === itemVal) || { value: itemVal, name: itemVal, icon: '🌤️' }
    default:
      return { value: itemVal, name: itemVal }
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
  let color = '#3B82F6'
  
  if (config.metric === 'rates') {
    seriesName = `Kurs ${config.item}`
    unit = 'IDR'
    color = '#3B82F6'
    
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
      color = '#3B82F6'
      
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
      color = '#3B82F6'
      
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
      color = '#3B82F6'
      
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
    itemName = c ? `${c.value} ke IDR (${c.name})` : config.item
  } else if (config.metric === 'commodities') {
    const c = COMMODITIES.find(comm => comm.value === config.item)
    itemName = c ? `${c.name}` : config.item
  } else if (config.metric === 'market') {
    const m = MARKET_STOCKS.find(st => st.value === config.item)
    itemName = m ? `${m.name}` : config.item
  } else if (config.metric === 'weather') {
    const c = CITIES.find(ct => ct.value === config.item)
    const m = WEATHER_METRICS.find(met => met.value === config.weatherMetric)
    itemName = `${config.item} (${m ? m.label : ''})`
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
        const dateStr = params[0].axisValueLabel || params[0].name
        let result = `<div style="font-family: Inter, sans-serif; font-size: 12px; line-height: 1.5; padding: 6px 8px;">`
        result += `<div style="font-weight: 700; margin-bottom: 8px; color: var(--tooltip-title-color, #8B949E);">Tanggal: ${dateStr}</div>`
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
            <span style="font-weight: 700; color: var(--tooltip-val-color, #ffffff);">${val}</span>
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
        color: '#52525B'
      }
    },
    yAxis: {
      type: 'value',
      scale: true,
      min: (value: any) => {
        if (!value || typeof value.min !== 'number' || isNaN(value.min) || !isFinite(value.min)) {
          return 'dataMin'
        }
        const range = value.max - value.min
        if (range === 0) return value.min * 0.95
        const minVal = value.min - range * 0.05
        return value.min >= 0 && minVal < 0 ? 0 : minVal
      },
      max: (value: any) => {
        if (!value || typeof value.max !== 'number' || isNaN(value.max) || !isFinite(value.max)) {
          return 'dataMax'
        }
        const range = value.max - value.min
        if (range === 0) return value.max * 1.05
        return value.max + range * 0.05
      },
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        fontFamily: 'Inter, sans-serif',
        fontSize: 10,
        color: '#52525B',
        formatter: (value: number) => {
          const unitStr = chart.unit ? ' ' + chart.unit : ''
          if (chart.unit === 'Rp') {
            const formatted = Math.abs(value) >= 1000000
              ? new Intl.NumberFormat('id-ID', { notation: 'compact', style: 'currency', currency: 'IDR', maximumFractionDigits: 1 }).format(value)
              : new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(value)
            return formatted
          } else if (chart.unit === 'USD') {
            const formatted = Math.abs(value) >= 1000000
              ? new Intl.NumberFormat('en-US', { notation: 'compact', style: 'currency', currency: 'USD', maximumFractionDigits: 1 }).format(value)
              : new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(value)
            return formatted
          }
          const formatted = Math.abs(value) >= 1000000
            ? new Intl.NumberFormat('id-ID', { notation: 'compact', maximumFractionDigits: 1 }).format(value)
            : new Intl.NumberFormat('id-ID', { maximumFractionDigits: 2 }).format(value)
          return formatted + unitStr
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
  
  // Add global click listener to close dropdowns when clicking outside
  if (typeof window !== 'undefined') {
    document.addEventListener('click', closeAllDropdowns)
  }

  // Finalize options mapping
  applySettings()
})

onUnmounted(() => {
  if (typeof window !== 'undefined') {
    document.removeEventListener('click', closeAllDropdowns)
  }
})
</script>

<template>
  <div class="widget-card">

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
          class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-lg bg-transparent border border-white/10 hover:border-white/20 text-zinc-400 hover:text-white cursor-pointer transition-all"
        >
          <span v-if="copied" class="flex items-center gap-1.5 text-blue-400 font-bold">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
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
          class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-lg bg-transparent border border-white/10 hover:border-white/20 text-zinc-400 hover:text-white cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed transition-all"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
          Export CSV
        </button>
      </div>
    </div>

    <!-- Controls Selector Grid -->
    <div class="grid grid-cols-1 md:grid-cols-12 gap-4 p-4 mb-6 rounded-xl bg-bg-subcard border border-white/5 transition-all duration-300">
      <!-- 1. Metric Category -->
      <div class="col-span-1 md:col-span-3 relative">
        <label class="block text-[10px] font-bold text-text-tertiary uppercase tracking-wider mb-1.5">
          Kategori Metrik
        </label>
        <button
          type="button"
          @click.stop="toggleDropdown('metric')"
          class="w-full flex items-center justify-between text-xs rounded-lg border border-white/5 bg-bg-card px-3 py-2 text-text-primary focus:border-blue-500/40 focus:outline-none transition-all cursor-pointer hover:border-white/10 hover:bg-bg-subcard/50 active:scale-[0.99]"
        >
          <span>{{ METRIC_CATEGORIES.find(m => m.value === selectedMetric)?.label }}</span>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-text-secondary transition-transform duration-200" :class="{ 'rotate-180': activeDropdown === 'metric' }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </button>
        <div
          v-if="activeDropdown === 'metric'"
          class="absolute left-0 right-0 mt-1.5 z-50 max-h-60 overflow-y-auto rounded-lg border border-white/10 bg-bg-card/95 backdrop-blur-xl shadow-xl py-1 scrollbar-thin"
        >
          <button
            v-for="m in METRIC_CATEGORIES"
            :key="m.value"
            type="button"
            @click="selectMetric(m.value)"
            class="w-full text-left px-3 py-2 text-xs text-text-secondary hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
          >
            <span>{{ m.label }}</span>
            <span v-if="selectedMetric === m.value" class="text-blue-400 font-bold">✓</span>
          </button>
        </div>
      </div>

      <!-- 2. Sub Option (Item) -->
      <div class="col-span-1 relative" :class="selectedMetric === 'weather' ? 'md:col-span-3' : 'md:col-span-5'">
        <label class="block text-[10px] font-bold text-text-tertiary uppercase tracking-wider mb-1.5">
          Pilih Item
        </label>
        <button
          type="button"
          @click.stop="toggleDropdown('item')"
          class="w-full flex items-center justify-between text-xs rounded-lg border border-white/5 bg-bg-card px-3 py-2 text-text-primary focus:border-blue-500/40 focus:outline-none transition-all cursor-pointer hover:border-white/10 hover:bg-bg-subcard/50 active:scale-[0.99]"
        >
          <span class="flex items-center gap-2">
            <!-- Selected Item flag/icon/details -->
            <template v-if="selectedMetric === 'rates'">
              <img 
                :src="'https://flagcdn.com/w40/' + (selectedItemDetails.flagCode || 'un') + '.png'" 
                class="w-5 h-3.5 object-cover rounded-xs border border-white/10" 
              />
              <span class="font-bold text-text-primary">{{ selectedItemDetails.value }}</span>
              <span class="text-text-tertiary hidden md:inline">— {{ selectedItemDetails.name }}</span>
            </template>
            <template v-else-if="selectedMetric === 'market'">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 3 2" class="w-5 h-3.5 rounded-xs border border-white/10 shrink-0" aria-hidden="true">
                <rect width="3" height="1" fill="#FF0000"/>
                <rect y="1" width="3" height="1" fill="#FFFFFF"/>
              </svg>
              <span class="font-bold text-text-primary">{{ selectedItemDetails.name }}</span>
              <span class="text-text-tertiary font-mono text-[9px] uppercase">({{ selectedItemDetails.value }})</span>
            </template>
            <template v-else-if="selectedMetric === 'commodities'">
              <span class="text-sm">{{ selectedItemDetails.icon }}</span>
              <span class="font-bold text-text-primary">{{ selectedItemDetails.name }}</span>
            </template>
            <template v-else-if="selectedMetric === 'weather'">
              <span class="text-sm">{{ selectedItemDetails.icon }}</span>
              <span class="font-bold text-text-primary">{{ selectedItemDetails.name }}</span>
            </template>
          </span>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-text-secondary transition-transform duration-200" :class="{ 'rotate-180': activeDropdown === 'item' }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </button>
        <div
          v-if="activeDropdown === 'item'"
          class="absolute left-0 right-0 mt-1.5 z-50 max-h-60 overflow-y-auto rounded-lg border border-white/10 bg-bg-card/95 backdrop-blur-xl shadow-xl py-1 scrollbar-thin"
        >
          <!-- Rates options -->
          <template v-if="selectedMetric === 'rates'">
            <button
              v-for="c in CURRENCIES"
              :key="c.value"
              type="button"
              @click="selectItem(c.value)"
              class="w-full text-left px-3 py-2 text-xs text-text-secondary hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
            >
              <span class="flex items-center gap-2">
                <img :src="'https://flagcdn.com/w40/' + c.flagCode + '.png'" class="w-5 h-3.5 object-cover rounded-xs border border-white/10" />
                <span class="font-bold text-text-primary">{{ c.value }}</span>
                <span class="text-text-tertiary">— {{ c.name }}</span>
              </span>
              <span v-if="selectedItem === c.value" class="text-blue-450 font-bold">✓</span>
            </button>
          </template>

          <!-- Stocks options -->
          <template v-else-if="selectedMetric === 'market'">
            <button
              v-for="m in MARKET_STOCKS"
              :key="m.value"
              type="button"
              @click="selectItem(m.value)"
              class="w-full text-left px-3 py-2.5 text-xs text-text-secondary hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
            >
              <span class="flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 3 2" class="w-5 h-3.5 rounded-xs border border-white/10 shrink-0" aria-hidden="true">
                  <rect width="3" height="1" fill="#FF0000"/>
                  <rect y="1" width="3" height="1" fill="#FFFFFF"/>
                </svg>
                <span class="font-bold text-text-primary">{{ m.name }}</span>
                <span class="text-text-tertiary font-mono text-[9px] uppercase">({{ m.value }})</span>
              </span>
              <span v-if="selectedItem === m.value" class="text-blue-450 font-bold">✓</span>
            </button>
          </template>

          <!-- Commodities options -->
          <template v-else-if="selectedMetric === 'commodities'">
            <button
              v-for="c in COMMODITIES"
              :key="c.value"
              type="button"
              @click="selectItem(c.value)"
              class="w-full text-left px-3 py-2 text-xs text-text-secondary hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
            >
              <span class="flex items-center gap-2">
                <span class="text-sm">{{ c.icon }}</span>
                <span class="font-bold text-text-primary">{{ c.name }}</span>
                <span class="text-text-tertiary font-mono text-[9px] uppercase">({{ c.value }})</span>
              </span>
              <span v-if="selectedItem === c.value" class="text-blue-450 font-bold">✓</span>
            </button>
          </template>

          <!-- Cities options -->
          <template v-else-if="selectedMetric === 'weather'">
            <button
              v-for="c in CITIES"
              :key="c.value"
              type="button"
              @click="selectItem(c.value)"
              class="w-full text-left px-3 py-2 text-xs text-text-secondary hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
            >
              <span class="flex items-center gap-2">
                <span class="text-sm">{{ c.icon }}</span>
                <span class="font-bold text-text-primary">{{ c.name }}</span>
              </span>
              <span v-if="selectedItem === c.value" class="text-blue-450 font-bold">✓</span>
            </button>
          </template>
        </div>
      </div>

      <!-- 3. Weather Metric Parameter (Cuaca only) -->
      <div v-if="selectedMetric === 'weather'" class="col-span-1 md:col-span-2 relative">
        <label class="block text-[10px] font-bold text-text-tertiary uppercase tracking-wider mb-1.5">
          Parameter Cuaca
        </label>
        <button
          type="button"
          @click.stop="toggleDropdown('weatherMetric')"
          class="w-full flex items-center justify-between text-xs rounded-lg border border-white/5 bg-bg-card px-3 py-2 text-text-primary focus:border-blue-500/40 focus:outline-none transition-all cursor-pointer hover:border-white/10 hover:bg-bg-subcard/50 active:scale-[0.99]"
        >
          <span>{{ WEATHER_METRICS.find(w => w.value === selectedWeatherMetric)?.label }}</span>
          <svg xmlns="http://www.w3.org/2500/svg" class="h-3.5 w-3.5 text-text-secondary transition-transform duration-200" :class="{ 'rotate-180': activeDropdown === 'weatherMetric' }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </button>
        <div
          v-if="activeDropdown === 'weatherMetric'"
          class="absolute left-0 right-0 mt-1.5 z-50 max-h-60 overflow-y-auto rounded-lg border border-white/10 bg-bg-card/95 backdrop-blur-xl shadow-xl py-1 scrollbar-thin"
        >
          <button
            v-for="w in WEATHER_METRICS"
            :key="w.value"
            type="button"
            @click="selectWeatherMetric(w.value)"
            class="w-full text-left px-3 py-2 text-xs text-text-secondary hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
          >
            <span>{{ w.label }} ({{ w.unit }})</span>
            <span v-if="selectedWeatherMetric === w.value" class="text-blue-450 font-bold">✓</span>
          </button>
        </div>
      </div>

      <!-- 4. Chart Type Selection -->
      <div class="col-span-1 md:col-span-2 relative">
        <label class="block text-[10px] font-bold text-text-tertiary uppercase tracking-wider mb-1.5">
          Tipe Grafik
        </label>
        <button
          type="button"
          @click.stop="toggleDropdown('chartType')"
          class="w-full flex items-center justify-between text-xs rounded-lg border border-white/5 bg-bg-card px-3 py-2 text-text-primary focus:border-blue-500/40 focus:outline-none transition-all cursor-pointer hover:border-white/10 hover:bg-bg-subcard/50 active:scale-[0.99]"
        >
          <span>{{ selectedChartType === 'line' ? 'Line Chart' : selectedChartType === 'bar' ? 'Bar Chart' : 'Area Chart' }}</span>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-text-secondary transition-transform duration-200" :class="{ 'rotate-180': activeDropdown === 'chartType' }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </button>
        <div
          v-if="activeDropdown === 'chartType'"
          class="absolute left-0 right-0 mt-1.5 z-50 max-h-60 overflow-y-auto rounded-lg border border-white/10 bg-bg-card/95 backdrop-blur-xl shadow-xl py-1 scrollbar-thin"
        >
          <button
            type="button"
            @click="selectChartType('line')"
            class="w-full text-left px-3 py-2 text-xs text-text-secondary hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
          >
            <span>Line Chart</span>
            <span v-if="selectedChartType === 'line'" class="text-blue-450 font-bold">✓</span>
          </button>
          <button
            type="button"
            @click="selectChartType('bar')"
            class="w-full text-left px-3 py-2 text-xs text-text-secondary hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
          >
            <span>Bar Chart</span>
            <span v-if="selectedChartType === 'bar'" class="text-blue-450 font-bold">✓</span>
          </button>
          <button
            type="button"
            @click="selectChartType('area')"
            class="w-full text-left px-3 py-2 text-xs text-text-secondary hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
          >
            <span>Area Chart</span>
            <span v-if="selectedChartType === 'area'" class="text-blue-450 font-bold">✓</span>
          </button>
        </div>
      </div>

      <!-- 5. Time Range Selection -->
      <div class="col-span-1 relative" :class="selectedMetric === 'weather' ? 'md:col-span-1.5' : 'md:col-span-2'">
        <label class="block text-[10px] font-bold text-text-tertiary uppercase tracking-wider mb-1.5">
          Waktu
        </label>
        <button
          type="button"
          @click.stop="toggleDropdown('range')"
          class="w-full flex items-center justify-between text-xs rounded-lg border border-white/5 bg-bg-card px-3 py-2 text-text-primary focus:border-blue-500/40 focus:outline-none transition-all cursor-pointer hover:border-white/10 hover:bg-bg-subcard/50 active:scale-[0.99]"
        >
          <span>{{ selectedRange }} Hari</span>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-text-secondary transition-transform duration-200" :class="{ 'rotate-180': activeDropdown === 'range' }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </button>
        <div
          v-if="activeDropdown === 'range'"
          class="absolute left-0 right-0 mt-1.5 z-50 max-h-60 overflow-y-auto rounded-lg border border-white/10 bg-bg-card/95 backdrop-blur-xl shadow-xl py-1 scrollbar-thin"
        >
          <button
            type="button"
            @click="selectRange(3)"
            class="w-full text-left px-3 py-2 text-xs text-text-secondary hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
          >
            <span>3 Hari</span>
            <span v-if="Number(selectedRange) === 3" class="text-blue-450 font-bold">✓</span>
          </button>
          <button
            type="button"
            @click="selectRange(7)"
            class="w-full text-left px-3 py-2 text-xs text-text-secondary hover:text-white hover:bg-white/4 transition-colors flex items-center justify-between cursor-pointer"
          >
            <span>7 Hari</span>
            <span v-if="Number(selectedRange) === 7" class="text-blue-450 font-bold">✓</span>
          </button>
        </div>
      </div>

      <!-- 6. Apply Button -->
      <div class="col-span-1 md:col-span-2 flex items-end">
        <button
          @click="applySettings()"
          class="w-full h-[34px] inline-flex items-center justify-center gap-1.5 px-4 py-2 text-xs font-semibold rounded-lg bg-blue-500 hover:bg-blue-400 text-white shadow-md shadow-blue-500/10 hover:shadow-blue-500/20 active:scale-[0.98] transition-all cursor-pointer"
        >
          <span>⚡</span> Terapkan
        </button>
      </div>
    </div>

    <!-- Chart Panel -->
    <div class="mt-6">
      <div class="flex items-center justify-between mb-4 border-b border-white/5 pb-3">
        <div class="flex items-center gap-2">
          <!-- Currency Flag -->
          <img 
            v-if="appliedConfig.metric === 'rates'"
            :src="'https://flagcdn.com/w40/' + (appliedItemDetails.flagCode || 'un') + '.png'" 
            class="w-5 h-3.5 object-cover rounded-xs border border-white/10 shrink-0" 
          />
          <!-- Stock Indonesian Flag SVG -->
          <svg v-else-if="appliedConfig.metric === 'market'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 3 2" class="w-5 h-3.5 rounded-xs border border-white/10 shrink-0" aria-hidden="true">
            <rect width="3" height="1" fill="#FF0000"/>
            <rect y="1" width="3" height="1" fill="#FFFFFF"/>
          </svg>
          <!-- Commodity Icon -->
          <span v-else-if="appliedConfig.metric === 'commodities'" class="text-sm">
            {{ appliedItemDetails.icon }}
          </span>
          <!-- Weather City Icon -->
          <span v-else-if="appliedConfig.metric === 'weather'" class="text-sm">
            {{ appliedItemDetails.icon }}
          </span>

          <h4 class="text-xs font-extrabold text-text-secondary uppercase tracking-wide">
            {{ activeTitle }}
          </h4>
        </div>
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
      <div v-if="!isLoading && chartData.dataPoints.length > 0" class="mt-6 overflow-x-auto border border-white/5 rounded-xl bg-bg-subcard/50 scrollbar-thin">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="bg-bg-subcard border-b border-white/5 text-text-tertiary font-bold tracking-wider uppercase">
              <th class="py-2 px-4">Tanggal</th>
              <th class="py-2 px-4 text-right">Nilai ({{ chartData.unit }})</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-white/5">
            <tr v-for="(val, idx) in chartData.dataPoints" :key="'preview-'+idx" class="hover:bg-white/2 text-text-secondary">
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
  border-color: rgba(59, 130, 246, 0.4);
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.1);
}
</style>
