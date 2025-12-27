import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vuetify from 'vite-plugin-vuetify'
import path from 'path'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  const base = env.VITE_BASE_ROUTE || '/'

  return {
    base: base.endsWith('/') ? base : `${base}/`,
    plugins: [
      vue(),
      vuetify({ autoImport: true })
    ],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src')
      },
      extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.vue']
    },
    build: {
      rollupOptions: {
        output: {
          manualChunks(id) {
            if (!id.includes('node_modules')) return
            if (id.includes('vuetify')) return 'vuetify'
            if (id.includes('vue')) return 'vue'
            if (id.includes('codemirror')) return 'codemirror'
            if (id.includes('@xterm')) return 'xterm'
            return 'vendor'
          }
        }
      }
    }
  }
})
