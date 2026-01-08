import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    host: true,
    proxy: {
      // Config Manager API (Go) - runs on port 8090
      // More specific path first to avoid conflicts
      '/api/v1/devices': {
        target: 'http://localhost:8090',
        changeOrigin: true,
      },
      // Engine API (Rust) - runs on port 8081
      '/api/v1': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
      // Alarm API (Go) - runs on port 8083
      '/alarms': {
        target: 'http://localhost:8083',
        changeOrigin: true,
      },
      // Audit API (Go) - runs on port 8082
      '/logs': {
        target: 'http://localhost:8082',
        changeOrigin: true,
      },
      // Auth API (Go) - runs on port 8080
      '/auth': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
