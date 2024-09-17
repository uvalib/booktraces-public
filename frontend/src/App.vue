<template>
   <ConfirmDialog position="top" :closable="false">
      <template #message="slotProps">
         <div class="confirm">
            <i :class="slotProps.message.icon" class="warn"></i>
            <div v-html="slotProps.message.message"></div>
         </div>
      </template>
   </ConfirmDialog>
   <Toast position="top-center" />
   <BookTracesHeader/>
   <template v-if="system.adminMode">
      <router-view/>
   </template>
   <template v-else>
      <div class="main">
         <BookTracesSidebar/>
         <div class="content-wrap">
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
import Toast from 'primevue/toast'

const system = useSystemStore()

</script>

<style lang="scss">
@media only screen and (max-width: 1050px) {
   #app .main {
      width: 100%;
   }
   #app div.content {
      margin: 25px;
   }
}
@media only screen and (max-width: 768px) {
   #app div.content {
      margin: 0 !important;
   }
}
html,
body {
   margin: 0;
   padding: 0;
   background-color: #222;
}
#app .main {
   padding-top:0;
   display: flex;
   flex-flow: row nowrap;
   .content-wrap {
      flex:1;
      display: flex;
      flex-direction: column;
      background-image: url(./assets/books-bkg.jpg);
      position: relative;
      .content {
         flex:1;
         position: relative;
         padding: 3vw 3vw 2vw 3vw;
         min-height: 600px;
         background: white;
         font-weight: 400;
         color: #444;
         border: 1px solid #222;
         margin: 25px;
      }
   }
}
#app {
   font-family: "Avenir", Helvetica, Arial, sans-serif;
   -webkit-font-smoothing: antialiased;
   -moz-osx-font-smoothing: grayscale;
   background-color: #222;
   border-right: 1px solid #222;
}
#app h2 {
   margin: 0;
   position: relative;
}
a {
   color: #24890d;
   font-weight: 500;
   text-decoration: none;
}
a:hover {
   text-decoration: underline;
}
#app a.back {
   position: absolute;
   top: 15px;
   left: 15px;
}
#app h1 {
   margin: 0 0 5px 0;
}
#app .content h2 {
   font-size: 1.5em;
   font-weight: bold;
   border-bottom: 1px dashed #666;
   font-family: 'Special Elite', cursive;
   padding-bottom: 5px;
   margin-bottom: 15px;
}
td.nowrap {
   white-space: nowrap;
}
.confirm {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-start;
   align-items: flex-start;
   gap: 20px;
   .warn {
      width: 32px;
      height: 32px;
      font-size: 32px;
   }
}
</style>

