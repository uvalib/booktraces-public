import Vue from 'vue'
import Vuex from 'vuex'
import auth from './modules/auth'
import unauth from './modules/unauth'
import axios from 'axios'

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
  // Global state/mutations/actions: for stuff that is used across modules
  state: {
    tags: [],
    events: [],
    error: null,
    loading: true
  },
  getters: {
    hasError: state => {
      return state.error != null
    }
  },
  mutations: {
    setTags (state, tags) {
      state.tags = tags
    },
    setEvents (state, events) {
      state.events = events
    },
    setError (state, error) {
      state.error = error
    },
    setLoading (state, isLoading) {
      state.loading = isLoading
    }
  },
  actions: {
    getTags( ctx ) {
      axios.get("/api/tags").then((response)  =>  {
        ctx.commit('setTags', response.data )
      }).catch((error) => {
        ctx.commit('setTags', []) 
        ctx.commit('setError', "Unable to get tags: "+error.response.data) 
      })
    },
    getEvents( ctx ) {
      axios.get("/api/events").then((response)  =>  {
        ctx.commit('setEvents', response.data )
      }).catch((error) => {
        ctx.commit('setEvents', []) 
        ctx.commit('setError', "Unable to get events: "+error.response.data) 
      })
    },
  },

  modules: {
    public: unauth,
    admin: auth,
  },

  plugins: [errorPlugin] 
})