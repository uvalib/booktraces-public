<template>
   <div class="header">
      <router-link :to="{ name: 'home'}">
         <span class="site-name">Book Traces</span>
      </router-link>
      <a target="_blank" href="https://library.virginia.edu"><img class="uva-library" src="../assets/uvalogo.png"/></a>
   </div>
</template>

<script>
export default {
}
</script>

<style scoped>
div.header {
   background-color: black;
   color: white;
   padding:15px;
   text-align: left;
}
#app div.header a {
  color: white;
  text-decoration: none;
}
div.header span.site-name {
   margin: 0;
   font-size: 1.4em;
   position: relative;
   font-weight: bold;
}
img.uva-library {
  position: absolute;
  top: 10px;
  right: 10px;
  max-height: 40px;
}
</style>
