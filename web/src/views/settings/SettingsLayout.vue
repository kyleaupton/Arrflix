<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from 'primevue/button'
import SettingsSidebar from '@/components/settings/SettingsSidebar.vue'

const route = useRoute()
const router = useRouter()

const currentTab = computed(() => {
  const path = route.path
  if (path.endsWith('/libraries')) return 'libraries'
  if (path.endsWith('/indexers')) return 'indexers'
  if (path.endsWith('/name-templates')) return 'name-templates'
  if (path.endsWith('/downloaders')) return 'downloaders'
  if (path.endsWith('/policies')) return 'policies'
  return 'general'
})

const navigateToTab = (
  tab: 'general' | 'libraries' | 'indexers' | 'name-templates' | 'downloaders' | 'policies',
) => {
  router.push(`/settings/${tab}`)
}
</script>

<template>
  <div class="settings-layout">
    <!-- Desktop Sidebar -->
    <SettingsSidebar />

    <!-- Main Content Area -->
    <div class="settings-content">
      <div v-if="false" class="settings-header">
        <h2 class="text-2xl font-semibold">Settings</h2>

        <!-- Mobile Navigation Tabs -->
        <div class="mobile-tabs">
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
          <Button
            :label="'Indexers'"
            :severity="currentTab === 'indexers' ? 'primary' : 'secondary'"
            :text="currentTab !== 'indexers'"
            @click="navigateToTab('indexers')"
          />
          <Button
            :label="'Name Templates'"
            :severity="currentTab === 'name-templates' ? 'primary' : 'secondary'"
            :text="currentTab !== 'name-templates'"
            @click="navigateToTab('name-templates')"
          />
          <Button
            :label="'Downloaders'"
            :severity="currentTab === 'downloaders' ? 'primary' : 'secondary'"
            :text="currentTab !== 'downloaders'"
            @click="navigateToTab('downloaders')"
          />
          <Button
            :label="'Policies'"
            :severity="currentTab === 'policies' ? 'primary' : 'secondary'"
            :text="currentTab !== 'policies'"
            @click="navigateToTab('policies')"
          />
        </div>
      </div>

      <!-- Content Area -->
      <div class="settings-main">
        <router-view />
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-layout {
  display: flex;
  flex: 1;
  /* min-height: 0; */
  height: 100%;
  gap: 1rem;
}

.settings-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  min-width: 0;
}

.settings-header {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.mobile-tabs {
  display: flex;
  gap: 0.5rem;
  border-bottom: 1px solid var(--p-border-color);
  padding-bottom: 1rem;
  background: var(--p-card-background);
  border-radius: 8px;
  padding: 1rem;
  box-shadow: var(--p-shadow-1);
}

.settings-main {
  flex: 1;
}

/* Desktop: Show sidebar, hide mobile tabs */
@media (min-width: 1024px) {
  .mobile-tabs {
    display: none;
  }
}

/* Mobile: Hide sidebar, show mobile tabs */
@media (max-width: 1023px) {
  .settings-layout {
    flex-direction: column;
    padding: 0.75rem;
  }

  .settings-content {
    width: 100%;
  }
}
</style>
