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
         <div class="archive" v-for="(archive,idx) in submissionsStore.archives" :key="idx">
            <span @click="archiveClicked(archive.dataset.date)" :data-date="archive.internalDate" class="archive-date">
               {{archive.displayDate}}
            </span>
         </div>
      </div>
   </div>
</template>

<script setup>
import { useSystemStore } from "@/stores/system"
import { useSubmissionsStore } from "@/stores/submissions"
import { useRouter } from 'vue-router'

const system = useSystemStore()
const submissionsStore = useSubmissionsStore()
const router = useRouter()

const recentClicked = ((id) => {
   router.push("/submissions/" + id)
   system.getSubmissionDetail(id)
})

const archiveClicked = ((tgtDate) => {
   submissionsStore.getArchive(tgtDate)
   router.push("/results")
})

</script>

<style scoped>
@media only screen and (max-width: 1050px) {
   .sidebar {
      display:none;
   }
}
@media only screen and (max-width: 1150px) {
   #tfeed {
      display:none;
   }
}
div.recent .recent-link:hover {
   cursor:pointer;
   color: #34991d !important;
}
div.section {
   margin-top: 20px;
   position: relative;
}
.recent, .archive {
   font-size: 0.9em;
   margin: 0 0 5px 8px;
}
.recent p {
   margin:0;
}
.recent p.title {
   font-weight: bold;
   font-style: italic;
}
div.sidebar {
   padding: 0 15px 15px 15px;
   background-color: black;
   color: #ccc;
   font-size: 0.9em;
}
#app div.sidebar a {
  color: #ccc;
  text-decoration: none;
}
.subtitle {
   margin:0;
   padding:0;
   font-family: 'Special Elite', cursive;
   font-weight: 500;
}
.subtitle.pad {
   margin-bottom: 5px;
}
.archive-date {
   cursor:pointer;
}
.archive-date:hover {
   font-weight: bold;
   text-decoration: underline;
}
</style>
