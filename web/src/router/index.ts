import { createRouter, createWebHistory } from 'vue-router'
import Discover from '../views/Discover.vue'
import Login from '../views/Login.vue'
import AuthCallback from '../views/AuthCallback.vue'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: Discover,
    },
    {
      path: '/login',
      component: Login,
      meta: { public: true, layout: 'auth' },
    },
    {
      path: '/auth/callback',
      component: AuthCallback,
      meta: { public: true, layout: 'auth' },
    },
  ],
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (!auth.token) {
    await auth.rehydrate()
  }
  // Public routes
  if (to.meta.public) {
    return true
  }
  // Require auth for everything else (for now)
  if (!auth.isAuthenticated) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  return true
})

export default router
