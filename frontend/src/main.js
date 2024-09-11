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

// Primevue setup
import PrimeVue from 'primevue/config'
import BookTraces from '@/assets/theme/booktraces'
import 'primeicons/primeicons.css'
import Button from 'primevue/button'
import ConfirmationService from 'primevue/confirmationservice'
import ConfirmDialog from 'primevue/confirmdialog'

app.use(ConfirmationService)
app.component("Button", Button)
app.component("ConfirmDialog", ConfirmDialog)

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