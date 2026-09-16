/** A colour palette, mirrored from the Go `theme.Theme` struct. */
export interface Theme {
  name: string
  bg: string
  surface: string
  border: string
  muted: string
  subtle: string
  text: string
  dim: string
  cpu_accent: string
  mem_accent: string
  disk_accent: string
  warn: string
  danger: string
}

/** The colours a single metric card needs to render itself. */
export interface CardColors {
  accent: string
  warn: string
  danger: string
}

/** One metric's current reading plus its downsampled 3-hour trend. */
export interface Series {
  percent: number
  trend: number[] | null
  max: number
  avg: number
}

/** Everything the backend pushes once per tick. */
export interface Snapshot {
  cpu: Series
  memory: Series
  disk: Series
  cpuCores: number
  memoryUsed: number
  memoryTotal: number
  diskUsed: number
  diskTotal: number
  uptimeSecs: number
}

/** One-off state fetched before the first snapshot arrives. */
export interface Bootstrap {
  theme: string
  themes: Theme[]
  version: string
  alwaysOnTop: boolean
}

export const EVENT_SNAPSHOT = 'pulse:snapshot'
