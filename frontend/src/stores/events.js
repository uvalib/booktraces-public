import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useEventsStore = defineStore('events', {
	state: () => ({
      list: []
   }),
   actions: {
      addEventPlaceholder() {
         this.list.unshift({ date: "", description: "" })
      },
      cancelAddEvent() {
         this.list.shift()
      },
      addEvent(newEvent) {
         this.list[0] = Object.assign({}, newEvent)
      },
      // updateEvent(modified) {
      //    this.list = this.list.map(n => {
      //       if (n.id === modified.id) {
      //          return Object.assign({}, n, modified)
      //       }
      //       return n
      //    })
      // },

      getAll() {
         this.list = []
         axios.get("/api/events").then((response) => {
            response.data.forEach( evt => {
               evt.date = evt.date.split("T")[0]
               this.list.push(evt)
            })
         }).catch((error) => {
            useSystemStore().setError("Unable to get events: " + error.response.data)
         })
      },
      deleteEvent(id) {
         axios.delete("/api/admin/events/" + id).then(() => {
            let delIdx = this.list.findIndex( n => n.id == id)
            if (delIdx > -1) {
               this.list.splice(delIdx, 1)
            }
         }).catch((error) => {
            useSystemStore().setError("Unable to delete event: " + error.response.data)
         })
      },
      // FIXME
      // updateEvent(modified) {
      //    return new Promise((resolve, reject) => {
      //       axios.put("/api/admin/events/" + modified.id, modified).then(() => {
      //          this.commit("updateEvent", modified)
      //          resolve()
      //       }).catch((error) => {
      //          this.commit("setError", error.response.data, { root: true })
      //          reject(error)
      //       })
      //    })
      // },
      // addEvent(newEvent) {
      //    return new Promise((resolve, reject) => {
      //       axios.post("/api/admin/events/", newEvent).then((response) => {
      //          this.commit("addEvent", response.data)
      //          resolve()
      //       }).catch((error) => {
      //          this.commit("setError", error.response.data, { root: true })
      //          reject(error)
      //       })
      //    })
      // },
   }
})