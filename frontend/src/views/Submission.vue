<template>
   <div class="submission content">
      <template v-if="system.loading || details.submission == null">
         <h1>Loading Details...</h1>
      </template>
      <template v-else>
          <div class="paging" v-if="!details.isTranscribing">
            <Button :disabled="!details.hasPrev" @click="prevClicked" label="Prior Submission"/>
            <Button :disabled="!details.hasNext" @click="nextClicked" label="Next Submission"/>
         </div>
         <div class="submit-header">
            <h3><b>BOOK SUBMISSION:</b> {{submission.title}}</h3>
            <div class="submit-time">
               <i class="pi pi-clock"></i><span>{{submitDate}}</span>
            </div>
         </div>
         <div class="details">
            <div><label>Title: </label><span class="value">{{submission.title}}</span></div>
            <div><label>Author: </label><span class="value">{{submission.author}}</span></div>
            <div><label>Publication date:</label><span class="value">{{submission.publication}}</span></div>
            <div><label>Institution: </label><span class="value">{{submission.institution}}</span></div>
            <div><label>Call number: </label><span class="value">{{submission.callNumber}}</span></div>
            <div><label>Submitted by: </label><span class="value">{{submission.submitter}}</span></div>
            <div><label>Description: </label><p class="value" v-html="formatDescription(submission.description)"></p></div>
         </div>
         <div class="images-wrapper">
            <div class="image-controls">
               <div class="section">
                  <Button id="rotate-left" icon="pi pi-undo" rounded severity="secondary"/>
                  <Button id="rotate-right" icon="pi pi-undo" class="rotated" rounded severity="secondary"/>
                  <Button id="zoom-in" icon="pi pi-search-plus" rounded severity="secondary"/>
                  <Button id="zoom-out" icon="pi pi-search-minus" rounded severity="secondary"/>
               </div>
               <div class="section wide" v-if="submission.files.length > 1 && details.isTranscribing == false">
                  <Button icon="pi pi-chevron-left" rounded severity="secondary"
                     :disabled="currImageIdx==0" @click="prevImage"/>
                  <span>Image {{ currImageIdx+1 }} of {{ submission.files.length }}</span>
                  <Button icon="pi pi-chevron-right" rounded severity="secondary"
                     :disabled="currImageIdx==submission.files.length-1" @click="nextImage"/>
               </div>
            </div>
            <div class="image-content">
               <div class="zoom-wrap">
                  <div id="dragon" class="viewer"></div>
                  <Button severity="info" @click="transcribeClicked(file)" label="Transcribe" :disabled="details.isTranscribing"/>
               </div>
               <div class="transcription-wrap">
                  <div class="head">Transcription</div>
                  <div class="trans-form" v-if="details.isTranscribing">
                     <Textarea rows="6" v-model="newTranscription.text" />
                     <div class="row">
                        <label for="submitter">Transcribed By<span class="required">(required)</span></label>
                        <InputText id="submitted" v-model="newTranscription.name" />
                     </div>
                     <div class="row">
                        <label for="email">Email Address<span class="required">(required)</span></label>
                        <InputText id="email" v-model="newTranscription.email"/>
                        <p class="note">We will keep your email address private and will only use if if we need to contact you about your transcription.</p>
                     </div>
                     <div class="controls">
                        <Button @click="cancelTranscribeClicked" severity="secondary" label="Cancel"/>
                        <Button @click="submitTranscribeClicked" label="Submit" :disabled="!canSubmit" />
                     </div>
                  </div>
                  <div v-else class="transcription">
                     <div class="pending" v-if="hasPendingTranscription">
                        Transcription under review.<br/>Please check back in a few days.
                     </div>
                     <pre v-else>{{transcription}}</pre>
                  </div>
               </div>
            </div>
         </div>
         <div class="tags">
            <Button v-for="tag in submission.tags" severity="info" rounded small :label="tag"  @click="tagClicked(tag)" :disabled="details.isTranscribing"/>
         </div>
      </template>
   </div>
</template>

<script setup>
import { onMounted, computed, ref, onUnmounted } from 'vue'
import { useDetailsStore } from "@/stores/details"
import { useSystemStore } from "@/stores/system"
import { useSubmissionsStore } from "@/stores/submissions"
import { useRoute, useRouter } from 'vue-router'
import OpenSeadragon from "openseadragon"
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import { useToast } from 'primevue/usetoast'

const toast = useToast()
const system = useSystemStore()
const submissions = useSubmissionsStore()
const details = useDetailsStore()
const route = useRoute()
const router = useRouter()
const viewer = ref()
const currImageIdx = ref(0)
const newTranscription = ref({text: "", name: "", email: ""})

const submitDate = computed( () => {
   return details.submission.submittedAt.split("T")[0]
})

const submission = computed( () => {
   return details.submission
})

const canSubmit = computed(() => {
   if (details.working) return false
   return newTranscription.value.text != "" && newTranscription.value.name != "" && newTranscription.value.email != ""
})

onMounted( async () => {
   await details.getSubmission( route.params.id )
   initViewer()
})

const initViewer = (() => {
   let sources = []
   details.submission.files.forEach( f => {
      sources.push( {type: 'image',url:  f.url} )
   })
   viewer.value = OpenSeadragon({
      id: `dragon`,
      animationTime: 0.1,
      autoResize: true,
      constrainDuringPan: true,
      imageSmoothingEnabled: true,
      smoothTileEdgesMinZoom: Infinity,
      maxZoomPixelRatio: 10.0,
      showSequenceControl: false,
      zoomInButton:   "zoom-in",
      zoomOutButton:  "zoom-out",
      showNavigator: true,
      showRotationControl: true,
      rotateLeftButton: "rotate-left",
      rotateRightButton: "rotate-right",
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

onUnmounted( async () => {
   if ( viewer.value ) {
      viewer.value.destroy()
   }
})

const nextImage = ( async () => {
   currImageIdx.value++
   viewer.value.goToPage( currImageIdx.value )
})

const prevImage = ( async () => {
   currImageIdx.value--
   viewer.value.goToPage( currImageIdx.value )
})

const transcribeClicked = (() => {
   let imgFile = details.submission.files[currImageIdx.value]
   details.setTranscriptionTarget( imgFile )
})

const hasPendingTranscription = computed(() => {
   let imgFile = details.submission.files[currImageIdx.value]
   if (imgFile.transcriptions.length == 0 ) return false
   let pending = true
   imgFile.transcriptions.forEach( t => {
      if (t.approved === true) {
         pending = false
      }
   })
   return pending
})

const transcription  = computed(() => {
   let imgFile = details.submission.files[currImageIdx.value]
   let t = imgFile.transcriptions.find( t=> t.approved == true)
   if ( t!=null) {
      return t.text
   }
   return ""
})

const tagClicked  = ((tag) => {
   submissions.getTaggedSubmissions( tag )
   router.push("/results")
})

const formatDescription = (( desc ) => {
   let out = desc.replace(/\r|\r\n/gm, '\n').replace(/\n+/gm, "<br/><br/>")
   return out
})

const prevClicked = (async () => {
   if (details.submission.nextId > 0) {
      await details.getSubmission( details.submission.nextId )
      viewer.value.destroy()
      initViewer()
      router.push("/submissions/" + details.submission.nextId )

   }
})

const nextClicked = (async () => {
   if ( details.submission.previousId > 0) {
      await details.getSubmission( details.submission.previousId )
      viewer.value.destroy()
      initViewer()
      router.push("/submissions/" + details.submission.previousId)
   }
})

const cancelTranscribeClicked = (() => {
   details.cancelTranscription()
})

const submitTranscribeClicked = (async () => {
   await details.submitTranscription(newTranscription.value.text, newTranscription.value.name, newTranscription.value.email)
   if (details.transcribeError == "" ) {
      toast.add( {
         severity: 'success', summary: 'Transcription Submitted',
         detail: `The transcription has been submitted. It will appear online after it has been reviewd and approved. Check back later.`
      })
   } else {
      toast.add( {
         severity: 'error', summary: 'Transcription Error',
         detail: `Transcription submission failed: ${details.transcribeError}`
      })
   }
})

</script>

<style scoped lang="scss">
h3 {
   font-family: 'Special Elite', cursive;
   margin: 10px 0;
   font-weight: 500;
   font-size: 1.25em;
}
div.submit-time {
   color: #aaa;
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-start;
   align-items: flex-start;
   gap: 10px;
}
div.paging {
   display: flex;
   flex-flow: row nowrap;
   align-items: flex-start;
   justify-content: flex-end;
   gap: 5px;
}

div.details {
   margin-top:20px;
   padding-bottom: 20px;
   border-bottom: 1px solid #ccc;
   div {
      margin-bottom: 3px;
   }
   .value {
      font-weight: 200;
      color: #444;
   }
   p {
      margin: 3px 30px;
   }
   label {
      font-weight: bold;
      margin-right: 5px;
   }
}

.images-wrapper {
   padding-top: 10px;
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
         }
      }
   }
   div.transcription-wrap {
      flex: 1;
      margin: 0;
      border:1px solid #ddd;
      border-radius: 4px;
      .pending {
         font-size: 1.3em;
         text-align: center;
         margin: 15% auto;
      }
      .head {
         background: #fafafa;
         padding: 5px 10px;
         border-bottom: 1px solid #ccc;
      }
      .transcription {
         padding: 5px 10px;
         pre {
            white-space: pre-wrap;       /* Since CSS 2.1 */
            white-space: -moz-pre-wrap;  /* Mozilla, since 1999 */
            white-space: -pre-wrap;      /* Opera 4-6 */
            white-space: -o-pre-wrap;    /* Opera 7 */
            word-wrap: break-word;       /* Internet Explorer 5.5+ */
         }
      }
      .trans-form {
         padding: 20px;
         display: flex;
         flex-direction: column;
         gap: 20px;
         .row {
            display: flex;
            flex-direction: column;
            gap: 5px;
            .required {
               display: inline-block;
               margin-left: 5px;
               color: #aaa;
            }
         }
         .controls {
            display: flex;
            flex-flow: row nowrap;
            justify-content: flex-end;
            gap: 10px;
         }
      }
   }
}
div.tags {
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;
   align-items: flex-start;
   gap: 5px;
   padding-top: 20px;
}

@media only screen and (min-width: 768px) {
   .zoom-wrap {
      width: 50%;
   }
   div.transcription-wrap {
      width: 50%;
   }
   .image-controls {
      justify-content: space-between;
      align-items: flex-start;
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
   div.transcription-wrap {
      width: 100%;
   }
}
</style>