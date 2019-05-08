<template>
   <div class="bt-header">
      <div class="site-name">
         <router-link to="/">
            <span class="site-name">Book Traces</span>
         </router-link>
      </div>
      <div class="pure-menu pure-menu-horizontal menubar">
         <ul class="pure-menu-list">
            <li @click="adminClicked" class="pure-menu-item admin">Admin</li>
            <li class="pure-menu-item"><router-link to="/about">About</router-link></li>
            <li class="pure-menu-item"><router-link to="/press">Press</router-link></li>
            <li class="pure-menu-item"><router-link to="/events">Events</router-link></li>
            <li class="pure-menu-item"><router-link to="/faq">FAQ</router-link></li>
            <li class="pure-menu-item"><router-link to="/submit">Submit a Book</router-link></li>
         </ul>
         <span @click="showSearchClick" class="search">
            <i v-bind:class="{selected: showSearch}" class="fas fa-search"></i>
         </span>
      </div>
      <div class="uva-icon">
         <a target="_blank" href="https://library.virginia.edu"><img class="uva-library" src="../assets/uvalogo.png"/></a>
      </div>
   </div>
</template>

<script>
import { mapState } from 'vuex'
export default {
   computed: {
      ...mapState({
         showSearch: state => state.public.showSearch
      }),
   },
   methods: {
      adminClicked() {
         window.location.href = "/authenticate?url=/admin"
      },
      showSearchClick() {
         this.$store.commit('public/showSearch', !this.showSearch )
      },
   }
}
</script>

<style scoped>
#app .menubar .pure-menu-item .router-link-active {
   color: #55d737;
   border-bottom: 2px solid green;
}
.search {
   font-size: 1.25em;
}
.fas.selected {
   color: #55d737;  
}
.search:hover {
   color: #55d737;  
   cursor: pointer;
}
.admin:hover {
   color:palevioletred;
   cursor: pointer;
}
div.bt-header {
   background-color: black;
   color: white;
   padding:15px;
   position: relative;
   font-family: 'Special Elite', cursive;
}
#app div.bt-header a {
  color: white;
  text-decoration: none;
}
div.bt-header span.site-name {
   margin: 0;
   font-size: 24px;
   position: relative;
}
img.uva-library {
  position: absolute;
  top: 6px;
  right: 6px;
  max-height: 40px;
}
.pure-menu-horizontal.menubar {
   position: absolute;
   bottom: 7px;
   right: 110px;
   white-space: inherit;
   width: auto;
   display: inline-block;
}
li.pure-menu-item {
   padding: 0 10px;
}
.pure-menu-item a:hover {
   border-bottom: 2px solid white;
}
.bt-header {
   position: fixed !important;
   width: 100%;
   z-index: 1000;
   max-width: 1230px;
}
</style>
