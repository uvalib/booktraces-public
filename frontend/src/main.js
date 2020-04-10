import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

import PinchZoom from 'vue-pinch-zoom'
Vue.component('pinch-zoom', PinchZoom)

import VueCookies from 'vue-cookies'
Vue.use(VueCookies)

Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')