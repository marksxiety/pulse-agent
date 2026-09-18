import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  build: {
    // main.go embeds frontend/dist, and the tracked `gitkeep` placeholder in
    // that directory is what lets `go build ./...` succeed on a fresh clone.
    // Emptying the output directory would delete that placeholder and leave a
    // phantom deletion in git status after every frontend build.
    emptyOutDir: false,
  },
})
