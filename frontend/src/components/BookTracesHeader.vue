<template>
   <div class="bt-header">
      <div class="site-name">
         <router-link to="/">
            <span class="site-name">Book Traces</span>
         </router-link>
      </div>
      <div v-if="!system.adminMode" class="pure-menu pure-menu-horizontal menubar">
         <ul class="pure-menu-list">
            <li @click="adminClicked" v-bind:class="{active: system.adminMode}" class="pure-menu-item admin">Admin</li>
            <li class="pure-menu-item"><router-link to="/" exact>Home</router-link></li>
            <li class="pure-menu-item"><router-link to="/about">About</router-link></li>
            <li class="pure-menu-item"><router-link to="/press">Press</router-link></li>
            <li class="pure-menu-item"><router-link to="/news">News</router-link></li>
            <li class="pure-menu-item"><router-link to="/pedagogy">Pedagogy</router-link></li>
            <li class="pure-menu-item"><router-link to="/events">Events</router-link></li>
            <li class="pure-menu-item"><router-link to="/faq">FAQ</router-link></li>
            <li class="pure-menu-item"><router-link to="/submit">Submit a Book</router-link></li>
         </ul>
         <span @click="showSearchClick" class="search">
            <i v-bind:class="{selected: unAuth.showSearch}" class="fas fa-search"></i>
         </span>
      </div>
      <div v-else class="pure-menu pure-menu-horizontal menubar">
         <ul class="pure-menu-list">
            <li class="public pure-menu-item"><router-link to="/">Public</router-link></li>
            <li class="pure-menu-item"><router-link to="/admin/submissions">Submissions</router-link></li>
            <li class="pure-menu-item"><router-link to="/admin/events">Events</router-link></li>
            <li class="pure-menu-item"><router-link to="/admin/news">News</router-link></li>
            <li class="pure-menu-item"><router-link to="/admin/pedagogy">Pedagogy</router-link></li>
         </ul>
      </div>
      <div  v-if="!system.adminMode" class="hmenu">
         <div @click="toggleHMenu" class="hmenu-button"><i class="fas fa-bars"></i></div>
         <ul id="hmenu" class="hmenu-items hidden">
            <li @click="adminClicked" v-bind:class="{active: system.adminMode}" class="admin">Admin</li>
            <li @click="toggleHMenu"><router-link to="/" exact>Home</router-link></li>
            <li @click="toggleHMenu"><router-link to="/about">About</router-link></li>
            <li @click="toggleHMenu"><router-link to="/press">Press</router-link></li>
            <li @click="toggleHMenu"><router-link to="/news">News</router-link></li>
            <li @click="toggleHMenu"><router-link to="/pedagogy">Pedagogy</router-link></li>
            <li @click="toggleHMenu"><router-link to="/events">Events</router-link></li>
            <li @click="toggleHMenu"><router-link to="/faq">FAQ</router-link></li>
            <li @click="toggleHMenu"><router-link to="/submit">Submit a Book</router-link></li>
            <li><span @click="showSearchClick" class="search small">Search</span></li>
         </ul>
      </div>
   </div>
</template>

<script setup>
import { useSystemStore } from "@/stores/system"
import { useUnauthStore } from "@/stores/unauth"

const system = useSystemStore()
const unAuth = useUnauthStore()

const adminClicked = (() => {
   window.location.href = "/authenticate?url=/admin/submissions"
   hideHMenu()
})

const showSearchClick = (() => {
   unAuth.showSearch = !unAuth.showSearch
   hideHMenu()
})

const hideHMenu = (() => {
   let items = document.getElementById("hmenu")
   if (!items.classList.contains("hidden")) {
      items.classList.add("hidden")
   }
})

const toggleHMenu = (() => {
   let items = document.getElementById("hmenu")
   if (items.classList.contains("hidden")) {
      items.classList.remove("hidden")
   } else {
      items.classList.add("hidden")
   }
})
</script>

<style scoped>
span.search.small {
   font-size: 1.05em;
}
.hmenu-button {
   position: absolute;
   right: 16px;
   bottom: 12px;
   font-size: 1.5em;
   cursor: pointer;
}
.hmenu-items.hidden {
   height: none;
}
.hmenu-items {
   position: absolute;
   right: 15px;
   z-index: 1000;
   background: black;
   width: 100px;
   list-style: none;
   padding: 10px 10px;
   margin: 0;
   text-align: right;
}
.hmenu li {
   padding: 4px 0;
}
.hmenu li a:hover {
   color: #55d737 !important;
}
.hmenu li .router-link-active {
   color: #55d737 !important;
}
@media only screen and (max-width: 768px) {
   div.pure-menu.pure-menu-horizontal.menubar {
      display:none;
   }
}
@media only screen and (min-width: 768px) {
   div.hmenu {
      display:none;
   }
}
#app .menubar .pure-menu-item .router-link-active {
   color: #55d737;
   border-bottom: 2px solid green;
}
#app .menubar .pure-menu-item.public .router-link-active {
   color:white;
   border-bottom: 0;
}
#app .menubar .pure-menu-item.public .router-link-active:hover {
   color:cornflowerblue !important;
   border-bottom: 0;
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
   opacity: 1;
}
.admin {
   opacity: 0.5;
   font-size: 0.85em;
}
.admin.active {
   color: #55d737;
   opacity: 1;
}
div.bt-header {
   background-color: black;
   color: white;
   padding:15px;
   position: relative;
   font-family: 'Special Elite', serif;
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
