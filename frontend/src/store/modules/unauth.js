import axios from 'axios'
import { getField, updateField } from 'vuex-map-fields'

const unauth = {
  namespaced: true,
  
  // root state object. Holds all of the state for the system
  state: {
    uploadID: null,
    uploadedFiles: [],
    totalSubmissions: 0,
    thumbs: [],
    archives: [],
    recents: [],
    currPage: 0,
    showSearch: false,
    searchResults: [],
    query: "",
    archiveDate: "",
    targetInstitution: "",
    tgtTag: "",
    news: []
  },

  // state getter functions. All are functions that take state as the first param 
  // and the getters themselves as the second param. Getter params are passed 
  // as a function. Access as a property like: this.$store.getters.NAME
  getters: {
    getField,
    thumbsCount: state => {
      return state.thumbs.length
    },
    searchHitCount: state => {
      return state.searchResults.length
    },
    searchType: state => {
      if (state.tgtTag.length > 0) {
        return "tag"
      }
      if (state.archiveDate.length > 0) {
        return "archive"
      }
      if ( state.targetInstitution != "" ) {
         return "institution"   
      }
      if (state.query.length > 0) {
        return "query"
      }
      return "none"
    },
  },

  // Synchronous updates to the state. Can be called directly in components like this:
  // this.$store.commit('mutation_name') or called from asynchronous actions
  mutations: {
    updateField,
    setArchiveDate (state, date) {
      state.archiveDate = date
      state.searchResults =[]
      state.query = ""
      state.tgtTag = ""
      state.targetInstitution = ""
    },
    setTargetTag (state, t) {
      state.searchResults =[]
      state.query = ""
      state.archiveDate = ""
      state.targetInstitution = ""
      state.tgtTag = t
    },
    initInstitutionSearch (state) {
      state.searchResults =[]
      state.query = ""
      state.archiveDate = ""
      state.tgtTag = ""
    },
    updateSearchQuery(state, q) {
      state.query = q
      state.searchResults = []
      state.archiveDate = ""
      state.tgtTag = ""
      state.targetInstitution = ""
    },
    nextThumbsPage(state) {
      state.currPage++
    },
    showSearch(state, show) {
      state.showSearch = show
    },
    setUploadID (state, uploadID) {
      state.uploadID = uploadID
    },
    addThumbs(state, submissionInfo) {
      state.totalSubmissions = submissionInfo.total
      submissionInfo.thumbs.forEach( function(thumb) {
        state.thumbs.push(thumb)
      })
    },
    clearThumbs(state) {
      state.currPage = 0
      state.totalSubmissions = 0
      state.thumbs = []
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
    setArchives (state, archives) {
      state.archives = archives
    },
    setRecents (state, recents) {
      state.recents = recents
    },
    setSearchResults(state, results) {
      state.searchResults = results
    },
    setNews (state, news) {
      state.news = news
    },
  },

  // Actions are asynchronous calls that commit mutatations to the state.
  // All actions get context as a param which is essentially the entirety of the 
  // Vuex instance. It has access to all getters, setters and commit. They are 
  // called from components like: this.$store.dispatch('action_name', data_object)
  actions: {
    search( ctx ) {
      ctx.commit('setLoading', true, {root: true}) 
      ctx.commit('showSearch', false) 
      axios.get("/api/search?q="+ctx.state.query).then((response)  =>  {
        ctx.commit('setSearchResults', response.data )
        ctx.commit('setLoading', false, {root: true}) 
      }).catch((error) => {
        ctx.commit('setError', "Unable to get search results: "+error.response.data, {root: true}) 
        ctx.commit('setSearchResults', [] )
        ctx.commit('setLoading', false, {root: true}) 
      })
    },
    searchInstitutions( ctx) {
      ctx.commit('setLoading', true, {root: true}) 
      ctx.commit('showSearch', false) 
      ctx.commit('initInstitutionSearch')
      axios.get("/api/search?i="+ctx.state.targetInstitution).then((response)  =>  {
        ctx.commit('setSearchResults', response.data )
        ctx.commit('setLoading', false, {root: true}) 
      }).catch((error) => {
        ctx.commit('setError', "Unable to get search results: "+error.response.data, {root: true}) 
        ctx.commit('setSearchResults', [] )
        ctx.commit('setLoading', false, {root: true}) 
      })
    },
    getArchive( ctx, date ) {
      ctx.commit('setLoading', true, {root: true}) 
      ctx.commit('showSearch', false) 
      ctx.commit('setArchiveDate', date)
      axios.get("/api/search?a="+date).then((response)  =>  {
        ctx.commit('setSearchResults', response.data )
        ctx.commit('setLoading', false, {root: true}) 
      }).catch((error) => {
        ctx.commit('setError', "Unable to get archives: "+error.response.data, {root: true}) 
        ctx.commit('setSearchResults', [] )
        ctx.commit('setLoading', false, {root: true}) 
      })
    },
    getTaggedSubmissions( ctx, tag ) {
      ctx.commit('setLoading', true, {root: true}) 
      ctx.commit('showSearch', false) 
      ctx.commit('setTargetTag', tag)
      axios.get("/api/search?t="+tag).then((response)  =>  {
        ctx.commit('setSearchResults', response.data )
        ctx.commit('setLoading', false, {root: true}) 
      }).catch((error) => {
        ctx.commit('setError', "Unable to get archives: "+error.response.data, {root: true}) 
        ctx.commit('setSearchResults', [] )
        ctx.commit('setLoading', false, {root: true}) 
      })
    },
    getArchiveDates( ctx ) {
      if (ctx.state.archives.length > 0) {
        return
      }
      axios.get("/api/archives").then((response)  =>  {
        ctx.commit('setArchives', response.data )
      }).catch((error) => {
        ctx.commit('setError', "Unable to get archives: "+error.response.data, {root: true}) 
      })
    },
    getRecentSubmissions( ctx ) {
      if (ctx.state.recents.length > 0) {
        return
      }
      axios.get("/api/recents").then((response)  =>  {
        ctx.commit('setRecents', response.data )
      }).catch((error) => {
        ctx.commit('setError', "Unable to get recent submissions: "+error.response.data, {root: true}) 
      })
    },
    getRecentThumbs( ctx ) {
      ctx.commit('setLoading', true, {root: true}) 
      ctx.commit('nextThumbsPage')
      axios.get("/api/thumbs?page="+ctx.state.currPage).then((response)  =>  {
        ctx.commit('addThumbs', response.data )
        ctx.commit('setLoading', false, {root: true}) 
      }).catch((error) => {
        ctx.commit('setError', "Unable to get recent thumbs: "+error.response.data, {root: true}) 
        ctx.commit('setLoading', false, {root: true}) 
      })
    },
    getUploadID( ctx ) {
      axios.get("/api/identifier").then((response)  =>  {
        ctx.commit('setUploadID', response.data )
      }).catch((error) => {
        ctx.commit('setUploadID', []) 
        ctx.commit('setError', "Unable to get uploadID: "+error.response.data, {root: true}) 
      })
    },
    removeUploadedFile( ctx, filename ) {
      ctx.commit("removeUploadedFile",filename)
      axios.delete("/api/upload/"+filename+"?key="+ctx.getters.uploadID) 
    }
  }
}

export default unauth
