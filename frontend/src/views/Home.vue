<template>
   <div class="home content">
      <div class="bt-banner">
         <h1>Book Traces</h1>
      </div>
      <div class="info">
         <h3>Find unique copies of 19th- and early 20th-century books on library shelves</h3>
         <p>
            <Button as="router-link" label="SUBMIT A BOOK" to="/submit" />
         </p>
         <p>
            Thousands of <b>old library books</b> bear fascinating traces of the past.
            <b>Readers wrote in their books, and left pictures, letters, flowers, locks of hair, and other
            things between their pages.</b> We need your help identifying them in the stacks of academic
            libraries.  Together we can find out more about what books were and how they were used by
            their original owners, while also proving the value of maintaining rich print collections in our
            libraries.
         </p>
         <p>
            CURRENTLY COLLECTING IMAGES and CITATIONS of MARKED COPIES OF <b>LIBRARY
            BOOKS PUBLISHED BEFORE 1923.</b>  We are focusing on CIRCULATING AND RESEARCH COLLECTIONS
            (<b>not</b> rare books or special collections).  <b>Join the search!</b>
         </p>
         <p>
            Download the Data Set for Book Traces: Civil War Era Readers and Their Books in Virginia Libraries from UVA's LibraData:<br/>
            <a href="https://dataverse.lib.virginia.edu/dataset.xhtml?persistentId=doi:10.18130/V3/XVEQMH" target="_blank">
               https://dataverse.lib.virginia.edu/dataset.xhtml?persistentId=doi:10.18130/V3/XVEQMH
            </a>
         </p>
         <p>
            Explore our CLIR-funded project at the University of Virginia:<br/>
            <a href="http://booktraces.library.virginia.edu" target="_blank">
               http://booktraces.library.virginia.edu
            </a>
         </p>
         <p>
            Search the Book Traces database of over 12,000 University of Virginia books:<br/>
            <a href="https://booktraces.lib.virginia.edu/" target="_blank">
               https://booktraces.lib.virginia.edu/
            </a>
         </p>
         <div>
            <div class="book-info"><strong>Now Available!</strong></div>
            <div class="book">
               <div class="cover">
                  <a href="https://www.upenn.edu/pennpress/book/16176.html" target="_blank">
                     <img src="../assets/booktraces_book.jpg"/>
                  </a>
               </div>
               <div class="blurb">
                  <div>
                     <a href="https://www.upenn.edu/pennpress/book/16176.html" target="_blank">
                        Book Traces: Nineteenth-Century Readers and the Future of the Library
                     </a>
                  </div>
                  <p>
                     In Book Traces, Andrew M. Stauffer adopts what he calls "guided serendipity" as a tactic in pursuit of two goals:
                     first, to read nineteenth-century poetry through the clues and objects earlier readers left in their books and, second,
                     to defend the value of keeping the physical volumes on the shelves. Finding in such books of poetry the inscriptions, annotations,
                     and insertions made by their original owners, and using them as exemplary case studies, Stauffer shows how the physical,
                     historical book enables a modern reader to encounter poetry through the eyes of someone for whom it was personal.
                  </p>
               </div>
            </div>
         </div>
      </div>
      <div class="recents">
         <h3>Recently Submitted Books</h3>
         <div class="controls">
            <span>Total Submissions: {{submissionsStore.total}}</span>
             <InstitutionSearch style="margin-left:auto" />
         </div>
         <div class="thumbs">
            <img v-for="thumb in submissionsStore.thumbs" :key="thumb.submissionID"
               class="thumb" :src="thumb.url" @click="thumbClicked(thumb)" />
         </div>
         <h4 v-if="system.loading===true">Loading...</h4>
         <div class="more" v-if="submissionsStore.thumbsCount < submissionsStore.total">
            <Button  @click="moreClicked" label="Load More"/>
         </div>
      </div>
   </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useSubmissionsStore } from "@/stores/submissions"
import { useSystemStore } from "@/stores/system"
import InstitutionSearch from "@/components/InstitutionSearch.vue"
import { useRouter } from 'vue-router'

const router = useRouter()
const submissionsStore = useSubmissionsStore()
const system = useSystemStore()

const thumbClicked = ((thumb) => {
  router.push(`/submissions/${thumb.submissionID}`)
})
const moreClicked = (() => {
   submissionsStore.getRecentThumbs()
})

onMounted(() => {
   system.getInstitutions()
   submissionsStore.clearThumbs()
   submissionsStore.getRecentThumbs()
   submissionsStore.getArchiveDates()
   submissionsStore.getRecentSubmissions()
})
</script>

<style scoped lang="scss">
h3, h1 {
   font-family: 'Special Elite', cursive;
   margin-bottom: 5px;
}
div.book-info {
   font-weight: bold;
   font-size: 1.1em;
   margin: 25px 0 10px 0;
}
div.book {
   display: flex;
   flex-flow: row wrap;
   justify-content: center;
}
div.book img {
   flex: 1 1 10%;
   max-height: 200px;
   margin-bottom: 20px;
}
div.blurb {
   margin-left: 20px;
   flex: 1 1 70%;
}
div.book a {
   margin: 0 !important;
}
div.blurb p {
   padding: 0;
   margin: 10px 0;
}
h3 {
  margin-bottom: 10px;
  font-weight: 500;
}
.bt-banner {
   position: relative;
}
div.info {
   margin-bottom: 30px;
}
div.info a.inline {
   margin-left:0;
}
div.info a {
   margin-left: 20px;
}
.recents {
   display: flex;
   flex-direction: column;
   gap: 15px;
   .thumbs {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(0, min(100%/2, max(150px, 100%/6))));
      grid-gap: 1.25rem;
      padding: 0;
      justify-content: center;
      .thumb {
         position: relative;
         width: 150px;
         height: 150px;
         background-color: #efefef;
         &:hover {
            top: -2px;
            left: -2px;
            box-shadow: 0 3px 6px #666;
         }
      }
   }
   .more {
      text-align: center;
   }
   .controls {
      display:flex;
      flex-flow: row wrap;
      align-items: center;
      border-bottom: 1px dashed #666;
      padding-bottom: 10px;
      margin-bottom: 10px;
   }
}
</style>