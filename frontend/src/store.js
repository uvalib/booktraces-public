import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

// root state object. Holds all of the state for the system
const state = {
  tags: [],
  events: [],
  uploadID: null,
  error: null,
  uploadedFiles: [],
  totalSubmissions: 0,
  submissions: [],
  user: null
}

// state getter functions. All are functions that take state as the first param 
// and the getters themselves as the second param. Getter params are passed 
// as a function. Access as a property like: this.$store.getters.NAME
const getters = {
}

// Synchronous updates to the state. Can be called directly in components like this:
// this.$store.commit('mutation_name') or called from asynchronous actions
const mutations = {
  setTags (state, tags) {
    state.tags = tags
  },
  setEvents (state, events) {
    state.events = events
  },
  setUser (state, user) {
    state.user = user
  },
  clearUser(state) {
    state.user = {firstName: "", lastName:"", title:"", affiliation:"", email:"", phone:""}
  },
  setError (state, error) {
    state.error = error
  },
  setUploadID (state, uploadID) {
    state.uploadID = uploadID
  },
  addSubmissions(state, submissionInfo) {
    state.totalSubmissions = submissionInfo.total
    submissionInfo.thumbs.forEach( function(thumb) {
      state.submissions.push(thumb)
    })
  },
  clearSubmissions(state) {
    state.totalSubmissions = 0
    state.submissions = []
  },
  addUploadedFile (state, filename) {
    state.uploadedFiles.push(filename)
  },
  removeUploadedFile (state, filename) {
    let index = state.uploadedFiles.indexOf(filename)
    if (index !== -1) {
      state.uploadedFiles.splice(index, 1)
    }
  }
}

// Actions are asynchronous calls that commit mutatations to the state.
// All actions get context as a param which is essentially the entirety of the 
// Vuex instance. It has access to all getters, setters and commit. They are 
// called from components like: this.$store.dispatch('action_name', data_object)
const actions = {
  getSubmissions( ctx ) {
    ctx.commit('clearSubmissions' )
    axios.get("/api/submissions").then((response)  =>  {
      ctx.commit('addSubmissions', response.data )
    }).catch((error) => {
      ctx.commit('setError', "Unable to get submissions: "+error.response.data) 
    })
  },
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
  getUploadID( ctx ) {
    axios.get("/api/identifier").then((response)  =>  {
      ctx.commit('setUploadID', response.data )
    }).catch((error) => {
      ctx.commit('setUploadID', []) 
      ctx.commit('setError', "Unable to get uploadID: "+error.response.data) 
    })
  },
  removeUploadedFile( ctx, filename ) {
    ctx.commit("removeUploadedFile",filename)
    axios.delete("/api/upload/"+filename+"?key="+ctx.getters.uploadID)
  }
}

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

// A Vuex instance is created by combining state, getters, actions and mutations
export default new Vuex.Store({
  state,
  getters,
  actions,
  mutations,
  plugins: [errorPlugin] 
})