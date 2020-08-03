import axios from 'axios'

const transcribe = {
   namespaced: true,

   // root state object. Holds all of the state for the system
   state: {
      file: null,
      error: "",
      submitting: false
   },

   getters: {
      isTranscribing: state => {
         return state.file != null
      },
   },

   mutations: {
      setFile(state, file) {
         state.file = file
         state.error = ""
         state.submitting = false 
      },
      setTranscribeError(state, err) {
         state.error = err
      },
      cancel(state) {
         state.file = null
         state.submitting = false 
         state.error = ""
      },
      setSubmitting(state, flag) {
         state.submitting = flag
         state.error = ""
      },
   },

   actions: {
      approve(ctx, id) {
         let subID = ctx.rootState.submissionDetail.id
         ctx.commit("setSubmitting", true)
         axios.post(`/api/admin/submissions/${subID}/transcription/${id}/approve`).then((response) => {
            ctx.commit("setSubmitting", false) 
            ctx.commit("setSubmissionDetail", response.data, {root: true})
         }).catch ( err => {
            ctx.commit("setError", err.response.data, {root: true})
            ctx.commit("setSubmitting", false) 
         })
      },
      async update(ctx, data) {
         let subID = data.submissionID
         let transID = data.transcriptionID
         let txt = data.transcription
         ctx.commit("setSubmitting", true)
         try {
            await axios.put(`/api/admin/submissions/${subID}/transcription/${transID}`, {transcription: txt})
            ctx.commit("setSubmitting", false) 
            ctx.commit("updateTranscription", data, {root: true})
         } catch ( err ) {
            ctx.commit("setError", err.response.data, {root: true})
            ctx.commit("setSubmitting", false) 
         }
      },
      delete(ctx, id) {
         let subID = ctx.rootState.submissionDetail.id
         ctx.commit("setSubmitting", true)
         axios.delete(`/api/admin/submissions/${subID}/transcription/${id}`).then((response) => {
            ctx.commit("setSubmitting", false) 
            ctx.commit("setSubmissionDetail", response.data, {root: true})
         }).catch ( err => {
            ctx.commit("setError", err.response.data, {root: true})
            ctx.commit("setSubmitting", false) 
         })
      },
      submit(ctx, data) {
         ctx.commit("setTranscribeError", "")
         if (ctx.state.file == null) {
            ctx.commit("setTranscribeError", "File is missing")
            return
         }
         
         if (data.transcription == "") {
            ctx.commit("setTranscribeError", "Please enter a transcription")
            return
         }
         if (data.name == "") {
            ctx.commit("setTranscribeError", "Name is required")
            return
         }
         if (data.email == "") {
            ctx.commit("setTranscribeError", "Email is required")
            return
         }
         ctx.commit("setSubmitting", true)
         data.submissionID = ctx.rootState.submissionDetail.id
         data.fileID = ctx.state.file.id
         axios.post("/api/transcription", data).then((_response) => {
            ctx.commit("setFileTranscribed",  ctx.state.file.id, {root:true}) 
            ctx.commit("setFile", null) 
            ctx.commit("setSubmitting", false) 
         }).catch ( err => {
            ctx.commit("setTranscribeError", err.response.data)
            ctx.commit("setSubmitting", false) 
         })
      }
   }
}

export default transcribe
