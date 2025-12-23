import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin } from '@tanstack/vue-query'

import App from '@/App.vue'
import router from '@/router'
import { useAuthStore } from '@/stores/auth'
import { client } from '@/client/client.gen'
import '@/main.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(VueQueryPlugin, {
  enableDevtoolsV6Plugin: import.meta.env.DEV,
})

// Redirect to login on any 401 response
client.interceptors.response.use(async (response) => {
  if (response.status === 401) {
    const auth = useAuthStore()
    auth.logout()
    const current = router.currentRoute.value
    if (current.path !== '/login') {
      router.replace({ path: '/login', query: { redirect: current.fullPath } })
    }
  }
  return response
})

// Force dark mode globally
document.documentElement.classList.add('dark')

// Rehydrate auth before mount so guards/UI have token
const auth = useAuthStore()
await auth.rehydrate()

app.mount('#app')
