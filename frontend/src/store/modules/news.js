import axios from 'axios'

const news = {
   namespaced: true,
   state: {
      list: []
   },

   getters: {
      published: state => {
         let out = state.list.filter(news => news.published)
         return out
      }
   },

   mutations: {
      setNews(state, news) {
         state.list = news
      },
      addPlaceholder(state) {
         state.list.unshift({ content: "", title: "" })
      },
      cancelAdd(state) {
         state.list.shift()
      },
      addNews(state, newNews) {
         state.list[0] =  Object.assign({}, newNews)
      },
      updateNews(state, modified) {
         state.list = state.list.map(n => {
            if (n.id === modified.id) {
               return Object.assign({}, n, modified)
            }
            return n
         })
      },
      deleteNews(state, id) {
         state.list.some( function(e,idx) {
           if (e.id == id) {
             state.list.splice(idx, 1)
             return true
           }
           return false
         })
       },
   },

   actions: {
      getAll(ctx) {
         axios.get("/api/news").then((response) => {
            ctx.commit('setNews', response.data)
         }).catch((error) => {
            ctx.commit('setNews', [])
            ctx.commit('setError', "Unable to get news: " + error.response.data)
         })
      },
      togglePublication(ctx, index) {
         let news = Object.assign({}, ctx.state.list[index])
         news.published = !news.published
         axios.put("/api/admin/news/" + news.id, news).then((/*response*/) => {
            ctx.commit("updateNews", news)
         }).catch((error) => {
            ctx.commit("setError", error.response.data, { root: true })
         })
      },
      updateNews(ctx, modified) {
         return new Promise((resolve, reject) => {
           axios.put("/api/admin/news/"+modified.id, modified).then((/*response*/)  =>  {
             ctx.commit("updateNews", modified)
             resolve()
           }).catch((error) => {
             ctx.commit("setError", error.response.data, {root: true}) 
             reject(error)
           })  
         })
       },
       deleteNews( ctx, id ) {
         axios.delete("/api/admin/news/"+id).then((/*response*/)  =>  {
           ctx.commit("deleteNews", id)
         }).catch((error) => {
           ctx.commit("setError",error.response.data, {root: true}) 
         })
       },
       addNews(ctx, newNews) {
         return new Promise((resolve, reject) => {
           let data = {title: newNews.title, content: newNews.content}
           axios.post("/api/admin/news/", data).then((response)  =>  {
             ctx.commit("addNews", response.data)
             resolve()
           }).catch((error) => {
             ctx.commit("setError", error.response.data, {root: true}) 
             reject(error)
           })  
         })
       },
   }
}

export default news