<template>
   <div class="admin-edit">
      <div class="row">
         <label for="title">Title of book</label>
         <InputText id="title" v-model="edit.title" />
      </div>
      <div class="row">
         <label for="author">Author (Last name, first name)</label>
         <InputText id="author" v-model="edit.author" />
      </div>
      <div class="row">
         <label for="publication">Publication place/date (e.g., London, 1888)</label>
         <InputText id="publication" v-model="edit.publication" />
      </div>
      <div class="row">
         <label for="institution">Institution where found</label>
         <Select v-model="selectedInstitution" :options="system.institutions"
            optionLabel="name" :editable="true"
            placeholder="Select or create an institution" />
      </div>
      <div class="row">
         <label for="callnumber">Call Number</label>
         <InputText id="callnumber" v-model="edit.callNumber" />
      </div>
      <div class="row">
         <label for="description">Description</label>
         <Textarea id="description" v-model="edit.description" rows="4" />
      </div>
      <div class="row">
         <label for="submitter">Submitted By</label>
         <InputText id="submitter" v-model="edit.submitter" />
      </div>
        <!--
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
         -->
      <div class="buttons">
         <Button size="small" label="Cancel" @click="emit('cancel')" severity="secondary"/>
         <Button size="small" label="Save" @click="saveClicked" severity="info"/>
      </div>
   </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useDetailsStore } from "@/stores/details"
import { useSystemStore } from "@/stores/system"
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'

const emit = defineEmits( ['cancel', 'save'] )

const details = useDetailsStore()
const system = useSystemStore()

const selectedInstitution = ref()
const edit = ref({})

onMounted(() => {
   edit.value = Object.assign({}, details.submission)
   selectedInstitution.value = system.institutions.find( i => i.id == edit.value.institution_id)
   console.log(edit.value)
   console.log(selectedInstitution.value)
})

const saveClicked = (() => {
   // TODO move selectedInstiturion into edit and create as necessary
   emit('save', edit.value)
})
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
//       addTagClicked() {
//          this.showTagList = true;
//       },
</script>

<style scoped lang="scss">
.admin-edit {
   display: flex;
   flex-direction: column;
   gap: 20px;
   margin: 50px auto;

   .row {
      display: flex;
      flex-direction: column;
      gap: 5px;
   }
   .buttons {
      display: flex;
      flex-flow: row wrap;
      justify-content: flex-end;
      align-items: center;
      gap: 5px;
   }
}
@media only screen and (min-width: 768px) {
   .admin-edit {
      width: 70%;
   }
}
@media only screen and (max-width: 768px) {
   .admin-edit {
      width: 95%;
   }
}
</style>