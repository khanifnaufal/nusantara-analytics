<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from '#app'
import { useDark, useClipboard } from '@vueuse/core'
import { useCanvasStore } from '~/stores/canvas'
import { useRatesStore } from '~/stores/rates'
import { useWeatherStore } from '~/stores/weather'
import { useCommoditiesStore } from '~/stores/commodities'
import { useMarketStore } from '~/stores/market'
import { useQuakesStore } from '~/stores/quakes'
import AddChartModal from '~/components/canvas/AddChartModal.vue'
import CanvasChart from '~/components/canvas/CanvasChart.vue'
import WidgetConfigPanel from '~/components/canvas/WidgetConfigPanel.vue'

useHead({
  title: 'Nusantara Analytics - Custom Canvas'
})

// Initialize Stores
const canvasStore = useCanvasStore()
const ratesStore = useRatesStore()
const weatherStore = useWeatherStore()
const commoditiesStore = useCommoditiesStore()
const marketStore = useMarketStore()
const quakesStore = useQuakesStore()

const route = useRoute()
const router = useRouter()

// UI state
const isModalOpen = ref(false)
const isExporting = ref(false)
const isConfigOpen = ref(false)

// Active drag widget tracker
const activeDragWidgetId = ref<string | null>(null)

// 1-level deleted widget history for Ctrl+Z undo
const lastDeletedWidget = ref<any | null>(null)

const removeWidgetWithUndo = (id: string) => {
  const widget = canvasStore.widgets.find(w => w.id === id)
  if (widget) {
    lastDeletedWidget.value = JSON.parse(JSON.stringify(widget))
  }
  canvasStore.removeWidget(id)
}

const undoLastDelete = () => {
  if (lastDeletedWidget.value) {
    canvasStore.addWidget(lastDeletedWidget.value)
    lastDeletedWidget.value = null
  }
}

// Inline title editing
const editingWidgetId = ref<string | null>(null)
const editingTitle = ref('')
const titleInputRefs = ref<Record<string, HTMLInputElement | null>>({})

// Clipboard for sharing
const { copy, copied } = useClipboard()

// Keyboard shortcuts handler
const handleKeyDown = (event: KeyboardEvent) => {
  const target = event.target as HTMLElement
  if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) {
    return
  }

  // Delete/Backspace: hapus widget yang sedang selected
  if ((event.key === 'Delete' || event.key === 'Backspace') && canvasStore.selectedWidgetId) {
    removeWidgetWithUndo(canvasStore.selectedWidgetId)
  }

  // Escape: deselect widget / tutup config panel
  if (event.key === 'Escape') {
    canvasStore.selectedWidgetId = null
  }

  // Ctrl+Z: undo hapus widget terakhir
  if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 'z') {
    event.preventDefault()
    undoLastDelete()
  }
}

onMounted(async () => {
  if (typeof window !== 'undefined') {
    window.addEventListener('click', closeSidebarDropdowns)
    window.addEventListener('keydown', handleKeyDown)
  }

  // Pre-fetch all data sources to ensure components render correctly
  try {
    await Promise.all([
      ratesStore.fetchRates(),
      weatherStore.fetchWeather(),
      commoditiesStore.fetchCommodities(),
      marketStore.fetchMarket(),
      quakesStore.fetchQuakes()
    ])
  } catch (err) {
    console.error('Error fetching data sources on mount:', err)
  }

  // Load state from URL if present
  const canvasParam = route.query.canvas
  if (canvasParam && typeof canvasParam === 'string') {
    canvasStore.importFromBase64(canvasParam)
  } else {
    // If blank, add some default widgets for visualization guidance
    if (canvasStore.widgets.length === 0) {
      canvasStore.addWidget({
        title: 'Nilai Tukar USD ke IDR',
        dataSource: 'rates',
        metric: 'USD',
        chartType: 'area',
        position: { x: 40, y: 50 },
        size: { w: 420, h: 300 },
        config: { field: 'history', range: 7 }
      })
      canvasStore.addWidget({
        title: 'Suhu Kota Surabaya',
        dataSource: 'weather',
        metric: 'Surabaya',
        chartType: 'gauge',
        position: { x: 490, y: 50 },
        size: { w: 380, h: 300 },
        config: { field: 'temperature', range: 7 }
      })
    }
  }
})

onUnmounted(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('click', closeSidebarDropdowns)
    window.removeEventListener('keydown', handleKeyDown)
  }
})

// Drag logic
const startDrag = (event: PointerEvent, widget: any) => {
  // Drag only on left click
  if (event.button !== 0) return
  // Don't drag if editing the title
  if (editingWidgetId.value === widget.id) return

  canvasStore.selectedWidgetId = widget.id
  activeDragWidgetId.value = widget.id
  
  const startX = event.clientX
  const startY = event.clientY
  const startPos = { ...widget.position }
  
  const onPointerMove = (moveEvent: PointerEvent) => {
    const dx = moveEvent.clientX - startX
    const dy = moveEvent.clientY - startY
    
    // Bounds checking
    const newX = Math.max(0, startPos.x + dx)
    const newY = Math.max(0, startPos.y + dy)
    
    canvasStore.updatePosition(widget.id, newX, newY)
  }
  
  const onPointerUp = () => {
    activeDragWidgetId.value = null
    window.removeEventListener('pointermove', onPointerMove)
    window.removeEventListener('pointerup', onPointerUp)
  }
  
  window.addEventListener('pointermove', onPointerMove)
  window.addEventListener('pointerup', onPointerUp)
}

// Resize logic
const startResize = (event: PointerEvent, widget: any) => {
  if (event.button !== 0) return
  event.stopPropagation()
  
  canvasStore.selectedWidgetId = widget.id
  
  const startX = event.clientX
  const startY = event.clientY
  const startSize = { ...widget.size }
  
  const onPointerMove = (moveEvent: PointerEvent) => {
    const dx = moveEvent.clientX - startX
    const dy = moveEvent.clientY - startY
    
    // Min width and height boundaries
    const newW = Math.max(280, startSize.w + dx)
    const newH = Math.max(200, startSize.h + dy)
    
    canvasStore.updateSize(widget.id, newW, newH)
  }
  
  const onPointerUp = () => {
    window.removeEventListener('pointermove', onPointerMove)
    window.removeEventListener('pointerup', onPointerUp)
  }
  
  window.addEventListener('pointermove', onPointerMove)
  window.addEventListener('pointerup', onPointerUp)
}

// Inline Rename handlers
const startTitleEdit = (widget: any) => {
  editingWidgetId.value = widget.id
  editingTitle.value = widget.title
  nextTick(() => {
    const el = titleInputRefs.value[widget.id]
    if (el) {
      el.focus()
      el.select()
    }
  })
}

const saveTitle = (id: string) => {
  if (editingTitle.value.trim()) {
    const widget = canvasStore.widgets.find(w => w.id === id)
    if (widget) {
      widget.title = editingTitle.value.trim()
    }
  }
  editingWidgetId.value = null
}

const cancelTitleEdit = () => {
  editingWidgetId.value = null
}

// Handle adding widget from modal
const onWidgetAdd = (widgetData: any) => {
  canvasStore.addWidget(widgetData)
}

// Copy sharing link
const shareCanvas = () => {
  const base64 = canvasStore.exportToBase64()
  if (base64) {
    const shareUrl = `${window.location.origin}${window.location.pathname}?canvas=${base64}`
    copy(shareUrl)
  }
}

// Export canvas image
const exportAsPng = async () => {
  const canvasEl = document.getElementById('analytics-canvas-area')
  if (!canvasEl) return

  isExporting.value = true
  await nextTick()
  // Wait for selection rings / UI chrome to visually clear
  await new Promise(r => setTimeout(r, 200))

  try {
    const { toPng } = await import('html-to-image')

    // Determine the actual bounding box of all widgets so we crop tightly
    const widgetEls = canvasEl.querySelectorAll<HTMLElement>('.widget-card-container')
    let minX = Infinity, minY = Infinity, maxX = 0, maxY = 0
    widgetEls.forEach(el => {
      const rect = el.getBoundingClientRect()
      const parentRect = canvasEl.getBoundingClientRect()
      const x = rect.left - parentRect.left + canvasEl.scrollLeft
      const y = rect.top - parentRect.top + canvasEl.scrollTop
      minX = Math.min(minX, x)
      minY = Math.min(minY, y)
      maxX = Math.max(maxX, x + rect.width)
      maxY = Math.max(maxY, y + rect.height)
    })

    // Fall back to full area if no widgets found
    const padding = 32
    const hasWidgets = widgetEls.length > 0 && isFinite(minX)
    const cropX = hasWidgets ? Math.max(0, minX - padding) : 0
    const cropY = hasWidgets ? Math.max(0, minY - padding) : 0
    const cropW = hasWidgets ? maxX - minX + padding * 2 : canvasEl.scrollWidth
    const cropH = hasWidgets ? maxY - minY + padding * 2 : canvasEl.scrollHeight

    const dataUrl = await toPng(canvasEl, {
      backgroundColor: '#0A0A0A',
      pixelRatio: 2,
      width: cropW,
      height: cropH,
      style: {
        transform: `translate(-${cropX}px, -${cropY}px)`,
        transformOrigin: 'top left',
        width: canvasEl.scrollWidth + 'px',
        height: canvasEl.scrollHeight + 'px'
      },
      filter: (node: HTMLElement) => {
        // Exclude UI-only chrome that shouldn't appear in the export
        if (node.classList?.contains('resize-handle')) return false
        if (node.classList?.contains('widget-actions')) return false
        if (node.classList?.contains('edit-title-icon')) return false
        return true
      }
    })

    const link = document.createElement('a')
    const timestamp = new Date().toISOString().slice(0, 19).replace(/[:T]/g, '-')
    link.download = `nusantara-canvas-${timestamp}.png`
    link.href = dataUrl
    link.click()
  } catch (err) {
    console.error('Error exporting canvas PNG:', err)
  } finally {
    isExporting.value = false
  }
}

// Compute currently selected widget object for configuration panel
const selectedWidget = computed(() => {
  return canvasStore.widgets.find(w => w.id === canvasStore.selectedWidgetId) || null
})

const getWidgetIcon = (widget: any) => {
  switch (widget.dataSource) {
    case 'rates': return '💱'
    case 'weather': return '🌤'
    case 'commodities': return '📊'
    case 'market': return '📈'
    case 'quakes': return '🌋'
    default: return '📊'
  }
}

const toggleConfigPanel = () => {
  isConfigOpen.value = !isConfigOpen.value
}

// Custom dropdowns states
const activeSidebarDropdown = ref<string | null>(null)
const toggleSidebarDropdown = (name: string) => {
  activeSidebarDropdown.value = activeSidebarDropdown.value === name ? null : name
}
const closeSidebarDropdowns = () => {
  activeSidebarDropdown.value = null
}
</script>

<template>
  <div class="min-h-screen bg-[#0A0A0A] text-white flex">
    
    <!-- MOBILE FALLBACK VIEW -->
    <div class="md:hidden flex flex-col items-center justify-center min-h-screen w-full bg-[#0A0A0A] p-6 text-center select-none">
      <div class="bg-zinc-900/30 border border-white/5 p-8 rounded-3xl max-w-sm flex flex-col items-center gap-4">
        <span class="text-5xl">📊</span>
        <div>
          <h3 class="text-sm font-extrabold text-white">Canvas Builder Tidak Tersedia</h3>
          <p class="text-xs text-text-secondary mt-2 px-4 leading-relaxed">
            Canvas builder tersedia di desktop. Di mobile, gunakan dashboard utama.
          </p>
        </div>
        <NuxtLink
          to="/"
          class="w-full py-2.5 px-5 rounded-xl bg-blue-600 hover:bg-blue-500 text-white font-bold text-xs transition-all cursor-pointer shadow-md text-center"
        >
          Kembali ke Dashboard
        </NuxtLink>
      </div>
    </div>

    <!-- SIDEBAR LEFT (280px, fixed) -->
    <aside class="hidden md:flex w-[280px] border-r border-white/5 bg-[#111111] h-screen fixed left-0 top-0 flex-col z-30 select-none">
      
      <!-- Sidebar Header -->
      <div class="p-5 border-b border-white/5 flex items-center justify-between">
        <div>
          <h1 class="text-sm font-extrabold text-white tracking-wider flex items-center gap-1.5">
            <span>🎨</span> Analytics Canvas
          </h1>
          <p class="text-[10px] text-text-secondary mt-0.5 font-medium">Self-Service Workspace</p>
        </div>
        <button
          v-if="canvasStore.widgets.length > 0"
          @click="canvasStore.clearCanvas()"
          class="p-1.5 hover:bg-red-500/10 text-zinc-500 hover:text-red-400 rounded-lg transition-colors cursor-pointer"
          title="Bersihkan Canvas"
        >
          🗑️
        </button>
      </div>

      <!-- Scrollable sidebar content -->
      <div class="flex-1 overflow-y-auto p-4 space-y-6 scrollbar-thin">
        
        <!-- SECTION 1: TAMBAH WIDGET -->
        <div class="space-y-2">
          <h3 class="text-[10px] font-bold text-text-tertiary uppercase tracking-wider">Tambah Widget</h3>
          <button
            @click="isModalOpen = true"
            class="w-full py-2.5 px-4 rounded-xl bg-blue-600 hover:bg-blue-500 text-white font-bold text-xs transition-all flex items-center justify-center gap-2 active:scale-[0.98] cursor-pointer shadow-md shadow-blue-500/10 hover:shadow-blue-500/20"
          >
            <span>➕</span> Tambah Chart
          </button>
        </div>

        <!-- SECTION 2: WIDGET AKTIF -->
        <div class="space-y-2.5">
          <h3 class="text-[10px] font-bold text-text-tertiary uppercase tracking-wider">Widget Aktif ({{ canvasStore.widgets.length }})</h3>
          <div v-if="canvasStore.widgets.length === 0" class="text-center py-4 border border-dashed border-white/5 rounded-xl text-[10px] text-text-tertiary">
            Canvas kosong. Tambah widget untuk memulai.
          </div>
          <div v-else class="space-y-1.5 max-h-[45vh] overflow-y-auto pr-1 scrollbar-thin">
            <div
              v-for="w in canvasStore.widgets"
              :key="w.id"
              @click="canvasStore.selectedWidgetId = w.id"
              class="flex items-center justify-between p-2.5 rounded-lg border text-xs cursor-pointer transition-colors group"
              :class="[
                canvasStore.selectedWidgetId === w.id
                  ? 'border-blue-500/50 bg-blue-950/10 text-white'
                  : 'border-white/5 bg-[#161616] hover:bg-zinc-800 text-text-secondary'
              ]"
            >
              <span class="truncate pr-2 flex items-center gap-1.5 font-medium">
                <span>{{ getWidgetIcon(w) }}</span>
                {{ w.title }}
              </span>
              <button
                @click.stop="removeWidgetWithUndo(w.id)"
                class="text-zinc-500 hover:text-red-400 font-bold leading-none p-0.5 rounded hover:bg-white/5"
              >
                ×
              </button>
            </div>
          </div>
        </div>

        <!-- SECTION 4: EXPORT & SHARE -->
        <div class="space-y-2 border-t border-white/5 pt-4">
          <h3 class="text-[10px] font-bold text-text-tertiary uppercase tracking-wider">Export & Share</h3>
          <div class="grid grid-cols-2 gap-2">
            <button
              @click="exportAsPng"
              :disabled="canvasStore.widgets.length === 0"
              class="py-2 px-3 rounded-lg border border-white/8 hover:border-white/20 hover:bg-zinc-800 text-zinc-300 hover:text-white text-[11px] font-bold transition-all flex items-center justify-center gap-1 active:scale-[0.98] cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed"
            >
              📥 PNG
            </button>
            <button
              @click="shareCanvas"
              class="py-2 px-3 rounded-lg border border-white/8 hover:border-white/20 hover:bg-zinc-800 text-zinc-300 hover:text-white text-[11px] font-bold transition-all flex items-center justify-center gap-1 active:scale-[0.98] cursor-pointer"
            >
              <span v-if="copied" class="text-blue-400">✓ Tersalin</span>
              <span v-else class="flex items-center gap-1">🔗 Salin Link</span>
            </button>
          </div>
        </div>

      </div>

      <!-- Sidebar Footer / Back to Dashboard Link -->
      <div class="p-4 border-t border-white/5 bg-[#0A0A0A]">
        <NuxtLink
          to="/"
          class="w-full py-2 px-4 rounded-xl border border-white/5 hover:border-white/10 hover:bg-zinc-800 text-zinc-400 hover:text-white transition-all text-xs font-bold flex items-center justify-center gap-2"
        >
          <span>←</span> Dashboard Utama
        </NuxtLink>
      </div>

    </aside>

    <!-- CANVAS WORKSPACE AREA (flex-1) -->
    <main class="hidden md:flex ml-[280px] flex-1 min-h-screen bg-[#0A0A0A] flex-col relative overflow-hidden">
      
      <!-- Backdrop Dot Grid Pattern -->
      <div 
        class="absolute inset-0 pointer-events-none z-0"
        style="
          background-image: radial-gradient(circle, rgba(255,255,255,0.06) 1px, transparent 1px);
          background-size: 20px 20px;
        "
      ></div>

      <!-- Floating Canvas Info banner -->
      <div class="absolute top-4 right-4 z-20 flex items-center gap-3 select-none">
        <div class="bg-[#111111]/90 backdrop-blur-md px-3.5 py-1.5 rounded-full border border-white/8 text-[10px] font-bold text-text-secondary shadow-lg">
          <span>✨</span> Drag widget dari header, resize dari pojok kanan bawah
        </div>
      </div>

      <!-- Canvas container wrapper -->
      <div 
        id="analytics-canvas-area" 
        class="flex-1 relative overflow-auto p-8 min-h-screen z-10 scrollbar-thin cursor-default"
        :class="{ 'export-mode': isExporting }"
        @pointerdown="canvasStore.selectedWidgetId = null"
      >
        
        <!-- Empty State Canvas -->
        <div 
          v-if="canvasStore.widgets.length === 0" 
          class="absolute inset-0 flex flex-col items-center justify-center text-center p-8 select-none pointer-events-none z-10"
        >
          <div class="flex flex-col items-center max-w-sm">
            <span class="text-6xl mb-4">📊</span>
            <h3 class="text-lg font-bold text-white mb-2">Canvas Kosong</h3>
            <p class="text-xs text-text-secondary mb-6 px-4 leading-relaxed">
              Klik '+ Tambah Chart' untuk mulai membangun dashboard kustom
            </p>
            <button
              @click.stop="isModalOpen = true"
              class="pointer-events-auto py-2.5 px-6 rounded-xl bg-blue-600 hover:bg-blue-500 text-white font-bold text-xs transition-all cursor-pointer shadow-lg shadow-blue-500/10 hover:shadow-blue-500/20 active:scale-95"
            >
              + Tambah Chart Pertama
            </button>
          </div>
        </div>

        <!-- DRAGGABLE & RESIZABLE WIDGETS -->
        <div
          v-for="widget in canvasStore.widgets"
          :key="widget.id"
          class="absolute widget-card-container transition-all duration-150 ease-out"
          :class="[
            canvasStore.selectedWidgetId === widget.id ? 'ring-2 ring-blue-500/50 ring-offset-2 ring-offset-[#0A0A0A]' : '',
            activeDragWidgetId === widget.id ? 'opacity-90 shadow-2xl shadow-blue-500/10' : 'shadow-lg'
          ]"
          :style="{
            left: `${widget.position.x}px`,
            top: `${widget.position.y}px`,
            width: `${widget.size.w}px`,
            height: `${widget.size.h}px`,
            zIndex: activeDragWidgetId === widget.id ? 50 : (canvasStore.selectedWidgetId === widget.id ? 20 : 10)
          }"
          @pointerdown.stop="canvasStore.selectedWidgetId = widget.id"
        >
          <div class="h-full flex flex-col bg-[#111111] border border-white/8 rounded-2xl overflow-hidden shadow-2xl">
            
            <!-- Header (Drag Handle) -->
            <div
              class="drag-handle flex justify-between items-center px-4 py-3 border-b border-white/5 select-none"
              :class="[
                activeDragWidgetId === widget.id ? 'cursor-grabbing' : 'cursor-grab'
              ]"
              @pointerdown.stop="startDrag($event, widget)"
            >
              <div class="flex items-center gap-2 min-w-0 flex-1">
                <span class="text-sm shrink-0">{{ getWidgetIcon(widget) }}</span>
                
                <!-- Editable title wrapper -->
                <input
                  v-if="editingWidgetId === widget.id"
                  :ref="(el) => { if (el) titleInputRefs[widget.id] = el as HTMLInputElement }"
                  v-model="editingTitle"
                  class="bg-zinc-800 border border-blue-500/50 rounded px-1.5 py-0.5 text-xs text-white focus:outline-none w-full"
                  @blur="saveTitle(widget.id)"
                  @keydown.enter="saveTitle(widget.id)"
                  @keydown.esc="cancelTitleEdit()"
                />
                <h4
                  v-else
                  class="text-xs font-bold text-white truncate cursor-pointer hover:text-blue-400 flex items-center gap-1 group/title"
                  @click="startTitleEdit(widget)"
                >
                  <span>{{ widget.title }}</span>
                  <span class="text-[9px] text-text-tertiary opacity-0 group-hover/title:opacity-100 transition-opacity edit-title-icon">✏️</span>
                </h4>
              </div>
              
              <!-- Action tools -->
              <div class="flex items-center gap-1.5 ml-2 widget-actions">
                <button
                  @click.stop="canvasStore.selectedWidgetId = widget.id"
                  class="p-1 text-zinc-500 hover:text-white rounded hover:bg-white/5 transition-colors cursor-pointer text-xs"
                  title="Edit Widget"
                >
                  ⚙️
                </button>
                <button
                  @click.stop="removeWidgetWithUndo(widget.id)"
                  class="p-1 text-zinc-500 hover:text-red-400 rounded hover:bg-white/5 transition-colors cursor-pointer font-bold text-sm leading-none"
                  title="Hapus"
                >
                  ×
                </button>
              </div>
            </div>
            
            <!-- Body area -->
            <div class="flex-1 p-3 min-h-0 relative">
              <CanvasChart :widget="widget" />
            </div>
            
            <!-- Resize Handle -->
            <div
              class="absolute bottom-1 right-1 w-4 h-4 cursor-se-resize flex items-end justify-end select-none resize-handle"
              @pointerdown.stop.prevent="startResize($event, widget)"
            >
              <svg width="8" height="8" viewBox="0 0 8 8" class="text-zinc-500 opacity-40 hover:opacity-100 transition-opacity">
                <path d="M6 0 L8 0 L8 8 L0 8 L0 6 L4 6 L4 4 L6 4 Z" fill="currentColor"/>
              </svg>
            </div>
          </div>
        </div>

      </div>

    </main>

    <!-- SLIDE-IN CONFIGURATION PANEL -->
    <WidgetConfigPanel :is-open="canvasStore.selectedWidgetId !== null" />

    <!-- MULTI-STEP ADD CHART MODAL -->
    <AddChartModal
      :is-open="isModalOpen"
      @close="isModalOpen = false"
      @add="onWidgetAdd"
    />

  </div>
</template>

<style>
/* CSS styles to toggle during high resolution html2canvas PNG export */
.export-mode .resize-handle {
  display: none !important;
}
.export-mode .widget-actions {
  display: none !important;
}
.export-mode .edit-title-icon {
  display: none !important;
}
.export-mode .widget-card-container {
  box-shadow: none !important;
  outline: none !important;
}
</style>
