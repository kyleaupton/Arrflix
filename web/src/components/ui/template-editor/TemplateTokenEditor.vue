<script setup lang="ts">
import { ref, watch, onBeforeUnmount } from 'vue'
import { useEditor, EditorContent, type Editor } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import { VariableMention } from './extensions/VariableMention'
import { variableSuggestion } from './extensions/variableSuggestion'
import TokenContextMenu from './TokenContextMenu.vue'
import { cn } from '@/lib/utils'

export interface Token {
  type: 'text' | 'variable'
  value: string
  func?: 'clean' | 'sanitize'
}

interface Props {
  modelValue: string
  placeholder?: string
  mediaType?: 'movie' | 'series'
  disabled?: boolean
  class?: string
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: 'Type { to insert a variable...',
  mediaType: 'movie',
  disabled: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

// Context menu state
const showContextMenu = ref(false)
const contextMenuPosition = ref({ x: 0, y: 0 })
const selectedNodePos = ref<number | null>(null)

// Template parsing regex - matches {{.Var}} or {{func .Var}}
const TEMPLATE_REGEX = /\{\{(clean|sanitize)?\s*([.\w]+)\}\}/g

/**
 * Parse template string to Tiptap JSON
 */
function parseTemplateToTiptap(template: string) {
  if (!template) {
    return {
      type: 'doc',
      content: [{ type: 'paragraph', content: [] }],
    }
  }

  const content: Array<Record<string, unknown>> = []
  let lastIndex = 0

  TEMPLATE_REGEX.lastIndex = 0
  let match

  while ((match = TEMPLATE_REGEX.exec(template)) !== null) {
    // Add text before match
    if (match.index > lastIndex) {
      const text = template.slice(lastIndex, match.index)
      if (text) {
        content.push({ type: 'text', text })
      }
    }

    // Add variable mention
    content.push({
      type: 'variableMention',
      attrs: {
        path: match[2],
        func: match[1] || null,
      },
    })

    lastIndex = match.index + match[0].length
  }

  // Add remaining text
  if (lastIndex < template.length) {
    const text = template.slice(lastIndex)
    if (text) {
      content.push({ type: 'text', text })
    }
  }

  return {
    type: 'doc',
    content: [
      {
        type: 'paragraph',
        content: content.length > 0 ? content : undefined,
      },
    ],
  }
}

/**
 * Serialize Tiptap content to template string
 */
function serializeTiptapToTemplate(editorInstance: Editor): string {
  const json = editorInstance.getJSON()
  const parts: string[] = []

  function processNode(node: {
    type?: string
    text?: string
    attrs?: Record<string, unknown>
    content?: Array<Record<string, unknown>>
  }) {
    if (node.type === 'text') {
      parts.push(node.text || '')
    } else if (node.type === 'variableMention') {
      const { path, func } = node.attrs || {}
      if (func) {
        parts.push(`{{${func} ${path}}}`)
      } else {
        parts.push(`{{${path}}}`)
      }
    } else if (node.content) {
      node.content.forEach(processNode)
    }
  }

  json.content?.forEach(processNode as (node: Record<string, unknown>) => void)
  return parts.join('')
}

// Initialize editor
const editor = useEditor({
  extensions: [
    StarterKit.configure({
      // Disable features we don't need
      heading: false,
      bulletList: false,
      orderedList: false,
      listItem: false,
      blockquote: false,
      codeBlock: false,
      horizontalRule: false,
      hardBreak: false,
      code: false,
      bold: false,
      italic: false,
      strike: false,
    }),
    VariableMention.configure({
      suggestion: variableSuggestion({ mediaType: props.mediaType }),
    }),
  ],
  content: parseTemplateToTiptap(props.modelValue),
  editable: !props.disabled,
  editorProps: {
    attributes: {
      class: cn(
        'template-editor min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background',
        'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
        'disabled:cursor-not-allowed disabled:opacity-50',
        props.class,
      ),
      'data-placeholder': props.placeholder,
    },
  },
  onUpdate: ({ editor: ed }) => {
    const serialized = serializeTiptapToTemplate(ed as unknown as Editor)
    emit('update:modelValue', serialized)
  },
})

// Watch for external changes to modelValue
watch(
  () => props.modelValue,
  (newValue) => {
    if (!editor.value) return

    const currentValue = serializeTiptapToTemplate(editor.value)
    if (currentValue !== newValue) {
      const content = parseTemplateToTiptap(newValue)
      editor.value.commands.setContent(content)
    }
  },
)

// Watch for disabled changes
watch(
  () => props.disabled,
  (disabled) => {
    if (editor.value) {
      editor.value.setEditable(!disabled)
    }
  },
)

// Handle context menu on variable mentions
function handleNodeContextMenu(event: MouseEvent, pos: number) {
  event.preventDefault()
  selectedNodePos.value = pos
  contextMenuPosition.value = { x: event.clientX, y: event.clientY }
  showContextMenu.value = true
}

// Set up context menu listener
if (editor.value) {
  editor.value.view.dom.addEventListener('contextmenu', (event) => {
    const target = event.target as HTMLElement
    if (target.classList.contains('variable-mention')) {
      // Find the node position
      const pos = editor.value?.view.posAtDOM(target, 0)
      if (pos !== undefined) {
        handleNodeContextMenu(event, pos)
      }
    }
  })
}

// Get current selected node
const selectedNode = ref<{
  type?: { name?: string }
  attrs?: Record<string, unknown>
  nodeSize?: number
} | null>(null)
watch([showContextMenu, selectedNodePos], () => {
  if (showContextMenu.value && selectedNodePos.value !== null && editor.value) {
    const node = editor.value.state.doc.nodeAt(selectedNodePos.value)
    selectedNode.value = node
  } else {
    selectedNode.value = null
  }
})

// Context menu handlers
function handleWrapWithFunction(funcName: 'clean' | 'sanitize') {
  if (selectedNodePos.value !== null && editor.value) {
    editor.value.commands.updateVariableMentionFunc(selectedNodePos.value, funcName)
  }
  showContextMenu.value = false
  selectedNodePos.value = null
}

function handleRemoveFunction() {
  if (selectedNodePos.value !== null && editor.value) {
    editor.value.commands.updateVariableMentionFunc(selectedNodePos.value, null)
  }
  showContextMenu.value = false
  selectedNodePos.value = null
}

function handleDeleteToken() {
  if (selectedNodePos.value !== null && editor.value) {
    const node = editor.value.state.doc.nodeAt(selectedNodePos.value)
    if (node) {
      editor.value
        .chain()
        .focus()
        .deleteRange({
          from: selectedNodePos.value,
          to: selectedNodePos.value + node.nodeSize,
        })
        .run()
    }
  }
  showContextMenu.value = false
  selectedNodePos.value = null
}

// Convert selected node to Token format for context menu
const selectedToken = ref<Token | null>(null)
watch(selectedNode, (node) => {
  if (node?.type?.name === 'variableMention' && node.attrs) {
    selectedToken.value = {
      type: 'variable',
      value: (node.attrs.path as string) || '',
      func: (node.attrs.func as 'clean' | 'sanitize' | undefined) || undefined,
    }
  } else {
    selectedToken.value = null
  }
})

// Cleanup
onBeforeUnmount(() => {
  editor.value?.destroy()
})
</script>

<template>
  <div class="template-editor-wrapper relative">
    <EditorContent :editor="editor" />

    <!-- Context Menu -->
    <TokenContextMenu
      v-if="showContextMenu && selectedToken"
      :position="contextMenuPosition"
      :token="selectedToken"
      @wrap="handleWrapWithFunction"
      @remove-function="handleRemoveFunction"
      @delete="handleDeleteToken"
      @close="showContextMenu = false"
    />
  </div>
</template>

<style>
/* ProseMirror base styles */
.ProseMirror {
  white-space: pre-wrap;
  word-wrap: break-word;
  overflow-wrap: break-word;
  line-height: 1.6;
}

.ProseMirror p {
  margin: 0;
}

/* Placeholder */
.ProseMirror p.is-editor-empty:first-child::before {
  content: attr(data-placeholder);
  color: var(--muted-foreground);
  pointer-events: none;
  float: left;
  height: 0;
}

/* Focus styles */
.ProseMirror:focus {
  outline: none;
}

/* Variable mention spacing */
.ProseMirror .variable-mention {
  user-select: none;
  vertical-align: middle;
}
</style>
