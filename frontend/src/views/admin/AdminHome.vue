<template>
   <div class="admin">
      <h2>System Admin Panel <span class="login"><b>Logged in as:</b>{{admin.loginName}}</span></h2>
      <div>
         <h3>Submissions</h3>
         <!-- <div class="error">{{error}}</div>
         <div class="list-controls">
            <div class="search pure-button-group" role="group">
               <input @input="updateSearchQuery" @keyup.enter="searchClicked" type="text" id="search" :value="queryStr">
               <button @click="searchClicked" class="search pure-button pure-button-primary">Search</button>
            </div>
            <span class="tag-filter" v-if="tgtTag.length > 0">
               <b>Items Tagged:</b> {{tgtTag}} <i @click="removeFilter" class="unfilter fas fa-times-circle"></i>
            </span>
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
               <td>
                  <span class="tag" v-for="(tag,idx) in tagList(sub)" :key="idx" @click="tagClicked">{{tag}}</span>
               </td>
               <td style="text-align:center;">{{ sub.submittedAt.split("T")[0] }}</td>
               <td class="centered">
                  <span @click="togglePublishClicked" :data-id="sub.id" :data-published="isPublished(sub)" v-html="publishIcon(sub)"></span>
               </td>
               <td class="centered">
                  <i :data-id="sub.id" title="Delete" class="action fas fa-trash-alt" @click="deleteClicked"></i>
               </td>
            </tr>
         </table> -->
      </div>
   </div>
</template>

<script setup>
import { onMounted } from 'vue'
import AdminPager from "@/components/AdminPager.vue"
import { useSystemStore } from "@/stores/system"
import { useAdminStore } from "@/stores/admin"

const system = useSystemStore()
const admin = useAdminStore()

onMounted(() => {
   system.getInstitutions()
   system.getTags()
   admin.getSubmissions()
})

//       ...mapState({
//          total: state => state.admin.totalSubmissions,
//          submissions: state => state.admin.submissions,
//          error: state => state.error,
//          loading: state => state.loading,
//          tgtTag: state => state.admin.tgtTag,
//          queryStr: state => state.admin.queryStr,
//       }),
//       ...mapGetters({
//          loginName: 'admin/loginName',


//       tagList( sub ) {
//          if (sub.tags) {
//             return sub.tags.split(",")
//          }
//          return []
//       },
//       updateSearchQuery(e) {
//          this.$store.commit('admin/updateSearchQuery', e.target.value)
//       },
//       searchClicked() {
//          this.$store.dispatch("admin/getSubmissions")
//       },
//       isPublished(sub) {
//          return sub.published
//       },
//       publishIcon(sub) {
//          if (sub.published) {
//             return '<i class="published fas fa-check-circle"></i>'
//          } else {
//             return '<i class="unpublished fas fa-times-circle"></i>'
//          }
//       },
//       submissionClicked(event) {
//          let tgt = event.currentTarget
//          let subID = tgt.dataset.id
//          this.$router.push("/admin/submissions/" + subID)
//       },
//       deleteClicked(event) {
//          event.stopPropagation()
//          let resp = confirm("Delete this submission? All data and unloaded files will be permanently lost. Are you sure?")
//          if (resp) {
//             let id = event.currentTarget.dataset.id
//             this.$store.dispatch("admin/deleteSubmission", {id:id})
//          }
//       },
//       togglePublishClicked(event) {
//          event.stopPropagation()
//          let published = event.currentTarget.dataset.published
//          if (published ) {
//             let resp = confirm("Unpublish this submission?")
//             if (!resp) {
//                return
//             }
//          } else {
//             let resp = confirm("Publish this submission?")
//             if (!resp) {
//                return
//             }
//          }
//          let id = event.currentTarget.dataset.id
//          this.$store.dispatch("admin/updatePublicationStatus", {id:id, public: !published})
//       },
//       tagClicked(event) {
//          event.stopPropagation()
//          // textContent may return whitepace before/after tag. Strip it
//          let tag = event.currentTarget.textContent.replace(/^\s+|\s+$/g, '')
//          this.$store.commit('admin/setTagFilter', tag)
//          this.$store.dispatch("admin/getSubmissions")
//       },
//       removeFilter() {
//          this.$store.commit('admin/setTagFilter', "")
//          this.$store.dispatch("admin/getSubmissions")
//       }
//    },
//    created() {
//       this.$store.commit("clearSubmissionDetail")
//       this.$store.dispatch("admin/getSubmissions")
//    },
// };
</script>

<style scoped>
.tag-filter {
   font-size: 0.8em;
   border: 1px solid #ccc;
   padding: 2px 4px 0 10px;
   border-radius: 20px;
   cursor: pointer;
}
span.tag {
   display: inline-block;
   margin: 0 4px 4px 0;
   font-size: 0.9em;
   background: #0078e7;
   color: white;
   padding: 2px 10px 2px 10px;
   font-weight: 500;
   opacity: 0.6;
}
span.tag:hover {
   cursor: pointer;
   opacity: 1;
}
i.fas.unfilter {
   color: firebrick;
   margin-left: 5px;
   opacity: 0.6;
}
i.fas.unfilter:hover {
   cursor: pointer;
   opacity: 1;
}
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
  display: inline-block;
  margin-right: 10px;
}
#search {
   border: 1px solid #ccc;
   padding: 2px 4px;
   outline: none;
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