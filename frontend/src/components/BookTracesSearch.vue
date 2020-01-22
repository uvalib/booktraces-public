<template>
   <div v-if="showSearch===true" class="searchbar pure-form">
      <div class="controls">
         <div class="institutions">
            <select @change="doInstitutionSearch" v-model="selectedInstitution">
               <option value="" disabled>Browse by institution</option>
               <option v-for="(institution) in institutions"
                  :key="institution.id" :value="institution.id">
                  {{ institution.name }}
               </option>
            </select>
         </div>
         <div class="query">
            <input type="text" id="search" @input="updateSearchQuery" @keyup.enter="doSearch" >
            <button @click="doSearch" class="search pure-button pure-button-primary">
               Search
            </button>
         </div>
      </div>
   </div>
</template>

<script>
import { mapState } from 'vuex'
export default {
   computed: {
      ...mapState({
         showSearch: state => state.public.showSearch,
         institutions: state => state.institutions,
      }),
   },
   data: function() {
      return {
         selectedInstitution: "",
      }
   },
   methods: {
      updateSearchQuery (e) {
         this.selectedInstitution = ""
         this.$store.commit('public/updateSearchQuery', e.target.value)
      },
      doSearch() {
         this.selectedInstitution = ""
         this.$store.dispatch("public/search")
         this.$router.push("/results")
      },
      doInstitutionSearch() {
         let tgt = this.institutions.find( i => i.id == this.selectedInstitution)
         this.$store.dispatch("public/searchInstitutions", tgt )
         this.$router.push("/results")   
      }
   }
}
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
.institutions select {
   width: 100%;
   outline: none;
   border: none;
   margin-bottom: 10px;

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
