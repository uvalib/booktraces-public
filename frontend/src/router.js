import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import About from './views/About.vue'
import Press from './views/Press.vue'
import Events from './views/Events.vue'
import FAQ from './views/FAQ.vue'
import Submit from './views/Submit.vue'
import Submission from './views/Submission.vue'
import Thanks from './views/Thanks.vue'
import Admin from './views/Admin.vue'
import Forbidden from './views/Forbidden.vue'
import store from './store'

Vue.use(Router)

export default new Router({
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
      path: '/submission/:id',
      name: 'submission',
      component: Submission
    },
    {
      path: '/thanks',
      name: 'thanks',
      component: Thanks
    },
    {
      path: '/admin',
      name: 'admin',
      component: Admin
    },
    {
      path: '/forbidden',
      name: 'forbidden',
      component: Forbidden
    }
  ]
})
