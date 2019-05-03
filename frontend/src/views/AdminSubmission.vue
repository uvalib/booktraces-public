<template>
   <div class="admin-submission">
      <h2>Book Traces System Admin Panel <span class="login"><b>Logged in as:</b>{{loginName}}</span></h2>
      <div class="back">
         <router-link to="/admin"><i class="fas fa-arrow-left"></i>&nbsp;Back to Submissions</router-link>
      </div>
      <template v-if="loading">
         <h1>Loading Details...</h1>
      </template>
      <template v-else-if="hasError">
        <div class="error">{{error}}</div>
      </template>
      <template v-else>
         <div class="details">
            <div><label>Title: </label><span class="value">{{details.title}}</span></div>
            <div><label>Author: </label><span class="value">{{details.author}}</span></div>
            <div><label>Publication date:</label><span class="value">{{details.publication}}</span></div>
            <div><label>Library: </label><span class="value">{{details.library}}</span></div>
            <div><label>Call number: </label><span class="value">{{details.callNumber}}</span></div>
            <div><label>Submitted by: </label><span class="value">{{details.submitter}}</span></div>
            <div><label>Submitted on: </label><span class="value">{{submitDate}}</span></div>
            <div><label>Tags: </label><span class="value">{{details.tags.join(", ")}}</span></div>
            <div><label>Description: </label><p class="value">{{details.description}}</p></div>
         </div>
         <div class="thumbs">
            <div class="thumb" v-for="(url,idx) in details.files" :key="idx">
               <img class="thumb" :src="url"/>
            </div>
         </div>
      </template>
   </div>
</template>

<script>
import { mapState } from 'vuex'
import { mapGetters } from 'vuex'
export default {
   name: "admin-submission",
   computed: {
      ...mapState({
         details: state => state.public.details,
         loading: state => state.loading,
         error: state => state.error,
      }),
      ...mapGetters({
         loginName: 'admin/loginName',
      }),
      submitDate() {
         return this.details.submittedAt.split("T")[0]
      },
      hasError() {
         return this.$store.getters.hasError
      }
   },
   methods: {
   },
   created() {
      this.$store.dispatch("public/getSubmissionDetail", this.$route.params.id)
   }
};
</script>

<style scoped>
span.login {
   font-family: sans-serif;
   font-size: 0.5em;
   float: right;
   font-weight: 100;
}
span.login b {
   margin-right: 5px;
}
h2 {
  font-size: 1.5em;
  font-weight: bold;
  border-bottom: 1px dashed #666;
  font-family: 'Special Elite', cursive;
  padding-bottom: 5px;
  margin-bottom: 15px;
}
div.admin-submission {
   padding: 15px 25px;
   min-height: 600px;
   background: white;
   color: #444;
   position: relative;
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
   border-top: 1px solid #ccc;
   padding-top: 20px;
}
div.back {
   position: relative;
}
a {
   position: absolute;
   right:0;
   top: 10px;
   font-size: 0.8em;
   
}
.error {
  margin: 5px 0 10px 0;
  color: firebrick;
  font-style: italic;
}
</style>