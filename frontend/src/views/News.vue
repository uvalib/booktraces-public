<template>
   <div class="content events">
      <h2>News</h2>
      <div class="news-item" v-for="item in system.publishedNews" >
         <h3>
            <span>{{item.title}}</span>
            <span class="date">{{formatDate(item.createdAt)}}</span>
         </h3>
         <div class="text" v-html="item.content"></div>
      </div>
   </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useSystemStore } from "@/stores/system"

const system = useSystemStore()

onMounted(() => {
   system.getNews()
})

const formatDate = ((date) => {
   return date.split("T")[0]
})
</script>

<style scoped lang="scss">
.news-item {
   margin: 10px 0 25px 0;
   padding: 0px;
}
h3 {
   margin: 0;
   position: relative;
   padding: 5px 10px;
   background: #666;
   color: white;
   display: flex;
   flex-flow: row wrap;
   justify-content: space-between;
   align-items: flex-start;
   gap: 10px;

   span.date {
      font-weight: 100;
      font-size: 0.8em;
   }
}
div.text {
   padding: 10px;
   :deep(p) {
      margin: 0 0;
   }
}
</style>