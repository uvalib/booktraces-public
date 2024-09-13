import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useDetailsStore = defineStore('details', {
	state: () => ({
      working: false,
      submission: null,
      transcribeFile: null,
      transcribeError: "",
   }),
   getters: {
      isTranscribing: state => {
         return state.transcribeFile != null
      },
      hasNext: state => {
         return state.submission.previousId > 0
      },
      hasPrev: state => {
         return state.submission.nextId > 0
      },
   },
   actions: {
      setTranscriptionTarget( imgFile ) {
         this.transcribeFile = imgFile
         this.transcribeError = ""
      },

      cancelTranscription() {
         this.transcribeFile = null
         this.transcribeError = ""
      },

      getSubmission(id) {
         const system = useSystemStore()
         system.loading = true
         this.submission = null
         axios.get("/api/submissions/" + id).then((response) => {
            this.submission = response.data
            system.loading  = false
         }).catch((error) => {
            system.setError("Unable to get submission detail: " + error.response.data)
            system.loading  = false
         })
      },

      async submitTranscription(transcription, name, email) {
         this.transcribeError = ""
         this.working = true
         let req = {
            transcription: transcription,
            name: name,
            email: email,
            submissionID: this.submission.id,
            fileID:  this.transcribeFile.id
         }
         await axios.post("/api/transcription", req).then(() => {
            let idx = this.submission.files.findIndex( f =>f.id = this.transcribeFile.id )
            if (idx > -1 ) {
               let pend = {text: "pending", approved: false}
               this.submission.files[idx].transcriptions.push( pend )
            }
            this.transcribeFile  = null
            this.working = false
         }).catch ( err => {
            this.transcribeError = err.response.data
            this.working = false
         })
      },
   },
})