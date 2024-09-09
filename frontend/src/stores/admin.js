import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useAdminStore = defineStore('admin', {
   state: () => ({
      working: false,
      submissions: {
         total: 0,
         start: 0,
         hits: [],
         query: "",
         tagFilter: "",
      },
      user: null,

      news: []
   }),
   getters: {
      isAuthenticated: state => {
         if (state.user == null) {
            return false
         }
         return state.user.authenticated
      },
      loginName: state => {
         if (state.user) {
            return state.user.firstName + " (" + state.user.email + ")"
         }
         return ""
      },
   },
   actions: {
      clearUser() {
         this.user = { firstName: "", lastName: "", title: "", affiliation: "", email: "", phone: "" }
      },

      rotateImage(data) {
         this.working = true
         let url = `/api/admin/submissions/${data.submissionID}/rotate?url=${data.imgURL}`
         axios.put(url, { withCredentials: true }).then((_response) => {
            this.working = false
         }).catch((error) => {
            this.working = false
            useSystemStore().setError("Unable to get rotate image: " + error.response.data)
         })
      },

      getSubmissions() {
         this.working = true
         let url = "/api/admin/submissions?start=" + this.submissions.start
         if (this.submissions.query.length > 0) {
            url = url + "&q=" + this.submissions.query
         }
         if (this.submissions.tagFilter.length > 0) {
            url = url + "&t=" + this.submissions.tagFilter
         }
         axios.get(url, { withCredentials: true }).then((response) => {
            this.submissions.total = response.data.total
            this.submissions.start = response.data.start
            this.submissions.hits = response.data.submissions
            this.working = false
            console.log("GOT ADMIN SUBMISSIONS")
         }).catch((error) => {
            this.working = false
            useSystemStore().setError("Unable to get recent submissions: " + error.response.data)
            if (error.response.status == 403) {
               router.push("/forbidden")
            }
         })
      },

      updatePublicationStatus(payload) {
         let id = payload.id
         let published = payload.public
         let url = "/api/admin/submissions/" + id + "/publish"
         if (!published) {
            url = "/api/admin/submissions/" + id + "/unpublish"
         }
         this.working = true
         axios.post(url).then(() => {
            this.commit("setPublished", payload)
            let idx = this.submissions.findIndex( sub => sub.id == id)
            if (idx > -1) {
               this.submissions[idx].public = published
            }
            this.working = false
         }).catch((error) => {
            useSystemStore.setError(error.response.data)
            this.working = false
         })
      },

      deleteSubmission(payload) {
         let tgtID = payload.id
         this.working = true
         axios.delete("/api/admin/submissions/" + tgtID).then(() => {
            let idx = this.submissions.findIndex( sub => sub.id == tgtID)
            if (idx > -1) {
               this.submissions.splice(idx, 1)
            }
            this.working = false
            if (payload.backToIndex) {
               router.push("/admin")
            }
         }).catch((error) => {
            useSystemStore.setError(error.response.data)
            this.working = false
         })
      },

      updateSubmission(modified) {
         // update wants tags to be a list of IDs, but currently has names...
         // clone the original list of tag names
         var tagNames = []
         if (modified.tags) {
            tagNames = modified.tags.slice(0)
            for (let i = 0; i < modified.tags.length; i++) {
               let tn = modified.tags[i]
               this.rootState.tags.some(function (t) {
                  if (t.name == tn) {
                     modified.tags[i] = t.id
                  }
                  return t.name == tn
               })
            }
         }
         return new Promise((resolve, reject) => {
            axios.post("/api/admin/submissions/" + modified.id, modified).then((/*response*/) => {
               // put the string tags back in place and update the model in root state
               modified.tags = tagNames.slice(0)
               this.commit("setSubmissionDetail", modified, { root: true })
               resolve()
            }).catch((error) => {
               useSystemStore.setError(error.response.data)
               reject(error)
            })
         })
      },

      // getNews() {
      //    axios.get("/api/admin/news").then((response) => {
      //       this.commit('setNews', response.data)
      //    }).catch((error) => {
      //       this.commit('setNews', [])
      //       this.commit('setError', "Unable to get news: " + error.response.data)
      //    })
      // },
   }
})