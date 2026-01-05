import { mergeAttributes, Node } from '@tiptap/core'
import { VueNodeViewRenderer } from '@tiptap/vue-3'
import ConditionalBlockView from './ConditionalBlockView.vue'

export interface ConditionalBlockOptions {
  HTMLAttributes: Record<string, any>
}

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    conditionalBlock: {
      /**
       * Insert a conditional block
       */
      insertConditionalBlock: (attributes: { field: string }) => ReturnType
      /**
       * Update the field of a conditional block
       */
      updateConditionalField: (pos: number, field: string) => ReturnType
      /**
       * Delete a conditional block
       */
      deleteConditionalBlock: (pos: number) => ReturnType
    }
  }
}

export const ConditionalBlock = Node.create<ConditionalBlockOptions>({
  name: 'conditionalBlock',

  addOptions() {
    return {
      HTMLAttributes: {},
    }
  },

  group: 'inline',

  inline: true,

  content: 'inline*',

  addAttributes() {
    return {
      field: {
        default: '',
        parseHTML: (element) => element.getAttribute('data-field'),
        renderHTML: (attributes) => ({
          'data-field': attributes.field,
        }),
      },
    }
  },

  parseHTML() {
    return [
      {
        tag: 'span[data-type="conditional-block"]',
      },
    ]
  },

  renderHTML({ node, HTMLAttributes }) {
    return [
      'span',
      mergeAttributes(
        { 'data-type': 'conditional-block' },
        this.options.HTMLAttributes,
        HTMLAttributes,
      ),
      ['span', { class: 'conditional-indicator' }, `if ${node.attrs.field}`],
      ['span', { class: 'conditional-content' }, 0],
    ]
  },

  addNodeView() {
    return VueNodeViewRenderer(ConditionalBlockView)
  },

  addCommands() {
    return {
      insertConditionalBlock:
        (attributes) =>
        ({ commands }) => {
          return commands.insertContent({
            type: this.name,
            attrs: attributes,
            // Don't specify content - let it be empty, user can type into it
          })
        },
      updateConditionalField:
        (pos, field) =>
        ({ tr }) => {
          const node = tr.doc.nodeAt(pos)
          if (node && node.type.name === this.name) {
            tr.setNodeMarkup(pos, undefined, {
              ...node.attrs,
              field,
            })
            return true
          }
          return false
        },
      deleteConditionalBlock:
        (pos) =>
        ({ tr, dispatch }) => {
          const node = tr.doc.nodeAt(pos)
          if (node && node.type.name === this.name) {
            if (dispatch) {
              tr.delete(pos, pos + node.nodeSize)
            }
            return true
          }
          return false
        },
    }
  },
})

