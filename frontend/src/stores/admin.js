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
      error: "",
      pedagogy: []
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

      updatePublicationStatus(id, published) {
         let url = "/api/admin/submissions/" + id + "/publish"
         if (!published) {
            url = "/api/admin/submissions/" + id + "/unpublish"
         }
         this.working = true
         axios.post(url).then(() => {
            let idx = this.submissions.hits.findIndex( sub => sub.id == id)
            if (idx > -1) {
               this.submissions.hits[idx].published = published
            }
            this.working = false
         }).catch((error) => {
            useSystemStore().setError(error.response.data)
            this.working = false
         })
      },

      deleteSubmission( subID ) {
         this.working = true
         axios.delete("/api/admin/submissions/" + subID).then(() => {
            let idx = this.submissions.hits.findIndex( sub => sub.id == subID)
            if (idx > -1) {
               this.submissions.hits.splice(idx, 1)
            }
            this.working = false
         }).catch((error) => {
            this.error = error.response.data
            this.working = false
         })
      },

      deleteEvent(id) {
         axios.delete("/api/admin/events/" + id).then(() => {
            const system = useSystemStore()
            let delIdx = system.events.findIndex( n => n.id == id)
            if (delIdx > -1) {
               system.events.splice(delIdx, 1)
            }
         }).catch((error) => {
            this.error = error.response.data
         })
      },
      async updateEvent(idx, updated) {
         const system = useSystemStore()
         let origEvt = system.events[idx]
         this.working = true
         await axios.put(`/api/admin/events/${origEvt.id}`, updated).then(() => {
            origEvt.date = updated.date
            origEvt.description = updated.description
            this.working = false
            console.log("UPDATED")
         }).catch((error) => {
            this.error = error.response.data
            this.working = false
         })
      },
      async addEvent(newDate, newDesc) {
         const system = useSystemStore()
         let newEvent = {date: newDate, description: newDesc}
         this.working = true
         await axios.post("/api/admin/events/", newEvent).then((response) => {
            system.events.push(response.data)
            system.events.sort( (a,b) => {
               if ( a.date < b.date ){
                  return 1
                }
                if ( a.date > b.date ){
                  return -1
                }
                return 0
            })
            this.working = false
         }).catch((error) => {
            this.error = error.response.data
            this.working = false
         })
      },

      async addNews(newTitle, newContent) {
         const system = useSystemStore()
         this.working = true
         this.error = ""
         let data = { title: newTitle, content: newContent }
         axios.post("/api/admin/news/", data).then((response) => {
            system.news.unshift(response.data)
            this.working = false
         }).catch((error) => {
            this.error = error.response.data
            this.working = false
         })
      },
      async updateNews(id, title, content) {
         const system = useSystemStore()
         let newsRec = system.news.find( n=>n.id = id)
         if ( !newsRec ) {
            this.error = "Unable to find news record "+id
            return
         }
         newsRec.title = title
         newsRec.content = content
         this.working = true
         this.error = ""
         await axios.put("/api/admin/news/" + id, newsRec).catch((error) => {
            this.error = error.response.data
         }).finally( () =>  this.working = false )
      },
      async deleteNews(id) {
         const system = useSystemStore()
         return axios.delete("/api/admin/news/" + id).then(() => {
            let delIdx = system.news.findIndex( n => n.id == id)
            if (delIdx > -1) {
               system.news.splice(delIdx, 1)
            }
         }).catch((error) => {
            this.error = "Unable to delete news: "+error.response.data
         })
      },

      updateSubmission(modified) {
         // FIXME
         // // update wants tags to be a list of IDs, but currently has names...
         // // clone the original list of tag names
         // var tagNames = []
         // if (modified.tags) {
         //    tagNames = modified.tags.slice(0)
         //    for (let i = 0; i < modified.tags.length; i++) {
         //       let tn = modified.tags[i]
         //       this.rootState.tags.some(function (t) {
         //          if (t.name == tn) {
         //             modified.tags[i] = t.id
         //          }
         //          return t.name == tn
         //       })
         //    }
         // }
         // return new Promise((resolve, reject) => {
         //    axios.post("/api/admin/submissions/" + modified.id, modified).then((/*response*/) => {
         //       // put the string tags back in place and update the model in root state
         //       modified.tags = tagNames.slice(0)
         //       this.commit("setSubmissionDetail", modified, { root: true })
         //       resolve()
         //    }).catch((error) => {
         //       useSystemStore().setError(error.response.data)
         //       reject(error)
         //    })
         // })
      },

      getPedagogyDocuments() {
         this.working = true
         this.error = ""
         this.pedagogy = []
         axios.get(`/api/admin/pedagogy`).then((response) => {
            this.pedagogy = response.data
         }).catch((error) => {
            this.error = "Unable to get pedagogy: " + error.response.data
         }).finally( ()=>{
            this.working = false
         })
      },
      deleteDocument( key ) {
         this.working = true
         axios.delete(`/api/admin/pedagogy/${key}`).then(() => {
            let idx = this.pedagogy.findIndex( d => d.key == key)
            if ( idx > -1) {
               this.pedagogy.splice(idx,1)
            }
         }).catch((error) => {
            this.error = error
         }).finally( ()=>{
            this.working = false
         })
      },
      async updateDocument( doc ) {
         this.working = true
         let docIdx = this.pedagogy.findIndex( d => d.id == doc.id)
         await axios.put(`/api/admin/pedagogy/${this.pedagogy[docIdx].key}`, doc).then(() => {
            this.pedagogy[docIdx] = doc
         }).catch((error) => {
            this.error = error
         }).finally( ()=>{
            this.working = false
         })
      },
      async addDocument( doc ) {
         this.working = true
         await axios.post(`/api/admin/pedagogy`, doc).then((response) => {
            this.pedagogy.push( response.data )
         }).catch((error) => {
            this.error = error
         }).finally( ()=>{
            this.working = false
         })
      }
   }
})