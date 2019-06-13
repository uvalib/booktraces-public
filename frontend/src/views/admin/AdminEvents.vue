<template>
   <div class="admin">
      <h2>System Admin Panel <span class="login"><b>Logged in as:</b>{{loginName}}</span></h2>
      <div>
         <h3>Events <span @click="addEventClick" class="add pure-button pure-button-primary">Add Event</span></h3>
         <div class="error">{{error}}</div>
         <table class="events pure-table">
            <thead>
               <th>Date</th>
               <th>Description</th>
               <th style="width:50px;"></th>
            </thead>
            <tr v-for="event in events" :key="event.id">
               <template v-if="editingEvent(event.id)">
                  <td class="date"><input type="text" id="date" v-model="editDetails.date"></td>
                  <td><textarea id="description" v-model="editDetails.description"></textarea></td>
                  <td class="centered">
                     <i class="action cancel fas fa-times-circle" @click="cancelClicked"></i>
                     <i class="action save fas fa-check-circle" @click="saveClicked"></i>
                  </td>
               </template>
               <template v-else>
                  <td class="date">{{event.date}}</td>
                  <td><span v-html="event.description"></span></td>
                  <td class="centered">
                     <i :data-id="event.id" title="Delete" class="action fas fa-trash-alt" @click="deleteClicked"></i>
                     <i :data-id="event.id" title="Edit" class="action fa fa-edit" @click="editClicked"></i>
                  </td>
               </template>
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
   data: function () {
      return {
         edit: false,
         editDetails: null,
         addingNew: false,
      }
   },
   computed: {
      ...mapState({
         error: state => state.error,
         loading: state => state.loading,
         events: state => state.events.list,
      }),
      ...mapGetters({
         loginName: 'admin/loginName',
      })
   },
   methods: {
      cancelClicked() {
         this.edit = false
         this.editDetails = null
         if (this.addingNew) {
            this.addingNew = false
            this.$store.commit("events/cancelAddEvent")
         }
      },
      saveClicked() {
         if (this.addingNew) {
            this.$store.dispatch('events/addEvent',this.editDetails).then((/*response*/) => {
               this.edit=false
               this.editDetails = null
               this.addingNew = false
            })
         } else {
            this.$store.dispatch('events/updateEvent',this.editDetails).then((/*response*/) => {
               this.edit=false
               this.editDetails = null
               this.addingNew = false
            })
         }
      },
      editingEvent(id) {
         if (this.edit == false) return false 
         return this.editDetails.id == id
      },
      deleteClicked(event) {
         let tgt = event.currentTarget
         let eventID = tgt.dataset.id
         let resp = confirm("Delete this event? All data will be permanently lost. Continue?")
         if (resp) {
             this.$store.dispatch('events/deleteEvent', eventID)
         }
      },
      editClicked(event) {
         let tgt = event.currentTarget
         var eventID = tgt.dataset.id
         var evtIdx = -1
         this.events.some( function(e,idx) {
          if (e.id == eventID) {
            evtIdx = idx
            return true
          }
          return false
        })
        if (evtIdx > -1) {
           this.editDetails = Object.assign({}, this.events[evtIdx])
           this.edit = true
        }
      },
      addEventClick() {
         this.$store.commit("events/addEventPlaceholder")
         this.editDetails = Object.assign({}, this.events[0])
         this.edit = true
         this.addingNew = true
      },
   },
   created() {
      this.$store.dispatch('events/getAll')
   },
};
</script>

<style scoped>
i.action.cancel, .error {
   color: firebrick;
}
i.action.save {
   color: green;
}
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
table tr:hover {
   background: #f0f0f0;
}
table td input, table td textarea {
   width: 100%;
}
td.centered {
   text-align: center;
}
h3 {
   position: relative;
}
h3 span.pure-button.add {
   font-size: 0.6em;
   font-weight: 100;
   position: absolute;
   right: 0;
   bottom:0;
}
</style>