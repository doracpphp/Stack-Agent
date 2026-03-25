import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    // Wails embeds the dist directory
    outDir: 'dist',
  },
  server: {
    // Wails dev server proxy
    port: 34115,
    strictPort: true,
  },
})
