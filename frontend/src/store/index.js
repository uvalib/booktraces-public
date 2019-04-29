import Vue from 'vue'
import Vuex from 'vuex'
import admin from './modules/admin'
import core from './modules/core'

Vue.use(Vuex)

// Plugin to listen for error messages being set. After a delay, clear them
const errorPlugin = store => {
   store.subscribe((mutation) => {
     if (mutation.type === "setError") {
       if ( mutation.payload != null ) {
         setTimeout( ()=>{ store.commit('setError', null)}, 10000)
       }
     }
   })
 }

export default new Vuex.Store({
  modules: {
    admin: admin,
    core: core,
    plugins: [errorPlugin] 
  },
})