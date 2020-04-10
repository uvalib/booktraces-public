<template>
   <div class="transcribe">
      <div class="pure-form data">
         <div class="image">
            <pinch-zoom v-bind:limitZoom="200">
               <img :src="transcribeFile.url">
            </pinch-zoom>
         </div>
         <textarea v-model="transcription"></textarea>
      </div>
      <div class="controls">
         <span @click="cancelClicked" class="cancel pure-button pure-button-primary">Cancel</span>
         <span @click="submitClicked" class="pure-button pure-button-primary">Submit</span>
      </div>
   </div>
</template>

<script>
import { mapState } from 'vuex'

export default {
   computed: {
      ...mapState({
         details: state => state.submissionDetail,
         loading: state => state.loading,
         error: state => state.error,
         transcribeFile: state => state.public.transcribeFile
      }),
   },
   data: function () {
      return {
         transcription: ""
      }
   },
   methods: {
      cancelClicked() {
         this.$store.commit("public/cancelTranscribe", false)
      },
      submitClicked() {
         this.$store.commit("public/cancelTranscribe", false)
      }
   }
}
</script>

<style scoped>
div.data {
   display: flex;
   flex-flow: row wrap;
}
div.data textarea {
   flex-grow: 1;
   margin-left:10px;
}
.image {
   max-width: 50%;
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
.controls .pure-button.pure-button-primary.cancel  {
   background-color: #999;
}
</style>
