<template>
   <div class="institutions">
      <select @change="doInstitutionSearch" v-model="submissionsStore.targetInstitution">
         <option value="" disabled>Browse by institution</option>
         <option v-for="(institution) in system.institutions"
            :key="institution.id" :value="institution.id">
            {{ institution.name }}
         </option>
      </select>
   </div>
</template>

<script setup>
import { useSubmissionsStore } from "@/stores/submissions"
import { useSystemStore } from "@/stores/system"
import { useRouter } from 'vue-router'

const router = useRouter()
const submissionsStore = useSubmissionsStore()
const system = useSystemStore()

const doInstitutionSearch = (() => {
   submissionsStore.searchInstitutions()
   router.push("/results")
})

</script>

<style scoped>
select {
   padding: 5px 10px;
   border-color: #999;
   border-radius: 4px;
}
</style>
