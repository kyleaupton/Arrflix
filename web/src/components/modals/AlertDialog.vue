<script setup lang="ts">
import { inject } from 'vue'
import Button from 'primevue/button'

interface Props {
  message: string
  severity?: 'info' | 'warning' | 'error' | 'success'
  okLabel?: string
}

const props = withDefaults(defineProps<Props>(), {
  severity: 'info',
  okLabel: 'OK',
})

const dialogRef = inject('dialogRef') as { value: { close: (data?: unknown) => void } }

const handleOk = () => {
  dialogRef.value.close()
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <p class="text-base">{{ message }}</p>
    <div class="flex justify-end pt-2">
      <Button
        :label="okLabel"
        :severity="severity === 'error' ? 'danger' : severity"
        @click="handleOk"
      />
    </div>
  </div>
</template>

