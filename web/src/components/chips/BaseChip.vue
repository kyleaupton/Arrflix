<template>
  <div>
    <Chip>{{ label }}</Chip>
  </div>
</template>

<script setup lang="ts">
import Chip from 'primevue/chip'

defineProps<{
  label: string
}>()
</script>

<style scoped></style>
