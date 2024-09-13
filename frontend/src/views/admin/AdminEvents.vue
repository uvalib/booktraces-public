<template>
   <div class="admin">
      <h2>
         <span>System Admin Panel</span>
         <span class="login">
            <b>Logged in as:</b>{{admin.loginName}}
         </span>
      </h2>
      <div>
         <h3>Events</h3>
         <div class="error" v-if="admin.error">{{admin.error}}</div>
         <DataTable :value="system.events" ref="eventsTable" dataKey="id" size="small" stripedRows
            :lazy="false" :paginator="true" :rows="30" :rowsPerPageOptions="[30,50,100]"
            paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
            currentPageReportTemplate="{first} - {last} of {totalRecords}" paginatorPosition="top"
            editMode="row" v-model:editingRows="editRows" @row-edit-save="saveEdits" :loading="system.loading"
         >
            <template #paginatorstart>
               <Button label="Add Event" severity="secondary" size="small" @click="addClicked"/>
            </template>
            <Column field="date" header="Date" class="nowrap">
               <template #editor="{ data, field }">
                  <InputText v-model="data[field]" fluid />
                </template>
            </Column>
            <Column field="description" header="Description">
               <template #body="slotProps">
                  <div v-html="slotProps.data.description"></div>
               </template>
               <template #editor="{ data, field }">
                    <Textarea v-model="data[field]" :rows="2" fluid />
                </template>
            </Column>
            <Column :rowEditor="true"  bodyStyle="text-align:center" />
            <Column  style="width: 20px">
               <template #body="slotProps">
                  <Button icon="pi pi-trash" rounded text severity="secondary" @click="deleteEvent(slotProps.data)"/>
               </template>
            </Column>
         </DataTable>
      </div>
   </div>
   <Dialog v-model:visible="showAdd" :style="{width: '450px'}" header="Add Event" :modal="true" position="top">
      <div class="add-form">
         <div class="row">
            <label for="new-date">Date</label>
            <InputText id="new-date" v-model="newDate" />
         </div>
         <div class="row">
            <label for="new-desc">Description</label>
            <Textarea id="new-desc" v-model="newDesc" :rows="4"/>
         </div>
      </div>
      <template #footer>
         <Button label="Cancel" severity="secondary" @click="showAdd = false" autofocus />
         <Button label="Add Event"  @click="addEvent" autofocus :disabled="newDate == '' || newDesc=='' || admin.working"/>
      </template>
   </Dialog>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useSystemStore } from "@/stores/system"
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import Textarea from 'primevue/textarea'
import InputText from 'primevue/inputtext'
import { useAdminStore } from "@/stores/admin"
import { useConfirm } from "primevue/useconfirm"

const admin = useAdminStore()
const system = useSystemStore()
const confirm = useConfirm()

const editRows = ref([])
const showAdd = ref(false)
const newDate = ref("")
const newDesc = ref("")

onMounted(() => {
   system.getEvents()
})

const deleteEvent = ((evt) => {
   confirm.require({
      message: 'Delete this event? All data will be permanently lost. Are you sure?',
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
         admin.deleteEvent(evt.id)
      }
   })
})

const saveEdits = ( async ( event) => {
   // TODO error handling
   let { newData, index } = event
   await admin.updateEvent( index, newData )
})

const addClicked = (() => {
   newDate.value = ""
   newDesc.value = ""
   showAdd.value = true
})

const addEvent = ( async () => {
   // TODO error handling
   await admin.addEvent( newDate.value, newDesc.value )
   showAdd.value = false
})
</script>

<style scoped lang="scss">
div.admin {
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
}
.add-form {
   display: flex;
   flex-direction: column;
   gap: 15px;
   .row {
      display: flex;
      flex-direction: column;
      gap: 5px;
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