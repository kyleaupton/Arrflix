<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from 'primevue/button'

const route = useRoute()
const router = useRouter()

const currentTab = computed(() => {
  const path = route.path
  if (path.endsWith('/libraries')) return 'libraries'
  return 'general'
})

const navigateToTab = (tab: 'general' | 'libraries') => {
  router.push(`/settings/${tab}`)
}
</script>

<template>
  <div class="p-4 space-y-4">
    <h2 class="text-2xl font-semibold">Settings</h2>

    <!-- Navigation Tabs -->
    <div class="flex gap-2 border-b">
      <Button
        :label="'General'"
        :severity="currentTab === 'general' ? 'primary' : 'secondary'"
        :text="currentTab !== 'general'"
        @click="navigateToTab('general')"
      />
      <Button
        :label="'Libraries'"
        :severity="currentTab === 'libraries' ? 'primary' : 'secondary'"
        :text="currentTab !== 'libraries'"
        @click="navigateToTab('libraries')"
      />
    </div>

    <!-- Content Area -->
    <router-view />
  </div>
</template>
