<template>
   <Menubar :model="btMenu">
      <template #start>
         <div class="site-name">
            <router-link to="/" class="site-name">Book Traces</router-link>
         </div>
      </template>
      <template #end v-if="!system.adminMode">
         <Button icon="pi pi-search" @click="showSearchClick" text rounded class="search-icon" severity="secondary"/>
      </template>
   </Menubar>
</template>

<script setup>
import { useSystemStore } from "@/stores/system"
import { useSubmissionsStore } from "@/stores/submissions"
import { computed } from 'vue'
import Menubar from 'primevue/menubar'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()
const system = useSystemStore()
const submissionsStore = useSubmissionsStore()

const btMenu = computed( () => {
   if (system.adminMode) {
      return [
         {label: "Public", class: "public", command: ()=>menuLinkClicked("/")},
         {label: "Submissions", class: `${linkClass("/admin/submissions")}`, command: ()=>menuLinkClicked("/admin/submissions")},
         {label: "Events", class: `${linkClass("/admin/events")}`, command: ()=>menuLinkClicked("/admin/events")},
         {label: "News", class: `${linkClass("/admin/news")}`, command: ()=>menuLinkClicked("/admin/news")},
         {label: "Pedagogy", class: `${linkClass("/admin/pedagogy")}`, command: ()=>menuLinkClicked("/admin/pedagogy")},
      ]
   } else {
      return [
         {label: "Admin", class: "admin", command: ()=>adminClicked()},
         {label: "Home", class: `${linkClass("/")}`, command: ()=>menuLinkClicked("/")},
         {label: "About", class: `${linkClass("/about")}`, command: ()=>menuLinkClicked("/about")},
         {label: "Press", class: `${linkClass("/press")}`, command: ()=>menuLinkClicked("/press")},
         {label: "News", class: `${linkClass("/news")}`, command: ()=>menuLinkClicked("/news")},
         {label: "Pedagogy", class: `${linkClass("/pedagogy")}`, command: ()=>menuLinkClicked("/pedagogy")},
         {label: "Events", class: `${linkClass("/events")}`, command: ()=>menuLinkClicked("/events")},
         {label: "FAQ", class: `${linkClass("/faq")}`, command: ()=>menuLinkClicked("/faq")},
         {label: "Submit a Book", class: `${linkClass("/submit")}`, command: ()=>menuLinkClicked("/submit")},
      ]
   }
})

const linkClass = ( (tgtRoute) => {
   if (tgtRoute == route.path ) return "active"
   return ""
})

const menuLinkClicked = ( (tgtRoute) => {
   router.push(tgtRoute)
})

const adminClicked = (() => {
   window.location.href = "/authenticate?url=/admin/submissions"
})

const showSearchClick = (() => {
   submissionsStore.showSearch = !submissionsStore.showSearch
})

</script>

<style scoped>
:deep(.admin .p-menubar-item-link) {
   color: #aaa;
   background-color: transparent;
}
:deep(.active .p-menubar-item-link) {
   color: #55d737;
}
a.site-name {
   font-family: Special Elite, cursive;
   font-size: 24px;
   color: white;
}
button.search-icon {
   color: white;
}
</style>
