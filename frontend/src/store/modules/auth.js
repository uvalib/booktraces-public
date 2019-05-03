import axios from 'axios'
import router from '../../router'

const auth = {
   namespaced: true,
  // root state object. Holds all of the state for the system
  state: {
    totalSubmissions: 0,
    submissions: [],
    user: null,
  },

  // state getter functions. All are functions that take state as the first param 
  // and the getters themselves as the second param. Getter params are passed 
  // as a function. Access as a property like: this.$store.getters.NAME
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

  // Synchronous updates to the state. Can be called directly in components like this:
  // this.$store.commit('mutation_name') or called from asynchronous actions
  mutations: {
    setUser (state, user) {
      state.user = user
    },
    clearUser(state) {
      state.user = {firstName: "", lastName:"", title:"", affiliation:"", email:"", phone:""}
    },
    addSubmissions(state, submissionInfo) {
      state.totalSubmissions = submissionInfo.total
      submissionInfo.submissions.forEach( function(sub) {
        state.submissions.push(sub)
      })
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
    clearSubmissions(state) {
      state.totalSubmissions = 0
      state.submissions = []
    },
  },

  // Actions are asynchronous calls that commit mutatations to the state.
  // All actions get context as a param which is essentially the entirety of the 
  // Vuex instance. It has access to all getters, setters and commit. They are 
  // called from components like: this.$store.dispatch('action_name', data_object)
  actions: {
    getSubmissions( ctx ) {
      ctx.commit('clearSubmissions' )
      axios.get("/api/admin/submissions",{ withCredentials: true }).then((response)  =>  {
        ctx.commit('addSubmissions', response.data )
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
        ctx.commit("setLoading", false, {root: true})
      }).catch((error) => {
        ctx.commit("setError",error.response.data, {root: true}) 
        ctx.commit("setLoading", false, {root: true})
      })
    },
    deleteSubmission( ctx, id ) {
      ctx.commit("setLoading", true, {root: true})
      axios.delete("/api/admin/submissions/"+id).then((/*response*/)  =>  {
        ctx.commit("deleteSubmission", id)
        ctx.commit("setLoading", false, {root: true})
      }).catch((error) => {
        ctx.commit("setError",error.response.data, {root: true}) 
        ctx.commit("setLoading", false, {root: true})
      })
    }
  }
}

export default auth