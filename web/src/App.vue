<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useLayoutStore } from '@/stores/layout'
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
    <div v-if="route.meta.layout !== 'auth'" class="app-body">
      <AppSidebar v-model:mobileVisible="mobileSidebarVisible" />

      <main class="app-main">
        <RouterView />
      </main>
    </div>

    <main v-else class="app-main auth">
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

.app-body {
  flex: 1;
  display: flex;
  min-height: 0;
  gap: var(--layout-gap, 1rem);
  padding: var(--layout-padding, 1rem);
}

.app-main {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
  min-width: 0;
}

.app-main.auth {
  max-width: 100%;
  margin: 0 auto;
  padding: 0;
}

.brand {
  font-weight: 600;
}

@media (max-width: 1023px) {
  .app-body {
    padding: var(--layout-padding-mobile, 0.75rem);
    padding-top: 0;
  }
}
</style>
