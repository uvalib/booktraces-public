<template>
   <div class="submission content">
      <template v-if="system.loading">
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
         <Transcribe v-if="details.isTranscribing" />
         <div v-else class="thumbs">
            <div class="thumb" v-for="file in submission.files">
               <div class="zoom-wrap">
                  <vue-image-zoomer :regular="file.url" />
                  <Button severity="info" @click="transcribeClicked(file)" :disabled="hasPendingTranscription(file)" label="Transcribe"/>
               </div>
               <div class="transcription-wrap">
                  <div class="head">Transcription</div>
                  <div class="transcription">
                     <div class="pending" v-if="hasPendingTranscription(file)">
                        Transcription under review.<br/>Please check back in a few days.
                     </div>
                     <pre v-else>{{transcription(file)}}</pre>
                  </div>
               </div>
            </div>
         </div>
         <div class="tags" v-if="!details.isTranscribing">
            <Button v-for="tag in submission.tags" severity="info" rounded small :label="tag"  @click="tagClicked(tag)"/>
         </div>
      </template>
   </div>
</template>

<script setup>
import { onBeforeMount, computed } from 'vue'
import { useDetailsStore } from "@/stores/details"
import { useSystemStore } from "@/stores/system"
import { useSubmissionsStore } from "@/stores/submissions"
import Transcribe from "@/components/Transcribe.vue"
import { useRoute, useRouter } from 'vue-router'

const system = useSystemStore()
const submissions = useSubmissionsStore()
const details = useDetailsStore()
const route = useRoute()
const router = useRouter()

const submitDate = computed( () => {
   return details.submission.submittedAt.split("T")[0]
})

const submission = computed( () => {
   return details.submission
})

onBeforeMount(() => {
   details.getSubmission( route.params.id )
})

const transcribeClicked = ((imgFile) => {
   if ( hasPendingTranscription(imgFile) == false) {
      details.setTranscriptionTarget( imgFile )
   }
})

const hasPendingTranscription = ((imgFile) => {
   if (imgFile.transcriptions.length == 0 ) return false
   let pending = true
   imgFile.transcriptions.forEach( t => {
      if (t.approved === true) {
         pending = false
      }
   })
   return pending
})

const transcription = ((imgFile) => {
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

const prevClicked = (() => {
   if (details.submission.nextId > 0) {
      router.push("/submissions/" + details.submission.nextId )
      details.getSubmission( details.submission.nextId )
   }
})

const nextClicked = (() => {
   if ( details.submission.previousId > 0) {
      router.push("/submissions/" + details.submission.previousId)
      details.getSubmission( details.submission.previousId )
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

.thumbs {
   margin-top: 20px;
   padding-top: 20px;
   display: flex;
   flex-direction: column;
   gap: 20px;
   div.thumb {
      padding-bottom: 20px;
      border-bottom: 1px solid #ccc;
      display: flex;
      flex-flow: row wrap;
      align-items: stretch;
      justify-content: flex-start;
      gap: 15px;
      div.zoom-wrap {
         flex: 1;
         display: flex;
         flex-direction: column;
         justify-content: stretch;
         align-items: stretch;
         gap: 15px;
         background: #fafafa;
         border: 1px solid #ddd;
         padding: 10px;
         border-radius: 4px;
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
      max-width: 50%;
   }
   div.transcription-wrap {
      flex-basis: 50%;
   }
}
@media only screen and (max-width: 768px) {
   .zoom-wrap {
      max-width: 100%;
   }
   div.transcription-wrap {
      max-width: 100%;
   }
}
</style>