<template>
   <div id="app">
      <BookTracesHeader/>
      <div class="main pure-g">
         <div class="pure-u-1-6">
            <BookTracesSidebar/>
         </div>
         <div class="pure-u-5-6 bkg">
            <router-view/>
         </div>
      </div>
      <BookTracesFooter/>
   </div>
</template>


<script>
import BookTracesSidebar from "@/components/BookTracesSidebar";
import BookTracesHeader from "@/components/BookTracesHeader";
import BookTracesFooter from "@/components/BookTracesFooter";

export default {
   components: {
      BookTracesHeader,
      BookTracesSidebar,
      BookTracesFooter
   },
   data: function() {
      return {};
   },
   created: function() {}
};
</script>

<style>
html,
body {
   height: 100%;
   margin: 0;
   padding: 0;
   background-color: #0e2b56;
}
#app {
   font-family: "Avenir", Helvetica, Arial, sans-serif;
   -webkit-font-smoothing: antialiased;
   -moz-osx-font-smoothing: grayscale;
   min-height: 100%;
   position: relative;
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
div.bkg {
   background-image: url(./assets/main-bkg.jpg);
}
#app div.content {
   position: relative;
   min-width: 800px;
   padding: 30px 60px 250px 60px;
   min-height: 600px;
   background: white;
   font-weight: 400;
   color: #444;
   width: 80%;
   border-right: 1px solid #666;
}
#app .content h2 {
   font-size: 1.5em;
   font-weight: 500;
   color: #eb5f0c;
   border-bottom: 1px dashed #eb5f0c;
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

