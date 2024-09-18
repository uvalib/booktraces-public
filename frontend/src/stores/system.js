import { defineStore } from 'pinia'
import axios from 'axios'

export const useSystemStore = defineStore('system', {
	state: () => ({
      tags: [],
      institutions: [],
      events: [],
      news: [],
      document: null,
      error: null,
      loading: true,
      adminMode: false,
      version: "2.0.0"
   }),
   getters: {
      hasError: state => {
         return state.error != null
      },
      publishedNews: state => {
         let out = state.news.filter(news => news.published)
         return out
      }
   },
   actions: {
      setError(error) {
         this.error = error
         setTimeout( ()=>{ this.error = ""}, 10000)
      },
      getTags() {
         axios.get("/api/tags").then((response) => {
            this.tags = response.data
         }).catch((error) => {
            this.tags = []
            this.setError("Unable to get tags: " + error.response.data)
         })
      },
      getInstitutions() {
         axios.get("/api/institutions").then((response) => {
            this.institutions =  response.data
         }).catch((error) => {
            this.institutions = []
            this.setError("Unable to get institutions: " + error.response.data)
         })
      },
      getEvents() {
         this.loading = true
         this.events = []
         axios.get("/api/events").then((response) => {
            response.data.forEach( evt => {
               evt.date = evt.date.split("T")[0]
               this.events.push(evt)
            })
            this.loading = false
         }).catch((error) => {
            this.setError("Unable to get events: " + error.response.data)
            this.loading = true
         })
      },
      getNews() {
         this.loading = true
         this.news = []
         axios.get("/api/news").then((response) => {
            this.news = response.data
            this.loading = false
         }).catch((error) => {
            this.setError("Unable to get news: " + error.response.data)
            this.loading = false
         })
      },
      async getPedagogyDocument( key ) {
         this.loading = true
         this.document = null
         return axios.get(`/api/pedagogy/${key}`).then((response) => {
            let doc = response.data
            let content = doc.content
            let idx = content.indexOf("$DOC[")
            while (idx > -1) {
               let idx2 = content.indexOf("]", idx)
               if (idx2 == -1) {
                  break
               }
               let substr = content.substring(idx+5,idx2)
               let data = substr.split("::")
               let label = data[1]
               var fakeDoc = new DOMParser().parseFromString(label, 'text/html');
               let cleanLabel =  fakeDoc.body.textContent || "";
               let url = `<span class="pedagogy-link" data-link="${data[0]}">${cleanLabel}</span>`
               content = content.substring(0,idx) + url + content.substring(idx2+1)
               idx = content.indexOf("$DOC[")
            }
            doc.content = content
            this.document = doc
         }).catch((error) => {
            this.setError("Unable to get pedagogy dicument: " + error.response.data)
         }).finally( ()=>{
            this.loading = false
         })
      },
      async addInstitution(name) {
         let url = `/api/institutions`
         await axios.post(url, { name: name }).then((response) => {
            this.institutions.push( response.data )
         })
      },
   },
})