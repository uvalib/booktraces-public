<template>
   <div class="institutions">
      <select @change="doInstitutionSearch" v-model="targetInstitution">
         <option value="" disabled>Browse by institution</option>
         <option v-for="(institution) in institutions"
            :key="institution.id" :value="institution.id">
            {{ institution.name }}
         </option>
      </select>
   </div>
</template>

<script>
import { mapState } from 'vuex'
import { mapFields } from 'vuex-map-fields'
export default {
   computed: {
      ...mapState({
         institutions: state => state.institutions,
      }),
      ...mapFields('public',[
           'targetInstitution'
      ]),
   },
   methods: {
      doInstitutionSearch() {
         this.$store.dispatch("public/searchInstitutions" )
         let curr = this.$router.currentRoute 
         if (curr.name != "results") {
            this.$router.push("/results")   
         }
      }
   }
}
</script>

<style scoped>

</style>
