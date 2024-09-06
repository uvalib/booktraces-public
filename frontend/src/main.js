import { createApp, markRaw } from 'vue'
import { createPinia } from 'pinia'
import App from '@/App.vue'
import router from '@/router'

const app = createApp(App)

const pinia = createPinia()
pinia.use(({ store }) => {
   // all stores can access router with this.router
   store.router = markRaw(router)
})

// bind store and router to all componens as $store and $router
app.use(pinia)
app.use(router)

import BTSpinner from "@/components/BTSpinner.vue"
app.component('BTSpinner', BTSpinner)

import VueImageZoomer from 'vue-image-zoomer'
import 'vue-image-zoomer/dist/style.css'
app.use(VueImageZoomer)

import VueCookies from 'vue-cookies'
app.use(VueCookies)

// Styles
import '@fortawesome/fontawesome-free/css/all.css'

// Primevue setup
import PrimeVue from 'primevue/config'
import BookTraces from '@/assets/theme/booktraces'
import 'primeicons/primeicons.css'
import Button from 'primevue/button'

app.component("Button", Button)

app.use(PrimeVue, {
   theme: {
      preset: BookTraces,
      options: {
         prefix: 'p',
         darkModeSelector: '.bt-dark'
      }
   }
})

app.mount('#app')