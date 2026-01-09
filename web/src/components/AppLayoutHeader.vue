<script setup lang="ts">
import { Search } from 'lucide-vue-next'
import { SidebarTrigger } from '@/components/ui/sidebar'
import Input from '@/components/ui/input/Input.vue'
</script>

<template>
  <header
    class="z-10 fixed w-full flex h-16 shrink-0 items-center gap-2 border-b bg-background transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12"
  >
    <div class="flex items-center gap-2 px-4">
      <SidebarTrigger class="-ml-1" />
      <div class="relative">
        <Input placeholder="Search..." class="pl-8 w-64" />
        <Search
          class="pointer-events-none absolute top-1/2 left-2 size-4 -translate-y-1/2 opacity-50 select-none"
        />
      </div>
    </div>
  </header>
</template>
