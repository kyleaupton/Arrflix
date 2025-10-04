import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import App from '@/App.vue'
import router from '@/router'
import SnagglePreset from '@/theme/preset'
import { VueQueryPlugin } from '@tanstack/vue-query'
import 'primeicons/primeicons.css'
import { useAuthStore } from '@/stores/auth'

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

// rehydrate auth before mount so guards/UI have token
const auth = useAuthStore()
await auth.rehydrate()

app.mount('#app')
