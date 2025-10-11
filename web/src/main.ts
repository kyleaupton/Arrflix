import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import App from '@/App.vue'
import router from '@/router'
import SnagglePreset from '@/theme/preset'
import { VueQueryPlugin } from '@tanstack/vue-query'
import 'primeicons/primeicons.css'
import { useAuthStore } from '@/stores/auth'
import { client } from '@/client/client.gen'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(PrimeVue, {
  theme: {
    preset: SnagglePreset,
    options: {
      darkModeSelector: '.dark',
    },
  },
})
app.use(VueQueryPlugin)

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

// rehydrate auth before mount so guards/UI have token
const auth = useAuthStore()
await auth.rehydrate()

app.mount('#app')
