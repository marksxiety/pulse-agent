<script setup lang="ts">
import type { Theme } from '../types'

defineProps<{
  themes: Theme[]
  active: string
}>()

const emit = defineEmits<{
  (event: 'select', name: string): void
}>()
</script>

<template>
  <div class="theme-panel">
    <button
      v-for="theme in themes"
      :key="theme.name"
      type="button"
      class="theme-row"
      :class="{ active: theme.name === active }"
      @click="emit('select', theme.name)"
    >
      <span class="swatches">
        <span class="swatch" :style="{ background: theme.cpu_accent }"></span>
        <span class="swatch" :style="{ background: theme.mem_accent }"></span>
        <span class="swatch" :style="{ background: theme.disk_accent }"></span>
      </span>
      <span class="theme-name">{{ theme.name }}</span>
    </button>
  </div>
</template>
