<template>
   <div class="admin">
      <h2>
         <span>System Admin Panel</span>
         <span class="login">
            <b>Logged in as:</b>{{admin.loginName}}
         </span>
      </h2>
      <div>
         <h3>News</h3>
         <div class="ctls">
            <Button label="Add News" @click="addNewsClicked" size = "small"/>
         </div>
         <div class="error" v-if="admin.error">{{admin.error}}</div>
         <div class="news">
            <Panel v-for="item in system.news" :key="item.id">
               <template #header>
                  <div class="header">
                     <div class="left" v-if="editID != item.id">
                        <span class="title">{{item.title}}</span>
                        <span class="date">{{ item.createdAt.split("T")[0] }}</span>
                        <div class="published">
                           <label>
                              <Checkbox v-model="item.published" :binary="true" @update:modelValue="publicToggled(item)" />
                              <span>Published</span>
                           </label>
                        </div>
                     </div>
                     <div class="edit-title" v-else-if="editID == item.id" >
                        <label>Title:</label>
                        <InputText v-model="editTitle" />
                     </div>
                     <div class="ctls">
                        <template v-if="editID == item.id">
                           <Button label="Cancel" severity="secondary" size="small" @click="cancelEdit"  />
                           <Button label="Save" severity="primart" size="small" @click="saveEdit" />
                        </template>
                        <template v-else>
                           <Button rounded icon="pi pi-trash" severity="secondary" size="small" :disabled="editID != 0" @click="deleteNews(item)" />
                           <Button rounded icon="pi pi-file-edit" severity="secondary" size="small" :disabled="editID != 0"  @click="editNews(item) "/>
                        </template>
                     </div>
                  </div>
               </template>
               <div class="text" v-html="item.content" v-if="editID != item.id" />
               <Editor v-else v-model="editContent" editorStyle="height: 200px" />
            </Panel>
         </div>
      </div>
   </div>
   <Dialog v-model:visible="showAdd" :style="{width: '450px'}" header="Add News" :modal="true" position="top">
      <div class="add-form">
         <div class="row">
            <label for="new-title">Title</label>
            <InputText id="new-title" v-model="editTitle" />
         </div>
         <Editor v-model="editContent" editorStyle="height: 200px" />
      </div>
      <template #footer>
         <Button label="Cancel" severity="secondary" @click="showAdd = false" />
         <Button label="Add News"  @click="addNews" autofocus :disabled="editTitle == '' || editContent=='' || admin.working"/>
      </template>
   </Dialog>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useSystemStore } from "@/stores/system"
import { useAdminStore } from "@/stores/admin"
import Panel from 'primevue/panel'
import Checkbox from 'primevue/checkbox'
import InputText from 'primevue/inputtext'
import Editor from 'primevue/editor'
import Dialog from 'primevue/dialog'
import { useConfirm } from "primevue/useconfirm"

const admin = useAdminStore()
const system = useSystemStore()
const confirm = useConfirm()

const showAdd = ref(false)
const editID = ref(0)
const editTitle = ref("")
const editContent = ref("")

onMounted(() => {
   system.getNews()
})

const publicToggled = ((item) => {
   admin.updateNews(item.id, item.title, item.content)
})

const deleteNews = ((item) => {
   confirm.require({
      message: `Delete '${item.title}' All data will be permanently lost. Are you sure?`,
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
         admin.deleteNews(item.id)
      }
   })
})

const addNewsClicked = (() => {
   editID.value = 0
   editTitle.value = ""
   editContent.value = ""
   showAdd.value = true
})

const addNews = ( async () => {
   await admin.addNews( editTitle.value, editContent.value )
   if ( admin.error == "" ) {
      editID.value = 0
      editTitle.value = ""
      editContent.value = ""
      showAdd.value = false
   }
})

const editNews = ( (item) => {
   editTitle.value = item.title
   editContent.value = item.content
   editID.value = item.id
})

const cancelEdit = ( () => {
   editID.value = 0
   editTitle.value = ""
   editContent.value = ""
})

const saveEdit = ( async () => {
   await admin.updateNews(editID.value, editTitle.value, editContent.value)
   if ( admin.error == "" ) {
      editID.value = 0
      editTitle.value = ""
      editContent.value = ""
   }
})
</script>
//       togglePublication(idx) {
//          this.$store.dispatch("news/togglePublication", idx)
//       },

//       deleteClicked(event) {
//          let tgt = event.currentTarget
//          let newsID = tgt.dataset.id
//          let resp = confirm("Delete this news item? All data will be permanently lost. Continue?")
//          if (resp) {
//              this.$store.dispatch('news/deleteNews', newsID)
//          }
//       },

<style scoped lang="scss">
div.admin {
   padding: 15px 25px;
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
   .ctls {
      width: 100%;
      flex: 1;
      display: flex;
      flex-flow: row nowrap;
      gap: 10px;
      justify-content: flex-end;
   }
   :deep(p) {
      margin: 0 0;
   }
   :deep(a) {
      color: #24890d !important;
   }
   .news {
      margin-top: 10px;
      display: flex;
      flex-direction: column;
      gap: 15px;
      .text {
         padding-top: 10px;
      }
      .header {
         display: flex;
         flex-flow: row nowrap;
         justify-content: space-between;
         align-items: center;
         gap: 20px;
         width: 100%;
         .edit-title {
            width: 100%;
            display: flex;
            flex-flow: row nowrap;
            gap: 5px;
            align-items: center;
            input {
               flex:1;
            }
         }
         .left {
            display: flex;
            flex-direction: column;
            gap: 5px;
            .title {
               font-weight: bold;
            }
            .date {
               color: #aaa;
               font-size: 0.9em;
            }
            .published label{
               display: flex;
               flex-flow: row nowrap;
               gap: 10px;
            }
         }
      }
   }
}
.add-form {
   display:flex;
   flex-direction: column;
   gap: 15px;
   .row {
      display:flex;
      flex-direction: column;
      gap: 5px;
   }
}
</style>