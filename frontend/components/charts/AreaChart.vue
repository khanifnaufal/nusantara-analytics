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
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  height: '300px',
  loading: false,
  unit: ''
})

const COLORS = ['#10b981', '#6366f1', '#f59e0b', '#ec4899', '#3b82f6', '#8b5cf6']

const option = computed(() => {
  return {
    title: props.title
      ? {
          text: props.title,
          textStyle: {
            fontFamily: 'Inter, sans-serif',
            fontSize: 14,
            fontWeight: 600
          },
          top: 0,
          left: 0
        }
      : undefined,
    color: COLORS,
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
        result += `<div style="font-weight: 600; margin-bottom: 6px; color: var(--tooltip-title-color, #64748b);">${title}</div>`
        params.forEach((item: any) => {
          const val = item.value !== undefined && item.value !== null
            ? new Intl.NumberFormat().format(item.value)
            : '-'
          result += `<div style="display: flex; align-items: center; justify-content: space-between; gap: 16px; margin-top: 4px;">
            <div style="display: flex; align-items: center; gap: 6px;">
              ${item.marker}
              <span style="color: var(--tooltip-text-color, #475569); font-weight: 500;">${item.seriesName}</span>
            </div>
            <span style="font-weight: 600; color: var(--tooltip-val-color, #0f172a);">${val}${props.unit ? ' ' + props.unit : ''}</span>
          </div>`
        })
        result += `</div>`
        return result
      }
    },
    legend: {
      show: true,
      bottom: 0,
      left: 'center',
      icon: 'circle',
      textStyle: {
        fontFamily: 'Inter, sans-serif',
        fontSize: 11
      }
    },
    grid: {
      top: props.title ? 50 : 20,
      left: '2%',
      right: '2%',
      bottom: 35,
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: props.xAxis,
      axisLine: {
        show: true,
        lineStyle: {
          color: 'rgba(156, 163, 175, 0.2)'
        }
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        fontFamily: 'Inter, sans-serif',
        fontSize: 11
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
        fontSize: 11,
        formatter: (value: number) => {
          return new Intl.NumberFormat(undefined, { notation: 'compact' }).format(value) + (props.unit ? ' ' + props.unit : '')
        }
      },
      splitLine: {
        lineStyle: {
          color: 'rgba(156, 163, 175, 0.1)',
          type: 'dashed'
        }
      }
    },
    series: props.data.map((s, idx) => {
      const color = COLORS[idx % COLORS.length]
      return {
        name: s.name,
        type: 'line',
        data: s.data,
        smooth: true,
        showSymbol: false,
        symbolSize: 6,
        color: color,
        areaStyle: {
          opacity: 0.12,
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: color },
              { offset: 1, color: 'transparent' }
            ],
            global: false
          }
        },
        emphasis: {
          focus: 'series',
          itemStyle: {
            borderWidth: 2
          }
        },
        lineStyle: {
          width: 2.5
        }
      }
    })
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
