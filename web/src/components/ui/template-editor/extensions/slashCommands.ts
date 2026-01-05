import { Editor, Extension } from '@tiptap/core'
import type { SuggestionOptions } from '@tiptap/suggestion'

export interface CommandItem {
  title: string
  description: string
  icon: string
  command?: (editor: Editor) => void
  onSelect?: (editor: Editor) => void
}

// Export a function that returns suggestion options
// This will be used directly with the Suggestion plugin in VariableMention
export function createSlashCommandsSuggestion(
  onConditionalSelect: () => void
): Omit<SuggestionOptions, 'editor'> {
  // Instead of creating a separate extension, we'll return null
  // and handle slash commands through a keyboard shortcut
  return null as any
}

