<template>
  <header class="app-header shadow-[var(--p-shadow-1)]">
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
          :dt="{ root: { borderRadius: '32px' } }"
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
        <Avatar
          class="cursor-pointer"
          :label="avatarLabel"
          shape="circle"
          @click="handleAvatarClick"
        />

        <Menu ref="menu" :model="items" class="w-full md:w-60" :popup="true">
          <template #start>
            <div
              class="relative overflow-hidden w-full rounded-border bg-transparent flex items-center p-2 pl-4 transition-colors duration-200"
            >
              <Avatar :label="avatarLabel" class="mr-2" shape="circle" />
              <span class="inline-flex flex-col items-start">
                <span class="font-bold">{{ user?.name }}</span>
                <span class="text-sm">{{ user?.email }}</span>
              </span>
            </div>
          </template>

          <template #item="{ item, props }">
            <component
              :is="item.to ? 'router-link' : 'a'"
              class="flex items-center"
              v-bind="props.action"
              :to="item.to"
            >
              <span :class="item.icon" />
              <span>{{ item.label }}</span>
              <Badge v-if="item.badge" class="ml-auto" :value="item.badge" />
              <span
                v-if="item.shortcut"
                class="ml-auto border border-surface rounded bg-emphasis text-muted-color text-xs p-1"
                >{{ item.shortcut }}</span
              >
            </component>
          </template>
        </Menu>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, RouterLink, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import IconField from 'primevue/iconfield'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import Avatar from 'primevue/avatar'
import Button from 'primevue/button'
import Menu, { type MenuMethods } from 'primevue/menu'
import Badge from 'primevue/badge'
import { PrimeIcons } from '@/icons'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const { user } = storeToRefs(authStore)
const avatarLabel = computed(() => user.value?.name?.[0])

const query = ref('')
const windowWidth = ref(window.innerWidth)
const menu = ref<MenuMethods | null>(null)

const mobileSidebarVisible = defineModel<boolean>('mobileSidebarVisible', { default: false })

// Show mobile menu when window is too narrow to fit navigation links
const showMobileMenu = computed(() => {
  return windowWidth.value < 768 // Adjust breakpoint as needed
})

const links = ref([
  { label: 'Home', to: '/' },
  { label: 'Library', to: '/library' },
  { label: 'Downloads', to: '/downloads' },
  { label: 'Requests', to: '/requests' },
])

const items = ref([
  {
    separator: true,
  },
  {
    label: 'Messages',
    icon: PrimeIcons.INBOX,
    badge: 2,
  },
  {
    label: 'Settings',
    icon: PrimeIcons.COG,
    to: '/settings',
  },
  {
    label: 'Logout',
    icon: PrimeIcons.SIGN_OUT,
    command: () => {
      const auth = useAuthStore()
      auth.logout()
      router.push('/login')
    },
  },
])

const emit = defineEmits<{
  (e: 'search', query: string): void
}>()

const handleResize = () => {
  windowWidth.value = window.innerWidth
}

const handleAvatarClick = (event: MouseEvent) => {
  menu.value?.toggle(event)
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
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.75rem 1rem;
  background: var(--bg);
  border-radius: 32px;
  margin: 1rem;
  margin-bottom: 0;
  flex-shrink: 0;
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
