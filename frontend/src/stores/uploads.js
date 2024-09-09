import { defineStore } from 'pinia'
import { useSystemStore } from "@/stores/system"
import axios from 'axios'

export const useUploadStore = defineStore('uploads', {
   state: () => ({
      uploadID: null,
      title: "",
      author: "",
      publication: "",
      callNumber: "",
      description: "",
      submitter: "",
      email: "",
      uploadedFiles: [],
      tags: [],
      institution: null,
      working: false,
      submitted: false,
      error: "",
   }),

   getters: {
      canSubmit: state => {
         return state.uploadID != null && state.title != "" &&  state.author != "" &&
            state.institution != null && state.uploadedFiles.length > 0 && state.submitter != "" &&
            state.email != "" && state.tags.length > 0
      }
   },

   actions: {
      async submit() {
         this.working = true
         this.error = ""
         if ( typeof this.institution == 'string' ) {
            // if a new institution has been added, institution will just be the string name
            // in other cases, institution is an object with id and name
            const system = useSystemStore()
            await system.addInstitution( this.institution )
            let newInstitution = system.institutions.find( i => i.name == this.institution )
            this.institution = newInstitution
         }
         let form = {
            uploadID: this.uploadID,
            title: this.title,
            author: this.title,
            publication: this.title,
            institution_id: this.institution.id,
            institution: this.institution.name,
            callNumber: this.title,
            description: this.title,
            files: this.uploadedFiles,
            submitter: this.submitter,
            email: this.email,
            tags: this.tags
         }

         axios.post("/api/submit", form).then(() =>  {
            this.working = false
            this.submitted = true
            this.router.push("/thanks")
         }).catch((error) => {
            this.error = error.response.data
         })
      },
      addImage( file ) {
         this.uploadedFiles.push( file.name )
      },
      getUploadID() {
         this.$reset()
         axios.get("/api/identifier").then((response) => {
            this.uploadID = response.data
         }).catch((error) => {
            this.error = "Unable to get prepare submission: " + error.response.data
         })
      },
      removeUploadedImage(filename) {
         let index = this.uploadedFiles.indexOf(filename)
         if (index !== -1) {
            this.uploadedFiles.splice(index, 1)
         }
         axios.delete("/api/upload/" + filename + "?key=" + this.uploadID)
      },
   }
})
