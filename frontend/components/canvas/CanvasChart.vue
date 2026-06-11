<script setup lang="ts">
import { computed } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import {
  LineChart,
  BarChart,
  PieChart,
  ScatterChart,
  GaugeChart,
  TreemapChart,
  RadarChart,
  CandlestickChart,
  HeatmapChart
} from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent,
  VisualMapComponent,
  RadarComponent,
  CalendarComponent
} from 'echarts/components'
import VChart from 'vue-echarts'
import { useDark } from '@vueuse/core'
import type { CanvasWidget } from '~/stores/canvas'

import { useRatesStore } from '~/stores/rates'
import { useWeatherStore } from '~/stores/weather'
import { useCommoditiesStore } from '~/stores/commodities'
import { useMarketStore } from '~/stores/market'
import { useQuakesStore } from '~/stores/quakes'

// Register ECharts modules
use([
  CanvasRenderer,
  LineChart,
  BarChart,
  PieChart,
  ScatterChart,
  GaugeChart,
  TreemapChart,
  RadarChart,
  CandlestickChart,
  HeatmapChart,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent,
  VisualMapComponent,
  RadarComponent,
  CalendarComponent
])

interface Props {
  widget: CanvasWidget
}

const props = defineProps<Props>()

const isDark = useDark()
const theme = computed(() => (isDark.value ? 'dark' : 'light'))

// Access Pinia Stores
const ratesStore = useRatesStore()
const weatherStore = useWeatherStore()
const commoditiesStore = useCommoditiesStore()
const marketStore = useMarketStore()
const quakesStore = useQuakesStore()

// Check loading state dynamically
const loading = computed(() => {
  switch (props.widget.dataSource) {
    case 'rates': return ratesStore.loading
    case 'weather': return weatherStore.loading
    case 'commodities': return commoditiesStore.loading
    case 'market': return marketStore.loading
    case 'quakes': return quakesStore.loading
    default: return false
  }
})

// Deterministic hashing helper for data generation
const getHash = (str: string): number => {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = (hash << 5) - hash + str.charCodeAt(i)
    hash |= 0
  }
  return Math.abs(hash)
}

// Deterministic Currency History
const getRatesHistory = (currency: string, currentValInIdr: number, rangeDays = 7) => {
  const history = []
  const today = new Date()
  let val = currentValInIdr
  for (let i = 0; i < rangeDays; i++) {
    const d = new Date()
    d.setDate(today.getDate() - i)
    const dateStr = d.toISOString().split('T')[0] || ''
    history.push({ date: dateStr, value: val })
    const hash = getHash(currency + dateStr)
    const changePercent = (hash % 100) / 180 - 0.27 // vol ~0.3%
    val = val / (1 + changePercent / 100)
  }
  return history.reverse()
}

// Deterministic Weather History
const getWeatherHistory = (city: string, metric: string, currentVal: number, rangeDays = 7) => {
  const history = []
  const today = new Date()
  let val = currentVal
  for (let i = 0; i < rangeDays; i++) {
    const d = new Date()
    d.setDate(today.getDate() - i)
    const dateStr = d.toISOString().split('T')[0] || ''
    history.push({ date: dateStr, value: Number(val.toFixed(1)) })
    const hash = getHash(city + metric + dateStr)
    if (metric === 'temperature') {
      const change = (hash % 200) / 100 - 1.0 // -1 to +1 deg C
      val = val - change
    } else if (metric === 'humidity') {
      const change = (hash % 10) - 5
      val = Math.max(30, Math.min(100, val - change))
    } else {
      const change = (hash % 40) / 10 - 2
      val = Math.max(0, val - change)
    }
  }
  return history.reverse()
}

// Deterministic Stock/Commodity History
const getQuoteHistoryData = (symbol: string, currentPrice: number, historyList: any[] | undefined, rangeDays = 7) => {
  if (historyList && historyList.length > 0) {
    const sorted = [...historyList].sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime())
    const startIdx = Math.max(0, sorted.length - rangeDays)
    return sorted.slice(startIdx).map(h => ({
      date: new Date(h.timestamp).toISOString().split('T')[0] || '',
      value: h.close
    }))
  }
  // Generate fallback
  const history = []
  const today = new Date()
  let val = currentPrice
  for (let i = 0; i < rangeDays; i++) {
    const d = new Date()
    d.setDate(today.getDate() - i)
    const dateStr = d.toISOString().split('T')[0] || ''
    history.push({ date: dateStr, value: Number(val.toFixed(2)) })
    const hash = getHash(symbol + dateStr)
    const changePercent = (hash % 160) / 100 - 0.8 // -0.8% to +0.8%
    val = val / (1 + changePercent / 100)
  }
  return history.reverse()
}

// Generate Mock OHLC Candlestick data deterministically
const getCandlestickData = (history: { date: string; value: number }[], seedSymbol: string) => {
  return history.map((item, idx) => {
    const close = item.value
    const hash = getHash(seedSymbol + item.date + idx)
    
    // Deterministic open, high, low
    const open = close * (1 + ((hash % 20) - 10) / 1000)
    const maxVal = Math.max(open, close)
    const minVal = Math.min(open, close)
    const high = maxVal * (1 + (hash % 10) / 1000)
    const low = minVal * (1 - (hash % 10) / 1000)
    
    return [
      Number(open.toFixed(2)),
      Number(close.toFixed(2)),
      Number(low.toFixed(2)),
      Number(high.toFixed(2))
    ] // format: [open, close, lowest, highest]
  })
}

// ECharts configurations computed property
const option = computed(() => {
  const type = props.widget.chartType
  const source = props.widget.dataSource
  const metric = props.widget.metric
  const config = props.widget.config || {}
  const range = config.range || 7

  // Custom accent color preset resolution
  const colorPresetMap: Record<string, string> = {
    blue: '#3B82F6',
    green: '#10B981',
    red: '#EF4444',
    amber: '#F59E0B',
    violet: '#8B5CF6',
    cyan: '#06B6D4'
  }
  const accentColorPreset = config.accentColor || 'blue'
  const primaryBlue = colorPresetMap[accentColorPreset] || '#3B82F6'

  // Visibility toggles
  const showLabels = config.showLabels !== false
  const showLegend = config.showLegend !== false

  const isDarkTheme = theme.value === 'dark'
  const labelColor = isDarkTheme ? '#A1A1AA' : '#52525B'
  const splitLineColor = isDarkTheme ? 'rgba(255, 255, 255, 0.04)' : 'rgba(0, 0, 0, 0.05)'

  // Base options template
  let chartOption: Record<string, any> = {
    backgroundColor: 'transparent',
    textStyle: { fontFamily: 'Inter, sans-serif' },
    color: [
      primaryBlue,
      '#10B981',
      '#F59E0B',
      '#8B5CF6',
      '#06B6D4',
      '#EF4444',
      '#EC4899'
    ],
    tooltip: {
      trigger: 'axis',
      backgroundColor: isDarkTheme ? '#111111' : '#FFFFFF',
      borderColor: isDarkTheme ? 'rgba(255, 255, 255, 0.08)' : 'rgba(0, 0, 0, 0.08)',
      textStyle: { color: isDarkTheme ? '#FFFFFF' : '#111111', fontSize: 11 }
    },
    legend: {
      show: showLegend,
      textStyle: { color: labelColor, fontSize: 10 },
      top: 5
    },
    grid: { top: showLegend ? 45 : 30, right: 15, bottom: 35, left: 45, containLabel: true }
  }

  // --- DATA EXTRACTION AND TRANSFORMATION ---
  let xAxisData: any[] = []
  let seriesData: any[] = []
  let legendData: any[] = []
  let chartTitle = ''
  let unit = ''

  try {
    if (source === 'rates') {
      const dataObj = ratesStore.data
      const rawVal = dataObj?.rates[metric]
      const currentInIdr = rawVal ? 1 / rawVal : 14000
      unit = 'IDR'

      // Check if displaying current list or single item history
      if (config.field === 'nilai') {
        // Current values of main rates comparison
        if (type === 'pie' || type === 'donut' || type === 'treemap') {
          xAxisData = ['USD', 'EUR', 'SGD', 'JPY', 'MYR', 'CNY', 'AUD']
          seriesData = xAxisData.map(c => {
            const r = dataObj?.rates[c]
            return { name: c, value: r ? Math.round(1 / r) : 10000 }
          })
        } else if (type === 'radar') {
          legendData = ['Kurs Valas']
          seriesData = [{
            value: ['USD', 'EUR', 'SGD', 'AUD'].map(c => {
              const r = dataObj?.rates[c]
              return r ? Math.round(1 / r) : 10000
            }),
            name: 'Kurs Valas (IDR)'
          }]
          chartOption.radar = {
            indicator: [
              { name: 'USD', max: 18000 },
              { name: 'EUR', max: 20000 },
              { name: 'SGD', max: 13000 },
              { name: 'AUD', max: 12000 }
            ],
            axisName: { show: showLabels, color: showLabels ? labelColor : 'transparent', fontSize: 10 }
          }
        } else if (type === 'gauge') {
          seriesData = [{ name: metric, value: Math.round(currentInIdr) }]
          chartOption.series = [{
            type: 'gauge',
            min: Math.round(currentInIdr * 0.8),
            max: Math.round(currentInIdr * 1.2),
            detail: { formatter: '{value} IDR', fontSize: 14, color: isDarkTheme ? '#FFFFFF' : '#111111' },
            data: seriesData
          }]
          return chartOption
        } else {
          // Fallback to history trend
          const hist = getRatesHistory(metric, currentInIdr, range)
          xAxisData = hist.map(h => h.date)
          seriesData = hist.map(h => Number(h.value.toFixed(2)))
        }
      } else {
        // history field
        const hist = getRatesHistory(metric, currentInIdr, range)
        xAxisData = hist.map(h => h.date)
        seriesData = hist.map(h => Number(h.value.toFixed(2)))

        if (type === 'pie' || type === 'donut') {
          seriesData = hist.map(h => ({ name: h.date, value: Math.round(h.value) }))
        } else if (type === 'gauge') {
          chartOption.series = [{
            type: 'gauge',
            min: Math.round(currentInIdr * 0.8),
            max: Math.round(currentInIdr * 1.2),
            detail: { formatter: '{value} IDR', fontSize: 14, color: isDarkTheme ? '#FFFFFF' : '#111111' },
            data: [{ name: metric, value: Math.round(currentInIdr) }]
          }]
          return chartOption
        }
      }
    } else if (source === 'weather') {
      const cities = weatherStore.data?.cities || []
      const targetCity = cities.find(c => c.city === metric) || {
        city: metric, temperature: 30, humidity: 80, windSpeed: 10
      }
      const valField = config.field || 'temperature'
      const val = targetCity[valField as 'temperature' | 'humidity' | 'windSpeed'] || 0
      unit = valField === 'temperature' ? '°C' : valField === 'humidity' ? '%' : 'km/h'

      if (type === 'pie' || type === 'donut' || type === 'treemap') {
        // Compare parameter across all cities
        xAxisData = cities.map(c => c.city)
        seriesData = cities.map(c => ({
          name: c.city,
          value: c[valField as 'temperature' | 'humidity' | 'windSpeed'] || 0
        }))
      } else if (type === 'radar') {
        legendData = [metric]
        seriesData = [{
          value: [targetCity.temperature, targetCity.humidity, targetCity.windSpeed],
          name: metric
        }]
        chartOption.radar = {
          indicator: [
            { name: 'Suhu (°C)', max: 45 },
            { name: 'Kelembapan (%)', max: 100 },
            { name: 'Kec. Angin (km/h)', max: 50 }
          ],
          axisName: { show: showLabels, color: showLabels ? labelColor : 'transparent', fontSize: 10 }
        }
      } else if (type === 'gauge') {
        const maxLimit = valField === 'temperature' ? 50 : valField === 'humidity' ? 100 : 80
        chartOption.series = [{
          type: 'gauge',
          min: 0,
          max: maxLimit,
          detail: { formatter: `{value}${unit}`, fontSize: 14, color: isDarkTheme ? '#FFFFFF' : '#111111' },
          data: [{ name: valField, value: val }]
        }]
        return chartOption
      } else if (type === 'heatmap') {
        // Grid of cities vs days of the week (last 5 days)
        const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri']
        chartOption.xAxis = { type: 'category', data: days, axisLabel: { color: labelColor } }
        chartOption.yAxis = { type: 'category', data: cities.map(c => c.city), axisLabel: { color: labelColor } }
        
        const hmData: any[] = []
        cities.forEach((c, cIdx) => {
          days.forEach((d, dIdx) => {
            const hVal = c[valField as 'temperature' | 'humidity' | 'windSpeed'] || 25
            const hash = getHash(c.city + d + valField)
            const simulated = valField === 'temperature' 
              ? hVal + (hash % 6 - 3) 
              : valField === 'humidity' 
                ? Math.max(30, Math.min(100, hVal + (hash % 16 - 8)))
                : Math.max(0, hVal + (hash % 8 - 4))
            hmData.push([dIdx, cIdx, Number(simulated.toFixed(1))])
          })
        })
        
        chartOption.visualMap = {
          min: valField === 'temperature' ? 20 : valField === 'humidity' ? 40 : 0,
          max: valField === 'temperature' ? 38 : valField === 'humidity' ? 95 : 30,
          calculable: true,
          orient: 'horizontal',
          left: 'center',
          bottom: 0,
          textStyle: { color: labelColor, fontSize: 9 }
        }
        
        chartOption.series = [{
          name: valField,
          type: 'heatmap',
          data: hmData,
          label: { show: true, fontSize: 9, color: isDarkTheme ? '#FFF' : '#333' }
        }]
        return chartOption
      } else {
        // trend lines
        const hist = getWeatherHistory(metric, valField, val, range)
        xAxisData = hist.map(h => h.date)
        seriesData = hist.map(h => h.value)
      }
    } else if (source === 'commodities' || source === 'market') {
      const itemsList = source === 'commodities' 
        ? commoditiesStore.data?.commodities 
        : marketStore.data?.stocks
      
      const quote = itemsList?.find(q => q.symbol === metric) || {
        symbol: metric, price: 1000, history: [], currency: 'USD', name: metric
      }
      unit = source === 'market' && metric === '^JKSE' ? 'pt' : quote.currency || 'USD'

      if (config.field === 'nilai' && (type === 'pie' || type === 'donut' || type === 'treemap')) {
        // Compare list of items
        const rawItems = itemsList || []
        xAxisData = rawItems.map(i => i.name)
        seriesData = rawItems.map(i => ({ name: i.name, value: i.price }))
      } else if (type === 'gauge') {
        const price = quote.price
        chartOption.series = [{
          type: 'gauge',
          min: Math.round(price * 0.5),
          max: Math.round(price * 1.5),
          detail: { formatter: `${unit} {value}`, fontSize: 13, color: isDarkTheme ? '#FFFFFF' : '#111111' },
          data: [{ name: quote.name, value: Number(price.toFixed(2)) }]
        }]
        return chartOption
      } else {
        // history trend
        const hist = getQuoteHistoryData(metric, quote.price, quote.history, range)
        xAxisData = hist.map(h => h.date)
        seriesData = hist.map(h => h.value)

        if (type === 'candlestick') {
          // Map to open, close, lowest, highest
          const candleData = getCandlestickData(hist, metric)
          chartOption.xAxis = {
            type: 'category',
            data: xAxisData,
            axisLabel: { color: labelColor },
            axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.08)' } }
          }
          chartOption.yAxis = {
            type: 'value',
            scale: true,
            axisLabel: { color: labelColor },
            splitLine: { lineStyle: { color: splitLineColor } }
          }
          chartOption.series = [{
            type: 'candlestick',
            name: quote.name,
            data: candleData,
            itemStyle: {
              color: '#EF4444',
              color0: '#10B981',
              borderColor: '#EF4444',
              borderColor0: '#10B981'
            }
          }]
          return chartOption
        } else if (type === 'pie' || type === 'donut') {
          seriesData = hist.map(h => ({ name: h.date, value: Math.round(h.value) }))
        }
      }
    } else if (source === 'quakes') {
      const quakes = quakesStore.data?.quakes || []
      
      if (metric === 'all' || !metric) {
        const param = config.field || 'timeline'
        
        if (param === 'magnitude distribution') {
          // Count magnitude intervals
          const levels = {
            'Rendah (M < 3)': 0,
            'Sedang (3 <= M < 5)': 0,
            'Tinggi (M >= 5)': 0
          }
          quakes.forEach(q => {
            if (q.magnitude < 3) levels['Rendah (M < 3)']++
            else if (q.magnitude < 5) levels['Sedang (3 <= M < 5)']++
            else levels['Tinggi (M >= 5)']++
          })

          seriesData = Object.entries(levels).map(([name, value]) => ({ name, value }))
          xAxisData = Object.keys(levels)
          
        } else if (param === 'depth distribution') {
          // Scatter plot: Depth (x) vs Magnitude (y)
          chartOption.tooltip = {
            trigger: 'item',
            formatter: (params: any) => {
              const d = params.value
              return `<div style="padding: 4px 6px; font-size:11px;">
                <strong>Lokasi:</strong> ${d[2]}<br/>
                <strong>Kedalaman:</strong> ${d[0]} km<br/>
                <strong>Magnitudo:</strong> M ${d[1]}
              </div>`
            }
          }
          chartOption.xAxis = {
            type: 'value',
            name: showLabels ? 'Kedalaman (km)' : '',
            nameGap: 20,
            nameLocation: 'middle',
            nameTextStyle: { color: showLabels ? labelColor : 'transparent', fontSize: 10 },
            axisLabel: { show: showLabels, color: labelColor },
            splitLine: { lineStyle: { color: splitLineColor } }
          }
          chartOption.yAxis = {
            type: 'value',
            name: showLabels ? 'Magnitudo' : '',
            axisLabel: { show: showLabels, color: labelColor },
            splitLine: { lineStyle: { color: splitLineColor } }
          }
          chartOption.series = [{
            type: 'scatter',
            symbolSize: (data: any) => data[1] * 2.5,
            data: quakes.map(q => [q.depth, q.magnitude, q.place]),
            itemStyle: {
              color: primaryBlue,
              opacity: 0.75,
              borderColor: 'rgba(255,255,255,0.1)',
              borderWidth: 1
            }
          }]
          return chartOption
        } else {
          // timeline: plot magnitude of quakes over time
          const sorted = [...quakes].sort((a, b) => new Date(a.time).getTime() - new Date(b.time).getTime())
          xAxisData = sorted.map(q => {
            const d = new Date(q.time)
            return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short' }) + ' ' + d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
          })
          seriesData = sorted.map(q => q.magnitude)
          unit = 'M'
        }
      }
    }
  } catch (err) {
    console.error('Error parsing widget data:', err)
  }

  // --- CHART RENDERING OPTION BUILDER ---
  if (type === 'line' || type === 'area' || type === 'bar') {
    chartOption.xAxis = {
      type: 'category',
      data: xAxisData,
      axisLabel: { show: showLabels, color: labelColor, fontSize: 9, rotate: range > 10 ? 30 : 0 },
      axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.08)' } },
      axisTick: { show: false }
    }
    chartOption.yAxis = {
      type: 'value',
      scale: true,
      axisLabel: {
        show: showLabels,
        color: labelColor,
        fontSize: 9,
        formatter: (v: number) => {
          if (unit === 'IDR' || unit === 'Rp') {
            return new Intl.NumberFormat('id-ID', { notation: 'compact', style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(v)
          }
          return v + (unit ? ' ' + unit : '')
        }
      },
      splitLine: { lineStyle: { color: splitLineColor, type: 'dashed' } }
    }

    const isBar = type === 'bar'
    chartOption.series = [{
      name: props.widget.title,
      type: isBar ? 'bar' : 'line',
      data: seriesData,
      smooth: !isBar,
      showSymbol: false,
      itemStyle: isBar ? { borderRadius: [3, 3, 0, 0], color: primaryBlue } : { color: primaryBlue },
      lineStyle: !isBar ? { width: type === 'area' ? 2 : 2.5, color: primaryBlue } : undefined,
      areaStyle: type === 'area' ? {
        opacity: 0.12,
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: primaryBlue },
            { offset: 1, color: 'transparent' }
          ]
        }
      } : undefined
    }]
  } else if (type === 'pie' || type === 'donut') {
    chartOption.tooltip = {
      trigger: 'item',
      formatter: '{b}: <strong>{c}</strong> ({d}%)'
    }
    chartOption.series = [{
      name: props.widget.title,
      type: 'pie',
      radius: type === 'donut' ? ['40%', '70%'] : '70%',
      avoidLabelOverlap: true,
      itemStyle: {
        borderRadius: 4,
        borderColor: isDarkTheme ? '#111111' : '#FFFFFF',
        borderWidth: 1
      },
      label: {
        show: showLabels,
        fontSize: 9,
        color: labelColor,
        formatter: '{b}'
      },
      labelLine: {
        lineStyle: { color: 'rgba(255, 255, 255, 0.1)' },
        smooth: 0.2,
        length: 8,
        length2: 10
      },
      data: seriesData
    }]
  } else if (type === 'treemap') {
    chartOption.tooltip = { trigger: 'item' }
    chartOption.series = [{
      type: 'treemap',
      data: seriesData,
      breadcrumb: { show: false },
      label: { show: showLabels, fontSize: 10, formatter: '{b}\n{c}' },
      itemStyle: {
        borderColor: isDarkTheme ? '#111111' : '#FFFFFF',
        borderWidth: 1
      }
    }]
  } else if (type === 'radar') {
    chartOption.tooltip = { trigger: 'item' }
    chartOption.legend = {
      show: showLegend,
      data: legendData,
      bottom: 0,
      textStyle: { color: labelColor, fontSize: 9 }
    }
    chartOption.series = [{
      type: 'radar',
      data: seriesData,
      itemStyle: { color: primaryBlue },
      areaStyle: { opacity: 0.1 }
    }]
  }

  return chartOption
})
</script>

<template>
  <div class="relative w-full h-full">
    <!-- Skeleton loader -->
    <div
      v-if="loading"
      class="w-full h-full flex flex-col items-center justify-center rounded-2xl bg-zinc-900/30 animate-pulse p-6"
    >
      <div class="flex flex-col items-center gap-3">
        <div class="relative w-8 h-8">
          <div class="absolute inset-0 border-2 border-white/5 rounded-full"></div>
          <div class="absolute inset-0 border-2 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
        </div>
        <p class="text-[10px] text-text-tertiary font-medium tracking-wide">
          Memuat data canvas...
        </p>
      </div>
    </div>

    <!-- Empty State -->
    <div
      v-else-if="!option || (!option.series && !option.radar)"
      class="w-full h-full flex flex-col items-center justify-center p-6 text-center text-text-tertiary"
    >
      <span class="text-xl mb-1.5">📊</span>
      <p class="text-[10px]">Data tidak tersedia untuk konfigurasi ini.</p>
    </div>

    <!-- Chart -->
    <VChart
      v-else
      class="w-full h-full"
      :option="option"
      :theme="theme"
      :autoresize="true"
    />
  </div>
</template>
