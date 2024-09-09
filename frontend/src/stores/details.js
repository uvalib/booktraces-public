import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useDetailsStore = defineStore('details', {
	state: () => ({
      working: true,
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
         this.transcribeFile = file
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

      approveTranscription() {
         this.working = true
         axios.post(`/api/admin/submissions/${subID}/transcription/${this.submission.id}/approve`).then((response) => {
            this.working = false
            this.submission = response.data
         }).catch ( err => {
            useSystemStore().setError(err.response.data)
            this.working = false
         })
      },

      async updateTranscription( data ) {
         let subID = data.submissionID
         let transID = data.transcriptionID
         let txt = data.transcription
         this.working = true
         try {
            await axios.put(`/api/admin/submissions/${subID}/transcription/${transID}`, {transcription: txt})
            this.working = false
            let f = this.submission.files.find( f=> f.id == data.fileID)
            if ( f ) {
               let t = f.transcriptions.find( t => t.id == data.transcriptionID )
               if ( t ) {
                  t.text = data.transcription
               }
            }
         } catch ( err ) {
            useSystemStore().setError(err.response.data)
            this.working = false
         }
      },

      deleteTranscription( id ) {
         this.working = true
         axios.delete(`/api/admin/submissions/${this.submission.id}/transcription/${id}`).then((response) => {
            this.working = false
            this.submission = response.data
         }).catch ( err => {
            useSystemStore().setError(err.response.data)
            this.working = false
         })
      },

      submitTranscription(data) {
         this.transcribeError = ""
         if (this.transcribeFile == null) {
            this.transcribeError =  "File is missing"
            return
         }

         if (data.transcription == "") {
            this.transcribeError =  "Please enter a transcription"
            return
         }
         if (data.name == "") {
            this.transcribeError =  "Name is required"
            return
         }
         if (data.email == "") {
            this.transcribeError =  "Email is required"
            return
         }
         this.working = true
         data.submissionID = this.submission.id
         data.fileID = this.transcribeFile.id
         axios.post("/api/transcription", data).then((_response) => {
            ctx.commit("setFileTranscribed",  ctx.state.file.id, {root:true})
            let idx = this.submission.files.findIndex( f =>f.id = data.fileID )
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

      // FIXME? NOT SURE WHICH TO USE
      // THIS WAS A DUPLICATE FROM SUBMISSIONS.JS
      // submitTranscription(transcription) {
      //    this.transcribeError = ""
      //    if (this.transcribeFile == null) {
      //        this.transcribeError = "File is missing"
      //       return
      //    }
      //    axios.post("/api/transcription", {fileID: this.transcribeFile.id, transcription: transcription}).then((_response) => {
      //       this.transcribeFile = null
      //    }).catch ( err => {
      //       this.transcribeError = err.response.data
      //    })
      // }

   },
})