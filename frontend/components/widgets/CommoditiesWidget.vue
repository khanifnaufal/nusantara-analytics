<script setup lang="ts">
import { computed } from "vue";
import { useCommoditiesStore } from "~/stores/commodities";
import { usePolling } from "~/composables/usePolling";
import { useExportCsv } from "~/composables/useExportCsv";
import AreaChart from "../charts/AreaChart.vue";

const props = withDefaults(defineProps<{
  disableAutoPoll?: boolean
}>(), {
  disableAutoPoll: false
})

const commoditiesStore = useCommoditiesStore();

// Auto-fetch and poll commodities every 60s
if (!props.disableAutoPoll) {
  usePolling(() => commoditiesStore.fetchCommodities(), 60000);
}

// Helper to format prices according to currency
const formatPrice = (value: number, currency: string) => {
  if (!Number.isFinite(value)) {
    return "-";
  }
  if (currency === "USD") {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }).format(value);
  } else if (currency === "MYR") {
    return new Intl.NumberFormat("en-MY", {
      style: "currency",
      currency: "MYR",
    })
      .format(value)
      .replace("MYR", "RM");
  } else if (currency === "IDR") {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      maximumFractionDigits: 0,
    }).format(value);
  }
  return `${currency} ${new Intl.NumberFormat(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(value)}`;
};

// Generate chart data for a commodity
const getChartProps = (
  history: { timestamp: string; close: number }[],
  name: string,
) => {
  if (!history || history.length === 0) {
    return {
      xAxis: [] as string[],
      data: [] as { name: string; data: number[] }[],
    };
  }

  // Sort history by timestamp ascending to make sure chart draws correctly left-to-right
  const sortedHistory = [...history].sort(
    (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime(),
  );

  const xAxis = sortedHistory.map((h) => {
    const d = new Date(h.timestamp);
    return d.toLocaleDateString("id-ID", { day: "2-digit", month: "short" });
  });

  const data = [
    {
      name,
      data: sortedHistory.map((h) => parseFloat(h.close.toFixed(2))),
    },
  ];

  return { xAxis, data };
};

// Export CSV Setup
const csvData = computed(() => {
  if (!commoditiesStore.data?.commodities) return [];
  return commoditiesStore.data.commodities.map((c) => ({
    Simbol: c.symbol,
    Nama: c.name,
    Harga: c.price,
    MataUang: c.currency,
    Perubahan: c.change,
    "Persen Perubahan (%)": c.changePercent.toFixed(2) + "%",
    "Update Terakhir": c.lastUpdated
      ? new Date(c.lastUpdated).toLocaleString("id-ID")
      : "-",
  }));
});

const { exportCsv } = useExportCsv(csvData, "komoditas_harga_analytics");

// Dynamic Y-axis min/max calculation with 5% padding
const yAxisMinFn = (value: any) => {
  if (!value || typeof value.min !== 'number' || isNaN(value.min) || !isFinite(value.min)) {
    return 'dataMin'
  }
  const range = value.max - value.min
  if (range === 0) return value.min * 0.95
  const minVal = value.min - range * 0.05
  return value.min >= 0 && minVal < 0 ? 0 : minVal
}

const yAxisMaxFn = (value: any) => {
  if (!value || typeof value.max !== 'number' || isNaN(value.max) || !isFinite(value.max)) {
    return 'dataMax'
  }
  const range = value.max - value.min
  if (range === 0) return value.max * 1.05
  return value.max + range * 0.05
}

interface CommodityDisplayItem {
  symbol: string
  name: string
  price: number
  change: number
  changePercent: number
  currency: string
  history: { timestamp: string; close: number }[]
  lastUpdated?: string
  isEmptyPlaceholder?: boolean
}

const displayCommodities = computed((): CommodityDisplayItem[] => {
  if (!commoditiesStore.data?.commodities || commoditiesStore.data.commodities.length === 0) {
    return [];
  }
  
  const list = [...commoditiesStore.data.commodities] as CommodityDisplayItem[];
  const hasCpo = list.some(c => c.symbol === 'PALM.KL' || c.name.toLowerCase() === 'cpo');
  
  if (!hasCpo) {
    list.push({
      symbol: 'PALM.KL',
      name: 'CPO',
      price: 0,
      change: 0,
      changePercent: 0,
      currency: 'MYR',
      history: [],
      isEmptyPlaceholder: true
    });
  }
  
  return list;
});
</script>

<template>
  <div class="widget-card">

    <!-- Header -->
    <div class="flex items-center justify-between mb-6 mt-1">
      <div>
        <h3
          class="text-base font-bold text-text-primary flex items-center gap-2"
        >
          <span class="text-emerald-400">📈</span> Komoditas Global
        </h3>
        <p class="text-xs text-text-tertiary mt-1">
          Harga pasar komoditas dunia real-time • Update terakhir:
          {{
            commoditiesStore.data?.lastUpdated
              ? new Date(commoditiesStore.data.lastUpdated).toLocaleTimeString(
                  "id-ID",
                )
              : "-"
          }}
        </p>
      </div>

      <!-- Export Button -->
      <button
        @click="exportCsv()"
        :disabled="
          !commoditiesStore.data?.commodities ||
          commoditiesStore.data.commodities.length === 0
        "
        class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-lg bg-transparent border border-white/10 hover:border-white/20 text-zinc-400 hover:text-white cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed transition-all"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-3.5 w-3.5"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="2"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
          />
        </svg>
        Export CSV
      </button>
    </div>

    <!-- Error State -->
    <div
      v-if="commoditiesStore.error"
      class="bg-rose-500/10 border border-rose-500/20 rounded-xl p-4 mb-4 text-xs text-rose-400"
    >
      <div class="flex items-center gap-2 font-semibold">
        <span>⚠️</span> Gagal memuat data komoditas
      </div>
      <p class="mt-1 opacity-90 font-mono">{{ commoditiesStore.error }}</p>
    </div>

    <!-- Commodities Grid -->
    <div class="grid grid-cols-1 gap-4">
      <!-- Loading Skeleton State -->
      <template
        v-if="
          commoditiesStore.loading &&
          (!commoditiesStore.data || !commoditiesStore.data.commodities)
        "
      >
        <div
          v-for="n in 3"
          :key="'skeleton-comm-' + n"
          class="min-h-[120px] rounded-xl border border-white/5 bg-bg-subcard animate-pulse p-4 flex gap-4 justify-between"
        >
          <!-- Left Skeleton Info -->
          <div class="w-28 shrink-0 flex flex-col justify-between">
            <div class="space-y-2">
              <div class="h-4 bg-white/5 rounded w-2/3"></div>
              <div class="h-3 bg-white/5 rounded w-1/2"></div>
            </div>
            <div class="space-y-2 mt-4">
              <div class="h-6 bg-white/5 rounded w-3/4"></div>
              <div class="h-4 bg-white/5 rounded w-1/2"></div>
            </div>
          </div>
          <!-- Right Skeleton Chart -->
          <div class="grow bg-white/5 rounded-lg h-[80px] self-center"></div>
        </div>
      </template>

      <!-- Empty State -->
      <div
        v-else-if="
          !displayCommodities ||
          displayCommodities.length === 0
        "
        class="col-span-full py-12 text-center text-text-tertiary"
      >
        Tidak ada data komoditas tersedia
      </div>

      <!-- Real Cards -->
      <template v-else>
        <div
          v-for="item in displayCommodities"
          :key="item.symbol"
          class="rounded-xl border border-white/5 bg-bg-subcard p-4 flex gap-4 transition-all duration-300 hover:border-blue-500/20 min-h-[120px]"
        >
          <template v-if="item.isEmptyPlaceholder">
            <!-- Left Info -->
            <div class="flex flex-col justify-between w-28 shrink-0">
              <div>
                <h4 class="font-bold text-sm text-text-primary tracking-tight">
                  {{ item.name }}
                </h4>
                <p class="text-[9px] text-text-tertiary font-semibold uppercase font-mono tracking-wider">
                  {{ item.symbol }}
                </p>
              </div>
            </div>
            <!-- Right Centered Message -->
            <div class="grow flex items-center justify-center border-l border-white/5 pl-4 py-8">
              <span class="text-xs text-text-tertiary font-medium text-center">
                Data CPO sedang tidak tersedia
              </span>
            </div>
          </template>

          <template v-else>
            <!-- Left: Info -->
            <div class="flex flex-col justify-between w-28 shrink-0">
              <div>
                <div class="flex items-center justify-between mb-1">
                  <h4 class="font-bold text-sm text-text-primary tracking-tight truncate">
                    {{ item.name }}
                  </h4>
                </div>
                <p class="text-[9px] text-text-tertiary font-semibold uppercase font-mono tracking-wider">
                  {{ item.symbol }}
                </p>
              </div>

              <div class="mt-2">
                <div class="text-lg font-mono font-bold text-text-primary tracking-tight">
                  {{ formatPrice(item.price, item.currency) }}
                </div>
                <div class="mt-1">
                  <span
                    :class="[
                      'inline-flex items-center px-1.5 py-0.5 rounded text-[9px] font-bold font-mono',
                      item.changePercent > 0
                        ? 'bg-green-500/15 text-green-400'
                        : item.changePercent < 0
                          ? 'bg-red-500/15 text-red-400'
                          : 'bg-zinc-800 text-zinc-500',
                    ]"
                  >
                    {{ item.changePercent > 0 ? "+" : "" }}{{ item.changePercent.toFixed(2) }}%
                  </span>
                </div>
              </div>
            </div>

            <!-- Right: Sparkline Chart -->
            <div class="grow h-[80px] self-center border-l border-white/5 pl-4">
              <p class="text-[9px] text-text-tertiary font-bold uppercase tracking-wider mb-1">
                TREN 7 HARI
              </p>
              <div class="h-[65px] w-full">
                <div
                  v-if="!commoditiesStore.loading && (!item.history || item.history.length === 0)"
                  class="w-full h-full flex items-center justify-center rounded-lg bg-white/1 text-center"
                >
                  <span class="text-[10px] text-text-tertiary font-medium">Data tidak tersedia</span>
                </div>
                <AreaChart
                  v-else
                  v-bind="getChartProps(item.history, item.name)"
                  :loading="commoditiesStore.loading"
                  :unit="item.currency"
                  accentColor="#3B82F6"
                  height="100%"
                  :lineWidth="1.5"
                  :smooth="true"
                  :areaOpacity="0.1"
                  :showOnlyFirstLastX="true"
                  :yAxisMin="yAxisMinFn"
                  :yAxisMax="yAxisMaxFn"
                  :yAxisSplitNumber="2"
                  :isSparkline="true"
                />
              </div>
            </div>
          </template>
        </div>
      </template>
    </div>
  </div>
</template>
