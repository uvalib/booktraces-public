<template>
   <div class="admin">
      <h2>Book Traces System Admin Panel <span class="login"><b>Logged in as:</b>{{loginName}}</span></h2>
      <div>
         <h3>Submissions</h3>
         <table class="pure-table">
            <thead>
               <th>Title</th>
               <th>Author</th>
               <th>Tags</th>
               <th>Submitted</th>
               <th class="checkbox">Public</th>
               <th>Actions</th>
            </thead>
            <tr v-for="sub in submissions" :key="sub.id" :data-id="sub.id" @click="submissionClicked">
               <td>{{ sub.title }}</td>
               <td>{{ sub.author }}</td>
               <td>{{ sub.tags }}</td>
               <td>{{ sub.submittedAt.split("T")[0] }}</td>
               <td class="centered"><span v-html="publishIcon(sub)"></span></td>
               <td>
                  <i title="delete" class="action fas fa-trash-alt" @click="deleteClicked"></i>
                  <i v-bind:class="{disabled: isPublished(sub)}" title="publish" class="action fas fa-thumbs-up" @click="publishClicked"></i>
               </td>
            </tr>
         </table>
      </div>
   </div>
</template>

<script>
import { mapState } from 'vuex'
import { mapGetters } from 'vuex'
export default {
   name: "admin",
   computed: {
      ...mapState({
         total: state => state.admin.totalSubmissions,
         submissions: state => state.admin.submissions,
      }),
      ...mapGetters({
         loginName: 'admin/loginName',
      })
   },
   methods: {
      isPublished(sub) {
         return sub.published === 1
      },
      publishIcon(sub) {
         if (sub.published === 1) {
            return '<i style="color:green" class="fas fa-check-circle"></i>'
         } else {
            return ""
         }
      },
      submissionClicked(event) {
         let tgt = event.currentTarget
         let subID = tgt.dataset.id
         this.$router.push("admin/submissions/" + subID)
      },
      deleteClicked(event) {
         event.stopPropagation()
         confirm("Delete this submission?")
      },
      publishClicked(event) {
         event.stopPropagation()
         if (event.currentTarget.classList.contains("disabled")) {
            return
         }
         confirm("Publish this submission?")
      }
   },
   created() {
      this.$store.dispatch("admin/getSubmissions")
   }
};
</script>

<style scoped>
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