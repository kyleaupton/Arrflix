<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import AppSidebar from '@/components/AppSidebar.vue'
import AppLayoutHeader from '@/components/AppLayoutHeader.vue'
import DialogContainer from '@/components/DialogContainer.vue'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
import { TooltipProvider } from '@/components/ui/tooltip'

import 'vue-sonner/style.css'
import { Toaster } from '@/components/ui/sonner'

const authStore = useAuthStore()
const appStore = useAppStore()
const router = useRouter()
const route = useRoute()

// Bootstrap already ran in main.ts — just handle setup redirect
if (appStore.needsSetup && route.path !== '/setup') {
  router.push('/setup')
}
</script>

<template>
  <TooltipProvider>
    <Toaster position="top-center" />
    <DialogContainer />
    <div
      v-if="!appStore.isReady"
      class="flex min-h-svh items-center justify-center"
    >
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
