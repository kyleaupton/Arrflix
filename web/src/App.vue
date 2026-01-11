<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import AppSidebar from '@/components/AppSidebar.vue'
import AppLayoutHeader from '@/components/AppLayoutHeader.vue'
import DialogContainer from '@/components/DialogContainer.vue'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
import { TooltipProvider } from '@/components/ui/tooltip'
import { client } from '@/client/client.gen'

import 'vue-sonner/style.css'
import { Toaster } from '@/components/ui/sonner'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const isCheckingAuth = ref(true)
const isCheckingSetup = ref(true)
const needsSetup = ref(false)

onMounted(async () => {
  // Check setup status once on page load
  try {
    const response = await client.get({ url: '/v1/setup/status' })
    needsSetup.value = !(response as any).data.initialized
    if (needsSetup.value && route.path !== '/setup') {
      router.push('/setup')
    }
  } catch {
    // If setup check fails, assume initialized
    needsSetup.value = false
  }
  isCheckingSetup.value = false

  // Rehydrate auth token from localStorage
  if (!authStore.token) {
    await authStore.rehydrate()
  }
  isCheckingAuth.value = false
})
</script>

<template>
  <TooltipProvider>
    <Toaster position="top-center" />
    <DialogContainer />
    <div v-if="isCheckingAuth || isCheckingSetup" class="flex min-h-svh items-center justify-center">
      <div class="text-muted-foreground">Loading...</div>
    </div>
    <router-view v-else-if="route.meta.public" />
    <SidebarProvider v-else-if="authStore.isAuthenticated">
      <AppSidebar />
      <SidebarInset>
        <AppLayoutHeader />
        <div class="flex flex-1 flex-col gap-4 p-4 pt-19 overflow-y-auto min-w-0">
          <router-view />
        </div>
      </SidebarInset>
    </SidebarProvider>
    <router-view v-else />
  </TooltipProvider>
</template>
