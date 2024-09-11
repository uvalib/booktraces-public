<template>
   <div class="admin">
      <h2>
         <span>System Admin Panel</span>
         <span class="login">
            <b>Logged in as:</b>{{admin.loginName}}
         </span>
      </h2>
      <div>
         <h3>Submissions</h3>
         <div class="list-controls">
            <div class="left">
               <div class="search">
                  <input v-model="admin.submissions.query">
                  <Button @click="admin.getSubmissions()" size="small" label="Search"/>
               </div>
               <Button v-if="admin.submissions.tagFilter" @click="removeFilter" size="small" severity="secondary"
                  outlined rounded icon="pi pi-times-circle"
                  :label="`Items Tagged: ${admin.submissions.tagFilter}`"/>
            </div>
         </div>
         <DataTable :value="admin.submissions.hits" ref="adminSubTable" dataKey="id"
            stripedRows showGridlines responsiveLayout="scroll" :loading="admin.working"
            :lazy="true" :paginator="true" @page="onPageClicked($event)"
            :rows="50" :totalRecords="admin.submissions.total"
            paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
            :first="admin.submissions.start" currentPageReportTemplate="{first} - {last} of {totalRecords}" paginatorPosition="top"
         >
            <Column field="id" header="ID" />
            <Column field="title" header="Title" />
            <Column field="author" header="Author"/>
            <Column header="Tags" >
               <template #body="slotProps">
                  <div class="tags">
                     <Button v-for="tag in tagsList(slotProps.data)" :key="tag"
                        size="small" severity="info" :label="tag" @click="tagClicked(tag)"/>
                  </div>
               </template>
            </Column>
            <Column field="submittedAt" header="Submitted" class="nowrap" >
               <template #body="slotProps">
                  {{ slotProps.data.submittedAt.split("T")[0] }}
               </template>
            </Column>
            <Column header="Actions" class="icon" >
               <template #body="slotProps">
                  <div class="acts">
                     <Button @click="viewSubmission(slotProps.data)" size="small" severity="info" label="View" />
                     <Button @click="togglePublish(slotProps.data)" size="small"
                        :severity="publishSeverity(slotProps.data)"
                        :label="publishLabel(slotProps.data)"/>
                     <Button @click="deleteSubmisson(slotProps.data)" size="small" severity="danger" label="Delete" />
                  </div>
               </template>
            </Column>
         </DataTable>
      </div>
   </div>
</template>

<script setup>
import { onMounted } from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { useSystemStore } from "@/stores/system"
import { useAdminStore } from "@/stores/admin"
import { useRouter } from 'vue-router'
import { useConfirm } from "primevue/useconfirm"

const system = useSystemStore()
const admin = useAdminStore()
const router = useRouter()
const confirm = useConfirm()

onMounted(() => {
   system.getInstitutions()
   system.getTags()
   admin.getSubmissions()
})

const tagsList = ( (sub) => {
   if ( sub.tags) {
      return sub.tags.split(',')
   }
   return []
})

const removeFilter = (() => {
   admin.submissions.tagFilter = ""
   admin.getSubmissions()
})

const onPageClicked = ((event) => {
   admin.submissions.start = event.first
   admin.getSubmissions()
})

const tagClicked = ((tag) => {
   event.stopPropagation()
   admin.submissions.tagFilter = tag
   admin.getSubmissions()
})

const publishLabel = ( (submission) => {
   if (submission.published) {
      return "Unpublish"
   }
   return  "Publish"
 })

 const publishSeverity = ( (submission) => {
   if (submission.published) {
      return "secondary"
   }
   return  "primary"
})

const viewSubmission = ( (submission) => {
   router.push("/admin/submissions/" + submission.id)
})

const togglePublish = ( (submission) => {
   let act = "Publish"
   if (submission.published) {
      act = "Unpublish"
   }
   confirm.require({
      message: `${act} this submission?`,
      header: `Confirm ${act}`,
      icon: 'pi pi-exclamation-triangle',
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: act
      },
      accept: () => {
         admin.updatePublicationStatus(submission.id, !submission.published )
      }
   })
})

const deleteSubmisson = ( (submission) => {
   confirm.require({
      message: 'Delete this submission? All data and unloaded files will be permanently lost. Are you sure?',
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
         admin.deleteSubmission(submission.id)
      }
   })
})
</script>

<style scoped lang="scss">
div.admin {
   padding: 15px 25px;
   min-height: 600px;
   background: white;
   color: #444;
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
.list-controls {
   display: flex;
   flex-flow: row wrap;
   justify-content: space-between;
   align-items: flex-start;
   margin-bottom: 15px;
   .left {
      display: flex;
      flex-flow: row nowrap;
      gap: 15px;
      .search {
         display: flex;
         flex-flow: row nowrap;
         align-items: stretch;
         gap: 10px;
         input {
            border: 1px solid #ccc;
            border-radius: 4px;
            width: 300px;
         }
      }
      .tag-filter {
         font-size: 0.8em;
         border: 1px solid #ccc;
         padding: 2px 4px 0 10px;
         border-radius: 20px;
         cursor: pointer;
      }
   }
}
:deep(td.icon) {
   text-align: center;
}
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: space-between;
   gap: 5px;
   button {
      flex: 1;
   }
}
.tags {
   display: flex;
   flex-flow: row wrap;
   gap: 5px;
}
</style>