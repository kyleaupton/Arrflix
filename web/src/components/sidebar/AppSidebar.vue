<script setup lang="ts">
import Drawer from 'primevue/drawer'
import AppSidebarMenu from './AppSidebarMenu.vue'

const mobileVisible = defineModel<boolean>('mobileVisible', { default: false })
</script>

<template>
  <aside class="sidebar-desktop">
    <AppSidebarMenu />
  </aside>

  <Drawer v-model:visible="mobileVisible" position="left" class="sidebar-mobile">
    <AppSidebarMenu />
  </Drawer>
</template>

<style scoped>
.sidebar-desktop {
  display: none;
}

/* Desktop and up */
@media (min-width: 1024px) {
  .sidebar-desktop {
    display: block;
    width: 280px;
    flex: 0 0 280px;
    padding: 0.75rem;
  }
}

/* Ensure Sidebar panel uses full height on mobile */
.sidebar-mobile :deep(.p-sidebar-content) {
  padding: 0.75rem;
}
</style>
