<template>
   <div class="submission content">
      <template v-if="loading">
         <h1>Loading Details...</h1>
      </template>
      <template v-else-if="hasError">
        <div class="error">{{error}}</div>
      </template>
      <template v-else>
         <div class="paging">
            <button v-bind:class="{disabled: hasPrev == false}" @click="prevClicked" class="prev pure-button pure-button-primary">Prior Submission</button>
            <button v-bind:class="{disabled: hasNext == false}" @click="nextClicked" class="next pure-button pure-button-primary">Next Submission</button>
         </div>
         <div class="submit-header">
            <h3><b>BOOK SUBMISSION:</b> {{details.title}}</h3>
            <div class="submit-time">
               <i class="far fa-clock"></i><span>{{submitDate}}</span>
            </div>
         </div>
         <div class="details">
            <div><label>Title: </label><span class="value">{{details.title}}</span></div>
            <div><label>Author: </label><span class="value">{{details.author}}</span></div>
            <div><label>Publication date:</label><span class="value">{{details.publication}}</span></div>
            <div><label>Library: </label><span class="value">{{details.library}}</span></div>
            <div><label>Call number: </label><span class="value">{{details.callNumber}}</span></div>
            <div><label>Submitted by: </label><span class="value">{{details.submitter}}</span></div>
            <div><label>Description: </label><p class="value" v-html="formatDescription(details.description)"></p></div>
         </div>
         <div class="thumbs">
            <div class="thumb" v-for="(url,idx) in details.files" :key="idx">
               <a :href="url" target="_blank"><img class="thumb" :src="url"/></a>
            </div>
         </div>
         <div class="tags">
            <div @click="tagClicked" class="tag" v-for="(tag,idx) in details.tags" :key="idx">
               {{tag}}
            </div>
         </div>
      </template>
   </div>
</template>

<script>
import { mapState } from 'vuex'
export default {
   name: "submission",
   computed: {
      ...mapState({
         details: state => state.submissionDetail,
         loading: state => state.loading,
         error: state => state.error,
      }),
      submitDate() {
         return this.details.submittedAt.split("T")[0]
      },
      hasError() {
         return this.$store.getters.hasError
      },
      hasPrev() {
         return this.details.nextId > 0
      },
      hasNext() {
         return this.details.previousId > 0
      }

   },
   methods: {
      tagClicked(event) {
         let t = event.currentTarget.textContent.replace(/^\s+|\s+$/g, '')
         this.$store.dispatch("public/getTaggedSubmissions", t)
         this.$router.push("/results")
      },
      formatDescription( desc ) {
         let out = desc.replace(/\r|\r\n/gm, '\n').replace(/\n+/gm, "<br/><br/>")
         return out
      },
      prevClicked() {
         if (this.details.nextId > 0) {
            this.$router.push("/submissions/" + this.details.nextId)
            this.$store.dispatch("getSubmissionDetail", this.details.nextId)
         }
      },
      nextClicked() {
         if (this.details.previousId > 0) {
            this.$router.push("/submissions/" + this.details.previousId)
            this.$store.dispatch("getSubmissionDetail", this.details.previousId)
         }
      }
   },
   created: function() {
      this.$store.dispatch("getSubmissionDetail", this.$route.params.id)
   }
};
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
img.thumb {
   max-width: 250px;
   max-height: 250px;
}
div.thumb {
   /* float: left; */
   display: inline-block;
   margin: 5px 10px;
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
</style>