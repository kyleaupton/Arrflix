<script setup lang="ts">
import { markRaw, provide, ref } from 'vue'
import { Dialog, DialogContent } from '@/components/ui/dialog'
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

// Lazy load the component
// const LazyComponent = computed(() => {
//   const comp = props.instance.component
//   if (typeof comp === 'function') {
//     return defineAsyncComponent(comp)
//   }
//   return comp
// })

// Handle dialog close
const handleOpenChange = (open: boolean) => {
  if (!open) {
    isOpen.value = false
    setTimeout(() => {
      props.instance.onClose?.()
    }, 200)
  }
}
</script>

<template>
  <Dialog :open="isOpen" @update:open="handleOpenChange">
    <DialogContent
      :class="instance.props?.class"
      :style="instance.props?.style"
      :show-close-button="(instance.props?.showCloseButton as boolean) ?? true"
    >
      <component :is="component" v-bind="props.instance.props" />
    </DialogContent>
  </Dialog>
</template>
