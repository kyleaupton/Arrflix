<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useTemplateVariables, type TemplateVariable } from '@/composables/useTemplateVariables'
import { ScrollArea } from '@/components/ui/scroll-area'

interface Props {
  mediaType?: 'movie' | 'series'
}

const props = withDefaults(defineProps<Props>(), {
  mediaType: 'movie',
})

const emit = defineEmits<{
  select: [field: string]
  close: []
}>()

const { variablesByNamespace, isLoading } = useTemplateVariables({
  mediaType: props.mediaType,
})

const searchQuery = ref('')

// Filter variables based on search
const filteredVariables = computed(() => {
  if (!searchQuery.value) {
    return variablesByNamespace.value
  }

  const q = searchQuery.value.toLowerCase()
  const filtered: Record<string, TemplateVariable[]> = {}

  for (const [namespace, variables] of Object.entries(variablesByNamespace.value)) {
    const matched = variables.filter(
      (v) => v.label.toLowerCase().includes(q) || v.path.toLowerCase().includes(q),
    )
    if (matched.length > 0) {
      filtered[namespace] = matched
    }
  }

  return filtered
})

// Order of namespaces for display
const namespaceOrder = ['Media', 'Quality', 'Release', 'Candidate', 'MediaInfo']

const orderedNamespaces = computed(() => {
  const groups = filteredVariables.value
  return namespaceOrder.filter((ns) => groups[ns] && groups[ns].length > 0)
})

const hasResults = computed(() => {
  return orderedNamespaces.value.length > 0
})

function selectField(variable: TemplateVariable, event?: MouseEvent) {
  // Stop event propagation to prevent closing the dialog
  if (event) {
    event.stopPropagation()
    event.preventDefault()
  }
  emit('select', variable.path)
}

function handleKeyDown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    event.stopPropagation()
    event.preventDefault()
    emit('close')
  }
  // Let other keys pass through normally for the search input
}

function handleClickOutside(event: MouseEvent) {
  const target = event.target as HTMLElement
  // Check if clicking outside the field picker (but allow clicking inside)
  if (!target.closest('.field-picker')) {
    event.stopPropagation()
    emit('close')
  }
}

onMounted(() => {
  document.addEventListener('mousedown', handleClickOutside)
  document.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('mousedown', handleClickOutside)
  document.removeEventListener('keydown', handleKeyDown)
})
</script>

<template>
  <div
    class="field-picker w-80 rounded-md border bg-popover p-2 text-popover-foreground shadow-md"
    @click.stop
    @mousedown.stop
  >
    <div class="mb-2">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search fields..."
        class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        autofocus
        @click.stop
        @keydown.stop
      />
    </div>

    <div v-if="isLoading" class="py-6 text-center text-sm text-muted-foreground">
      Loading fields...
    </div>

    <div v-else-if="!hasResults" class="py-6 text-center text-sm text-muted-foreground">
      No fields found.
    </div>

    <ScrollArea v-else class="h-[300px]">
      <div v-for="namespace in orderedNamespaces" :key="namespace" class="mb-3 last:mb-0">
        <div class="px-2 py-1.5 text-xs font-semibold text-muted-foreground">
          {{ namespace }}
        </div>
        <button
          v-for="variable in filteredVariables[namespace]"
          :key="variable.path"
          type="button"
          class="flex w-full items-center justify-between gap-2 rounded-sm px-2 py-1.5 text-sm hover:bg-accent cursor-pointer text-left"
          @click="selectField(variable, $event)"
          @mousedown.stop
        >
          <div class="flex flex-col gap-0.5">
            <span class="font-medium">{{ variable.label }}</span>
            <code class="text-xs text-muted-foreground">{{ variable.path }}</code>
          </div>
          <span
            v-if="variable.postDownloadOnly"
            class="text-xs text-amber-500"
            title="Only available after download"
          >
            post-DL
          </span>
        </button>
      </div>
    </ScrollArea>
  </div>
</template>

<style scoped>
.field-picker {
  animation: fadeIn 0.15s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>

