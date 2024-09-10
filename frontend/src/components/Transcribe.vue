<template>
   <div class="transcribe">
      <div id="view-ctls">
         <Button id="rotate-left" icon="pi pi-undo" rounded size="large" severity="secondary"/>
         <Button id="rotate-right" icon="pi pi-undo" class="rotated" rounded size="large" severity="secondary"/>
         <Button id="zoom-in" icon="pi pi-search-plus" rounded size="large" severity="secondary"/>
         <Button id="zoom-out" icon="pi pi-search-minus" rounded size="large" severity="secondary"/>
      </div>
      <div id="dragon"></div>
      <textarea rows="5" v-model="transcription"></textarea>
      <div class="row">
         <label for="submitter">Transcribed By<span class="required">(required)</span></label>
         <input v-model="name"  class="pure-u-1-1" type="text">
      </div>
      <div class="row">
         <label for="email">Email Address<span class="required">(required)</span></label>
         <input v-model="email" class="pure-u-1-1" type="email">
         <p class="note">We will keep your email address private and will only use if if we need to contact you about your transcription.</p>
      </div>
      <p class="error" v-if="details.transcribeError">{{details.transcribeError}}</p>
      <div class="controls">
         <Button @click="cancelClicked" severity="secondary" label="Cancel"/>
         <Button @click="submitClicked" label="Submit" :disabled=" details.working" />
      </div>
   </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useDetailsStore } from "@/stores/details"
import OpenSeadragon from "openseadragon"

const details = useDetailsStore()

const transcription = ref("")
const name = ref("")
const email = ref("")
const viewer = ref()

onMounted(() => {
   viewer.value = OpenSeadragon({
      id: "dragon",
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
      tileSources: {
         type: 'image',
         url:  details.transcribeFile.url
      },
   });
})

const cancelClicked = (() => {
   details.cancelTranscription()
})

const submitClicked = (() => {
   details.submitTranscription(transcription.value, name.value, email.value)
})
</script>
<style scoped lang="scss">
.transcribe {
   margin-top: 20px;
   display: flex;
   flex-direction: column;
   gap: 15px;
   #view-ctls {
      margin-top:10px;
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-end;
      align-items: center;
      gap: 5px;
      button {
         display: inherit !important;
      }
      .rotated {
         transform: scaleX(-1);
      }
   }
   #dragon {
      height: 450px;
      background: #fafafa;
      border: 1px solid #ddd;
      border-radius: 5px;
   }
   textarea, input {
      border: 1px solid #ddd;
      border-radius: 4px;
      padding: 10px;
   }
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
   .note {
      padding: 0;
      margin: 0;
      color: #aaa;
   }
   .error {
      font-style: italic;
      color: firebrick;
   }
   .controls {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-end;
      gap: 10px;
   }
}
</style>
