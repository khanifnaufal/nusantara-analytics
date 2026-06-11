<script setup lang="ts">
import { computed } from 'vue'
import BaseChart from './BaseChart.vue'

interface SeriesData {
  name: string
  data: number[]
}

interface Props {
  data: SeriesData[]
  xAxis: string[]
  title?: string
  height?: string
  loading?: boolean
  unit?: string
  accentColor?: string
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  height: '300px',
  loading: false,
  unit: '',
  accentColor: ''
})

const option = computed(() => {
  const colors = props.accentColor
    ? [props.accentColor, '#60A5FA', '#93C5FD', '#2563EB', '#1D4ED8']
    : ['#3B82F6', '#60A5FA', '#93C5FD', '#2563EB', '#1D4ED8']

  return {
    backgroundColor: 'transparent',
    title: props.title
      ? {
          text: props.title,
          textStyle: {
            fontFamily: 'Inter, sans-serif',
            fontSize: 14,
            fontWeight: 600,
            color: '#F0F6FC'
          },
          top: 0,
          left: 0
        }
      : undefined,
    color: colors,
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
          const val = item.value !== undefined && item.value !== null
            ? new Intl.NumberFormat().format(item.value)
            : '-'
          result += `<div style="display: flex; align-items: center; justify-content: space-between; gap: 16px; margin-top: 4px;">
            <div style="display: flex; align-items: center; gap: 6px;">
              ${item.marker}
              <span style="color: var(--tooltip-text-color, #8B949E); font-weight: 500;">${item.seriesName}</span>
            </div>
            <span style="font-weight: 600; color: var(--tooltip-val-color, #ffffff);">${val}${props.unit ? ' ' + props.unit : ''}</span>
          </div>`
        })
        result += `</div>`
        return result
      }
    },
    legend: {
      show: props.data.length > 1,
      bottom: 0,
      left: 'center',
      icon: 'circle',
      textStyle: {
        fontFamily: 'Inter, sans-serif',
        fontSize: 10,
        color: '#8B949E'
      }
    },
    grid: {
      top: props.title ? 50 : 15,
      left: '2%',
      right: '2%',
      bottom: props.data.length > 1 ? 35 : 15,
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: true,
      data: props.xAxis,
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
          return new Intl.NumberFormat(undefined, { notation: 'compact' }).format(value) + (props.unit ? ' ' + props.unit : '')
        }
      },
      splitLine: {
        lineStyle: {
          color: 'rgba(255, 255, 255, 0.04)',
          type: 'dashed'
        }
      }
    },
    series: props.data.map(s => ({
      name: s.name,
      type: 'bar',
      data: s.data,
      barMaxWidth: 24,
      itemStyle: {
        borderRadius: [4, 4, 0, 0]
      },
      emphasis: {
        focus: 'series'
      }
    }))
  }
})
</script>

<template>
  <BaseChart
    :option="option"
    :height="height"
    :loading="loading"
  />
</template>

