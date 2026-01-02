<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { cn } from '@/lib/utils'
import VariableAutocomplete from './VariableAutocomplete.vue'
import TokenContextMenu from './TokenContextMenu.vue'
import type { TemplateVariable } from '@/composables/useTemplateVariables'

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
  placeholder: 'Type {{ to insert a variable...',
  mediaType: 'movie',
  disabled: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

// Refs
const editorRef = ref<HTMLDivElement | null>(null)
const tokens = ref<Token[]>([])
const isEditing = ref(false)

// Autocomplete state
const showAutocomplete = ref(false)
const autocompletePosition = ref({ top: 0, left: 0 })
const autocompleteQuery = ref('')
const pendingAutocompleteRange = ref<Range | null>(null)

// Context menu state
const showContextMenu = ref(false)
const contextMenuPosition = ref({ x: 0, y: 0 })
const selectedTokenIndex = ref<number | null>(null)

// Template parsing regex
// Matches: {{.Var}}, {{func .Var}}
const TEMPLATE_REGEX = /\{\{(clean|sanitize)?\s*([.\w]+)\}\}/g

/**
 * Parse a template string into tokens
 */
function parseTemplate(template: string): Token[] {
  const result: Token[] = []
  let lastIndex = 0

  // Reset regex
  TEMPLATE_REGEX.lastIndex = 0

  let match
  while ((match = TEMPLATE_REGEX.exec(template)) !== null) {
    // Add text before the match
    if (match.index > lastIndex) {
      result.push({
        type: 'text',
        value: template.slice(lastIndex, match.index),
      })
    }

    // Add the variable token
    result.push({
      type: 'variable',
      value: match[2], // The variable path
      func: match[1] as 'clean' | 'sanitize' | undefined,
    })

    lastIndex = match.index + match[0].length
  }

  // Add remaining text
  if (lastIndex < template.length) {
    result.push({
      type: 'text',
      value: template.slice(lastIndex),
    })
  }

  return result
}

/**
 * Serialize tokens back to a template string
 */
function serializeTokens(tokenList: Token[]): string {
  return tokenList
    .map((token) => {
      if (token.type === 'text') {
        return token.value
      }
      if (token.func) {
        return `{{${token.func} ${token.value}}}`
      }
      return `{{${token.value}}}`
    })
    .join('')
}

/**
 * Get display label for a variable
 */
function getVariableLabel(token: Token): string {
  if (token.func) {
    return `${token.func} ${token.value}`
  }
  return token.value
}

// Initialize tokens from modelValue
watch(
  () => props.modelValue,
  (newValue) => {
    if (!isEditing.value) {
      tokens.value = parseTemplate(newValue || '')
    }
  },
  { immediate: true },
)

// Update modelValue when tokens change
function emitUpdate() {
  const serialized = serializeTokens(tokens.value)
  emit('update:modelValue', serialized)
}

/**
 * Handle input in the contenteditable
 */
function handleInput(event: Event) {
  isEditing.value = true
  const target = event.target as HTMLDivElement

  // Get the current text content and cursor position
  const selection = window.getSelection()
  if (!selection || selection.rangeCount === 0) return

  const range = selection.getRangeAt(0)

  // Check if we should trigger autocomplete
  const textBeforeCursor = getTextBeforeCursor(range)
  const autocompleteMatch = textBeforeCursor.match(/\{\{([^}]*)$/)

  if (autocompleteMatch) {
    // Show autocomplete
    autocompleteQuery.value = autocompleteMatch[1].trim()
    showAutocomplete.value = true
    pendingAutocompleteRange.value = range.cloneRange()

    // Position autocomplete near cursor
    const rect = range.getBoundingClientRect()
    const editorRect = editorRef.value?.getBoundingClientRect()
    if (editorRect) {
      autocompletePosition.value = {
        top: rect.bottom - editorRect.top + 4,
        left: rect.left - editorRect.left,
      }
    }
  } else {
    showAutocomplete.value = false
    pendingAutocompleteRange.value = null
  }

  // Parse the current content to update tokens
  syncFromDOM()
}

/**
 * Get the text content before the cursor
 */
function getTextBeforeCursor(range: Range): string {
  const container = range.startContainer
  if (container.nodeType === Node.TEXT_NODE) {
    return container.textContent?.slice(0, range.startOffset) || ''
  }
  return ''
}

/**
 * Sync tokens from the DOM state
 */
function syncFromDOM() {
  if (!editorRef.value) return

  const newTokens: Token[] = []

  editorRef.value.childNodes.forEach((node) => {
    if (node.nodeType === Node.TEXT_NODE) {
      const text = node.textContent || ''
      if (text) {
        // Parse any template syntax in the text
        const parsed = parseTemplate(text)
        newTokens.push(...parsed)
      }
    } else if (node.nodeType === Node.ELEMENT_NODE) {
      const el = node as HTMLElement
      if (el.classList.contains('template-token')) {
        const variable = el.dataset.variable || ''
        const func = el.dataset.func as 'clean' | 'sanitize' | undefined
        newTokens.push({
          type: 'variable',
          value: variable,
          func,
        })
      } else {
        // Handle other elements (like pasted content)
        const text = el.textContent || ''
        if (text) {
          const parsed = parseTemplate(text)
          newTokens.push(...parsed)
        }
      }
    }
  })

  tokens.value = newTokens
  emitUpdate()

  nextTick(() => {
    isEditing.value = false
  })
}

/**
 * Handle variable selection from autocomplete
 */
function handleVariableSelect(variable: TemplateVariable) {
  showAutocomplete.value = false

  if (!editorRef.value) return

  // Get the current content and find the {{ trigger
  const currentContent = serializeTokens(tokens.value)

  // Find the last {{ in the content (this is the trigger we want to replace)
  const triggerPos = currentContent.lastIndexOf('{{')

  if (triggerPos !== -1) {
    // Build new content: everything before {{ + the new token
    const beforeTrigger = currentContent.slice(0, triggerPos)
    const newContent = beforeTrigger + `{{${variable.path}}}`

    // Parse and update tokens
    tokens.value = parseTemplate(newContent)
    emitUpdate()

    // Re-render and focus
    nextTick(() => {
      renderTokens()
      editorRef.value?.focus()

      // Move cursor to end
      const selection = window.getSelection()
      if (selection && editorRef.value) {
        const range = document.createRange()
        range.selectNodeContents(editorRef.value)
        range.collapse(false)
        selection.removeAllRanges()
        selection.addRange(range)
      }
    })
  }

  pendingAutocompleteRange.value = null
}

/**
 * Create a token DOM element
 */
function createTokenElement(token: Token): HTMLSpanElement {
  const span = document.createElement('span')
  span.className = 'template-token'
  span.contentEditable = 'false'
  span.dataset.variable = token.value
  if (token.func) {
    span.dataset.func = token.func
  }

  const label = document.createElement('span')
  label.className = 'token-label'
  label.textContent = getVariableLabel(token)
  span.appendChild(label)

  return span
}

/**
 * Handle right-click on a token
 */
function handleTokenContextMenu(event: MouseEvent, index: number) {
  event.preventDefault()
  selectedTokenIndex.value = index
  contextMenuPosition.value = { x: event.clientX, y: event.clientY }
  showContextMenu.value = true
}

/**
 * Handle function wrap from context menu
 */
function handleWrapWithFunction(funcName: 'clean' | 'sanitize') {
  if (selectedTokenIndex.value === null) return

  const token = tokens.value[selectedTokenIndex.value]
  if (token && token.type === 'variable') {
    token.func = funcName
    emitUpdate()
    renderTokens()
  }

  showContextMenu.value = false
  selectedTokenIndex.value = null
}

/**
 * Handle remove function from context menu
 */
function handleRemoveFunction() {
  if (selectedTokenIndex.value === null) return

  const token = tokens.value[selectedTokenIndex.value]
  if (token && token.type === 'variable') {
    token.func = undefined
    emitUpdate()
    renderTokens()
  }

  showContextMenu.value = false
  selectedTokenIndex.value = null
}

/**
 * Handle delete token from context menu
 */
function handleDeleteToken() {
  if (selectedTokenIndex.value === null) return

  tokens.value.splice(selectedTokenIndex.value, 1)
  emitUpdate()
  renderTokens()

  showContextMenu.value = false
  selectedTokenIndex.value = null
}

/**
 * Render tokens to the DOM
 */
function renderTokens() {
  if (!editorRef.value) return

  // Save cursor position
  const selection = window.getSelection()
  let savedOffset = 0
  let savedInToken = false

  if (selection && selection.rangeCount > 0) {
    const range = selection.getRangeAt(0)
    // Try to calculate offset
    const walker = document.createTreeWalker(editorRef.value, NodeFilter.SHOW_ALL)
    let node: Node | null
    let offset = 0
    while ((node = walker.nextNode())) {
      if (node === range.startContainer) {
        savedOffset = offset + range.startOffset
        break
      }
      if (node.nodeType === Node.TEXT_NODE) {
        offset += node.textContent?.length || 0
      } else if ((node as HTMLElement).classList?.contains('template-token')) {
        savedInToken = true
        offset += 1 // Count token as 1 character for positioning
      }
    }
  }

  // Clear and rebuild
  editorRef.value.innerHTML = ''

  let tokenIndex = 0
  tokens.value.forEach((token, idx) => {
    if (token.type === 'text') {
      const textNode = document.createTextNode(token.value)
      editorRef.value!.appendChild(textNode)
    } else {
      const tokenEl = createTokenElement(token)
      tokenEl.addEventListener('contextmenu', (e) => handleTokenContextMenu(e, idx))
      editorRef.value!.appendChild(tokenEl)
      tokenIndex++
    }
  })

  // Restore cursor position
  if (selection && !savedInToken) {
    try {
      const range = document.createRange()
      const walker = document.createTreeWalker(editorRef.value, NodeFilter.SHOW_TEXT)
      let node: Node | null
      let currentOffset = 0

      while ((node = walker.nextNode())) {
        const nodeLength = node.textContent?.length || 0
        if (currentOffset + nodeLength >= savedOffset) {
          range.setStart(node, savedOffset - currentOffset)
          range.collapse(true)
          selection.removeAllRanges()
          selection.addRange(range)
          break
        }
        currentOffset += nodeLength
      }
    } catch {
      // Cursor restoration failed, that's okay
    }
  }
}

/**
 * Handle keydown events
 */
function handleKeyDown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    showAutocomplete.value = false
  }

  // Handle backspace on token
  if (event.key === 'Backspace') {
    const selection = window.getSelection()
    if (selection && selection.rangeCount > 0) {
      const range = selection.getRangeAt(0)
      if (range.collapsed) {
        // Check if we're right after a token
        const prevSibling = range.startContainer.previousSibling
        if (prevSibling && (prevSibling as HTMLElement).classList?.contains('template-token')) {
          event.preventDefault()
          prevSibling.remove()
          syncFromDOM()
        }
      }
    }
  }
}

/**
 * Handle paste events
 */
function handlePaste(event: ClipboardEvent) {
  event.preventDefault()
  const text = event.clipboardData?.getData('text/plain') || ''

  // Insert as plain text
  const selection = window.getSelection()
  if (selection && selection.rangeCount > 0) {
    const range = selection.getRangeAt(0)
    range.deleteContents()

    // Parse the pasted text for any template syntax
    const parsed = parseTemplate(text)

    parsed.forEach((token) => {
      if (token.type === 'text') {
        const textNode = document.createTextNode(token.value)
        range.insertNode(textNode)
        range.setStartAfter(textNode)
      } else {
        const tokenEl = createTokenElement(token)
        range.insertNode(tokenEl)
        range.setStartAfter(tokenEl)
      }
    })

    range.collapse(true)
    selection.removeAllRanges()
    selection.addRange(range)
  }

  syncFromDOM()
}

/**
 * Handle focus
 */
function handleFocus() {
  isEditing.value = true
}

/**
 * Handle blur
 */
function handleBlur() {
  // Delay to allow autocomplete clicks
  setTimeout(() => {
    if (!showAutocomplete.value) {
      isEditing.value = false
    }
  }, 150)
}

/**
 * Close context menu when clicking outside
 */
function handleDocumentClick(event: MouseEvent) {
  if (showContextMenu.value) {
    showContextMenu.value = false
  }
}

// Computed
const isEmpty = computed(() => {
  return (
    tokens.value.length === 0 ||
    (tokens.value.length === 1 && tokens.value[0].type === 'text' && !tokens.value[0].value)
  )
})

const selectedToken = computed(() => {
  if (selectedTokenIndex.value === null) return null
  return tokens.value[selectedTokenIndex.value]
})

// Lifecycle
onMounted(() => {
  renderTokens()
  document.addEventListener('click', handleDocumentClick)
})

onUnmounted(() => {
  document.removeEventListener('click', handleDocumentClick)
})

// Watch for external changes
watch(
  () => props.modelValue,
  () => {
    if (!isEditing.value) {
      renderTokens()
    }
  },
)
</script>

<template>
  <div class="template-editor-wrapper relative">
    <div
      ref="editorRef"
      :class="
        cn(
          'template-editor min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background',
          'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
          'disabled:cursor-not-allowed disabled:opacity-50',
          isEmpty && 'is-empty',
          props.class,
        )
      "
      contenteditable="true"
      :data-placeholder="placeholder"
      @input="handleInput"
      @keydown="handleKeyDown"
      @paste="handlePaste"
      @focus="handleFocus"
      @blur="handleBlur"
    />

    <!-- Autocomplete Popover -->
    <VariableAutocomplete
      v-if="showAutocomplete"
      :query="autocompleteQuery"
      :position="autocompletePosition"
      :media-type="mediaType"
      @select="handleVariableSelect"
      @close="showAutocomplete = false"
    />

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
.template-token {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.125rem 0.5rem;
  margin: 0 0.125rem;
  background-color: var(--secondary);
  border: 1px solid var(--border);
  border-radius: 0.375rem;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.8125rem;
  color: var(--secondary-foreground);
  cursor: default;
  user-select: none;
  vertical-align: middle;
  white-space: nowrap;
  transition: background-color 0.15s ease;
}

.template-token:hover {
  background-color: var(--accent);
}

.template-token .token-label {
  white-space: nowrap;
}

/* Tokens with a function wrapper (clean, sanitize) get a distinct style */
.template-token[data-func] {
  background-color: color-mix(in oklch, var(--chart-2) 20%, transparent);
  border-color: color-mix(in oklch, var(--chart-2) 40%, transparent);
  color: var(--foreground);
}

.template-token[data-func]:hover {
  background-color: color-mix(in oklch, var(--chart-2) 30%, transparent);
}
</style>

<style scoped>
.template-editor {
  white-space: pre-wrap;
  word-wrap: break-word;
  overflow-wrap: break-word;
  line-height: 1.6;
}

.template-editor:empty::before,
.template-editor.is-empty::before {
  content: attr(data-placeholder);
  color: var(--muted-foreground);
  pointer-events: none;
}

.template-editor:focus:empty::before,
.template-editor:focus.is-empty::before {
  content: '';
}
</style>
