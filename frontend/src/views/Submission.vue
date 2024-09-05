<template>
   <div class="submission content">
      <template v-if="system.loading">
         <h1>Loading Details...</h1>
      </template>
      <template v-else-if="system.hasError ">
        <div class="error">{{system.error}}</div>
      </template>
      <template v-else>
          <div class="paging" v-if="!details.isTranscribing">
            <button v-bind:class="{disabled: details.hasPrev == false}" @click="prevClicked" class="prev pure-button pure-button-primary">Prior Submission</button>
            <button v-bind:class="{disabled: details.hasNext == false}" @click="nextClicked" class="next pure-button pure-button-primary">Next Submission</button>
         </div>
         <div class="submit-header">
            <h3><b>BOOK SUBMISSION:</b> {{submission.title}}</h3>
            <div class="submit-time">
               <i class="far fa-clock"></i><span>{{submitDate}}</span>
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
                  <span @click="transcribeClicked(file)" class="ctls pure-button pure-button-primary"
                     :class="{disbled: hasPendingTranscription(file)}">Transcribe</span>
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
         <div class="tags">
            <div @click="tagClicked(tag)" class="tag" v-for="tag in submission.tags">{{ tag }}</div>
         </div>
      </template>
   </div>
</template>

<script setup>
import { onMounted, computed } from 'vue'
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

onMounted(() => {
   details.getSubmission( route.params.id )
})

const transcribeClicked = ((imgFile) => {
   if (this.hasPendingTranscription(imgFile) == false) {
      details.setTranscriptionTarget( imgFile )
   }
})

const hasPendingTranscription = ((imgFile) => {
   return imgFile.transcriptions.length > 0 && imgFile.transcriptions[0].approved == false
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

<style scoped>
h1 {
   font-family: 'Special Elite', cursive;
}
h3 {
   font-family: 'Special Elite', cursive;
   margin: 10px 0;
   font-weight: 500;
   font-size: 1.25em;
}
div.submit-time {
   color: #aaa;
   font-size: 0.85em;
}
div.submit-time span {
   display: inline-block;
   margin-left: 10px;
}
div.details {
   margin-top:20px;
}
div.details div {
   margin-bottom: 3px;
}
div.details .value {
   font-weight: 200;
   color: #444;
}
div.details p {
   margin: 3px 30px;
}
div.details label {
   font-weight: bold;
   margin-right: 5px;
}
div.thumb {
   margin: 0 0 15px 0;
   padding-bottom: 15px;
   border-bottom: 1px solid #ccc;
   display: flex;
   flex-flow: row wrap;
   align-content: flex-start;
}
@media only screen and (min-width: 768px) {
   .zoom-wrap {
      flex-basis: 45%;
      margin-right: 15px;
   }
   div.transcription-wrap {
      flex-basis: 50%;
   }
}
@media only screen and (max-width: 768px) {
   .zoom-wrap {
      flex-basis: 95%;
      margin-bottom: 15px;
   }
   div.transcription-wrap {
      flex-basis: 95%;
   }
}
div.transcription-wrap {
   box-sizing: border-box;
   margin: 0;
   border:1px solid #ccc;
}
.zoom-wrap .ctls.pure-button.pure-button-primary.disbled {
   opacity: 0.2;
   cursor: default;
}
div.thumb:first-of-type {
   border-top: 1px solid #ccc;
   padding-top: 15px;
}
.thumbs {
   margin-top: 20px;
   padding-top: 20px;
}
div.tags {
   clear:both;
   margin-top: 20px;
   padding-top: 20px;
}
div.tag {
   display: inline-block;
   margin: 0 10px 5px 0;
   font-size: 0.85em;
   background: #68c;
   color: white;
   padding: 4px 20px 3px 20px;
   border-radius: 10px;
   text-transform: uppercase;
   font-weight: 500;
}
div.paging {
   font-size: 0.8em;
   text-align: right;
   position: relative;
   top: -10px;
   right: -5px;
}
div.paging .pure-button {
   margin-left: 10px;
   background: #24890d;
}
div.paging .pure-button.disabled {
   margin-left: 10px;
   background: #24890d;
   opacity: 0.5;
   cursor: default;
}
div.tag:hover {
   cursor:pointer;
   background: #79d;
}
.error {
  margin: 5px 0 10px 0;
  color: firebrick;
  font-style: italic;
}
.ctls {
   margin: 5px 0;
   color: white !important;
   width:100%;
}
.toolbar .ctls.pure-button.pure-button-primary.disbled {
   opacity: 0.5;
   cursor: default;
   background: #aaa;
}
.pending {
   font-size: 1.3em;
   text-align: center;
   margin: 15% auto;
}
.transcription-wrap .head {
   font-size: 1.25em;
   background: #dadada;
   padding: 4px 8px;
}
.transcription-wrap .transcription {
   padding: 5px 10px;
}
div.transcription pre {
   white-space: pre-wrap;       /* Since CSS 2.1 */
   white-space: -moz-pre-wrap;  /* Mozilla, since 1999 */
   white-space: -pre-wrap;      /* Opera 4-6 */
   white-space: -o-pre-wrap;    /* Opera 7 */
   word-wrap: break-word;       /* Internet Explorer 5.5+ */
}
</style>