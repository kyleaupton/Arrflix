<script setup lang="ts">
import { SidebarTrigger } from '@/components/ui/sidebar'
import AppBreadcrumbs from '@/components/AppBreadcrumbs.vue'
</script>

<template>
  <header
    class="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12"
  >
    <div class="flex items-center gap-2 px-4">
      <SidebarTrigger class="-ml-1" />
      <AppBreadcrumbs />
    </div>
  </header>
</template>

