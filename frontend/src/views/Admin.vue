<template>
   <div class="admin">
      <h2>System Admin Panel <span class="login"><b>Logged in as:</b>{{loginName}}</span></h2>
      <div>
         <h3>Submissions</h3>
         <div class="error">{{error}}</div>
         <div class="list-controls">
            <div class="search pure-button-group" role="group">
               <input @input="updateSearchQuery" @keyup.enter="searchClicked" type="text" id="search">
               <button @click="searchClicked" class="search pure-button pure-button-primary">Search</button>
            </div>
            <AdminPager/>
         </div>
         <table class="pure-table">
            <thead>
               <th>ID</th>
               <th>Title</th>
               <th style="width:200px;">Author</th>
               <th style="width:200px;">Tags</th>
               <th style="width:75px;">Submitted</th>
               <th class="checkbox">Public</th>
               <th style="width:35px;"></th>
            </thead>
            <tr v-for="sub in submissions" :key="sub.id" :data-id="sub.id" @click="submissionClicked">
               <td>{{ sub.id }}</td>
               <td>{{ sub.title }}</td>
               <td>{{ sub.author }}</td>
               <td>{{ tagList(sub) }}</td>
               <td style="text-align:center;">{{ sub.submittedAt.split("T")[0] }}</td>
               <td class="centered">
                  <span @click="togglePublishClicked" :data-id="sub.id" :data-published="isPublished(sub)" v-html="publishIcon(sub)"></span>
               </td>
               <td class="centered">
                  <i :data-id="sub.id" title="Delete" class="action fas fa-trash-alt" @click="deleteClicked"></i>
               </td>
            </tr>
         </table>
      </div>
   </div>
</template>

<script>
import { mapState } from 'vuex'
import { mapGetters } from 'vuex'
import AdminPager from "@/components/AdminPager"
export default {
   name: "admin",
   components: {
      AdminPager
   },
   computed: {
      ...mapState({
         total: state => state.admin.totalSubmissions,
         submissions: state => state.admin.submissions,
         error: state => state.error,
      }),
      ...mapGetters({
         loginName: 'admin/loginName',
      })
   },
   methods: {
      tagList( sub ) {
         if (sub.tags) {
            return sub.tags.split(",").join(", ")
         } 
         return ""
      },
      updateSearchQuery(e) {
         this.$store.commit('admin/updateSearchQuery', e.target.value)
      },
      searchClicked() {
         this.$store.dispatch("admin/getSubmissions")
      },
      isPublished(sub) {
         return sub.published
      },
      publishIcon(sub) {
         if (sub.published) {
            return '<i class="published fas fa-check-circle"></i>'
         } else {
            return '<i class="unpublished fas fa-times-circle"></i>'
         }
      },
      submissionClicked(event) {
         let tgt = event.currentTarget
         let subID = tgt.dataset.id
         this.$router.push("admin/submissions/" + subID)
      },
      deleteClicked(event) {
         event.stopPropagation()
         let resp = confirm("Delete this submission? All data and unloaded files will be permanently lost. Are you sure?")
         if (resp) {
            let id = event.currentTarget.dataset.id
            this.$store.dispatch("admin/deleteSubmission", {id:id})
         }
      },
      togglePublishClicked(event) {
         event.stopPropagation()
         let published = event.currentTarget.dataset.published
         if (published ) {
            let resp = confirm("Unpublish this submission?")
            if (!resp) {
               return
            }
         } else {
            let resp = confirm("Publish this submission?")
            if (!resp) {
               return
            }
         }
         let id = event.currentTarget.dataset.id
         this.$store.dispatch("admin/updatePublicationStatus", {id:id, public: !published})
      }
   },
   created() {
      this.$store.commit("clearSubmissionDetail")
      this.$store.dispatch("admin/getSubmissions")
   },
};
</script>

<style scoped>
td span >>> .published {
   color:green;font-size:1.25em;
   opacity: 0.5;
}
td span >>> .published:hover,td span >>> .unpublished:hover {
   opacity: 1;
   cursor: pointer;
}
td span >>> .unpublished {
   color:firebrick;font-size:1.25em;
   opacity: 0.5;
}
div.list-controls {
  position:relative;
  margin: 25px 0 5px 0;
}
div.search {
  font-size: 14px;
}
div.search button.search.pure-button {
  padding: 3px 15px;
}
div.search input {
  width:300px;
  outline:none;
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
i.fas.action.disabled {
   cursor: default;
   opacity: 0.2;
}
i.fas.action {
   cursor: pointer;
   opacity: 0.5;
   margin: 0 5px;
   font-size: 1.1em;
}
i.fas:hover {
   opacity: 1;
}
h2 {
  font-size: 1.5em;
  font-weight: bold;
  border-bottom: 1px dashed #666;
  font-family: 'Special Elite', cursive;
  padding-bottom: 5px;
  margin-bottom: 15px;
}
div.admin {
   padding: 15px 25px;
   min-height: 600px;
   background: white;
   color: #444;
}
table {
   width: 100%;
   font-size: 0.85em;
   color: #444;
}
table tr:hover {
   cursor: pointer;
   background: #f5f5f5;
}
th.checkbox {
   width: 40px;
}
td.centered {
   text-align: center;
}
</style>