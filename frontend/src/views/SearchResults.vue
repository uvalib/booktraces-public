<template>
   <div class="search content">
      <div class="bt-banner">
         <h1>Book Traces</h1>
         <h3 v-if="submissionsStore.searchType=='query'">Search Results for: {{ submissionsStore.query }}</h3>
         <h3 v-if="submissionsStore.searchType=='archive'">Submissions from: {{ submissionsStore.archiveDate }}</h3>
         <h3 v-if="submissionsStore.searchType=='institution'">Submissions from: {{ institutionName }}</h3>
         <h3 v-if="submissionsStore.searchType=='tag'">Submissions tagged: {{ submissionsStore.tgtTag }}</h3>
      </div>
      <h4 v-if="system.loading===true">Loading...</h4>
      <template v-else>
         <div v-if="submissionsStore.searchHitCount==0">
            <h4>Sorry, but nothing matched your search terms. Please try again with some different keywords.</h4>
         </div>
         <div v-else class="hits">
            <div class="controls">
               <span>Total Matches Found: {{ submissionsStore.searchHitCount }}</span>
               <InstitutionSearch />
            </div>
            <div class="hits">
               <div v-for="hit in submissionsStore.searchResults" :key="hit.id">
                  <router-link class="hit" :to="`/submissions/${hit.id}`">
                     <div class="hit">
                        <img class="thumb" :src="hit.url"/>
                        <div class="data">
                           <p><b>Title:</b> {{ hit.title }}</p>
                           <p><b>Submitted:</b> {{ hit.submittedOn }}</p>
                           <p><b>Institution:</b> {{ hit.institution }}</p>
                           <p><b>Tags: </b>{{ formatTags(hit.tags) }}</p>
                           <div class="desc">
                              <label>Description:</label>
                              <p class="indent">{{ formatDescription(hit.description) }}</p>
                           </div>
                        </div>
                     </div>
                  </router-link>
               </div>
            </div>
         </div>
      </template>
   </div>
</template>

<script setup>
import { onMounted, computed } from 'vue'
import InstitutionSearch from "@/components/InstitutionSearch.vue"
import { useSubmissionsStore } from "@/stores/submissions"
import { useSystemStore } from "@/stores/system"
import { useRouter } from 'vue-router'

const submissionsStore = useSubmissionsStore()
const system = useSystemStore()
const router = useRouter()

onMounted(() => {
   if (submissionsStore.searchType == "none") {
      router.replace("/")
   }
})

const institutionName = computed(() => {
   let tgt = system.institutions.find( i => i.id == submissionsStore.targetInstitution)
   if (tgt) {
      return tgt.name
   }
   return ""
})

const formatDescription = ((text) => {
   if (text.length < 250) {
      return text;
   }
   return text.substr(0,249)+"..."
})

const formatTags = ((text) => {
   if (!text) {
      return "None"
   }
   return text.split(",").join(", ")
})
</script>

<style scoped lang="scss">
@media only screen and (min-width: 768px) {
   div.hit {
      display: flex;
      flex-flow: row wrap;
      gap: 20px;
   }
}
@media only screen and (max-width: 768px) {
   div.hit {
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      gap: 10px;
   }
}
.bt-banner {
   h3, h1 {
      font-family: 'Special Elite', cursive;
      margin-bottom: 5px;
   }
   h3 {
      margin: 10px 0 10px 0;
      font-weight: 500;
      padding-bottom:10px;
      border-bottom: 1px dashed #666;
   }
}
div.small-img {
   text-align: center;
   display:none;
}
div.hits {
   display: flex;
   flex-direction: column;
   gap: 20px;
   a.hit {
      color: #444 !important;
      text-decoration: none !important;
   }
   div.hit {
      border: 1px solid #ccc;
      padding: 15px;
      border-radius: 4px;
      img.thumb {
         width: 150px;
         height: 150px;
      }
      div.data {
         flex:1;
         p {
            margin: 5px 0;
         }
         label {
            font-weight:bold;
         }
         p.indent {
            margin-left: 25px;
         }
      }
   }
   .controls {
      display:flex;
      flex-flow: row wrap;
      align-items: center;
      justify-content: space-between;
   }
}
</style>