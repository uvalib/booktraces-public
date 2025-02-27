<template>
   <div class="admin-submission">
      <h2>
         <span>System Admin Panel</span>
         <span class="login">
            <b>Logged in as:</b>{{admin.loginName}}
         </span>
      </h2>
      <BTSpinner v-if="system.loading==true || details.submission == null" message="Loading details..." />
      <template v-else>
         <div class="actions" v-if="!edit">
            <router-link to="/admin">
               <i class="pi pi-arrow-left"></i>&nbsp;Back to Submissions
            </router-link>
            <div class="buttons">
               <AdminEditSubmission @save="saveEdits" :details="details.submission" />
               <Button @click="togglePublish" size="small" :severity="publishSeverity" :label="publishLabel"/>
               <Button size="small" label="Delete" @click="deleteSubmission" severity="danger"/>
            </div>
         </div>
         <div v-if="!edit" class="details">
            <table class="submit-info">
               <tbody>
                  <tr>
                     <td class="label">Visible to public:</td>
                     <td class="value">{{ published }}</td>
                  </tr>
                  <tr>
                     <td class="label">Title:</td>
                     <td class="value">{{details.submission.title}}</td>
                  </tr>
                  <tr>
                     <td class="label">Author:</td>
                     <td class="value">{{details.submission.author}}</td>
                  </tr>
                  <tr>
                     <td class="label">Publication details:</td>
                     <td class="value">{{details.submission.publication}}</td>
                  </tr>
                  <tr>
                     <td class="label">Institution:</td>
                     <td class="value">{{details.submission.institution}}</td>
                  </tr>
                  <tr>
                     <td class="label">Call number:</td>
                     <td class="value">{{details.submission.callNumber}}</td>
                  </tr>
                  <tr>
                     <td class="label">Submitted by:</td>
                     <td class="value">{{details.submission.submitter}} ( {{details.submission.email}} )</td>
                  </tr>
                  <tr>
                     <td class="label">Submitted on:</td>
                     <td class="value">{{ details.submission.submittedAt.split("T")[0] }}</td>
                  </tr>
                  <tr>
                     <td class="label">Tags:</td>
                     <td class="value">{{submissionTagCSV}}</td>
                  </tr>
                  <tr>
                     <td class="label">Description:</td>
                     <td class="value">
                        <span v-html="description"></span>
                     </td>
                  </tr>
               </tbody>
            </table>
         </div>
         <div class="images-wrapper">
            <div class="image-controls">
               <div class="section">
                  <Button id="zoom-in" icon="pi pi-search-plus" rounded severity="secondary"/>
                  <Button id="zoom-out" icon="pi pi-search-minus" rounded severity="secondary"/>
               </div>
               <div class="section wide" v-if="details.submission.files.length > 1">
                  <Button icon="pi pi-chevron-left" rounded severity="secondary"
                     :disabled="currImageIdx==0" @click="prevImage"/>
                  <span>Image {{ currImageIdx+1 }} of {{ details.submission.files.length }}</span>
                  <Button icon="pi pi-chevron-right" rounded severity="secondary"
                     :disabled="currImageIdx==details.submission.files.length-1" @click="nextImage"/>
               </div>
            </div>
            <div class="image-content">
               <div class="zoom-wrap">
                  <div id="dragon" class="viewer"></div>
                  <Button severity="info" label="Rotate Right" @click="rotateClicked" />
               </div>
               <div class="transcription-panel">
                  <div class="head">Transcriptions</div>
                  <div v-if="hasTranscription == false" class="none">None</div>
                  <template v-else>
                     <div class="transcribe-acts">
                        <div class="buttons">
                           <Button v-if="isTranscriptionApproved" aria-label="approve transcription"
                              icon="pi pi-check-circle"  size="small" @click="approveTranscription(currImage)"/>
                           <Button icon="pi pi-file-edit"  severity="secondary" size="small" @click="editTranscription(currImage)"/>
                           <Button icon="pi pi-trash"  severity="danger" size="small" @click="deleteTranscription(currImage)"/>
                        </div>
                        <span class="paging">
                           <Button icon="pi pi-chevron-left" rounded severity="secondary" size="small"
                              :disabled="transcriptionIdx == 0" @click="priorTran()"/>
                           <span>{{transcriptionIdx+1}} of {{currImage.transcriptions.length}}</span>
                           <Button icon="pi pi-chevron-right" rounded severity="secondary" size="small"
                              :disabled="transcriptionIdx == currImage.transcriptions.length-1" @click="nextTran()"/>
                        </span>
                     </div>
                     <div class="transcription-info" v-if="currImage.transcriptions.length > 0">
                        <table>
                           <tbody>
                              <tr>
                                 <td class="label">Date:</td>
                                 <td>{{ transcribeDate }}</td>
                              </tr>
                              <tr>
                                 <td class="label">Submitter:</td>
                                 <td>{{ transcriber }}</td>
                              </tr>
                              <tr>
                                 <td class="label">Status:</td>
                                 <td>{{ transcribeStatus }}</td>
                              </tr>
                           </tbody>
                        </table>
                        <div class="edit-trans" v-if="editTrans">
                           <Textarea id="description" v-model="workingTrans" rows="4" />
                           <div class="actions">
                              <span class="buttons">
                                 <Button severity="secondary" label="Cancel" @click="cancelEditTranscription" />
                                 <Button severity="info" label="Submit" @click="submitEditTranscription" />
                              </span>
                           </div>
                        </div>
                        <template v-else>
                           <pre class="transcription">{{currImage.transcriptions[transcriptionIdx].text}}</pre>
                        </template>
                     </div>
                  </template>
               </div>
            </div>
         </div>
      </template>
   </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import { useAdminStore } from "@/stores/admin"
import { useSystemStore } from "@/stores/system"
import { useDetailsStore } from "@/stores/details"
import { useRoute, useRouter } from 'vue-router'
import { useConfirm } from "primevue/useconfirm"
import AdminEditSubmission from "@/components/AdminEditSubmission.vue"
import Textarea from 'primevue/textarea'
import OpenSeadragon from "openseadragon"
import { useToast } from 'primevue/usetoast'

const toast = useToast()
const confirm = useConfirm()
const system = useSystemStore()
const admin = useAdminStore()
const details = useDetailsStore()
const route = useRoute()
const router = useRouter()
const viewer = ref()
const currImageIdx = ref(0)

const edit = ref(false)
const editTrans = ref(false)
const workingTrans = ref("")
const transcriptionIdx = ref(0)

const published = computed(() => {
   if (details.submission.published) return "YES"
   return "NO"
})

const submissionTagCSV = computed(() => {
   if (details.submission.tags) {
      return details.submission.tags.join(", ")
   }
   return ""
})

const currImage = computed(() => {
   return details.submission.files[currImageIdx.value]
})

const description = computed(() => {
   return details.submission.description.replace(/\r|\r\n/gm, "\n").replace(/\n+/gm, "<br/><br/>")
})
const publishLabel = computed(() => {
   if ( details.submission.published) return "Unpublish"
   return  "Publish"
 })
 const publishSeverity = computed(() => {
   if ( details.submission.published) return "secondary"
   return  "primary"
 })

onMounted(async () => {

   system.getTags()
   system.getInstitutions()

   await details.getSubmission( route.params.id )
   let sources = []
   details.submission.files.forEach( f => {
      sources.push( {type: 'image',url:  f.url} )
   })
   viewer.value = OpenSeadragon({
      id: `dragon`,
      immediateRender: true,
      animationTime: 0.1,
      autoResize: true,
      constrainDuringPan: true,
      imageSmoothingEnabled: true,
      smoothTileEdgesMinZoom: Infinity,
      maxZoomPixelRatio: 10.0,
      showSequenceControl: false,
      zoomInButton:   "zoom-in",
      zoomOutButton:  "zoom-out",
      showNavigator: false,
      showRotationControl: false,
      tileSources: sources,
      sequenceMode: true,
      preserveViewport: true,
      showReferenceStrip: true,
      referenceStripScroll: 'vertical',
   })
   viewer.value.addHandler("page", (data) => {
      currImageIdx.value = data.page
   })
})

const saveEdits = ( async ( edits ) => {
   await admin.updateSubmission( details.submission.id, edits )
   details.getSubmission( details.submission.id )
})

const togglePublish = ( () => {
   let act = "Publish"
   if (details.submission.published) {
      act = "Unpublish"
   }
   confirm.require({
      message: `${act} this submission?`,
      header: `Confirm ${act}`,
      icon: 'pi pi-exclamation-triangle',
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: act
      },
      accept: async () => {
         await admin.updatePublicationStatus(details.submission.id, !details.submission.published )
         details.submission.published = !details.submission.published
      }
   })
})

const deleteSubmission = ( () => {
   confirm.require({
      message: 'Delete this submission?<br/>All data and unloaded files will be permanently lost.<br/><br/>Are you sure?',
      header: 'Confirm Delete',
      icon: 'pi pi-exclamation-triangle',
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: 'Delete'
      },
      accept: async () => {
         await admin.deleteSubmission(details.submission.id)
         router.push("/admin/submissions")
      }
   })
})

const nextImage = ( async () => {
   currImageIdx.value++
   transcriptionIdx.value = 0
   await viewer.value.goToPage( currImageIdx.value )
   console.log( details.submission )
   console.log( details.submission.files.length )
})

const prevImage = ( async () => {
   currImageIdx.value--
   transcriptionIdx.value = 0
   await viewer.value.goToPage( currImageIdx.value )
})

const rotateClicked = ( async () => {
   let imgFile = details.submission.files[currImageIdx.value]
   await admin.rotateImage( details.submission.id, imgFile.url)
   toast.add({
      severity: 'success', summary: 'Rotate Success',
      detail: 'The image has been rotated, but is cached in your browser. Wait a few seconds, clear the cache and reload the page to see the rotated image.' })
})

const hasTranscription = computed(() => {
   let imgFile = details.submission.files[currImageIdx.value]
   if (imgFile.transcriptions.length == 0 ) return false
   console.log("HAS T")
   console.log(imgFile)
   return true
})

const isTranscriptionApproved  = computed(() => {
   let imgFile = details.submission.files[currImageIdx.value]
   if (imgFile ==  null ) return false
   if ( imgFile.transcriptions.length == 0 ) return false
   return imgFile.transcriptions[transcriptionIdx.value].approved==false
})

const nextTran = (( ) => {
   transcriptionIdx.value++
})
const priorTran = (() => {
   transcriptionIdx.value--
})
const transcribeDate = computed( () => {
   let f = details.submission.files[currImageIdx.value]
   if (f.transcriptions.length == 0) return ""
   let t = f.transcriptions[transcriptionIdx.value]
   return t.transcribed_at.split("T")[0]
})
const transcriber = computed( () => {
   let f = details.submission.files[currImageIdx.value]
   if (f.transcriptions.length == 0) return ""
   let t = f.transcriptions[transcriptionIdx.value]
   return `${t.transcriber_email} (${t.transcriber})`
})
const transcribeStatus = computed( () => {
   let f = details.submission.files[currImageIdx.value]
   if (f.transcriptions.length == 0) return ""
   let t = f.transcriptions[transcriptionIdx.value]
   if (t.approved) {
      return "Approved"
   }
   return "Pending"
})
const editTranscription = ((image) => {
   editTrans.value = true
   workingTrans.value = image.transcriptions[transcriptionIdx.value].text
})
const cancelEditTranscription = (() => {
   editTrans.value = false
   workingTrans.value = ""
})
const submitEditTranscription = ( async () => {
   let image = details.submission.files[currImageIdx.value]
   let trans = image.transcriptions[transcriptionIdx.value]
   await admin.updateTranscription(details.submission.id, trans.id, workingTrans.value)
   if (admin.error != "" ) {
      toast.add( {severity: 'error', summary: 'Updae failed',
         detail: `Unable to update transcription: ${admin.error}`})
   } else {
      editTrans.value = false
      trans.text = workingTrans.value
      workingTrans.value = ""
      toast.add( {severity: 'success', summary: 'Updated',
                  detail: `Transcription has been updated`, life: 3000})
   }
})
const deleteTranscription = ((image) => {
   confirm.require({
      message: 'Delete this transcription?<br/>All data will be permanently lost.<br/><br/>Are you sure?',
      header: 'Confirm Delete',
      icon: 'pi pi-exclamation-triangle',
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: 'Delete'
      },
      accept: async () => {
         let transcription = image.transcriptions[transcriptionIdx.value]
         if ( transcription ) {
            await admin.deleteTranscription(details.submission.id, transcription.id)
            if ( admin.error == "") {
               image.transcriptions.splice(transcriptionIdx.value, 1)
               transcriptionIdx.value = 0
               toast.add( {severity: 'success', summary: 'Updated',
                  detail: `Transcription has been deleted`, life: 3000})
            } else {
               toast.add( {severity: 'error', summary: 'Delete failed',
                  detail: `Unable to delete transcription: ${admin.error}`})
            }
         }
      }
   })
})

const approveTranscription = ( async (image) => {
   let t = image.transcriptions[transcriptionIdx.value]
   await admin.approveTranscription(details.submission.id, t.id)
   if ( admin.error == "") {
      image.transcriptions.forEach(trans => trans.approved = false)
      t.approved = true
      toast.add( {severity: 'success', summary: 'Approved',
         detail: `Transcription has been approved and will be visisble to the public`, life: 3000})
   } else {
      toast.add( {severity: 'error', summary: 'Delete failed',
         detail: `Unable to approve transcription: ${admin.error}`})
   }
})
</script>

<style scoped lang="scss">
div.admin-submission {
   min-height: 600px;
   background: white;
   color: #444;

   h3 {
      margin-bottom: 10px;
   }
   h2 {
      display: flex;
      flex-flow: row wrap;
      justify-content: space-between;
      align-items: flex-start;
      font-size: 1.5em;
      font-weight: bold;
      border-bottom: 1px dashed #666;
      font-family: 'Special Elite', cursive;
      padding-bottom: 5px;
      margin-bottom: 15px;
      .login {
         color: #999;
         font-family: sans-serif;
         font-size: 0.8em;
         float: right;
         font-weight: normal;
         b {
            color: #444;
            margin-right: 10px;
         }
      }
   }
   .edit-trans {
      display: flex;
      flex-direction: column;
      gap: 10px;
      padding: 10px;
   }
   .error {
      color: firebrick;
      padding: 5px 10px;
      background-color: #ffeded;
      border: 2px solid firebrick;
      border-radius: 5px;
      margin: 10px 0 20px 0;
   }
   .buttons {
      display: flex;
      flex-flow: row wrap;
      justify-content: flex-start;
      align-items: center;
      gap: 10px;
   }
   div.actions, .transcribe-acts {
      padding: 10px;
      display: flex;
      flex-flow: row wrap;
      justify-content: space-between;
      align-content: center;
      gap: 20px;
      border-bottom: 1px solid #ddd;
      margin-bottom: 15px;
   }
   .paging {
      display: flex;
      flex-flow: row nowrap;
      align-items: center;
      gap: 10px;
   }
   div.details {
      margin: 0 25px 0 25px;
   }

   .images-wrapper {
      border-top: 1px solid #ccc;
      padding: 10px 0;
      margin-top: 10px;
      display: flex;
      flex-direction: column;
      gap: 10px;
      .image-controls {
         display: flex;
         flex-flow: row wrap;
         .section {
            display: flex;
            flex-flow: row wrap;
            justify-content: flex-start;
            align-items: center;
            gap: 5px;
         }
         .section.wide  {
            gap: 10px;
         }
         button {
            display: inherit !important;
         }
         .rotated {
            transform: scaleX(-1);
         }
      }
      div.image-content {
         flex: 1;
         padding-bottom: 20px;
         border-bottom: 1px solid #ccc;
         display: flex;
         flex-flow: row wrap;
         align-items: stretch;
         justify-content: flex-start;
         gap: 15px;
         div.zoom-wrap {
            display: flex;
            flex-direction: column;
            justify-content: stretch;
            align-items: stretch;
            gap: 15px;
            background: #fafafa;
            border: 1px solid #ddd;
            padding: 10px;
            border-radius: 4px;
            .viewer {
               height: 600px;
               width: 100%;
               background-image: url('@/assets/spinner.gif');
               background-repeat:no-repeat;
               background-position: center center;
            }
         }
      }
   }

   div.transcription-panel {
      flex:1;
      margin: 0;
      border: 1px solid #ddd;
      border-radius: 4px;
      .none {
         font-size: 1.3em;
         text-align: center;
         margin: 15% auto;
      }
      .head {
         background: #fafafa;
         padding: 5px 10px;
         display: flex;
         flex-flow: row wrap;
         justify-content: space-between;
         align-items: center;
         gap: 20px;
         border-bottom: 1px solid #ccc;
         border-radius: 4px 4px 0 0;
      }
   }
}

pre.transcription  {
   margin: 15px;
   padding: 20px;
   border: 1px solid #ddd;
   border-radius: 4px;
   white-space: pre-wrap;       /* Since CSS 2.1 */
   white-space: -moz-pre-wrap;  /* Mozilla, since 1999 */
   white-space: -pre-wrap;      /* Opera 4-6 */
   white-space: -o-pre-wrap;    /* Opera 7 */
   word-wrap: break-word;       /* Internet Explorer 5.5+ */
}

table {
   width: 100%;
   td {
      padding: 2px 0;
   }
   td.label {
      width: 150px;
      text-align: right;
      font-weight: bold;
      color: #444;
      padding: 0 10px 0 0;
      vertical-align: text-top;
      white-space:  nowrap;
   }
}

@media only screen and (min-width: 768px) {
   .image-controls {
      justify-content: space-between;
      align-items: flex-start;
   }
   .zoom-wrap {
      width: 50%;
   }
   div.transcription-panel {
      width: 50%;
   }
   .admin-submission {
      padding: 15px 25px;
   }
}
@media only screen and (max-width: 768px) {
   .image-controls {
      justify-content: center;
      align-items: flex-start;
      gap: 10px 0;
   }
   .zoom-wrap {
      width: 100%;
   }
   div.transcription-panel {
      width: 100%;
   }
   .admin-submission {
      padding: 10px;
   }
}
</style>