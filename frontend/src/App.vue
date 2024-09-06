<template>
   <BookTracesHeader/>
   <template v-if="system.adminMode">
      <router-view/>
   </template>
   <template v-else>
      <div class="main pure-g">
         <div class="pure-u-1-6">
            <BookTracesSidebar/>
         </div>
         <div class="pure-u-5-6 bkg">
            <BookTracesSearch/>
            <router-view/>
         </div>
      </div>
   </template>
   <BookTracesFooter/>
</template>

<script setup>
import BookTracesSearch from "@/components/BookTracesSearch.vue"
import BookTracesSidebar from "@/components/BookTracesSidebar.vue"
import BookTracesHeader from "@/components/BookTracesHeader.vue"
import BookTracesFooter from "@/components/BookTracesFooter.vue"
import { useSystemStore } from "@/stores/system"
import { useSubmissionsStore } from "@/stores/submissions"
import { onBeforeMount } from 'vue'

const system = useSystemStore()
const submissionsStore = useSubmissionsStore()

onBeforeMount( () => {
   submissionsStore.getArchiveDates()
   submissionsStore.getRecentSubmissions()
})
</script>

<style>
@media only screen and (max-width: 1050px) {
   #app .main div.pure-u-5-6.bkg {
      width: 100%;
   }
}
html,
body {
   margin: 0;
   padding: 0;
   /* background-image: url(./assets/main-bkg.jpg); */
   background-color: black;
}
#app .main {
   padding-top:0;
}
#app {
   font-family: "Avenir", Helvetica, Arial, sans-serif;
   -webkit-font-smoothing: antialiased;
   -moz-osx-font-smoothing: grayscale;
   background-color: black;
   border-right: 1px solid black;
}
#app h2 {
   margin: 0;
   position: relative;
}
#app a {
   color: #24890d;
   font-weight: 500;
   text-decoration: none;
}
#app a:hover {
   text-decoration: underline;
}
#app a.back {
   position: absolute;
   top: 15px;
   left: 15px;
}
div.bkg {
   background-image: url(./assets/books-bkg.jpg);
   position: relative;
}
#app h1 {
   margin: 0 0 5px 0;
}
#app div.content {
   position: relative;
   padding: 3vw 3vw 2vw 3vw;
   min-height: 600px;
   background: white;
   font-weight: 400;
   color: #444;
   border: 1px solid black;
   /* width: 80%; */
   margin: 25px;
}
#app .content h2 {
   font-size: 1.5em;
   font-weight: bold;
   border-bottom: 1px dashed #666;
   font-family: 'Special Elite', cursive;
   padding-bottom: 5px;
   margin-bottom: 15px;
}
</style>

