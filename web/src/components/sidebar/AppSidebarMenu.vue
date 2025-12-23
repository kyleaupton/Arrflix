<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import Menu, { type MenuMethods } from 'primevue/menu'
import Badge from 'primevue/badge'
import IconField from 'primevue/iconfield'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import Avatar from 'primevue/avatar'
import { navigationItems } from '@/config/navigation'
import { PrimeIcons } from '@/icons'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const { user } = storeToRefs(authStore)
const avatarLabel = computed(() => user.value?.name?.[0] || 'U')
const query = ref('')
const menu = ref<MenuMethods | null>(null)

// Make navigation items reactive
const items = ref(navigationItems)

// Avatar menu items
const avatarMenuItems = ref([
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

// Emit event to close mobile drawer when navigating
const emit = defineEmits<{
  (e: 'navigate'): void
  (e: 'search', query: string): void
}>()

// Reactive function that checks if a route is active
// Using route.path directly makes it reactive to route changes
const isActive = (routePath?: string): boolean => {
  if (!routePath) return false
  if (routePath === '/') return route.path === '/'
  return route.path.startsWith(routePath)
}

const handleNavigation = (routePath: string) => {
  // Only navigate if we're not already on this route
  if (route.path !== routePath) {
    router.push(routePath).catch(() => {
      // Ignore navigation errors (e.g., navigating to same route)
    })
  }
  
  // Always emit navigate to close mobile drawer
  emit('navigate')
}

const handleAvatarClick = (event: MouseEvent) => {
  menu.value?.toggle(event)
}

const handleSearch = () => {
  emit('search', query.value)
}
</script>

<template>
  <div class="sidebar-menu">
    <!-- Brand -->
    <div class="sidebar-brand">
      <svg
        width="35"
        height="40"
        viewBox="0 0 35 40"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        class="brand-icon"
      >
        <path
          d="M25.87 18.05L23.16 17.45L25.27 20.46V29.78L32.49 23.76V13.53L29.18 14.73L25.87 18.04V18.05ZM25.27 35.49L29.18 31.58V27.67L25.27 30.98V35.49ZM20.16 17.14H20.03H20.17H20.16ZM30.1 5.19L34.89 4.81L33.08 12.33L24.1 15.67L30.08 5.2L30.1 5.19ZM5.72 14.74L2.41 13.54V23.77L9.63 29.79V20.47L11.74 17.46L9.03 18.06L5.72 14.75V14.74ZM9.63 30.98L5.72 27.67V31.58L9.63 35.49V30.98ZM4.8 5.2L10.78 15.67L1.81 12.33L0 4.81L4.79 5.19L4.8 5.2ZM24.37 21.05V34.59L22.56 37.29L20.46 39.4H14.44L12.34 37.29L10.53 34.59V21.05L12.42 18.23L17.45 26.8L22.48 18.23L24.37 21.05ZM22.85 0L22.57 0.69L17.45 13.08L12.33 0.69L12.05 0H22.85Z"
          fill="var(--p-primary-color)"
        />
        <path
          d="M30.69 4.21L24.37 4.81L22.57 0.69L22.86 0H26.48L30.69 4.21ZM23.75 5.67L22.66 3.08L18.05 14.24V17.14H19.7H20.03H20.16H20.2L24.1 15.7L30.11 5.19L23.75 5.67ZM4.21002 4.21L10.53 4.81L12.33 0.69L12.05 0H8.43002L4.22002 4.21H4.21002ZM21.9 17.4L20.6 18.2H14.3L13 17.4L12.4 18.2L12.42 18.23L17.45 26.8L22.48 18.23L22.5 18.2L21.9 17.4ZM4.79002 5.19L10.8 15.7L14.7 17.14H14.74H15.2H16.85V14.24L12.24 3.09L11.15 5.68L4.79002 5.2V5.19Z"
          fill="var(--p-text-color)"
        />
      </svg>
      <span class="brand-text">SNAGGLE</span>
    </div>

    <!-- Search -->
    <div class="sidebar-search">
      <IconField class="w-full">
        <InputIcon :class="PrimeIcons.SEARCH" />
        <InputText
          v-model="query"
          class="w-full"
          placeholder="Search"
          variant="filled"
          size="small"
          @keyup.enter="handleSearch"
        />
      </IconField>
    </div>

    <!-- Navigation Menu -->
    <Menu :model="items" class="w-full sidebar-menu-component flex-1">
      <template #item="{ item, props }">
        <a
          v-if="item.route"
          v-ripple
          v-bind="props.action"
          class="menu-item-link"
          :class="{ 'menu-item-active': isActive(item.route) }"
          href="#"
          @click.prevent="handleNavigation(item.route)"
        >
          <span v-if="item.icon" :class="item.icon" />
          <span class="menu-item-label">{{ item.label }}</span>
          <Badge v-if="item.badge" class="ml-auto" :value="item.badge" />
          <span
            v-if="item.shortcut"
            class="ml-auto border border-surface rounded bg-emphasis text-muted-color text-xs p-1"
            >{{ item.shortcut }}</span
          >
        </a>
        <a
          v-else-if="item.url"
          v-ripple
          v-bind="props.action"
          class="menu-item-link"
          :href="item.url"
          :target="item.target"
        >
          <span v-if="item.icon" :class="item.icon" />
          <span class="menu-item-label">{{ item.label }}</span>
          <Badge v-if="item.badge" class="ml-auto" :value="item.badge" />
        </a>
        <a
          v-else
          v-ripple
          v-bind="props.action"
          class="menu-item-link"
        >
          <span v-if="item.icon" :class="item.icon" />
          <span class="menu-item-label">{{ item.label }}</span>
          <Badge v-if="item.badge" class="ml-auto" :value="item.badge" />
        </a>
      </template>
    </Menu>

    <!-- Avatar at bottom -->
    <div class="sidebar-footer">
      <div class="avatar-container">
        <Avatar
          class="cursor-pointer"
          :label="avatarLabel"
          shape="circle"
          @click="handleAvatarClick"
        />

        <Menu ref="menu" :model="avatarMenuItems" class="w-full md:w-60" :popup="true">
          <template #start>
            <div class="avatar-menu-header">
              <Avatar :label="avatarLabel" class="mr-2" shape="circle" />
              <span class="inline-flex flex-col items-start">
                <span class="font-bold">{{ user?.name }}</span>
                <span class="text-sm">{{ user?.email }}</span>
              </span>
            </div>
          </template>

          <template #item="{ item, props }">
            <a
              v-if="item.to"
              class="flex items-center"
              v-bind="props.action"
              href="#"
              @click.prevent="handleNavigation(item.to)"
            >
              <span :class="item.icon" />
              <span>{{ item.label }}</span>
              <Badge v-if="item.badge" class="ml-auto" :value="item.badge" />
              <span
                v-if="item.shortcut"
                class="ml-auto border border-surface rounded bg-emphasis text-muted-color text-xs p-1"
                >{{ item.shortcut }}</span
              >
            </a>
            <a
              v-else
              class="flex items-center"
              v-bind="props.action"
            >
              <span :class="item.icon" />
              <span>{{ item.label }}</span>
              <Badge v-if="item.badge" class="ml-auto" :value="item.badge" />
            </a>
          </template>
        </Menu>
      </div>
    </div>
  </div>
</template>

<style scoped>
.sidebar-menu {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: transparent;
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem;
  border-bottom: 1px solid var(--p-border-color);
  flex-shrink: 0;
  background: transparent;
}

.sidebar-search {
  padding: 1rem;
  border-bottom: 1px solid var(--p-border-color);
  flex-shrink: 0;
  background: transparent;
}

.brand-icon {
  height: 2rem;
}

.brand-text {
  font-size: 1.25rem;
  font-weight: 600;
}

.menu-item-link {
  display: flex;
  align-items: center;
  padding: 0.75rem 1rem;
  text-decoration: none;
  color: inherit;
  transition: all var(--p-transition-duration);
  cursor: pointer;
}

.menu-item-link:hover {
  background: var(--p-emphasis-background);
}

.menu-item-active {
  background: var(--p-primary-color);
  color: var(--p-primary-contrast-color);
}

.menu-item-active:hover {
  background: var(--p-primary-color-hover);
}

.menu-item-label {
  margin-left: 0.5rem;
}

/* Menu component styling */
.sidebar-menu-component {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: transparent;
}

.sidebar-menu-component :deep(.p-menu-root) {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: transparent;
}

.sidebar-menu-component :deep(.p-menu) {
  background: transparent;
}

.sidebar-menu-component :deep(.p-menu-list) {
  padding: 0;
  gap: 0.25rem;
  flex: 1;
  background: transparent;
}

.sidebar-menu-component :deep(.p-menu-item) {
  margin: 0;
  background: transparent;
}

.sidebar-menu-component :deep(.p-submenu-header) {
  padding: 0.75rem 1rem;
  font-weight: 600;
  color: var(--p-text-muted-color);
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 0.05em;
  background: transparent;
}

.sidebar-menu-component :deep(.p-submenu-content) {
  padding: 0.25rem 0;
  background: transparent;
}

.sidebar-menu-component :deep(.p-submenu-content .p-menu-item) {
  padding-left: 2rem;
}

.sidebar-footer {
  margin-top: auto;
  padding: 1rem;
  border-top: 1px solid var(--p-border-color);
  flex-shrink: 0;
  background: transparent;
}

.avatar-container {
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-menu-header {
  display: flex;
  align-items: center;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--p-border-color);
}
</style>
