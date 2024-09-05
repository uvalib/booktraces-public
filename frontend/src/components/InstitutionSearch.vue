<template>
   <div class="institutions">
      <select @change="doInstitutionSearch" v-model="submissionsStore.targetInstitution">
         <option value="" disabled>Browse by institution</option>
         <option v-for="(institution) in submissionsStore.institutions"
            :key="institution.id" :value="institution.id">
            {{ institution.name }}
         </option>
      </select>
   </div>
</template>

<script setup>
import { useSubmissionsStore } from "@/stores/submissions"
import { useRouter } from 'vue-router'

const router = useRouter()
const submissionsStore = useSubmissionsStore()

const doInstitutionSearch = (() => {
   submissionsStore.searchInstitutions()
   router.push("/results")
})

</script>

<style scoped>
</style>
