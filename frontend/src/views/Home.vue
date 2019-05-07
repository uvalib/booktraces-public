<template>
   <div class="home content">
      <div class="bt-banner">
         <h1>Book Traces</h1>
         <span class="admin">
            <button @click="adminClicked" class="pure-button pure-button-primary">Admin Access</button>
         </span>
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
            Explore our CLIR-funded project at the University of Virginia:<br/>
            <a href="http://booktraces.library.virginia.edu" target="_blank">
               http://booktraces.library.virginia.edu
            </a>
         <p>
            Search the Book Traces database of over 12,000 University of Virginia books:<br/>
            <a href="https://booktraces.lib.virginia.edu/" target="_blank">
               https://booktraces.lib.virginia.edu/
            </a>
         </p>
      </div>
      <div class="recents">
         <h3>Recently Submitted Books</h3>
         <p>Total Submission: {{total}}</p>
         <div class="pure-g thumbs" style="text-align:center;">
            <div class="pure-u-sm-1-3 pure-u-md-1-4 pure-u-lg-1-5  pure-u-xl-1-5" v-for="thumb in submissions" :key="thumb.submissionID">
               <router-link :to="thumbURL(thumb.submissionID)"><img class="pure-img thumb" :src="thumb.url"/></router-link>
            </div>
         </div>
      </div>
   </div>
</template>

<script>
import { mapState } from 'vuex'
export default {
   name: "home",
   computed: {
      ...mapState({
         total: state => state.public.totalSubmissions,
         submissions: state => state.public.submissions,
      }),
   },
   methods: {
      thumbURL(id) {
         return "/submissions/"+id
      },
      adminClicked() {
         window.location.href = "/authenticate?url=/admin"
      }
   },
   created() {
      this.$store.dispatch("public/getRecentThumbs")
   }
};
</script>

<style scoped>
h3, h1 {
   font-family: 'Special Elite', cursive;
   margin-bottom: 5px;
}
h3 {
  margin-bottom: 20px;
  font-weight: 500;
}
.pure-g.thumbs a {
   display: inline-block;
   padding:5px;
}
.pure-img.thumb:hover {
   box-shadow: 0 0 5px green;
}
.pure-button.submit {
   background: #24890d;
   font-weight: bold;
}
.bt-banner {
   position: relative;
}
span.admin {
   font-size: 11px;
   position: absolute;
   right: 0;
   top:0;
}
div.info {
   margin-bottom: 30px;
}
div.info a {
   margin-left: 20px;
}
</style>