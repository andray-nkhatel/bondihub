import { fileURLToPath, URL } from 'node:url';

import { PrimeVueResolver } from '@primevue/auto-import-resolver';
import vue from '@vitejs/plugin-vue';
import Components from 'unplugin-vue-components/vite';
import { defineConfig } from 'vite';

// https://vitejs.dev/config/
export default defineConfig({
    optimizeDeps: {
        noDiscovery: true,
         // Include jQuery to ensure proper loading
        include: ['jquery']
    },
    plugins: [
        vue(),
        Components({
            resolvers: [PrimeVueResolver()]
        })
    ],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url))
        }
    },
     // Define global variables for external libraries
  define: {
    global: 'globalThis',
  },
  // Development server proxy to bypass CORS issues
  // To use the proxy in development, set VITE_API_BASE_URL=/api in .env.local
  // The proxy will forward /api/* requests to the backend server
  server: {
    proxy: {
      '/api': {
        target: 'https://bluebirdhub.somee.com',
        changeOrigin: true,
        secure: true,
        // Don't rewrite - keep /api in the path
        configure: (proxy, _options) => {
          proxy.on('error', (err, _req, _res) => {
            console.log('Proxy error:', err);
          });
          proxy.on('proxyReq', (proxyReq, req, _res) => {
            console.log('Proxying request:', req.method, req.url, '->', proxyReq.path);
          });
        }
      }
    }
  },
  // Ensure external scripts load properly
  build: {
    rollupOptions: {
      external: ['jquery'],
      output: {
        globals: {
          jquery: 'jQuery'
        }
      }
    }
  }
});
