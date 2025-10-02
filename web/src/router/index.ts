import { createRouter, createWebHistory } from 'vue-router'
import Discover from '../views/Discover.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: Discover,
    },
  ],
})

export default router
