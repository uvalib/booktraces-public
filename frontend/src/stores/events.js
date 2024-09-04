import { defineStore } from 'pinia'
import axios from 'axios'

export const useEventsStore = defineStore('events', {
	state: () => ({
      list: []
   }),
   actions: {
      setEvents(events) {
         this.list = this.list.splice(0,  this.list.length)
         events.forEach( evt => {
            evt.date = evt.date.split("T")[0]
            this.list.push(evt )
         })
      },
      addEventPlaceholder() {
         this.list.unshift({ date: "", description: "" })
      },
      cancelAddEvent() {
         this.list.shift()
      },
      addEvent(newEvent) {
         this.list[0] = Object.assign({}, newEvent)
      },
      deleteEvent(id) {
         this.list.some(function (e, idx) {
            if (e.id == id) {
               this.list.splice(idx, 1)
               return true
            }
            return false
         })
      },
      updateEvent(modified) {
         this.list = this.list.map(n => {
            if (n.id === modified.id) {
               return Object.assign({}, n, modified)
            }
            return n
         })
      },

      getAll() {
         axios.get("/api/events").then((response) => {
            this.commit('setEvents', response.data)
         }).catch((error) => {
            this.commit('setEvents', [])
            this.commit('setError', "Unable to get events: " + error.response.data)
         })
      },
      deleteEvent(id) {
         axios.delete("/api/admin/events/" + id).then((/*response*/) => {
            this.commit("deleteEvent", id)
         }).catch((error) => {
            this.commit("setError", error.response.data, { root: true })
         })
      },
      updateEvent(modified) {
         return new Promise((resolve, reject) => {
            axios.put("/api/admin/events/" + modified.id, modified).then((/*response*/) => {
               this.commit("updateEvent", modified)
               resolve()
            }).catch((error) => {
               this.commit("setError", error.response.data, { root: true })
               reject(error)
            })
         })
      },
      addEvent(newEvent) {
         return new Promise((resolve, reject) => {
            axios.post("/api/admin/events/", newEvent).then((response) => {
               this.commit("addEvent", response.data)
               resolve()
            }).catch((error) => {
               this.commit("setError", error.response.data, { root: true })
               reject(error)
            })
         })
      },
   }
})