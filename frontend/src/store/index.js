import Vue from 'vue'
import Vuex from 'vuex'
import auth from './modules/auth'
import unauth from './modules/unauth'
import events from './modules/events'
import news from './modules/news'
import axios from 'axios'
import {version} from '../../package.json'

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
    institutions: [],
    error: null,
    loading: true,
    submissionDetail: null,
    adminMode: false,
    appVersion: version
  },
  getters: {
    hasError: state => {
      return state.error != null
    }
  },
  mutations: {
    setAdminMode (state, mode) {
      state.adminMode = mode
    },
    setTags (state, tags) {
      state.tags = tags
    },
    setInstitutions (state, institutions) {
      state.institutions = institutions
    },
    addInstitution(state, newInst) {
      state.institutions.push(newInst)
    },
    setError (state, error) {
      state.error = error
    },
    setLoading (state, isLoading) {
      state.loading = isLoading
    },
    clearSubmissionDetail (state) {
      state.submissionDetail = null
    },
    setSubmissionDetail (state, details) {
      state.submissionDetail = details
    },
    setCurrSubPublished (state, pub) {
      if (state.submissionDetail) {
        state.submissionDetail.published = pub
      }
    },
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
    getInstitutions( ctx ) {
      axios.get("/api/institutions").then((response)  =>  {
        ctx.commit('setInstitutions', response.data )
      }).catch((error) => {
        ctx.commit('setInstitutions', []) 
        ctx.commit('setError', "Unable to get institutions: "+error.response.data) 
      })
    },
    addInstitution(ctx, name) {
      let url = `/api/institutions`
      return axios.post(url, {name: name}).then((response) => {
         ctx.commit('addInstitution', response.data)
      })
    },
    getSubmissionDetail( ctx, id ) {
      ctx.commit("setLoading", true)
      ctx.commit('clearSubmissionDetail' )
      axios.get("/api/submissions/"+id).then((response)  =>  {
        ctx.commit('setSubmissionDetail', response.data )
        ctx.commit("setLoading", false) 
      }).catch((error) => {
        ctx.commit('setError', "Unable to get submission detail: "+error.response.data) 
        ctx.commit("setLoading", false) 
      })
    },
  },

  modules: {
    public: unauth,
    admin: auth,
    events: events,
    news: news
  },

  plugins: [errorPlugin] 
})