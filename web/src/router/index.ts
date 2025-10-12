import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: () => import('@/views/Home.vue'),
    },
    {
      path: '/library',
      component: () => import('@/views/Library.vue'),
    },
    {
      path: '/requests',
      component: () => import('@/views/Requests.vue'),
    },
    {
      path: '/users',
      component: () => import('@/views/Users.vue'),
    },
    {
      path: '/settings',
      component: () => import('@/views/Settings.vue'),
    },
    {
      path: '/login',
      component: () => import('@/views/Login.vue'),
      meta: { public: true, layout: 'auth' },
    },
    {
      path: '/auth/callback',
      component: () => import('@/views/AuthCallback.vue'),
      meta: { public: true, layout: 'auth' },
    },

    // Media
    {
      path: '/movie/:id',
      component: () => import('@/views/Movie.vue'),
    },
    {
      path: '/series/:id',
      component: () => import('@/views/Series.vue'),
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
