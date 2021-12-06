<template>
   <div class="home content">
      <div class="bt-banner">
         <h1>Book Traces</h1>
      </div>
      <div class="info">
         <h3>Find unique copies of 19th- and early 20th-century books on library shelves</h3>
         <p><router-link to="/submit"><button class="submit pure-button pure-button-primary">SUBMIT A BOOK</button></router-link></p>
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
            <span>Total Submission: {{total}}</span>
             <InstitutionSearch style="margin-left:auto" />
         </div>
         <div class="pure-g thumbs">
            <div class="pure-u-sm-1-3 pure-u-md-1-4 pure-u-lg-1-5  pure-u-xl-1-5" v-for="thumb in thumbs" :key="thumb.submissionID">
               <router-link :to="thumbURL(thumb.submissionID)"><img class="pure-img thumb" :src="thumb.url"/></router-link>
            </div>
         </div>
         <h4 v-if="loading===true">Loading...</h4>
         <button v-if="thumbsCount < total" @click="moreClicked" class="more pure-button pure-button-primary">More</button>
      </div>
   </div>
</template>

<script>
import { mapState } from 'vuex'
import { mapGetters } from 'vuex'
import InstitutionSearch from "@/components/InstitutionSearch"
export default {
   name: "home",
   components: {
      InstitutionSearch
   },
   computed: {
      ...mapState({
         total: state => state.public.totalSubmissions,
         thumbs: state => state.public.thumbs,
         loading: state => state.loading,
      }),
      ...mapGetters({
         thumbsCount: 'public/thumbsCount',
      }),
   },
   methods: {
      thumbURL(id) {
         return "/submissions/"+id
      },
      adminClicked() {
         window.location.href = "/authenticate?url=/admin"
      },
      moreClicked() {
         this.$store.dispatch("public/getRecentThumbs")
      }
   },
   created() {
      this.$store.commit('public/clearThumbs' )
      this.$store.dispatch("public/getRecentThumbs")
      this.$store.dispatch('getInstitutions')
      this.$store.commit("transcribe/cancel")
   }
};
</script>

<style scoped>
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
.pure-g.thumbs a {
   display: inline-block;
   padding:5px;
}
.pure-img.thumb:hover {
   box-shadow: 0 0 5px green;
}
.pure-button.submit, .pure-button.more  {
   background: #24890d;
   font-weight: bold;
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
.controls {
   display:flex;
   flex-flow: row wrap;
   align-items: center;
   border-bottom: 1px dashed #666;
   padding-bottom: 10px;
   margin-bottom: 10px;
}
</style>