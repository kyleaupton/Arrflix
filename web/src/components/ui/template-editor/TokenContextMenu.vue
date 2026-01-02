<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { TEMPLATE_FUNCTIONS } from '@/composables/useTemplateVariables'
import type { Token } from './TemplateTokenEditor.vue'
import { WrapText, X, Trash2, Sparkles, ChevronRight } from 'lucide-vue-next'

interface Props {
  position: { x: number; y: number }
  token: Token
}

const props = defineProps<Props>()

const emit = defineEmits<{
  wrap: [funcName: 'clean' | 'sanitize']
  removeFunction: []
  delete: []
  close: []
}>()

const hasFunction = computed(() => !!props.token.func)
const showFunctionSubmenu = ref(false)

function handleWrap(funcName: 'clean' | 'sanitize') {
  emit('wrap', funcName)
}

function handleRemoveFunction() {
  emit('removeFunction')
}

function handleDelete() {
  emit('delete')
}

function handleKeyDown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    emit('close')
  }
}

function handleClickOutside(event: MouseEvent) {
  const target = event.target as HTMLElement
  if (!target.closest('.token-context-menu')) {
    emit('close')
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeyDown)
  // Delay adding click listener to avoid immediate close
  setTimeout(() => {
    document.addEventListener('click', handleClickOutside)
  }, 10)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown)
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <Teleport to="body">
    <div
      class="token-context-menu fixed z-[200] min-w-56 rounded-md border bg-popover p-1 text-popover-foreground shadow-md"
      :style="{
        top: `${position.y}px`,
        left: `${position.x}px`,
      }"
    >
      <!-- Token label -->
      <div class="px-2 py-1.5 text-xs font-semibold font-mono text-muted-foreground">
        {{
          token.func
            ? `${token.func} ${token.value.replace(/^\./, '')}`
            : `${token.value.replace(/^\./, '')}`
        }}
      </div>

      <div class="h-px bg-border my-1" />

      <!-- Wrap with function -->
      <div
        class="relative"
        @mouseenter="showFunctionSubmenu = true"
        @mouseleave="showFunctionSubmenu = false"
      >
        <button
          class="flex w-full cursor-default items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground"
        >
          <Sparkles class="mr-2 h-4 w-4" />
          <span class="flex-1 text-left">Wrap with function</span>
          <ChevronRight class="h-4 w-4" />
        </button>

        <!-- Submenu -->
        <div
          v-if="showFunctionSubmenu"
          class="absolute left-full top-0 ml-1 min-w-48 rounded-md border bg-popover p-1 text-popover-foreground shadow-md"
        >
          <button
            v-for="func in TEMPLATE_FUNCTIONS"
            :key="func.name"
            class="flex w-full cursor-default items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground"
            @click="handleWrap(func.name as 'clean' | 'sanitize')"
          >
            <WrapText class="mr-2 h-4 w-4" />
            <div class="flex flex-col gap-0.5">
              <span class="font-mono">{{ func.name }}</span>
              <span class="text-xs text-muted-foreground">{{ func.description }}</span>
            </div>
          </button>
        </div>
      </div>

      <!-- Remove function (only if wrapped) -->
      <button
        v-if="hasFunction"
        class="flex w-full cursor-default items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground"
        @click="handleRemoveFunction"
      >
        <X class="mr-2 h-4 w-4" />
        <span>Remove function</span>
      </button>

      <div class="h-px bg-border my-1" />

      <!-- Delete token -->
      <button
        class="flex w-full cursor-default items-center rounded-sm px-2 py-1.5 text-sm outline-none text-destructive hover:bg-destructive/10"
        @click="handleDelete"
      >
        <Trash2 class="mr-2 h-4 w-4" />
        <span>Delete</span>
      </button>
    </div>
  </Teleport>
</template>

<style scoped>
.token-context-menu {
  animation: fadeIn 0.1s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}
</style>
