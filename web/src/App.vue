<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useLayoutStore } from '@/stores/layout'
import AppHeader from '@/components/AppHeader.vue'
import AppSidebar from '@/components/sidebar/AppSidebar.vue'
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
    <AppSidebar v-if="route.meta.layout !== 'auth'" v-model:mobileVisible="mobileSidebarVisible" />

    <div class="app-body">
      <AppHeader
        v-if="route.meta.layout !== 'auth'"
        @toggle-sidebar="mobileSidebarVisible = true"
      />

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
  overflow-x: hidden; /* prevent stray horizontal scroll from flex gaps or sticky elements */
}

.app-header {
  position: sticky;
  top: 0;
  z-index: 10;
}

.app-main {
  flex: 1 1 auto;
  min-width: 0; /* allow children to shrink without forcing horizontal scroll */
}

.app-body {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  flex: 1 1 auto; /* take remaining width next to sidebar without overflowing */
  min-width: 0; /* critical for flex layouts to avoid overflow */
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

.auth .app-body {
  max-width: 100%;
  margin: 0 auto;
}
</style>
