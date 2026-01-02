<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { TEMPLATE_FUNCTIONS } from '@/composables/useTemplateVariables'
import type { Token } from './TemplateTokenEditor.vue'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuLabel,
} from '@/components/ui/dropdown-menu'
import { WrapText, X, Trash2, Sparkles } from 'lucide-vue-next'

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

onMounted(() => {
  document.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown)
})
</script>

<template>
  <div
    class="token-context-menu fixed z-50"
    :style="{
      top: `${position.y}px`,
      left: `${position.x}px`,
    }"
  >
    <DropdownMenu :open="true" @update:open="(open) => !open && emit('close')">
      <DropdownMenuContent class="w-56" align="start" :side-offset="0">
        <DropdownMenuLabel class="font-mono text-xs">
          {{ token.func ? `${token.func} ${token.value}` : token.value }}
        </DropdownMenuLabel>
        <DropdownMenuSeparator />

        <!-- Wrap with function submenu -->
        <DropdownMenuSub>
          <DropdownMenuSubTrigger>
            <Sparkles class="mr-2 h-4 w-4" />
            <span>Wrap with function</span>
          </DropdownMenuSubTrigger>
          <DropdownMenuSubContent>
            <DropdownMenuItem
              v-for="func in TEMPLATE_FUNCTIONS"
              :key="func.name"
              @click="handleWrap(func.name as 'clean' | 'sanitize')"
            >
              <WrapText class="mr-2 h-4 w-4" />
              <div class="flex flex-col gap-0.5">
                <span class="font-mono">{{ func.name }}</span>
                <span class="text-xs text-muted-foreground">{{ func.description }}</span>
              </div>
            </DropdownMenuItem>
          </DropdownMenuSubContent>
        </DropdownMenuSub>

        <!-- Remove function (only if wrapped) -->
        <DropdownMenuItem v-if="hasFunction" @click="handleRemoveFunction">
          <X class="mr-2 h-4 w-4" />
          <span>Remove function</span>
        </DropdownMenuItem>

        <DropdownMenuSeparator />

        <!-- Delete token -->
        <DropdownMenuItem class="text-destructive focus:text-destructive" @click="handleDelete">
          <Trash2 class="mr-2 h-4 w-4" />
          <span>Delete</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  </div>
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


