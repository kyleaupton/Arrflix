<script setup lang="ts">
import Drawer from 'primevue/drawer'
import Button from 'primevue/button'
import AppSidebarMenu from './AppSidebarMenu.vue'

const mobileVisible = defineModel<boolean>('mobileVisible', { default: false })
</script>

<template>
  <!-- Mobile Menu Button -->
  <Button
    v-if="!mobileVisible"
    class="mobile-menu-button"
    icon="pi pi-bars"
    variant="text"
    @click="mobileVisible = true"
  />

  <!-- Desktop: Persistent Sidebar -->
  <aside class="app-sidebar-desktop">
    <div class="sidebar-wrapper">
      <AppSidebarMenu />
    </div>
  </aside>

  <!-- Mobile: Drawer Sidebar -->
  <Drawer v-model:visible="mobileVisible" position="left" class="sidebar-mobile">
    <AppSidebarMenu @navigate="mobileVisible = false" />
  </Drawer>
</template>

<style scoped>
.app-sidebar-desktop {
  display: none;
  width: var(--sidebar-width, 280px);
  flex-shrink: 0;
}

@media (min-width: 1024px) {
  .app-sidebar-desktop {
    display: block;
  }
}

.sidebar-wrapper {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--p-card-background);
  border-radius: var(--p-border-radius-lg);
  box-shadow: var(--p-shadow-1);
  overflow: hidden;
}

/* Mobile Menu Button */
.mobile-menu-button {
  position: fixed;
  top: 1rem;
  left: 1rem;
  z-index: 1000;
  display: block;
}

@media (min-width: 1024px) {
  .mobile-menu-button {
    display: none;
  }
}

/* Mobile Drawer */
.sidebar-mobile :deep(.p-drawer-content) {
  padding: 0;
  background: var(--p-card-background);
}

.sidebar-mobile :deep(.p-drawer) {
  background: var(--p-card-background);
}
</style>
