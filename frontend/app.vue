<script setup lang="ts">
import { ref } from 'vue'
import LineChart from '~/components/charts/LineChart.vue'
import BarChart from '~/components/charts/BarChart.vue'
import AreaChart from '~/components/charts/AreaChart.vue'
import { useDark, useToggle } from '@vueuse/core'

const isDark = useDark()
const toggleDark = useToggle(isDark)

const loading = ref(false)

const xAxis = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
const chartData = [
  {
    name: 'Revenue',
    data: [12000, 19000, 3000, 5000, 2000, 3000, 20000, 35000, 25000, 40000, 45000, 60000]
  },
  {
    name: 'Expenses',
    data: [8000, 12000, 4000, 6000, 1500, 2000, 15000, 25000, 22000, 30000, 32000, 40000]
  }
]
</script>

<template>
  <div class="min-h-screen bg-slate-50 dark:bg-slate-950 p-8 text-slate-900 dark:text-slate-100 transition-colors duration-300 font-sans">
    <div class="max-w-6xl mx-auto space-y-8">
      
      <!-- Header Area -->
      <div class="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-5">
        <div>
          <h1 class="text-3xl font-extrabold tracking-tight bg-gradient-to-r from-emerald-500 to-indigo-600 bg-clip-text text-transparent">
            Nusantara Analytics Charts Sandbox
          </h1>
          <p class="text-sm text-slate-500 dark:text-slate-400 mt-1">
            Testing the reusable ECharts wrapper components
          </p>
        </div>
        
        <div class="flex items-center gap-4">
          <!-- Toggle loading state -->
          <button 
            @click="loading = !loading"
            class="px-4 py-2 text-xs font-semibold rounded-lg border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer shadow-sm"
          >
            Toggle Loading: {{ loading ? 'ON' : 'OFF' }}
          </button>
          
          <!-- Toggle Dark Mode -->
          <button 
            @click="toggleDark()"
            class="px-4 py-2 text-xs font-semibold rounded-lg bg-emerald-600 hover:bg-emerald-500 text-white transition-colors cursor-pointer shadow-md shadow-emerald-500/20"
          >
            Switch to {{ isDark ? 'Light' : 'Dark' }} Mode
          </button>
        </div>
      </div>

      <!-- Charts Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        
        <!-- Line Chart -->
        <div class="bg-white dark:bg-slate-900 p-6 rounded-2xl border border-slate-100 dark:border-slate-800 shadow-sm space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-bold text-slate-800 dark:text-slate-200">Line Chart</h2>
            <span class="text-xs px-2.5 py-0.5 rounded-full font-medium bg-emerald-50 text-emerald-600 dark:bg-emerald-950/30 dark:text-emerald-400">Smooth Line</span>
          </div>
          <LineChart 
            :data="chartData" 
            :xAxis="xAxis" 
            title="Monthly Financial Overview" 
            unit="USD"
            :loading="loading"
          />
        </div>

        <!-- Bar Chart -->
        <div class="bg-white dark:bg-slate-900 p-6 rounded-2xl border border-slate-100 dark:border-slate-800 shadow-sm space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-bold text-slate-800 dark:text-slate-200">Bar Chart</h2>
            <span class="text-xs px-2.5 py-0.5 rounded-full font-medium bg-indigo-50 text-indigo-600 dark:bg-indigo-950/30 dark:text-indigo-400">Rounded Top Corners</span>
          </div>
          <BarChart 
            :data="chartData" 
            :xAxis="xAxis" 
            title="Revenue vs Expenses Comparison" 
            unit="USD"
            :loading="loading"
          />
        </div>

        <!-- Area Chart (Full Width) -->
        <div class="md:col-span-2 bg-white dark:bg-slate-900 p-6 rounded-2xl border border-slate-100 dark:border-slate-800 shadow-sm space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-bold text-slate-800 dark:text-slate-200">Area Chart</h2>
            <span class="text-xs px-2.5 py-0.5 rounded-full font-medium bg-amber-50 text-amber-600 dark:bg-amber-950/30 dark:text-amber-400">Gradient Fill</span>
          </div>
          <AreaChart 
            :data="chartData" 
            :xAxis="xAxis" 
            title="Historical Price Trend" 
            unit="USD"
            height="350px"
            :loading="loading"
          />
        </div>

      </div>

    </div>
  </div>
</template>
