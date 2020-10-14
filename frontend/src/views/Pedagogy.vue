<template>
   <div class="content pedagogy">
      <h2>Pedagogy</h2>
      <BTSpinner v-if="loading==true" message="Loading pedagogy content..." />
      <div v-else class="pedagogy-content">
         <template v-if="document">
            <h3 v-if="document.key != 'index'">{{document.title}}</h3>
            <div class="text" v-html="document.content"></div>
         </template>
         <div class="not-found" v-else>
            <h4>Page Not Found!</h4>
            <div>
               Sorry, this page cannot be found.
            </div>
         </div>
      </div>
   </div>
</template>

<script>
import { mapState } from 'vuex'
export default {
   name: "pedagogy",
   computed: {
      ...mapState({
         document: state => state.pedagogy.document,
         loading: state => state.loading
      })
   },
   mounted: async function () {
      let docKey = this.$route.params.id
      if (!docKey) {
         docKey = "index"
      }
      await this.$store.dispatch('pedagogy/get', docKey)
   },
};
</script>

<style lang="scss" scoped>
.content.pedagogy {
   h3 {
      margin: 0;
      position: relative;
      padding: 5px 10px;
      background: #666;
      color: white;
   }
   div.text {
      padding: 10px;
   }
}
::v-deep p,  ::v-deep ul{
   margin: 0 0;
}
::v-deep div.text  a {
   color: #24890d !important;
   font-weight: 500 !important;
   text-decoration: none !important;
   cursor: pointer !important;
}
::v-deep div.text  a:hover {
   text-decoration: underline !important;
}
</style>