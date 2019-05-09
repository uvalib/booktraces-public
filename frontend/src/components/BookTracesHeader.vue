<template>
   <div class="bt-header">
      <div class="site-name">
         <router-link to="/">
            <span class="site-name">Book Traces</span>
         </router-link>
      </div>
      <div v-if="!adminMode" class="pure-menu pure-menu-horizontal menubar">
         <ul class="pure-menu-list">
            <li @click="adminClicked" v-bind:class="{active: adminMode}" class="pure-menu-item admin">Admin</li>
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
   </div>
</template>

<script>
import { mapState } from 'vuex'
export default {
   computed: {
      ...mapState({
         showSearch: state => state.public.showSearch,
         adminMode: state => state.adminMode
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
.admin.active {
   color: #55d737;
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
.pure-menu-horizontal.menubar {
   position: absolute;
   bottom: 7px;
   right: 15px;
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
</style>
