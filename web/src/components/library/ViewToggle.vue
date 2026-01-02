<template>
  <div class="flex items-center border rounded-lg p-1 bg-muted/50">
    <Button
      variant="ghost"
      size="sm"
      :class="[
        'h-8 w-8 p-0',
        modelValue === 'grid' && 'bg-background shadow-sm',
      ]"
      @click="$emit('update:modelValue', 'grid')"
    >
      <LayoutGrid class="h-4 w-4" />
      <span class="sr-only">Grid view</span>
    </Button>
    <Button
      variant="ghost"
      size="sm"
      :class="[
        'h-8 w-8 p-0',
        modelValue === 'table' && 'bg-background shadow-sm',
      ]"
      @click="$emit('update:modelValue', 'table')"
    >
      <Table class="h-4 w-4" />
      <span class="sr-only">Table view</span>
    </Button>
  </div>
</template>

<script setup lang="ts">
import { LayoutGrid, Table } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'

export type ViewMode = 'grid' | 'table'

defineProps<{
  modelValue: ViewMode
}>()

defineEmits<{
  (e: 'update:modelValue', value: ViewMode): void
}>()
</script>

