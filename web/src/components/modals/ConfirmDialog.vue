<script setup lang="ts">
import { computed, inject } from 'vue'
import { Button } from '@/components/ui/button'
import BaseDialog from './BaseDialog.vue'

interface Props {
  title?: string
  message: string
  confirmLabel?: string
  cancelLabel?: string
  severity?: 'danger' | 'warning' | 'info' | 'success' | 'secondary'
}

const props = withDefaults(defineProps<Props>(), {
  title: 'Confirm',
  confirmLabel: 'Confirm',
  cancelLabel: 'Cancel',
  severity: 'danger',
})

const dialogRef = inject('dialogRef') as { value: { close: (data?: unknown) => void } }

const handleConfirm = () => {
  dialogRef.value.close({ confirmed: true })
}

const handleCancel = () => {
  dialogRef.value.close({ confirmed: false })
}

// Map severity to button variant
const confirmButtonVariant = computed(() => {
  switch (props.severity) {
    case 'danger':
      return 'destructive'
    case 'warning':
      return 'default'
    case 'success':
      return 'default'
    case 'info':
      return 'default'
    case 'secondary':
      return 'secondary'
    default:
      return 'destructive'
  }
})
</script>

<template>
  <BaseDialog :title="title">
    <p class="text-base">{{ message }}</p>
    <template #footer>
      <Button variant="outline" @click="handleCancel">
        {{ cancelLabel }}
      </Button>
      <Button :variant="confirmButtonVariant" @click="handleConfirm">
        {{ confirmLabel }}
      </Button>
    </template>
  </BaseDialog>
</template>
