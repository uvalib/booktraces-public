import axios from 'axios'
import router from '../../router'

const auth = {
  namespaced: true,
  state: {
    totalSubmissions: 0,
    pageSize: 0,
    page: 1,
    submissions: [],
    user: null,
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
    setSubbmissionPage(state, submissionInfo) {
      state.totalSubmissions = submissionInfo.total
      state.pageSize = submissionInfo.pageSize
      state.page = submissionInfo.page
      state.submissions = submissionInfo.submissions
    },
    setPublished(state, subID) {
      state.submissions.some( function(sub) {
        if (sub.id == subID) {
          sub.published = 1
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
    getSubmissions( ctx ) {
      axios.get("/api/admin/submissions?page="+ctx.state.page,{ withCredentials: true }).then((response)  =>  {
        ctx.commit('setSubbmissionPage', response.data )
      }).catch((error) => {
        ctx.commit('setError', "Unable to get recent submissions: "+error.response.data, {root: true}) 
        if (error.response.status == 403) {
          router.push("/forbidden")
        }
      })
    },
    publishSubmission( ctx, id ) {
      ctx.commit("setLoading", true, {root: true})
      axios.post("/api/admin/submissions/"+id+"/publish", {userID:ctx.state.user.id }).then((/*response*/)  =>  {
        ctx.commit("setPublished", id)
        ctx.commit("setCurrSubPublished", null, {root: true})
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
    }
  }
}

export default auth