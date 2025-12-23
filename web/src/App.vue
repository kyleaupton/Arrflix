<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useLayoutStore } from '@/stores/layout'
import AppHeader from '@/components/AppHeader.vue'
import AppSidebar from '@/components/sidebar/AppSidebar.vue'
import DynamicDialog from 'primevue/dynamicdialog'
import '@/main.css'

const mobileSidebarVisible = ref(false)
const route = useRoute()

onMounted(async () => {
  const layoutStore = useLayoutStore()

  // watch window resize
  window.addEventListener('resize', () => {
    layoutStore.screenWidth = window.innerWidth
  })
})
</script>

<template>
  <!-- <Toast /> -->

  <div class="app-shell" :class="{ auth: route.meta.layout === 'auth' }">
    <AppHeader
      v-if="route.meta.layout !== 'auth'"
      v-model:mobileSidebarVisible="mobileSidebarVisible"
    />

    <AppSidebar v-if="route.meta.layout !== 'auth'" v-model:mobileVisible="mobileSidebarVisible" />

    <main class="app-main">
      <RouterView />
    </main>
  </div>

  <DynamicDialog />
</template>

<style scoped>
.app-shell {
  height: 100dvh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.app-main {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
  padding: 0.75rem;
}

@media (min-width: 1024px) {
  .app-main {
    padding: 1rem 1.25rem 1.25rem 1rem;
  }
}

.brand {
  font-weight: 600;
}

.auth .app-main {
  max-width: 100%;
  margin: 0 auto;
}
</style>
