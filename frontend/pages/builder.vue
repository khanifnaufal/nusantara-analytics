<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import CustomVizBuilder from '~/components/widgets/CustomVizBuilder.vue'

// Live clock state
const liveTime = ref('')
let clockTimer: any = null

const updateClock = () => {
  liveTime.value = new Date().toLocaleString('id-ID', {
    dateStyle: 'medium',
    timeStyle: 'medium'
  })
}

onMounted(() => {
  updateClock()
  clockTimer = setInterval(updateClock, 1000)
})

onUnmounted(() => {
  if (clockTimer) clearInterval(clockTimer)
})
</script>

<template>
  <div class="relative min-h-screen bg-slate-950 text-slate-100 p-4 sm:p-6 md:p-8">
    <!-- Premium Decorative background glow blobs wrapped to prevent page overflow -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none z-0">
      <div class="absolute -top-40 -left-40 w-96 h-96 bg-indigo-500/10 rounded-full blur-3xl"></div>
      <div class="absolute top-1/2 -right-40 w-96 h-96 bg-violet-500/5 rounded-full blur-3xl"></div>
      <div class="absolute -bottom-40 left-1/3 w-96 h-96 bg-emerald-500/5 rounded-full blur-3xl"></div>
    </div>

    <div class="max-w-7xl mx-auto space-y-8 relative z-10">
      
      <!-- Top Header Area -->
      <header class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-6 border-b border-slate-800/60">
        <div>
          <div class="flex items-center gap-2">
            <span class="text-3xl" aria-hidden="true">🇮🇩</span>
            <h1 class="text-2xl md:text-3xl font-extrabold tracking-tight bg-linear-to-r from-white via-slate-100 to-slate-400 bg-clip-text text-transparent">
              Nusantara Analytics
            </h1>
          </div>
          <p class="text-xs md:text-sm text-slate-400 mt-1.5 font-medium">
            Visualisasi Grafik Kustom untuk Komoditas, Kurs Mata Uang, Pasar Saham, & Cuaca Regional.
          </p>
        </div>

        <!-- System Navigation & Clock -->
        <div class="flex flex-wrap items-center gap-3 self-start sm:self-auto">
          <!-- Back to Dashboard -->
          <NuxtLink
            to="/"
            class="inline-flex items-center gap-2 px-4 py-2 text-xs font-bold rounded-xl border border-slate-800 bg-slate-900/40 hover:bg-slate-900 hover:border-slate-700 text-slate-200 transition-all shadow-sm"
          >
            <span>🏠</span> Dashboard Utama
          </NuxtLink>

          <div class="inline-flex items-center gap-3 bg-slate-900/60 border border-slate-850/80 backdrop-blur-md py-2 px-4 rounded-xl text-xs font-semibold text-slate-300 shadow-sm">
            <span class="relative flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            </span>
            <span class="tracking-wide">Sistem Online</span>
            <span class="text-slate-700" aria-hidden="true">|</span>
            <span class="font-mono text-slate-200">{{ liveTime }}</span>
          </div>
        </div>
      </header>

      <!-- Custom Builder Widget -->
      <div class="w-full">
        <CustomVizBuilder />
      </div>

    </div>
  </div>
</template>
