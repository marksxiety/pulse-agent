/**
 * Formats a byte count the same way the terminal UI does: decimal (1e9) units
 * with two decimals, so both frontends report identical numbers.
 */
export function bytesToGB(bytes: number): string {
  return `${(bytes / 1e9).toFixed(2)} GB`
}

/**
 * Formats an elapsed number of seconds to match `utils.FormatUptime` exactly —
 * seconds are always shown, and minutes are zero-padded only once an hour is in
 * play. Both frontends report the same uptime string for the same duration.
 */
export function formatUptime(totalSeconds: number): string {
  const seconds = Math.max(0, Math.floor(totalSeconds))
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor(seconds / 60) % 60
  const secs = seconds % 60

  const tail = `${String(secs).padStart(2, '0')}s`
  if (hours > 0) {
    return `${hours}h ${String(minutes).padStart(2, '0')}m ${tail}`
  }
  return `${minutes}m ${tail}`
}

/** Expands #rrggbb into an rgba() string with the given alpha. */
export function withAlpha(hex: string, alpha: number): string {
  const value = hex.replace('#', '')
  if (value.length !== 6) {
    return hex
  }
  const r = parseInt(value.slice(0, 2), 16)
  const g = parseInt(value.slice(2, 4), 16)
  const b = parseInt(value.slice(4, 6), 16)
  return `rgba(${r}, ${g}, ${b}, ${alpha})`
}
