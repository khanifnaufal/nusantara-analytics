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
  <div class="relative min-h-screen bg-[#0A0A0A] text-slate-100 p-4 sm:p-6 md:p-8">
    <!-- Decorative background glow blob -->
    <div class="absolute top-0 left-1/2 -translate-x-1/2 w-[600px] h-[300px] bg-blue-500/5 rounded-full blur-3xl pointer-events-none z-0"></div>

    <div class="max-w-7xl mx-auto space-y-8 relative z-10">
      
      <!-- Top Header Area - Sticky Navbar style -->
      <header class="sticky top-4 z-40 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 p-4 md:p-5 rounded-2xl border-b border-white/5 bg-[#0A0A0A]/90 backdrop-blur-xl shadow-lg shadow-black/10">
        <div>
          <div class="flex items-center gap-2.5">
            <span class="text-2xl" aria-hidden="true">🇮🇩</span>
            <h1 class="text-xl md:text-2xl font-extrabold tracking-tight text-transparent bg-clip-text bg-gradient-to-r from-white to-slate-400">
              Nusantara Analytics
            </h1>
          </div>
          <p class="text-[11px] md:text-xs text-text-secondary mt-1 font-medium">
            Visualisasi Grafik Kustom untuk Komoditas, Kurs Mata Uang, Pasar Saham, & Cuaca Regional.
          </p>
        </div>

        <!-- System Navigation & Status -->
        <div class="flex flex-wrap items-center gap-3 self-start sm:self-auto">
          <!-- Back to Dashboard -->
          <NuxtLink
            to="/"
            class="inline-flex items-center gap-2 px-4 py-2 text-xs font-bold rounded-xl border border-white/10 bg-transparent hover:border-white/20 text-zinc-400 hover:text-white transition-all shadow-sm"
          >
            <span>🏠</span> Dashboard Utama
          </NuxtLink>

          <!-- Live status pill -->
          <div class="inline-flex items-center gap-2 bg-[#111111] border border-white/5 py-1.5 px-3.5 rounded-xl text-xs font-bold text-text-secondary shadow-sm">
            <span class="relative flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-500 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-blue-500"></span>
            </span>
            <span class="tracking-wide">Live</span>
            <span class="text-white/10" aria-hidden="true">|</span>
            <span class="font-mono text-blue-400 text-[11px]">{{ liveTime }}</span>
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
