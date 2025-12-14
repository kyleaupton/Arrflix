<script setup lang="ts">
import { inject } from 'vue'
import Button from 'primevue/button'

interface Props {
  message: string
  confirmLabel?: string
  cancelLabel?: string
  severity?: 'danger' | 'warning' | 'info' | 'success' | 'secondary'
}

withDefaults(defineProps<Props>(), {
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
</script>

<template>
  <div class="flex flex-col gap-4">
    <p class="text-base">{{ message }}</p>
    <div class="flex justify-end gap-2 pt-2">
      <Button :label="cancelLabel" severity="secondary" variant="outlined" @click="handleCancel" />
      <Button :label="confirmLabel" :severity="severity" @click="handleConfirm" />
    </div>
  </div>
</template>
