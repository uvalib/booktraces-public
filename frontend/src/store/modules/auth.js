import axios from 'axios'
import router from '../../router'

const auth = {
  namespaced: true,
  state: {
    totalSubmissions: 0,
    filteredTotal: 0,
    pageSize: 0,
    page: 1,
    submissions: [],
    user: null,
    queryStr: "",
    tgtTag: "",
    news: []
  },

  // Properties that are computed based on state
  getters: {
    isAuthenticated(state) {
      if (state.user == null) {
         return false
      }
      return state.user.authenticated
    },
    loginName: state => {
      if (state.user) {
        return state.user.firstName + " ("+state.user.email+")"
      }
      return ""
    }
  },

  // Synchronous updates to the state
  mutations: {
    setTagFilter(state, val) {
      state.tgtTag = val
    },
    updateSearchQuery(state, val) {
      state.queryStr = val
    },
    gotoFirstPage(state) {
      state.page = 1
    },
    gotoLastPage(state) {
      state.page = Math.floor( state.totalSubmissions / state.pageSize)+1
    },
    nextPage(state) {
      state.page++
    },
    prevPage(state) {
      state.page--
    },
    setUser (state, user) {
      state.user = user
    },
    clearUser(state) {
      state.user = {firstName: "", lastName:"", title:"", affiliation:"", email:"", phone:""}
    },
    setSubbmissionPage(state, resp) {
      state.filteredTotal = resp.filteredTotal
      state.totalSubmissions = resp.total
      state.pageSize = resp.pageSize
      state.page = resp.page
      state.submissions = resp.submissions
    },
    setPublished(state, payload) {
      if (!payload) {
        return
      }
      let subID = payload.id
      state.submissions.some( function(sub) {
        if (sub.id == subID) {
          sub.published = payload.public
        }
        return sub.id == subID
      })
    },
    deleteSubmission(state, subID) {
      state.submissions.some( function(sub,idx) {
        if (sub.id == subID) {
          state.submissions.splice(idx, 1)
          return true
        }
        return false
      })
    },
  },

  // Asynchronous updates to the state
  actions: {
    firstPage( ctx ) {
      ctx.commit('gotoFirstPage')
      ctx.dispatch("getSubmissions")
    },
    prevPage( ctx ) {
      ctx.commit('prevPage')
      ctx.dispatch("getSubmissions")
    },
    nextPage( ctx ) {
      ctx.commit('nextPage')
      ctx.dispatch("getSubmissions")
    },
    lastPage( ctx ) {
      ctx.commit('gotoLastPage')
      ctx.dispatch("getSubmissions")
    },

    rotateImage(ctx, data ) {
      ctx.commit("setLoading", true, {root: true})
      let url = `/api/admin/submissions/${data.submissionID}/rotate?url=${data.imgURL}`
      axios.put(url,{withCredentials: true }).then((_response)  =>  {
        ctx.commit("setLoading", false, {root: true})
      }).catch((error) => {
        ctx.commit("setLoading", false, {root: true})
        ctx.commit('setError', "Unable to get rotate image: "+error.response.data, {root: true}) 
      })
    },

    getSubmissions( ctx ) {
      ctx.commit("setLoading", true, {root: true})
      let url = "/api/admin/submissions?page="+ctx.state.page
      if (ctx.state.queryStr.length > 0 ) {
        url = url +"&q="+ctx.state.queryStr
      }
      if (ctx.state.tgtTag.length > 0 ) {
        url = url +"&t="+ctx.state.tgtTag
      }
      axios.get(url,{ withCredentials: true }).then((response)  =>  {
        ctx.commit('setSubbmissionPage', response.data )
        ctx.commit("setLoading", false, {root: true})
      }).catch((error) => {
        ctx.commit("setLoading", false, {root: true})
        ctx.commit('setError', "Unable to get recent submissions: "+error.response.data, {root: true}) 
        if (error.response.status == 403) {
          router.push("/forbidden")
        }
      })
    },
    updatePublicationStatus( ctx, payload ) {
      let id = payload.id 
      let published = payload.public
      let url = "/api/admin/submissions/"+id+"/publish"
      if (!published) {
        url = "/api/admin/submissions/"+id+"/unpublish"
      }
      ctx.commit("setLoading", true, {root: true})
      axios.post(url).then((/*response*/)  =>  {
        ctx.commit("setPublished", payload )
        ctx.commit("setCurrSubPublished", published, {root: true})
        ctx.commit("setLoading", false, {root: true})
      }).catch((error) => {
        ctx.commit("setError",error.response.data, {root: true}) 
        ctx.commit("setLoading", false, {root: true})
      })
    },
    deleteSubmission( ctx, payload ) {
      let id=payload.id
      ctx.commit("setLoading", true, {root: true})
      axios.delete("/api/admin/submissions/"+id).then((/*response*/)  =>  {
        ctx.commit("deleteSubmission", id)
        ctx.commit("setLoading", false, {root: true})
        if (payload.backToIndex) {
          router.push("/admin")
        }
      }).catch((error) => {
        ctx.commit("setError",error.response.data, {root: true}) 
        ctx.commit("setLoading", false, {root: true})
      })
    },
    updateSubmission(ctx, modified) {
      // update wants tags to be a list of IDs, but currently has names...
      // clone the original list of tag names
      var tagNames = []
      if (modified.tags) {
        tagNames = modified.tags.slice(0)
        for ( let i=0; i<modified.tags.length; i++) {
          let tn = modified.tags[i]
          ctx.rootState.tags.some( function(t) {
            if (t.name == tn) {
              modified.tags[i] = t.id
            }
            return t.name == tn
          })
        }
      }
      return new Promise((resolve, reject) => {
        axios.post("/api/admin/submissions/"+modified.id, modified).then((/*response*/)  =>  {
          // put the string tags back in place and update the model in root state
          modified.tags = tagNames.slice(0)
          ctx.commit("setSubmissionDetail", modified, {root: true}) 
          resolve()
        }).catch((error) => {
          ctx.commit("setError", error.response.data, {root: true}) 
          reject(error)
        })  
      })
    },
    getNews( ctx ) {
      axios.get("/api/admin/news").then((response)  =>  {
        ctx.commit('setNews', response.data )
      }).catch((error) => {
        ctx.commit('setNews', []) 
        ctx.commit('setError', "Unable to get news: "+error.response.data) 
      })
    },
  }
}

export default auth