<template>
   <div class="transcribe">
      <div class="pure-form transcription">
         <div class="image">
            <UVAViewer :url="transcribeFile.url" :height="400"/>
         </div>
         <textarea rows="5" v-model="transcription"></textarea>
      </div>
      <div class="pure-form">
         <div class="pure-u-1-1 gap">
          <label for="submitter">Transcribed By<span class="required">(required)</span></label>
          <input v-model="name"  class="pure-u-1-1" type="text">
        </div>
        <div class="pure-u-1-1 gap">
          <label for="email">Email Address<span class="required">(required)</span></label>
          <input v-model="email" class="pure-u-1-1" type="email">
          <p class="note">We will keep your email address private and will only use if if we need to contact you about your transcription.</p>
        </div>
      </div>
      <p class="error">{{transcribeError}}</p>
      <div class="controls">
         <span @click="cancelClicked" class="cancel pure-button pure-button-primary">Cancel</span>
         <span @click="submitClicked" :class="{disabled: submitting}" class="pure-button pure-button-primary">Submit</span>
      </div>
   </div>
</template>

<script>
import { mapState } from 'vuex'
import UVAViewer from "@/components/UVAViewer"
export default {
   components: {
      UVAViewer
   },
   computed: {
      ...mapState({
         submitting: state => state.transcribe.submitting,
         transcribeFile: state => state.transcribe.file,
         transcribeError: state => state.transcribe.error,
      }),
   },
   data: function () {
      return {
         transcription: "",
         name: "",
         email: "",
         zoomProp: {
            zoomControlScale: 10
         }
      }
   },
   methods: {
      cancelClicked() {
         this.$store.commit("transcribe/cancel")
      },
      submitClicked() {
         this.$store.dispatch("transcribe/submit", {transcription: this.transcription,
            name: this.name, email: this.email})
      }
   }
}
</script>

<style scoped>
.error {
   font-style: italic;
   color: firebrick;
}
div.transcription {
   display: flex;
   flex-direction: column;
}
.transcription .image {
   overflow: hidden;
   margin-bottom: 10px;
}
.transcribe {
   padding: 15px;
}
.controls {
   margin-top: 10px;
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-end;
}
.controls .pure-button.pure-button-primary {
   margin-left: 10px;
}
.controls .pure-button.pure-button-primary.disabled {
   cursor: default;
   opacity: 0.2;
}
.controls .pure-button.pure-button-primary.cancel  {
   background-color: #999;
}
.gap {
   margin-top: 10px;
}
.required {
   font-size: 0.8em;
   font-weight: 100;
   margin-left: 10px;
}
p.note {
   padding:0;
   margin:2px 15px;
   font-size: 0.8em;
   font-weight: 100;
}
</style>
