<template>
   <div class="content events">
      <h2>Events</h2>
      <table>
         <tbody>
            <tr v-for="event in system.events">
               <td class="date">{{event.date}}</td>
               <td><span v-html="event.description"></span></td>
            </tr>
         </tbody>
      </table>
   </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useSystemStore } from "@/stores/system"

const system = useSystemStore()

onMounted(() => {
   system.getEvents()
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