import { Extension } from '@tiptap/core'

export interface SlashCommandOptions {
  onShowCommandList: (position: { top: number; left: number }) => void
}

export const SlashCommand = Extension.create<SlashCommandOptions>({
  name: 'slashCommand',

  addOptions() {
    return {
      onShowCommandList: () => {},
    }
  },

  addKeyboardShortcuts() {
    return {
      '/': ({ editor }) => {
        // Get cursor position for popup
        const { view } = editor
        const { from } = view.state.selection
        const coords = view.coordsAtPos(from)
        
        // Show command list at cursor position (after a small delay to let / be typed)
        setTimeout(() => {
          this.options.onShowCommandList({
            top: coords.bottom + 4,
            left: coords.left,
          })
        }, 10)
        
        // Let the / character be typed
        return false
      },
    }
  },
})

