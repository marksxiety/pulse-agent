<script setup lang="ts">
import { computed } from 'vue'
import TrendChart from './TrendChart.vue'
import type { CardColors, Series } from '../types'

const props = defineProps<{
  label: string
  icon: string
  series: Series
  colors: CardColors
  detail: string
  tooltip: string
}>()

const clamped = computed(() => Math.min(100, Math.max(0, props.series.percent)))

// Mirrors the terminal UI so the same machine reads the same way in both
// frontends: amber at 75%, red at 90%.
const activeColor = computed(() => {
  if (props.series.percent >= 90) {
    return props.colors.danger
  }
  if (props.series.percent >= 75) {
    return props.colors.warn
  }
  return props.colors.accent
})

const hasTrend = computed(() => (props.series.trend?.length ?? 0) >= 2)
</script>

<template>
  <div class="card" :title="tooltip">
    <div class="card-head">
      <span class="card-label">{{ icon }} {{ label }}</span>
      <span class="card-detail">{{ detail }}</span>
      <span class="card-percent" :style="{ color: activeColor }">
        {{ series.percent.toFixed(1) }}%
      </span>
    </div>

    <div class="card-bar">
      <div
        class="card-bar-fill"
        :style="{ width: `${clamped}%`, backgroundColor: activeColor }"
      ></div>
    </div>

    <div class="card-chart">
      <TrendChart v-if="hasTrend" :data="series.trend" :accent="colors.accent" />
    </div>
  </div>
</template>
