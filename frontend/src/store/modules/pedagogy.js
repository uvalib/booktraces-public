import axios from 'axios'

const pedagogy = {
   namespaced: true,
   state: {
      document: null,
      list: []
   },

   getters: {
   },

   mutations: {
      setDocument(state, doc) {
         let content = doc.content
         let idx = content.indexOf("$DOC[")
         while (idx > -1) {
            let idx2 = content.indexOf("]", idx)
            if (idx2 == -1) {
               break
            }
            let data = content.substring(idx+5,idx2).split(":")
            let url = `<span class="pedagogy-link" data-link="${data[0]}">${data[1]}</span>`
            content = content.substring(0,idx) + url + content.substring(idx2+1)
            idx = content.indexOf("$DOC[")
         }
         doc.content = content
         state.document = doc
      },
      clearList(state) {
         state.list = []
      },
      setDocumentList(state, list) {
         list.forEach( doc => {
            doc.createdAt = doc.createdAt.split("T")[0]
            state.list.push( doc )
         })
      },
      deleteDocument(state, key) {
         let idx = state.list.findIndex( d => d.key == key)
         if ( idx > -1) {
            state.list.splice(idx,1)
         }
      },
      updateDocument(state, doc) {
         let idx = state.list.findIndex( d => d.id == doc.id)
         state.list.splice(idx, 1)
         state.list.push(doc)
      },
      addDocument(state, doc) {
         state.list.push(doc)
      }
   },

   actions: {
      getList(ctx) {
         ctx.commit("setLoading", true, {root: true})
         ctx.commit("clearList")
         axios.get(`/api/admin/pedagogy`).then((response) => {
            ctx.commit('setDocumentList', response.data)
         }).catch((error) => {
            ctx.commit("setError", error, {root: true})
         }).finally( ()=>{
            ctx.commit("setLoading", false, {root: true})
         })
      },
      get(ctx, key) {
         ctx.commit("setLoading", true, {root: true})
         return axios.get(`/api/pedagogy/${key}`).then((response) => {
            ctx.commit('setDocument', response.data)
         }).catch((_error) => {
            ctx.commit('setDocument', null)
         }).finally( ()=>{
            ctx.commit("setLoading", false, {root: true})
         })
      },
      deleteDocument(ctx, key) {
         ctx.commit("setLoading", true, {root: true})
         axios.delete(`/api/admin/pedagogy/${key}`).then((_response) => {
            ctx.commit('deleteDocument', key)
         }).catch((error) => {
            ctx.commit("setError", error, {root: true})
         }).finally( ()=>{
            ctx.commit("setLoading", false, {root: true})
         })
      },
      updateDocument(ctx, doc) {
         ctx.commit("setLoading", true, {root: true})
         let old = ctx.state.list.find( d => d.id == doc.id)
         return axios.put(`/api/admin/pedagogy/${old.key}`, doc).then((response) => {
            ctx.commit('updateDocument', response.data)
         }).catch((error) => {
            ctx.commit("setError", error, {root: true})
         }).finally( ()=>{
            ctx.commit("setLoading", false, {root: true})
         })
      },
      addDocument(ctx, doc) {
         ctx.commit("setLoading", true, {root: true})
         return axios.post(`/api/admin/pedagogy`, doc).then((response) => {
            ctx.commit('addDocument', response.data)
         }).catch((error) => {
            ctx.commit("setError", error, {root: true})
         }).finally( ()=>{
            ctx.commit("setLoading", false, {root: true})
         })
      }
   }
}

export default pedagogy