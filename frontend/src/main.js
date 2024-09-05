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

import '@fortawesome/fontawesome-free/css/all.css'

app.mount('#app')