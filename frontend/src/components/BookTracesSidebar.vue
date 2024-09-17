<template>
   <div class="sidebar">
      <p class="subtitle">Find unique copies of old library books</p>
      <div class="section">
         <p class="subtitle pad">Recent Submissions</p>
         <div class="recent" v-for="recent in submissionsStore.recents" :key="recent.id">
            <div @click="recentClicked(recent.id)" class="recent-link">
               <p class="title">{{recent.title}}</p>
               <p>{{recent.submittedAt}}</p>
            </div>
         </div>
      </div>
      <div class="section">
         <p class="subtitle pad">Archives</p>
         <div class="archive" v-for="archive in submissionsStore.archives">
            <span @click="archiveClicked(archive.internalDate)" class="archive-date">
               {{archive.displayDate}}
            </span>
         </div>
      </div>
   </div>
</template>

<script setup>
import { useSubmissionsStore } from "@/stores/submissions"
import { useDetailsStore } from "@/stores/details"
import { useRouter } from 'vue-router'

const details = useDetailsStore()
const submissionsStore = useSubmissionsStore()
const router = useRouter()

const recentClicked = ((id) => {
   router.push("/submissions/" + id)
   details.getSubmission(id)
})

const archiveClicked = ((tgtDate) => {
   submissionsStore.getArchive(tgtDate)
   router.push("/results")
})

</script>

<style scoped lang="scss">
@media only screen and (max-width: 1050px) {
   .sidebar {
      display:none;
   }
}
@media only screen and (min-width: 1050px) {
   .sidebar {
      display: flex;
      flex-direction: column;
      gap: 25px;
   }
}
div.sidebar {
   padding: 10px;
   background-color: #222;
   color: #dadada;
   width: 18%;
   a, .recent-link, .archive-date {
      color: #dadada;
      text-decoration: none;
      &:hover {
         text-decoration: underline;
         color: #34991d !important;
         cursor: pointer;
      }
   }
   .recent, .archive {
      margin: 0 0 10px 15px;
      font-size: 0.9rem;
   }
   .recent p {
      margin:0;
   }

   .subtitle {
      margin:0;
      padding:0;
      font-family: 'Special Elite', cursive;
      font-weight: 500;
   }
   .subtitle.pad {
      margin-bottom: 10px;
   }
}
</style>
