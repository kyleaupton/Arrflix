<template>
  <Card>
    <template #title>
      {{ item.title }}
    </template>
    <template #content>
      <div>content here</div>
    </template>
  </Card>
</template>

<script setup lang="ts">
import Card from 'primevue/card'

defineProps<{
  item: any // eslint-disable-line @typescript-eslint/no-explicit-any
}>()
</script>

<style scoped></style>
