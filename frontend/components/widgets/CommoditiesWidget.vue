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
</script>

<template>
  <div
    class="relative w-full overflow-hidden rounded-2xl border border-white/6 bg-[#161B22] p-5 md:p-6 transition-all duration-300 hover:border-emerald-500/30 hover:shadow-[0_0_20px_rgba(16,185,129,0.08)]"
  >
    <!-- Gradient Accent Top Border -->
    <div class="absolute top-0 left-0 h-1 w-full bg-gradient-to-r from-emerald-500 to-teal-500"></div>

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
        class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-lg border border-white/5 bg-[#1C2128] hover:bg-[#21262D] text-text-secondary hover:text-text-primary cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed transition-all"
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
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
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
          class="h-[320px] rounded-xl border border-white/5 bg-[#1C2128] animate-pulse p-4 flex flex-col justify-between"
        >
          <div>
            <div class="flex justify-between items-center mb-4">
              <div
                class="h-4 bg-white/5 rounded w-1/3"
              ></div>
              <div
                class="h-4 bg-white/5 rounded w-1/4"
              ></div>
            </div>
            <div
              class="h-6 bg-white/5 rounded w-1/2 mb-4"
            ></div>
          </div>
          <!-- Chart Placeholder -->
          <div
            class="h-[160px] bg-white/5 rounded-lg w-full flex items-center justify-center"
          >
            <span class="text-xs text-text-tertiary">Memuat grafik...</span>
          </div>
        </div>
      </template>

      <!-- Empty State -->
      <div
        v-else-if="
          !commoditiesStore.data ||
          !commoditiesStore.data.commodities ||
          commoditiesStore.data.commodities.length === 0
        "
        class="col-span-full py-12 text-center text-text-tertiary"
      >
        Tidak ada data komoditas tersedia
      </div>

      <!-- Real Cards -->
      <template v-else>
        <div
          v-for="item in commoditiesStore.data.commodities"
          :key="item.symbol"
          class="rounded-xl border border-white/5 bg-[#1C2128] p-4 flex flex-col justify-between transition-all duration-300 hover:border-emerald-500/20"
        >
          <!-- Price and Badge -->
          <div class="mb-2">
            <div class="flex items-center justify-between mb-1">
              <div>
                <h4
                  class="font-bold text-sm text-text-primary tracking-tight"
                >
                  {{ item.name }}
                </h4>
                <p
                  class="text-[9px] text-text-tertiary font-semibold uppercase font-mono tracking-wider"
                >
                  {{ item.symbol }}
                </p>
              </div>

              <!-- Change Badge -->
              <span
                :class="[
                  'inline-flex items-center px-2 py-0.5 rounded text-[10px] font-bold border font-mono',
                  item.changePercent >= 0
                    ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20'
                    : 'bg-rose-500/10 text-rose-400 border-rose-500/20',
                ]"
              >
                {{ item.changePercent >= 0 ? "+" : ""
                }}{{ item.changePercent.toFixed(2) }}%
              </span>
            </div>

            <!-- Price -->
            <div
              class="text-xl font-extrabold text-text-primary tracking-tight font-mono mt-1"
            >
              {{ formatPrice(item.price, item.currency) }}
            </div>
          </div>

          <!-- History Chart (Last 7 Days) -->
          <div class="w-full mt-2 border-t border-white/5 pt-3">
            <p
              class="text-[9px] text-text-tertiary font-bold uppercase tracking-wider mb-2"
            >
              Tren 7 Hari Terakhir
            </p>
            <div class="h-[160px] w-full">
              <AreaChart
                v-bind="getChartProps(item.history, item.name)"
                :loading="commoditiesStore.loading"
                :unit="item.currency"
                accentColor="#10B981"
                height="100%"
              />
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
