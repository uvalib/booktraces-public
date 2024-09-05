<template>
   <div class="content events">
      <h2>Events</h2>
      <table>
         <tr v-for="event in eventsStore.list">
            <td class="date">{{event.date}}</td>
            <td><span v-html="event.description"></span></td>
         </tr>
      </table>
   </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useEventsStore } from "@/stores/events"

const eventsStore = useEventsStore()

onMounted(() => {
   eventsStore.getAll()
})
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