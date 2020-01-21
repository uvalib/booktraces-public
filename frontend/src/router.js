import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import About from './views/About.vue'
import Press from './views/Press.vue'
import Events from './views/Events.vue'
import News from './views/News.vue'
import FAQ from './views/FAQ.vue'
import Submit from './views/Submit.vue'
import SearchResults from './views/SearchResults.vue'
import Submission from './views/Submission.vue'
import Thanks from './views/Thanks.vue'
import Forbidden from './views/Forbidden.vue'
import AdminHome from './views/admin/AdminHome.vue'
import AdminEvents from './views/admin/AdminEvents.vue'
import AdminNews from './views/admin/AdminNews.vue'
import AdminSubmission from './views/admin/AdminSubmission.vue'
import store from './store'

Vue.use(Router)

const router = new Router({
  mode: 'history',
  base: process.env.BASE_URL,
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
      beforeEnter: (_to, _from, next) => {
        store.dispatch('public/getUploadID')
        next()
      }
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
      path: '/forbidden',
      name: 'forbidden',
      component: Forbidden
    }
  ],
  scrollBehavior(/*to, from, savedPosition*/) {
    return { x: 0, y: 0 }
  },
})

router.beforeEach((to, _from, next) => {
  if (to.meta.requiresAuth == true) {
    store.commit("setAdminMode", true)
    let getters = store.getters
    if (getters["admin/isAuthenticated"] == false) {
      let authUser = Vue.$cookies.get("bt_admin_user")
      if (authUser) {
        authUser.authenticated = true
        store.commit("admin/setUser", authUser)
        Vue.$cookies.remove("bt_admin_user")
      } else {
        let authURL =  "/authenticate?url=" + to.fullPath
        window.location.href = authURL
      }
    }
  } else {
    store.commit("setAdminMode", false)
    store.dispatch("public/getArchiveDates")
    store.dispatch("public/getRecentSubmissions")
  }
  next()
})

export default router