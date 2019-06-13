import axios from 'axios'

const events = {
   namespaced: true,
   state: {
      list: []
   },

   mutations: {
      setEvents(state, events) {
         state.list = events
      },
      addEventPlaceholder(state) {
         state.list.unshift({ date: "", description: "" })
      },
      cancelAddEvent(state) {
         state.list.shift()
      },
      addEvent(state, newEvent) {
         state.list[0] = Object.assign({}, newEvent)
      },
      deleteEvent(state, id) {
         state.list.some(function (e, idx) {
            if (e.id == id) {
               state.list.splice(idx, 1)
               return true
            }
            return false
         })
      },
      updateEvent(state, modified) {
         state.list = state.list.map(n => {
            if (n.id === modified.id) {
               return Object.assign({}, n, modified)
            }
            return n
         })
      },
   },

   actions: {
      getAll(ctx) {
         axios.get("/api/events").then((response) => {
            ctx.commit('setEvents', response.data)
         }).catch((error) => {
            ctx.commit('setEvents', [])
            ctx.commit('setError', "Unable to get events: " + error.response.data)
         })
      },
      deleteEvent(ctx, id) {
         axios.delete("/api/admin/events/" + id).then((/*response*/) => {
            ctx.commit("deleteEvent", id)
         }).catch((error) => {
            ctx.commit("setError", error.response.data, { root: true })
         })
      },
      updateEvent(ctx, modified) {
         return new Promise((resolve, reject) => {
            axios.put("/api/admin/events/" + modified.id, modified).then((/*response*/) => {
               ctx.commit("updateEvent", modified)
               resolve()
            }).catch((error) => {
               ctx.commit("setError", error.response.data, { root: true })
               reject(error)
            })
         })
      },
      addEvent(ctx, newEvent) {
         return new Promise((resolve, reject) => {
            axios.post("/api/admin/events/", newEvent).then((response) => {
               ctx.commit("addEvent", response.data)
               resolve()
            }).catch((error) => {
               ctx.commit("setError", error.response.data, { root: true })
               reject(error)
            })
         })
      },
   }
}

export default events