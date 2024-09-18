<template>
   <div class="content events">
      <h2>Events</h2>
      <div class="events">
         <div class="event" v-for="event in system.events">
            <span class="date nowrap">{{event.date}}</span>
            <span v-html="event.description"></span>
         </div>
      </div>
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

<style scoped lang="scss">
.events {
   width: 100%;
   display: flex;
   flex-direction: column;
   gap: 15px;
   div.event {
      display: flex;
      flex-flow: row nowrap;
      align-items: flex-start;
      gap: 15px;
      border-bottom: 1px solid #ddd;
      padding-bottom: 10px;
      .date {
         white-space: nowrap;
         font-weight: bold;
      }
   }
}
</style>