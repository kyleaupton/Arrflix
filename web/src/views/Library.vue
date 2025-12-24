<template>
  <div class="flex flex-col gap-6">
    <div>
      <h1 class="text-2xl font-semibold">Library</h1>
      <p class="text-sm text-muted-foreground">Browse and manage your media</p>
    </div>
    <div class="space-y-3">
      <!-- <LibraryItem v-for="item in query.data || []" :key="item.id" :item="item" /> -->
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { getV1LibraryOptions } from '@/client/@tanstack/vue-query.gen'
import LibraryItem from '@/components/library/LibraryItem.vue'

const query = useQuery(getV1LibraryOptions())

onMounted(async () => {
  const { data } = await query.refetch()
  console.log(data)
})
</script>

<style scoped></style>
