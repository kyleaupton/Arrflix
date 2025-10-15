<script setup lang="ts">
import Drawer from 'primevue/drawer'
import AppSidebarMenu from './AppSidebarMenu.vue'

const mobileVisible = defineModel<boolean>('mobileVisible', { default: false })
</script>

<template>
  <Drawer v-model:visible="mobileVisible" position="left" class="sidebar-mobile">
    <AppSidebarMenu />
  </Drawer>
</template>

<style scoped>
/* Ensure Sidebar panel uses full height on mobile */
.sidebar-mobile :deep(.p-sidebar-content) {
  padding: 0.75rem;
}
</style>
