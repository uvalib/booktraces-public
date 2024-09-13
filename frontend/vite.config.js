 /*global process */

 import { fileURLToPath, URL } from 'url'
 import { defineConfig } from 'vite'
 import vue from '@vitejs/plugin-vue'

 export default defineConfig({
   define: {
      // enable hydration mismatch details in production build
      __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: 'true'
    },
    plugins: [vue()],
    resolve: {
       alias: {
          '@': fileURLToPath(new URL('./src', import.meta.url))
       }
    },
    server: { // this is used in dev mode only
       port: 8080,
       proxy: {
         '/api': {
           target: process.env.BOOKTRACES_SVC, // export BOOKTRACES_SVC=http://localhost:8085
           changeOrigin: true,
           logLevel: 'debug'
         },
         '/authenticate': {
           target: process.env.BOOKTRACES_SVC, // or 'http://localhost:8095',
           changeOrigin: true,
           logLevel: 'debug'
         },
         '/uploads': {
           target: process.env.BOOKTRACES_SVC, // or 'http://localhost:8095',
           changeOrigin: true,
           logLevel: 'debug'
         }
       }
    },
 })


