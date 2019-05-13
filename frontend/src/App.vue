<template>
   <div id="app">
      <BookTracesHeader/>
      <template v-if="adminMode">
         <router-view/>
      </template>
      <template v-else>
         <div class="main pure-g">
            <div class="pure-u-1-6">
               <BookTracesSidebar/>
            </div>
            <div class="pure-u-5-6 bkg">
               <BookTracesSearch/>
               <router-view/>
            </div>
         </div>
      </template>
      <BookTracesFooter/>
   </div>
</template>


<script>
import BookTracesSearch from "@/components/BookTracesSearch"
import BookTracesSidebar from "@/components/BookTracesSidebar"
import BookTracesHeader from "@/components/BookTracesHeader"
import BookTracesFooter from "@/components/BookTracesFooter"
import { mapGetters } from "vuex"
import { mapState } from 'vuex'

export default {
   components: {
      BookTracesSearch,
      BookTracesHeader,
      BookTracesSidebar,
      BookTracesFooter
   },
   computed: {
      ...mapGetters({
         isAuthenticated: "admin/isAuthenticated"
      }),
      ...mapState({
         adminMode: state => state.adminMode,
      })
   },
   data: function() {
      return {}
   },
};
</script>

<style>
.tweet-not-found {
   color: #fe6d6d;
   font-size: 0.9em;
   padding: 2px 5px 2px 10px;
}
@media only screen and (max-width: 1050px) {
   #app .main div.pure-u-5-6.bkg {
      width: 100%;
   }
}
html,
body {
   margin: 0;
   padding: 0;
   /* background-image: url(./assets/main-bkg.jpg); */
   background-color: black;
}
#app .main { 
   padding-top:0;
}
#app {
   font-family: "Avenir", Helvetica, Arial, sans-serif;
   -webkit-font-smoothing: antialiased;
   -moz-osx-font-smoothing: grayscale;
   background-color: black;
   border-right: 1px solid black;
}
#app h2 {
   margin: 0;
   position: relative;
}
#app a {
   color: #24890d;
   font-weight: 500;
   text-decoration: none;
}
#app a:hover {
   text-decoration: underline;
}
#app a.back {
   position: absolute;
   top: 15px;
   left: 15px;
}
div.bkg {
   background-image: url(./assets/books-bkg.jpg);
   position: relative;
}
#app div.content {
   position: relative;
   padding: 30px 60px 30px 60px;
   min-height: 600px;
   background: white;
   font-weight: 400;
   color: #444;
   border: 1px solid black;
   /* width: 80%; */
   margin: 25px;
}
#app .content h2 {
   font-size: 1.5em;
   font-weight: bold;
   border-bottom: 1px dashed #666;
   font-family: 'Special Elite', cursive;
   padding-bottom: 5px;
   margin-bottom: 15px;
}
#customdropzone .dz-preview {
   display: inline-block;
   background: transparent;
   margin: 5px;
}
#customdropzone .dz-preview .dz-details {
   background-color: rgba(64, 64, 64, 0.5);
   transition: none;
}
#customdropzone .dz-preview .dz-image {
   width: 150px;
   height: 150px;
   border-radius: 50%;
}
#app a.dz-remove {
   color: white;
}
</style>

