<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import uPlot from 'uplot'
import 'uplot/dist/uPlot.min.css'
import { withAlpha } from '../format'

const props = defineProps<{
  data: number[] | null
  accent: string
}>()

const host = ref<HTMLDivElement | null>(null)
let plot: uPlot | null = null
let observer: ResizeObserver | null = null

function alignedData(): uPlot.AlignedData {
  const values = props.data ?? []
  // The x axis is a sample index; the trend is a fixed-length rolling window,
  // so real timestamps would add nothing the eye can use.
  return [values.map((_, index) => index), values]
}

function plotOptions(width: number, height: number): uPlot.Options {
  const accent = props.accent

  return {
    width,
    height,
    padding: [2, 0, 0, 0],
    cursor: { show: false },
    legend: { show: false },
    scales: {
      x: { time: false },
      y: { range: [0, 100] },
    },
    axes: [{ show: false }, { show: false }],
    series: [
      {},
      {
        stroke: accent,
        width: 1.5,
        points: { show: false },
        fill: (u: uPlot) => {
          const gradient = u.ctx.createLinearGradient(0, 0, 0, u.height)
          gradient.addColorStop(0, withAlpha(accent, 0.38))
          gradient.addColorStop(1, withAlpha(accent, 0.02))
          return gradient
        },
      },
    ],
  }
}

function dimensions(): { width: number; height: number } {
  const element = host.value
  return {
    width: Math.max(1, Math.floor(element?.clientWidth ?? 1)),
    height: Math.max(1, Math.floor(element?.clientHeight ?? 1)),
  }
}

function mount() {
  if (!host.value) {
    return
  }
  const { width, height } = dimensions()
  plot = new uPlot(plotOptions(width, height), alignedData(), host.value)
}

function destroy() {
  plot?.destroy()
  plot = null
}

function resize() {
  if (!plot) {
    return
  }
  const { width, height } = dimensions()
  plot.setSize({ width, height })
}

onMounted(() => {
  mount()
  observer = new ResizeObserver(resize)
  if (host.value) {
    observer.observe(host.value)
  }
})

onBeforeUnmount(() => {
  observer?.disconnect()
  observer = null
  destroy()
})

watch(
  () => props.data,
  () => {
    plot?.setData(alignedData())
  },
)

// The accent is baked into the series stroke and fill gradient, so a theme
// change means rebuilding the chart rather than mutating it.
watch(
  () => props.accent,
  () => {
    destroy()
    mount()
  },
)
</script>

<template>
  <div ref="host" class="chart-host"></div>
</template>

<style scoped>
.chart-host {
  width: 100%;
  height: 100%;
}
</style>
