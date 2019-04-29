<template>
   <div class="admin">
      <h2>Book Traces System Admin Panel</h2>
      <div>
         <h3>Submissions</h3>
         <table class="pure-table">
            <thead>
               <th style="width:12px"/>
               <th>Title</th>
               <th>Author</th>
               <th>Tags</th>
               <th>Submitted</th>
               <th class="checkbox">Public</th>
               <th>Actions</th>
            </thead>
            <tr v-for="sub in submissions" :key="sub.id">
               <td><input :data-id="sub.id" type="checkbox"/></td>
               <td>{{ sub.title }}</td>
               <td>{{ sub.author }}</td>
               <td>{{ sub.tags }}</td>
               <td>{{ sub.submittedAt.split("T")[0] }}</td>
               <td class="centered"><span v-html="publishIcon(sub)"></span></td>
               <td>
                  <i title="view" class="action fas fa-eye"></i>
                  <i title="edit" class="action fas fa-edit"></i>
                  <i title="delete" class="action fas fa-trash-alt"></i>
                  <i v-bind:class="{disabled: isPublished(sub)}" title="publish" class="action fas fa-thumbs-up"></i>
               </td>
            </tr>
         </table>
      </div>
   </div>
</template>

<script>
import { mapState } from 'vuex'
export default {
   name: "admin",
   computed: {
      ...mapState({
         total: state => state.admin.totalSubmissions,
         submissions: state => state.admin.submissions,
      }),
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
      }
   },
   created() {
      let authUser = this.$cookies.get("bt_admin_user")
      if (authUser) {
         this.$store.commit("admin/setUser", authUser)
         this.$store.dispatch("admin/getSubmissions")
      } else {
         this.$router.push("forbidden")
      }
   }
};
</script>

<style scoped>
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
th.checkbox {
   width: 40px;
}
td.centered {
   text-align: center;
}
</style>