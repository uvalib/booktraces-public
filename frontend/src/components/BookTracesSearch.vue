<template>
   <div v-if="unAuth.showSearch===true" class="searchbar pure-form">
      <div class="controls">
         <div class="query">
            <input type="text" id="search" v-model="unAuth.query" @keyup.enter="doSearch" >
            <button @click="doSearch" class="search pure-button pure-button-primary">
               Search
            </button>
         </div>
      </div>
   </div>
</template>

<script setup>
import { useUnauthStore } from "@/stores/unauth"
import { useRouter } from 'vue-router'

const unAuth = useUnauthStore()
const router = useRouter()

const doSearch = (() => {
   unAuth.search()
   router.push("/results")
})

</script>

<style scoped>
.controls {
   display: inline-block;
}
.searchbar {
   background: #555;
   padding: 8px;
   font-size: 0.9em;
   text-align: right;
   position: absolute;
   right: 0;
   left: 0;
   z-index: 1000;
   border-bottom: 2px solid black;
}
#search {
   width: 350px;
   outline: none;
   border: none;
   border-radius: 5px 0 0 5px;
}
button {
   border-radius: 0 5px 5px 0;
}
</style>
