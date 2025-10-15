<template>
  <header class="app-header">
    <div class="flex items-center gap-5">
      <Button
        v-if="showMobileMenu"
        class="menu-btn"
        icon="pi pi-bars"
        variant="text"
        @click="mobileSidebarVisible = true"
      />

      <IconField class="search-field">
        <InputIcon :class="PrimeIcons.SEARCH" />
        <InputText
          v-model="query"
          class="w-full"
          placeholder="Search"
          variant="filled"
          size="small"
          @keyup.enter="emit('search', query)"
        />
      </IconField>

      <nav v-if="!showMobileMenu" class="flex items-center gap-5">
        <RouterLink
          v-for="item in links"
          :key="item.to"
          :to="item.to"
          class="text-decoration-none text-inherit px-2 py-1 rounded"
          :class="{ 'bg-emphasis': route.path === item.to }"
        >
          {{ item.label }}
        </RouterLink>
      </nav>
    </div>

    <div class="flex items-center gap-2">
      <div class="avatar-wrap">
        <Avatar label="S" shape="circle" />
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import IconField from 'primevue/iconfield'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import Avatar from 'primevue/avatar'
import Button from 'primevue/button'
import { PrimeIcons } from '@/icons'

const route = useRoute()

const query = ref('')
const windowWidth = ref(window.innerWidth)

const mobileSidebarVisible = defineModel<boolean>('mobileSidebarVisible', { default: false })

// Show mobile menu when window is too narrow to fit navigation links
const showMobileMenu = computed(() => {
  return windowWidth.value < 768 // Adjust breakpoint as needed
})

const links = ref([
  { label: 'Home', to: '/' },
  { label: 'Library', to: '/library' },
  { label: 'Requests', to: '/requests' },
])

const emit = defineEmits<{
  (e: 'search', query: string): void
}>()

const handleResize = () => {
  windowWidth.value = window.innerWidth
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.app-header {
  position: sticky;
  top: 0;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.75rem;
  background: var(--bg);
}

.menu-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2.25rem;
  height: 2.25rem;
  border: none;
  background: transparent;
  font-size: 1.25rem;
  cursor: pointer;
}

.search-field {
  width: 200px; /* Mobile: smaller */
}

@media (min-width: 640px) {
  .search-field {
    width: 300px; /* Tablet: medium */
  }
}

@media (min-width: 1024px) {
  .search-field {
    width: 400px; /* Desktop: larger */
  }
}

.avatar-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
