import { createApp } from 'vue'
import App from '@/App.vue'
import router from '@/router'
import store from '@/store'

const app = createApp(App)

// provide store access to the rouer
store.router = router

// bind store and router to all componens as $store and $router
app.use(store)
app.use(router)

import BTSpinner from "@/components/BTSpinner.vue"
app.component('BTSpinner', BTSpinner)

import PinchZoom from 'vue-pinch-zoom'
app.component('pinch-zoom', PinchZoom)

import VueCookies from 'vue-cookies'
app.use(VueCookies)

import '@fortawesome/fontawesome-free/css/all.css'

app.mount('#app')