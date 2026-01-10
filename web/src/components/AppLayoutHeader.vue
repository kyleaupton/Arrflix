<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from 'lucide-vue-next'
import { KbdGroup, Kbd } from '@/components/ui/kbd'
import { SidebarTrigger } from '@/components/ui/sidebar'
import Input from '@/components/ui/input/Input.vue'
import { Popover, PopoverContent, PopoverAnchor } from '@/components/ui/popover'
import SearchPopover from '@/components/search/SearchPopover.vue'
import { useSearch } from '@/composables/useSearch'
import { isMac } from '@/lib/platform'

const router = useRouter()
const { query, results, totalResults, isLoading, clear } = useSearch()

const searchContainerRef = ref<HTMLElement | null>(null)
const isOpen = ref(false)
const inputHasFocus = ref(false)

// Get the actual input element
const getInputElement = () => {
  return searchContainerRef.value?.querySelector('input') as HTMLInputElement | null
}

// Handle input focus
const onFocus = () => {
  inputHasFocus.value = true
  isOpen.value = true
}

// Handle input blur
const onBlur = (e: FocusEvent) => {
  inputHasFocus.value = false

  // Check if focus is moving to the popover content
  const relatedTarget = e.relatedTarget as HTMLElement | null
  const popoverContent = document.querySelector('[data-slot="popover-content"]')

  // If focus is moving to popover content, keep it open
  if (relatedTarget && popoverContent?.contains(relatedTarget)) {
    return
  }

  // Small delay to allow click events on popover items to fire first
  setTimeout(() => {
    if (!inputHasFocus.value) {
      isOpen.value = false
    }
  }, 150)
}

// Close popover when navigating
watch(
  () => router.currentRoute.value.path,
  () => {
    isOpen.value = false
    clear()
  },
)

// Handle result selection
const onSelect = () => {
  isOpen.value = false
  clear()
}

// Handle Enter key to go to search page
const onKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Enter' && query.value.length >= 2) {
    router.push({ path: '/search', query: { q: query.value } })
    isOpen.value = false
    clear()
  }
  if (e.key === 'Escape') {
    isOpen.value = false
    getInputElement()?.blur()
  }
}

// Global keyboard shortcut (Cmd/Ctrl + K)
const handleGlobalKeydown = (e: KeyboardEvent) => {
  if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
    e.preventDefault()
    // Focus first, which will trigger onFocus and open the popover
    getInputElement()?.focus()
  }
}

// Prevent popover's interact-outside from closing - let blur handler manage it
const handleInteractOutside = (e: Event) => {
  e.preventDefault()
}

onMounted(() => {
  document.addEventListener('keydown', handleGlobalKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleGlobalKeydown)
})
</script>

<template>
  <header
    class="z-10 fixed w-full flex h-16 shrink-0 items-center gap-2 border-b bg-background transition-[width,height] ease-linear"
  >
    <div class="flex items-center gap-2 px-4">
      <SidebarTrigger class="-ml-1" />

      <Popover v-model:open="isOpen" :modal="false">
        <PopoverAnchor as-child>
          <div ref="searchContainerRef" class="relative">
            <Input
              v-model="query"
              placeholder="Search..."
              class="pl-8 w-64"
              @focus="onFocus"
              @blur="onBlur"
              @keydown="onKeydown"
            />
            <Search
              class="pointer-events-none absolute top-1/2 left-2 size-4 -translate-y-1/2 opacity-50 select-none"
            />
            <div class="absolute top-1/2 -translate-y-1/2 right-1.5 hidden gap-1 sm:flex">
              <KbdGroup>
                <Kbd class="border">{{ isMac ? '⌘' : 'Ctrl' }}</Kbd>
                <Kbd class="border">K</Kbd>
              </KbdGroup>
            </div>
          </div>
        </PopoverAnchor>
        <PopoverContent
          class="p-0 w-80"
          align="start"
          :side-offset="8"
          @open-auto-focus.prevent
          @close-auto-focus.prevent
          @interact-outside="handleInteractOutside"
        >
          <SearchPopover
            :query="query"
            :results="results"
            :total-results="totalResults"
            :is-loading="isLoading"
            @select="onSelect"
          />
        </PopoverContent>
      </Popover>
    </div>
  </header>
</template>
