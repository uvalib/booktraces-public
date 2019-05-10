<template>
   <div class="admin-submission">
      <h2>System Admin Panel<span class="login"><b>Logged in as:</b>{{loginName}}</span></h2>
      <template v-if="loading">
         <h1>Loading Details...</h1>
      </template>
      <template v-else>
         <div class="actions">
            <router-link to="/admin"><i class="fas fa-arrow-left"></i>&nbsp;Back to Submissions</router-link>
            <div v-if="!edit" class="buttons">
               <button @click="editClicked" class="admin pure-button edit pure-button-primary">Edit</button>
               <button v-if="!details.published" @click="publishClicked" class="admin pure-button publish pure-button-primary">Publish</button>
               <button v-else @click="unpublishClicked" class="admin pure-button unpublish pure-button-primary">Unpublish</button>
               <button @click="deleteClicked" class="admin pure-button delete pure-button-primary">Delete</button>
            </div>
            <div v-else class="buttons">
               <button @click="saveClicked" class="admin pure-button pure-button-primary">Save</button>
               <button @click="cancelClicked" class="admin pure-button pure-button-primary">Cancel</button>
            </div>
         </div>
         <div class="error">{{error}}</div>
         <div v-if="!edit" class="details">
            <table>
               <tr><td class="label">Visible to public:</td><td class="value">{{published}}</td></tr>
               <tr><td class="label">Title:</td><td class="value">{{details.title}}</td></tr>
               <tr><td class="label">Author:</td><td class="value">{{details.author}}</td></tr>
               <tr><td class="label">Publication details:</td><td class="value">{{details.publication}}</td></tr>
               <tr><td class="label">Library:</td><td class="value">{{details.library}}</td></tr>
               <tr><td class="label">Call number:</td><td class="value">{{details.callNumber}}</td></tr>
               <tr><td class="label">Submitted by:</td><td class="value">{{details.submitter}}</td></tr>
               <tr><td class="label">Submitted on:</td><td class="value">{{submitDate}}</td></tr>
               <tr><td class="label">Tags:</td><td class="value">{{details.tags.join(", ")}}</td></tr>
               <tr><td class="label">Description:</td><td class="value">{{details.description}}</td></tr>
            </table>
         </div>
         <div v-else class="edit details">
            <table>
               <tr>
                  <td class="label">Visible to public:</td>
                  <td class="value">{{published}}</td>
               </tr>
               <tr>
                  <td class="label">Title:</td>
                  <td class="value"><input id="title" type="text" v-model="editDetails.title"></td>
               </tr>
               <tr>
                  <td class="label">Author:</td>
                  <td class="value"><input id="author" type="text" v-model="editDetails.author"></td>
               </tr>
               <tr>
                  <td class="label">Publication details:</td>
                  <td class="value"><input id="publication" type="text" v-model="editDetails.publication"></td>
               </tr>
               <tr>
                  <td class="label">Library:</td>
                  <td class="value"><input id="library" type="text" v-model="editDetails.library"></td>
               </tr>
               <tr>
                  <td class="label">Call number:</td>
                  <td class="value"><input id="callNumber" type="text" v-model="editDetails.callNumber"></td>
               </tr>
               <tr>
                  <td class="label">Submitted by:</td>
                  <td class="value"><input id="submitter" type="text" v-model="editDetails.submitter"></td>
               </tr>
               <tr>
                  <td class="label">Submitted on:</td>
                  <td class="value">{{submitDate}}</td>
               </tr>
               <tr>
                  <td class="label">Tags:</td>
                  <td class="value" style="position:relative;">
                     <div v-if="showTagList" class="source-tags">
                        <p class="head">Available Tags</p>
                        <div class="list">
                           <p class="tag" v-for="(tag,idx) in tags" :key="idx" @click="addTag" :data-id="tag.id">{{tag.name}}</p>
                        </div>
                        <span @click="closeTagList" class="add-tag pure-button pure-button-primary">Done</span>
                     </div>
                     <span @click="addTagClicked" class="add-tag pure-button pure-button-primary">Add</span>
                     <span class="tag" v-for="(tag,idx) in editDetails.tags" :key="idx" @click="removeTag">
                        {{tag}} 
                        <i class="fas fa-times-circle"></i>
                     </span>   
                  </td>
               </tr>
               <tr>
                  <td class="label">Description:</td>
                  <td class="value"><textarea rows="5" id="description" v-model="editDetails.description"></textarea></td>
               </tr>
            </table>
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
   data: function () {
      return {
         edit: false,
         editDetails: null,
         showTagList: false,
      }
   },
   computed: {
      ...mapState({
         details: state => state.submissionDetail,
         loading: state => state.loading,
         error: state => state.error,
         tags: state => state.tags,
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
      addTag(event) {
         let addTag = event.currentTarget.textContent.replace(/^\s+|\s+$/g, '')
         if ( this.editDetails.tags.includes(addTag)) {
            return
         }
         this.editDetails.tags.push(addTag)
      },
      closeTagList() {
         this.showTagList = false
      },
      removeTag(event) {
         let delTag = event.currentTarget.textContent.replace(/^\s+|\s+$/g, '')
         var delIdx = -1
         this.editDetails.tags.some( function(tag,idx) {
            if (tag == delTag) {
               delIdx = idx
               return true
            }
            return false
         })
         this.editDetails.tags.splice(delIdx, 1)
      },
      addTagClicked() {
         this.showTagList = true
      },
      deleteClicked() {
         let resp = confirm("Delete this submission? All data and unloaded files will be permanently lost. Are you sure?")
         if (resp) {
            this.$store.dispatch("admin/deleteSubmission", {id:this.details.id, backToIndex:true})
         }
      },
      editClicked() {
         this.editDetails = Object.assign({},this.details)
         this.edit=true
      },
      saveClicked() {
         this.$store.dispatch('admin/updateSubmission',this.editDetails).then((/*response*/) => {
            this.edit=false
            this.editDetails = null
        })
      },
      cancelClicked() {
         this.edit=false
      },
      publishClicked() {
         this.$store.dispatch("admin/updatePublicationStatus", {id:this.details.id, public: true})
      },
      unpublishClicked() {
         this.$store.dispatch("admin/updatePublicationStatus", {id:this.details.id, public: false})
      }
   },
   created() {
      this.$store.dispatch("getSubmissionDetail", this.$route.params.id)
      this.$store.dispatch('getTags')
   }
};
</script>

<style scoped>
.source-tags {
   position: absolute;
   left: -100px;
   top: -100px;
   width:175px;
   font-size: 0.9em;
   border: 1px solid #cccc;
   display: inline-block;
   box-shadow: 2px 2px 10px #aaa;
   border-radius: 10px;
   text-align: center;
   background: #efefef;
   padding-bottom: 10px;
}
.source-tags p.head {
   padding: 4px 20px;
   font-weight: bold;
   margin: 0;
   background: #888;
   border-radius: 10px 10px 0 0;
   text-align: center;
   color: white;
}
.source-tags p {
   margin: 2px 0;
}
.source-tags p:hover {
   cursor: pointer;
   text-decoration: underline;
}
.source-tags .list {
   margin: 10px;
   max-height: 150px;
   overflow: scroll;
   text-align: left;
   border: 1px solid #ccc;
   padding: 5px;
   border-radius: 5px;
   background: white;
}
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
   vertical-align: text-top;
}
img.thumb {
   max-width: 250px;
   max-height: 250px;
}
div.thumb {
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
   top: 8px;
}
.error {
  margin: 5px 0 10px 0;
  color: firebrick;
  font-style: italic;
}
table {
   width:100%;
}
td {
   padding: 2px 0
}
td.label {
   width: 150px;
   text-align: right;
   font-weight: bold;
   color: #444;
   padding: 0 10px 0 0;
   vertical-align: text-top;
}
td input {
   border:1px solid #ccc;
   width: 90%;
   outline:none;  
}
td textarea {
   display: block;
   border: 1px solid #ccc;
   width:  90%;
}
span.add-tag.pure-button {
   margin-right: 10px;
   padding: 1px 20px 0px 20px;
   font-size: 0.9em;
}
span.tag {
   cursor: pointer;
   padding: 1px 4px 0px 12px;
   display:inline-block;
   border: 1px solid #ccc;
   border-radius: 20px;
   margin-right: 5px;
   font-size: 0.9em;
}
i.fas.fa-times-circle {
   color: firebrick;
   margin-left: 5px;
}
</style>