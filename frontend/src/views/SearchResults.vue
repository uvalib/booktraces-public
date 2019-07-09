<template>
   <div class="search content">
      <div class="bt-banner">
         <h1>Book Traces</h1>
         <h3 v-if="searchType=='query'">Search Results for: {{query}}</h3>
         <h3 v-if="searchType=='archive'">Submissions from: {{archiveDate}}</h3>
         <h3 v-if="searchType=='tag'">Submissions tagged: {{tgtTag}}</h3>
      </div>
      <h4 v-if="loading===true">Loading...</h4>
      <template v-else>
         <div v-if="searchHitCount==0">
            <h4>Sorry, but nothing matched your search terms. Please try again with some different keywords.</h4>
         </div>
         <div v-else class="hits">
            <p>Total Matches Found: {{searchHitCount}}</p>
            <div class="hits">
               <div v-for="hit in hits" :key="hit.id">
                  <router-link class="hit" :to="submissionURL(hit.id)">
                     <div class="hit pure-g">
                        <img class="pure-u-1-3 thumb" :src="hit.url"/>
                        <div class="pure-u-2-3 data">
                           <p><b>Title:</b> {{ hit.title }}</p>
                           <p><b>Submitted:</b> {{ hit.submittedOn }}</p>
                           <p><b>Tags: </b>{{ formatTags(hit.tags) }}</p>
                           <div class="desc">
                              <label>Description:</label>
                              <p class="indent">{{ formatDescription(hit.description) }}</p>
                           </div>
                           <div class="small-img">
                               <img class="thumb" :src="hit.url"/>
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

<script>
import { mapState } from 'vuex'
import { mapGetters } from 'vuex'
export default {
   name: "results",
   computed: {
      ...mapState({
         hits: state => state.public.searchResults,
         loading: state => state.loading,
         query: state => state.public.query,
         archiveDate: state => state.public.archiveDate,
         tgtTag: state => state.public.tgtTag
      }),
      ...mapGetters({
         searchHitCount: 'public/searchHitCount',
         searchType: 'public/searchType',
      }),
   },
   created() {
      if (this.searchType == "none") {
         this.$router.replace("/")
      }
   },
   methods: {
      submissionURL(id) {
         return "/submissions/"+id
      },
      formatDescription(text) {
         if (text.length < 250) {
            return text;
         }
         return text.substr(0,249)+"..."
      },
      formatTags(text) {
         if (!text) {
            return "None"
         }
         return text.split(",").join(", ")
      }
   },
};
</script>

<style scoped>
@media only screen and (max-width: 768px) {
   img.pure-u-1-3.thumb {
      display: none;
   }
   div.pure-u-2-3.data {
      width:100%;
   }
   div.desc {
      display:none;
   }
   div.small-img {
      display: block !important;
   }
}
div.small-img {
   text-align: center;
   display:none;
}
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
div.hits a.hit {
   color: #444 !important;
   text-decoration: none !important;
}
div.hit {
   border:1px solid #ccc;
   padding: 10px 0;
   margin: 5px 0;
   /* min-height:150px; */
   /* max-height:150px; */
   font-size: 0.9em;
}
img.thumb {
   max-width: 150px;
   max-height: 150px;
}
div.data {
   padding-left: 15px;
}
div.data p {
   margin: 5px 0;
}
div.data label {
   font-weight:bold;
}
div.data p.indent {
   margin-left: 25px;
}
</style>