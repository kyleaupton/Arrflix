<script setup lang="ts">
import { NodeViewWrapper, NodeViewContent } from '@tiptap/vue-3'
import { computed, ref } from 'vue'
import { X } from 'lucide-vue-next'
import FieldPicker from './FieldPicker.vue'

interface Props {
  node: {
    attrs: {
      field: string
    }
  }
  editor: any
  getPos: () => number
  updateAttributes: (attrs: Record<string, any>) => void
  deleteNode: () => void
}

const props = defineProps<Props>()

const showFieldPicker = ref(false)
const fieldPickerAnchor = ref<HTMLElement | null>(null)

/**
 * Format the field label for display
 * Removes leading dot: ".Release.Edition" -> "Release.Edition"
 */
const fieldLabel = computed(() => {
  return props.node.attrs.field.startsWith('.')
    ? props.node.attrs.field.slice(1)
    : props.node.attrs.field
})

function handleFieldClick(event: MouseEvent) {
  fieldPickerAnchor.value = event.target as HTMLElement
  showFieldPicker.value = true
}

function handleFieldSelect(field: string) {
  props.updateAttributes({ field })
  showFieldPicker.value = false
}

function handleDelete() {
  props.deleteNode()
}
</script>

<template>
  <NodeViewWrapper
    as="span"
    class="conditional-block inline-flex items-baseline gap-1 rounded-md border border-dashed border-primary/30 bg-primary/5 px-2 py-0.5 mx-0.5"
  >
    <button
      class="conditional-header inline-flex items-center gap-1 text-xs font-mono text-primary/70 hover:text-primary cursor-pointer"
      @click="handleFieldClick"
      type="button"
    >
      <span class="font-semibold">if</span>
      <span>{{ fieldLabel }}</span>
      <span class="text-[10px]">▼</span>
    </button>
    <span class="text-primary/40 text-xs">❰</span>
    <NodeViewContent class="conditional-content inline min-w-[3ch] px-1" data-placeholder="..." />
    <span class="text-primary/40 text-xs">❱</span>
    <button
      class="conditional-delete opacity-0 hover:opacity-100 transition-opacity text-destructive hover:text-destructive/80 ml-1"
      @click="handleDelete"
      type="button"
      title="Delete condition"
    >
      <X class="h-3 w-3" />
    </button>

    <!-- Field Picker Popover -->
    <Teleport to="body">
      <FieldPicker
        v-if="showFieldPicker"
        :anchor="fieldPickerAnchor"
        @select="handleFieldSelect"
        @close="showFieldPicker = false"
      />
    </Teleport>
  </NodeViewWrapper>
</template>

<style scoped>
.conditional-block:hover .conditional-delete {
  opacity: 1;
}

/* Show placeholder when content is empty */
.conditional-content:empty::before {
  content: attr(data-placeholder);
  color: hsl(var(--muted-foreground));
  opacity: 0.5;
  pointer-events: none;
}

/* Subtle highlight on hover to show interactivity */
.conditional-block:hover {
  background-color: hsl(var(--primary) / 0.08);
  border-color: hsl(var(--primary) / 0.4);
}
</style>

