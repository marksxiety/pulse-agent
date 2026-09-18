<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import MetricCard from './components/MetricCard.vue'
import ThemePicker from './components/ThemePicker.vue'
import { bytesToGB, formatUptime } from './format'
import type { Bootstrap, CardColors, Series, Snapshot, Theme } from './types'
import { EVENT_SNAPSHOT } from './types'
import { GetBootstrap, Quit, SetAlwaysOnTop, SetTheme } from '../wailsjs/go/main/App'
import { EventsOn } from '../wailsjs/runtime/runtime'

const EMPTY_SERIES: Series = { percent: 0, trend: null, max: 0, avg: 0 }

// Rendering from a zeroed snapshot keeps the card layout stable during the
// first second, before the backend has pushed anything.
const EMPTY_SNAPSHOT: Snapshot = {
  cpu: EMPTY_SERIES,
  memory: EMPTY_SERIES,
  disk: EMPTY_SERIES,
  cpuCores: 0,
  memoryUsed: 0,
  memoryTotal: 0,
  diskUsed: 0,
  diskTotal: 0,
  uptimeSecs: 0,
}

const bootstrap = ref<Bootstrap | null>(null)
const snapshot = ref<Snapshot>(EMPTY_SNAPSHOT)
const activeTheme = ref<Theme | null>(null)
const showThemes = ref(false)
const pinned = ref(true)

let stopSnapshots: (() => void) | null = null

// One computed for all three cards so each card receives a stable object
// identity, and so the warn/danger endpoints are resolved in one place.
const cardColors = computed<Record<'cpu' | 'memory' | 'disk', CardColors>>(() => {
  const theme = activeTheme.value
  const warn = theme?.warn ?? ''
  const danger = theme?.danger ?? ''

  return {
    cpu: { accent: theme?.cpu_accent ?? '', warn, danger },
    memory: { accent: theme?.mem_accent ?? '', warn, danger },
    disk: { accent: theme?.disk_accent ?? '', warn, danger },
  }
})

const memorySize = computed(() =>
  snapshot.value.memoryTotal > 0
    ? `${bytesToGB(snapshot.value.memoryUsed)} / ${bytesToGB(snapshot.value.memoryTotal)}`
    : '',
)

const diskSize = computed(() =>
  snapshot.value.diskTotal > 0
    ? `${bytesToGB(snapshot.value.diskUsed)} / ${bytesToGB(snapshot.value.diskTotal)}`
    : '',
)

function applyTheme(theme: Theme) {
  activeTheme.value = theme

  const root = document.documentElement.style
  root.setProperty('--bg', theme.bg)
  root.setProperty('--surface', theme.surface)
  root.setProperty('--border', theme.border)
  root.setProperty('--muted', theme.muted)
  root.setProperty('--subtle', theme.subtle)
  root.setProperty('--text', theme.text)
  root.setProperty('--dim', theme.dim)
  root.setProperty('--cpu', theme.cpu_accent)
  root.setProperty('--mem', theme.mem_accent)
  root.setProperty('--disk', theme.disk_accent)
  root.setProperty('--warn', theme.warn)
  root.setProperty('--danger', theme.danger)
}

async function chooseTheme(name: string) {
  const applied = await SetTheme(name)
  applyTheme(applied)
  showThemes.value = false
}

async function togglePinned() {
  pinned.value = !pinned.value
  await SetAlwaysOnTop(pinned.value)
}

function peakDetail(series: Series): string {
  return `pk ${series.max.toFixed(1)}%`
}

function seriesTooltip(name: string, series: Series, size: string): string {
  const parts = [size, `avg ${series.avg.toFixed(1)}%`, `peak ${series.max.toFixed(1)}%`]
  return `${name}\n${parts.filter(Boolean).join(' · ')}`
}

onMounted(async () => {
  const boot = await GetBootstrap()
  bootstrap.value = boot
  pinned.value = boot.alwaysOnTop

  const initial =
    boot.themes.find((candidate) => candidate.name === boot.theme) ?? boot.themes[0]
  if (initial) {
    applyTheme(initial)
  }

  stopSnapshots = EventsOn(EVENT_SNAPSHOT, (next: Snapshot) => {
    snapshot.value = next
  })
})

onBeforeUnmount(() => {
  stopSnapshots?.()
  stopSnapshots = null
})
</script>

<template>
  <div class="widget">
    <div class="titlebar">
      <span class="brand" :title="`Pulse Agent ${bootstrap?.version ?? ''}`"
        >&#9670; PULSE <span>AGENT</span></span
      >
      <span class="uptime">{{ formatUptime(snapshot.uptimeSecs) }}</span>

      <span class="actions">
        <button
          type="button"
          class="icon-btn"
          :class="{ active: pinned }"
          :title="pinned ? 'Unpin from top' : 'Pin on top'"
          @click="togglePinned"
        >
          {{ pinned ? '\u25C9' : '\u25CE' }}
        </button>
        <button
          type="button"
          class="icon-btn"
          :class="{ active: showThemes }"
          title="Theme"
          @click="showThemes = !showThemes"
        >
          &#9681;
        </button>
        <button type="button" class="icon-btn close" title="Close" @click="Quit()">
          &#10005;
        </button>
      </span>
    </div>

    <div v-if="showThemes" class="scrim" @click="showThemes = false"></div>
    <ThemePicker
      v-if="showThemes && bootstrap"
      :themes="bootstrap.themes"
      :active="activeTheme?.name ?? ''"
      @select="chooseTheme"
    />

    <div class="cards">
      <MetricCard
        label="CPU"
        icon="&#x2B21;"
        :series="snapshot.cpu"
        :colors="cardColors.cpu"
        :detail="peakDetail(snapshot.cpu)"
        :tooltip="seriesTooltip('CPU', snapshot.cpu, `${snapshot.cpuCores} cores`)"
      />
      <MetricCard
        label="MEM"
        icon="&#x25A3;"
        :series="snapshot.memory"
        :colors="cardColors.memory"
        :detail="peakDetail(snapshot.memory)"
        :tooltip="seriesTooltip('Memory', snapshot.memory, memorySize)"
      />
      <MetricCard
        label="DISK"
        icon="&#x25C8;"
        :series="snapshot.disk"
        :colors="cardColors.disk"
        :detail="peakDetail(snapshot.disk)"
        :tooltip="seriesTooltip('Disk', snapshot.disk, diskSize)"
      />
    </div>
  </div>
</template>
