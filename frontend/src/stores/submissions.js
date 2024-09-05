import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useSubmissionsStore = defineStore('submissions', {
   state: () => ({
      uploadID: null,
      uploadedFiles: [],
      total: 0,
      thumbs: [],
      archives: [],
      recents: [],
      currPage: 0,
      showSearch: false,
      searchResults: [],
      query: "",
      archiveDate: "",
      targetInstitution: "",
      tgtTag: "",
      news: [],
      transcribeFile: null,
      transcribeError: ""
   }),

   getters: {
      isTranscribing: state => {
         return state.transcribeFile != null
      },
      thumbsCount: state => {
         return state.thumbs.length
      },
      searchHitCount: state => {
         return state.searchResults.length
      },
      searchType: state => {
         if (state.tgtTag.length > 0) {
            return "tag"
         }
         if (state.archiveDate.length > 0) {
            return "archive"
         }
         if (state.targetInstitution != "") {
            return "institution"
         }
         if (state.query.length > 0) {
            return "query"
         }
         return "none"
      },
   },

   actions: {
      cancelTranscribe() {
         this.transcribeFile = null
      },
      updateSearchQuery(q) {
         this.query = q
         this.searchResults = []
         this.archiveDate = ""
         this.tgtTag = ""
         this.targetInstitution = ""
      },
      clearThumbs() {
         this.currPage = 0
         this.total = 0
         this.thumbs = []
      },
      clearUploadedFiles() {
         this.uploadedFiles = []
      },
      addUploadedFile(filename) {
         this.uploadedFiles.push(filename)
      },
      setNews(news) {
         this.news = news
      },

      search() {
         const system = useSystemStore()
         system.loading = true
         this.showSearch = false
         axios.get("/api/search?q=" + this.query).then((response) => {
            this.searchResults = response.data
            system.loading = false
         }).catch((error) => {
            system.setError("Unable to get search results: " + error.response.data)
            this.searchResults = []
            system.loading = false
         })
      },
      searchInstitutions() {
         const system = useSystemStore()
         system.loading = true
         this.showSearch = false
         this.searchResults = []
         this.query = ""
         this.archiveDate = ""
         this.tgtTag = ""
         axios.get("/api/search?i=" + this.targetInstitution).then((response) => {
            this.searchResults = response.data
            system.loading = false
         }).catch((error) => {
            system.setError("Unable to get institutions: " + error.response.data)
            this.searchResults = []
            system.loading = false
         })
      },
      getArchive(date) {
         const system = useSystemStore()
         system.loading = true
         this.showSearch = false
         this.archiveDate = date
         this.searchResults = []
         this.query = ""
         this.tgtTag = ""
         this.targetInstitution = ""
         axios.get("/api/search?a=" + date).then((response) => {
            this.searchResults = response.data
            system.loading = false
         }).catch((error) => {
            system.setError("Unable to get archives: " + error.response.data)
            this.searchResults = []
            system.loading = false
         })
      },
      getTaggedSubmissions(tag) {
         const system = useSystemStore()
         system.loading = true
         this.showSearch = false
         this.searchResults = []
         this.query = ""
         this.archiveDate = ""
         this.targetInstitution = ""
         this.tgtTag = tag
         axios.get("/api/search?t=" + tag).then((response) => {
            this.searchResults = response.data
            system.loading = false
         }).catch((error) => {
            system.setError("Unable to get archives: " + error.response.data)
            this.searchResults = []
            system.loading = false
         })
      },
      getArchiveDates() {
         if (this.archives.length > 0) {
            return
         }
         axios.get("/api/archives").then((response) => {
            this.archives = response.data
         }).catch((error) => {
            const system = useSystemStore()
            system.setError("Unable to get archives: " + error.response.data)
         })
      },
      getRecentSubmissions() {
         if (this.recents.length > 0) {
            return
         }
         axios.get("/api/recents").then((response) => {
            this.recents = response.data
         }).catch((error) => {
            const system = useSystemStore()
            system.setError("Unable to get recent submissions: " + error.response.data)
         })
      },
      getRecentThumbs() {
         const system = useSystemStore()
         system.loading = true
         this.currPage++
         axios.get("/api/thumbs?page=" + this.currPage).then((response) => {
            this.total = response.data.total
            response.data.thumbs.forEach((thumb) => {
               this.thumbs.push(thumb)
            })
            system.loading = false
         }).catch((error) => {
            system.setError("Unable to get recent thumbs: " + error.response.data)
            system.loading = false
         })
      },
      getUploadID() {
         this.uploadID = null
         axios.get("/api/identifier").then((response) => {
            this.uploadID = response.data
         }).catch((error) => {
            const system = useSystemStore()
            system.setError("Unable to get uploadID: " + error.response.data)
         })
      },
      removeUploadedFile(filename) {
         let index = this.uploadedFiles.indexOf(filename)
         if (index !== -1) {
            this.uploadedFiles.splice(index, 1)
         }
         axios.delete("/api/upload/" + filename + "?key=" + this.getters.uploadID)
      },
      submitTranscription(transcription) {
         this.transcribeError = ""
         if (this.transcribeFile == null) {
             this.transcribeError = "File is missing"
            return
         }
         axios.post("/api/transcription", {fileID: this.transcribeFile.id, transcription: transcription}).then((_response) => {
            this.transcribeFile = null
         }).catch ( err => {
            this.transcribeError = err.response.data
         })
      }
   }
})
