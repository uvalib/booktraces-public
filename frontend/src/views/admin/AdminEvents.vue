<template>
   <div class="admin">
      <h2>System Admin Panel <span class="login"><b>Logged in as:</b>{{loginName}}</span></h2>
      <div>
         <h3>Events</h3>
         <table class="events pure-table">
            <thead>
               <th>Date</th>
               <th>Description</th>
               <th style="width:50px;"></th>
            </thead>
            <tr v-for="event in events" :key="event.id">
               <td class="date">{{event.date}}</td>
               <td><span v-html="event.description"></span></td>
               <td class="centered">
                  <i :data-id="event.id" title="Delete" class="action fas fa-trash-alt" @click="deleteClicked"></i>
                  <i :data-id="event.id" title="Delete" class="action fa fa-edit" @click="editClicked"></i>
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
   name: "admin-events",
   computed: {
      ...mapState({
         error: state => state.error,
         loading: state => state.loading,
         events: state => state.events,
      }),
      ...mapGetters({
         loginName: 'admin/loginName',
      })
   },
   methods: {
      deleteClicked(event) {
         let tgt = event.currentTarget
         let eventID = tgt.dataset.id
         let resp = confirm("Delete this event? All data will be permanently lost. Continue?")
         if (resp) {
             this.$store.dispatch('admin/deleteEvent', eventID)
         }
      },
      editClicked(event) {
         let tgt = event.currentTarget
         let eventID = tgt.dataset.id
         alert("edit "+eventID)
      }
   },
   created() {
      this.$store.dispatch('getEvents')
   },
};
</script>

<style scoped>
i.action {
   color: #666;
   opacity: 0.6;
   margin: 0 5px;
}
i.action:hover {
   cursor: pointer;
   opacity: 1;
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
td.centered {
   text-align: center;
}
</style>