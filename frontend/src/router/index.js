import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import About from '@/views/About.vue'
import Press from '@/views/Press.vue'
import Events from '@/views/Events.vue'
import News from '@/views/News.vue'
import Pedagogy from '@/views/Pedagogy.vue'
import FAQ from '@/views/FAQ.vue'
import Submit from '@/views/Submit.vue'
import SearchResults from '@/views/SearchResults.vue'
import Submission from '@/views/Submission.vue'
import Thanks from '@/views/Thanks.vue'
import Forbidden from '@/views/Forbidden.vue'
import AdminHome from '@/views/admin/AdminHome.vue'
import AdminEvents from '@/views/admin/AdminEvents.vue'
import AdminNews from '@/views/admin/AdminNews.vue'
import AdminPedagogy from '@/views/admin/AdminPedagogy.vue'
import AdminSubmission from '@/views/admin/AdminSubmission.vue'
import { useSystemStore } from "@/stores/system"
import { useAdminStore } from "@/stores/admin"
import VueCookies from 'vue-cookies'

const router = createRouter({
   history: createWebHistory(),
   routes: [
      {
         path: '/',
         name: 'home',
         component: Home
      },
      {
         path: '/about',
         name: 'about',
         component: About
      },
      {
         path: '/press',
         name: 'press',
         component: Press
      },
      {
         path: '/news',
         name: 'news',
         component: News
      },
      {
         path: '/pedagogy/:id?',
         name: 'pedagogy',
         component: Pedagogy
      },
      {
         path: '/events',
         name: 'events',
         component: Events
      },
      {
         path: '/faq',
         name: 'faq',
         component: FAQ
      },
      {
         path: '/submit',
         name: 'submit',
         component: Submit,
      },
      {
         path: '/submissions/:id',
         name: 'submission',
         component: Submission
      },
      {
         path: '/results',
         name: 'results',
         component: SearchResults
      },
      {
         path: '/thanks',
         name: 'thanks',
         component: Thanks
      },
      {
         path: '/admin',
         redirect: '/admin/submissions',
      },
      {
         path: '/admin/submissions',
         name: 'admin',
         component: AdminHome,
         meta: { requiresAuth: true }
      },
      {
         path: '/admin/submissions/:id',
         name: 'admin-submission',
         component: AdminSubmission,
         meta: { requiresAuth: true }
      },
      {
         path: '/admin/events',
         name: 'admin-events',
         component: AdminEvents,
         meta: { requiresAuth: true }
      },
      {
         path: '/admin/news',
         name: 'admin-news',
         component: AdminNews,
         meta: { requiresAuth: true }
      },
      {
         path: '/admin/pedagogy',
         name: 'admin-pedagogy',
         component: AdminPedagogy,
         meta: { requiresAuth: true }
      },
      {
         path: '/forbidden',
         name: 'forbidden',
         component: Forbidden
      }
   ],
   scrollBehavior(_to, _from, _savedPosition) {
      return new Promise(resolve => {
         setTimeout( () => {
            resolve({left: 0, top: 0})
         }, 100)
      })
   },
})

router.beforeEach(to => {
   const system = useSystemStore()
   if (to.meta.requiresAuth == true) {
      console.log("AUTH REQUIRED")
      const admin = useAdminStore()
      system.adminMode = true
      if (admin.isAuthenticated == false) {
         console.log("NOT AUTHENTICATED")
         let authUser = VueCookies.get("bt_admin_user")
         if (authUser) {
            console.log("HAS COOKED")
            authUser.authenticated = true
            admin.user = authUser
            VueCookies.remove("bt_admin_user")
         } else {
            console.log("REDIRECT TO AUTH")
            return "/authenticate?url=" + to.fullPath
         }
      } else {
         system.adminMode = false
      }
   }
})

export default router