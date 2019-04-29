<template>
   <div class="content events">
      <h2>Events</h2>
      <table>
         <tr v-for="event in events" :key="event.id">
            <td class="date">{{event.date}}</td>
            <td><span v-html="event.description"></span></td>
         </tr>
      </table>
   </div>
</template>

<script>
export default {
   name: "events",
   computed: {
      events() {
         return this.$store.state.core.events
      }
   },
   created: function () {
      this.$store.dispatch('core/getEvents')
   },
   methods: {
   }
};
</script>

<style scoped>
td.date {
   padding: 5px 20px;
   text-align: right;
}
td {
border-bottom: 1px solid #ddd;
}
</style>