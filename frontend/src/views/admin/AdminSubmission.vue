<template>
   <div class="admin-submission">
      <h2>
         <span>System Admin Panel</span>
         <span class="login">
            <b>Logged in as:</b>{{admin.loginName}}
         </span>
      </h2>
      <BTSpinner v-if="system.loading==true" message="Loading details..." />
      <template v-else>
         <div class="actions">
            <router-link to="/admin">
               <i class="pi pi-arrow-left"></i>&nbsp;Back to Submissions
            </router-link>
            <div class="buttons">
               <template v-if="!edit">
                  <Button size="small" label="Edit" @click="editClicked" severity="info"/>
                  <Button  v-if="!details.submission.published" size="small" label="Publish" @click="publishClicked" severity="info"/>
                  <Button  v-else size="small" label="Unpublish" @click="unpublishClicked" severity="info"/>
                  <Button size="small" label="Delete" @click="deleteClicked" severity="danger"/>
               </template>
               <template v-else>
                  <Button size="small" label="Cancel" @click="cancelClicked" severity="secondary"/>
                  <Button size="small" label="Save" @click="saveClicked" severity="info"/>
               </template>
            </div>
         </div>
         <div class="error" v-if="admin.error">{{admin.error}}</div>
         <div v-if="!edit" class="details">
            <table class="submit-info">
               <tr>
                  <td class="label">Visible to public:</td>
                  <td class="value">{{ published }}</td>
               </tr>
               <tr>
                  <td class="label">Title:</td>
                  <td class="value">{{details.submission.title}}</td>
               </tr>
               <tr>
                  <td class="label">Author:</td>
                  <td class="value">{{details.submission.author}}</td>
               </tr>
               <tr>
                  <td class="label">Publication details:</td>
                  <td class="value">{{details.submission.publication}}</td>
               </tr>
               <tr>
                  <td class="label">Institution:</td>
                  <td class="value">{{details.submission.institution}}</td>
               </tr>
               <tr>
                  <td class="label">Call number:</td>
                  <td class="value">{{details.submission.callNumber}}</td>
               </tr>
               <tr>
                  <td class="label">Submitted by:</td>
                  <td class="value">{{details.submission.submitter}} ( {{details.submission.email}} )</td>
               </tr>
               <tr>
                  <td class="label">Submitted on:</td>
                  <td class="value">{{ details.submission.submittedAt.split("T")[0] }}</td>
               </tr>
               <tr>
                  <td class="label">Tags:</td>
                  <td class="value">{{submissionTagCSV}}</td>
               </tr>
               <tr>
                  <td class="label">Description:</td>
                  <td class="value">
                     <span v-html="description"></span>
                  </td>
               </tr>
            </table>
            <div class="thumbs">
               <div class="thumb" v-for="file in details.submission.files">
                  <div class="zoom-wrap">
                     <vue-image-zoomer :regular="file.url" />
                     <!-- <p @click="rotateClicked(f.url)" class="pure-button rotate">Rotate Right</p> -->
                  </div>
                  <div class="transcription-panel">
                     <div class="head">Transcriptions</div>
                     <div v-if="file.transcriptions.length == 0" class="none">None</div>
                     <tempate v-else>
                        <div class="transcribe-acts">
                           <div class="buttons">
                              <Button v-if="file.transcriptions[transcriptionIdx].approved==false" aria-label="approve transcription"
                                 icon="pi pi-check-circle"  size="small" @click="approveTranscription(file)"/>
                              <Button icon="pi pi-file-edit"  severity="secondary" size="small" @click="editTranscription(file)"/>
                              <Button icon="pi pi-trash"  severity="danger" size="small" @click="deleteTranscription(file)"/>
                           </div>
                           <span class="paging">
                              <Button icon="pi pi-chevron-left" rounded severity="secondary" size="small"
                                 :disabled="transcriptionIdx == 0" @click="priorTran()"/>
                              <span>{{transcriptionIdx+1}} of {{file.transcriptions.length}}</span>
                              <Button icon="pi pi-chevron-right" rounded severity="secondary" size="small"
                                 :disabled="transcriptionIdx == file.transcriptions.length-1" @click="nextTran(file)"/>
                           </span>
                        </div>
                        <div class="transcription-info" v-if="file.transcriptions.length > 0">
                           <table>
                              <tr>
                                 <td class="label">Date:</td>
                                 <td>{{ getTranscribeDate(file) }}</td>
                              </tr>
                              <tr>
                                 <td class="label">Submitter:</td>
                                 <td>{{ getTranscriber(file) }}</td>
                              </tr>
                              <tr>
                                 <td class="label">Status:</td>
                                 <td>{{ getTranscribeStatus(file) }}</td>
                              </tr>
                           </table>
                           <template v-if="editTrans">
                              <!-- <textarea class="edit-trans" rows="3" v-model="workingTrans"></textarea>
                              <div v-if="error" class="error">{{error}}</div>
                              <div class="actions">
                                 <span class="buttons">
                                    <span @click="cancelEdit" class="pure-button trans">Cancel</span>
                                    <span @click="submitEdit(f)" class="pure-button trans">Submit</span>
                                 </span>
                              </div> -->
                           </template>
                           <template v-else>
                              <pre class="transcription">{{file.transcriptions[transcriptionIdx].text}}</pre>
                           </template>
                        </div>
                     </tempate>
                  </div>
               </div>
            </div>
         </div>
         <!-- <div v-else class="edit details  pure-form">
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
                        <i class="pi pi-trash"></i>
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

         </div> -->
      </template>
   </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import { useAdminStore } from "@/stores/admin"
import { useSystemStore } from "@/stores/system"
import { useDetailsStore } from "@/stores/details"
import { useRoute, useRouter } from 'vue-router'

const system = useSystemStore()
const admin = useAdminStore()
const details = useDetailsStore()
const route = useRoute()
const router = useRouter()

const edit = ref(false)
const editTrans = ref(false)
const transcriptionIdx = ref(0)

const published = computed(() => {
   if (details.submission.published) return "YES"
   return "NO"
})

const submissionTagCSV = computed(() => {
   if (details.submission.tags) {
      return details.submission.tags.join(", ")
   }
   return ""
})

const description = computed(() => {
   return details.submission.description.replace(/\r|\r\n/gm, "\n").replace(/\n+/gm, "<br/><br/>")
})

onMounted(() => {
   console.log("ID "+route.params.id )
   details.getSubmission( route.params.id )
   system.getTags()
   system.getInstitutions()
})
const nextTran = (( image ) => {
   if (transcriptionIdx.value == f.transcriptions.length -1) {
      return
   }
   transcriptionIdx.value++
})

const priorTran = (() => {
   if (transcriptionIdx.value == 0) {
      return
   }
   transcriptionIdx.value--
})

const getTranscribeDate = ((f) => {
   if (f.transcriptions.length == 0) return ""
   let t = f.transcriptions[transcriptionIdx.value]
   return t.transcribed_at.split("T")[0]
})
const getTranscriber = ((f) => {
   if (f.transcriptions.length == 0) return ""
   let t = f.transcriptions[transcriptionIdx.value]
   return `${t.transcriber_email} (${t.transcriber})`
})
const getTranscribeStatus = ((f) => {
   if (f.transcriptions.length == 0) return ""
   let t = f.transcriptions[transcriptionIdx.value]
   if (t.approved) {
      return "Approved"
   }
   return "Pending"
})

//    methods: {
//       cancelEdit() {
//          this.editTrans = false
//          this.workingTrans = ""
//       },
//       async submitEdit(f) {
//          let data = {submissionID: details.submission.id,
//             fileID: f.id,
//             transcriptionID: f.transcriptions[transcriptionIdx.value].id,
//             transcription: this.workingTrans}
//          await this.$store.dispatch("transcribe/update", data)
//          if (this.error == "" || this.error == null) {
//             this.editTrans = false
//             this.workingTrans = ""
//          }
//       },
//       approveTranscription(f) {
//          let t = f.transcriptions[transcriptionIdx.value]
//          this.$store.dispatch("transcribe/approve", t.id)
//       },
//       editTranscription(f) {
//          this.editTrans = true
//          this.workingTrans = f.transcriptions[transcriptionIdx.value].text
//       },
//       deleteTranscription(f) {
//          let resp = confirm("Are you sure you want to delete this transcrption?")
//          if (resp) {
//             let t = f.transcriptions[transcriptionIdx.value]
//             this.$store.dispatch("transcribe/delete", t.id)
//             transcriptionIdx.value = 0
//          }
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
//             submissionID: details.submission.id,
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
//                id: details.submission.id,
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
//             id: details.submission.id,
//             public: true
//          });
//       },
//       unpublishClicked() {
//          this.$store.dispatch("admin/updatePublicationStatus", {
//             id: details.submission.id,
//             public: false
//          });
//       }
</script>

<style scoped lang="scss">
@media only screen and (min-width: 768px) {
   .admin-submission {
      padding: 15px 25px;
   }
}
@media only screen and (max-width: 768px) {
   .admin-submission {
      padding: 10px;
   }
}
div.admin-submission {
   min-height: 600px;
   background: white;
   color: #444;

   h3 {
      margin-bottom: 10px;
   }
   h2 {
      display: flex;
      flex-flow: row wrap;
      justify-content: space-between;
      align-items: flex-start;
      font-size: 1.5em;
      font-weight: bold;
      border-bottom: 1px dashed #666;
      font-family: 'Special Elite', cursive;
      padding-bottom: 5px;
      margin-bottom: 15px;
      .login {
         color: #999;
         font-family: sans-serif;
         font-size: 0.8em;
         float: right;
         font-weight: normal;
         b {
            color: #444;
            margin-right: 10px;
         }
      }
   }
   .error {
      color: firebrick;
      padding: 5px 10px;
      background-color: #ffeded;
      border: 2px solid firebrick;
      border-radius: 5px;
      margin: 10px 0 20px 0;
   }
   .buttons {
      display: flex;
      flex-flow: row wrap;
      justify-content: flex-start;
      align-items: center;
      gap: 5px;
   }
   div.actions, .transcribe-acts {
      padding: 10px;
      display: flex;
      flex-flow: row wrap;
      justify-content: space-between;
      align-content: center;
      gap: 20px;
      border-bottom: 1px solid #ddd;
      margin-bottom: 15px;
   }
   .paging {
      display: flex;
      flex-flow: row nowrap;
      align-items: center;
      gap: 10px;
   }
   div.details {
      margin: 0 25px 0 25px;
   }
   .thumbs {
      border-top: 1px solid #ddd;
      margin-top: 20px;
      padding-top: 20px;
      display: flex;
      flex-direction: column;
      gap: 20px;
      div.thumb {
         padding-bottom: 20px;
         border-bottom: 1px solid #ccc;
         display: flex;
         align-items: stretch;
         gap: 15px;
         justify-content: flex-start;
         div.zoom-wrap {
            min-height: 400px;
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
      div.transcription-panel {
         margin: 0;
         border: 1px solid #ddd;
         border-radius: 4px;
         .none {
            font-size: 1.3em;
            text-align: center;
            margin: 15% auto;
         }
         .head {
            background: #fafafa;
            padding: 5px 10px;
            display: flex;
            flex-flow: row wrap;
            justify-content: space-between;
            align-items: center;
            gap: 20px;
            border-bottom: 1px solid #ccc;
            border-radius: 4px 4px 0 0;
         }
      }
   }
}

pre.transcription  {
   margin: 15px;
   padding: 20px;
   border: 1px solid #ddd;
   border-radius: 4px;
   white-space: pre-wrap;       /* Since CSS 2.1 */
   white-space: -moz-pre-wrap;  /* Mozilla, since 1999 */
   white-space: -pre-wrap;      /* Opera 4-6 */
   white-space: -o-pre-wrap;    /* Opera 7 */
   word-wrap: break-word;       /* Internet Explorer 5.5+ */
}

table {
   width: 100%;
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
      white-space:  nowrap;
   }
}
// td input {
//    border: 1px solid #ccc;
//    width: 100%;
//    outline: none;
//    box-sizing: border-box;
// }
// td textarea {
//    display: block;
//    border: 1px solid #ccc;
//    width: 90%;
// }
@media only screen and (min-width: 768px) {
   .zoom-wrap {
      flex-basis: 50%;
      max-width: 50%;
   }
   div.transcription-panel {
      flex-basis: 50%;
   }
}
@media only screen and (max-width: 768px) {
   .zoom-wrap {
      flex-basis: 100%;
   }
   div.transcription-panel {
      flex-basis: 100%;
   }
}
</style>