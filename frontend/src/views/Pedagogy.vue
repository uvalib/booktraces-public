<template>
   <div class="content pedagogy">
      <h2>
         <span>Pedagogy</span>
         <span v-if="loading==false && document.key != 'index'" class="back-link">
            <router-link to="/pedagogy"><i class="fas fa-arrow-left"></i>Back</router-link>
         </span>
      </h2>
      <BTSpinner v-if="loading==true" message="Loading pedagogy content..." />
      <div v-else class="pedagogy-content">
         <template v-if="document">
            <h3 v-if="document.key != 'index'">{{document.title}}</h3>
            <div @click="docClicked" class="text" v-html="document.content"></div>
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
   watch: {
        $route() {
           // this is needed to load details when a grouped image thumb has been clicked; new content
           // needs to be loaded, but the page remains the same (create not called)
           this.getDocument()
        }
   },
   methods: {
      async getDocument() {
         let docKey = this.$route.params.id
         if (!docKey) {
            docKey = "index"
         }
         await this.$store.dispatch('pedagogy/get', docKey)
      },
      docClicked(event) {
         if (event.target.className == 'pedagogy-link') {
            let docID = event.target.dataset.link
            if (docID) {
               this.$router.push(`/pedagogy/${docID}`)
            }
         }
      }
   },
   mounted() {
      this.getDocument()
   },
};
</script>

<style lang="scss" scoped>
.content.pedagogy {
   h2 {
      display: flex;
      flex-flow: row wrap;
      justify-content: space-between ;
      align-content: center;
      .back-link {
         font-size: 0.6em;;
      }
      span {
         display: inline-block;
      }
   }
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
::v-deep div.text  .pedagogy-link {
   color: #24890d !important;
   font-weight: 500 !important;
   text-decoration: none !important;
   cursor: pointer !important;
   &:hover {
      text-decoration: underline !important;
   }
}
::v-deep div.text  a:hover {
   text-decoration: underline !important;
}
</style>