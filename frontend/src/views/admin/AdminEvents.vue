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
         <!-- <span @click="addEventClick" class="add pure-button pure-button-primary">Add Event</span> -->
         <div class="error" v-if="admin.error">{{admin.error}}</div>
         <DataTable :value="eventsStore.list" ref="eventsTable" dataKey="id" size="small" stripedRows
            :lazy="false" :paginator="true" :rows="30" :rowsPerPageOptions="[30,50,100]" :alwaysShowPaginator="false"
            paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
            currentPageReportTemplate="{first} - {last} of {totalRecords}" paginatorPosition="top"
            editMode="row" v-model:editingRows="editRows" @row-edit-save="saveEdits"
         >
            <Column field="date" header="Date" class="nowrap"/>
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
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useEventsStore } from "@/stores/events"
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Textarea from 'primevue/textarea'
import { useAdminStore } from "@/stores/admin"
import { useConfirm } from "primevue/useconfirm"

const admin = useAdminStore()
const eventsStore = useEventsStore()
const editRows = ref([])
const confirm = useConfirm()

onMounted(() => {
   eventsStore.getAll()
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

const saveEdits = (() => {

})

// saveClicked() {
//    if (this.addingNew) {
//       this.$store.dispatch('events/addEvent',this.editDetails).then((/*response*/) => {
//          this.edit=false
//          this.editDetails = null
//          this.addingNew = false
//       })
//    } else {
//       this.$store.dispatch('events/updateEvent',this.editDetails).then((/*response*/) => {
//          this.edit=false
//          this.editDetails = null
//          this.addingNew = false
//       })
//    }
// },
// editingEvent(id) {
//    if (this.edit == false) return false
//    return this.editDetails.id == id
// },
// deleteClicked(event) {
//    let tgt = event.currentTarget
//    let eventID = tgt.dataset.id
//    let resp = confirm("Delete this event? All data will be permanently lost. Continue?")
//    if (resp) {
//          this.$store.dispatch('events/deleteEvent', eventID)
//    }
// },
// addEventClick() {
//    this.$store.commit("events/addEventPlaceholder")
//    this.editDetails = Object.assign({}, this.events[0])
//    this.edit = true
//    this.addingNew = true
// }
</script>

<style scoped lang="scss">
div.admin {
   padding: 15px 25px;
   min-height: 600px;
   background: white;
   color: #444;
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
</style>