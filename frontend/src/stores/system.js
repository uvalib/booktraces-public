import { defineStore } from 'pinia'
import axios from 'axios'

export const useSystemStore = defineStore('system', {
	state: () => ({
      tags: [],
      institutions: [],
      error: null,
      loading: true,
      submissionDetail: null,
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
      clearSubmissionDetail() {
         this.submissionDetail = null
      },
      setFileTranscribed(fileID) {
         let idx = this.submissionDetail.files.findIndex( f =>f.id = fileID )
         if (idx > -1 ) {
            let pend = {text: "pending", approved: false}
            this.submissionDetail.files[idx].transcriptions.push( pend )
         }
      },
      updateTranscription(data) {
         let f = this.submissionDetail.files.find( f=> f.id == data.fileID)
         if ( f ) {
            let t = f.transcriptions.find( t => t.id == data.transcriptionID )
            if ( t ) {
               t.text = data.transcription
            }
         }
      },
      setCurrSubPublished(pub) {
         if (this.submissionDetail) {
            this.submissionDetail.published = pub
         }
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
      getSubmissionDetail(id) {
         this.loading = true
         this.clearSubmissionDetail()
         axios.get("/api/submissions/" + id).then((response) => {
            this.submissionDetail = response.data
            this.loading = false
         }).catch((error) => {
            this.setError("Unable to get submission detail: " + error.response.data)
            this.loading = false
         })
      },
   },
})