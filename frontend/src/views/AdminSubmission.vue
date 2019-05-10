<template>
   <div class="admin-submission">
      <h2>Book Traces System Admin Panel <span class="login"><b>Logged in as:</b>{{loginName}}</span></h2>
      <template v-if="loading">
         <h1>Loading Details...</h1>
      </template>
      <template v-else>
         <div class="actions">
            <router-link to="/admin"><i class="fas fa-arrow-left"></i>&nbsp;Back to Submissions</router-link>
            <div class="buttons">
               <button @click="editClicked" class="admin pure-button edit pure-button-primary">Edit</button>
               <button v-if="!details.published" @click="publishClicked" class="admin pure-button publish pure-button-primary">Publish</button>
               <button @click="deleteClicked" class="admin pure-button delete pure-button-primary">Delete</button>
            </div>
         </div>
         <div class="error">{{error}}</div>
         <div class="details">
            <div><label>Visible to public: </label><span class="value">{{published}}</span></div>
            <div><label>Title: </label><span class="value">{{details.title}}</span></div>
            <div><label>Author: </label><span class="value">{{details.author}}</span></div>
            <div><label>Publication details:</label><span class="value">{{details.publication}}</span></div>
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
         details: state => state.submissionDetail,
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
      },
      published() {
         if (this.details.published) {
            return "YES"
         } else {
            return "NO"
         }
      }
   },
   methods: {
      deleteClicked() {
         let resp = confirm("Delete this submission? All data and unloaded files will be permanently lost. Are you sure?")
         if (resp) {
            this.$store.dispatch("admin/deleteSubmission", {id:this.details.id, backToIndex:true})
         }
      },
      editClicked() {
      },
      publishClicked() {
         this.$store.dispatch("admin/publishSubmission", this.details.id)
      }
   },
   created() {
      this.$store.dispatch("getSubmissionDetail", this.$route.params.id)
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
div.actions {
   position: relative;
   padding: 5px 0;
   border-bottom: 1px dashed;
   background-color: #fafafa;
   font-size: 0.9em;
}
div.actions button.admin.pure-button {
   padding: 2px 10px;
   margin: 0 4px;
   opacity: 0.6;
}
div.actions button.admin.pure-button:hover {
   opacity: 1;
}
div.actions button.admin.pure-button.delete {
   background-color: firebrick;
}
div.buttons {
   position: absolute;
   right: 0;
   top: 3px;
}
.error {
  margin: 5px 0 10px 0;
  color: firebrick;
  font-style: italic;
}
</style>