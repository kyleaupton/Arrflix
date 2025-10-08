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
import type { HandlersMediaItemSwagger } from '@/client/types.gen'

const props = defineProps<{
  item: HandlersMediaItemSwagger
}>()
</script>

<style scoped></style>
