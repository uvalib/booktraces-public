<template>
   <div class="admin">
      <h2>
         <span>System Admin Panel</span>
         <span class="login">
            <b>Logged in as:</b>{{admin.loginName}}
         </span>
      </h2>
      <div>
         <h3>Pedagogy</h3>
         <div class="error" v-if="admin.error">{{ admin.error }}</div>
         <DataTable v-if="edit == false" :value="admin.pedagogy" ref="pedagogyTable" dataKey="id" size="small" stripedRows
            :lazy="false" :paginator="true" :rows="10" :rowsPerPageOptions="[10,25]"
            paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
            currentPageReportTemplate="{first} - {last} of {totalRecords}" paginatorPosition="top" :loading="admin.working"
         >
            <template #paginatorstart>
               <Button label="Add Document" severity="secondary" size="small" @click="addDocumentClicked"/>
            </template>
            <Column field="key" header="Key" class="nowrap"/>
            <Column field="title" header="Title" />
            <Column field="createdAt" header="Created">
               <template #body="slotProps">
                  {{ slotProps.data.createdAt.split("T")[0] }}
               </template>
            </Column>
            <Column>
               <template #body="slotProps">
                  <div class="ctls">
                     <Button icon="pi pi-trash" rounded text severity="secondary" size="large" @click="deleteClicked(slotProps.data)"/>
                     <Button icon="pi pi-file-edit" rounded text severity="secondary" size="large" @click="editClicked(slotProps.data)"/>
                  </div>
               </template>
            </Column>
         </DataTable>
         <div v-else class="doc-edit">
            <div class="row">
               <label for="doc-key">Key</label>
               <InputText id="doc-key" v-model="tgtDoc.key" :disabled="tgtDoc.key=='index'" />
            </div>
            <div class="row">
               <label for="doc-title">Title</label>
               <InputText id="doc-title" v-model="tgtDoc.title" />
            </div>
            <div class="note">
               <b>NOTE:</b>
               to add a link to a pedagogy document, add this to the document: $DOC[key::label]
               where key is the key of the target document and label is the text you would like to be displayed as the link.
            </div>
            <Editor v-model="tgtDoc.content" editorStyle="height: 500px" />
            <div class="ctls">
               <Button label="Cancel" severity="secondary"  @click="cancelEdit"/>
               <Button label="Submit"  @click="submitDocument" :disabled="!canSubmit" />
            </div>
         </div>
      </div>
   </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import { useAdminStore } from "@/stores/admin"
import { useConfirm } from "primevue/useconfirm"
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import Editor from 'primevue/editor'

const admin = useAdminStore()
const confirm = useConfirm()

const edit = ref(false)
const tgtDoc = ref({id: 0, createdAt: null, key: "", title: "", content: "", updatedAt: null})


const canSubmit = computed( () => {
   if (admin.working ) return false
   return (tgtDoc.value.key != "" || tgtDoc.value.title != "" || tgtDoc.value.content != "")
})
onMounted(() => {
   admin.getPedagogyDocuments()
})

const editClicked = ( (doc) => {
   edit.value = true
   tgtDoc.value = {
      id: doc.id, createdAt: doc.createdAt, key: doc.key,
      title: doc.title, content: doc.content, updatedAt: doc.updatedAt
   }
})

const cancelEdit = ( () => {
   edit.value =  false
   tgtDoc.value = {id: 0, createdAt: null, key: "", title: "", content: "", updatedAt: null}
})

const submitDocument = ( async () => {
   if ( tgtDoc.value.id > 0) {
      await admin.updateDocument( tgtDoc.value )
   } else {
      await admin.addDocument( tgtDoc.value )
   }
   edit.value = false
})

const addDocumentClicked = ( () => {
   edit.value = true
   tgtDoc.value = {id: 0, createdAt: null, key: "", title: "", content: "", updatedAt: null}
})

const deleteClicked = ( (doc) => {
   confirm.require({
      message: `Delete document '${doc.title}'?<br/>All data will be permanently lost.<br/><br/>Are you sure?`,
      header: 'Confirm Delete',
      icon: 'pi pi-exclamation-triangle',
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: 'Delete'
      },
      accept: () => {
         admin.deleteDocument( doc.key )
      }
   })
})
</script>

<style lang="scss" scoped>
div.admin {
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
   .doc-edit {
      display: flex;
      flex-direction: column;
      gap: 15px;
      .row {
         display: flex;
         flex-direction: column;
         gap: 5px;
      }
      .ctls {
         display: flex;
         flex-flow: row nowrap;
         gap: 5px;
         justify-content: flex-end;
      }
   }

   .note {
      margin: 20px 0;
      background-color: white;
      padding: 15px;
      border: 2px solid cornflowerblue;
   }
}
@media only screen and (min-width: 768px) {
   .admin {
      padding: 15px 25px;
   }
}
@media only screen and (max-width: 768px) {
   .admin {
      padding: 10px;
   }
}
</style>