import axios from 'axios'

const core = {
  namespaced: true,
  
  // root state object. Holds all of the state for the system
  state: {
    tags: [],
    events: [],
    uploadID: null,
    error: null,
    uploadedFiles: [],
    totalSubmissions: 0,
    submissions: [],
    details: null,
    loading: true
  },

  // state getter functions. All are functions that take state as the first param 
  // and the getters themselves as the second param. Getter params are passed 
  // as a function. Access as a property like: this.$store.getters.NAME
  getters: {
  },

  // Synchronous updates to the state. Can be called directly in components like this:
  // this.$store.commit('mutation_name') or called from asynchronous actions
  mutations: {
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
      state.user = null
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
    clearUploadedFiles(state) {
      state.uploadedFiles = []
    },
    addUploadedFile (state, filename) {
      state.uploadedFiles.push(filename)
    },
    removeUploadedFile (state, filename) {
      let index = state.uploadedFiles.indexOf(filename)
      if (index !== -1) {
        state.uploadedFiles.splice(index, 1)
      }
    },
    clearSubmissionDetail (state) {
      state.details = null
    },
    setSubmissionDetail (state, details) {
      state.details = details
    },
    setloading (state, isLoading) {
      state.loading = isLoading
    }
  },

  // Actions are asynchronous calls that commit mutatations to the state.
  // All actions get context as a param which is essentially the entirety of the 
  // Vuex instance. It has access to all getters, setters and commit. They are 
  // called from components like: this.$store.dispatch('action_name', data_object)
  actions: {
    getRecentThumbs( ctx ) {
      ctx.commit('clearSubmissions' )
      axios.get("/api/recents").then((response)  =>  {
        ctx.commit('addSubmissions', response.data )
      }).catch((error) => {
        ctx.commit('setError', "Unable to get recent submissions: "+error.response.data) 
      })
    },
    getSubmissionDetail( ctx, id ) {
      ctx.commit("setloading", true)
      ctx.commit('clearSubmissionDetail' )
      axios.get("/api/submissions/"+id).then((response)  =>  {
        ctx.commit('setSubmissionDetail', response.data )
        ctx.commit("setloading", false)
      }).catch((error) => {
        ctx.commit('setError', "Unable to get submission detail: "+error.response.data) 
        ctx.commit("setloading", false)
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
}

export default core