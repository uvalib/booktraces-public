import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useNewsStore = defineStore('news', {
	state: () => ({
      list: []
   }),

   getters: {
      published: state => {
         let out = state.list.filter(news => news.published)
         return out
      }
   },

   actions: {
      addPlaceholder() {
         this.list.unshift({ content: "", title: "" })
      },
      cancelAdd() {
         this.list.shift()
      },
      addNews(newNews) {
         this.list[0] = Object.assign({}, newNews)
      },
      getAll() {
         this.list = []
         axios.get("/api/news").then((response) => {
            this.list = response.data
         }).catch((error) => {
            useSystemStore().setError("Unable to get news: " + error.response.data)
         })
      },

      togglePublication( index) {
         this.list[index].published =  !this.list[index].published
         axios.put("/api/admin/news/" + news.id, news).catch((error) => {
            useSystemStore().setError("Unable to toggle news publication: " + error.response.data)
         })
      },

      // FIXME
      // updateNews(ctx, modified) {
      //    return new Promise((resolve, reject) => {
      //       axios.put("/api/admin/news/" + modified.id, modified).then(() => {
      //          ctx.commit("updateNews", modified)
      //          resolve()
      //       }).catch((error) => {
      //          useSystemStore().setError("Unable to get news: " + error.response.data)
      //          reject(error)
      //       })
      //    })
      // },

      deleteNews(id) {
         axios.delete("/api/admin/news/" + id).then(() => {
            let delIdx = this.list.findIndex( n => n.id == id)
            if (delIdx > -1) {
               this.list.splice(delIdx, 1)
            }
         }).catch((error) => {
            useSystemStore().setError("Unable to delete news: " + error.response.data)
         })
      },

      //FIXME
      // addNews(ctx, newNews) {
      //    return new Promise((resolve, reject) => {
      //       let data = { title: newNews.title, content: newNews.content }
      //       axios.post("/api/admin/news/", data).then((response) => {
      //          ctx.commit("addNews", response.data)
      //          resolve()
      //       }).catch((error) => {
      //          useSystemStore().setError("Unable to get news: " + error.response.data)
      //          reject(error)
      //       })
      //    })
      // },
   }
})