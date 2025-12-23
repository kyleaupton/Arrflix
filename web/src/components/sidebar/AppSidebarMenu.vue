<template>
  <Menu class="app-sidebar-menu h-full flex flex-col" :model="navigationItems" :dt>
    <template #start>
      <span class="inline-flex items-center gap-1 px-2 py-2">
        <svg
          width="35"
          height="40"
          viewBox="0 0 35 40"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
          class="h-8"
        >
          <path d="..." fill="var(--p-primary-color)" />
          <path d="..." fill="var(--p-text-color)" />
        </svg>
        <span class="text-xl font-semibold">Snaggle</span>
      </span>
    </template>

    <template #end>
      <span class="inline-flex items-center gap-1 px-2 py-2">
        <Avatar :label="user?.name?.[0]" class="mr-2" shape="circle" />
        <span class="inline-flex flex-col items-start">
          <span class="font-bold">{{ user?.name }}</span>
          <span class="text-sm">{{ user?.email }}</span>
        </span>
      </span>
    </template>
  </Menu>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import Menu from 'primevue/menu'
import Avatar from 'primevue/avatar'
import type { MenuDesignTokens, MenuTokenSections } from '@primeuix/themes/types/menu'
import { navigationItems } from '@/config/navigation'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const { user } = storeToRefs(authStore)

const root: MenuTokenSections.Root = {}

const list: MenuTokenSections.List = {
  gap: '0.5rem',
}

const item: MenuTokenSections.Item = {}

const dt = { root, list, item } satisfies MenuDesignTokens
</script>

<style>
.app-sidebar-menu .p-menu-list {
  flex-grow: 1;
}
</style>
