<template>
   <div v-if="submissionsStore.showSearch===true" class="searchbar">
      <InputText id="search" v-model="submissionsStore.query" @keyup.enter="doSearch" />
      <Button @click="doSearch" label="Search" severity="info" />
   </div>
</template>

<script setup>
import { useSubmissionsStore } from "@/stores/submissions"
import { useRouter } from 'vue-router'
import InputText from 'primevue/inputtext'

const submissionsStore = useSubmissionsStore()
const router = useRouter()

const doSearch = (() => {
   submissionsStore.search()
   router.push("/results")
})

</script>

<style scoped>
.searchbar {
   background: #555;
   padding: 10px;
   position: absolute;
   right: 0;
   left: 0;
   z-index: 1000;
   border-bottom: 2px solid black;
   display: flex;
   flex-flow: row nowrap;
   gap: 5px;
   align-items: stretch;
   justify-content: flex-end;
   input {
      flex:1;
   }
}
</style>
