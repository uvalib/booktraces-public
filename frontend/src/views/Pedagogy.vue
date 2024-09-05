<template>
   <div class="content pedagogy">
      <h2>
         <span>Pedagogy</span>
         <span v-if="system.loading==false && pedagogy.document.key != 'index'" class="back-link">
            <router-link to="/pedagogy"><i class="fas fa-arrow-left"></i>Back</router-link>
         </span>
      </h2>
      <BTSpinner v-if="system.loading==true" message="Loading pedagogy content..." />
      <div v-else class="pedagogy-content">
         <template v-if="pedagogy.document != null">
            <h3 v-if="pedagogy.document.key != 'index'">{{pedagogy.document.title}}</h3>
            <div @click="docClicked" class="text" v-html="pedagogy.document.content"></div>
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

<script setup>
import { onMounted } from 'vue'
import { usePedagogyStore } from "@/stores/pedagogy"
import { useSystemStore } from "@/stores/system"
import { useRoute, useRouter, onBeforeRouteUpdate } from 'vue-router'

const system = useSystemStore()
const pedagogy = usePedagogyStore()
const route = useRoute()
const router = useRouter()

onBeforeRouteUpdate( async ( to ) => {
   // this is needed to load details when a grouped image thumb has been clicked; new content
   // needs to be loaded, but the page remains the same (create not called)
   await getDocument( to.params.id)
})

onMounted( async () => {
   await getDocument( route.params.id)
})

const getDocument = ( async (docKey) => {
   if (!docKey) {
      docKey = "index"
   }
   await pedagogy.getDocument( docKey )
})

const docClicked = ((event) => {
   if (event.target.className == 'pedagogy-link') {
      let docID = event.target.dataset.link
      if (docID) {
         router.push(`/pedagogy/${docID}`)
      }
   }
})

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
:deep(p),  :deep(ul){
   margin: 0 0;
}
:deep(div.text) {
   .pedagogy-link {
      color: #24890d !important;
      font-weight: 500 !important;
      text-decoration: none !important;
      cursor: pointer !important;
      &:hover {
         text-decoration: underline !important;
      }
   }
   a {
      color: #24890d !important;
      font-weight: 500 !important;
      text-decoration: none !important;
      cursor: pointer !important;
   }
   a:hover {
      text-decoration: underline !important;
   }
}
:deep(.ql-align-center) {
   text-align: center;
}
:deep(.ql-align-right) {
   text-align: right;
}
:deep(blockquote) {
  border-left: 6px solid #ccc;
  margin: 0;
  padding-left: 25px;
}
:deep(.ql-indent-1) {
   margin-left: 20px;
}
:deep(.ql-indent-2) {
   margin-left: 40px;
}
:deep(.ql-indent-3) {
   margin-left: 60px;
}
:deep(.ql-indent-4) {
   margin-left: 80px;
}
:deep(.ql-indent-5) {
   margin-left: 100px;
}
</style>