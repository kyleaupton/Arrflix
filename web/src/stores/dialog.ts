import { markRaw, ref } from 'vue'
import type { Component } from 'vue'

export interface DialogInstance {
  id: string
  component: Component | (() => Promise<Component>)
  props?: Record<string, unknown>
  data?: unknown
  onClose?: (result?: { data?: unknown }) => void
}

const dialogs = ref<DialogInstance[]>([])

let nextId = 0

export function useDialogStore() {
  const openDialog = (
    component: Component,
    options?: {
      data?: unknown
      props?: Record<string, unknown>
      onClose?: (result?: { data?: unknown }) => void
    },
  ): DialogInstance => {
    const id = `dialog-${nextId++}`
    const instance: DialogInstance = {
      id,
      // Must mark as raw or we get a warning
      // about performance issues in the console
      component: markRaw(component),
      props: options?.props,
      data: options?.data,
      onClose: options?.onClose,
    }
    dialogs.value.push(instance)
    return instance
  }

  const closeDialog = (id: string, data?: unknown) => {
    const index = dialogs.value.findIndex((d) => d.id === id)
    if (index === -1) return

    const instance = dialogs.value[index]
    dialogs.value.splice(index, 1)

    // Call onClose callback if provided
    if (instance?.onClose) {
      instance.onClose({ data })
    }
  }

  const closeAll = () => {
    dialogs.value = []
  }

  return {
    dialogs,
    openDialog,
    closeDialog,
    closeAll,
  }
}
