<script setup lang="ts">
import { computed, markRaw, provide, ref } from 'vue'
import { Dialog, DialogContent } from '@/components/ui/dialog'
import { cn } from '@/lib/utils'
import type { Component } from 'vue'

interface Props {
  instance: {
    id: string
    component: Component | (() => Promise<Component>)
    props?: Record<string, unknown>
    data?: unknown
    onClose?: (result?: { data?: unknown }) => void
  }
}

const props = defineProps<Props>()

const isOpen = ref(true)
const component = markRaw(props.instance.component)

// Create dialogRef object that components can inject
// Components expect: dialogRef.value.close() and dialogRef.value.data
const dialogRef = {
  value: {
    close: (data?: unknown) => {
      isOpen.value = false
      // Small delay to allow animation to complete
      setTimeout(() => {
        props.instance.onClose?.({ data })
      }, 200)
    },
    data: props.instance.data,
  },
}

// Provide dialogRef for child components to inject
provide('dialogRef', dialogRef)

// Handle dialog close
const handleOpenChange = (open: boolean) => {
  if (!open) {
    isOpen.value = false
    setTimeout(() => {
      props.instance.onClose?.()
    }, 200)
  }
}

// Default dialog constraints
// These can be overridden via instance.props.class
const defaultDialogClasses = computed(() => {
  return cn(
    // Max width: responsive, defaults to lg on small screens, xl on larger screens
    'max-w-[calc(100vw-2rem)] sm:max-w-lg lg:max-w-2xl',
    // Max height: prevent overflow, allow scrolling
    'max-h-[calc(100vh-2rem)]',
    // Ensure content can scroll if needed
    'flex flex-col overflow-hidden',
  )
})

// Merge default classes with custom classes from props
const dialogClasses = computed(() => {
  const customClass = props.instance.props?.class as string | undefined
  return cn(defaultDialogClasses.value, customClass)
})
</script>

<template>
  <Dialog :open="isOpen" @update:open="handleOpenChange">
    <DialogContent
      :class="dialogClasses"
      :style="instance.props?.style"
      :show-close-button="(instance.props?.showCloseButton as boolean) ?? true"
    >
      <component :is="component" v-bind="props.instance.props" />
    </DialogContent>
  </Dialog>
</template>
