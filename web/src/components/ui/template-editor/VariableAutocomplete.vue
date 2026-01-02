<script setup lang="ts">
import { computed, ref, watch, onMounted, onUnmounted } from 'vue'
import { useTemplateVariables, type TemplateVariable } from '@/composables/useTemplateVariables'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command'
import { ScrollArea } from '@/components/ui/scroll-area'

interface Props {
  query: string
  position: { top: number; left: number }
  mediaType?: 'movie' | 'series'
}

const props = withDefaults(defineProps<Props>(), {
  mediaType: 'movie',
})

const emit = defineEmits<{
  select: [variable: TemplateVariable]
  close: []
}>()

const { variablesByNamespace, searchVariables, isLoading } = useTemplateVariables({
  mediaType: props.mediaType,
})

const searchQuery = ref(props.query)

// Update search when query prop changes
watch(
  () => props.query,
  (newQuery) => {
    searchQuery.value = newQuery
  },
)

// Filter variables based on search
const filteredVariables = computed(() => {
  if (!searchQuery.value) {
    return variablesByNamespace.value
  }

  const results = searchVariables(searchQuery.value)

  // Group by namespace
  const grouped: Record<string, TemplateVariable[]> = {}
  for (const variable of results) {
    if (!grouped[variable.namespace]) {
      grouped[variable.namespace] = []
    }
    grouped[variable.namespace].push(variable)
  }

  return grouped
})

// Order of namespaces for display
const namespaceOrder = ['Shortcuts', 'Media', 'Quality', 'Candidate', 'MediaInfo']

const orderedNamespaces = computed(() => {
  const groups = filteredVariables.value
  return namespaceOrder.filter((ns) => groups[ns] && groups[ns].length > 0)
})

const hasResults = computed(() => {
  return orderedNamespaces.value.length > 0
})

function handleSelect(variable: TemplateVariable) {
  emit('select', variable)
}

function handleKeyDown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    emit('close')
  }
}

// Close on click outside
function handleClickOutside(event: MouseEvent) {
  const target = event.target as HTMLElement
  if (!target.closest('.variable-autocomplete')) {
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
  <Teleport to="body">
    <div
      class="variable-autocomplete fixed z-[100] w-72 rounded-md border bg-popover p-0 text-popover-foreground shadow-md outline-none"
      :style="{
        top: `${position.top}px`,
        left: `${position.left}px`,
      }"
    >
      <Command class="rounded-lg border-0 shadow-none" :should-filter="false">
        <CommandInput
          v-model="searchQuery"
          placeholder="Search variables..."
          class="h-9"
          autofocus
        />
        <CommandList>
          <ScrollArea class="h-[300px]">
            <CommandEmpty v-if="!isLoading && !hasResults">
              No variables found.
            </CommandEmpty>

            <div v-if="isLoading" class="py-6 text-center text-sm text-muted-foreground">
              Loading variables...
            </div>

            <template v-for="namespace in orderedNamespaces" :key="namespace">
              <CommandGroup :heading="namespace">
                <CommandItem
                  v-for="variable in filteredVariables[namespace]"
                  :key="variable.path"
                  :value="variable.path"
                  class="flex items-center justify-between gap-2"
                  @select="handleSelect(variable)"
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
                </CommandItem>
              </CommandGroup>
            </template>
          </ScrollArea>
        </CommandList>
      </Command>
    </div>
  </Teleport>
</template>

<style scoped>
.variable-autocomplete {
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


