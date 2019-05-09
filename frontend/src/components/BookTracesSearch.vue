<template>
   <div v-if="showSearch===true" class="searchbar pure-form">
      <input type="text" id="search" @input="updateSearchQuery" @keyup.enter="doSearch" >
      <button @click="doSearch" class="search pure-button pure-button-primary">
         Search
      </button>
   </div>
</template>

<script>
import { mapState } from 'vuex'
export default {
   computed: {
      ...mapState({
         showSearch: state => state.public.showSearch
      }),
   },
   methods: {
      updateSearchQuery (e) {
         this.$store.commit('public/updateSearchQuery', e.target.value)
      },
      doSearch() {
         this.$store.dispatch("public/search")
         this.$router.push("/results")
      }
   }
}
</script>

<style scoped>
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
   border-radius: 0;
   border: none;
}
</style>
