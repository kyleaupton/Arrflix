import { ref } from 'vue'
import { useDialog } from 'primevue/usedialog'
import ConfirmDialog from '@/components/modals/ConfirmDialog.vue'
import AlertDialog from '@/components/modals/AlertDialog.vue'

export interface ConfirmOptions {
  title?: string
  message: string
  confirmLabel?: string
  cancelLabel?: string
  severity?: 'danger' | 'warning' | 'info' | 'success' | 'secondary'
}

export interface AlertOptions {
  title?: string
  message: string
  severity?: 'info' | 'warning' | 'error' | 'success'
  okLabel?: string
}

export interface PromptOptions {
  title?: string
  message: string
  placeholder?: string
  defaultValue?: string
  confirmLabel?: string
  cancelLabel?: string
}

/**
 * Composable for managing modals using PrimeVue DynamicDialog
 */
export function useModal() {
  const dialog = useDialog()
  const currentDialogInstance = ref<ReturnType<typeof dialog.open> | null>(null)

  /**
   * Show a confirmation dialog
   * @returns Promise that resolves to true if confirmed, false if cancelled
   */
  const confirm = (options: ConfirmOptions): Promise<boolean> => {
    return new Promise((resolve) => {
      const instance = dialog.open(ConfirmDialog, {
        props: {
          header: options.title || 'Confirm',
          modal: true,
          style: { width: '450px' },
          ...options,
        },
        onClose: (result) => {
          currentDialogInstance.value = null
          resolve((result?.data as { confirmed?: boolean })?.confirmed === true)
        },
      })
      currentDialogInstance.value = instance
    })
  }

  /**
   * Show an alert dialog
   * @returns Promise that resolves when the dialog is closed
   */
  const alert = (options: AlertOptions): Promise<void> => {
    return new Promise((resolve) => {
      const instance = dialog.open(AlertDialog, {
        props: {
          header: options.title || 'Alert',
          modal: true,
          style: { width: '450px' },
          ...options,
        },
        onClose: () => {
          currentDialogInstance.value = null
          resolve()
        },
      })
      currentDialogInstance.value = instance
    })
  }

  /**
   * Show a prompt dialog (input dialog)
   * Note: This requires a PromptDialog component to be created
   * @returns Promise that resolves to the input value or null if cancelled
   */
  const prompt = (options: PromptOptions): Promise<string | null> => {
    // For now, we'll use a simple implementation with window.prompt
    // A proper PromptDialog component can be added later
    return new Promise((resolve) => {
      const result = window.prompt(options.message, options.defaultValue || '')
      resolve(result)
    })
  }

  /**
   * Open a custom component as a modal
   * @returns DynamicDialogInstance with a close method
   */
  const open = <T = unknown>(
    component: unknown,
    options?: {
      data?: unknown
      props?: Record<string, unknown>
      onClose?: (result?: { data?: T }) => void
    },
  ) => {
    const instance = dialog.open(component, {
      data: options?.data,
      props: {
        modal: true,
        ...options?.props,
      },
      onClose: (result?: { data?: T }) => {
        currentDialogInstance.value = null
        options?.onClose?.(result)
      },
    })
    currentDialogInstance.value = instance
    return instance
  }

  /**
   * Close the currently open dialog
   * Note: This closes the last dialog opened via this composable instance.
   * For more control, use the instance returned from open() and call close() on it directly.
   * @param data Optional data to pass back to the onClose callback
   */
  const close = (data?: unknown) => {
    if (currentDialogInstance.value) {
      currentDialogInstance.value.close(data)
      currentDialogInstance.value = null
    }
  }

  return {
    confirm,
    alert,
    prompt,
    open,
    close,
  }
}
