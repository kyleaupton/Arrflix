<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import AppSidebar from '@/components/AppSidebar.vue'
import AppLayoutHeader from '@/components/AppLayoutHeader.vue'
import Login from '@/views/Login.vue'
import DialogContainer from '@/components/DialogContainer.vue'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'

const authStore = useAuthStore()
const isCheckingAuth = ref(true)

onMounted(async () => {
  if (!authStore.token) {
    await authStore.rehydrate()
  }
  isCheckingAuth.value = false
})
</script>

<template>
  <DialogContainer />
  <div v-if="isCheckingAuth" class="flex min-h-svh items-center justify-center">
    <div class="text-muted-foreground">Loading...</div>
  </div>
  <Login v-else-if="!authStore.isAuthenticated" />
  <SidebarProvider v-else>
    <AppSidebar />
    <SidebarInset>
      <AppLayoutHeader />
      <div class="flex flex-1 flex-col gap-4 p-4 pt-0 overflow-y-auto min-w-0">
        <router-view />
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
