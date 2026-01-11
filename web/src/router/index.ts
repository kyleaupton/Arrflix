import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { client } from '@/client/client.gen'

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
      path: '/search',
      component: () => import('@/views/Search.vue'),
    },
    {
      path: '/downloads',
      component: () => import('@/views/Downloads.vue'),
    },
    {
      path: '/users',
      component: () => import('@/views/Users.vue'),
    },
    {
      path: '/settings',
      component: () => import('@/views/settings/SettingsLayout.vue'),
      children: [
        {
          path: '',
          redirect: '/settings/general',
        },
        {
          path: 'general',
          component: () => import('@/views/settings/GeneralSettings.vue'),
        },
        {
          path: 'libraries',
          component: () => import('@/views/settings/LibrarySettings.vue'),
        },
        {
          path: 'indexers',
          component: () => import('@/views/settings/IndexersSettings.vue'),
        },
        {
          path: 'name-templates',
          component: () => import('@/views/settings/NameTemplateSettings.vue'),
        },
        {
          path: 'downloaders',
          component: () => import('@/views/settings/downloader/DownloaderSettings.vue'),
        },
        {
          path: 'policies',
          component: () => import('@/views/settings/PolicySettings.vue'),
        },
      ],
    },
    {
      path: '/login',
      component: () => import('@/views/Login.vue'),
      meta: { public: true, layout: 'auth' },
    },
    {
      path: '/setup',
      component: () => import('@/views/Setup.vue'),
      meta: { public: true, layout: 'auth', setup: true },
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
    {
      path: '/person/:id',
      component: () => import('@/views/Person.vue'),
    },
  ],
})

// Check setup status
async function checkSetupStatus(): Promise<boolean> {
  try {
    const response = await client.get<{ initialized: boolean }>({
      url: '/v1/setup/status',
    })
    return (response as any).data.initialized
  } catch {
    // If setup check fails, assume initialized (safer default)
    return true
  }
}

router.beforeEach(async (to) => {
  const auth = useAuthStore()

  // Rehydrate auth token if needed
  if (!auth.token) {
    await auth.rehydrate()
  }

  // Check if setup is complete
  const isInitialized = await checkSetupStatus()

  // If not initialized, redirect everything to setup (except setup page itself)
  if (!isInitialized && to.path !== '/setup') {
    return { path: '/setup' }
  }

  // If initialized, block setup page
  if (isInitialized && to.path === '/setup') {
    return { path: '/login' }
  }

  // Public routes allowed
  if (to.meta.public) {
    return true
  }

  // Require auth for protected routes
  if (!auth.isAuthenticated) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }

  return true
})

export default router
