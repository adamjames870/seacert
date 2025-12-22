import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'
import { fileURLToPath } from 'url'
import { loadEnv } from 'vite'

const __dirname = path.dirname(fileURLToPath(import.meta.url))

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');
  const apiUrl = env.VITE_API_URL || 'http://localhost:8080';

  return {
    plugins: [react()],
    server: {
      proxy: {
        '/api': {
          target: apiUrl,
          changeOrigin: true,
          secure: false,
        },
        '/admin': {
          target: apiUrl,
          changeOrigin: true,
          secure: false,
        },
      },
    },
    resolve: {
      dedupe: ['react', 'react-dom'],
      alias: {
        'react': path.resolve(__dirname, 'node_modules/react'),
        'react-dom': path.resolve(__dirname, 'node_modules/react-dom'),
      },
    },
  }
})
