<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

onMounted(async () => {
  const ok = await auth.completeSsoFromCallback(
    new URLSearchParams(route.fullPath.split('?')[1] ?? ''),
  )
  if (ok) {
    router.replace({ path: '/' })
  } else {
    router.replace({ path: '/login' })
  }
})
</script>

<template>
  <div class="wrap">Signing you in…</div>
</template>

<style scoped>
.wrap {
  display: grid;
  place-items: center;
  min-height: 60vh;
}
</style>
