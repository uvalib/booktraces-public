import { defineStore } from 'pinia'
import axios from 'axios'

export const useSystemStore = defineStore('system', {
	state: () => ({
      tags: [],
      institutions: [],
      error: null,
      loading: true,
      adminMode: false,
      appVersion: "2.0.0"
   }),
   getters: {
      hasError: state => {
         return state.error != null
      }
   },
   actions: {
      setAdminMode(mode) {
         this.adminMode = mode
      },
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
      async addInstitution(name) {
         let url = `/api/institutions`
         return axios.post(url, { name: name }).then((response) => {
            this.institutions.push( response.data )
         })
      },
   },
})