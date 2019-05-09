<template>
   <div class="sidebar">
      <p class="subtitle">Find unique copies of old library books</p>
      <div class="section">
         <p class="subtitle pad">Recent Submissions</p>
         <div class="recent" v-for="recent in recents" :key="recent.id">
            <router-link :to="thumbURL(recent.id)">
               <p class="title">{{recent.title}}</p>
               <p>{{recent.submittedAt}}</p>
            </router-link>
         </div>
      </div>
      <div class="section">
         <p class="subtitle pad">Archives</p>
         <div class="archive" v-for="(archive,idx) in archives" :key="idx">
            <span @click="archiveClicked" :data-date="archive.internalDate" class="archive-date">
               {{archive.displayDate}}
            </span>
         </div>
      </div>
   </div>
</template>

<script>
import { mapState } from 'vuex'
export default {
   computed: {
      ...mapState({
         recents: state => state.public.recents,
         archives: state => state.public.archives,
      }),
   },
    methods: {
      thumbURL(id) {
         return "/submissions/"+id
      },
      archiveClicked(event) {
         let tgtDate = event.currentTarget.dataset.date 
         this.$store.dispatch("public/getArchive", tgtDate)
         this.$router.push("/results")
      }
    }
}
</script>

<style scoped>
div.recent a:hover {
   color: #34991d !important;
}
div.section {
   margin-top: 20px;
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
