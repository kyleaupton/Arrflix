<template>
  <PanelMenu :model="navigationItems" :dt="sidebarDt" />
</template>

<script setup lang="ts">
import PanelMenu from 'primevue/panelmenu'
import { navigationItems } from '@/config/navigation'
import sidebarDt from './sidebar-dt'
</script>

<style scoped></style>
