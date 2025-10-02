<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useLayoutStore } from './stores/layout'

import AppHeader from './components/AppHeader.vue'
import AppSidebar from './components/sidebar/AppSidebar.vue'
import './main.css'

const mobileSidebarVisible = ref(false)

onMounted(() => {
  const layoutStore = useLayoutStore()

  // watch window resize
  window.addEventListener('resize', () => {
    layoutStore.screenWidth = window.innerWidth
  })
})
</script>

<template>
  <!-- <Toast /> -->

  <div class="app-shell">
    <AppSidebar v-model:mobileVisible="mobileSidebarVisible" />

    <div class="app-body">
      <AppHeader @toggle-sidebar="mobileSidebarVisible = true" />

      <main class="app-main">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<style scoped>
.app-shell {
  min-height: 100vh;
  display: flex;
}

.app-header {
  position: sticky;
  top: 0;
  z-index: 10;
}

.app-main {
  flex: 1 1 auto;
}

.app-body {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  width: 100%;
  padding: 0.75rem;
}

@media (min-width: 1024px) {
  .app-body {
    padding: 1rem 1.25rem 1.25rem 1rem;
    gap: 1rem;
  }
}

.brand {
  font-weight: 600;
}
</style>
