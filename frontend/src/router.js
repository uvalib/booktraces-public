import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import About from './views/About.vue'
import Press from './views/Press.vue'
import Events from './views/Events.vue'
import Contact from './views/Contact.vue'
import FAQ from './views/FAQ.vue'
import Submit from './views/Submit.vue'
import Admin from './views/Admin.vue'
import Forbidden from './views/Forbidden.vue'

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
      path: '/contact',
      name: 'contact',
      component: Contact
    },
    {
      path: '/faq',
      name: 'faq',
      component: FAQ
    },
    {
      path: '/submit',
      name: 'submit',
      component: Submit
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
