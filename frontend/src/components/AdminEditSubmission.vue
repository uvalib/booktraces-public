<template>
   <Button size="small" label="Edit" @click="showDialog=true" severity="info"/>
   <Dialog v-model:visible="showDialog" header="Edit Submission" :modal="true" position="top" @show="opened"
      maximizable :style="{ width: '75%' }" :blockScroll="true" :breakpoints="{'768px': '100%'}"
   >
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
         <div class="row">
            <label>Tags:</label>
            <div class="tags">
               <Button v-for="tag in editTagNames" icon="pi pi-trash" :label="tag" @click="removeTag(tag)"
                  severity="secondary" outlined rounded/>
               <Button v-if="!showTags" label="Add Tags" rounded severity="secondary" icon="pi pi-plus" @click="showTags=true"/>
            </div>
            <div v-if="showTags" class="source-tags">
               <div class="title">Available Tags:</div>
               <div class="tags">
                  <Button v-for="tag in availableTags" :label="tag.name" @click="addTag(tag)" severity="secondary" outlined rounded/>
               </div>
               <Button label="Done" rounded severity="secondary"  @click="showTags=false"/>
            </div>
         </div>
      </div>
      <template #footer>
         <Button label="Cancel" @click="showDialog=false" severity="secondary"/>
         <Button label="Save" @click="saveClicked" severity="info"/>
      </template>
   </Dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useSystemStore } from "@/stores/system"
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
import Dialog from 'primevue/dialog'

const emit = defineEmits( ['save'] )

const props = defineProps({
   details: {
      type: Object,
      required: true
   },
})

const system = useSystemStore()

const selectedInstitution = ref()
const edit = ref({
   title: "",
   author: "",
   publication: "",
   institution_id: 0,
   institution: "",
	callNumber: "",
	description: "",
	submitter: "",
   tags: []
})
const editTagNames = ref([])
const showDialog = ref(false)
const showTags = ref(false)

const availableTags = computed(() => {
   return system.tags.filter( st => !editTagNames.value.includes(st.name))
})

const opened = (() => {
   edit.value.title = props.details.title
   edit.value.author = props.details.author
   edit.value.publication = props.details.publication
	edit.value.callNumber = props.details.callNumber
	edit.value.description = props.details.description
	edit.value.submitter = props.details.submitter

   // NOTE: the back end sends up institution as two distinct fields in th details,
   // but the select tracks an institution object. Track that in selectedInstitution. When
   // saved, populate the institution_id and institution fields.
   selectedInstitution.value = system.institutions.find( i => i.id == props.details.institution_id)
   edit.value.institution_id = 0
   edit.value.institution = ""

   // NOTE: tags in the submission details is a list of names, but the backend wants a list
   // of tag IDs for the update. Leave tags array empty now and track edits in editTagNames.
   // When save is clicked, convert the edited list of names to IDs.
   edit.value.tags =  []
   editTagNames.value = props.details.tags
})

const addTag = ((tagObj) => {
   editTagNames.value.push(tagObj.name)
})

const removeTag = ((tagName) => {
   let idx = editTagNames.value.findIndex(t => t == tagName)
   if ( idx > -1 ) {
      editTagNames.value.splice(idx, 1)
   }
})

const saveClicked = ( async () => {
   if ( typeof selectedInstitution.value == 'string' ) {
      // if a new institution has been added, institution will just be the string name.
      // Create a new institution, then pull the newly created info into selectedInstitution
      await system.addInstitution( selectedInstitution.value )
      selectedInstitution.value = system.institutions.find( i => i.name == selectedInstitution.value )
   }

   // Popuate institution in the edit object as separate fields as the back end wants it
   edit.value.institution = selectedInstitution.value.name
   edit.value.institution_id = selectedInstitution.value.id

   // convert list of tag names to tags IDs as backend wants it
   edit.value.tags = []
   editTagNames.value.forEach( tagName => {
      let tag = system.tags.find( t => t.name == tagName)
      if ( tag ) {
         edit.value.tags.push( tag.id)
      }
   })
   emit('save', edit.value)
   showDialog.value = false
})
</script>

<style scoped lang="scss">
.admin-edit {
   display: flex;
   flex-direction: column;
   gap: 20px;

   .row {
      display: flex;
      flex-direction: column;
      gap: 5px;
   }
   .tags {
      display: flex;
      flex-flow: row wrap;
      gap: 5px 10px;
      margin-bottom: 10px;
   }
   .source-tags {
      border: 1px solid #ddd;
      padding: 10px;
      border-radius: 5px;
      .title {
         padding-bottom: 10px;
         border-bottom: 1px solid #ddd;
         margin-bottom: 20px;
      }
   }
}
</style>