import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const usePedagogyStore = defineStore('pedagogy', {
	state: () => ({
      document: null,
      list: []
   }),

   actions: {
      setDocument(doc) {
         if ( doc == null ) {
            this.document = null
            return
         }
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
      },

      getList() {
         const system = useSystemStore()
         system.loading = true
         this.list = []
         axios.get(`/api/admin/pedagogy`).then((response) => {
            response.data.forEach( doc => {
               doc.createdAt = doc.createdAt.split("T")[0]
               this.list.push( doc )
            })
         }).catch((error) => {
            system.setError(error)
         }).finally( ()=>{
            system.loading = false
         })
      },

      async getDocument(key) {
         const system = useSystemStore()
         system.loading = true
         return axios.get(`/api/pedagogy/${key}`).then((response) => {
            this.setDocument(response.data)
         }).catch((error) => {
            this.setDocument(null)
            system.setError(error)
         }).finally( ()=>{
            system.loading = false
         })
      },

      deleteDocument( key ) {
         const system = useSystemStore()
         system.loading = true
         axios.delete(`/api/admin/pedagogy/${key}`).then((_response) => {
            let idx = this.list.findIndex( d => d.key == key)
            if ( idx > -1) {
               this.list.splice(idx,1)
            }
         }).catch((error) => {
            ctx.commit("setError", error, {root: true})
         }).finally( ()=>{
            system.loading = false
         })
      },

      updateDocument( doc ) {
         const system = useSystemStore()
         system.loading = true
         let docIdx = this.list.findIndex( d => d.id == doc.id)
         return axios.put(`/api/admin/pedagogy/${this.list[docIdx].key}`, doc).then((response) => {
            this.list[docIdx] = doc
         }).catch((error) => {
            system.setError(error)
         }).finally( ()=>{
            system.loading = false
         })
      },

      addDocument( doc ) {
         const system = useSystemStore()
         system.loading = true
         return axios.post(`/api/admin/pedagogy`, doc).then((response) => {
            this.list.push( response.data )
         }).catch((error) => {
            system.setError(error)
         }).finally( ()=>{
            system.loading = false
         })
      }
   },
})