<template>
   <div class="admin">
      <h2>System Admin Panel <span class="login"><b>Logged in as:</b>{{loginName}}</span></h2>
      <div>
         <h3>Pedagogy</h3>
         <div class="error">{{error}}</div>
         <div class="edit pure-form" v-if="edit">
            <div>
               <label>Key</label>
               <input v-if="working.key!='index'" class="title" type="text" v-model="working.key">
               <input v-else readonly class="title" type="text" v-model="working.key">
            </div>
            <div><label>Title</label><input class="title" type="text" v-model="working.title"></div>
            <vue-editor v-model="working.content"></vue-editor>
            <div class="controls">
               <span @click="cancelEditClick" class="btn pure-button pure-button-secondary">Cancel</span>
               <span @click="submitEditClick" class="btn pure-button pure-button-primary">Submit</span>
            </div>
         </div>
         <table v-else class="pure-table">
            <thead>
               <th>Key</th>
               <th>Title</th>
               <th style="width:75px;">Created</th>
               <th style="width:35px;"></th>
            </thead>
            <tr v-for="doc in documents" :key="doc.id" @click="documentClicked(doc)">
               <td>{{ doc.key }}</td>
               <td>{{ doc.title }}</td>
               <td style="text-align:center;">{{ doc.createdAt }}</td>
               <td class="centered">
                  <i v-if="doc.key!='index'" title="Delete" class="action fas fa-trash-alt" @click="deleteClicked(doc)"></i>
               </td>
            </tr>
         </table>
      </div>
   </div>
</template>

<script>
import { mapState, mapGetters } from 'vuex'
import { VueEditor } from 'vue2-editor'
export default {
   name: "admin-pedagogy",
   data: function() {
      return {
         edit: false,
         working: {
            id: 0,
            key: "",
            title: "",
            content: ""
         }
      }
   },
   components: {
      "vue-editor": VueEditor
   },
   computed: {
      ...mapState({
         documents: state => state.pedagogy.list,
         error: state => state.error,
         loading: state => state.loading,
      }),
      ...mapGetters({
         loginName: 'admin/loginName',
      })
   },
   methods: {
      cancelEditClick() {
         this.edit = false
      },
      async submitEditClick() {
         await this.$store.dispatch("pedagogy/updateDocument", this.working)
         this.edit = false
      },
      documentClicked(doc) {
         this.edit = true
         this.working.id = doc.id
         this.working.key = doc.key
         this.working.title = doc.title
         this.working.content = doc.content
      },
      deleteClicked(doc) {
         event.stopPropagation()
         let resp = confirm("Permanently delete this document?")
         if (resp) {
            this.$store.dispatch("pedagogy/deleteDocument", doc.key)
         }
      },
   },
   created() {
      this.$store.commit("pedagogy/clearList")
      this.$store.dispatch("pedagogy/getList")
   },
};
</script>

<style lang="scss" scoped>
div.admin {
   padding: 15px 25px;
   min-height: 600px;
   background: white;
   color: #444;

   h2 {
      font-size: 1.5em;
      font-weight: bold;
      border-bottom: 1px dashed #666;
      font-family: 'Special Elite' serif;
      padding-bottom: 5px;
      margin-bottom: 15px;
   }

   span.login {
      font-family: sans-serif;
      font-size: 0.5em;
      float: right;
      font-weight: 100;
      b {
         margin-right: 5px;
      }
   }
   table {
      width: 100%;
      font-size: 0.85em;
      color: #444;
      td.centered {
         text-align: center;
      }
      tr:hover {
         cursor: pointer;
         background: #f3f3f3;
      }
   }
   .edit {
      input {
         width: 100%;
         box-sizing: border-box;
         margin-bottom: 5px;
      }
      label {
         font-weight: bold;
         margin-bottom: 2px;
         display:block;
      }
      .controls {
         margin-top: 15px;
         text-align: right;
         .btn {
            margin-left: 10px;
            display: inline-block;
         }
      }
   }
}
::v-deep .ql-container {
   min-height: 400px;
}
</style>