<script setup lang="ts">
import { DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog'

interface Props {
  title?: string
  description?: string
}

withDefaults(defineProps<Props>(), {
  title: undefined,
  description: undefined,
})
</script>

<template>
  <div class="flex flex-col gap-4 min-h-0 flex-1">
    <DialogHeader v-if="title || $slots.header" class="shrink-0">
      <DialogTitle>
        <slot name="title">{{ title }}</slot>
      </DialogTitle>
      <DialogDescription v-if="description || $slots.description">
        <slot name="description">{{ description }}</slot>
      </DialogDescription>
    </DialogHeader>

    <div class="flex-1 min-h-0 overflow-auto -mx-1 px-1">
      <slot />
    </div>

    <DialogFooter v-if="$slots.footer" class="shrink-0">
      <slot name="footer" />
    </DialogFooter>
  </div>
</template>
