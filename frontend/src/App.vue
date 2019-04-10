<template>
   <div id="app">
      <BookTracesHeader/>
      <router-view/>
   </div>
</template>


<script>
import BookTracesHeader from "@/components/BookTracesHeader";

export default {
   components: {
      BookTracesHeader
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
   background: #dedede;
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
   color: cornflowerblue;
   font-weight: 500;
   text-decoration: none;
}
#app a:hover {
   text-decoration: underline;
}
</style>

