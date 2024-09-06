<template>
   <!-- <div class="admin-submission">
      <h2>
         System Admin Panel
         <span class="login">
            <b>Logged in as:</b>
            {{loginName}}
         </span>
      </h2>
      <template v-if="loading">
         <h1>Loading Details...</h1>
      </template>
      <template v-else>
         <div class="actions">
            <router-link to="/admin">
               <i class="fas fa-arrow-left"></i>&nbsp;Back to Submissions
            </router-link>
            <div v-if="!edit" class="buttons">
               <button @click="editClicked" class="admin pure-button edit pure-button-primary">Edit</button>
               <button
                  v-if="!details.published"
                  @click="publishClicked"
                  class="admin pure-button publish pure-button-primary"
               >Publish</button>
               <button
                  v-else
                  @click="unpublishClicked"
                  class="admin pure-button unpublish pure-button-primary"
               >Unpublish</button>
               <button
                  @click="deleteClicked"
                  class="admin pure-button delete pure-button-primary"
               >Delete</button>
            </div>
            <div v-else class="buttons">
               <button @click="saveClicked" class="admin pure-button pure-button-primary">Save</button>
               <button @click="cancelClicked" class="admin pure-button pure-button-primary">Cancel</button>
            </div>
         </div>
         <div class="error">{{error}}</div>
         <div v-if="!edit" class="details">
            <table>
               <tr>
                  <td class="label">Visible to public:</td>
                  <td class="value">{{published}}</td>
               </tr>
               <tr>
                  <td class="label">Title:</td>
                  <td class="value">{{details.title}}</td>
               </tr>
               <tr>
                  <td class="label">Author:</td>
                  <td class="value">{{details.author}}</td>
               </tr>
               <tr>
                  <td class="label">Publication details:</td>
                  <td class="value">{{details.publication}}</td>
               </tr>
               <tr>
                  <td class="label">Institution:</td>
                  <td class="value">{{details.institution}}</td>
               </tr>
               <tr>
                  <td class="label">Call number:</td>
                  <td class="value">{{details.callNumber}}</td>
               </tr>
               <tr>
                  <td class="label">Submitted by:</td>
                  <td class="value">{{details.submitter}} ( {{details.email}} )</td>
               </tr>
               <tr>
                  <td class="label">Submitted on:</td>
                  <td class="value">{{submitDate}}</td>
               </tr>
               <tr>
                  <td class="label">Tags:</td>
                  <td class="value">{{submissionTagCSV}}</td>
               </tr>
               <tr>
                  <td class="label">Description:</td>
                  <td class="value">
                     <span v-html="formatDescription(details.description)"></span>
                  </td>
               </tr>
            </table>
         </div>
         <div v-else class="edit details  pure-form">
            <table style="width:75%; margin: 0 auto;">
               <tr>
                  <td class="label">Visible to public:</td>
                  <td class="value">{{published}}</td>
               </tr>
               <tr>
                  <td class="label">Title:</td>
                  <td class="value">
                     <input id="title" type="text" v-model="editDetails.title" />
                  </td>
               </tr>
               <tr>
                  <td class="label">Author:</td>
                  <td class="value">
                     <input id="author" type="text" v-model="editDetails.author" />
                  </td>
               </tr>
               <tr>
                  <td class="label">Publication details:</td>
                  <td class="value">
                     <input id="publication" type="text" v-model="editDetails.publication" />
                  </td>
               </tr>
               <tr>
                  <td class="label">Institution:</td>
                  <td class="value">
                     <multiselect v-model="selectedInstitution" class="folders"
                        placeholder="Select or create an institution"
                        :showLabels="false"
                        :searchable="true"
                        :taggable="true"
                        track-by="id" label="name"
                        tagPlaceholder="Press enter to create a new institution"
                        @tag="addInstitution"
                        :options="institutions">
                     </multiselect>
                  </td>
               </tr>
               <tr>
                  <td class="label">Call number:</td>
                  <td class="value">
                     <input id="callNumber" type="text" v-model="editDetails.callNumber" />
                  </td>
               </tr>
               <tr>
                  <td class="label">Submitted by:</td>
                  <td class="value">
                     <input id="submitter" type="text" v-model="editDetails.submitter" />
                  </td>
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
                           <p class="tag" v-for="(tag,idx) in tags" :key="idx" @click="addTag" :data-id="tag.id">
                              {{tag.name}}
                           </p>
                        </div>
                        <span @click="closeTagList" class="add-tag pure-button pure-button-primary">
                           Done
                        </span>
                     </div>
                     <span @click="addTagClicked" class="add-tag pure-button pure-button-primary">
                        Add
                     </span>
                     <span class="tag" v-for="(tag,idx) in editDetails.tags"
                        :key="idx" @click="removeTag">
                        {{tag}}
                        <i class="fas fa-times-circle"></i>
                     </span>
                  </td>
               </tr>
               <tr>
                  <td class="label">Description:</td>
                  <td class="value">
                     <textarea rows="5" id="description" v-model="editDetails.description"></textarea>
                  </td>
               </tr>
            </table>
         </div>
         <div class="thumbs" v-if="edit == false">
            <div class="thumb" v-for="(f,idx) in details.files" :key="idx">
               <div class="zoom-wrap">
                  <pinch-zoom v-bind:limitZoom="200">
                     <img class="thumb" :src="`${f.url}?v=${Math.floor(Math.random() * 1000)}`" />
                  </pinch-zoom>
                  <p @click="rotateClicked(f.url)" class="pure-button rotate">Rotate Right</p>
               </div>
               <div class="transcriptions">
                  <div class="transcription-title">
                     <span class="head">Transcriptions</span>
                     <span class="paging" v-if="f.transcriptions.length == 0">None</span>
                     <span v-else class="paging">
                        <i  @click="priorTran(f)" class="paging fas fa-chevron-left" :class="{disabled: transcriptionIdx == 0}"></i>
                        <span>{{transcriptionIdx+1}} of {{f.transcriptions.length}}</span>
                        <i @click="nextTran(f)" class="paging fas fa-chevron-right" :class="{disabled: transcriptionIdx == f.transcriptions.length-1}"></i>
                     </span>
                  </div>
                  <div class="transcription-info" v-if="f.transcriptions.length > 0">
                     <div>
                        <label>Date:</label>
                        <span class="data">{{getTranscribeDate(f)}}</span>
                     </div>
                     <div>
                        <label>Submitter:</label>
                        <span class="data">{{getTranscriber(f)}}</span>
                     </div>
                     <template v-if="editTrans">
                        <textarea class="edit-trans" rows="3" v-model="workingTrans"></textarea>
                        <div v-if="error" class="error">{{error}}</div>
                        <div class="actions">
                           <span class="buttons">
                              <span @click="cancelEdit" class="pure-button trans">Cancel</span>
                              <span @click="submitEdit(f)" class="pure-button trans">Submit</span>
                           </span>
                        </div>
                     </template>
                     <template v-else>
                        <pre class="transcription">{{f.transcriptions[transcriptionIdx].text}}</pre>
                        <div class="actions">
                           <span class="status">{{getTranscribeStatus(f)}}</span>
                           <span class="buttons">
                              <span @click="deleteTranscription(f)" class="pure-button trans">Delete</span>
                              <span @click="editTranscription(f)" class="pure-button trans">Edit</span>
                              <span v-if="f.transcriptions[transcriptionIdx].approved==false" @click="approveTranscription(f)"
                                 class="pure-button trans">Approve</span>
                           </span>
                        </div>
                     </template>
                  </div>
               </div>
            </div>
         </div>
      </template>
   </div> -->
</template>

<script>
// import { mapState } from "vuex";
// import { mapGetters } from "vuex";
// import Multiselect from "vue-multiselect";
// export default {
//    name: "admin-submission",
//    components: {
//       Multiselect
//    },
//    data: function() {
//       return {
//          editTrans: false,
//          workingTrans: "",
//          edit: false,
//          editDetails: null,
//          showTagList: false,
//          selectedInstitution: null,
//          transcriptionIdx: 0
//       };
//    },
//    computed: {
//       ...mapState({
//          details: state => state.submissionDetail,
//          loading: state => state.loading,
//          error: state => state.error,
//          tags: state => state.tags,
//          institutions: state => state.institutions
//       }),
//       ...mapGetters({
//          loginName: "admin/loginName"
//       }),
//       submitDate() {
//          return this.details.submittedAt.split("T")[0];
//       },
//       hasError() {
//          return this.$store.getters.hasError;
//       },
//       published() {
//          if (this.details.published) {
//             return "YES";
//          } else {
//             return "NO";
//          }
//       },
//       submissionTagCSV() {
//          if (this.details.tags) {
//             return this.details.tags.join(", ");
//          }
//          return "";
//       }
//    },
//    methods: {
//       cancelEdit() {
//          this.editTrans = false
//          this.workingTrans = ""
//       },
//       async submitEdit(f) {
//          let data = {submissionID: this.details.id,
//             fileID: f.id,
//             transcriptionID: f.transcriptions[this.transcriptionIdx].id,
//             transcription: this.workingTrans}
//          await this.$store.dispatch("transcribe/update", data)
//          if (this.error == "" || this.error == null) {
//             this.editTrans = false
//             this.workingTrans = ""
//          }
//       },
//       approveTranscription(f) {
//          let t = f.transcriptions[this.transcriptionIdx]
//          this.$store.dispatch("transcribe/approve", t.id)
//       },
//       editTranscription(f) {
//          this.editTrans = true
//          this.workingTrans = f.transcriptions[this.transcriptionIdx].text
//       },
//       deleteTranscription(f) {
//          let resp = confirm("Are you sure you want to delete this transcrption?")
//          if (resp) {
//             let t = f.transcriptions[this.transcriptionIdx]
//             this.$store.dispatch("transcribe/delete", t.id)
//             this.transcriptionIdx = 0
//          }
//       },
//       getTranscribeStatus(f) {
//          if (f.transcriptions.length == 0) return ""
//          let t = f.transcriptions[this.transcriptionIdx]
//          if (t.approved) {
//             return "Approved"
//          }
//          return "Pending"
//       },
//       getTranscribeDate(f) {
//          if (f.transcriptions.length == 0) return ""
//          let t = f.transcriptions[this.transcriptionIdx]
//          return t.transcribed_at.split("T")[0]
//       },
//       getTranscriber(f) {
//          if (f.transcriptions.length == 0) return ""
//          let t = f.transcriptions[this.transcriptionIdx]
//          return `${t.transcriber_email} (${t.transcriber})`
//       },
//       nextTran(f) {
//          if (this.transcriptionIdx == f.transcriptions.length -1) {
//             return
//          }
//          this.transcriptionIdx++
//       },
//       priorTran(_f) {
//          if (this.transcriptionIdx == 0) {
//             return
//          }
//          this.transcriptionIdx--
//       },
//       addInstitution(newInstitutionName) {
//          this.$store
//             .dispatch("addInstitution", newInstitutionName)
//             .then(() => {
//                this.institutions.some(i => {
//                   if (i.name == newInstitutionName) {
//                      this.selectedInstitution = i
//                      return true;
//                   }
//                   return false;
//                })
//             })
//             .catch(error => {
//                // TODO something else maybe?
//                alert(error)
//             });
//       },
//       formatDescription(desc) {
//          let out = desc
//             .replace(/\r|\r\n/gm, "\n")
//             .replace(/\n+/gm, "<br/><br/>");
//          return out;
//       },
//       addTag(event) {
//          let addTag = event.currentTarget.textContent.replace(/^\s+|\s+$/g, "");
//          if (this.editDetails.tags.includes(addTag)) {
//             return;
//          }
//          this.editDetails.tags.push(addTag);
//       },
//       closeTagList() {
//          this.showTagList = false;
//       },
//       removeTag(event) {
//          let delTag = event.currentTarget.textContent.replace(/^\s+|\s+$/g, "");
//          var delIdx = -1;
//          this.editDetails.tags.some(function(tag, idx) {
//             if (tag == delTag) {
//                delIdx = idx;
//                return true;
//             }
//             return false;
//          });
//          this.editDetails.tags.splice(delIdx, 1);
//       },
//       rotateClicked(imgURL) {
//          this.$store.dispatch("admin/rotateImage", {
//             submissionID: this.details.id,
//             imgURL: imgURL
//          });
//       },
//       addTagClicked() {
//          this.showTagList = true;
//       },
//       deleteClicked() {
//          let resp = confirm(
//             "Delete this submission? All data and unloaded files will be permanently lost. Are you sure?"
//          );
//          if (resp) {
//             this.$store.dispatch("admin/deleteSubmission", {
//                id: this.details.id,
//                backToIndex: true
//             });
//          }
//       },
//       editClicked() {
//          this.editDetails = Object.assign({}, this.details);
//          this.edit = true;
//          this.institutions.some(i => {
//             if (i.name == this.editDetails.institution) {
//                this.selectedInstitution = i
//                return true
//             }
//             return false
//          })
//       },
//       saveClicked() {
//          this.editDetails.institution_id = this.selectedInstitution.id
//          this.editDetails.institution = this.selectedInstitution.name
//          this.$store
//             .dispatch("admin/updateSubmission", this.editDetails)
//             .then((/*response*/) => {
//                this.edit = false;
//                this.editDetails = null;
//             });
//       },
//       cancelClicked() {
//          this.edit = false;
//       },
//       publishClicked() {
//          this.$store.dispatch("admin/updatePublicationStatus", {
//             id: this.details.id,
//             public: true
//          });
//       },
//       unpublishClicked() {
//          this.$store.dispatch("admin/updatePublicationStatus", {
//             id: this.details.id,
//             public: false
//          });
//       }
//    },
//    created() {
//       this.$store.dispatch("getSubmissionDetail", this.$route.params.id)
//       this.$store.dispatch("getTags")
//       this.$store.dispatch("getInstitutions")
//    }
// };
</script>

<style scoped>
.edit-trans {
   box-sizing: border-box;
   width: 100%;
   border-color: #ccc;
}
.source-tags {
   position: absolute;
   left: -100px;
   top: -100px;
   width: 175px;
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
   font-family: "Special Elite", cursive;
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
   margin: 20px 25px 0 25px;
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
   display: block;
}
div.thumb {
   display: block;
   margin: 5px 10px;
   display: flex;
   flex-flow: row wrap;
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
   width: 100%;
}
td {
   padding: 2px 0;
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
   border: 1px solid #ccc;
   width: 100%;
   outline: none;
   box-sizing: border-box;
}
td textarea {
   display: block;
   border: 1px solid #ccc;
   width: 90%;
}
span.add-tag.pure-button {
   margin-right: 10px;
   padding: 1px 20px 0px 20px;
   font-size: 0.9em;
}
span.tag {
   cursor: pointer;
   padding: 1px 4px 0px 12px;
   display: inline-block;
   border: 1px solid #ccc;
   border-radius: 20px;
   margin-right: 5px;
   font-size: 0.9em;
}
i.fas.fa-times-circle {
   color: firebrick;
   margin-left: 5px;
}
p.pure-button.rotate {
   padding: 4px 20px;
   margin: 5px 0 0 0;
   width: 100%;
}
div.transcriptions {
   margin-left: 10px;
   border: 1px solid #ccc;
   position: relative;
}
div.transcription-title{
   margin:0 0 10px 0;
   padding: 4px 8px;
   border-bottom: 1px solid #ccc;
   display: flex;
   flex-flow: row nowrap;
}
div.transcription-title .head {
   font-size: 1.15em;
   font-weight: bold;
}
.transcription-info {
   margin: 10px;
   padding-bottom: 50px;
}
.transcription-info .actions {
   padding: 5px 5px 5px 10px;
   border-top: 1px solid #ccc;
   display: flex;
   flex-flow: row nowrap;
   position: absolute;
   bottom:0;
   left:0;
   right:0;
   align-items: center;;
}
.transcription-info .buttons {
   margin-left: auto;
}
.transcription-info .buttons {
   margin-left: auto;
}
.transcription-info .buttons .trans {
   margin-left: 10px;
}
.transcription-info label {
   font-weight: bold;
   margin-right: 10px;
}
i.paging {
   margin: 0 10px;
   cursor: pointer;
}
i.paging.disabled {
   opacity: 0.2;
   cursor: default;
}
.paging {
   margin-left: auto;
}
pre {
   font-size: 0.9em;
   font-family: sans-serif;
   white-space: pre-wrap;       /* Since CSS 2.1 */
   white-space: -moz-pre-wrap;  /* Mozilla, since 1999 */
   white-space: -pre-wrap;      /* Opera 4-6 */
   white-space: -o-pre-wrap;    /* Opera 7 */
}@media only screen and (min-width: 768px) {
   .zoom-wrap {
      flex-basis: 45%;
      max-width: 45%;
      margin-right: 15px;
   }
   div.transcriptions {
      flex-basis: 50%;
      max-width: 50%;
   }
}
@media only screen and (max-width: 768px) {
   .zoom-wrap {
      flex-basis: 95%;
      max-width: 95%;
      margin-bottom: 15px;
   }
   div.transcriptions {
      flex-basis: 95%;
      max-width: 95%;
   }
}
</style>