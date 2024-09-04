import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useAuthStore = defineStore('auth', {
   state: () => ({
      totalSubmissions: 0,
      filteredTotal: 0,
      pageSize: 0,
      page: 1,
      submissions: [],
      user: null,
      queryStr: "",
      tgtTag: "",
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
      }
   },
   actions: {
      setTagFilter(val) {
         this.tgtTag = val
      },
      updateSearchQuery(val) {
         this.queryStr = val
      },
      gotoFirstPage() {
         this.page = 1
      },
      gotoLastPage() {
         this.page = Math.floor(this.totalSubmissions / this.pageSize) + 1
      },
      nextPage() {
         this.page++
      },
      prevPage() {
         this.page--
      },
      setUser(user) {
         this.user = user
      },
      clearUser() {
         this.user = { firstName: "", lastName: "", title: "", affiliation: "", email: "", phone: "" }
      },
      setSubbmissionPage(resp) {
         this.filteredTotal = resp.filteredTotal
         this.totalSubmissions = resp.total
         this.pageSize = resp.pageSize
         this.page = resp.page
         this.submissions = resp.submissions
      },
      setPublished(payload) {
         if (!payload) {
            return
         }
         let subID = payload.id
         this.submissions.some(function (sub) {
            if (sub.id == subID) {
               sub.published = payload.public
            }
            return sub.id == subID
         })
      },
      deleteSubmission(subID) {
         this.submissions.some(function (sub, idx) {
            if (sub.id == subID) {
               this.submissions.splice(idx, 1)
               return true
            }
            return false
         })
      },

      firstPage() {
         this.commit('gotoFirstPage')
         this.dispatch("getSubmissions")
      },
      prevPage() {
         this.commit('prevPage')
         this.dispatch("getSubmissions")
      },
      nextPage() {
         this.commit('nextPage')
         this.dispatch("getSubmissions")
      },
      lastPage() {
         this.commit('gotoLastPage')
         this.dispatch("getSubmissions")
      },

      rotateImage(data) {
         const system = useSystemStore()
         system.loading = true
         let url = `/api/admin/submissions/${data.submissionID}/rotate?url=${data.imgURL}`
         axios.put(url, { withCredentials: true }).then((_response) => {
            system.loading = false
         }).catch((error) => {
            system.loading = false
            this.commit('setError', "Unable to get rotate image: " + error.response.data, { root: true })
         })
      },

      getSubmissions() {
         const system = useSystemStore()
         system.loading = true
         let url = "/api/admin/submissions?page=" + this.state.page
         if (this.state.queryStr.length > 0) {
            url = url + "&q=" + this.state.queryStr
         }
         if (this.state.tgtTag.length > 0) {
            url = url + "&t=" + this.state.tgtTag
         }
         axios.get(url, { withCredentials: true }).then((response) => {
            this.commit('setSubbmissionPage', response.data)
            system.loading = false
         }).catch((error) => {
            system.loading = false
            this.commit('setError', "Unable to get recent submissions: " + error.response.data, { root: true })
            if (error.response.status == 403) {
               router.push("/forbidden")
            }
         })
      },
      updatePublicationStatus(payload) {
         const system = useSystemStore()
         let id = payload.id
         let published = payload.public
         let url = "/api/admin/submissions/" + id + "/publish"
         if (!published) {
            url = "/api/admin/submissions/" + id + "/unpublish"
         }
         system.loading = true
         axios.post(url).then((/*response*/) => {
            this.commit("setPublished", payload)
            this.commit("setCurrSubPublished", published, { root: true })
            system.loading = false
         }).catch((error) => {
            this.commit("setError", error.response.data, { root: true })
            system.loading = false
         })
      },
      deleteSubmission(payload) {
         const system = useSystemStore()
         let id = payload.id
         system.loading = true
         axios.delete("/api/admin/submissions/" + id).then((/*response*/) => {
            this.commit("deleteSubmission", id)
            system.loading = false
            if (payload.backToIndex) {
               router.push("/admin")
            }
         }).catch((error) => {
            this.commit("setError", error.response.data, { root: true })
            system.loading = false
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
               this.commit("setError", error.response.data, { root: true })
               reject(error)
            })
         })
      },
      getNews() {
         axios.get("/api/admin/news").then((response) => {
            this.commit('setNews', response.data)
         }).catch((error) => {
            this.commit('setNews', [])
            this.commit('setError', "Unable to get news: " + error.response.data)
         })
      },
   }
})